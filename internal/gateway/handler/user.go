package handler

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService service.UserService
	errResponse *helper.ErrResponse
}

func (u *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	strId := r.PathValue("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		u.errResponse.BadRequestResponse(w, r, err)
		return
	}

	user, err := u.userService.GetUserById(r.Context(), id)
	if err != nil {
		u.errResponse.NotFoundResponse(w, r)
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"user": user}, nil); err != nil {
		u.errResponse.ServerErrorResponse(w, r, err)
	}
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
