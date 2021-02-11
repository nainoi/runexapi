package paymethod

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

// Get api godoc
// @Summary get paymethod datas list
// @Description get paymethod datas list
// @Consume application/x-www-form-urlencoded
// @Tags paymethod
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.PaymentMethod}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /paymethods  [get]
func Get(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)

	p := repository.Get()
	res.Response(http.StatusOK, "success", p)
	return
}