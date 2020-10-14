package workouts

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
	"thinkdev.app/think/runex/runexapi/utils"
)

// WorkoutsAPI struct repo
type WorkoutsAPI struct {
	WorkoutsRepository repository.WorkoutsRepository
}

// AddWorkout api godoc
// @Summary Add workout
// @Description save workout API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags workouts
// @Accept  application/json
// @Produce application/json
// @Param payload body model.AddWorkoutForm true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workout/add [post]
func (api WorkoutsAPI) AddWorkout(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	var form model.AddWorkoutForm
	fmt.Println(form)
	if err := c.ShouldBind(&form); err != nil {
		appG.Response(http.StatusBadRequest, e.ERROR, nil)
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _, _ := utils.GetTokenValue(c)
	time1, err := time.Parse(time.RFC3339, form.WorkoutDate)
	if err != nil {
		fmt.Println(err)
		time1 = time.Now()
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	time2, err := time.Parse(time.RFC3339, form.StartDate)
	if err != nil {
		fmt.Println(err)
		time2 = time.Now()
	}

	time3, err := time.Parse(time.RFC3339, form.EndDate)
	if err != nil {
		fmt.Println(err)
		time3 = time.Now()
	}

	form.UserID = userID

	userObjectID, err := primitive.ObjectIDFromHex(userID)

	workInfo := model.WorkoutActivityInfo{
		ActivityType:     form.ActivityType,
		App:              form.App,
		Calory:           form.Calory,
		Caption:          form.Caption,
		Distance:         form.Distance,
		Pace:             form.Pace,
		Duration:         form.Duration,
		TimeString:       form.TimeString,
		EndDate:          time3,
		StartDate:        time2,
		WorkoutDate:      time1,
		NetElevationGain: form.NetElevationGain,
		IsSync:           form.IsSync,
		Locations:        form.Locations,
	}

	workoutModel := model.AddWorkout{
		UserID:              userObjectID,
		WorkoutActivityInfo: workInfo,
	}

	err2 := api.WorkoutsRepository.AddWorkout(workoutModel)
	if err2 != nil {
		log.Println("error AddWorkout", err2.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
