package sectionwriter_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kdar/stringio"
	sw "github.com/wjessop/sectionwriter"
)

func isSectionWriter(t interface{}) bool {
	switch t.(type) {
	case *sw.SectionWriter:
		return true
	default:
		return false
	}
}

func TestNewSectionWriter(t *testing.T) {
	sio := stringio.New()
	s := sw.NewSectionWriter(sio, 10, 100)
	assert.True(t, isSectionWriter(s), "should be a SectionWriter")
}

func TestWrite(t *testing.T) {
	data := []byte("0123456789abcdefghij")

	sio := stringio.New()
	s1 := sw.NewSectionWriter(sio, 0, 10)
	s2 := sw.NewSectionWriter(sio, 10, 10)

	n1, err1 := s1.Write(data[0:10])
	n2, err2 := s2.Write(data[10:len(data)])

	assert.NoError(t, err1, "Write 1 should not error")
	assert.NoError(t, err2, "Write 2 should not error")

	assert.Equal(t, 10, n1, "Write 1 should write 10 bytes")
	assert.Equal(t, 10, n2, "Write 2 should write 10 bytes")

	assert.Equal(t, data, sio.GetValueBytes(), "output IO object should match input")
}

type NopeWriter struct {
}

func (n *NopeWriter) WriteAt(p []byte, off int64) (int, error) {
	return 0, io.EOF
}

func TestWriteError(t *testing.T) {
	s := sw.NewSectionWriter(&NopeWriter{}, 0, 10)
	n, err := s.Write([]byte("foo"))

	assert.Equal(t, 0, n, "bytes written should be 0")
	assert.EqualError(t, err, "EOF", "error should be EOF")
}

func TestReturnsShortWriteWhenPassedTooMuchData(t *testing.T) {
	sio := stringio.New()
	s := sw.NewSectionWriter(sio, 5, 5)

	n, err := s.Write([]byte("1234567"))

	assert.Equal(t, 5, n, "bytes written should be 5")
	assert.EqualError(t, err, "short write", "error should be ErrShortWrite")
}

func TestSerialWritesAdvancesOffset(t *testing.T) {
	sio := stringio.New()

	s := sw.NewSectionWriter(sio, 0, 10)
	s.Write([]byte("01234"))
	s.Write([]byte("56789"))

	assert.Equal(t, "0123456789", sio.GetValueString())
}
