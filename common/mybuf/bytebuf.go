package mybuf

import (
	"bytes"
	"io"
)

// Concat concatenates one or more bytes.Buffers to the base buffer.  The
// function returns the total number of bytes concatenated and any errors.
func Concat(base io.ReaderFrom, add ...*bytes.Buffer) (int64, error) {
	var total int64
	for _, a := range add {
		n, e := base.ReadFrom(a)
		if e != nil {
			return 0, e
		}
		total += n
	}
	return total, nil
}

// ToBSlice reads the unread contents of a bytes.Buffer and returns it as a slice of bytes.
func ToBSlice(b *bytes.Buffer) []byte {
	return b.Next(b.Len())
}

// Copy returns a copy of a buffer.  The length of the buffer and any errors
// are also returned.
func Copy(orig *bytes.Buffer) (*bytes.Buffer, int, error) {
	var copy1, copy2 bytes.Buffer

	n, err := io.Copy(&copy1, io.TeeReader(orig, &copy2))
	if err != nil {
		return nil, 0, err
	}
	*orig = copy1

	return &copy2, int(n), nil
}
