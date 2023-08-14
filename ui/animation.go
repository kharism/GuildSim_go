package main

import "time"

type MoveAnimation struct {
	tx        float64
	ty        float64
	Speed     float64
	SleepPre  time.Duration
	SleepPost time.Duration
	DoneFunc  func()
}
