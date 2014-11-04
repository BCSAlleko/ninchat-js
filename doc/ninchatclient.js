/**
 * Ninchat API connection library.
 *
 * @see https://ninchat.com/api/v2
 *
 * @module NinchatClient
 */

/**
 * Call the sessionless API.  Implemented using the GET method.  The
 * returned promise will be resolved with an event header array as a
 * parameter to the callback functions, or rejected without a parameter.
 *
 * @see https://ninchat.com/api/v2#sessionless-http-calling
 *
 * @param {object}   header     Action parameters to send.
 * @param {function} [onLog]    Message logger.
 * @param {string}   [address]  Alternative API endpoint.
 *
 * @return {NinchatClient~Promise}
 */
function call(header, onLog, address) {}

/**
 * Create an uninitialized Session object.
 *
 * @return {NinchatClient~Session}
 */
function newSession() {}

/**
 * Convert an event's payload part to a string.
 *
 * @param {object} data
 *
 * @return {string}
 */
function stringifyFrame(data) {}

/**
 * Session hides the details of API connection management.  It needs to be
 * initialized by calling at least the onSessionEvent, onEvent and
 * setParams methods.  After that the open method is used to make a
 * connection to the server.  Finally, the close method disconnects from
 * the server.
 *
 * Session objects may be instantiated only via the newSession function.
 *
 * @class
 */
function Session() {

	/**
	 * Set the session creation handler.  It will be invoked with a
	 * "session_created" or an "error" event header as a parameter.
	 *
	 * If another "session_created" event is received, it means that the
	 * previous session was lost.  If an "error" event is received, it
	 * means that a session can't be established (without help from the
	 * client code), and the Session object is closed.
	 *
	 * @param {function}  callback
	 */
	this.onSessionEvent = function(callback) {}

	/**
	 * Set the handler for in-session events. It will be invoked with an
	 * event header and a payload array parameter.
	 *
	 * "error" events received via this callback are not fatal.
	 *
	 * @param {function}  callback
	 */
	this.onEvent = function(callback) {}

	/**
	 * Set an optional connection state change monitor.  It will be called with
	 * one of the following strings:
	 *
	 * - "connecting"
	 * - "connected"
	 * - "disconnected"
	 *
	 * @param {function}  callback
	 */
	this.onConnState = function(callback) {}

	/**
	 * Set an optional connection activity monitor.  It will be called whenever
	 * data has been received on the connection.
	 *
	 * @param {function}  callback
	 */
	this.onConnActive = function(callback) {}

	/**
	 * Set an optional message logger.  It will be called with a single string
	 * argument.
	 *
	 * @param {function}  callback
	 */
	this.onLog = function(callback) {}

	/**
	 * Set "create_session" action parameters.
	 *
	 * @param {object}  params
	 */
	this.setParams = function(params) {}

	/**
	 * Force a specific network transport implementation to be used.
	 * Currently only "longpoll" may be specified.
	 *
	 * @param {string}  name
	 */
	this.setTransport = function(name) {}

	/**
	 * Use an alternative API endpoint.
	 *
	 * @param {string}  address
	 */
	this.setAddress = function(address) {}

	/**
	 * Create a session on the server.
	 */
	this.open = function() {}

	/**
	 * Close the session on the server.
	 */
	this.close = function() {}

	/**
	 * Does the currently active transport implementation support binary
	 * payloads?  The result is undefined until a session has been established.
	 *
	 * @return {boolean}
	 */
	this.binarySupported = function() {}

	/**
	 * Send an action.
	 *
	 * To send an action without an "action_id" parameter, specify it as
	 * null.  Otherwise an "action_id" is generated automatically.
	 *
	 * If an "action_id" is used, a promise is returned.  It may be used to
	 * wait for a reply from the server; the promise will be resolved with
	 * an event header and a payload array parameter.  If the Session
	 * object is closed before a reply is received, the promise will be
	 * rejected without a parameter.
	 *
	 * @param {object}  header     Action parameters to send.
	 * @param {array}   [payload]  Consists of (already encoded) data
	 *                             frames.
	 *
	 * @return {NinchatClient~Promise}
	 */
	this.send = function(header, payload) {}

}

/**
 * @class
 */
function Promise() {

	/**
	 * Add callback(s) to be called when the promise is resolved or
	 * rejected.
	 *
	 * Promise objects may not be instantiated directly.
	 *
	 * @param {function}  [success]
	 * @param {function}  [failure]
	 */
	this.then = function(success, failure) {};

}
