package multipublisher

import "sync"

type MultiPublisher struct {
    closed bool
    limit int
    in chan<- interface{}
    subscribers []*chan interface{}
    mutex sync.Mutex 
}

func CreateMultiPublisher(limit int) (publisher MultiPublisher) {
    publisher.closed = false
    publisher.limit = limit
    publisher.in = make(chan<- interface{}, limit)
    return
}

func (this *MultiPublisher) Subscribe() (outChan <-chan interface{}) {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    this.checkClosed()

    channel := make(chan interface{}, this.limit)
    outChan = channel

    this.subscribers = append(this.subscribers, &channel)
    return
}

func (this *MultiPublisher) Unsubscribe(c *<-chan interface{}) {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    pos := -1
    for i, subscriber := range this.subscribers {
        if *subscriber == *c {
            pos = i
            break
        }
    }

    if pos != -1 {
        close(*this.subscribers[pos])
        copy(this.subscribers[pos:], this.subscribers[pos + 1:])
        this.subscribers[len(this.subscribers) - 1] = nil
        this.subscribers = this.subscribers[:len(this.subscribers) - 1]
    }
}

func (this *MultiPublisher) Push(value interface{}) {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    this.checkClosed()

    for _, c := range this.subscribers {
        (*c) <- value
    }
}

func (this *MultiPublisher) Close() {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    this.checkClosed()
    this.closed = true

    for _, c := range this.subscribers {
        close(*c)
    }
    this.subscribers = this.subscribers[:0]
}

func (this *MultiPublisher) checkClosed() {
    if this.closed {
        panic("Publisher is closed")
    }
}

