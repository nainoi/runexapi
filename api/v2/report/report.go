package report

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/repository"
)

type ReportAPI struct {
	ReportRepository repository.ReportRepository
}

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
