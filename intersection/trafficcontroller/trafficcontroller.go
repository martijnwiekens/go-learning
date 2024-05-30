package trafficcontroller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gookit/event"
	"github.com/martijnwiekens/gointersection/collisonwarning"
	"github.com/martijnwiekens/gointersection/intersection"
)

type CurrentCall struct {
	roadName   string
	laneName   string
	greenTick  int
	orangeTick int
	redTick    int
}

type TrafficController struct {
	intersection        IntersectionBridge
	currentPatternIndex int
	pendingCalls        [][2]string
	currentCalls        []*CurrentCall
}

const ORANGE_WAIT_TIME int = 10
const RED_WAIT_TIME int = 11

func NewTrafficController(mode string, in *intersection.Intersection) *TrafficController {
	// Find the right intersection connection
	var intersectionConnection IntersectionBridge
	if mode == "INTEGRATED" {
		intersectionConnection = &IntersectionDirectConnection{in: in}
	} else {
		intersectionConnection = &IntersectionApiConnection{}
	}
	// Create the traffic controller
	t := &TrafficController{intersection: intersectionConnection}

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
			removedItems := 0
			for index, call := range t.pendingCalls {
				// Get data
				roadName := call[0]
				laneName := call[1]

				// Check for collison warning
				if t.intersection.CollisionWarningOnGreen(roadName, laneName) {
					continue
				}

				// Remove the pending call from the list
				t.pendingCalls = append(t.pendingCalls[:(index-removedItems)], t.pendingCalls[(index-removedItems)+1:]...)
				removedItems++

				// Create a new call
				newCall := &CurrentCall{
					roadName:   roadName,
					laneName:   laneName,
					greenTick:  currentTick,
					orangeTick: currentTick + ORANGE_WAIT_TIME,
					redTick:    currentTick + RED_WAIT_TIME,
				}
				t.currentCalls = append(t.currentCalls, newCall)

				// Turn on the lights
				t.SetLaneState(newCall, "GREEN")
			}
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
		result := t.intersection.CollisionWarningOnGreen(roadName, laneName)
		if result {
			log.Default().Fatal()
		}
	}

	// Check if we should set the crosswalk
	if roadName == "CROSSWALK" || roadName == "BICYCLE" {
		// Set the crosswalk
		var roadNames []string = []string{"NORTH", "SOUTH", "EAST", "WEST"}
		for _, roadName := range roadNames {
			t.intersection.SetLightState(roadName, roadName, state)
		}

		// Check if we should set all the lanes (ALL, FORWARD, LEFT, RIGHT)
	} else {
		t.intersection.SetLightState(roadName, laneName, state)
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
			call.orangeTick = t.intersection.GetCurrentTick()
			call.redTick = t.intersection.GetCurrentTick() + 1
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

type IntersectionBridge interface {
	GetCurrentTick() int
	SetCurrentTick(tick int)
	FullStopLights()
	CollisionWarningOnGreen(roadName string, laneName string) bool
	HasLane(roadName string, laneName string) bool
	GetWaitingTrafficByLane(roadName string, laneName string) int
	SetLightState(roadName string, laneName string, state string) bool
	GetLaneState(roadName string, laneName string) string
}

type IntersectionDirectConnection struct {
	in *intersection.Intersection
}

type IntersectionApiConnection struct {
	currentTick int
}

func (ic *IntersectionDirectConnection) GetCurrentTick() int {
	return ic.in.CurrentTick
}

func (ic *IntersectionDirectConnection) SetCurrentTick(tick int) {
	// Can't do it here
}

func (ic *IntersectionDirectConnection) FullStopLights() {
	ic.in.FullStopLights()
}

func (ic *IntersectionDirectConnection) CollisionWarningOnGreen(roadName string, laneName string) bool {
	return collisonwarning.CollisionWarningOnGreen(ic.in, roadName, laneName)
}

func (ic *IntersectionDirectConnection) HasLane(roadName string, laneName string) bool {
	return ic.in.HasLane(roadName, laneName)
}

func (ic *IntersectionDirectConnection) GetWaitingTrafficByLane(roadName string, laneName string) int {
	return ic.in.GetWaitingTrafficByLane(roadName, laneName)
}

func (ic *IntersectionDirectConnection) SetLightState(roadName string, laneName string, state string) bool {
	return ic.in.SetLights(roadName, laneName, state)
}

func (ic *IntersectionDirectConnection) GetLaneState(roadName string, laneName string) string {
	return ic.in.GetLaneState(roadName, laneName)
}

func (ic *IntersectionApiConnection) FullStopLights() {
	http.Post("http://localhost:8080/stop", "application/json", nil)
}

func (ic *IntersectionApiConnection) CollisionWarningOnGreen(roadName string, laneName string) bool {
	result, err := http.Get("http://localhost:8080/road/lane/state?road=" + roadName + "&lane=" + laneName)
	if err != nil {
		return false
	}
	var data bool
	json.NewDecoder(result.Body).Decode(&data)
	return data
}

func (ic *IntersectionApiConnection) HasLane(roadName string, laneName string) bool {
	result, err := http.Get("http://localhost:8080/road/lane?road=" + roadName + "&lane=" + laneName)
	if err != nil {
		return false
	}
	if result.StatusCode == 200 {
		return true
	}
	return false
}

func (ic *IntersectionApiConnection) GetWaitingTrafficByLane(roadName string, laneName string) int {
	result, err := http.Get("http://localhost:8080/road/lane?road=" + roadName + "&lane=" + laneName)
	if err != nil {
		return 0
	}
	type resultDataType struct {
		Traffic int `json:"traffic"`
	}
	var data resultDataType
	json.NewDecoder(result.Body).Decode(&data)
	return data.Traffic
}

func (ic *IntersectionApiConnection) SetLightState(roadName string, laneName string, state string) bool {
	type inputDataType struct {
		Road  string `json:"road"`
		Lane  string `json:"lane"`
		State string `json:"state"`
	}
	inputData := inputDataType{Road: roadName, Lane: laneName, State: state}
	jsonData, _ := json.Marshal(inputData)
	result, err := http.Post(
		"http://localhost:8080/road/lane/state",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return false
	}
	if result.StatusCode == 200 {
		return true
	}
	return false
}

func (ic *IntersectionApiConnection) GetLaneState(roadName string, laneName string) string {
	result, err := http.Get("http://localhost:8080/road/lane?road=" + roadName + "&lane=" + laneName)
	if err != nil {
		return "FLASH"
	}
	type resultDataType struct {
		State string `json:"state"`
	}
	var data resultDataType
	json.NewDecoder(result.Body).Decode(&data)
	return data.State
}

func (ic *IntersectionApiConnection) GetCurrentTick() int {
	return ic.currentTick
}

func (ic *IntersectionApiConnection) SetCurrentTick(tick int) {
	ic.currentTick = tick
}

func StartTrafficControllerSeperated(tickSpeed time.Duration) {
	// Create the TrafficController
	tc := NewTrafficController("SEPERATED", nil)

	// Create the loop
	currentTick := 0
	for {
		// Set the current tick in the intersection connection
		tc.intersection.SetCurrentTick(currentTick)

		// Update the TrafficController
		tc.Tick(currentTick, 0)

		// Increase the tick
		currentTick++

		// Wait a little
		time.Sleep(tickSpeed)
	}
}
