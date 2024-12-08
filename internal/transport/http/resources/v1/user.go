package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/example/gophers/libs/errors/httperrors"
	"gitlab.com/example/gophers/libs/logger"

	"github.com/alisher-99/LomBarter/internal/domain/entity"
	"github.com/alisher-99/LomBarter/internal/domain/form"
	"github.com/alisher-99/LomBarter/internal/service"
	"github.com/alisher-99/LomBarter/internal/transport/http/resources/detector"
)

// UserResource представляет собой обработчик для пользователей.
type UserResource struct {
	userService service.UserService // Сервис для работы с пользователями
	logger      logger.Logger       // Логирование запросов и ошибок обработчиков
	json        jsoniter.API        // JSON-парсер
}

// NewUserHandler создает новый экземпляр UserResource.
func NewUserHandler(userService service.UserService, log logger.Logger) *UserResource {
	return &UserResource{
		userService: userService,
		logger:      log,
		json:        jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

// Routes возвращает роутер для обработчика пользователей.
func (vr UserResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", vr.getUsers)
	r.Get("/{id}", vr.getByID)
	r.Post("/", vr.createUser)

	return r
}

// getUsers возвращает список пользователей по фильтру.
// @Summary Получение списка пользователей
// @Description Получение списка пользователей
// @Tags users
// @Accept json
// @Produce json
// @Param filter query form.UsersGetByBio false "Фильтр"
// @Success 200 {object} entity.List{items=entity.Users}
// @Failure 400 {object} swagger.HTTPResponse400 "Код ошибки"
// @Failure 500 {object} swagger.HTTPResponse500 "Внутренняя ошибка сервера"
// @Router /v1/users [get]
func (vr UserResource) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := form.UsersGetByBio{
		Bio: r.URL.Query().Get("bio"),
	}

	users, err := vr.userService.GetUsersByBio(ctx, filter)
	if err != nil {
		vr.logger.Errorf("Ошибка при получении списка пользователей: %v", err)
		_ = render.Render(w, r, detector.Error(err))

		return
	}

	render.JSON(w, r, entity.List{
		Items: users,
	})
}

// createUser создает нового пользователя.
// @Summary Создание пользователя
// @Description Создание пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body form.UserCreate true "Пользователь"
// @Success 200 {string} string
// @Failure 400 {object} swagger.HTTPResponse400 "Код ошибки"
// @Failure 500 {object} swagger.HTTPResponse500 "Внутренняя ошибка сервера"
// @Router /v1/users [post]
func (vr UserResource) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createForm form.UserCreate
	if err := vr.json.NewDecoder(r.Body).Decode(&createForm); err != nil {
		_ = render.Render(w, r, httperrors.BadRequest(err, entity.UserDecodeCode))

		return
	}

	createdUser, err := vr.userService.CreateUser(ctx, createForm, time.Now().UTC())
	if err != nil {
		vr.logger.Errorf("Ошибка при создании пользователя: %v", err)
		_ = render.Render(w, r, detector.Error(err))

		return
	}

	render.JSON(w, r, createdUser)
}

// getByID возвращает пользователя по его идентификатору.
// @Summary Получение пользователя по идентификатору
// @Description Получение пользователя по идентификатору
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Идентификатор пользователя"
// @Success 200 {object} entity.User
// @Failure 400 {object} swagger.HTTPResponse400 "Код ошибки"
// @Failure 500 {object} swagger.HTTPResponse500 "Внутренняя ошибка сервера"
// @Router /v1/users/{id} [get]
func (vr UserResource) getByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		_ = render.Render(w, r, detector.Error(entity.ErrUserIDEmpty))

		return
	}

	user, err := vr.userService.GetUserByID(ctx, id)
	if err != nil {
		vr.logger.Errorf("Ошибка при получении пользователя по идентификатору %s: %v", id, err)
		_ = render.Render(w, r, detector.Error(err))

		return
	}

	render.JSON(w, r, user)
}
