package subscription

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ent1k1377/subscriptions/internal/service"
	"github.com/ent1k1377/subscriptions/internal/transport/http/common"
	"github.com/ent1k1377/subscriptions/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	logger              *slog.Logger
	subscriptionService *service.Subscription
}

func NewHandler(baseLogger *slog.Logger, subscriptionService *service.Subscription) *Handler {
	logger := baseLogger.WithGroup("subscription handler")

	return &Handler{
		logger:              logger,
		subscriptionService: subscriptionService,
	}
}

// Create создает запись в бд на основе запроса CreateSubscriptionRequest
//
//	@Summary		Создает подписку
//	@Description	Создает подписку для пользователя
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateSubscriptionRequest	true	"Данные подписки"
//	@Success		201		{object}	common.SuccessfulResponse	"Successfully created"
//	@Failure		400		{object}	common.ErrorResponse
//	@Failure		500		{object}	common.ErrorResponse
//	@Router			/subscriptions [post]
func (h *Handler) Create(c *gin.Context) {
	logger := h.logger.With(
		slog.String("request_id", c.MustGet(middleware.RequestIDKey).(string)),
		slog.String("func", "Create"),
	)

	logger.Info("Start create subscription")
	var request CreateSubscriptionRequest
	if err := c.ShouldBind(&request); err != nil {
		logger.Warn("Failed to bind the body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.ToErrorResponse("json body is not valid"))
		return
	}

	if err := ValidateCreateSubscriptionRequest(request); err != nil {
		logger.Warn("Failed to validate the request", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.ToErrorResponse(err.Error()))
		return
	}

	params, err := ToCreateSubscriptionParams(&request)
	if err != nil {
		logger.Warn("Failed to convert the request to create subscription params", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.ToErrorResponse("subscription is not valid")) // TODO поменять ошибку
		return
	}

	ct := context.WithValue(c.Request.Context(), middleware.RequestIDKey, c.MustGet(middleware.RequestIDKey).(string))
	err = h.subscriptionService.CreateSubscription(ct, params)
	if err != nil {
		logger.Warn("Failed to create subscription", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to create the subscription"))
		return
	}

	logger.Info("Create subscription successfully")
	c.JSON(http.StatusCreated, common.ToSuccessfulResponse("created the subscription"))
}

// GetSubscription возвращает подписку по UUID пользователя
//
//	@Summary		Получить подписку
//	@Description	Возвращает информацию о подписке по UUID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string	true	"UUID подписки"	Format(uuid)	Example(f81d4fae-7dec-11d0-a765-00a0c91e6bf6)
//	@Success		200		{object}	GetSubscriptionResponse
//	@Failure		400		{object}	common.ErrorResponse
//	@Failure		500		{object}	common.ErrorResponse
//	@Router			/subscriptions/{uuid} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	logger := h.logger.With(
		slog.String("func", "GetSubscription"),
	)

	logger.Info("Start get subscription")
	uuidParam := c.Param("uuid")
	uuidParse, err := uuid.Parse(uuidParam)
	if err != nil {
		logger.Warn("Failed to parse the uuid", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.ToErrorResponse("uuidParse is not valid"))
		return
	}

	subscription, err := h.subscriptionService.GetSubscription(uuidParse)
	if err != nil {
		logger.Error("Failed to get subscription", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to get the subscription"))
		return
	}

	logger.Info("Get subscription successfully")
	response := ToGetSubscriptionResponse(subscription)
	c.JSON(http.StatusOK, response)
}

// UpdateSubscription обновляет данные подписки по UUID
//
//	@Summary		Обновить подписку
//	@Description	Обновляет информацию о подписке по UUID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string						true	"UUID подписки"	Format(uuid)	Example(f81d4fae-7dec-11d0-a765-00a0c91e6bf6)
//	@Param			request	body		UpdateSubscriptionRequest	true	"Данные для обновления подписки"
//	@Success		200		{object}	common.SuccessfulResponse
//	@Failure		400		{object}	common.ErrorResponse
//	@Failure		500		{object}	common.ErrorResponse
//	@Router			/subscriptions/{uuid} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	logger := h.logger.With(
		slog.String("func", "UpdateSubscription"),
	)

	uuidParam := c.Param("uuid")
	uuidParse, err := uuid.Parse(uuidParam)
	if err != nil {
		logger.Warn("Failed to parse the uuid", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.ToErrorResponse("uuid is not valid"))
		return
	}

	var request UpdateSubscriptionRequest
	if err := c.ShouldBind(&request); err != nil {
		logger.Warn("Failed to bind the body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.ToErrorResponse("json body is not valid"))
		return
	}

	params := ToUpdateSubscriptionParams(&request)
	err = h.subscriptionService.UpdateSubscription(uuidParse, params)
	if err != nil {
		logger.Error("Failed to update subscription", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to update the subscription"))
		return
	}

	logger.Info("Update subscription successfully")
	c.JSON(http.StatusOK, common.ToSuccessfulResponse("updated the subscription"))
}

// DeleteSubscription удаляет подписку по UUID
//
//	@Summary		Удалить подписку
//	@Description	Удаляет подписку по UUID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string						true	"UUID подписки"	Format(uuid)	Example(f81d4fae-7dec-11d0-a765-00a0c91e6bf6)
//	@Success		200		{object}	common.SuccessfulResponse	"Успешное удаление"
//	@Failure		400		{object}	common.ErrorResponse		"Неверный UUID"
//	@Failure		500		{object}	common.ErrorResponse		"Ошибка сервера"
//	@Router			/subscriptions/{uuid} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	logger := h.logger.With(
		slog.String("func", "DeleteSubscription"),
	)

	uuidParam := c.Param("uuid")
	uuidParse, err := uuid.Parse(uuidParam)
	if err != nil {
		logger.Warn("Failed to parse the uuid", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.ToErrorResponse("uuid is not valid"))
		return
	}

	err = h.subscriptionService.DeleteSubscription(uuidParse)
	if err != nil {
		logger.Error("Failed to delete subscription", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to delete the subscription"))
		return
	}

	logger.Info("Delete subscription successfully")
	c.JSON(http.StatusOK, common.ToSuccessfulResponse("deleted the subscription"))
}

// ListSubscriptions возвращает список подписок с возможностью фильтрации и пагинации
//
//	@Summary		Список подписок
//	@Description	Возвращает список подписок с фильтрацией и пагинацией (через query-параметры)
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			user_id	query		string					false	"Фильтр по UUID пользователя"		Format(uuid)	Example(f81d4fae-7dec-11d0-a765-00a0c91e6bf6)
//	@Param			service	query		string					false	"Фильтр по названию сервиса"		Example(Netflix)
//	@Param			limit	query		int						false	"Количество записей на странице"	minimum(1)	maximum(100)	Example(10)
//	@Param			offset	query		int						false	"Смещение для пагинации"			minimum(0)	Example(0)
//	@Success		200		{array}		GetSubscriptionResponse	"Список подписок"
//	@Failure		400		{object}	common.ErrorResponse	"Неверный запрос"
//	@Failure		500		{object}	common.ErrorResponse	"Ошибка сервера"
//	@Router			/subscriptions/list [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	logger := h.logger.With(
		slog.String("func", "ListSubscriptions"),
	)

	var request ListSubscriptionRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Warn("Failed to bind the body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.ToErrorResponse("query line is not valid"))
		return
	}

	params, err := ToListSubscriptionParams(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ToErrorResponse("query is not valid"))
		return
	}

	subscription, err := h.subscriptionService.ListSubscriptions(params)
	if err != nil {
		logger.Error("Failed to list subscriptions", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to list the subscription "+err.Error()))
		return
	}

	logger.Info("List subscriptions successfully")
	c.JSON(200, ToListSubscriptionResponse(subscription))
}

// TotalCostSubscriptions считает общую стоимость подписок по заданным параметрам
//
//	@Summary		Общая стоимость подписок
//	@Description	Возвращает суммарную стоимость подписок для пользователя (по фильтрам из тела запроса)
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		TotalCostSubscriptionsRequest	true	"Параметры для подсчета стоимости"
//	@Success		200		{object}	map[string]int					"Общая стоимость подписок"
//	@Failure		400		{object}	common.ErrorResponse			"Неверный запрос"
//	@Failure		500		{object}	common.ErrorResponse			"Ошибка сервера"
//	@Router			/subscriptions/total [post]
func (h *Handler) TotalCostSubscriptions(c *gin.Context) {
	logger := h.logger.With(
		slog.String("func", "ListSubscriptions"),
	)

	var request TotalCostSubscriptionsRequest
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		logger.Warn("Failed to bind the body", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, common.ToErrorResponse("json body is not valid"))
		return
	}

	params := ToTotalCostSubscriptionsParams(&request)
	sum, err := h.subscriptionService.TotalCostSubscriptions(params)
	if err != nil {
		logger.Error("Failed to list subscriptions", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to get the total cost subscriptions"))
		return
	}

	logger.Info("Total cost subscriptions successfully")
	c.JSON(http.StatusOK, gin.H{"total": sum})
}
