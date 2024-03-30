package controller

import "github.com/jellydator/validation"

// SampleGreetingRequest is the request body for the greeting endpoint on SampleController.
type SampleGreetingRequest struct {
	Name string `json:"name" validate:"required" example:"Xavier"`
}

// Validate implements validation.Validatable for SampleGreetingRequest. It validates the content
// of the request.
func (greetReq SampleGreetingRequest) Validate() error {
	return validation.ValidateStruct(&greetReq,
		validation.Field(&greetReq.Name, validation.Required),
	)
}

// SampleGreetingResponse is the response for the greeting endpoint on SampleController
type SampleGreetingResponse struct {
	Greeting string `json:"greeting" validate:"required" example:"Hello Xavier"`
}

// NewGreetingRequest is the request body for the "add greeting" endpoint on SampleController
type NewGreetingRequest struct {
	Greeting string `json:"greeting" validate:"required" example:"Hello"`
}

// Validate implements validation.Validatable for NewGreetingRequest.
func (greetReq NewGreetingRequest) Validate() error {
	return validation.ValidateStruct(&greetReq,
		validation.Field(&greetReq.Greeting, validation.Required, validation.Length(0, 32)),
	)
}
