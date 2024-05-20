package utils

import "log"

func AssertErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func AssertConv(ok bool) {
	if !ok {
		log.Panic("invalid widget conversion")
	}
}
