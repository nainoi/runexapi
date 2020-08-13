package board

import (
	_ "image/png"
	_ "io"
	"log"
	"net/http"

	"bitbucket.org/suthisakch/runex/model"
	"bitbucket.org/suthisakch/runex/pkg/app"
	"bitbucket.org/suthisakch/runex/pkg/e"
	"bitbucket.org/suthisakch/runex/repository"
	"bitbucket.org/suthisakch/runex/utils"
	"github.com/gin-gonic/gin"
)

type BoardAPI struct {
	BoardRepository repository.BoardRepository
}

type BoardEvent struct {
	EventID string `json:"event_id" bson:"event_id" binding:"required"`
}

type BoardResponse struct {
	AllRank []model.Ranking `json:"ranks"`
	MyRank  []model.Ranking `json:"myrank`
}

func (api BoardAPI) GetBoardByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	eventID := c.Param("eventID")
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _, _ := utils.GetTokenValue(c)

	allActivities, myActivities, err := api.BoardRepository.GetBoardByEvent(eventID, userID)

	if err != nil {
		log.Println("error Get Event info", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ranks := BoardResponse{
		AllRank: allActivities,
		MyRank:  myActivities,
	}

	appG.Response(http.StatusOK, e.SUCCESS, ranks)
}
