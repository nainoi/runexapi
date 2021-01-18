package workouts

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository"
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
// @Success 200 {object} response.Response{data=model.WorkoutActivityInfo}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workout [post]
func (api WorkoutsAPI) AddWorkout(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var form model.AddWorkoutForm
	fmt.Println(form)
	if err := c.ShouldBind(&form); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _ := oauth.GetValuesToken(c)
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
		APP:              form.APP,
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

	workoutInfo, err2 := api.WorkoutsRepository.AddWorkout(workoutModel)
	if err2 != nil {
		log.Println("error AddWorkout", err2.Error())
		res.Response(http.StatusInternalServerError, err2.Error(), nil)
		return
	}

	res.Response(http.StatusOK, "success", workoutInfo)
}

// AddMultiWorkout api godoc
// @Summary Add workouts multiple
// @Description save workout API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags workouts
// @Accept  application/json
// @Produce application/json
// @Param payload body model.AddMultiWorkout true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workouts [post]
func (api WorkoutsAPI) AddMultiWorkout(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var form model.AddMultiWorkout
	fmt.Println(form)
	if err := c.ShouldBind(&form); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _ := oauth.GetValuesToken(c)
	/*time1, err := time.Parse(time.RFC3339, form.WorkoutDate)
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

	form.UserID = userID*/

	//userObjectID, err := primitive.ObjectIDFromHex(userID)

	// workInfo := model.WorkoutActivityInfo{
	// 	ActivityType:     form.ActivityType,
	// 	APP:              form.APP,
	// 	Calory:           form.Calory,
	// 	Caption:          form.Caption,
	// 	Distance:         form.Distance,
	// 	Pace:             form.Pace,
	// 	Duration:         form.Duration,
	// 	TimeString:       form.TimeString,
	// 	EndDate:          time3,
	// 	StartDate:        time2,
	// 	WorkoutDate:      time1,
	// 	NetElevationGain: form.NetElevationGain,
	// 	IsSync:           form.IsSync,
	// 	Locations:        form.Locations,
	// }

	// workoutModel := model.AddWorkout{
	// 	UserID:              userObjectID,
	// 	WorkoutActivityInfo: workInfo,
	// }

	err := api.WorkoutsRepository.AddMultiWorkout(userID, form.WorkoutActivityInfos)
	if err != nil {
		log.Println("error Add multi Workout", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res.Response(http.StatusOK, "success", nil)
}

// GetWorkouts api godoc
// @Summary Get workouts list
// @Description list workouts API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags workouts
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.Workouts}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workouts [get]
func (api WorkoutsAPI) GetWorkouts(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	isNotHas, workout, err := api.WorkoutsRepository.GetWorkouts(userObjectID)
	if isNotHas {
		res.Response(http.StatusNoContent, "status no content", workout)
		c.Abort()
		return
	} else if err != nil {
		log.Println("error get work", err.Error())
		res.Response(http.StatusInternalServerError, "get workout fail", workout)
		c.Abort()
		return
	}
	res.Response(http.StatusOK, "success", workout)
	c.Abort()
}

// GetWorkoutsHistoryMonth api godoc
// @Summary Get workouts history list by month
// @Description list workouts API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags workouts
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.Workouts}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workouts/history [post]
func (api WorkoutsAPI) GetWorkoutsHistoryMonth(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var form model.WorkoutHistoryMonthFilter
	if err := c.ShouldBind(&form); err != nil {
		log.Println(err.Error())
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}
	userID, _ := oauth.GetValuesToken(c)
	//userID := "5d8820749c3f42e4088c980f"
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	isNotHas, workout, err := api.WorkoutsRepository.HistoryMonth(userObjectID, form.Year)
	if isNotHas {
		res.Response(http.StatusNoContent, "status no content", workout)
		c.Abort()
		return
	} else if err != nil {
		log.Println("error get work", err.Error())
		res.Response(http.StatusInternalServerError, "get workout fail", workout)
		c.Abort()
		return
	}
	res.Response(http.StatusOK, "success", workout)
	c.Abort()
}

// GetWorkoutsHistoryMonth api godoc
// @Summary Get workouts history list by month
// @Description list workouts API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags workouts
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.Workouts}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workouts/historyAll [get]
func (api WorkoutsAPI) GetWorkoutsHistoryAll(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)

	userID, _ := oauth.GetValuesToken(c)
	//userID := "5d8820749c3f42e4088c980f"
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	isNotHas, workout, err := api.WorkoutsRepository.HistoryAll(userObjectID)
	if isNotHas {
		res.Response(http.StatusNoContent, "status no content", workout)
		c.Abort()
		return
	} else if err != nil {
		log.Println("error get work", err.Error())
		res.Response(http.StatusInternalServerError, "get workout fail", workout)
		c.Abort()
		return
	}
	res.Response(http.StatusOK, "success", workout)
	c.Abort()
}

// GetWorkoutsDetail api godoc
// @Summary Get workouts detail by id
// @Description list workouts API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags workouts
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.Workouts}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workoutDetail [get]
func (api WorkoutsAPI) GetWorkoutDetail(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)

	userID, _ := oauth.GetValuesToken(c)
	//userID := "5d8820749c3f42e4088c980f"
	//userID := "5d8aca21d950a8181151aab9"
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	id := c.Param("id")
	workoutID, _ := primitive.ObjectIDFromHex(id)

	isNotHas, workout, err := api.WorkoutsRepository.WorkoutInfo(userObjectID, workoutID)
	if isNotHas {
		res.Response(http.StatusNoContent, "status no content", workout)
		c.Abort()
		return
	} else if err != nil {
		log.Println("error get work", err.Error())
		res.Response(http.StatusInternalServerError, "get workout fail", workout)
		c.Abort()
		return
	}
	res.Response(http.StatusOK, "success", workout)
	c.Abort()
}
