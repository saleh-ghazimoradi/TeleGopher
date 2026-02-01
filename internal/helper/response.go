package helper

import (
	"fmt"
	"log/slog"
	"net/http"
)

type ErrResponse struct {
	logger *slog.Logger
}

func (e *ErrResponse) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	e.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (e *ErrResponse) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}

	if err := WriteJSON(w, status, env, nil); err != nil {
		e.LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (e *ErrResponse) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.LogError(r, err)
	message := "the server encountered a problem and could not process your request"
	e.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (e *ErrResponse) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	e.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (e *ErrResponse) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	e.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (e *ErrResponse) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (e *ErrResponse) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	e.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (e *ErrResponse) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	e.ErrorResponse(w, r, http.StatusConflict, message)
}

func (e *ErrResponse) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded, please try again"
	e.ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func (e *ErrResponse) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	e.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (e *ErrResponse) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	e.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func NewErrResponse(logger *slog.Logger) *ErrResponse {
	return &ErrResponse{logger: logger}
}
