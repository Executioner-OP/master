package queue

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Executioner-OP/master/db"
)

type Queue struct {
	Pending []TaskWrapper
	Taken   []TaskWrapper
	Mutex   sync.Mutex
}

type TaskWrapper struct {
	Task         db.ExecutionRequest
	TimeoutTimer *time.Timer
}

func (q *Queue) AddToPending(task db.ExecutionRequest, timeout time.Duration) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	timer := time.AfterFunc(timeout, func() {
		q.MoveToPending(task.ID.String())
	})

	q.Pending = append(q.Pending, TaskWrapper{Task: task, TimeoutTimer: timer})
	log.Printf("Task %v added to pending queue", task.ID)
}

func (q *Queue) MoveToPending(taskID string) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	for i, t := range q.Taken {
		if t.Task.ID.String() == taskID {
			t.TimeoutTimer.Stop()
			q.Taken = append(q.Taken[:i], q.Taken[i+1:]...)
			q.AddToPending(t.Task, time.Second*10)
			log.Printf("Task %v moved back to pending queue", taskID)
			return
		}
	}
}

func (q *Queue) PopFromPending() (db.ExecutionRequest, error) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if len(q.Pending) == 0 {
		return db.ExecutionRequest{}, errors.New("empty pending queue")
	}
	taskWrapper := q.Pending[0]
	q.Pending = q.Pending[1:]
	q.Taken = append(q.Taken, taskWrapper)
	return taskWrapper.Task, nil
}

func (q *Queue) PopFromTaken(taskID string) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	for i, t := range q.Taken {
		if t.Task.ID.String() == taskID {
			t.TimeoutTimer.Stop()
			q.Taken = append(q.Taken[:i], q.Taken[i+1:]...)
			log.Printf("Task %v completed and removed from taken queue", taskID)
			return
		}
	}
}

func (q *Queue) GetPendingLength() int {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	return len(q.Pending)
}

func (q *Queue) IsPendingEmpty() bool {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	return len(q.Pending) == 0
}
