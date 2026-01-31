package helper

import "net/http"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

type PaginatedResponse struct {
	Response Response      `json:"response"`
	Meta     PaginatedMeta `json:"meta"`
}

type PaginatedMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

func SuccessResponse(w http.ResponseWriter, message string, data any, headers http.Header) error {
	return writeJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	}, headers)
}

func CreatedResponse(w http.ResponseWriter, message string, data any, headers http.Header) error {
	return writeJSON(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	}, headers)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) error {
	response := Response{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return writeJSON(w, statusCode, response, http.Header{})
}

func BadRequestResponse(w http.ResponseWriter, message string, err error) error {
	return ErrorResponse(w, http.StatusBadRequest, message, err)
}

func NotFoundResponse(w http.ResponseWriter, message string) error {
	return ErrorResponse(w, http.StatusNotFound, message, nil)
}

func InternalServerResponse(w http.ResponseWriter, message string, err error) error {
	return ErrorResponse(w, http.StatusInternalServerError, message, err)
}

func UnauthorizedResponse(w http.ResponseWriter, message string) error {
	return ErrorResponse(w, http.StatusUnauthorized, message, nil)
}

func ForbiddenResponse(w http.ResponseWriter, message string) error {
	return ErrorResponse(w, http.StatusForbidden, message, nil)
}

func PaginatedSuccessResponse(w http.ResponseWriter, message string, data any, meta PaginatedMeta, header http.Header) error {
	return writeJSON(w, http.StatusOK, PaginatedResponse{
		Response: Response{
			Success: true,
			Message: message,
			Data:    data,
		},
		Meta: meta,
	}, header)
}
