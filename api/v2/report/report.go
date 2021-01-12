package report

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/repository"
	repo2 "thinkdev.app/think/runex/runexapi/repository/v2"
)

//ReportAPI struct
type ReportAPI struct {
	ReportRepository repository.ReportRepository
}

// GetDashboardByEvent godoc
// @Summary Get report dashboard event
// @Description get report dashboard event API calls
// @Security bearerAuth
// @Tags report
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.ReportDashboard}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /report/dashboard/{eventID} [get]
func (api ReportAPI) GetDashboardByEvent(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)

	eventID := c.Param("eventID")
	log.Println("eventID :", eventID)
	dashboard, err := api.ReportRepository.GetDashboardByEvent(eventID)

	if err != nil {
		log.Println("error get dashboard", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", dashboard)
}

// GetDashboardByEventCode godoc
// @Summary Get report dashboard event
// @Description get report dashboard event API calls
// @Security bearerAuth
// @Tags report
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.ReportDashboard}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /report/dashboard/{eventID} [get]
func (api ReportAPI) GetDashboardByEventCode(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)

	code := c.Param("code")
	event, err := repo2.DetailEventByCode(code)

	dashboard, err := api.ReportRepository.GetDashboardByEventCode(event)

	if err != nil {
		log.Println("error get dashboard", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", dashboard)
}
