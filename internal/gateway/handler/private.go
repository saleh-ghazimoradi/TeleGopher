package handler

import (
	"errors"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"net/http"
)

type PrivateHandler struct {
	privateService service.PrivateService
}

// CreatePrivate godoc
// @Summary      Create a private conversation
// @Description  Create a new private conversation between the authenticated user and another user
// @Tags         Private Conversations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        request body dto.PrivateRequest true "Receiver user ID"
// @Success      201 {object} helper.Response{data=dto.PrivateResponse} "Private conversation successfully created"
// @Failure      400 {object} helper.Response "Invalid request data or trying to create conversation with yourself"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      409 {object} helper.Response "Private conversation already exists"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /conversations/privates [post]
func (p *PrivateHandler) CreatePrivate(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	var payload dto.PrivateRequest
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	private, err := p.privateService.CreatePrivate(r.Context(), userId, payload.ReceiverId)
	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, repository.ErrSameUser):
			helper.BadRequestResponse(w, "Cannot create conversation with yourself", err)
		case errors.Is(err, repository.ErrPrivateAlreadyExists):
			helper.EditConflictResponse(w, "Private conversation already exists", err)
		default:
			helper.InternalServerError(w, "Failed to create private", err)
		}
		return
	}

	helper.CreatedResponse(w, "Private successfully created", private)
}

// GetPrivateById godoc
// @Summary      Get private conversation by ID
// @Description  Get a specific private conversation by its ID (user must be a participant)
// @Tags         Private Conversations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        id path int true "Private conversation ID"
// @Success      200 {object} helper.Response{data=dto.PrivateResponse} "Private conversation successfully retrieved"
// @Failure      400 {object} helper.Response "Invalid private conversation ID"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      403 {object} helper.Response "Forbidden - User doesn't have access to this conversation"
// @Failure      404 {object} helper.Response "Private conversation not found"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /conversations/privates/{id} [get]
func (p *PrivateHandler) GetPrivateById(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	privateId, err := helper.ReadParams(r)
	if err != nil {
		helper.BadRequestResponse(w, "Invalid private ID", err)
		return
	}

	private, err := p.privateService.GetPrivateById(r.Context(), privateId, userId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "Private conversation not found")
		case err.Error() == "unauthorized to access this private conversation":
			helper.ForbiddenResponse(w, "You don't have access to this conversation")
		default:
			helper.InternalServerError(w, "Failed to get private", err)
		}
		return
	}

	helper.SuccessResponse(w, "Private successfully retrieved", private)
}

// GetConversations godoc
// @Summary      Get user's private conversations
// @Description  Get all private conversations for the authenticated user
// @Tags         Private Conversations
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Success      200 {object} helper.Response{data=[]dto.PrivateResponse} "Conversations successfully retrieved"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      404 {object} helper.Response "User not found"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /conversations [get]
func (p *PrivateHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	privates, err := p.privateService.GetPrivatesForUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			helper.NotFoundResponse(w, "User not found")
			return
		}
		helper.InternalServerError(w, "Failed to get conversations", err)
		return
	}

	helper.SuccessResponse(w, "Conversations successfully retrieved", privates)
}

func NewPrivateHandler(privateService service.PrivateService) *PrivateHandler {
	return &PrivateHandler{
		privateService: privateService,
	}
}
