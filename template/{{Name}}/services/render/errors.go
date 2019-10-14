package render

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Status      int             `json:"-"`
	Error       string          `json:"error"`
	Details     string          `json:"details,omitempty"`
	DetailsJSON json.RawMessage `json:"detailsJSON,omitempty"`
}

func (er *ErrorResponse) SetDetails(err error) {
	jErr, ok := err.(json.Marshaler)

	if !ok {
		er.Details = err.Error()
		return
	}

	data, e := jErr.MarshalJSON()
	if e != nil {
		er.Details = err.Error()
		return
	}

	er.DetailsJSON = data
}

var (
	// HTTPNotFound prepared NotFound response
	HTTPNotFound = &ErrorResponse{
		Status: http.StatusNotFound,
		Error:  http.StatusText(http.StatusNotFound),
	}

	// HTTPMethodNotAllowed prepared MethodNotAllowed response
	HTTPMethodNotAllowed = &ErrorResponse{
		Status: http.StatusMethodNotAllowed,
		Error:  http.StatusText(http.StatusMethodNotAllowed),
	}

	// HTTPConflict prepared Conflict response
	HTTPConflict = &ErrorResponse{
		Status: http.StatusConflict,
		Error:  http.StatusText(http.StatusConflict),
	}

	// HTTPForbidden prepared Forbidden response
	HTTPForbidden = &ErrorResponse{
		Status: http.StatusForbidden,
		Error:  http.StatusText(http.StatusForbidden),
	}

	// HTTPUnauthorized prepared Unauthorized response
	HTTPUnauthorized = &ErrorResponse{
		Status: http.StatusUnauthorized,
		Error:  http.StatusText(http.StatusUnauthorized),
	}
)

// HTTPBadRequest prepared BadRequest response
func HTTPBadRequest(err error) *ErrorResponse {
	response := &ErrorResponse{
		Status: http.StatusBadRequest,
		Error:  http.StatusText(http.StatusBadRequest),
	}

	response.SetDetails(err)

	return response
}

// HTTPInternalServerError prepared InternalServerError response
func HTTPInternalServerError(err error) *ErrorResponse {
	log.Error(err.Error())
	return &ErrorResponse{
		Status:  http.StatusInternalServerError,
		Error:   http.StatusText(http.StatusInternalServerError),
		Details: err.Error(),
	}
}

// HTTPConflictError prepared ConflictError response
func HTTPConflictError(err error) *ErrorResponse {
	response := &ErrorResponse{
		Status: http.StatusConflict,
		Error:  http.StatusText(http.StatusConflict),
	}

	response.SetDetails(err)

	return response
}

// HTTPServiceUnavailableError prepared HTTPServiceUnavailableError response
func HTTPServiceUnavailableError(err error) *ErrorResponse {
	log.Error(err.Error())
	return &ErrorResponse{
		Status:  http.StatusServiceUnavailable,
		Error:   http.StatusText(http.StatusServiceUnavailable),
		Details: err.Error(),
	}
}
