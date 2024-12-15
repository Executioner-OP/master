package queue

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Executioner-OP/master/db"
)

type Queue struct {
	Elements []TaskWrapper
	Mutex    sync.Mutex
}

type TaskWrapper struct {
	Task         db.ExecutionRequest
	TimeoutTimer *time.Timer
}

func (q *Queue) Add(task db.ExecutionRequest, timeout time.Duration) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	timer := time.AfterFunc(timeout, func() {
		q.Mutex.Lock()
		defer q.Mutex.Unlock()
		for i, t := range q.Elements {
			if t.Task.ID == task.ID {
				q.Elements = append(q.Elements[:i], q.Elements[i+1:]...)
				q.Add(task, timeout)
				log.Printf("Task %v re-added to the queue after timeout", task.ID)
				break
			}
		}
	})

	q.Elements = append(q.Elements, TaskWrapper{Task: task, TimeoutTimer: timer})
	log.Printf("Task %v added to the queue", task.ID)
}

func (q *Queue) Pop() (db.ExecutionRequest, error) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if q.IsEmpty() {
		return db.ExecutionRequest{}, errors.New("empty queue")
	}
	taskWrapper := q.Elements[0]
	taskWrapper.TimeoutTimer.Stop()
	q.Elements = q.Elements[1:]
	return taskWrapper.Task, nil
}

func (q *Queue) GetLength() int {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	return len(q.Elements)
}

func (q *Queue) IsEmpty() bool {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	return len(q.Elements) == 0
}
