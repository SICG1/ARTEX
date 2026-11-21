package server

import (
	"io"
	"net/http"
)

// maxJSONBody caps a JSON request body at 16 MiB. Uploads keep their own, larger limits.
const maxJSONBody = 16 << 20

// jsonBody returns the request body capped at maxJSONBody, so one request cannot exhaust memory.
// A body over the cap is cut short and fails to decode, which handlers answer with 400.
func jsonBody(r *http.Request) io.Reader {
	return io.LimitReader(r.Body, maxJSONBody)
}
