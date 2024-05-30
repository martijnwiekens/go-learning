package ui

import (
	"fmt"
	"strings"

	"github.com/martijnwiekens/go-learning/gointersection/intersection"
	"github.com/martijnwiekens/go-learning/gointersection/road"
)

func PrintTick(currentTick int) {
	fmt.Println("=========================================")
	fmt.Println("Tick: ", currentTick)
}

func PrintTotalCarsWaiting(i *intersection.Intersection) {
	var totalCarsWaiting int
	for _, r := range i.GetRoads() {
		for _, l := range r.GetLanes() {
			totalCarsWaiting += l.GetWaitingTrafficCount()
		}
	}
	fmt.Println("Total cars waiting: ", totalCarsWaiting)
}

func PrintIntersection(i *intersection.Intersection) {
	/**
		2 lanes (west, east)
	    ---------|        |--------
	      [2] L  |        | [2] F
	    ---------|        |--------
	      [1] F  |        | [1] F
	    ---------|        |--------

		3 lanes (north, east, west)
		        |     |
		        | [1] |
				|  A  |
				|     |
		--------|-----|--------
	     [2] L  |     | [2] F
	    --------|     |--------
	     [1] F  |     | [1] F
	    --------|     |--------
	*/

	// Empty lines to frame the intersection
	fmt.Println()

	// Figure out how many roads we have
	var hasNorth, hasEast, hasSouth, hasWest bool
	var northLanes, eastLanes, southLanes, westLanes uint8
	var northRoad, eastRoad, southRoad, westRoad *road.Road
	for _, r := range i.GetRoads() {
		if r.GetName() == "NORTH" {
			hasNorth = true
			northRoad = r
			northLanes = r.GetLanesCount()
		} else if r.GetName() == "EAST" {
			hasEast = true
			eastRoad = r
			eastLanes = r.GetLanesCount()
		} else if r.GetName() == "SOUTH" {
			hasSouth = true
			southRoad = r
			southLanes = r.GetLanesCount()
		} else if r.GetName() == "WEST" {
			hasWest = true
			westRoad = r
			westLanes = r.GetLanesCount()
		}
	}

	// Figure out the spaces on the left
	var leftSpaces int
	if hasWest {
		leftSpaces = 9
	}

	// Check if we have a north road
	if hasNorth {
		// Get the road
		northOutput := getNorthRoadOutput(northRoad)

		// Print the north road
		for i := 0; i < 4; i++ {
			fmt.Print(strings.Repeat(" ", leftSpaces))
			fmt.Println(northOutput[i])
		}
	}

	// Create west road output
	var eastRoadOutput, westRoadOutput []string
	if hasWest {
		westRoadOutput = getWestRoadOutput(westRoad)
	}
	if hasEast {
		eastRoadOutput = getEastRoadOutput(eastRoad)
	}

	// Create a intersection box
	maxVerticalLanes := int(max(northLanes, southLanes))
	maxHorizontalLanes := max(int(max(eastLanes, westLanes)), 1)
	for i := 0; i < maxHorizontalLanes*4; i++ {
		// Print west lane
		if hasWest && len(westRoadOutput) > i {
			fmt.Print(westRoadOutput[i])
		} else {
			fmt.Print(strings.Repeat(" ", leftSpaces))
		}
		fmt.Print("|")

		// Print inner box
		if i == 0 || i == (maxHorizontalLanes*4)-1 {
			fmt.Print(strings.Repeat("------", int(maxVerticalLanes)))
		} else {
			fmt.Print(strings.Repeat("      ", int(maxVerticalLanes)))
		}

		// Print east lanes
		fmt.Print("|")
		if hasEast && len(eastRoadOutput) > i {
			fmt.Print(eastRoadOutput[i])
		} else {
			fmt.Print(strings.Repeat(" ", leftSpaces))
		}
		fmt.Println()
	}

	// Check if we have an south road
	if hasSouth {
		// Get the road
		southOutput := getSouthRoadOutput(southRoad)

		// Print the south road
		for i := 0; i < 4; i++ {
			fmt.Print(strings.Repeat(" ", leftSpaces))
			fmt.Println(southOutput[i])
		}
	}

}

func getNorthRoadOutput(r *road.Road) []string {
	var output []string
	for i := 0; i < 10; i++ {
		var outputLine string
		outputLine += "|"
		lanes := r.GetLanes()
		for il := len(lanes) - 1; il > -1; il-- {
			l := lanes[il]
			if i == 1 {
				// Output traffic count
				if l.GetDirection() == "OUTPUT" {
					outputLine += "  ^  "
				} else {
					outputLine += fmt.Sprintf(" [%d] ", l.GetWaitingTrafficCount())
				}
			} else if i == 2 && l.GetDirection() != "OUTPUT" {
				// Output direction
				outputLine += fmt.Sprintf(" (%v) ", string(l.GetDirection()[0]))
			} else if i == 3 && l.GetDirection() != "OUTPUT" {
				// Traffic light
				trafficLight := getTrafficLight(l)
				outputLine += fmt.Sprintf("  %v  ", trafficLight)
			} else {
				outputLine += "     "
			}
			outputLine += "|"
		}
		output = append(output, outputLine)
	}
	return output
}

func getSouthRoadOutput(r *road.Road) []string {
	var output []string
	for i := 0; i < 10; i++ {
		var outputLine string
		outputLine += "|"
		for _, l := range r.GetLanes() {
			if i == 0 && l.GetDirection() != "OUTPUT" {
				// Traffic light
				trafficLight := getTrafficLight(l)
				outputLine += fmt.Sprintf("  %v  ", trafficLight)
			} else if i == 1 && l.GetDirection() != "OUTPUT" {
				// Output direction
				outputLine += fmt.Sprintf(" (%v) ", string(l.GetDirection()[0]))

			} else if i == 2 {
				// Output traffic count
				if l.GetDirection() == "OUTPUT" {
					outputLine += "  V  "
				} else {
					outputLine += fmt.Sprintf(" [%d] ", l.GetWaitingTrafficCount())
				}
			} else {
				outputLine += "     "
			}
			outputLine += "|"
		}
		output = append(output, outputLine)
	}
	return output
}

func getWestRoadOutput(r *road.Road) []string {
	var output []string
	output = append(output, "---------")
	lanes := r.GetLanes()
	for i := 0; i < len(lanes); i++ {
		l := lanes[i]
		trafficLight := getTrafficLight(l)
		output = append(output, "         ")
		if l.GetDirection() == "OUTPUT" {
			output = append(output, "   <     ")
		} else {
			output = append(output, fmt.Sprintf(" [%d](%v)%v ", l.GetWaitingTrafficCount(), string(l.GetDirection()[0]), trafficLight))
		}
		output = append(output, "         ")
		output = append(output, "---------")
	}
	return output
}

func getEastRoadOutput(r *road.Road) []string {
	var output []string
	output = append(output, "---------")
	lanes := r.GetLanes()
	for i := len(lanes) - 1; i > -1; i-- {
		l := lanes[i]
		trafficLight := getTrafficLight(l)
		output = append(output, "         ")
		if l.GetDirection() == "OUTPUT" {
			output = append(output, "     >   ")
		} else {
			output = append(output, fmt.Sprintf(" %v(%v)[%d] ", trafficLight, string(l.GetDirection()[0]), l.GetWaitingTrafficCount()))
		}
		output = append(output, "         ")
		output = append(output, "---------")
	}
	return output
}

func getTrafficLight(l *road.Lane) string {
	trafficLight := "F"
	if l.GetState() == "GREEN" {
		trafficLight = "O"
	} else if l.GetState() == "ORANGE" {
		trafficLight = "E"
	} else if l.GetState() == "RED" {
		trafficLight = "X"
	}
	return trafficLight
}
