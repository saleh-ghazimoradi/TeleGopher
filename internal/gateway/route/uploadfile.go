package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type UploadFileRoute struct {
	middleware        *middleware.Middleware
	uploadFileHandler *handler.UploadFileHandler
}

func (u *UploadFileRoute) UploadFileRoutes(mux *http.ServeMux) {
	mux.Handle("POST /v1/files/{id}", u.middleware.WrapAuth(u.uploadFileHandler.UploadFile))
	mux.Handle("GET /v1/files/", u.middleware.WrapAuth(u.uploadFileHandler.GetFile().ServeHTTP))
}

func NewUploadFileRoute(middleware *middleware.Middleware, uploadFileHandler *handler.UploadFileHandler) *UploadFileRoute {
	return &UploadFileRoute{
		middleware:        middleware,
		uploadFileHandler: uploadFileHandler,
	}
}
