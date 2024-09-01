package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func loadRoutesFromPath(path string) ([]*load, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ' ' // Assuming space delimiter
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var loads []*load
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		if len(record) != 3 {
			return nil, fmt.Errorf("wrong number of fields in line %d: %v", i+1, record)
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid load number in line %d: %v", i+1, record[0])
		}

		pickup, err := parsePoint(record[1])
		if err != nil {
			return nil, fmt.Errorf("invalid pickup point in line %d: %v", i+1, record[1])
		}

		dropoff, err := parsePoint(record[2])
		if err != nil {
			return nil, fmt.Errorf("invalid dropoff point in line %d: %v", i+1, record[2])
		}

		loads = append(loads, &load{
			id:               id,
			pickup:           pickup,
			dropoff:          dropoff,
			deliveryDistance: distanceBetweenPoints(pickup, dropoff),
		})
	}

	return loads, nil
}

func distanceBetweenPoints(p1, p2 point) float64 {
	xDiff := p1.X - p2.X
	yDiff := p1.Y - p2.Y
	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

func parsePoint(s string) (point, error) {
	s = strings.Trim(s, "()")
	coords := strings.Split(s, ",")
	if len(coords) != 2 {
		return point{}, fmt.Errorf("invalid number of coordinates: %v", s)
	}

	x, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return point{}, fmt.Errorf("invalid x coordinate: %v", coords[0])
	}

	y, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return point{}, fmt.Errorf("invalid y coordinate: %v", coords[1])
	}

	return point{x, y}, nil
}
