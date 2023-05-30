package main

type MoveAnimation struct {
	tx       float64
	ty       float64
	Speed    float64
	DoneFunc func()
}
