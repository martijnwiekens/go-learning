package intersection

import (
	"github.com/martijnwiekens/go-learning/intersection/road"
)

type Intersection struct {
	CurrentTick int
	roads       []*road.Road
}

func NewIntersection(roads []*road.Road) *Intersection {
	return &Intersection{
		roads: roads,
	}
}

func (i *Intersection) GetRoads() []*road.Road {
	return i.roads
}

func (i *Intersection) GetRoadByName(name string) *road.Road {
	for _, road := range i.roads {
		if road.GetName() == name {
			return road
		}
	}
	return nil
}

func (i *Intersection) Tick(currentTick int) {
	// Remember the tick
	i.CurrentTick = currentTick

	// Tick each road
	for _, r := range i.roads {
		r.Tick(currentTick)
	}
}

func (i *Intersection) FullStopLights() {
	for _, r := range i.roads {
		for _, l := range r.GetLanes() {
			l.SetState("RED")
		}
	}
}

func (i *Intersection) DisableLights() {
	for _, r := range i.roads {
		for _, l := range r.GetLanes() {
			l.SetState("FLASH")
		}
	}
}

func (i *Intersection) EnableLights(roadName string, direction string) {
	road := i.GetRoadByName(roadName)
	if road != nil {
		lanes := road.GetLanesByName(direction)
		for _, l := range lanes {
			l.SetState("GREEN")
		}
	}
}

func (i *Intersection) SetLights(roadName string, direction string, state string) bool {
	// Use all the lanes
	lightsSet := false
	road := i.GetRoadByName(roadName)
	if road != nil {
		laneNames := []string{direction}
		if direction == "ALL" {
			laneNames = []string{"LEFT", "FORWARD", "RIGHT", "ALL"}
		}
		for _, laneName := range laneNames {
			lanes := road.GetLanesByName(laneName)
			for _, lane := range lanes {
				lane.SetState(state)
				lightsSet = true
			}
		}
	}
	return lightsSet
}

func (i *Intersection) GetWaitingTrafficByLane(roadName string, direction string) int {
	var waitingTraffic int = 0
	road := i.GetRoadByName(roadName)
	if road != nil {
		lanes := road.GetLanesByName(direction)
		for _, lane := range lanes {
			waitingTraffic += lane.GetWaitingTrafficCount()
		}
	}
	return waitingTraffic
}

func (i *Intersection) GetWaitingTraffic() int {
	var waitingTraffic int = 0
	for _, r := range i.roads {
		for _, l := range r.GetLanes() {
			waitingTraffic += l.GetWaitingTrafficCount()
		}
	}
	return waitingTraffic
}

func (i *Intersection) HasLane(roadName string, direction string) bool {
	road := i.GetRoadByName(roadName)
	if road != nil {
		lanes := road.GetLanesByName(direction)
		if len(lanes) > 0 {
			return true
		}
	}
	return false
}

func (i *Intersection) GetLaneState(roadName string, direction string) string {
	road := i.GetRoadByName(roadName)
	if road != nil {
		lanes := road.GetLanesByName(direction)
		if len(lanes) > 0 {
			return lanes[0].GetState()
		}
	}
	return "FLASH"
}
