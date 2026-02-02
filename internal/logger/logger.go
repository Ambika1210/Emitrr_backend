package logger

import (
	"fmt"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func Init() {
	// No special initialization needed for fmt.Printf colors
}

func Info(msg string) {
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s %sINFO%s: %s\n", t, colorGreen, colorReset, msg)
}

func Error(msg string, err error) {
	t := time.Now().Format("2006-01-02 15:04:05")
	if err != nil {
		fmt.Printf("%s %sERROR%s: %s | Error: %v\n", t, colorRed, colorReset, msg, err)
	} else {
		fmt.Printf("%s %sERROR%s: %s\n", t, colorRed, colorReset, msg)
	}
}

func Warn(msg string) {
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s %sWARN%s: %s\n", t, colorYellow, colorReset, msg)
}

func Debug(msg string) {
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s %sDEBUG%s: %s\n", t, colorBlue, colorReset, msg)
}
