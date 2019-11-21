package base

import (
	"fmt"
	"io"
)

// netstringMissingComma indicates that a comma is missing at the end of the string.
type netstringMissingComma struct {
}

var _ error = (*netstringMissingComma)(nil)

func (netstringMissingComma) Error() string {
	return "Invalid NetString (missing ,)"
}

// netstringMissingColon indicates that there is a colon missing after the length of the string.
type netstringMissingColon struct {
}

var _ error = (*netstringMissingColon)(nil)

func (netstringMissingColon) Error() string {
	return "Invalid NetString (missing :)"
}

// netstringNoLengthSpecifier indicates that the length specifier is missing.
type netstringNoLengthSpecifier struct {
}

var _ error = (*netstringNoLengthSpecifier)(nil)

func (netstringNoLengthSpecifier) Error() string {
	return "Invalid NetString (no length specifier)"
}

// netstringLeadingZero indicates that there is a leading zero at the beginning of the string.
type netstringLeadingZero struct {
}

var _ error = (*netstringLeadingZero)(nil)

func (netstringLeadingZero) Error() string {
	return "Invalid NetString (leading zero)"
}

// netstringMaxLengthExceeded indicates that the maximum length of the string is exceeded.
type netstringMaxLengthExceeded struct {
	maxMessageLength int
}

var _ error = (*netstringMaxLengthExceeded)(nil)

func (mle netstringMaxLengthExceeded) Error() string {
	return fmt.Sprintf("Max data length exceeded: %d KB", mle.maxMessageLength/1024)
}

// netstringLengthSpecifierTooLarge indicates that the length specifier of the string is too large.
type netstringLengthSpecifierTooLarge struct {
}

var _ error = (*netstringLengthSpecifierTooLarge)(nil)

func (netstringLengthSpecifierTooLarge) Error() string {
	return "Length specifier must not exceed 9 characters"
}

// ReadNetStringsFromStream picks byte by byte aiming to split all bytes into length, ':', the
// actual message and ',' and return the message.
func ReadNetStringFromStream(stream io.Reader, maxMessageLength int) ([]byte, error) {
	length := 0
	leadingZero := false
	for readBytes := 0; ; readBytes++ {
		p := make([]byte, 1)
		if _, err := stream.Read(p); err != nil {
			return nil, err
		}

		if '0' <= p[0] && p[0] <= '9' {
			if readBytes == 9 {
				return nil, netstringLengthSpecifierTooLarge{}
			}

			if leadingZero {
				return nil, netstringLeadingZero{}
			}

			length = length*10 + int(p[0]-'0')

			if readBytes == 0 && p[0] == '0' {
				leadingZero = true
			}

		} else if p[0] == ':' {
			if readBytes == 0 {
				return nil, netstringNoLengthSpecifier{}
			}

			break
		} else {
			return nil, netstringMissingColon{}
		}
	}

	if maxMessageLength >= 0 && length > maxMessageLength {
		return nil, netstringMaxLengthExceeded{maxMessageLength}
	}

	p2 := make([]byte, length)
	if _, err := stream.Read(p2); err != nil {
		return nil, err
	}

	p := make([]byte, 1)
	if _, err := stream.Read(p); err != nil {
		return nil, err
	}

	if p[0] != ',' {
		return nil, netstringMissingComma{}
	}

	return p2, nil
}

// WriteNetStringToStream writes len(str):str, to the data stream.
func WriteNetStringToStream(stream io.Writer, str []byte) error {
	if _, err := fmt.Fprintf(stream, "%d:", len(str)); err != nil {
		return err
	}

	if _, err := stream.Write(str); err != nil {
		return err
	}

	_, err := stream.Write([]byte{','})
	return err
}
