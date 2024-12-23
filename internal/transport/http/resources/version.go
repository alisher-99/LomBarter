package resources

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const APIVersion = "v1"

// VersionResponse - ответ на запрос версии.
type VersionResponse struct {
	API     string `json:"api"`
	Version string `json:"version"`
}

// VersionResource - структура содержащая версию API и приложения.
type VersionResource struct {
	Version string
}

// Routes возвращает роутер для раздачи версии API.
func (vr VersionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", vr.Get)

	return r
}

// Get возвращает версию API.
func (vr VersionResource) Get(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, VersionResponse{
		API:     APIVersion,
		Version: vr.Version,
	})
}
