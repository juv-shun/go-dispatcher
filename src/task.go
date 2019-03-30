package main

import (
	"time"
)

func SampleFetchTask(job interface{}) error {
	processor.Add(job)
	return nil
}

func SampleProcessTask(job interface{}) error {
	time.Sleep(5 * time.Second)
	return nil
}
