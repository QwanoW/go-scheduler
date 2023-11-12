package main

import (
	"scheduler/internal"
)

func main() {
	// keep app alive
	go internal.KeepAlive()

	// start bot
	internal.StartBot()
}
