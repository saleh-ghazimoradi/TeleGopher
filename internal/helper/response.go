package helper

import "net/http"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

type PaginatedMeta struct {
	Page      int64 `json:"page"`
	Limit     int64 `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
}

type PaginatedResponse struct {
	Response Response
	Meta     PaginatedMeta `json:"meta"`
}

func SuccessResponse(w http.ResponseWriter, message string, data any) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	writeJSON(w, http.StatusOK, response)
}

func CreatedResponse(w http.ResponseWriter, message string, data any) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	writeJSON(w, http.StatusCreated, response)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	writeJSON(w, statusCode, response)
}

func BadRequestResponse(w http.ResponseWriter, message string, err error) {
	ErrorResponse(w, http.StatusBadRequest, message, err)
}

func UnauthorizedResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusUnauthorized, message, nil)
}

func ForbiddenResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusForbidden, message, nil)
}

func NotFoundResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusNotFound, message, nil)
}

func InternalServerError(w http.ResponseWriter, message string, err error) {
	ErrorResponse(w, http.StatusInternalServerError, message, err)
}

func FailedValidationResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusUnprocessableEntity, message, nil)
}

func EditConflictResponse(w http.ResponseWriter, message string, err error) {
	ErrorResponse(w, http.StatusConflict, message, err)
}

func RateLimitExceededResponse(w http.ResponseWriter, message string) {
	ErrorResponse(w, http.StatusTooManyRequests, message, nil)
}

func HTTPRouterNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusNotFound, "Not found", nil)
}

func HTTPRouterMethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
}

func PaginatedSuccessResponse(w http.ResponseWriter, message string, data any, meta PaginatedMeta) {
	paginatedResponse := PaginatedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		},
		Meta: meta,
	}
	writeJSON(w, http.StatusOK, paginatedResponse)
}
