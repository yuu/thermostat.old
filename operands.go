package main

import (
	"github.com/goadesign/goa"
	"github.com/yuu/thermostat/app"
	"fmt"
	"math"
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
	res := &app.JSON{
		CurrentHeatingCoolingState: c.status.CurrentHeatingCoolingState,
		CurrentRelativeHumidity:    c.status.CurrentRelativeHumidity,
		CurrentTemperature:         c.status.CurrentTemperature,
		TargetHeatingCoolingState:  c.status.TargetHeatingCoolingState,
		TargetRelativeHumidity:     c.status.TargetRelativeHumidity,
		TargetTemperature:          c.status.TargetTemperature,
	}

	return ctx.OK(res)
}

// TargetHeatingCoolingState runs the targetHeatingCoolingState action.
func (c *OperandsController) TargetHeatingCoolingState(ctx *app.TargetHeatingCoolingStateOperandsContext) error {
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

	return ctx.OK(nil)
}

// TargetRelativeHumidity runs the targetRelativeHumidity action.
func (c *OperandsController) TargetRelativeHumidity(ctx *app.TargetRelativeHumidityOperandsContext) error {
	fmt.Println("not implment")

	return nil
}

// TargetTemperature runs the targetTemperature action.
func (c *OperandsController) TargetTemperature(ctx *app.TargetTemperatureOperandsContext) error {
	current := c.status.TargetTemperature
	future  := ctx.Value
	numbers := int(math.Abs(float64(future - current)))

	var data []byte
	if current < future {
		data, _ = temp_upBytes()
	}

	if current > future {
		data, _ = temp_downBytes()
	}

	for index := 0; index <= numbers; index++ {
		c.ir.Write(IR_FREQ_DEFAULT, data)
	}

	c.status.TargetTemperature = ctx.Value
	c.status.CurrentTemperature = ctx.Value

	return ctx.OK(nil)
}
