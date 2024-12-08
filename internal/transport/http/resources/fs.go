package resources

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

// FilesResource для раздачи статичных файлов
type FilesResource struct {
	FilesDir string
}

// Routes возвращает роутер для раздачи статичных файлов
func (fr FilesResource) Routes() chi.Router {
	r := chi.NewRouter()
	filesRoot := http.Dir(fr.FilesDir)

	NewFileServer(r, "/", filesRoot)

	return r
}

// NewFileServer устанавливает обработчик для раздачи статичных файлов
func NewFileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}

	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
