package response

import (
	"net/http"

	"example.com/sample/commonlib/response/dtos"
)

// BadRequest creates a dtos.APIErrorHelper with a canned message for a 400 bad request.
func BadRequest(err error) *dtos.APIErrorHelper {
	return &dtos.APIErrorHelper{
		Status:      http.StatusBadRequest,
		Description: "Could not parse the submitted data.",
		Error:       err,
	}
}

// Conflict creates a dtos.APIErrorHelper with a canned message for a 409 conflict.
func Conflict(err error) *dtos.APIErrorHelper {
	return &dtos.APIErrorHelper{
		Status:      http.StatusConflict,
		Error:       err,
		Description: "Existing data in the system prevented the operation from completing.",
	}
}

// UnauthorizedWithErr creates a dtos.APIErrorHelper with a canned message for a 401 unauthorized response.
func UnauthorizedWithErr(err error) *dtos.APIErrorHelper {
	helper := Unauthorized()
	helper.Error = err
	return helper
}

// Unauthorized creates a dtos.APIErrorHelper with a canned message for a 401 unauthorized response.
// This function does not require an error, unlike UnauthorizedWithErr
func Unauthorized() *dtos.APIErrorHelper {
	return &dtos.APIErrorHelper{
		Status:      http.StatusUnauthorized,
		Description: "Failed to authorize you with the provided credentials.",
	}
}

// Forbidden creates a dtos.APIErrorHelper with a canned message for a 403 forbidden response.
func Forbidden(err error) *dtos.APIErrorHelper {
	return &dtos.APIErrorHelper{
		Status:      http.StatusForbidden,
		Error:       err,
		Description: "You do not have permission to access this data.",
	}
}

// NotFound creates a dtos.APIErrorHelper with a canned message for a 404 not found.
func NotFound() *dtos.APIErrorHelper {
	return &dtos.APIErrorHelper{
		Status:      http.StatusNotFound,
		Description: "Could not find the requested data.",
	}
}

// InternalServerError creates a dtos.APIErrorHelper with a canned message for a 500 internal server error.
func InternalServerError(err error) *dtos.APIErrorHelper {
	return &dtos.APIErrorHelper{
		Status:      http.StatusInternalServerError,
		Description: "Something went wrong while serving your request.",
		Error:       err,
	}
}
