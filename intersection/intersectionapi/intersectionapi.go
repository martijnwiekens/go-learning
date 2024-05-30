package intersectionapi

import (
	"github.com/gin-gonic/gin"
	"github.com/martijnwiekens/gointersection/collisonwarning"
	"github.com/martijnwiekens/gointersection/intersection"
)

var GLOBAL_INTERSECTION *intersection.Intersection = nil

func StartApi(in *intersection.Intersection) {
	// Save the intersection
	GLOBAL_INTERSECTION = in

	// Create the API
	router := gin.Default()
	router.POST("/stop", setFullStop)
	router.GET("/road/lane", getLane)
	router.GET("/road/lane/state", validateNewLaneState)
	router.POST("/road/lane/state", setNewLaneState)
	router.Run("localhost:8080")
}

func getLane(c *gin.Context) {
	// Retrieve details
	roadName := c.Query("road")
	laneName := c.Query("lane")

	// Check if lane exists
	if !GLOBAL_INTERSECTION.HasLane(roadName, laneName) {
		c.AbortWithStatus(404)
		return
	}

	// Build the JSON
	type result struct {
		roadName string
		laneName string
		state    string
		traffic  int
	}
	road := GLOBAL_INTERSECTION.GetRoadByName(roadName)
	lane := road.GetLanesByName(laneName)[0]
	r := result{roadName: roadName, laneName: laneName, state: lane.GetState(), traffic: lane.GetWaitingTrafficCount()}

	// Return the JSON
	c.JSON(200, r)
}

func setNewLaneState(c *gin.Context) {
	// Retrieve details
	type inputData struct {
		Road  string `json:"road"`
		Lane  string `json:"lane"`
		State string `json:"state"`
	}
	var input inputData

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&input); err != nil {
		return
	}
	roadName := input.Road
	laneName := input.Lane
	state := input.State

	// Check if lane exists
	if !GLOBAL_INTERSECTION.HasLane(roadName, laneName) {
		c.AbortWithStatus(404)
		return
	}

	// Set the state
	GLOBAL_INTERSECTION.SetLights(roadName, laneName, state)
	c.AbortWithStatus(200)
}

func validateNewLaneState(c *gin.Context) {
	// Retrieve details
	roadName := c.Query("road")
	laneName := c.Query("lane")

	// Check if lane exists
	if !GLOBAL_INTERSECTION.HasLane(roadName, laneName) {
		c.AbortWithStatus(404)
		return
	}

	// Build the JSON
	collision := collisonwarning.CollisionWarningOnGreen(GLOBAL_INTERSECTION, roadName, laneName)

	// Return the result
	c.JSON(200, collision)
}

func setFullStop(c *gin.Context) {
	GLOBAL_INTERSECTION.FullStopLights()
	c.AbortWithStatus(200)
}
