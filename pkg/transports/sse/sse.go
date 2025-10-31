package sse

import (
	"net/http"

	"gopkg.in/antage/eventsource.v1"
)

func NewSSE() eventsource.EventSource {
	return eventsource.New(
		eventsource.DefaultSettings(),
		func(req *http.Request) [][]byte {
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
			}
		},
	)
}
