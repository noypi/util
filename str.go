package util

import (
	"fmt"
	"unicode"
)

func ToHex(bb []byte) string {
	return fmt.Sprintf("%x", bb)
}

func StrFieldsQuoted(bQuotedOnly bool) func(r rune) bool {
	bInQuote := false
	return func(r rune) bool {
		if unicode.IsSpace(r) && !bInQuote {
			return true
		} else if unicode.Is(unicode.Quotation_Mark, r) {
			bInQuote = !bInQuote
			return true
		} else if bInQuote {
			return false
		}

		return bQuotedOnly
	}
}
