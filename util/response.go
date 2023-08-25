package util

import (
	"net/http"
)

type Response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func InternalServerError(message *string) Response {
	if message == nil {
		m := "internal server error"
		message = &m
	}
	return Response{
		Message:    *message,
		StatusCode: http.StatusInternalServerError,
	}
}

func BadRequest(message *string) Response {
	if message == nil {
		m := "bad request"
		message = &m
	}
	return Response{
		Message:    *message,
		StatusCode: http.StatusBadRequest,
	}
}

func Ok(message *string) Response {
	if message == nil {
		m := "ok"
		message = &m
	}
	return Response{
		Message:    *message,
		StatusCode: http.StatusOK,
	}
}
