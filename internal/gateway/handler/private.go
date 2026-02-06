package handler

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"net/http"
	"strconv"
)

type PrivateHandler struct {
	privateService service.PrivateService
	errResponse    *helper.ErrResponse
}

func (p *PrivateHandler) GetPrivate(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		p.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	privateIdStr := r.PathValue("private_id")
	privateId, err := strconv.ParseInt(privateIdStr, 10, 64)
	if err != nil {
		p.errResponse.BadRequestResponse(w, r, err)
		return
	}

	private, err := p.privateService.GetPrivate(r.Context(), privateId, userId)
	if err != nil {
		p.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"private": private}, nil); err != nil {
		p.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (p *PrivateHandler) CreatePrivate(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		p.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	var req struct {
		ReceiverId int64 `json:"receiver_id"`
	}

	if err := helper.ReadJSON(w, r, &req); err != nil {
		p.errResponse.BadRequestResponse(w, r, err)
		return
	}

	private, err := p.privateService.CreatePrivate(r.Context(), userId, req.ReceiverId)
	if err != nil {
		p.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusCreated, helper.Envelope{"private": private}, nil); err != nil {
		p.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (p *PrivateHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		p.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	privates, err := p.privateService.GetUserPrivates(r.Context(), userId)
	if err != nil {
		p.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"privates": privates}, nil); err != nil {
		p.errResponse.ServerErrorResponse(w, r, err)
	}
}

func NewPrivateHandler(errResponse *helper.ErrResponse, privateService service.PrivateService) *PrivateHandler {
	return &PrivateHandler{
		errResponse:    errResponse,
		privateService: privateService,
	}
}
