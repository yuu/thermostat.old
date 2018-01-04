package main

import (
	"github.com/goadesign/goa"
	"github.com/yuu/thermostat/app"
	"fmt"
)

const (
	MODE_OFF  = 0
	MODE_HEAT = 1
	MODE_COOL = 2
	MODE_AUTO = 3
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
	var need_boot bool = c.status.CurrentHeatingCoolingState == MODE_OFF
	var mode []byte
	switch ctx.Value {
	case MODE_OFF:
		mode, _ = Asset("off")
		c.status.CurrentHeatingCoolingState = 0
	case MODE_HEAT:
		mode, _ = modeHeatBytes()
		c.status.CurrentHeatingCoolingState = 1
	case MODE_COOL:
		mode, _ = modeCoolBytes()
		c.status.CurrentHeatingCoolingState = 2
	case MODE_AUTO:
		mode, _ = modeHumidityBytes()
		c.status.CurrentHeatingCoolingState = 2
	default:
		return fmt.Errorf("not supported")
	}

	if need_boot {
		on, _ := Asset("on")
		c.ir.Write(IR_FREQ_DEFAULT, on)
	}

	c.ir.Write(IR_FREQ_DEFAULT, mode)
	c.status.TargetHeatingCoolingState = ctx.Value

	// OperandsController_TargetHeatingCoolingState: end_implement
	return nil
}

// TargetRelativeHumidity runs the targetRelativeHumidity action.
func (c *OperandsController) TargetRelativeHumidity(ctx *app.TargetRelativeHumidityOperandsContext) error {
	// OperandsController_TargetRelativeHumidity: start_implement

	// Put your logic here
	fmt.Println("not implment")

	// OperandsController_TargetRelativeHumidity: end_implement
	return nil
}

// TargetTemperature runs the targetTemperature action.
func (c *OperandsController) TargetTemperature(ctx *app.TargetTemperatureOperandsContext) error {
	ct := c.status.TargetTemperature
	tt := ctx.Status.TargetTemperature

	for index := tt - ct {
		c.ir.Write(IR_FREQ_DEFAULT, temp_up())
	}

	c.status.TargetTemperature = ctx.Status.TargetTemperature

	return ctx.OK(nil)
}
