package handler

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"net/http"
)

type MessageHandler struct {
	messageService service.MessageService
}

// SendMessage godoc
// @Summary      Send a new message
// @Description  Send a new message in a private conversation
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        request body dto.MessageRequest true "Message details"
// @Success      201 {object} helper.Response{data=dto.MessageResponse} "Message successfully created"
// @Failure      400 {object} helper.Response "Invalid request data"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      404 {object} helper.Response "Private conversation not found"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /messages [post]
func (m *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	var payload dto.MessageRequest
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateMessageRequest(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "input's not valid")
		return
	}

	message, err := m.messageService.SendMessage(r.Context(), &payload, userId)
	if err != nil {
		helper.InternalServerError(w, "failed to send message", err)
		return
	}

	helper.CreatedResponse(w, "Message successfully created", message)
}

// GetMessage godoc
// @Summary      Get a specific message
// @Description  Get a specific message by ID
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        id path int true "Message ID"
// @Success      200 {object} helper.Response{data=dto.MessageResponse} "Message successfully retrieved"
// @Failure      400 {object} helper.Response "Invalid message ID"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      403 {object} helper.Response "Forbidden - User not authorized to view this message"
// @Failure      404 {object} helper.Response "Message not found"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /messages/{id} [get]
func (m *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	id, err := helper.ReadParams(r)
	if err != nil {
		helper.BadRequestResponse(w, "invalid id", err)
		return
	}

	message, err := m.messageService.GetMessage(r.Context(), id, userId)
	if err != nil {
		helper.InternalServerError(w, "failed to get message", err)
		return
	}

	helper.SuccessResponse(w, "Message successfully retrieved", message)
}

// GetPrivateMessages godoc
// @Summary      Get private conversation messages
// @Description  Get all messages from a specific private conversation with pagination
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        id path int true "Private conversation ID"
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(20) maximum(100)
// @Success      200 {object} helper.Response{data=dto.MessageListResponse} "Private messages successfully fetched"
// @Failure      400 {object} helper.Response "Invalid conversation ID or pagination parameters"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      403 {object} helper.Response "Forbidden - User not authorized to view this conversation"
// @Failure      404 {object} helper.Response "Private conversation not found"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /conversations/privates/{id}/messages [get]
func (m *MessageHandler) GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	id, err := helper.ReadParams(r)
	if err != nil {
		helper.BadRequestResponse(w, "invalid id", err)
		return
	}

	page, limit := helper.ParsePagination(r)

	messages, err := m.messageService.GetPrivateMessages(r.Context(), id, userId, page, limit)
	if err != nil {
		helper.InternalServerError(w, "failed to get messages", err)
		return
	}

	helper.SuccessResponse(w, "Private messages successfully fetched ", messages)
}

// GetUndeliveredMessages godoc
// @Summary      Get undelivered messages
// @Description  Get all undelivered messages from a specific private conversation
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        id path int true "Private conversation ID"
// @Success      200 {object} helper.Response{data=[]dto.MessageResponse} "Undelivered messages successfully fetched"
// @Failure      400 {object} helper.Response "Invalid conversation ID"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      403 {object} helper.Response "Forbidden - User not authorized to view this conversation"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /conversations/privates/{id}/messages/undelivered [get]
func (m *MessageHandler) GetUndeliveredMessages(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	id, err := helper.ReadParams(r)
	if err != nil {
		helper.BadRequestResponse(w, "invalid id", err)
		return
	}

	messages, err := m.messageService.GetUndeliveredMessages(r.Context(), id, userId)
	if err != nil {
		helper.InternalServerError(w, "failed to get messages", err)
		return
	}

	helper.SuccessResponse(w, "Undelivered messages successfully fetched", messages)
}

// MarkMessageAsRead godoc
// @Summary      Mark message as read
// @Description  Mark a specific message as read by the recipient
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        id path int true "Message ID"
// @Success      200 {object} helper.Response "Message successfully marked as read"
// @Failure      400 {object} helper.Response "Invalid message ID"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      403 {object} helper.Response "Forbidden - User not authorized to mark this message as read"
// @Failure      404 {object} helper.Response "Message not found"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /messages/{id}/read [patch]
func (m *MessageHandler) MarkMessageAsRead(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	id, err := helper.ReadParams(r)
	if err != nil {
		helper.BadRequestResponse(w, "invalid id", err)
		return
	}

	if err := m.messageService.MarkMessageAsRead(r.Context(), id, userId); err != nil {
		helper.InternalServerError(w, "failed to mark message as read", err)
		return
	}

	helper.SuccessResponse(w, "Message successfully marked as read", nil)
}

// MarkMessageAsDelivered godoc
// @Summary      Mark message as delivered
// @Description  Mark a specific message as delivered to the recipient
// @Tags         Messages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        id path int true "Message ID"
// @Success      200 {object} helper.Response "Message successfully marked as delivered"
// @Failure      400 {object} helper.Response "Invalid message ID"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      403 {object} helper.Response "Forbidden - User not authorized to mark this message as delivered"
// @Failure      404 {object} helper.Response "Message not found"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /messages/{id}/delivered [patch]
func (m *MessageHandler) MarkMessageAsDelivered(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	id, err := helper.ReadParams(r)
	if err != nil {
		helper.BadRequestResponse(w, "invalid id", err)
		return
	}

	if err := m.messageService.MarkMessageAsDelivered(r.Context(), id, userId); err != nil {
		helper.InternalServerError(w, "failed to mark message as delivered", err)
		return
	}

	helper.SuccessResponse(w, "Message successfully marked as delivered", nil)
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}
