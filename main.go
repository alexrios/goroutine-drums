package main

import (
	"sync"
)

func main() {
	// - - - - - - - -
	// 1 2 3 4 5 6 7 8

	drum := Drum{
		Tempo: 500,
	}

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go drum.hihat("xxxxxxxx", wg)
	go drum.kick("x---x---", wg)
	go drum.snare("--x---x-", wg)
	wg.Wait()

}
