package fallback

import (
	"io"
)

type (
	Fallback interface {
		Do(sharedHandler func(Helper))
		io.Closer
	}

	Helper interface {
		Exclusive(exclusiveHandler func(), slowAsync func())
	}
)
