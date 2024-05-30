package traffic

type Car struct {
	FirstTick int
}

type Human struct {
	FirstTick int
}

type Bycycle struct {
	FirstTick int
}

func (c *Car) GetFirstTick() int {
	return c.FirstTick
}

func (h *Human) GetFirstTick() int {
	return h.FirstTick
}

func (b *Bycycle) GetFirstTick() int {
	return b.FirstTick
}

func (c *Car) CrossRoad() {
	// Cross the road
}

func (h *Human) CrossRoad() {
	// Cross the road
}

func (b *Bycycle) CrossRoad() {
	// Cross the road
}
