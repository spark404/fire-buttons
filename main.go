package main

import "github.com/dordille/artnet"
import . "github.com/splace/joysticks"
import "fmt"
import "time"
import "log"
import "sync"

var buddyNode *artnet.Node

var buttonState1 bool
var buttonState2 bool
var buttonState3 bool
var buttonState4 bool
var buttonState5 bool
var buttonState6 bool

var universe uint8 = 0x20

var wg sync.WaitGroup

func main() {
	fmt.Println("Press button the joystick to get started")

	device := Connect(1)

	if device == nil {
		panic("no HIDs")
	}
	fmt.Printf("HID#1:- Buttons:%d, Hats:%d\n", len(device.Buttons), len(device.HatAxes)/2)

	buttonState1 = false
	buttonState2 = false
	buttonState3 = false
	buttonState4 = false
	buttonState5 = false
	buttonState6 = false

	// make channels for specific events
	b1press := device.OnClose(1)
	b2press := device.OnClose(2)
	b3press := device.OnClose(3)
	b4press := device.OnClose(4)
	b5press := device.OnClose(5)
	b6press := device.OnClose(6)
	b1release := device.OnOpen(1)
	b2release := device.OnOpen(2)
	b3release := device.OnOpen(3)
	b4release := device.OnOpen(4)
	b5release := device.OnOpen(5)
	b6release := device.OnOpen(6)

	h1move := device.OnMove(1)

	// feed OS events onto the event channels.
	go device.ParcelOutEvents()

	err, buddyNode := artnet.NewNode("2.49.8.143:6454")
	if err != nil {
		log.Fatal(err)
	}

	// handle event channels
	go func() {
		for {
			select {
			case <-b1press:
				fmt.Println("button #1 pressed")
				buttonState1 = true
			case <-b2press:
				fmt.Println("button #2 pressed")
				buttonState2 = true
			case <-b3press:
				fmt.Println("button #3 pressed")
				buttonState3 = true
			case <-b4press:
				fmt.Println("button #4 pressed")
				buttonState4 = true
			case <-b5press:
				fmt.Println("button #5 pressed")
				buttonState5 = true
			case <-b6press:
				fmt.Println("button #6 pressed")
				buttonState6 = true
			case <-b1release:
				fmt.Println("button #1 released")
				buttonState1 = false
			case <-b2release:
				fmt.Println("button #2 released")
				buttonState2 = false
			case <-b3release:
				fmt.Println("button #3 released")
				buttonState3 = false
			case <-b4release:
				fmt.Println("button #4 released")
				buttonState4 = false
			case <-b5release:
				fmt.Println("button #5 released")
				buttonState5 = false
			case <-b6release:
				fmt.Println("button #6 released")
				buttonState6 = false
			case h := <-h1move:
				hpos := h.(CoordsEvent)
				fmt.Println("hat #1 moved too:", hpos.X, hpos.Y)
			}
		}
	}()

	ticker := time.NewTicker(30 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				var dmxData [512]uint8
				if buttonState1 {
					dmxData[402] = 0xff
				}
				if buttonState2 {
					dmxData[403] = 0xff
				}
				if buttonState3 {
					dmxData[404] = 0xff
				}
				if buttonState4 {
					dmxData[405] = 0xff
				}
				if buttonState5 {
					dmxData[402] = 0xff
					dmxData[403] = 0xff
					dmxData[404] = 0xff
					dmxData[405] = 0xff
				}
                                if buttonState6 {
					dmxData[180] = 0xff;
					dmxData[182] = 0x7f;
				}

				buddyNode.Dmx(universe, dmxData)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	fmt.Println("Blocking for ever")
	select {}
}
