package road

import (
	"log"
	"math/rand"

	"github.com/gookit/event"
	"github.com/martijnwiekens/gointersection/traffic"
)

type Road struct {
	name             string  // Location of the road in the intersection NORTH, SOUTH, EAST, WEST, OUTPUT
	lanes            []*Lane // Lanes of the road
	newCarSpeed      uint8   // How fast new cars are coming in the road
	newHumanSpeed    uint8   // How fast new humans are coming in the road
	newBycycleSpeed  uint8   // How fast new bycycles are coming in the road
	crossWalkEnabled bool    // Whether or not the road has a crosswalk
	bycyclesEnabled  bool    // Whether or not the road has bycycles
}

type Lane struct {
	road            string    // The road this lane belongs to
	waitingTraffic  []Traffic // Traffic that is waiting in the lane
	direction       string    // Direction of the lane LEFT, RIGHT, FORWARD, ALL
	state           string    // State of the lane RED, ORANGE, GREEN
	notifiedTraffic bool      // Whether or not the road has been notified of much traffic to the traffic controller
}

type Traffic interface {
	GetFirstTick() int
	CrossRoad()
}

type NewRoadInput struct {
	Name             string
	Left             uint8
	Right            uint8
	Forward          uint8
	All              uint8
	CrossWalkEnabled bool
	BycyclesEnabled  bool
	NewCarSpeed      uint8
	NewHumanSpeed    uint8
	NewBycycleSpeed  uint8
}

func NewRoad(input NewRoadInput) *Road {
	// Create the lanes of the road
	var lanes []*Lane

	// There is always an output lane
	lanes = append(lanes, &Lane{
		road:           input.Name,
		waitingTraffic: []Traffic{},
		direction:      "OUTPUT",
		state:          "FLASH",
	})

	// Create left lanes
	for i := uint8(0); i < input.Left; i++ {
		lanes = append(lanes, &Lane{
			road:           input.Name,
			waitingTraffic: []Traffic{},
			direction:      "LEFT",
			state:          "FLASH",
		})
	}

	// Create forward lanes
	for i := uint8(0); i < input.Forward; i++ {
		lanes = append(lanes, &Lane{
			road:           input.Name,
			waitingTraffic: []Traffic{},
			direction:      "FORWARD",
			state:          "FLASH",
		})
	}

	// Create lanes for all directions
	for i := uint8(0); i < input.All; i++ {
		lanes = append(lanes, &Lane{
			road:           input.Name,
			waitingTraffic: []Traffic{},
			direction:      "ALL",
			state:          "FLASH",
		})
	}

	// Create right lanes
	for i := uint8(0); i < input.Right; i++ {
		lanes = append(lanes, &Lane{
			road:           input.Name,
			waitingTraffic: []Traffic{},
			direction:      "RIGHT",
			state:          "FLASH",
		})
	}

	// Check if we need an cross walk
	if input.CrossWalkEnabled {
		lanes = append(lanes, &Lane{
			road:           input.Name,
			waitingTraffic: []Traffic{},
			direction:      "CROSSWALK",
			state:          "FLASH",
		})
	}
	if input.BycyclesEnabled {
		lanes = append(lanes, &Lane{
			road:           input.Name,
			waitingTraffic: []Traffic{},
			direction:      "BICYCLE",
			state:          "FLASH",
		})
	}

	// Check if we have speed
	if input.NewCarSpeed == 0 {
		input.NewCarSpeed = uint8(rand.Intn(15))
	}
	if input.NewHumanSpeed == 0 {
		input.NewHumanSpeed = uint8(rand.Intn(60))
	}
	if input.NewBycycleSpeed == 0 {
		input.NewBycycleSpeed = uint8(rand.Intn(15))
	}

	// Create the road
	r := &Road{
		name:             input.Name,
		lanes:            lanes,
		newCarSpeed:      max(5, input.NewCarSpeed),
		newHumanSpeed:    max(10, input.NewHumanSpeed),
		newBycycleSpeed:  max(5, input.NewBycycleSpeed),
		crossWalkEnabled: input.CrossWalkEnabled,
		bycyclesEnabled:  input.BycyclesEnabled,
	}

	// Handle empty road
	event.On("gointersection-road-empty", event.ListenerFunc(func(e event.Event) error {
		roadName := e.Get("road").(string)
		laneName := e.Get("lane").(string)
		if roadName == r.name {
			lanes := r.GetLanesByName(laneName)
			if len(lanes) > 0 {
				lanes[0].notifiedTraffic = false
			}
		}
		return nil
	}), event.Normal)

	return r
}

func (r *Road) GetName() string {
	return r.name
}

func (r *Road) GetLanes() []*Lane {
	return r.lanes
}

func (r *Road) GetLanesCount() uint8 {
	return uint8(len(r.lanes))
}

func (r *Road) GetLane(index uint8) *Lane {
	return r.lanes[index]
}

func (r *Road) GetLanesByName(name string) []*Lane {
	var lanes []*Lane
	for _, lane := range r.lanes {
		if lane.direction == name {
			lanes = append(lanes, lane)
		}
	}
	return lanes
}

func (r *Road) Tick(currentTick int) {
	// One loop
	if currentTick%int(r.newCarSpeed) == 0 {
		// Choose a random lane
		laneIndex := max(1, rand.Intn(int(r.GetLanesCount())))
		lane := r.GetLane(uint8(laneIndex))

		// Add new cars, there is a 10% chance of a double car
		var amountCars uint8 = 1
		if rand.Intn(5) == 0 {
			amountCars = uint8(rand.Intn(10))
		}
		log.Default().Println("R: +", amountCars, "cars incoming at", r.name, ":", lane.direction)
		for i := uint8(0); i < amountCars; i++ {
			lane.waitingTraffic = append(lane.waitingTraffic, &traffic.Car{FirstTick: currentTick})
		}
	}
	if r.crossWalkEnabled && currentTick%int(r.newHumanSpeed) == 0 {
		// Create a new human on crosswalk lane
		lanes := r.GetLanesByName("CROSSWALK")
		if len(lanes) > 0 {
			log.Default().Println("R: +1 human incoming at", r.name, ":CROSSWALK")
			lanes[0].waitingTraffic = append(lanes[0].waitingTraffic, &traffic.Human{FirstTick: currentTick})
		}
	}
	if r.bycyclesEnabled && currentTick%int(r.newBycycleSpeed) == 0 {
		// Create a new bycycle on bycycle lane
		lanes := r.GetLanesByName("BICYCLE")
		if len(lanes) > 0 {
			log.Default().Println("R: +1 bicycle incoming at", r.name, ":BICYCLE")
			lanes[0].waitingTraffic = append(lanes[0].waitingTraffic, &traffic.Bycycle{FirstTick: currentTick})
		}
	}

	// Check how many items in the lane
	for _, lane := range r.GetLanes() {
		if lane.GetWaitingTrafficCount() > 3 && !lane.notifiedTraffic {
			lane.notifiedTraffic = true
			event.MustFire("gointersection-road-traffic", event.M{"road": r.name, "lane": lane.direction})
			break
		}
	}

	// Tick each of the lanes
	for _, lane := range r.GetLanes() {
		lane.Tick()
	}
}

func (r *Road) GetLongestWaitingTraffic(laneName string) int {
	var longest int = 0
	lanes := r.GetLanesByName(laneName)
	for _, lane := range lanes {
		if lane.GetWaitingTrafficCount() > longest {
			longest = lane.GetWaitingTrafficCount()
		}
	}
	return longest
}

func (l *Lane) GetWaitingTrafficCount() int {
	return len(l.waitingTraffic)
}

func (l *Lane) GetDirection() string {
	return l.direction
}

func (l *Lane) GetState() string {
	return l.state
}

func (l *Lane) SetState(state string) {
	// Update state
	l.state = state
}

func (l *Lane) Tick() {
	if l.GetState() == "FLASH" || l.GetState() == "GREEN" {
		// Let out some traffic
		if len(l.waitingTraffic) > 0 {
			// There is a 20% chance of letting double cars out
			var amountCars uint8 = 1
			if rand.Intn(20) == 0 {
				amountCars = 2
			}
			log.Default().Println("R: -", amountCars, "cars leaving on", l.road, ":", l.direction)
			for i := uint8(0); i < amountCars; i++ {
				if len(l.waitingTraffic) > 0 {
					l.waitingTraffic[0].CrossRoad()
					l.waitingTraffic = l.waitingTraffic[1:]
				}
			}

			// Check if we have traffic
			if len(l.waitingTraffic) == 0 {
				event.MustFire("gointersection-road-empty", event.M{"road": l.road, "lane": l.direction})
			}
		}
	}
}

func (l *Lane) GetRoad() string {
	return l.road
}

func (l *Lane) GetLongestWaitingTraffic() int {
	if len(l.waitingTraffic) > 0 {
		return l.waitingTraffic[0].GetFirstTick()
	}
	return 0
}

func (l *Lane) GetNotified() bool {
	return l.notifiedTraffic
}
