package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	kittpx "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

type CustomError interface {
	Error() string
	StatusCode() int
	Unwrap() error
	GRPCStatus() *status.Status
}

type NotFoundError struct {
	Err error
}

func NewNotFoundError(e error) NotFoundError {
	return NotFoundError{
		Err: e,
	}
}

func (e NotFoundError) Error() string {
	return e.Err.Error()
}

func (e NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

func (e NotFoundError) Unwrap() error {
	return e.Err
}

func (e NotFoundError) GRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

type BadRequestError struct {
	Err error
}

func NewBadRequestError(e error) BadRequestError {
	return BadRequestError{
		Err: e,
	}
}

func (e BadRequestError) Error() string {
	return e.Err.Error()
}

func (e BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}

func (e BadRequestError) Unwrap() error {
	return e.Err
}

func (e BadRequestError) GRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

type InternalServerError struct {
	Err error
}

func NewInternalServerError(e error) InternalServerError {
	return InternalServerError{
		Err: e,
	}
}

func (e InternalServerError) Error() string {
	return e.Err.Error()
}

func (e InternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}

func (e InternalServerError) Unwrap() error {
	return e.Err
}

func (e InternalServerError) GRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

type UnauthorizedError struct {
	Err error
}

func NewUnauthorizedError(e error) UnauthorizedError {
	return UnauthorizedError{
		Err: e,
	}
}

func (e UnauthorizedError) Error() string {
	return e.Err.Error()
}

func (e UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}

func (e UnauthorizedError) Unwrap() error {
	return e.Err
}

func (e UnauthorizedError) GRPCStatus() *status.Status {
	return status.New(codes.Unauthenticated, e.Error())
}

type NotImplementedError struct {
	Err error
}

func NewNotImplementedError(e error) NotImplementedError {
	return NotImplementedError{
		Err: e,
	}
}

func (e NotImplementedError) Error() string {
	return e.Err.Error()
}

func (e NotImplementedError) StatusCode() int {
	return http.StatusNotImplemented
}

func (e NotImplementedError) Unwrap() error {
	return e.Err
}

func (e NotImplementedError) GRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

type errorWrapper struct {
	ErrCode    int    `json:"error_code"`
	ErrMessage string `json:"error"`
}

func castByCode(code codes.Code, message string) error {
	switch code {
	case codes.NotFound:
		return NewNotFoundError(errors.New(message))
	case codes.InvalidArgument:
		return NewBadRequestError(errors.New(message))
	case codes.Unauthenticated:
		return NewUnauthorizedError(errors.New(message))
	case codes.Unimplemented:
		return NewNotImplementedError(errors.New(message))
	default:
		return NewInternalServerError(errors.New(message))
	}
}

func errToStatusCoder(err error) error {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		return err
	}

	return castByCode(grpcStatus.Code(), grpcStatus.Message())
}

// MakeHTTPErrorEncoder _
func MakeHTTPErrorEncoder(logger log.Logger) kittpx.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		_ = logger.Log("transport", "HTTP", "error", err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if headerer, ok := err.(kittpx.Headerer); ok {
			for k, values := range headerer.Headers() {
				for _, v := range values {
					w.Header().Add(k, v)
				}
			}
		}

		err = errToStatusCoder(err)

		code := http.StatusInternalServerError
		if sc, ok := err.(kittpx.StatusCoder); ok {
			code = sc.StatusCode()
		}
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(errorWrapper{
			ErrCode:    code,
			ErrMessage: err.Error(),
		})
	}
}
