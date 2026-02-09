package handler

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

type UploadFileHandler struct {
	errResponse *helper.ErrResponse
}

func (u *UploadFileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		u.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	privateIdStr := r.PathValue("private_id")
	privateId, err := strconv.ParseInt(privateIdStr, 10, 64)
	if err != nil {
		u.errResponse.BadRequestResponse(w, r, err)
		return
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		u.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		u.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	defer file.Close()

	dirPath := filepath.Join("files", "chats", fmt.Sprintf("%d", privateId), fmt.Sprintf("%d", userId))

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		u.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	filePath := filepath.Join(dirPath, header.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		u.errResponse.ServerErrorResponse(w, r, err)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		u.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	fileUrl := fmt.Sprintf("/v1/files/chats/%d/%d/%s", privateId, userId, url.PathEscape(header.Filename))

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"data": fileUrl}, nil); err != nil {
		u.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (u *UploadFileHandler) GetFile() http.Handler {
	fs := http.FileServer(http.Dir("./files"))
	return http.StripPrefix("/v1/files", fs)
}

func NewUploadFileHandler(errResponse *helper.ErrResponse) *UploadFileHandler {
	return &UploadFileHandler{
		errResponse: errResponse,
	}
}
