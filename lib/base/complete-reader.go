package base

import "io"

type CompleteReader struct {
	reader io.Reader
}

var _ io.Reader = CompleteReader{}

func (cr CompleteReader) Read(p []byte) (n int, err error) {
	for len(p) > 0 {
		m, errRd := cr.reader.Read(p)
		n += m

		if errRd != nil {
			err = errRd
			break
		}

		p = p[m:]
	}

	return
}
