package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type UploadFileRoute struct {
	uploadFileHandler *handler.UploadFileHandler
	middleware        *middleware.Middleware
}

func (u *UploadFileRoute) UploadFileRoutes(mux *http.ServeMux) {
	mux.Handle("POST /v1/files/{private_id}",
		u.middleware.AuthMiddleware(http.HandlerFunc(u.uploadFileHandler.UploadFile)))

	mux.Handle("/v1/files/",
		u.middleware.AuthMiddleware(u.uploadFileHandler.GetFile()))
}

func NewUploadFileRoute(uploadFileHandler *handler.UploadFileHandler, middleware *middleware.Middleware) *UploadFileRoute {
	return &UploadFileRoute{
		uploadFileHandler: uploadFileHandler,
		middleware:        middleware,
	}
}
