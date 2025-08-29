package subscription

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"subscriptions/internal/service"
	"subscriptions/internal/transport/http/common"
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
