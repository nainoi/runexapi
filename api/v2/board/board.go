package board

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

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
// @Param payload body model.RankingRequest true "payload"
// @Success 200 {object} response.Response{data=BoardResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/board/ranking [post]
func GetBoardByEvent(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	var req model.RankingRequest
	if err := c.ShouldBind(&req); err != nil {
		appG.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _ := oauth.GetValuesToken(c)

	event, count, allActivities, myActivities, err := repository.GetBoardByEvent(req, userID)

	if err != nil {
		log.Println("error Get Event info", err.Error())
		appG.Response(http.StatusBadRequest, "Error Get Event info", BoardResponse{
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

	appG.Response(http.StatusOK, "success", ranks)
}

// GetAllBoardByEvent api godoc
// @Summary get leader board activty
// @Description get leader board activty API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags board
// @Accept  application/json
// @Produce application/json
// @Param payload body model.AllRankingRequest true "payload"
// @Success 200 {object} response.Response{data=BoardResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/board/rankings [post]
func GetAllBoardByEvent(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	var req model.AllRankingRequest
	if err := c.ShouldBind(&req); err != nil {
		appG.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}

	event, count, allActivities, err := repository.GetAllBoardByEvent(req)

	if err != nil {
		log.Println("error Get Event info", err.Error())
		appG.Response(http.StatusBadRequest, "Error Get Event info", BoardResponse{
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
	}

	appG.Response(http.StatusOK, "success", ranks)
}

func GetBoardUpdateActivity(c *gin.Context) {
	repository.UpdateActivityV3()
}