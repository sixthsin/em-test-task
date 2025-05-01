// @title Enrichment API
// @version 1.0
// @description API для обогащения данных о людях. Позволяет создавать, редактировать, удалять и получать информацию о людях с возможностью фильтрации и пагинации.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@enrichment.api
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8081
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name ApiKey
package handlers

import (
	"em-task/cfg"
	"em-task/internal/repository"
	"em-task/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	_ "em-task/docs"

	"em-task/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// EnrichHandlerDeps содержит зависимости для обработчика EnrichHandler
type EnrichHandlerDeps struct {
	*cfg.Config
	*services.EnrichService
	Logger *logging.Logger
}

// EnrichHandler реализует обработчики для API обогащения данных о людях
type EnrichHandler struct {
	*cfg.Config
	*services.EnrichService
	Logger *logging.Logger
}

// NewEnrichHandler создает новый экземпляр обработчика и регистрирует маршруты
func NewEnrichHandler(router *gin.Engine, deps EnrichHandlerDeps) {
	handler := &EnrichHandler{
		Config:        deps.Config,
		EnrichService: deps.EnrichService,
		Logger:        deps.Logger,
	}
	router.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/api/v1/create", handler.CreatePerson)
	router.DELETE("/api/v1/delete/:id", handler.DeletePerson)
	router.GET("/api/v1/get", handler.GetWithParams)
	router.PATCH("/api/v1/edit/:id", handler.EditPerson)
}

// CreatePerson создает новую персону
// @Summary Создать персону
// @Description Создание новой персоны с обогащением данных
// @Tags people
// @Accept json
// @Produce json
// @Param input body repository.CreatePersonDTO true "Данные для создания персоны"
// @Success 200 {object} PersonResponse "Успешное создание"
// @Failure 400 {object} BadRequestResponse "Ошибка валидации"
// @Failure 500 {object} InternalServerErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/create [post]
func (h *EnrichHandler) CreatePerson(c *gin.Context) {
	var requestData repository.CreatePersonDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestData); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	person, err := h.EnrichService.EnrichPersonData(requestData)
	if err != nil {
		h.Logger.ErrorLogger.Errorf("CreatePerson error: %v", err)
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Status:  "error",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, person)
}

// DeletePerson удаляет персону по ID
// @Summary Удалить персону
// @Description Удаление персоны по ID
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID персоны" example(1)
// @Success 200 {object} SuccessfullyDeletedResponse "Успешное удаление"
// @Failure 400 {object} BadRequestResponse "Неверный ID"
// @Failure 500 {object} InternalServerErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/delete/{id} [delete]
func (h *EnrichHandler) DeletePerson(c *gin.Context) {
	personIdString := c.Param("id")
	id, err := strconv.Atoi(personIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID format",
			Status:  "error",
		})
		return
	}

	if err := h.EnrichService.DeletePersonById(id); err != nil {
		h.Logger.ErrorLogger.Errorf("DeletePerson error: %v", err)
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Status:  "error",
		})
		return
	}

	c.JSON(http.StatusOK, ApiResponse{
		Code:    http.StatusOK,
		Message: "Person successfully deleted",
		Status:  "ok",
	})
}

// GetWithParams получает список персон с параметрами
// @Summary Получить список персон
// @Description Получение списка персон с использованием параметров фильтрации, лимита и смещения
// @Tags people
// @Accept json
// @Produce json
// @Param limit query int false "Лимит" example(10)
// @Param offset query int false "Смещение" example(0)
// @Param filters query repository.PersonParams false "Параметры фильтрации"
// @Success 200 {object} PersonResponse "Успешное получение"
// @Failure 400 {object} BadRequestResponse "Неверные параметры запроса"
// @Failure 500 {object} InternalServerErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/get [get]
func (h *EnrichHandler) GetWithParams(c *gin.Context) {
	limitString := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitString)
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid limit parameter",
		})
		return
	}

	offsetString := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetString)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid offset parameter",
		})
		return
	}

	var filters repository.PersonParams
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid query parameters",
		})
		return
	}

	people := h.EnrichService.GetPeopleWithParams(limit, offset, filters)
	c.JSON(http.StatusOK, people)
}

// EditPerson редактирует данные персоны
// @Summary Редактировать персону
// @Description Редактирование данных персоны по ID
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID персоны" example(1)
// @Param input body repository.FullPersonDTO true "Обновленные данные"
// @Success 200 {object} PersonResponse "Обновленная персона"
// @Failure 400 {object} BadRequestResponse "Ошибка валидации"
// @Failure 500 {object} InternalServerErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/edit/{id} [patch]
func (h *EnrichHandler) EditPerson(c *gin.Context) {
	personIdString := c.Param("id")
	id, err := strconv.Atoi(personIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: "Invalid ID format",
		})
		return
	}

	var requestData repository.FullPersonDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestData); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	editedPerson, err := h.EnrichService.EditPerson(id, requestData)
	if err != nil {
		h.Logger.ErrorLogger.Errorf("EditPerson error: %v", err)
		c.JSON(http.StatusBadRequest, ApiResponse{
			Status:  "error",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, editedPerson)
}
