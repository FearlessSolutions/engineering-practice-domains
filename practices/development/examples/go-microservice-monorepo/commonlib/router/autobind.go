package router

import (
	"reflect"

	"example.com/sample/commonlib/logger"
	"example.com/sample/commonlib/request"
	"example.com/sample/commonlib/response"
	"github.com/jellydator/validation"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AutoBindAndValidate can wrap a controller function to accept an additional "requestBody" parameter, which will be
// automatically deserialized. It will also be automatically validated if it implements validation.Validatable.
func AutoBindAndValidate[BodyType any](next func(ctx echo.Context, requestBody BodyType) error) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var bodyTarget BodyType

		// Deserialize the body
		if deserializeErr := request.Bind.BindBody(ctx, &bodyTarget); deserializeErr != nil {
			logger.Log.Error("Received malformed JSON.", zap.Error(deserializeErr), zap.String("deserializedType", reflect.TypeOf(bodyTarget).String()))
			return response.BadRequest(deserializeErr).Respond(ctx)
		}

		// If the deserialized body implements `validation.Validatable`, validate it. To type-assert here,
		// we need to cast to "any", then back into validation.Validatable.
		var maybeValidatable any = bodyTarget
		if validatable, isValidatable := maybeValidatable.(validation.Validatable); isValidatable {
			if validationErr := validatable.Validate(); validationErr != nil {
				logger.Log.Warn("Received invalid data.", zap.Error(validationErr), zap.String("validatedType", reflect.TypeOf(bodyTarget).String()))
				return response.BadRequest(validationErr).Respond(ctx)
			}
		}

		return next(ctx, bodyTarget)
	}
}
