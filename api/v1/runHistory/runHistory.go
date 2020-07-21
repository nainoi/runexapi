package runhistory

import (
	"log"
	"net/http"

	"bitbucket.org/suthisakch/runex/model"
	"bitbucket.org/suthisakch/runex/pkg/app"
	"bitbucket.org/suthisakch/runex/pkg/e"
	"bitbucket.org/suthisakch/runex/repository"
	"bitbucket.org/suthisakch/runex/utils"
	"github.com/gin-gonic/gin"
)

type RunHistoryAPI struct {
	RunHistoryRepository repository.RunHistoryRepository
}

func (api RunHistoryAPI) AddRunHistory(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	var form model.AddHistoryForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _, _ := utils.GetTokenValue(c)

	err2 := api.RunHistoryRepository.AddHistory(userID, form)

	if err2 != nil {
		log.Println("error AddRunHistory", err2.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api RunHistoryAPI) MyRunHistory(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	userID, _, _ := utils.GetTokenValue(c)

	history, err := api.RunHistoryRepository.GetHistoryByUser(userID)
	if err != nil {
		log.Println("error AddEMyRunHistoryvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, history)
}

func (api RunHistoryAPI) DeleteActivityHistory(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	activityID := c.Param("activityID")

	log.Printf("[info] ActivityID %s", activityID)

	userID, _, _ := utils.GetTokenValue(c)

	err := api.RunHistoryRepository.DeleteActivity(userID, activityID)

	if err != nil {
		log.Println("error Delete Activity", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}
