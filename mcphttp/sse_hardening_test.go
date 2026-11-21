package mcphttp

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// stubTransport answers requests in memory, so the test opens no connection.
type stubTransport func(*http.Request) (*http.Response, error)

func (f stubTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// A call waiting on the SSE stream must fail as soon as the stream is closed.
func TestLegacyCallFailsWhenStreamClosed(t *testing.T) {
	done := make(chan struct{})
	close(done)
	c := &Client{
		messageURL: "http://mcp.test/message",
		http: &http.Client{Transport: stubTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusAccepted, Body: io.NopCloser(strings.NewReader(""))}, nil
		})},
		pending:    map[int]chan *rpcResponse{},
		streamErr:  make(chan error, 1),
		streamDone: done,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := c.legacyRoundTrip(ctx, rpcRequest{JSONRPC: "2.0", ID: 1, Method: "tools/list"}, true); err == nil {
		t.Fatal("expected an error once the SSE stream is closed")
	}
	if ctx.Err() != nil {
		t.Fatal("the call waited for its own timeout instead of seeing the closed stream")
	}
}

// The POST endpoint may only be on the configured origin: another host or scheme is
// refused, and so is an endpoint carrying userinfo (v0.7.0 SSE hardening).
func TestSSEEndpointSameOriginOnly(t *testing.T) {
	base, _ := url.Parse("https://mcp.example.com/sse")
	okCases := map[string]string{
		"data: /message?session=1":              "https://mcp.example.com/message?session=1",
		"data: https://mcp.example.com/message": "https://mcp.example.com/message",
	}
	for line, want := range okCases {
		got, err := readSSEEndpoint(bufio.NewReader(strings.NewReader(line+"\n")), base)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", line, err)
		}
		if got != want {
			t.Fatalf("%s: endpoint %q, want %q", line, got, want)
		}
	}
	badCases := []string{
		"data: https://evil.example.net/message",
		"data: http://mcp.example.com/message",
		"data: https://mcp.example.com:8443/message",
		"data: https://user:pass@mcp.example.com/message",
	}
	for _, line := range badCases {
		if got, err := readSSEEndpoint(bufio.NewReader(strings.NewReader(line+"\n")), base); err == nil {
			t.Fatalf("%s: endpoint %q accepted, want a refusal", line, got)
		}
	}
}
