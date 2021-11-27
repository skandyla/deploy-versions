package transport

import (
	"context"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/skandyla/deploy-versions/models"
)

type Versions interface {
	Create(ctx context.Context, version models.CreateVersionRequest) error
	GetByID(ctx context.Context, id int) (models.VersionDBModel, error)
	GetAll(ctx context.Context) ([]models.VersionResponse, error)
	Delete(ctx context.Context, id int) error
	Health(ctx context.Context) error
	//Update(ctx context.Context, id int64, req models.UpdateVersionRequest) error
}

type Handler struct {
	versionsService Versions
}

func NewHandler(versions Versions) *Handler {
	return &Handler{
		versionsService: versions,
	}
}

func (h *Handler) InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(loggingMiddleware) //test our own middleware implementation

	r.Route("/info", func(r chi.Router) {
		r.Get("/", h.info)
	})

	r.Route("/versions", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", h.getAllVersions)
	})

	r.Route("/version", func(r chi.Router) {
		r.Post("/", h.createVersion)
		r.Route("/{buildID}", func(r chi.Router) {
			r.Get("/", h.getVersionByID)
			r.Put("/", h.updateVersionByID) //update entity
			r.Delete("/", h.deleteVersionByID)
		})
	})

	return r
}
