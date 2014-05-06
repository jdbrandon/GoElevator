package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const SIZE int = 1000
const sleepdur time.Duration = time.Millisecond
const testdur time.Duration = time.Minute

var START time.Time = time.Now()
var prq [SIZE]int

func main() {
	done := make(chan bool)
	go insertstuff()
	go servicerequests(0, done)
	go servicerequests(1, done)
	<-done
}

func insertstuff() {
	for time.Since(START) < testdur {
		i := int(rand.Float32() * float32(SIZE))
		prq[i] += 1
		if prq[i] < 1{
			fmt.Println("There's a race afoot!!")
			os.Exit(1)
		}
		time.Sleep(time.Duration(rand.Float32()*float32(sleepdur)) * 3) // sleep 0 to 3 units of duration
	}
}

func servicerequests(id int, done chan<- bool) {
	up := true
	for i := 0; time.Since(START) < testdur; {
		for ; prq[i] > 0; prq[i]-- {
			fmt.Println("Goroutine", id, "Servicing", i, ",", prq[i], "requests remain")
			time.Sleep(time.Duration(rand.Float32() * float32(sleepdur) * 6)) // sleep 0 to 2 units of duraton
		}
		if up {
			i++
			if i == SIZE {
				up = false
				i -= 2 //we just serviced i-1 so skip it
			}
		} else {
			i--
			if i == -1 {
				up = true
				i += 2 //we just serviced 0 so skip it
			}
		}
	}
	done <- true
}
