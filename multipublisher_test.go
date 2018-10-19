package multipublisher

import "testing"

func TestMultiPublisher(t *testing.T) {
    publisher := CreateMultiPublisher(5)

    ready := make(chan bool, 2)
    done := make(chan bool, 2)
    f := func () {
        c := publisher.Subscribe()
        defer publisher.Unsubscribe(&c)

        ready <- true

        prevValue := (<-c).(int)
        for value := range c {
            if value.(int) != prevValue + 1 {
                t.Error("Wrong value")
            }
            prevValue = value.(int)
        }

        done <- true
    }

    go f()
    go f()
    <-ready
    <-ready

    for i := 0; i < 5; i++ {
        publisher.Push(i)
    }
    publisher.Close()

    <-done
    <-done
}
