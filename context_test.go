package vodka

import "testing"

func TestAbortStopsChain(t *testing.T) {
	calls := []string{}

	// First handler aborts
	handler1 := func(c *Context) {
		calls = append(calls, "h1")
		c.Abort()
	}

	// This should NOT run
	handler2 := func(c *Context) {
		calls = append(calls, "h2")
	}

	c := &Context{handlers: []HandlerFunc{handler1, handler2}, index: -1}
	c.Next()

	if len(calls) != 1 || calls[0] != "h1" {
		t.Errorf("got %v, want [h1]", calls)
	}
}