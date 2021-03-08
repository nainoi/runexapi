package eventv2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	//handle_user "thinkdev.app/think/runex/runexapi/api/v1/user"

	"thinkdev.app/think/runex/runexapi/api/v2/response"

	//"thinkdev.app/think/runex/runexapi/repository"
	repo2 "thinkdev.app/think/runex/runexapi/repository/v2"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"thinkdev.app/think/runex/runexapi/model"
)

// EventStatus struct
type EventStatus struct {
	Status string `json:"status" bson:"status" binding:"required"`
}

// GetAll api godoc
// @Summary Get event all
// @Description get Get event all API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags event
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.EventList}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /event/all [get]
func GetAll(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	urlS := fmt.Sprintf("%s/events", viper.GetString("events.api"))
	var bearer = "Bearer olcgZVpqDXQikRDG"
	//reqURL, _ := url.Parse(urlS)
	req, err := http.NewRequest("GET", urlS, nil)
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//log.Println(err)
			appG.Response(http.StatusOK, err.Error(), model.EventList{})
			c.Abort()
			return
		}
		var events model.EventLists
		err = json.Unmarshal(body, &events)
		if err != nil {
			//log.Println(err)
			appG.Response(http.StatusOK, err.Error(), model.EventList{})
			c.Abort()
			return
		}
		if events.Events != nil {
			appG.Response(http.StatusOK, "success", events.Events)
			c.Abort()
			return
		}
		appG.Response(http.StatusOK, "success", model.EventList{})
		c.Abort()
		return
	}

	appG.Response(http.StatusInternalServerError, err.Error(), nil)

}

// GetDetail api godoc
// @Summary Get event detail
// @Description get Get event detail API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags event
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.EventList}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /event/detail/{code} [get]
func GetDetail(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	code := c.Param("code")
	if code == "" {
		appG.Response(http.StatusInternalServerError, "code not found", nil)
		c.Abort()
		return
	}
	event, err := repo2.DetailEventByCode(code)

	if err != nil {
		appG.Response(http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return

	}

	appG.Response(http.StatusOK, "success", event)
	c.Abort()
	return

}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
	return nil
}
