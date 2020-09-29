package strava

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/firebase"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository"
)

//API struct for user repository
type API struct {
	Repo repository.StravaRepository
}

// AddStravaActivity godoc
// @Summary Add new activity from strava
// @Description add new activity from strava API calls
// @Consume application/x-www-form-urlencoded
// @Accept  application/json
// @Produce application/json
// @Tags sync
// @Param payload body model.StravaAddRequest true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /strava/activity [post]
func (api API) AddStravaActivity(c *gin.Context) {
	var activity model.StravaAddRequest
	// span, err := tracer.CreateTracerAndSpan("login", c)
	// if err != nil {
	// 	logger.Logger.Errorf(err.Error())
	// }
	var (
		res = response.Gin{C: c}
	)
	err := c.ShouldBindJSON(&activity)
	if err != nil {
		fmt.Println(err.Error())
		//tracer.OnErrorLog(span, err)
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	err = api.Repo.AddActivity(activity)
	if err != nil {
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return
	}
	res.Response(http.StatusOK, "Add strava activity success", nil)
	fcm := firebase.InitializeAppWithServiceAccount()
	firebase.SendToToken(fcm, "", "")
	return
}

// GetStravaActivities godoc
// @Summary get activities from strava
// @Description get activities from strava API calls
// @Consume application/x-www-form-urlencoded
// @Accept  application/json
// @Produce application/json
// @Security bearerAuth
// @Tags sync
// @Success 200 {object} response.Response{data=[]model.StravaActivity}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /strava/activities [get]
func (api API) GetStravaActivities(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	if userID != "" {
		datas, err := api.Repo.GetActivities(userID)
		if err != nil {
			res.Response(http.StatusInternalServerError, err.Error(), datas)
			c.Abort()
			return
		}
		res.Response(http.StatusOK, "success", datas)
		c.Abort()
		return
	}
	res.Response(http.StatusBadRequest, "user not found", nil)
	c.Abort()
	return
}
