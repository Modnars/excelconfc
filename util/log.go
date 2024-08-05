package util

import "log"

func LogInfo(format string, args ...any) {
	log.SetPrefix("[INFO] ")
	log.Printf(format+"\n", args...)
}

func LogError(format string, args ...any) {
	log.SetPrefix("[ERROR] ")
	log.Printf(format+"\n", args...)
}
