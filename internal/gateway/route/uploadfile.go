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
	mux.HandleFunc("POST /v1/files/{private_id}",
		u.middleware.Authenticate(u.uploadFileHandler.UploadFile))

	fileHandler := u.middleware.Authenticate(u.uploadFileHandler.GetFile().ServeHTTP)
	mux.Handle("GET /v1/files/", fileHandler)
}

func NewUploadFileRoute(uploadFileHandler *handler.UploadFileHandler, middleware *middleware.Middleware) *UploadFileRoute {
	return &UploadFileRoute{
		uploadFileHandler: uploadFileHandler,
		middleware:        middleware,
	}
}
