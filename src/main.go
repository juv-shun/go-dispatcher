package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var logger = func() *Logger {
	l, err := NewLogger(logLevel)
	if err != nil {
		log.Fatalln(err)
	}
	return l
}()

var fetcher = NewDispatcher("fetcher", SampleFetchTask, fetcherNum, ContinualMode)
var processor = NewDispatcher("processor", SampleProcessTask, processorNum, QueueMode)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	processor.Start()
	logger.Infoln("processor started")

	fetcher.Start()
	logger.Infoln("fetcher started")

	exitChan := make(chan struct{})
	go waitComplete(exitChan)
	go waitStopSignal(processor, exitChan)

	<-exitChan
	fetcher.Stop()
	processor.Stop()
	logger.Infoln("successfully stopped")
}

func waitComplete(exitChan chan<- struct{}) {
	time.Sleep(1 * time.Minute)
	logger.Infoln("all tasks completed")

	exitChan <- struct{}{}
}

func waitStopSignal(processor *Dispatcher, exitChan chan<- struct{}) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan

	logger.Infoln("stop signal fired")
	exitChan <- struct{}{}
}
