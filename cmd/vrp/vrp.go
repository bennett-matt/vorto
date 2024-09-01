package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"slices"
	"sort"
	"strings"
)

func (l load) isAssigned() bool {
	return l.assigned != nil
}

func (v *vrp) loadRoutes(path string) error {
	loads, err := loadRoutesFromPath(path)
	if err != nil {
		return err
	}

	for _, l := range loads {
		v.loadMap[l.id] = l
	}

	return nil
}

func (v *vrp) calculateSavings() []saving {
	var savings []saving
	for i := range v.loadMap {
		for j := range v.loadMap {
			if i != j {
				load1 := v.loadMap[i]
				load2 := v.loadMap[j]
				key := [2]int{i, j}
				s := saving{
					key: key,
					value: distanceBetweenPoints(load1.dropoff, v.depot) +
						distanceBetweenPoints(v.depot, load2.pickup) -
						distanceBetweenPoints(load1.dropoff, load2.pickup),
				}
				savings = append(savings, s)
			}
		}
	}
	sort.Slice(savings, func(i, j int) bool {
		return savings[i].value > savings[j].value
	})
	return savings
}

func (v *vrp) calculateDistance(loads []*load) float64 {
	if len(loads) == 0 {
		return 0.0
	}

	totalDistance := 0.0
	for i := 0; i < len(loads); i++ {
		totalDistance += loads[i].deliveryDistance
		if i != (len(loads) - 1) {
			totalDistance += distanceBetweenPoints(loads[i].dropoff, loads[i+1].pickup)
		}
	}
	totalDistance += distanceBetweenPoints(v.depot, loads[0].pickup)
	totalDistance += distanceBetweenPoints(loads[len(loads)-1].dropoff, v.depot)
	return totalDistance
}

func (v *vrp) assignLoads() {
	savings := v.calculateSavings()
	for _, s := range savings {
		load1 := v.loadMap[s.key[0]]
		load2 := v.loadMap[s.key[1]]
		if !load1.isAssigned() && !load2.isAssigned() {
			v.assignNewDriver(load1, load2)
		} else if load1.isAssigned() && !load2.isAssigned() {
			v.assignLoadToDriver(load1, load2)
		} else if !load1.isAssigned() && load2.isAssigned() {
			v.prependLoadToDriver(load1, load2)
		} else {
			v.mergeDrivers(load1, load2)
		}
	}

	v.assignUnassignedLoads()
}

func (v *vrp) assignNewDriver(l1, l2 *load) {
	cost := v.calculateDistance([]*load{l1, l2})
	if cost <= v.maxDistance {
		d := driver{id: uuid.NewString()}
		d.routes = append(d.routes, l1, l2)
		v.drivers = append(v.drivers, &d)
		l1.assigned = &d
		l2.assigned = &d
	}
}

func (v *vrp) assignLoadToDriver(l1, l2 *load) {
	d := l1.assigned
	i := slices.IndexFunc(d.routes, func(l *load) bool { return l1.id == l.id })
	if i == len(d.routes)-1 {
		cost := v.calculateDistance(append(d.routes, l2))
		if cost <= v.maxDistance {
			d.routes = append(d.routes, l2)
			l2.assigned = d
		}
	}
}

func (v *vrp) prependLoadToDriver(l1, l2 *load) {
	d := l2.assigned
	i := slices.IndexFunc(d.routes, func(l *load) bool { return l2.id == l.id })
	if i == 0 {
		cost := v.calculateDistance(append([]*load{l1}, d.routes...))
		if cost <= v.maxDistance {
			d.routes = append([]*load{l1}, d.routes...)
			l1.assigned = d
		}
	}
}

func (v *vrp) mergeDrivers(l1, l2 *load) {
	d1 := l1.assigned
	i1 := slices.IndexFunc(d1.routes, func(l *load) bool { return l1.id == l.id })
	d2 := l2.assigned
	i2 := slices.IndexFunc(d2.routes, func(l *load) bool { return l2.id == l.id })
	if (i1 == len(d1.routes)-1) && (i2 == 0) && d1.id != d2.id {
		cost := v.calculateDistance(append(d1.routes, d2.routes...))
		if cost <= v.maxDistance {
			d1.routes = append(d1.routes, d2.routes...)
			for _, l := range d2.routes {
				l.assigned = d1
			}
			driverIndex := slices.IndexFunc(v.drivers, func(d *driver) bool { return d.id == d2.id })
			if driverIndex != -1 {
				v.drivers = append(v.drivers[:driverIndex], v.drivers[driverIndex+1:]...)
			}
		}
	}
}

func (v *vrp) assignUnassignedLoads() {
	for _, l := range v.loadMap {
		if l.assigned == nil {
			d := driver{id: uuid.NewString()}
			d.routes = append(d.routes, l)
			v.drivers = append(v.drivers, &d)
			l.assigned = &d
		}
	}
}

func (v *vrp) output() {
	for _, d := range v.drivers {
		var ids []string
		for _, l := range d.routes {
			ids = append(ids, fmt.Sprint(l.id))
		}
		fmt.Printf("[%s]\n", strings.Join(ids, ","))
	}
}

func (v *vrp) getSolutionCost() float64 {
	totalDrivenMinutes := 0.0
	for _, d := range v.drivers {
		driverRoutes := []point{{0.0, 0.0}}
		for _, route := range d.routes {
			driverRoutes = append(driverRoutes, route.pickup, route.dropoff)
		}
		driverRoutes = append(driverRoutes, point{0.0, 0.0})
		for i := 1; i < len(driverRoutes); i++ {
			totalDrivenMinutes += distanceBetweenPoints(driverRoutes[i-1], driverRoutes[i])
		}
	}
	return float64(500*len(v.drivers)) + totalDrivenMinutes
}

func (v *vrp) debugDrivers() {
	log.Println("Drivers list:")
	for _, d := range v.drivers {
		log.Printf("Driver ID: %s, Loads: %v\n", d.id, d.routes)
	}
	log.Println("Load map:")
	for id, l := range v.loadMap {
		log.Printf("Load ID: %d, Assigned Driver: %s\n", id, l.assigned.id)
	}
}
