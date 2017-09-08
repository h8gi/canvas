// draw context
package canvas

import (
	"sync"

	"github.com/fogleman/gg"
)

type Context struct {
	gg.Context
	mu sync.Mutex
}

func NewContext(width, height int) *Context {
	return &Context{
		*gg.NewContext(width, height),
		*(new(sync.Mutex)),
	}
}
