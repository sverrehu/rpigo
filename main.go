package main

import (
	"time"

	"github.com/sverrehu/rpigo/gpio"
)

func main() {
	g := gpio.NewGPIO()
	g.SetMode(gpio.Bcm)
	g.Setup(18, gpio.Out)
	println("LED on")
	g.Output(18, gpio.High)
	time.Sleep(time.Duration(1) * time.Second)
	println("LED off")
	g.Output(18, gpio.Low)
}
