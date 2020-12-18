package notification

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/firebase"
	"thinkdev.app/think/runex/runexapi/model"
)

// SendOneNotification godoc
// @Summary send notification one token one time
// @Description send notification one token one time
// @Consume application/x-www-form-urlencoded
// @Accept  application/json
// @Produce application/json
// @Tags notification
// @Param payload body model.NotificationRequest true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /notificationOne [post]
func SendOneNotification(c *gin.Context) {
	var noti model.NotificationRequest
	// span, err := tracer.CreateTracerAndSpan("login", c)
	// if err != nil {
	// 	logger.Logger.Errorf(err.Error())
	// }
	var (
		res = response.Gin{C: c}
	)
	err := c.ShouldBindJSON(&noti)
	if err != nil {
		fmt.Println(err.Error())
		//tracer.OnErrorLog(span, err)
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	fcm := firebase.InitializeServiceAccountID()
	ctx := context.Background()
	client, err := fcm.Messaging(ctx)
	if err != nil {
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return
	}
	firebase.SendMulticastAndHandleErrors(ctx, client, []string{noti.Token}, noti.Title, noti.Body)
	res.Response(http.StatusOK, "send noti", nil)
	return
}
