package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("thermostat", func() {
	Title("The thermostat API")
	Description("Controll IR")
	Host("localhost:9999")
	Scheme("http")
})

var _ = Resource("operands", func() {
	Action("status", func() {
		Routing(GET("status"))
		Description("Get any thermostat info")
		Response(OK, "application/json")
	})
})
