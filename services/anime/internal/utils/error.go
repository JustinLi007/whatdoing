package utils

import "log"

func RequireNoError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
