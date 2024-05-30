package collisonwarning

import (
	"testing"

	"github.com/martijnwiekens/go-learning/gointersection/intersection"
	"github.com/martijnwiekens/go-learning/gointersection/road"
)

func prepareIntersection() *intersection.Intersection {
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
	return in
}

func TestGetOppositeDirection(t *testing.T) {
	var result string
	result = GetOppositeDirection("NORTH")
	if result != "SOUTH" {
		t.Fatalf(`North direction should be SOUTH, got %s`, result)
	}
	result = GetOppositeDirection("SOUTH")
	if result != "NORTH" {
		t.Fatalf(`South direction should be NORTH, got %s`, result)
	}
	result = GetOppositeDirection("EAST")
	if result != "WEST" {
		t.Fatalf(`East direction should be WEST, got %s`, result)
	}
	result = GetOppositeDirection("WEST")
	if result != "EAST" {
		t.Fatalf(`West direction should be EAST, got %s`, result)
	}
}

func TestGetDirectionLeft(t *testing.T) {
	var result string
	result = GetDirectionLeft("NORTH")
	if result != "WEST" {
		t.Fatalf(`North direction should be WEST, got %s`, result)
	}
	result = GetDirectionLeft("SOUTH")
	if result != "EAST" {
		t.Fatalf(`South direction should be EAST, got %s`, result)
	}
	result = GetDirectionLeft("EAST")
	if result != "NORTH" {
		t.Fatalf(`East direction should be NORTH, got %s`, result)
	}
	result = GetDirectionLeft("WEST")
	if result != "SOUTH" {
		t.Fatalf(`West direction should be SOUTH, got %s`, result)
	}
}

func TestGetDirectionRight(t *testing.T) {
	var result string
	result = GetDirectionRight("NORTH")
	if result != "EAST" {
		t.Fatalf(`North direction should be EAST, got %s`, result)
	}
	result = GetDirectionRight("SOUTH")
	if result != "WEST" {
		t.Fatalf(`South direction should be WEST, got %s`, result)
	}
	result = GetDirectionRight("EAST")
	if result != "SOUTH" {
		t.Fatalf(`East direction should be SOUTH, got %s`, result)
	}
	result = GetDirectionRight("WEST")
	if result != "NORTH" {
		t.Fatalf(`West direction should be NORTH, got %s`, result)
	}
}

func TestCollisionWarningOnGreenNorthAll(t *testing.T) {
	// Create an intersection
	in := prepareIntersection()

	// Everything should be RED
	in.FullStopLights()

	// Set NORTH ALL road to GREEN, this should not be a problem
	var result bool
	result = CollisionWarningOnGreen(in, "NORTH", "ALL")
	if result == true {
		t.Fatalf(`CollisionWarningOnGreen should be false, got %v`, result)
	}

	// Test if we can set NORTH ALL to GREEN, when SOUTH LEFT doesn't exist
	in.EnableLights("SOUTH", "LEFT")
	result = CollisionWarningOnGreen(in, "NORTH", "ALL")
	if result == true {
		t.Fatalf(`CollisionWarningOnGreen should be true, got %v`, result)
	}
	in.FullStopLights()

	// Test if we can set NORTH ALL to GREEN, when EAST LEFT
	in.EnableLights("EAST", "LEFT")
	result = CollisionWarningOnGreen(in, "NORTH", "ALL")
	if result == false {
		t.Fatalf(`CollisionWarningOnGreen should be true, got %v`, result)
	}
	in.FullStopLights()

	// Test if we can set NORTH ALL to GREEN, when SOUTH ALL
	in.EnableLights("SOUTH", "ALL")
	result = CollisionWarningOnGreen(in, "NORTH", "ALL")
	if result == false {
		t.Fatalf(`CollisionWarningOnGreen should be true, got %v`, result)
	}
	in.FullStopLights()

	// Test if we can set NORTH ALL to GREEN, when EAST RIGHT
	in.EnableLights("EAST", "LEFT")
	result = CollisionWarningOnGreen(in, "NORTH", "ALL")
	if result == false {
		t.Fatalf(`CollisionWarningOnGreen should be true, got %v`, result)
	}
	in.FullStopLights()

	// Test if we can set NORTH ALL to GREEN, when SOUTH FORWARD
	in.EnableLights("WEST", "FORWARD")
	result = CollisionWarningOnGreen(in, "NORTH", "ALL")
	if result == false {
		t.Fatalf(`CollisionWarningOnGreen should be true, got %v`, result)
	}
	in.FullStopLights()
}

func TestCollisionWarningOnGreenNorthLeft(t *testing.T) {
	// Create an intersection
	in := prepareIntersection()

	// Everything should be RED
	in.FullStopLights()

	// Set NORTH LEFT road to GREEN, this should not be a problem
	var result bool
	result = CollisionWarningOnGreen(in, "NORTH", "LEFT")
	if result == true {
		t.Fatalf(`CollisionWarningOnGreen should be false, got %v`, result)
	}

	// Test if we can set NORTH LEFT to GREEN, when SOUTH LEFT doesn't exist
	in.EnableLights("SOUTH", "LEFT")
	result = CollisionWarningOnGreen(in, "NORTH", "LEFT")
	if result == true {
		t.Fatalf(`CollisionWarningOnGreen should be true, got %v`, result)
	}
	in.FullStopLights()
}

func TestCollisionWarningOnGreenNorthRight(t *testing.T) {
	// Create an intersection
	in := prepareIntersection()

	// Everything should be RED
	in.FullStopLights()

	// Set NORTH RIGHT road to GREEN, this should not be a problem
	var result bool
	result = CollisionWarningOnGreen(in, "NORTH", "RIGHT")
	if result == true {
		t.Fatalf(`CollisionWarningOnGreen should be false, got %v`, result)
	}

	// Test if we can set NORTH RIGHT to GREEN, when SOUTH RIGHT doesn't exist
	in.EnableLights("SOUTH", "RIGHT")
	result = CollisionWarningOnGreen(in, "NORTH", "RIGHT")
	if result == true {
		t.Fatalf(`CollisionWarningOnGreen should be true, got %v`, result)
	}
	in.FullStopLights()
}

func TestCollisionWarningOnGreenNorthForward(t *testing.T) {
	// Create an intersection
	in := prepareIntersection()

	// Everything should be RED
	in.FullStopLights()

	// Set NORTH FORWARD road to GREEN, this should not be a problem
	var result bool
	result = CollisionWarningOnGreen(in, "NORTH", "FORWARD")
	if result == true {
		t.Fatalf(`CollisionWarningOnGreen should be false, got %v`, result)
	}

	// Test if we can set NORTH FORWARD to GREEN, when SOUTH FORWARD doesn't exist
	in.EnableLights("SOUTH", "FORWARD")
	result = CollisionWarningOnGreen(in, "NORTH", "FORWARD")
	if result == true {
		t.Fatalf(`CollisionWarningOnGreen should be true, got %v`, result)
	}
	in.FullStopLights()

	// Test if we can set NORTH FORWARD to GREEN, when SOUTH FORWARD doesn't exist
	in.EnableLights("NORTH", "FORWARD")
	in.EnableLights("NORTH", "RIGHT")
	result = CollisionWarningOnGreen(in, "SOUTH", "FORWARD")
	if result == true {
		t.Fatalf(`CollisionWarningOnGreen should be true, got %v`, result)
	}
	in.FullStopLights()
}
