package simpleQueue

import (
	"log"
	"time"
)

//Item a delegate function that is added to the queue
type Item struct {
	Fn      func()
	Name    string
	StartAt time.Time
}

//Queue - basic queue object
type Queue struct {
	Items []Item
	Every time.Duration
}

//NewQueue - constructor
func NewQueue(seconds time.Duration) *Queue {
	queue := new(Queue)
	queue.Items = make([]Item, 0)
	queue.Every = seconds

	return queue
}

//Push an item onto the Queue - if re-adding then adjust the start time
func (q *Queue) Push(item Item) bool {

	for i := 0; i < len(q.Items); i++ {
		if q.Items[i].Name == item.Name {
			q.Items[i].StartAt = item.StartAt

			log.Printf("Delaying %v", item.Name)
			return true
		}
	}
	q.Items = append(q.Items, item)
	return true
}

func (q *Queue) next() (item Item, ok bool) {

	items := q.Items
	ok = false

	if len(items) > 0 {

		if items[0].StartAt.Before(time.Now()) {
			item, items = items[0], items[1:]
			q.Items = items
			ok = true
		} else {
			log.Printf("Waiting")
		}
	}

	return
}

//Run - Start Processing the Queue
func (q *Queue) Run() {
	time.AfterFunc(q.Every, func() {
		item, ok := q.next()
		if ok {
			item.Fn()
		}
		q.Run()
	})
}
