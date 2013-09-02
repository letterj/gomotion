// The package gomotion defines a concurrent Go library that can connect to a Leap motion device over a WebSocket conection.
// By default, the LeapMotion exposes a JSON WebSocket that pumps out messages near 30 to 50 fps.
package gomotion

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

// The LeapMotionDevice definition. Connecting to a device will return an instance of this struct.
type LeapMotionDevice struct {
	Pipe       chan *Frame
	Connection *websocket.Conn
}

// This function acts as a constructor and connector for the gomotion package.
func GetDevice(url string) *LeapMotionDevice {
	pipe := make(chan *Frame)
	connection, err := websocket.Dial(url, "", "http://localhost")
	if err != nil {
		log.Fatal(err)
	}
	return &LeapMotionDevice{pipe, connection}
}

// This function starts the listening on the WebSocket. By default it enables Gestures on the LeapMotionDevice.
func (device *LeapMotionDevice) Listen() {
	var config struct {
		enableGestures bool `json:"enableGestures"`
	}
	config.enableGestures = true
	err := websocket.JSON.Send(device.Connection, &config)
	if err != nil {
		log.Fatal(err)
	}
	go device.listenRead()
}

func (device *LeapMotionDevice) listenRead() {
	for {
		var frame Frame
		err := websocket.JSON.Receive(device.Connection, &frame)
		if err == nil {
			device.Pipe <- &frame
		} else {
			log.Fatal(err)
		}
	}
}

// This function closes the internal WebSocket connection on a LeapMotionDevice
func (device *LeapMotionDevice) Close() {
	device.Connection.Close()
}