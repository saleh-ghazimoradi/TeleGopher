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
)

type UploadFileHandler struct{}

// UploadFile godoc
// @Summary      Upload a file
// @Description  Upload a file to a specific chat conversation. Supports images, documents, and other file types up to 50MB.
// @Tags         Files
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        id path int true "Chat conversation ID (private conversation ID)"
// @Param        file formData file true "File to upload (max 50MB)"
// @Success      200 {object} helper.Response{data=string} "File URL returned in data field"
// @Failure      400 {object} helper.Response "Invalid conversation ID or missing file"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      403 {object} helper.Response "Forbidden - User not authorized to upload to this conversation"
// @Failure      413 {object} helper.Response "File too large (max 50MB)"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /files/{id} [post]
func (u *UploadFileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
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

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		helper.BadRequestResponse(w, "File too large or invalid form data", err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		helper.InternalServerError(w, "Failed to get file", err)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			helper.InternalServerError(w, "Failed to close file", err)
			return
		}
	}()

	if header.Size > 50<<20 {
		helper.BadRequestResponse(w, "File size exceeds 50MB limit", fmt.Errorf("file size: %d bytes", header.Size))
		return
	}

	dirPath := filepath.Join("files", "chats", fmt.Sprintf("%d", id), fmt.Sprintf("%d", userId))

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		helper.InternalServerError(w, "Failed to create dir", err)
		return
	}

	filePath := filepath.Join(dirPath, header.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		helper.InternalServerError(w, "Failed to create file", err)
		return
	}

	defer func() {
		if err := dst.Close(); err != nil {
			helper.InternalServerError(w, "Failed to close file", err)
			return
		}
	}()

	if _, err := io.Copy(dst, file); err != nil {
		helper.InternalServerError(w, "Failed to copy file", err)
		return
	}

	fileUrl := fmt.Sprintf("/v1/files/chats/%d/%d/%s", id, userId, url.PathEscape(header.Filename))

	helper.SuccessResponse(w, "File successfully uploaded", fileUrl)
}

// GetFile godoc
// @Summary      Serve a file
// @Description  Retrieve and serve an uploaded file
// @Tags         Files
// @Produce      application/octet-stream
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Success      200 {file} binary "File content"
// @Failure      400 {object} helper.Response "Invalid parameters"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      403 {object} helper.Response "Forbidden - User not authorized to access this file"
// @Failure      404 {object} helper.Response "File not found"
// @Router       /files/ [get]
func (u *UploadFileHandler) GetFile() http.Handler {
	fs := http.FileServer(http.Dir("./files"))
	return http.StripPrefix("/v1/files", fs)
}

func NewUploadFileHandler() *UploadFileHandler {
	return &UploadFileHandler{}
}
