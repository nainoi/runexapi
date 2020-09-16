package preorder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository"
)

//API struct for user repository
type API struct {
	PreRepo repository.PreorderRepositoryMongo
}

// SearchPreOrder godoc
// @Summary Search pre order detail
// @Description Search pre order detail API calls
// @Tags preorder
// @Accept  application/json
// @Produce application/json
// @Param payload body model.FindPreOrderRequest true "payload"
// @Success 200 {object} response.Response{data=[]model.PreOrder}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /searchPreOrder [get]
func (api API) SearchPreOrder(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)

	var request model.FindPreOrderRequest
	// span, err := tracer.CreateTracerAndSpan("login", c)
	// if err != nil {
	// 	logger.Logger.Errorf(err.Error())
	// }
	err := c.ShouldBindJSON(&request)
	if err != nil {
		//tracer.OnErrorLog(span, err)
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	preOrders, err := api.PreRepo.FindPreorder(request)
	if err != nil {
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return
	}
	res.Response(http.StatusOK, "success", preOrders)
}
