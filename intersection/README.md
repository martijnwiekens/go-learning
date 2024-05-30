# Intersection
*Go Learning Project - by Martijn Wiekens*

Complicated road intersection with traffic light controller.
This project contains an intersection. The intersection has 4 roads, each with multiple lanes, each lane has random traffic and a traffic light.

The traffic controller has a default pattern through which it runs, each lane gets an opportunity for traffic to leave. Some lanes will be skipped if there is no traffic.

When the lane has a lot of traffic (more than 3 cars waiting) it will signal the TrafficController that a lot of traffic is waiting. The traffic controller will than override the default pattern to serve this lane first.

The traffic controller has a default amount of time how long it waits before cycling through orange and red. But when a lane is empty it can signal the traffic controller the lane is empty which speeds up this process.

The system runs on a loop, with each tick (single run) in 4 seconds.

**Traffic Controller**

The Traffic Controller has 2 modes, `INTEGRATED` and `SEPERATED`. 
In integrated mode the traffic controller runs in the same loop as the intersection and calls the intersection directly.

In seperated mode the traffic controller runs in a goroutine with an API. The Traffic Controller controls the intersection through the API and not directly. It also has its own loop, with the same duraction as the intersection. 
You can change this behavior in `TRAFFIC_CONTROLLER_MODE` in [main.go](main.go)

## Install
1. `go mode download`

## Run
1. `go run main.go`

If you run in `SEPERATED` mode, you can access [http://localhost:8080](http://localhost:8080) to see the status of the intersection in your browser. It refreshes automatically.

In `INTEGRATED` mode you will see a visual representation of the intersection in the CLI.
