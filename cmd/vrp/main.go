package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./vrp <input.csv>")
		return
	}
	v := newVRP()
	err := v.loadRoutes(os.Args[1])
	if err != nil {
		fmt.Println("Error loading problem:", err)
		return
	}
	v.assignLoads()
	v.output()
	end := time.Now()

	log.Println("total cost: ", v.getSolutionCost())
	log.Println("time taken: ", end.Sub(start).Milliseconds())
}
