package main

import (
	"testing"

	"github.com/bmizerany/assert"
)

func TestSampleFetchTask(t *testing.T) {
	jobChan := make(chan interface{})

	go getTaskQueue(jobChan)
	err := SampleFetchTask(1)
	assert.Equal(t, err, nil)
	job := <-jobChan
	assert.Equal(t, job, 1)
}

func getTaskQueue(jobChan chan interface{}) {
	job := <-processor.taskQueue
	jobChan <- job
}

func TestSampleProcessTask(t *testing.T) {
	err := SampleProcessTask(1)
	assert.Equal(t, err, nil)
}
