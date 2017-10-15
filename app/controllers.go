// Code generated by goagen v1.2.0-dirty, DO NOT EDIT.
//
// API "thermostat": Application Controllers
//
// Command:
// $ goagen
// --design=thermostat/design
// --out=$(GOPATH)/src/thermostat
// --version=v1.2.0-dirty

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Encoder.Register(goa.NewGobEncoder, "application/gob", "application/x-gob")
	service.Encoder.Register(goa.NewXMLEncoder, "application/xml")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")
	service.Decoder.Register(goa.NewGobDecoder, "application/gob", "application/x-gob")
	service.Decoder.Register(goa.NewXMLDecoder, "application/xml")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// OperandsController is the controller interface for the Operands actions.
type OperandsController interface {
	goa.Muxer
	Status(*StatusOperandsContext) error
	TargetHeatingCoolingState(*TargetHeatingCoolingStateOperandsContext) error
	TargetRelativeHumidity(*TargetRelativeHumidityOperandsContext) error
	TargetTemperature(*TargetTemperatureOperandsContext) error
}

// MountOperandsController "mounts" a Operands resource controller on the given service.
func MountOperandsController(service *goa.Service, ctrl OperandsController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewStatusOperandsContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Status(rctx)
	}
	service.Mux.Handle("GET", "/status", ctrl.MuxHandler("status", h, nil))
	service.LogInfo("mount", "ctrl", "Operands", "action", "Status", "route", "GET /status")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewTargetHeatingCoolingStateOperandsContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.TargetHeatingCoolingState(rctx)
	}
	service.Mux.Handle("GET", "/targetHeatingCoolingState/:value", ctrl.MuxHandler("targetHeatingCoolingState", h, nil))
	service.LogInfo("mount", "ctrl", "Operands", "action", "TargetHeatingCoolingState", "route", "GET /targetHeatingCoolingState/:value")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewTargetRelativeHumidityOperandsContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.TargetRelativeHumidity(rctx)
	}
	service.Mux.Handle("GET", "/targetRelativeHumidity/:value", ctrl.MuxHandler("targetRelativeHumidity", h, nil))
	service.LogInfo("mount", "ctrl", "Operands", "action", "TargetRelativeHumidity", "route", "GET /targetRelativeHumidity/:value")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewTargetTemperatureOperandsContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.TargetTemperature(rctx)
	}
	service.Mux.Handle("GET", "/targetTemperature/:value", ctrl.MuxHandler("targetTemperature", h, nil))
	service.LogInfo("mount", "ctrl", "Operands", "action", "TargetTemperature", "route", "GET /targetTemperature/:value")
}
