package queue

import (
	"errors"
	"log"

	"github.com/Executioner-OP/master/db"
)

type Queue struct {
	Elements []db.ExecutionRequest
}

func (q *Queue) Add(task db.ExecutionRequest) {
	q.Elements = append(q.Elements, task)
	log.Printf("Task %v, added in the queue", task.ID)
}

func (q *Queue) Pop() error {
	if q.IsEmpty() {
		return errors.New("empty queue")
	}
	q.Elements = q.Elements[1:]
	return nil
}

func (q *Queue) GetLength() int {
	return len(q.Elements)
}

func (q *Queue) IsEmpty() bool {
	return len(q.Elements) == 0
}
