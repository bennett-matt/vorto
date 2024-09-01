package main

type point struct {
	X, Y float64
}

type load struct {
	id               int
	pickup           point
	dropoff          point
	assigned         *driver
	deliveryDistance float64
}

type driver struct {
	id       string
	capacity int
	routes   []*load
}

type vrp struct {
	drivers     []*driver
	loadMap     map[int]*load
	depot       point
	maxDistance float64
}

type saving struct {
	key   [2]int
	value float64
}

func newVRP() *vrp {
	return &vrp{
		loadMap:     make(map[int]*load),
		depot:       point{0, 0},
		maxDistance: 12 * 60,
	}
}
