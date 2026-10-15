package lab

import (
	"encoding/binary"
	"errors"
	"io"
)

const maxFrame = 1 << 20

func writeAll(w io.Writer, p []byte) error {
	for len(p) > 0 {
		n, err := w.Write(p)
		if n < 0 || n > len(p) {
			return io.ErrShortWrite
		}
		p = p[n:]
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}
func WriteFrame(w io.Writer, payload []byte) error {
	if len(payload) > maxFrame {
		return errors.New("frame too large")
	}
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))
	if err := writeAll(w, header[:]); err != nil {
		return err
	}
	return writeAll(w, payload)
}
func ReadFrame(r io.Reader) ([]byte, error) {
	var header [4]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(header[:])
	if n > maxFrame {
		return nil, errors.New("frame too large")
	}
	p := make([]byte, int(n))
	if _, err := io.ReadFull(r, p); err != nil {
		return nil, err
	}
	return p, nil
}
