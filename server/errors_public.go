package server

import "log"

// publicErr logs the real error and returns a text for the client that holds no server paths.
func publicErr(err error) string {
	log.Printf("[http] %v", err)
	return "error interno; revisa el registro del servidor"
}
