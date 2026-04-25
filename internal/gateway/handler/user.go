package handler

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

// GetUserById godoc
// @Summary      Get user by ID
// @Description  Retrieve a user's information by their ID. User must be authenticated.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        id path int true "User ID"
// @Success      200 {object} helper.Response{data=dto.UserResponse} "User successfully retrieved"
// @Failure      400 {object} helper.Response "Invalid ID format"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      404 {object} helper.Response "User not found"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /users/{id} [get]
func (u *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadParams(r)
	if err != nil {
		helper.BadRequestResponse(w, "invalid id format", err)
		return
	}

	user, err := u.userService.GetUserById(r.Context(), id)
	if err != nil {
		helper.NotFoundResponse(w, "User not found")
		return
	}

	helper.SuccessResponse(w, "User successfully retrieved", user)
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
