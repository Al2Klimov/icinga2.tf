package base

import "io"

func ReadNetStringFromStream(stream io.Reader, maxMessageLength int) ([]byte, error) {
	panic("implement me")
}

func WriteNetStringToStream(stream io.Writer, str []byte) error {
	panic("implement me")
}
