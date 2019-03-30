package main

import (
	"log"
	"os"
	"strconv"
)

var processorNum = getIntEnv("ProcessorNum", 1)
var fetcherNum = getIntEnv("FetcherNum", 1)
var logLevel = os.Getenv("LogLevel")

func getIntEnv(envName string, defaultNum int) (num int) {
	envVariable := os.Getenv(envName)
	if envVariable == "" {
		return defaultNum
	}
	num, err := strconv.Atoi(envVariable)
	if err != nil {
		log.Fatalln(err)
	}
	return num
}
