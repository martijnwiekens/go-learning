package main

import (
	"log"
	"time"

	"github.com/martijnwiekens/gointersection/intersection"
	"github.com/martijnwiekens/gointersection/intersectionapi"
	"github.com/martijnwiekens/gointersection/road"
	"github.com/martijnwiekens/gointersection/trafficcontroller"
	"github.com/martijnwiekens/gointersection/ui"
)

// const TRAFFIC_CONTROLLER_MODE = "INTEGRATED"
const TRAFFIC_CONTROLLER_MODE = "SEPERATED"
const TICK_SPEED = 4 * time.Second

func main() {
	// Remember the roads
	var roads []*road.Road

	// Create North Road
	var northRoad = road.NewRoad(road.NewRoadInput{Name: "NORTH", Left: 1, Right: 1, Forward: 2})
	roads = append(roads, northRoad)

	// Create South Road
	var southRoad = road.NewRoad(road.NewRoadInput{Name: "SOUTH", All: 2})
	roads = append(roads, southRoad)

	// Create East Road
	var eastRoad = road.NewRoad(road.NewRoadInput{Name: "EAST", Left: 1, All: 2})
	roads = append(roads, eastRoad)

	// Create West Road
	var westRoad = road.NewRoad(road.NewRoadInput{Name: "WEST", Left: 1, Right: 1, Forward: 2})
	roads = append(roads, westRoad)

	// Create Intersection
	in := intersection.NewIntersection(roads)

	// Create the traffic controller
	var tc *trafficcontroller.TrafficController
	if TRAFFIC_CONTROLLER_MODE == "INTEGRATED" {
		// Create traffic controller in the intersection
		tc = trafficcontroller.NewTrafficController(TRAFFIC_CONTROLLER_MODE, in)
	} else {
		// Start the API
		go intersectionapi.StartApi(in)

		// Create seperate traffic controller
		go trafficcontroller.StartTrafficControllerSeperated(TICK_SPEED)
	}

	// Remember the current tick
	currentTick := 0

	// Never ending loop
	for {
		// Log tick
		log.Default().Println("------", "Tick", currentTick, "------")

		// Tick the intersection
		in.Tick(currentTick)

		// Tick the traffic controller
		if TRAFFIC_CONTROLLER_MODE == "INTEGRATED" {
			tc.Tick(currentTick, 0)
		}

		// Print the intersection
		ui.PrintTick(currentTick)
		ui.PrintTotalCarsWaiting(in)
		ui.PrintIntersection(in)

		// Increment the tick
		currentTick++

		// Sleep for 1 second
		time.Sleep(TICK_SPEED)
	}
}
