package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
	fmt.Println(t.Format(time.TimeOnly))
}
