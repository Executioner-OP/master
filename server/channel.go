package server

import (
	"log"
	"time"

	"github.com/Executioner-OP/master/db"
)

var TaskChannel chan db.ExecutionRequest
var PendingChannel chan PendingTask

type PendingTask struct {
	// don't know any data type helpful to track timestamp
	TimeStamp time.Time
	Task      db.ExecutionRequest
}

type PendingTaskQueue struct {
	Elements []PendingTask
}

func (pq *PendingTaskQueue) AddPendingTask(pendingTask PendingTask) {
	// Calculate the difference between the current time and the timestamp of the pending task
	currentTime := time.Now()
	timeDiff := currentTime.Sub(pendingTask.TimeStamp)

	// If the difference is less than 15 seconds, wait for the remaining time
	if timeDiff < 15*time.Second {
		waitTime := 15*time.Second - timeDiff
		time.Sleep(waitTime)
	}
	res, err := db.CheckPendingTask(pendingTask.Task.ID)
	if err != nil {
		log.Fatalf("Failed in checking pending task: %v", err)
	}
	if res == false {
		TaskChannel <- pendingTask.Task
		log.Printf("Pending Task added to task_queue again")
	}
}
