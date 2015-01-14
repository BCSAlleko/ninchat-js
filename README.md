JavaScript utilities for use with the [Ninchat](https://ninchat.com) API.

Points of interest:

- `src/ninchatclient/` contains Go sources for NinchatClient, a library for
  accessing `api.ninchat.com` from a web browser.

- `gen/ninchatclient.js` contains JavaScript sources generated with
  [GopherJS](https://github.com/gopherjs/gopherjs).  Regenerate with `make`
  (requires [Go](https://golang.org)).

- `doc/ninchatclient.js` contains API documentation.

- `example/test.js` demonstrates usage.

NinchatClient supports IE8 and later, but the code generated by GopherJS
requires ES5 and typed array shims.  When WebSocket transport is not available,
NinchatClient doesn't actually use the typed arrays, so stubs are enough.
