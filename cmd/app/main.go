package main

import (
	"log"

	"github.com/alisher-99/LomBarter/internal/app"
	"github.com/alisher-99/LomBarter/internal/config"
	docs "github.com/alisher-99/LomBarter/swagger"
)

// @title           ServiceName API
// @version         1.0
// @description     Сервис для работы с ServiceName
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Конфигурация приложения
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Ошибка инициализации конфигурации: %s", err)
	}

	// Документация Swagger
	docs.SwaggerInfo.Host = cfg.Host

	// Запуск приложения
	err = app.Run(cfg)
	if err != nil {
		// позволяет дать знать кубернетесу, что процесс завершился с ошибкой
		log.Fatalf("ошибка при запуске: %s", err)
	}
}
