package multipublisher

import "sync"

type MultiPublisher struct {
    limit int
    in chan<- interface{}
    subscribersIn []*chan<- interface{}
    subscribersOut []*<-chan interface{}
    mutex sync.Mutex 
}

func createMultiPublisher(limit int) (publisher MultiPublisher) {
    publisher.limit = limit
    publisher.in = make(chan<- interface{}, limit)
    return
}

func (this *MultiPublisher) subscribe() (outChan <-chan interface{}) {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    channel := make(chan interface{}, this.limit)

    var inChan chan<- interface{}
    inChan = channel
    outChan = channel

    this.subscribersIn = append(this.subscribersIn, &inChan)
    this.subscribersOut = append(this.subscribersOut, &outChan)
    return
}

func (this *MultiPublisher) unsubscribe(c *<-chan interface{}) {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    pos := -1
    for i, subscriber := range this.subscribersOut {
        if subscriber == c {
            pos = i
            break
        }
    }

    if pos != -1 {
        copy(this.subscribersIn[pos:], this.subscribersIn[pos + 1:])
        this.subscribersIn[len(this.subscribersIn) - 1] = nil
        this.subscribersIn = this.subscribersIn[:len(this.subscribersIn) - 1]

        copy(this.subscribersOut[pos:], this.subscribersOut[pos + 1:])
        this.subscribersOut[len(this.subscribersOut) - 1] = nil
        this.subscribersOut = this.subscribersOut[:len(this.subscribersOut) - 1]
    }
}

func (this *MultiPublisher) push(value interface{}) {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    for _, c := range this.subscribersIn {
        (*c) <- value
    }
}

func (this *MultiPublisher) close() {
    this.mutex.Lock()
    defer this.mutex.Unlock()

    for _, c := range this.subscribersIn {
        close(*c)
    }    
}
