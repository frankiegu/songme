package utility

import (
	"crypto/sha256"
	"fmt"
)

// SHA256String returns SHA256 string.
func SHA256String(input string) string {
	sum := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", sum)
}
