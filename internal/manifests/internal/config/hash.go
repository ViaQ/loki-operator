package config

import (
	"crypto/sha1"
	"fmt"
)

// Sha1sum returns the sha1 checksum for a byte slice
// or an error if the contents cannot be written into
// the internal buffer.
func Sha1sum(c []byte) (string, error) {
	s := sha1.New()
	_, err := s.Write(c)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", s.Sum(nil)), nil
}
