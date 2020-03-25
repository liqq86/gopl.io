// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 246.

// Countdown implements the countdown for a rocket launch.
package main

// NOTE: the ticker goroutine never terminates if the launch is aborted.
// This is a "goroutine leak".

import (
	"fmt"
	"os"
	"time"
)

//!+

func main() {
	// ...create abort channel...

	//!-

	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()

	//!+
	fmt.Println("Commencing countdown.  Press return to abort.")
	//! it could cause go routine leak, use time.NewTicker(1*time.Second) instead
	//tick := time.Tick(1 * time.Second)
	tick := time.NewTicker(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick.C:
			// Do nothing.
		case <-abort:
			fmt.Println("Launch aborted!")
			return
		}
	}
	tick.Stop()
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
