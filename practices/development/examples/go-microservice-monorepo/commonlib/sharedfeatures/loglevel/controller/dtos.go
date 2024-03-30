package controller

import "github.com/jellydator/validation"

// ChangeLogLevelRequest is the format of the request body for adjusting the log level
type ChangeLogLevelRequest struct {
	NewLevel string
}

// Validate implements validation.Validatable for ChangeLogLevelRequest
func (req ChangeLogLevelRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.NewLevel,
			validation.Required,
			validation.In("debug", "info", "warn", "error", "panic", "fatal").
				Error("must be one of debug, info, warn, error, panic, or fatal"),
		),
	)
}
