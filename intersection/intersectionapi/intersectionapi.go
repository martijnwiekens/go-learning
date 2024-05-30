package intersectionapi

import (
	"fmt"

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
	router.GET("/", outputIntersection)
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
	type outputData struct {
		RoadName string `json:"road"`
		LaneName string `json:"lane"`
		State    string `json:"state"`
		Traffic  int    `json:"traffic"`
	}
	road := GLOBAL_INTERSECTION.GetRoadByName(roadName)
	lane := road.GetLanesByName(laneName)[0]
	r := outputData{RoadName: roadName, LaneName: laneName, State: lane.GetState(), Traffic: lane.GetWaitingTrafficCount()}

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

func outputIntersection(c *gin.Context) {
	// This will output HTML
	output := "<html><head><title>GO Intersection</title></head><body>"
	output += "<h1>GO Intersection</h1>"
	output += "<p>Current tick: " + fmt.Sprint(GLOBAL_INTERSECTION.CurrentTick) + "</p>"
	output += "<p>Total cars waiting: " + fmt.Sprint(GLOBAL_INTERSECTION.GetWaitingTraffic()) + "</p>"
	output += "<table style='width: 100%'>"
	output += "<thead>"
	output += "<tr>"
	output += "<th>Road</th>"
	output += "<th>Lane</th>"
	output += "<th>Traffic</th>"
	output += "<th>Waiting since</th>"
	output += "<th>State</th>"
	output += "</tr>"
	output += "</thead>"
	output += "<tbody>"
	for _, road := range GLOBAL_INTERSECTION.GetRoads() {
		for _, lane := range road.GetLanes() {
			// Figure out if we want to see this line
			opacity := "1"
			if lane.GetWaitingTrafficCount() == 0 {
				opacity = "0.5"
			}

			// Create the table row
			output += "<tr style='opacity: " + opacity + "'>"
			output += "<td>" + road.GetName() + "</td>"
			output += "<td>" + lane.GetDirection() + "</td>"
			if lane.GetDirection() == "OUTPUT" {
				output += "<td />"
				output += "<td />"
			} else {
				// Add traffic
				output += "<td>" + fmt.Sprint(lane.GetWaitingTrafficCount()) + "</td>"

				// Add waiting since
				output += "<td>" + fmt.Sprint(lane.GetLongestWaitingTraffic()) + "</td>"

				// Find traffic light color
				lightColor := "orange"
				if lane.GetState() == "GREEN" {
					lightColor = "green"
				} else if lane.GetState() == "RED" {
					lightColor = "red"
				}

				// Add traffic light
				output += "<td style='color: " + lightColor + "'>" + lane.GetState() + "</td>"
			}
			output += "</tr>"
		}
	}
	output += "</tbody>"
	output += "</table>"
	output += "<script>setTimeout(() => { window.location.reload() }, 2000);</script>"
	output += "</body></html>"
	c.Data(200, "text/html; charset=utf-8", []byte(output))
}
