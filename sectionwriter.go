package sectionwriter

import "io"

// NewSectionWriter returns a SectionWriter that writes to w
// starting at offset off and stops after n bytes.
func NewSectionWriter(w io.WriterAt, off int64, n int64) *SectionWriter {
	return &SectionWriter{w, off, off + n}
}

// SectionWriter implements Write, Seek, and WriteAt on a section
// of an underlying WriterAt.
type SectionWriter struct {
	w     io.WriterAt
	off   int64
	limit int64
}

// Write writes the data in p to the underlying WriterAt (w).
// When len(p) is greater than the space left in w, bytes that
// can be written in the remaining space will be, and the
// number of bytes written along with an io.ErrShortWrite will
// be returned
func (s *SectionWriter) Write(p []byte) (n int, err error) {
	if max := s.limit - s.off; int64(len(p)) > max {
		p = p[0:max]
		err = io.ErrShortWrite
	}

	n, writeErr := s.w.WriteAt(p, s.off)
	s.off += int64(n)
	if writeErr != nil {
		err = writeErr
		return
	}

	return
}
