package board

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
)

// BoardAPI struct ref repo
type BoardAPI struct {
	BoardRepository repository.BoardRepository
}

// BoardEvent struct
type BoardEvent struct {
	EventID string `json:"event_id" bson:"event_id" binding:"required"`
}

// BoardResponse response struct
type BoardResponse struct {
	Event         model.Event     `json:"event"`
	AllRank       []model.Ranking `json:"ranks"`
	MyRank        []model.Ranking `json:"myrank"`
	TotalActivity int64           `json:"total_activity"`
}

// GetBoardByEvent api godoc
// @Summary get leader board activty
// @Description get leader board activty API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags board
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=BoardResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/board/ranking/:{eventID} [get]
func (api BoardAPI) GetBoardByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	eventID := c.Param("eventID")
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _ := oauth.GetValuesToken(c)

	event, count, allActivities, myActivities, err := api.BoardRepository.GetBoardByEvent(eventID, userID)

	if err != nil {
		log.Println("error Get Event info", err.Error())
		appG.Response(http.StatusBadRequest, e.ERROR, BoardResponse{
			Event:         event,
			TotalActivity: count,
			AllRank:       []model.Ranking{},
			MyRank:        []model.Ranking{},
		})
		return
	}

	ranks := BoardResponse{
		Event:         event,
		TotalActivity: count,
		AllRank:       allActivities,
		MyRank:        myActivities,
	}

	appG.Response(http.StatusOK, e.SUCCESS, ranks)
}
