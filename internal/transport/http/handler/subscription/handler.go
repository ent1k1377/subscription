package subscription

import (
	"net/http"
	"subscriptions/internal/service"
	"subscriptions/internal/transport/http/common"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	subscriptionService *service.Subscription
}

func NewHandler(subscriptionService *service.Subscription) *Handler {
	return &Handler{
		subscriptionService: subscriptionService,
	}
}

func (h *Handler) Create(ctx *gin.Context) {
	var request CreateSubscriptionRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ToErrorResponse("json body is not valid"))
		return
	}

	if err := ValidateCreateSubscriptionRequest(request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ToErrorResponse(err.Error()))
		return
	}

	params, err := ToCreateSubscriptionParams(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ToErrorResponse("subscription is not valid")) // TODO поменять ошибку
		return
	}

	err = h.subscriptionService.CreateSubscription(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to create the subscription"))
		return
	}

	ctx.JSON(http.StatusCreated, common.ToSuccessfulResponse("created the subscription"))
}

func (h *Handler) GetSubscription(ctx *gin.Context) {
	uuidParam := ctx.Param("uuid")
	uuidParse, err := uuid.Parse(uuidParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ToErrorResponse("uuidParse is not valid"))
		return
	}

	subscription, err := h.subscriptionService.GetSubscription(uuidParse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to get the subscription"))
		return
	}

	response := ToGetSubscriptionResponse(subscription)
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateSubscription(ctx *gin.Context) {
	uuidParam := ctx.Param("uuid")
	uuidParse, err := uuid.Parse(uuidParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ToErrorResponse("uuid is not valid"))
		return
	}

	var request UpdateSubscriptionRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ToErrorResponse("json body is not valid"))
		return
	}

	params := ToUpdateSubscriptionParams(&request)
	err = h.subscriptionService.UpdateSubscription(uuidParse, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to update the subscription"))
		return
	}

	ctx.JSON(http.StatusOK, common.ToSuccessfulResponse("updated the subscription"))
}

func (h *Handler) DeleteSubscription(ctx *gin.Context) {
	uuidParam := ctx.Param("uuid")
	uuidParse, err := uuid.Parse(uuidParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ToErrorResponse("uuid is not valid"))
		return
	}

	err = h.subscriptionService.DeleteSubscription(uuidParse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to delete the subscription"))
		return
	}

	ctx.JSON(http.StatusOK, common.ToSuccessfulResponse("deleted the subscription"))
}

func (h *Handler) ListSubscriptions(ctx *gin.Context) {
	var request ListSubscriptionRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ToErrorResponse("query line is not valid"))
		return
	}

	params, err := ToListSubscriptionParams(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ToErrorResponse("query is not valid"))
		return
	}

	subscription, err := h.subscriptionService.ListSubscriptions(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ToErrorResponse("failed to list the subscription "+err.Error()))
		return
	}

	ctx.JSON(200, ToListSubscriptionResponse(subscription))
}
