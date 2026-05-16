package vodka

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// SSEHandlerFunc is the handler signature for Server-Sent Events connections.
type SSEHandlerFunc func(*SSEContext)

// SSEContext wraps an active SSE connection with framework utilities.
// Use Send() to push events and Done() to detect client disconnection.
type SSEContext struct {
	Writer  http.ResponseWriter
	flusher http.Flusher
	Keys    map[string]any
	Params  httprouter.Params
	Request *http.Request
}

// Send pushes a named event with JSON-encoded data to the client.
// Returns an error if the connection is no longer writable.
//
//	c.Send("update", vodka.M{"value": 42})
func (sc *SSEContext) Send(event string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(sc.Writer, "event: %s\ndata: %s\n\n", event, payload)
	if err != nil {
		return err
	}
	sc.flusher.Flush()
	return nil
}

// SendData pushes data without a named event (browsers receive it as a "message" event).
func (sc *SSEContext) SendData(data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(sc.Writer, "data: %s\n\n", payload)
	if err != nil {
		return err
	}
	sc.flusher.Flush()
	return nil
}

// SendComment writes an SSE comment line, useful as a keep-alive ping.
func (sc *SSEContext) SendComment(comment string) error {
	_, err := fmt.Fprintf(sc.Writer, ": %s\n\n", comment)
	if err != nil {
		return err
	}
	sc.flusher.Flush()
	return nil
}

// Done returns a channel that is closed when the client disconnects.
// Use this to stop your event loop cleanly.
//
//	for {
//	    select {
//	    case <-c.Done():
//	        return
//	    default:
//	        c.Send("ping", vodka.M{"ts": time.Now()})
//	        time.Sleep(time.Second)
//	    }
//	}
func (sc *SSEContext) Done() <-chan struct{} {
	return sc.Request.Context().Done()
}

// Set stores a value in the context by key.
func (sc *SSEContext) Set(key string, value any) {
	if sc.Keys == nil {
		sc.Keys = make(map[string]any)
	}
	sc.Keys[key] = value
}

// Get retrieves a value previously stored via Set.
func (sc *SSEContext) Get(key string) (value any, exists bool) {
	value, exists = sc.Keys[key]
	return
}

// Param returns a URL path parameter by name.
func (sc *SSEContext) Param(key string) string {
	return sc.Params.ByName(key)
}

// Query returns a URL query parameter by name.
func (sc *SSEContext) Query(key string) string {
	return sc.Request.URL.Query().Get(key)
}

// IP returns the remote address of the client.
func (sc *SSEContext) IP() string {
	return sc.Request.RemoteAddr
}
