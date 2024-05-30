package collisonwarning

import (
	"log"

	"github.com/martijnwiekens/gointersection/intersection"
	"github.com/martijnwiekens/gointersection/road"
)

func CollisionWarningOnGreen(in *intersection.Intersection, roadName string, laneName string) bool {
	// Create a list of lanes to check
	var checkLanes []*road.Lane

	// Find all the roads
	roadOpposite := in.GetRoadByName(GetOppositeDirection(roadName))
	roadLeft := in.GetRoadByName(GetDirectionLeft(roadName))
	roadRight := in.GetRoadByName(GetDirectionRight(roadName))

	// Check if we enable all lanes
	if laneName == "ALL" {
		// Nothing else should be green
		var checkRoads []*road.Road
		checkRoads = append(checkRoads, in.GetRoadByName(GetOppositeDirection(roadName)))
		checkRoads = append(checkRoads, in.GetRoadByName(GetDirectionLeft(roadName)))
		checkRoads = append(checkRoads, in.GetRoadByName(GetDirectionRight(roadName)))
		for _, road := range checkRoads {
			if road != nil {
				checkLanes = append(checkLanes, road.GetLanes()...)
			}
		}

		// Check if direction is LEFT
	} else if laneName == "LEFT" {
		if roadOpposite != nil {
			// FORWARD, RIGHT, but not LEFT
			checkLanes = append(checkLanes, roadOpposite.GetLanesByName("FORWARD")...)
			checkLanes = append(checkLanes, roadOpposite.GetLanesByName("RIGHT")...)
		}
		if roadLeft != nil {
			// FORWARD, LEFT roads, but not RIGHT
			checkLanes = append(checkLanes, roadLeft.GetLanesByName("FORWARD")...)
			checkLanes = append(checkLanes, roadLeft.GetLanesByName("LEFT")...)
		}
		if roadRight != nil {
			// FORWARD, LEFT roads, but not RIGHT
			checkLanes = append(checkLanes, roadLeft.GetLanesByName("FORWARD")...)
			checkLanes = append(checkLanes, roadLeft.GetLanesByName("LEFT")...)
		}

		// Check if direction is RIGHT
	} else if laneName == "RIGHT" {
		if roadOpposite != nil {
			// LEFT, but not FORWARD, RIGHT
			checkLanes = append(checkLanes, roadOpposite.GetLanesByName("LEFT")...)
		}
		if roadLeft != nil {
			// FORWARD roads, but not LEFT, RIGHT
			checkLanes = append(checkLanes, roadLeft.GetLanesByName("FORWARD")...)
		}

		// Check if direction is FORWARD
	} else if laneName == "FORWARD" {
		if roadOpposite != nil {
			// RIGHT, but not FORWARD, LEFT
			checkLanes = append(checkLanes, roadOpposite.GetLanesByName("LEFT")...)
		}
		if roadLeft != nil {
			// FORWARD, LEFT roads, but not RIGHT
			checkLanes = append(checkLanes, roadLeft.GetLanesByName("FORWARD")...)
			checkLanes = append(checkLanes, roadLeft.GetLanesByName("LEFT")...)
		}
		if roadRight != nil {
			// FORWARD, LEFT roads, but not RIGHT
			checkLanes = append(checkLanes, roadRight.GetLanesByName("FORWARD")...)
			checkLanes = append(checkLanes, roadRight.GetLanesByName("LEFT")...)
		}
	}

	// Check all the lanes
	for _, lane := range checkLanes {
		if lane.GetState() == "GREEN" || lane.GetState() == "ORANGE" {
			log.Default().Println(lane.GetRoad(), ":", lane.GetDirection(), "is still GREEN! COLLISION!")
			return true
		}
	}
	return false
}

func GetOppositeDirection(direction string) string {
	if direction == "NORTH" {
		return "SOUTH"
	} else if direction == "SOUTH" {
		return "NORTH"
	} else if direction == "EAST" {
		return "WEST"
	}
	return "EAST"
}

func GetDirectionLeft(direction string) string {
	if direction == "NORTH" {
		return "WEST"
	} else if direction == "SOUTH" {
		return "EAST"
	} else if direction == "EAST" {
		return "NORTH"
	}
	return "SOUTH"
}

func GetDirectionRight(direction string) string {
	if direction == "NORTH" {
		return "EAST"
	} else if direction == "SOUTH" {
		return "WEST"
	} else if direction == "EAST" {
		return "SOUTH"
	}
	return "NORTH"
}
