package main

import (
	"fmt"
	"time"
)

func mySleep(t int) {
	<-time.After(time.Second * time.Duration(t))
}

func main() {
	mySleep(5)

	fmt.Print("done")
}
