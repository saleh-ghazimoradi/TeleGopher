package handler

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	errResponse    *helper.ErrResponse
	validator      *helper.Validator
	messageService service.MessageService
}

func (m *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		m.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	var payload dto.CreateMessageRequest
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		m.errResponse.BadRequestResponse(w, r, err)
		return
	}

	payload.Validate(m.validator)
	if !m.validator.Valid() {
		m.errResponse.FailedValidationResponse(w, r, m.validator.Errors)
		return
	}

	message, err := m.messageService.SendMessage(r.Context(), &payload, userId)
	if err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusCreated, helper.Envelope{"message": message}, nil); err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (m *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		m.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	messageIdstr := r.PathValue("id")
	messageId, err := strconv.ParseInt(messageIdstr, 10, 64)
	if err != nil {
		m.errResponse.BadRequestResponse(w, r, err)
		return
	}

	message, err := m.messageService.GetMessage(r.Context(), messageId, userId)
	if err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": message}, nil); err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (m *MessageHandler) GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		m.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	privateIdStr := r.PathValue("private_id")
	privateId, err := strconv.ParseInt(privateIdStr, 10, 64)
	if err != nil {
		m.errResponse.BadRequestResponse(w, r, err)
		return
	}

	page := 1
	limit := 20

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	messages, err := m.messageService.GetPrivateMessages(r.Context(), privateId, userId, page, limit)
	if err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"data": messages}, nil); err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (m *MessageHandler) MarkMessageAsRead(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		m.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	messageIdStr := r.PathValue("id")
	messageId, err := strconv.ParseInt(messageIdStr, 10, 64)
	if err != nil {
		m.errResponse.BadRequestResponse(w, r, err)
		return
	}

	if err := m.messageService.MarkMessageAsRead(r.Context(), messageId, userId); err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "Message marked as read"}, nil); err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (m *MessageHandler) MarkMessageAsDelivered(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		m.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	messageIdStr := r.PathValue("id")
	messageId, err := strconv.ParseInt(messageIdStr, 10, 64)
	if err != nil {
		m.errResponse.BadRequestResponse(w, r, err)
		return
	}

	if err := m.messageService.MarkMessageAsDelivered(r.Context(), messageId, userId); err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "Message marked as delivered"}, nil); err != nil {
		m.errResponse.ServerErrorResponse(w, r, err)
	}
}

func NewMessageHandler(errResponse *helper.ErrResponse, validator *helper.Validator, messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		errResponse:    errResponse,
		validator:      validator,
		messageService: messageService,
	}
}
