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

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/v2/kao"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/api/v2/upload"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
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
		Distance:     form.WorkoutActivityInfo.Distance,
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
	eventObjectID, err := primitive.ObjectIDFromHex(form.EventID)
	activityModel := model.AddActivityV2{
		UserID:       userObjectID,
		EventID:      eventObjectID,
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

		if each.Partner.PartnerName != "" {
			if each.Partner.PartnerName == config.PartnerKao {
				body, err := kao.KaoActivity(path, workoutInfo.Distance, workoutInfo.Duration, each.Partner.Slug, each.Partner.RefActivityValue, each.Partner.RefEventValue, each.Partner.RefPhoneValue)
				if err != nil {
					log.Println("error AddActivity", err.Error())
					res.Response(http.StatusInternalServerError, err.Error(), nil)
					c.Abort()
					return
				}
				if err == nil {
					eventObjectID, err := primitive.ObjectIDFromHex(eventActivity.EventID)

					if err != nil {
						log.Println(err)
					}

					activityInfo := model.ActivityInfo{
						ID:           primitive.NewObjectID(),
						Caption:      workoutInfo.Caption,
						Distance:     workoutInfo.Distance,
						ImageURL:     resObject.URL,
						Time:         workoutInfo.Duration,
						APP:          workoutInfo.APP,
						ActivityDate: workoutInfo.WorkoutDate,
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					}

					activityModel := model.AddActivityV2{
						UserID:       userObjectID,
						EventID:      eventObjectID,
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
						EventID:        eventObjectID,
						ActivityInfoID: activityInfo.ID,
						Distance:       workoutInfo.Distance,
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
				log.Println(string(body))
			}
		} else {
			eventObjectID, err := primitive.ObjectIDFromHex(eventActivity.EventID)

			if err != nil {
				log.Println(err)
			}

			activityInfo := model.ActivityInfo{
				ID:           primitive.NewObjectID(),
				Caption:      workoutInfo.Caption,
				Distance:     workoutInfo.Distance,
				ImageURL:     resObject.URL,
				Time:         workoutInfo.Duration,
				APP:          workoutInfo.APP,
				ActivityDate: workoutInfo.WorkoutDate,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			activityModel := model.AddActivityV2{
				UserID:       userObjectID,
				EventID:      eventObjectID,
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

	}
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
// @Param payload body model.AddActivityForm true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /activity [post]
func (api ActivityV2API) AddActivity(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	var form model.AddActivityForm
	if err := c.ShouldBind(&form); err != nil {
		appG.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _ := oauth.GetValuesToken(c)
	path := ""
	file, header, err := c.Request.FormFile("image")

	var resObject = model.UploadResponse{}
	if err != nil {
		fmt.Println("Error Retrieving the File")
		path = ""
	} else {
		img, str, err := image.Decode(file)
		log.Println(str)
		if err != nil {
			log.Println(err)
		}
		file.Close()

		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Resize(800, 0, img, resize.Lanczos3)

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
			log.Println(err)
		}
		defer out.Close()
		// write new image to file
		jpeg.Encode(out, m, nil)

		//_, err = io.Copy(out, file)

		resObject = upload.UploadWithFolderToS3(path, "activity", uniqidFilename.String()+".png")

		defer os.Remove(path)
	}
	fmt.Println(form)
	time1, err := time.Parse(time.RFC3339, form.ActivityDate)
	if err != nil {
		fmt.Println(err)
		time1 = time.Now()
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	form.UserID = userID
	form.ImageURL = resObject.URL

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	eventObjectID, err := primitive.ObjectIDFromHex(form.EventID)

	activityInfo := model.ActivityInfo{
		Caption:      form.Caption,
		Distance:     form.Distance,
		ImageURL:     form.ImageURL,
		ActivityDate: time1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	activityModel := model.AddActivityV2{
		UserID:       userObjectID,
		EventID:      eventObjectID,
		ActivityInfo: activityInfo,
	}

	err2 := api.ActivityV2Repository.AddActivity(activityModel)
	if err2 != nil {
		log.Println("error AddActivity", err2.Error())
		appG.Response(http.StatusInternalServerError, err2.Error(), gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, "success", nil)
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
