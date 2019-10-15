package render

import (
	"fmt"
	"testing"

	v "github.com/go-ozzo/ozzo-validation"
)

var (
	errString      = fmt.Errorf("cruel error")
	validationError = v.Errors{
		"someError": errString,
	}

	JSONErrors = []error{
		validationError,
	}
	stringErrors = []error{
		errString,
	}
)

func TestErrors(t *testing.T) {
	for _, err := range JSONErrors {
		t.Run("Testing JSON errors", testJSONError(err))
	}

	for _, err := range stringErrors {
		t.Run("Testing String errors", testStringError(err))
	}
}

func testJSONError(err error) func(t *testing.T) {
	return func(t *testing.T) {
		errResp := HTTPBadRequest(err)
		if len(errResp.DetailsJSON) <= 0 {
			t.Errorf("DetailsJSON is empty HTTPBadRequest - %+v", err)
		}

		errResp = HTTPConflictError(err)
		if len(errResp.DetailsJSON) <= 0 {
			t.Errorf("DetailsJSON is empty HTTPConflictError - %+v", err)
		}

		errResp = HTTPInternalServerError(err)
		if errResp.Details == "" {
			t.Errorf("Details is empty HTTPInternalServerError - %+v", err)
		}

		errResp = HTTPServiceUnavailableError(err)
		if errResp.Details == "" {
			t.Errorf("Details is empty HTTPServiceUnavailableError - %+v", err)
		}
	}
}

func testStringError(err error) func(t *testing.T) {
	return func(t *testing.T) {
		errResp := HTTPBadRequest(err)
		if errResp.Details == "" || len(errResp.DetailsJSON) > 0 {
			t.Errorf("DetailsJSON is not empty or Details is empty HTTPBadRequest - %+v", err)
		}

		errResp = HTTPConflictError(err)
		if errResp.Details == "" || len(errResp.DetailsJSON) > 0 {
			t.Errorf("DetailsJSON is not empty or Details is empty HTTPConflictError - %+v", err)
		}

		errResp = HTTPInternalServerError(err)
		if errResp.Details == "" {
			t.Errorf("Details is empty HTTPInternalServerError - %+v", err)
		}

		errResp = HTTPServiceUnavailableError(err)
		if errResp.Details == "" {
			t.Errorf("Details is empty HTTPServiceUnavailableError - %+v", err)
		}
	}
}
