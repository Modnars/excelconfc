package util

import (
	"log"
)

var (
	COLORFUL_LOG bool

	prefixInfo  = "[INFO] "
	prefixError = "[ERROR] "

	colorfulPrefixInfo  = "\033[32m[INFO]\033[0m "
	colorfulPrefixError = "\033[31m[ERROR]\033[0m "
)

func infoLogPrefix() string {
	if COLORFUL_LOG {
		return colorfulPrefixInfo
	}
	return prefixInfo
}

func errorLogPrefix() string {
	if COLORFUL_LOG {
		return colorfulPrefixError
	}
	return prefixError
}

func LogInfo(format string, args ...any) {
	log.SetPrefix(infoLogPrefix())
	log.Printf(format+"\n", args...)
}

func LogError(format string, args ...any) {
	log.SetPrefix(errorLogPrefix())
	log.Printf(format+"\n", args...)
}
