package util

import (
	"log"
)

var g_enabledebugging = false

func EnableDebugging() {
	g_enabledebugging = true
}

func DBG(as ...interface{}) {
	if g_enabledebugging {
		log.Println(as...)
	}
}

func DBGf(fmt string, as ...interface{}) {
	if g_enabledebugging {
		log.Printf(fmt, as...)
	}
}
