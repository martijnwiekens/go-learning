package trafficcontroller

import (
	"log"

	"github.com/gookit/event"
	"github.com/martijnwiekens/gointersection/collisonwarning"
	"github.com/martijnwiekens/gointersection/intersection"
	"github.com/martijnwiekens/gointersection/road"
)

type CurrentCall struct {
	roadName   string
	laneName   string
	greenTick  int
	orangeTick int
	redTick    int
}

type TrafficController struct {
	intersection        *intersection.Intersection
	mode                string
	currentPatternIndex int
	pendingCalls        [][2]string
	currentCalls        []*CurrentCall
}

const ORANGE_WAIT_TIME int = 10
const RED_WAIT_TIME int = 11

func NewTrafficController(mode string, in *intersection.Intersection) *TrafficController {
	// Create the traffic controller
	t := &TrafficController{mode: mode, intersection: in}

	// Register for events
	event.On("gointersection-road-traffic", event.ListenerFunc(func(e event.Event) error {
		log.Default().Println("RS: Lot of traffic on", e.Get("road"), ":", e.Get("lane"))
		t.OnRoadTraffic(e)
		return nil
	}), event.Normal)
	event.On("gointersection-road-empty", event.ListenerFunc(func(e event.Event) error {
		log.Default().Println("RS: Empty road at", e.Get("road"), ":", e.Get("lane"))
		t.OnRoadEmpty(e)
		return nil
	}), event.Normal)
	return t
}

func (t *TrafficController) Tick(currentTick int, try int) {
	// Stop the loop if we try more than the number of tries
	if try > len(TRAFFIC_PATTERN) {
		return
	}

	// Take control of the situation
	if currentTick == 0 {
		t.intersection.FullStopLights()
	}

	// Check if we have an current call
	if len(t.currentCalls) > 0 {
		// Check the stop call
		maxStopTick := 0
		for _, call := range t.currentCalls {
			// Check if we should make it ORANGE
			if currentTick == call.orangeTick {
				// Get the road and the lane
				t.SetLaneState(call, "ORANGE")
			}

			// Check if we should make it RED
			if currentTick == call.redTick {
				// Get the road and the lane
				t.SetLaneState(call, "RED")
			}

			// Check if this is the max stop tick
			if call.redTick > maxStopTick {
				maxStopTick = call.redTick
			}
		}

		// Check if we should stop
		if currentTick == maxStopTick {
			// Empty the current calls
			t.currentCalls = []*CurrentCall{}
		}
	}
	if len(t.currentCalls) == 0 {
		// Handle pending calls first
		if len(t.pendingCalls) > 0 {
			// Process the first call
			call := t.pendingCalls[0]
			t.pendingCalls = t.pendingCalls[1:]

			// Create a new call
			newCall := &CurrentCall{
				roadName:   call[0],
				laneName:   call[1],
				greenTick:  currentTick,
				orangeTick: currentTick + ORANGE_WAIT_TIME,
				redTick:    currentTick + RED_WAIT_TIME,
			}
			t.currentCalls = append(t.currentCalls, newCall)

			// Turn on the lights
			t.SetLaneState(newCall, "GREEN")
		} else {
			// Continue in the patterns
			// Go back to the first pattern
			if t.currentPatternIndex >= len(TRAFFIC_PATTERN) {
				t.currentPatternIndex = 0
			}

			// Get the pattern
			pattern := TRAFFIC_PATTERN[t.currentPatternIndex]

			// Make sure we go to the next pattern
			t.currentPatternIndex++

			// Remember if we enabled any lights
			hasNewCurrentCalls := false

			// Execute the pattern
			for _, call := range pattern {
				// Check if we have the road
				if !t.intersection.HasLane(call.roadName, call.laneName) {
					continue
				}
				if t.intersection.GetWaitingTrafficByLane(call.roadName, call.laneName) == 0 {
					continue
				}

				// Create the call
				newCall := &CurrentCall{
					roadName:   call.roadName,
					laneName:   call.laneName,
					greenTick:  currentTick,
					orangeTick: currentTick + ORANGE_WAIT_TIME,
					redTick:    currentTick + RED_WAIT_TIME,
				}
				t.currentCalls = append(t.currentCalls, newCall)

				// Turn on the lights
				t.SetLaneState(newCall, "GREEN")
				hasNewCurrentCalls = true
			}

			// Check if we have new calls
			if !hasNewCurrentCalls {
				// Try again
				try++
				t.Tick(currentTick, try)
			}
		}
	}
}

func (t *TrafficController) SetLaneState(call *CurrentCall, state string) {
	t.setLaneState(call.roadName, call.laneName, state)
}

func (t *TrafficController) setLaneState(roadName string, laneName string, state string) {
	// Log
	log.Default().Println("TC: Set lane state", roadName, ":", laneName, state)

	// Check for collision warning
	if state == "GREEN" {
		result := collisonwarning.CollisionWarningOnGreen(t.intersection, roadName, laneName)
		if result {
			log.Default().Fatal()
		}
	}

	// Check if we should set the crosswalk
	if roadName == "CROSSWALK" || roadName == "BICYCLE" {
		// Set the crosswalk
		var r *road.Road
		var roadNames []string = []string{"NORTH", "SOUTH", "EAST", "WEST"}
		for _, roadName := range roadNames {
			r = t.intersection.GetRoadByName(roadName)
			if r != nil {
				lanes := r.GetLanesByName(roadName)
				if len(lanes) > 0 {
					lanes[0].SetState(state)
				}
			}
		}

		// Check if we should set all the lanes (ALL, FORWARD, LEFT, RIGHT)
	} else {
		t.intersection.SetLights(roadName, laneName, state)
	}
}

func (t *TrafficController) OnRoadTraffic(e event.Event) {
	// Find the road
	roadName := e.Get("road").(string)
	laneName := e.Get("lane").(string)

	// Check if the lane is currently green
	if t.intersection.GetLaneState(roadName, laneName) == "GREEN" {
		return
	}

	// Check if we already planned the call
	isAlreadyPlanned := false
	for _, call := range t.currentCalls {
		if call.roadName == roadName && call.laneName == laneName {
			isAlreadyPlanned = true
			break
		}
	}

	// Plan the green call
	if !isAlreadyPlanned {
		t.pendingCalls = append(t.pendingCalls, [2]string{roadName, laneName})
		log.Default().Println("TC: Planned call for", roadName, ":", laneName)
	}
}

func (t *TrafficController) OnRoadEmpty(e event.Event) {
	// Find the road
	roadName := e.Get("road").(string)
	laneName := e.Get("lane").(string)

	// Check if other lanes have traffic
	if t.intersection.GetWaitingTrafficByLane(roadName, laneName) > 0 {
		// Don't do anything
		return
	}

	// Check if the lane is green
	if t.intersection.GetLaneState(roadName, laneName) != "GREEN" {
		// Don't need to do anything
		return
	}

	// Set the light to RED
	t.setLaneState(roadName, laneName, "ORANGE")

	// Check if we are in the current call
	for _, call := range t.currentCalls {
		if call.roadName == roadName && call.laneName == laneName {
			// Set the light to RED on next tick
			call.orangeTick = t.intersection.CurrentTick
			call.redTick = t.intersection.CurrentTick + 1
			break
		}
	}
}

type TrafficPattern struct {
	roadName string
	laneName string
}

var TRAFFIC_PATTERN [][]*TrafficPattern = [][]*TrafficPattern{
	{&TrafficPattern{roadName: "CROSSWALK"}, &TrafficPattern{roadName: "BICYCLE"}},
	{&TrafficPattern{roadName: "NORTH", laneName: "RIGHT"}, &TrafficPattern{roadName: "WEST", laneName: "RIGHT"}, &TrafficPattern{roadName: "WEST", laneName: "LEFT"}},
	{&TrafficPattern{roadName: "SOUTH", laneName: "ALL"}},
	{&TrafficPattern{roadName: "NORTH", laneName: "FORWARD"}, &TrafficPattern{roadName: "NORTH", laneName: "RIGHT"}, &TrafficPattern{roadName: "SOUTH", laneName: "FORWARD"}, &TrafficPattern{roadName: "SOUTH", laneName: "RIGHT"}},
	{&TrafficPattern{roadName: "EAST", laneName: "ALL"}},
	{&TrafficPattern{roadName: "WEST", laneName: "FORWARD"}, &TrafficPattern{roadName: "WEST", laneName: "RIGHT"}, &TrafficPattern{roadName: "EAST", laneName: "FORWARD"}, &TrafficPattern{roadName: "EAST", laneName: "RIGHT"}},
	{&TrafficPattern{roadName: "NORTH", laneName: "ALL"}},
	{&TrafficPattern{roadName: "NORTH", laneName: "LEFT"}, &TrafficPattern{roadName: "SOUTH", laneName: "RIGHT"}},
	{&TrafficPattern{roadName: "WEST", laneName: "ALL"}},
	{&TrafficPattern{roadName: "WEST", laneName: "LEFT"}, &TrafficPattern{roadName: "EAST", laneName: "LEFT"}},
}
