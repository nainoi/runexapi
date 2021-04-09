package activityV2

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	guuid "github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/api/v2/upload"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
	v2 "thinkdev.app/think/runex/runexapi/repository/v2"
	"thinkdev.app/think/runex/runexapi/utils"
)

//ActivityV2API ref struct
type ActivityV2API struct {
	ActivityV2Repository repository.ActivityV2Repository
}

//ActivityEvent ref struct
type ActivityEvent struct {
	EventID string `json:"event_id" bson:"event_id" binding:"required"`
}

// AddFromWorkout api godoc
// @Summary Add activity from workout
// @Description save  Add activity from workout API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags activity
// @Accept  application/json
// @Produce application/json
// @Param payload body model.AddActivityFormWorkout true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /activityWorkout [post]
func (api ActivityV2API) AddFromWorkout(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var form model.AddActivityFormWorkout
	if err := c.ShouldBind(&form); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	userID, _ := oauth.GetValuesToken(c)
	activityInfo := model.ActivityInfo{
		Caption:      form.WorkoutActivityInfo.Caption,
		Distance:     utils.ToFixed(form.WorkoutActivityInfo.Distance, 2),
		ImageURL:     "",
		APP:          form.WorkoutActivityInfo.APP,
		ActivityDate: form.WorkoutActivityInfo.WorkoutDate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}
	rID, _ := primitive.ObjectIDFromHex(form.RegID)
	pID, _ := primitive.ObjectIDFromHex(form.ParentRegID)

	activityModel := model.AddActivityV2{
		UserID:       userObjectID,
		EventCode:    form.EventCode,
		RegID:        rID,
		ParentRegID:  pID,
		OrderID:      form.OrderID,
		Ticket:       form.Ticket,
		ActivityInfo: activityInfo,
	}

	err2 := api.ActivityV2Repository.AddActivity(activityModel)
	if err2 != nil {
		log.Println("error AddActivity", err2.Error())
		res.Response(http.StatusInternalServerError, err2.Error(), nil)
		c.Abort()
		return
	}

	form.WorkoutActivityInfo.IsSync = true
	err = api.ActivityV2Repository.UpdateWorkout(form.WorkoutActivityInfo, userObjectID)

	if err != nil {
		log.Println(err.Error())
	}

	res.Response(http.StatusOK, "success", nil)
}

// AddMultipleFromWorkout api godoc
// @Summary Add activity from workout
// @Description save  Add activity from workout API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags activity
// @Accept  application/json
// @Produce application/json
// @Param payload body model.AddMultiActivityFormWorkout true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /activitiesWorkout [post]
func (api ActivityV2API) AddMultipleFromWorkout(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		log.Println(err)
	}
	workout := c.PostForm("workout_info")
	eventInfo := c.PostForm("event_activity")
	var workoutInfo model.WorkoutActivityInfo
	var eventActivity []model.EventActivity
	//b, err := json.Marshal()
	err = json.Unmarshal([]byte(workout), &workoutInfo)
	if err != nil {
		log.Println("workout")
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}
	err = json.Unmarshal([]byte(eventInfo), &eventActivity)
	if err != nil {
		log.Println("eventActivity")
		log.Println(err.Error())
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	// var form model.AddMultiActivityFormWorkout
	// if err := c.ShouldBind(&form); err != nil {
	// 	res.Response(http.StatusBadRequest, err.Error(), nil)
	// 	c.Abort()
	// 	return
	// }

	userID, _ := oauth.GetValuesToken(c)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		c.Abort()
		return
	}

	var resObject = model.UploadResponse{}

	path := ""
	if err != nil {
		fmt.Println("Error Retrieving the File")
		path = ""
	} else {
		img, str, err := image.Decode(file)
		log.Println(str)
		if err != nil {
			log.Println(err)
		}

		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(960, 0, img, resize.Lanczos3)

		filename := header.Filename
		fmt.Println(filename)

		uniqidFilename := guuid.New()
		fmt.Printf("github.com/google/uuid:         %s\n", uniqidFilename.String())

		t := time.Now()
		year := t.Year()
		month := t.Month()

		pathDir := "./upload/image/running/" + strconv.Itoa(year) + "_" + strconv.Itoa(int(month))
		if _, err := os.Stat(pathDir); os.IsNotExist(err) {
			os.MkdirAll(pathDir, os.ModePerm)
		}

		out, err := os.Create(pathDir + "/" + uniqidFilename.String() + ".png")

		path = "/upload/image/running/" + strconv.Itoa(year) + "_" + strconv.Itoa(int(month)) + "/" + uniqidFilename.String() + ".png"
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		// write new image to file
		jpeg.Encode(out, m, nil)

		resObject = upload.UploadWithFolderToS3(pathDir+"/"+uniqidFilename.String()+".png", "activty", uniqidFilename.String()+".png")
	}

	defer file.Close()
	defer os.Remove(path)

	for index, each := range eventActivity {
		fmt.Printf("EventID value [%d] is [%s]\n", index, each)
		eventActivity := each

		/*if each.Partner.PartnerName != "" {
			if each.Partner.PartnerName == config.PartnerKao {
				_, err := kao.KaoActivity(path, workoutInfo.Distance, workoutInfo.Duration, each.Partner.Slug, each.Partner.RefActivityValue, each.Partner.RefEventValue, each.Partner.RefPhoneValue)
				if err != nil {
					log.Println("error AddActivity", err.Error())
					res.Response(http.StatusInternalServerError, err.Error(), nil)
					c.Abort()
					return
				}
				if err == nil {

					activityInfo := model.ActivityInfo{
						ID:           primitive.NewObjectID(),
						Caption:      workoutInfo.Caption,
						Distance:     utils.ToFixed(workoutInfo.Distance, 2),
						ImageURL:     resObject.URL,
						Time:         workoutInfo.Duration,
						APP:          workoutInfo.APP,
						IsApprove:    true,
						Status:       config.ACTIVITY_STATUS_APPROVE,
						ActivityDate: workoutInfo.WorkoutDate,
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					}

					rID, _ := primitive.ObjectIDFromHex(eventActivity.RegID)
					pID, _ := primitive.ObjectIDFromHex(eventActivity.ParentRegID)

					activityModel := model.AddActivityV2{
						UserID:       userObjectID,
						EventCode:    eventActivity.EventCode,
						RegID:        rID,
						ParentRegID:  pID,
						OrderID:      eventActivity.OrderID,
						Ticket:       eventActivity.Ticket,
						ActivityInfo: activityInfo,
					}

					err2 := api.ActivityV2Repository.AddActivity(activityModel)
					if err2 != nil {
						log.Println("error AddActivity", err2.Error())
						res.Response(http.StatusInternalServerError, err2.Error(), nil)
						c.Abort()
						return
					}

					kaoActivity := model.LogSendKaoActivity{
						UserID:         userObjectID,
						EventCode:      eventActivity.EventCode,
						ActivityInfoID: activityInfo.ID,
						Distance:       utils.ToFixed(workoutInfo.Distance, 2),
						ImageURL:       resObject.URL,
						Time:           workoutInfo.Duration,
						APP:            workoutInfo.APP,
						ActivityDate:   workoutInfo.WorkoutDate,
						Slug:           each.Partner.Slug,
						Ebib:           each.Partner.RefEventKey,
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					}
					_ = api.ActivityV2Repository.AddKaoLogActivity(kaoActivity)
				}
				// log.Println(string(body))
			}
		} else {*/

		activityInfo := model.ActivityInfo{
			ID:           primitive.NewObjectID(),
			Caption:      workoutInfo.Caption,
			Distance:     utils.ToFixed(workoutInfo.Distance, 2),
			ImageURL:     resObject.URL,
			Time:         workoutInfo.Duration,
			APP:          workoutInfo.APP,
			IsApprove:    true,
			Status:       config.ACTIVITY_STATUS_APPROVE,
			ActivityDate: workoutInfo.WorkoutDate,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		rID, _ := primitive.ObjectIDFromHex(eventActivity.RegID)
		pID, _ := primitive.ObjectIDFromHex(eventActivity.ParentRegID)

		activityModel := model.AddActivityV2{
			UserID:       userObjectID,
			EventCode:    eventActivity.EventCode,
			RegID:        rID,
			ParentRegID:  pID,
			OrderID:      eventActivity.OrderID,
			Ticket:       eventActivity.Ticket,
			ActivityInfo: activityInfo,
		}

		err2 := api.ActivityV2Repository.AddActivity(activityModel)
		if err2 != nil {
			log.Println("error AddActivity", err2.Error())
			res.Response(http.StatusInternalServerError, err2.Error(), nil)
			c.Abort()
			return
		}
	}

	// }
	workoutInfo.IsSync = true
	err = api.ActivityV2Repository.UpdateWorkout(workoutInfo, userObjectID)

	if err != nil {
		log.Println(err.Error())
	}

	res.Response(http.StatusOK, "success", nil)
}

// AddActivity api godoc
// @Summary Add activity
// @Description save  Add activity API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags activity
// @Accept  application/json
// @Produce application/json
// @Param payload body model.AddActivityV2 true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /activity [post]
func (api ActivityV2API) AddActivity(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	var form model.AddActivityForm2
	if err := c.ShouldBind(&form); err != nil {
		appG.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _ := oauth.GetValuesToken(c)

	file, header, err := c.Request.FormFile("image")

	var resObject = model.UploadResponse{}
	path := ""

	if err != nil {
		fmt.Println("Error Retrieving the File")
		path = ""
	} else {
		img, str, err := image.Decode(file)
		log.Println(str)
		if err != nil {
			log.Println(err)
		}

		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(960, 0, img, resize.Lanczos3)

		filename := header.Filename
		fmt.Println(filename)

		uniqidFilename := guuid.New()
		fmt.Printf("github.com/google/uuid:         %s\n", uniqidFilename.String())

		t := time.Now()
		year := t.Year()
		month := t.Month()

		pathDir := "./upload/image/running/" + strconv.Itoa(year) + "_" + strconv.Itoa(int(month))
		if _, err := os.Stat(pathDir); os.IsNotExist(err) {
			os.MkdirAll(pathDir, os.ModePerm)
		}

		out, err := os.Create(pathDir + "/" + uniqidFilename.String() + ".png")

		path = "/upload/image/running/" + strconv.Itoa(year) + "_" + strconv.Itoa(int(month)) + "/" + uniqidFilename.String() + ".png"
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		// write new image to file
		jpeg.Encode(out, m, nil)

		resObject = upload.UploadWithFolderToS3(pathDir+"/"+uniqidFilename.String()+".png", "activty", uniqidFilename.String()+".png")
	}

	defer file.Close()
	defer os.Remove(path)
	// fmt.Println(form)
	// time1, err := time.Parse(time.RFC3339, form.ActivityInfo.ActivityDate)
	// if err != nil {
	// 	fmt.Println(err)
	// 	time1 = time.Now()
	// 	//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	// }

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	pRegID, err := primitive.ObjectIDFromHex(form.ParentRegID)
	regID, err := primitive.ObjectIDFromHex(form.RegID)

	log.Println(resObject.URL)

	form.ActivityInfo.ImageURL = resObject.URL

	form.ActivityInfo.CreatedAt = time.Now()
	form.ActivityInfo.UpdatedAt = time.Now()

	activityModel := model.AddActivityV2{
		UserID:       userObjectID,
		EventCode:    form.EventCode,
		ActivityInfo: form.ActivityInfo,
		ParentRegID:  pRegID,
		OrderID:      form.OrderID,
		RegID:        regID,
		Ticket:       form.Ticket,
	}

	err2 := api.ActivityV2Repository.AddActivity(activityModel)
	if err2 != nil {
		log.Println("error AddActivity", err2.Error())
		appG.Response(http.StatusInternalServerError, err2.Error(), gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, "success", nil)
}

func (api ActivityV2API) GetActivityByEvent2(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	eventID := c.Param("event")
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _ := oauth.GetValuesToken(c)

	activity, err := api.ActivityV2Repository.GetActivityByEvent2(eventID, userID)

	if err != nil {
		log.Println("error AddEvent Get Event info2", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, "success", activity)
}

func (api ActivityV2API) GetActivityByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	eventID := c.Param("event")
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _, _ := utils.GetTokenValue(c)

	activity, err := api.ActivityV2Repository.GetActivityByEvent(eventID, userID)

	if err != nil {
		log.Println("error Get Event info", err.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)
}

// GetDashboard api godoc
// @Summary get activity dashboard
// @Description get activity API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags activity
// @Accept  application/json
// @Produce application/json
// @Param payload body model.EventActivityDashboardReq true "payload"
// @Success 200 {object} response.Response{data=[]model.ActivityDashboard}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /activity/dashboard [post]
func (api ActivityV2API) GetDashboard(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	var req model.EventActivityDashboardReq
	if err := c.ShouldBind(&req); err != nil {
		appG.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _ := oauth.GetValuesToken(c)

	activity, err := repository.GetActivityEventDashboard(req, userID)

	regRequest := model.RegEventDashboardRequest{
		EventCode: req.EventCode,
		ParentRegID: req.ParentRegID,
		RegID: req.RegID,
	}

	reg, err := v2.GetRegEventDashboard(regRequest)

	if err != nil {
		log.Println("error Get Event info2", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, "success", model.ActivityDashboard{
		Activity: activity,
		RegisterData: reg,
	})
}

func (api ActivityV2API) GetHistoryDayByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var form model.HistoryDayFilter
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("form ", form)
	userID, _, _ := utils.GetTokenValue(c)
	eventID := form.EventID

	activity, err := api.ActivityV2Repository.GetHistoryDayByEvent(eventID, userID, form.Year, form.Month)

	if err != nil {
		log.Println("error GetHistoryDayByEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)

}

func (api ActivityV2API) GetHistoryMonthByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var form model.HistoryDayFilter
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _, _ := utils.GetTokenValue(c)
	eventID := form.EventID

	activity, err := api.ActivityV2Repository.HistoryMonthByEvent(eventID, userID, form.Year)

	if err != nil {
		log.Println("error GetHistoryMonthByEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)

}

func (api ActivityV2API) DeleteActivityEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	id := c.Param("id")
	activityID := c.Param("activityID")

	log.Printf("[info] id %s", id)
	log.Printf("[info] ActivityID %s", activityID)

	userID, _, _ := utils.GetTokenValue(c)
	eventID := id

	err := api.ActivityV2Repository.DeleteActivity(eventID, userID, activityID)

	if err != nil {
		log.Println("error Delete Activity", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

// GetActivityWaiting api godoc
// @Summary get activity waiting approve
// @Description get activity API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags activity
// @Accept  application/json
// @Produce application/json
// @Param payload body model.OwnerRequest true "payload"
// @Success 200 {object} response.Response{data=model.ActivityInfo}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /activity/waiting [post]
func (api ActivityV2API) GetActivityWaiting(c *gin.Context) {
	var res = response.Gin{C: c}
	var form model.OwnerRequest
	token := c.GetHeader("token")
	key := viper.GetString("public.token")
	if err := c.ShouldBindJSON(&form); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		return
	}
	if token != key {
		res.Response(http.StatusNotFound, "", nil)
		return
	}
	if !v2.IsOwner(form.EventCode, form.OwnerID) {
		res.Response(http.StatusUnauthorized, "You do not have access to the information.", nil)
		return
	}

	datas, err := repository.GetActivityWaitApprove(form.EventCode)
	if err != nil {
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res.Response(http.StatusOK, "success", datas)
}
