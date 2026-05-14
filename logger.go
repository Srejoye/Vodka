package vodka

import (
	"log"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		log.Printf(Blue+"%s %s"+Reset, c.Request.Method, c.Request.URL.Path)
	}
}

