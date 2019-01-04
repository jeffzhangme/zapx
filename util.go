package zapx

import (
	"bytes"

	"go.uber.org/zap/buffer"
)

var (
	_pool = buffer.NewPool()
	// GetBufPool retrieves a buffer from the pool, creating one if necessary.
	GetBufPool = _pool.Get
)

// MultiError multiple error
type MultiError []error

func (p MultiError) Error() string {
	var errBuf bytes.Buffer
	for _, err := range p {
		errBuf.WriteString(err.Error())
		errBuf.WriteByte('\n')
	}
	return errBuf.String()
}
