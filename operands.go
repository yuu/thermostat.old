package main

import (
	"github.com/goadesign/goa"
	"thermostat/app"
)

type Status struct {
	CurrentHeatingCoolingState int
	CurrentRelativeHumidity    int
	CurrentTemperature         int
	TargetHeatingCoolingState  int
	TargetRelativeHumidity     int
	TargetTemperature          int
}

// OperandsController implements the operands resource.
type OperandsController struct {
	*goa.Controller
	status Status
	ir     IR
}

// NewOperandsController creates a operands controller.
func NewOperandsController(service *goa.Service, s Status, i IR) *OperandsController {
	return &OperandsController{Controller: service.NewController("OperandsController"), status: s, ir: i}
}

// Status runs the status action.
func (c *OperandsController) Status(ctx *app.StatusOperandsContext) error {
	// OperandsController_Status: start_implement

	// Put your logic here
	res := &app.JSON{
		CurrentHeatingCoolingState: c.status.CurrentHeatingCoolingState,
		CurrentRelativeHumidity:    c.status.CurrentRelativeHumidity,
		CurrentTemperature:         c.status.CurrentTemperature,
		TargetHeatingCoolingState:  c.status.TargetHeatingCoolingState,
		TargetRelativeHumidity:     c.status.TargetRelativeHumidity,
		TargetTemperature:          c.status.TargetTemperature,
	}

	// OperandsController_Status: end_implement
	return ctx.OK(res)
}

// TargetHeatingCoolingState runs the targetHeatingCoolingState action.
func (c *OperandsController) TargetHeatingCoolingState(ctx *app.TargetHeatingCoolingStateOperandsContext) error {
	// OperandsController_TargetHeatingCoolingState: start_implement

	// Put your logic here

	// OperandsController_TargetHeatingCoolingState: end_implement
	return nil
}

// TargetRelativeHumidity runs the targetRelativeHumidity action.
func (c *OperandsController) TargetRelativeHumidity(ctx *app.TargetRelativeHumidityOperandsContext) error {
	// OperandsController_TargetRelativeHumidity: start_implement

	// Put your logic here

	// OperandsController_TargetRelativeHumidity: end_implement
	return nil
}

// TargetTemperature runs the targetTemperature action.
func (c *OperandsController) TargetTemperature(ctx *app.TargetTemperatureOperandsContext) error {
	// OperandsController_TargetTemperature: start_implement

	// Put your logic here

	// OperandsController_TargetTemperature: end_implement
	return nil
}
