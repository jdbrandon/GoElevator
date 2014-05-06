package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

const SIZE int = 1000
const sleepdur time.Duration = time.Nanosecond * 2
const testdur time.Duration = time.Second * 15

var mtx sync.Mutex
var prq [SIZE]int
var START time.Time = time.Now()

func main() {
	n := 0
	idone := make(chan bool)
	done := make(chan bool)
	go insertstuff(idone)
	go servicerequests(n++, done, idone)
	go servicerequests(n++, done, idone)
	go servicerequests(n++, done, idone)
	for i := 0; i < n; i++ {
		<-done
	}
	sum := 0
	for _, v := range prq {
		sum += v
	}

	fmt.Println("Remaining requests after", testdur, "is", sum)
}

func insertstuff(done chan<- bool) {
	var hold int
	for time.Since(START) < testdur {
		i := int(rand.Float32() * float32(SIZE))
		mtx.Lock()
		prq[i] += 1
		hold = prq[i]
		mtx.Unlock()
		if hold < 1 {
			fmt.Println("There's a race afoot!!")
			os.Exit(1)
		}
		time.Sleep(time.Duration(rand.Float32()*float32(sleepdur))) // sleep 0 to 1 units of duration
	}
	close(done)
}

func servicerequests(id int, done chan<- bool, waitfor <-chan bool) {
	up := true
	for i := 0; time.Since(START) < testdur; {
		for mtx.Lock(); prq[i] > 0; {
			fmt.Println("Goroutine", id, "Servicing", i, "in region", getRegion(i), "R:", prq[i]-1)
			prq[i]--
		}
		mtx.Unlock()
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
	<-waitfor
	done <- true
}

func getRegion(i int) (s string) {
	const degree int = 60
	s = "\t"
	for j := 0; j < degree; j++ {
		if int(i/(SIZE/degree)) == j {
			s += fmt.Sprintf("1")
		} else {
			s += fmt.Sprintf("0")
		}
	}
	s += "\t"
	return
}
