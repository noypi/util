package util

import (
	"fmt"
)

func ToHex(bb []byte) string {
	return fmt.Sprintf("%x", bb)
}
