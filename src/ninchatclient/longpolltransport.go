package main

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	minPollTimeout = Second * 64
	minSendTimeout = Second * 7

	maxJSONPSize = 2048
)

// LongPollTransport creates or resumes a session, and runs an I/O loop after
// that.
func LongPollTransport(s *Session, host string) (connWorked, gotOnline bool) {
	defer func() {
		if err := jsError(recover()); err != nil {
			s.log("poll:", err)
		}
	}()

	// assume that JSONP works always (endpoint discovery worked, after all)
	connWorked = true

	url := "https://" + host + pollPath

	if s.sessionId == nil {
		s.log("session creation")

		header := s.makeCreateSessionAction()

		response, err := DataJSONP(url, header, JitterDuration(sessionCreateTimeout, 0.2))
		if err != nil {
			s.log("session creation:", err)
			return
		}

		select {
		case array := <-response:
			if array == nil {
				s.log("session creation timeout")
				return
			}

			header := array.Index(0)
			if !s.handleSessionEvent(header) {
				return
			}

		case <-s.closeNotify:
			longPollClose(s, url)
			return
		}

		gotOnline = true
		s.connState("connected")
		s.connActive()
	} else {
		s.log("session resumption")

		// this ensures that we get an event quickly if the connection works,
		// and can update the status
		if err := longPollPing(s, url); err != nil {
			s.log("session resumption:", err)
			return
		}
	}

	if longPollTransfer(s, url) {
		gotOnline = true
	}

	return
}

// longPollTransfer sends buffered actions and polls for events.  It stops if
// two subsequent requests (of any kind) time out.
func longPollTransfer(s *Session, url string) (gotOnline bool) {
	var poller <-chan js.Object
	var sender <-chan js.Object
	var sendingId uint64
	var timeouts int

	s.numSent = 0

	for timeouts < 2 {
		if poller == nil {
			var err error

			header := s.makeResumeSessionAction(true)

			poller, err = DataJSONP(url, header, JitterDuration(minPollTimeout, 0.2))
			if err != nil {
				s.log("poll:", err)
				return
			}
		}

		if sender == nil && s.numSent < len(s.sendBuffer) {
			action := s.sendBuffer[s.numSent]

			if action.Payload != nil {
				var payload js.Object
				var err error

				frame := action.Payload.Index(0)

				if action.Header.Get("action").Str() == "update_user" {
					base64, err := ParseDataURI(frame)
					if err != nil {
						s.log("send:", err)
						return
					}

					payload = NewArray()
					payload.Call("push", base64)
				} else {
					if payload, err = ParseJSON(frame.Str()); err != nil {
						s.log("send:", err)
						return
					}
				}

				action.Header.Set("payload", payload)
			}

			action.Header.Set("session_id", s.sessionId)

			json, err := StringifyJSON(action.Header)

			action.Header.Delete("session_id")

			if err != nil {
				s.log("send:", err)
				return
			}

			var channel chan js.Object

			if len(json) <= maxJSONPSize {
				channel, err = StringJSONP(url, json, JitterDuration(minSendTimeout, 0.2))
			} else {
				action.Header.Set("caller_id", s.sessionParams.Get("user_id"))
				action.Header.Set("caller_auth", s.sessionParams.Get("user_auth"))

				channel, err = PostCall(action.Header, s.log, s.address)

				action.Header.Delete("caller_id")
				action.Header.Delete("caller_auth")
			}

			action.Header.Delete("payload")

			if err != nil {
				s.log("send:", err)
				return
			}

			if action.Id == 0 {
				s.sendBuffer = append(s.sendBuffer[:s.numSent], s.sendBuffer[s.numSent+1:]...)
			} else {
				sender = channel
				sendingId = action.Id
			}
		}

		var array js.Object

		select {
		case array = <-poller:
			if array == nil {
				s.log("poll timeout")
			}

			poller = nil
			s.connActive()

		case array = <-sender:
			if array == nil {
				s.log("send timeout")
			} else if sendingId > 0 {
				s.numSent++
			}

			sender = nil
			sendingId = 0

		case sending := <-s.sendNotify:
			if !sending {
				longPollClose(s, url)
				return
			}

			continue

		case <-s.closeNotify:
			longPollClose(s, url)
			return
		}

		if array == nil {
			timeouts++
			s.numSent = 0
			continue
		}

		timeouts = 0

		for i := 0; i < array.Length(); i++ {
			header := array.Index(i)
			payload := NewArray()

			if object := header.Get("payload"); !object.IsUndefined() {
				json, err := StringifyJSON(object)
				if err != nil {
					s.log("poll payload:", err)
					return
				}

				payload.Call("push", json)
			}

			ackedActionId, _, ok := s.handleEvent(header, payload)

			// poll acked the action being sent before we got send response?
			if sendingId > 0 && sendingId <= ackedActionId {
				sendingId = 0
				s.numSent++
			}

			if !ok {
				return
			}

			if !gotOnline {
				gotOnline = true
				s.connState("connected")
			}
		}
	}

	return
}

// longPollPing sends a ping action without an action_id.
func longPollPing(s *Session, url string) (err error) {
	header := map[string]interface{}{
		"action":     "ping",
		"session_id": s.sessionId,
	}

	_, err = DataJSONP(url, header, JitterDuration(minSendTimeout, 0.9))
	return
}

// longPollClose sends a close_session action without caring about the
// response.
func longPollClose(s *Session, url string) {
	header := map[string]interface{}{
		"action":     "close_session",
		"session_id": s.sessionId,
	}

	if _, err := DataJSONP(url, header, JitterDuration(minSendTimeout, 0.9)); err != nil {
		s.log("send:", err)
	}
}
