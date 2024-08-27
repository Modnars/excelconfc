package util

import (
	"log"
)

var (
	NO_COLORFUL_LOG bool

	prefixInfo  = "[INFO] "
	prefixError = "[ERROR] "

	colorfulPrefixInfo  = "\033[32m[INFO]\033[0m "
	colorfulPrefixError = "\033[31m[ERROR]\033[0m "
)

func infoLogPrefix() string {
	if NO_COLORFUL_LOG {
		return prefixInfo
	}
	return colorfulPrefixInfo
}

func errorLogPrefix() string {
	if NO_COLORFUL_LOG {
		return prefixError
	}
	return colorfulPrefixError
}

func LogInfo(format string, args ...any) {
	log.SetPrefix(infoLogPrefix())
	log.Printf(format+"\n", args...)
}

func LogError(format string, args ...any) {
	log.SetPrefix(errorLogPrefix())
	log.Printf(format+"\n", args...)
}
