package base

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestReadNetStringFromStream(t *testing.T) {
	assertReadNetStringFromStream(t, "", 13, "", io.EOF)
	assertReadNetStringFromStream(t, "6", 13, "", io.EOF)
	assertReadNetStringFromStream(t, "6:", 13, "", io.EOF)
	assertReadNetStringFromStream(t, "6:fg", 13, "", io.EOF)
	assertReadNetStringFromStream(t, "6:foobar", 13, "", io.EOF)
	assertReadNetStringFromStream(t, "66666666666666", 20, "", netstringLengthSpecifierTooLarge{})
	assertReadNetStringFromStream(t, "0006:foobar,", 15, "", netstringLeadingZero{})
	assertReadNetStringFromStream(t, ":foobar,", 13, "", netstringNoLengthSpecifier{})
	assertReadNetStringFromStream(t, "6foobar,", 12, "", netstringMissingColon{})
	assertReadNetStringFromStream(t, "6:foobar,", 0, "", netstringMaxLengthExceeded{})
	assertReadNetStringFromStream(t, "2:dog,", 13, "", netstringMissingComma{})
	assertReadNetStringFromStream(t, "6:foobar,", 6, "foobar", nil)
	assertReadNetStringFromStream(t, "0:,", 0, "", nil)
}

func assertReadNetStringFromStream(t *testing.T, in string, maxmsglen int, out string, expectederr error) {
	t.Helper()

	buf := &bytes.Buffer{}
	buf.Write([]byte(in))

	if ret, err := ReadNetStringFromStream(buf, maxmsglen); err != expectederr {
		t.Errorf("buf := &bytes.Buffer{}; buf.Write([]byte(%#v)); _, err := ReadNetStringFromStream(buf, %d); err != %#v", in, maxmsglen, expectederr)
	} else if err == nil && bytes.Compare(ret, []byte(out)) != 0 {
		t.Errorf("buf := &bytes.Buffer{}; buf.Write([]byte(%#v)); ret, _ := ReadNetStringFromStream(buf, %d); bytes.Compare(ret, []byte(%#v)) != 0", in, maxmsglen, out)
	}
}

func TestWriteNetStringToStream(t *testing.T) {
	assertWriteNetStringToStream(t, "foobar", "6:foobar,")
	assertWriteNetStringToStream(t, "", "0:,")

}

func assertWriteNetStringToStream(t *testing.T, in string, out string) {
	t.Helper()

	buf := &bytes.Buffer{}
	if WriteNetStringToStream(buf, []byte(in)) != nil {
		t.Errorf("WriteNetStringToStream(&bytes.Buffer{}, []byte(%#v)) != nil", in)
	} else if !reflect.DeepEqual(buf.Bytes(), []byte(out)) {
		t.Errorf("buf := &bytes.Buffer{}; WriteNetStringToStream(buf, []byte(%#v)); !reflect.DeepEqual(buf.Bytes(), []byte(%#v))", in, out)
	}
}
