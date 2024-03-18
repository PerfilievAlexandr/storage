package http

import (
	"encoding/json"
	"errors"
	"github.com/PerfilievAlexandr/storage/internal/api/http/dtoHttpStorage"
	customErrors "github.com/PerfilievAlexandr/storage/internal/errors"
	"github.com/PerfilievAlexandr/storage/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	storageService service.StorageService
}

func New(storageService service.StorageService) *Handler {
	return &Handler{storageService}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1")
	{
		lists := api.Group("/storage")
		{
			lists.POST("/", h.put)
			lists.GET("/:key", h.get)
			lists.DELETE("/:key", h.delete)
		}
	}

	return router
}

func (h *Handler) put(ctx *gin.Context) {
	var addRequest dtoHttpStorage.AddRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&addRequest); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.storageService.Put(ctx.Request.Context(), addRequest)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) get(ctx *gin.Context) {
	key := ctx.Param("key")

	value, err := h.storageService.Get(ctx.Request.Context(), key)

	if err != nil {
		if errors.Is(customErrors.NotFound, err) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, value)
}

func (h *Handler) delete(ctx *gin.Context) {
	key := ctx.Param("key")

	err := h.storageService.Delete(ctx, key)

	if err != nil {
		if errors.Is(customErrors.NotFound, err) {
			newErrorResponse(ctx, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
