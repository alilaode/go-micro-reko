package response

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	E_DUPLICATE            = "duplicate"
	E_NOT_FOUND            = "not_found"
	E_UNPROCESSABLE_ENTITY = "unprocessable_entity"
	E_UNAUTHORIZED         = "unauthorized"
	E_FORBIDDEN            = "forbidden"
	E_BAD_REQUEST          = "bad_request"
	E_SERVER_ERROR         = "server_error"
)

type errorConstant struct {
	Duplicate               Error
	NotFound                Error
	RouteNotFound           Error
	UnprocessableEntity     Error
	Unauthorized            Error
	BadRequest              Error
	Forbidden               Error
	Validation              Error
	MethodAllowedError      Error
	InternalServerError     Error
	ServiceUnavailableError Error
	NotFileUpload           Error
}

var (
	_validator    = validator.New()
	ErrorConstant = errorConstant{
		ServiceUnavailableError: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Service Unavailable",
				},
				Error: E_SERVER_ERROR,
			},
			Code: http.StatusServiceUnavailable,
		},

		MethodAllowedError: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Method Not Allowed",
				},
				Error: E_SERVER_ERROR,
			},
			Code: http.StatusMethodNotAllowed,
		},

		Duplicate: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Created value already exists",
				},
				Error: E_DUPLICATE,
			},
			Code: http.StatusConflict,
		},
		NotFound: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Data not found",
				},
				Error: E_NOT_FOUND,
			},
			Code: http.StatusNotFound,
		},

		RouteNotFound: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Route not found",
				},
				Error: E_NOT_FOUND,
			},
			Code: http.StatusNotFound,
		},
		UnprocessableEntity: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Invalid parameters or payload",
				},
				Error: E_UNPROCESSABLE_ENTITY,
			},
			Code: http.StatusUnprocessableEntity,
		},

		Unauthorized: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Unauthorized",
				},
				Error: E_UNAUTHORIZED,
			},
			Code: http.StatusUnauthorized,
		},
		Forbidden: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Forbidden access",
				},
				Error: E_FORBIDDEN,
			},
			Code: http.StatusForbidden,
		},
		BadRequest: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Bad Request",
				},
				Error: E_BAD_REQUEST,
			},
			Code: http.StatusBadRequest,
		},
		Validation: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Invalid parameters or payload",
				},
				Error: E_BAD_REQUEST,
			},
			Code: http.StatusBadRequest,
		},
		InternalServerError: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Something bad happened",
				},
				Error: E_SERVER_ERROR,
			},
			Code: http.StatusInternalServerError,
		},
		NotFileUpload: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "No files to upload",
				},
				Error: E_SERVER_ERROR,
			},
			Code: http.StatusInternalServerError,
		},
	}
)

type errorResponse struct {
	Meta  Meta        `json:"meta"`
	Error interface{} `json:"data"`
	//Description interface{} `json:"description,omitempty"`
}

type Error struct {
	Header       *http.Header
	Response     errorResponse `json:"response"`
	Code         int           `json:"code"`
	ErrorMessage error
}

func ErrorResponse(err error) *Error {
	re, ok := err.(*Error)
	if ok {
		return re
	} else {
		return ErrorBuilder(&ErrorConstant.InternalServerError, err)
	}
}

func ErrorBuilder(res *Error, message error, vals ...interface{}) *Error {
	res.ErrorMessage = message
	///res.Response.Description = vals

	return res
}

func CustomErrorBuilder(code int, err interface{}, message string, vals ...interface{}) *Error {
	return &Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: message,
			},
			Error: err,
			///Description: vals,
		},
		Code:         code,
		ErrorMessage: errors.New(message),
	}
}

func (e *Error) Error() string {
	if e.ErrorMessage == nil {
		e.ErrorMessage = errors.New(http.StatusText(e.Code))
	}
	return fmt.Sprintf("error code '%d' because: %s", e.Code, e.ErrorMessage.Error())
}

func (e *Error) Send(c echo.Context) error {
	var errorMessage string
	if e.ErrorMessage != nil {
		//apm.CaptureError(c.Request().Context(), e.ErrorMessage).Send()
		errorMessage = fmt.Sprintf("%+v", errors.WithStack(e.ErrorMessage))
	}
	logrus.Error(errorMessage)

	if e.Header != nil {
		for k, values := range *e.Header {
			for _, v := range values {
				c.Response().Header().Add(k, v)
			}
		}
	}

	return c.JSON(e.Code, e.Response)
}
