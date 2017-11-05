//go:generate goagen bootstrap -d thermostat/design

package main

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/yuu/thermostat/app"
)

var status = Status{
	CurrentHeatingCoolingState: 0,
	CurrentRelativeHumidity:    50,
	CurrentTemperature:         21,
	TargetHeatingCoolingState:  0,
	TargetRelativeHumidity:     50,
	TargetTemperature:          21,
}

func main() {
	ir, _ := CreateIR()
	ir.OpenDevice()

	// Create service
	service := goa.New("thermostat")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "operands" controller
	c := NewOperandsController(service, status, *ir)
	app.MountOperandsController(service, c)

	// Start service
	if err := service.ListenAndServe(":9999"); err != nil {
		service.LogError("startup", "err", err)
	}
}
