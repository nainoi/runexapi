package activity

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	_ "io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	guuid "github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
	"thinkdev.app/think/runex/runexapi/utils"
)

type ActivityAPI struct {
	ActivityRepository repository.ActivityRepository
}

type ActivityEvent struct {
	EventID string `json:"event_id" bson:"event_id" binding:"required"`
}

func (api ActivityAPI) GetActivityByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	event := c.Param("event")
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _, _ := utils.GetTokenValue(c)

	eventUser := event + "." + userID

	activity, err := api.ActivityRepository.GetActivityByEvent(eventUser)

	if err != nil {
		log.Println("error Get Event info", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)
}

func (api ActivityAPI) GetActivityByEvent2(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	event := c.Param("event")
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _, _ := utils.GetTokenValue(c)

	eventUser := event + "." + userID

	activity, err := api.ActivityRepository.GetActivityByEvent2(eventUser)

	if err != nil {
		log.Println("error Get Event info2", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)
}

func (api ActivityAPI) AddActivity(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	var form model.AddActivityForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _, _ := utils.GetTokenValue(c)
	path := ""
	file, header, err := c.Request.FormFile("image")
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
		m := resize.Resize(680, 0, img, resize.Lanczos3)

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
		//_, err = io.Copy(out, file)
	}
	fmt.Println(form)
	time1, err := time.Parse(time.RFC3339, form.ActivityDate)
	if err != nil {
		fmt.Println(err)
		time1 = time.Now()
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	form.UserID = userID
	form.ImageURL = path

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	eventUser := form.EventCode + "." + userID

	//event.UpdatedTime = time.Now()
	//loc, _ := time.LoadLocation("Asia/Bangkok")

	activityInfo := model.ActivityInfo{
		Caption:      form.Caption,
		Distance:     form.Distance,
		ImageURL:     form.ImageURL,
		ActivityDate: time1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	activityModel := model.AddActivity{
		EventUser:    eventUser,
		UserID:       userObjectID,
		EventCode:    form.EventCode,
		ActivityInfo: activityInfo,
	}

	err2 := api.ActivityRepository.AddActivity(activityModel)
	if err2 != nil {
		log.Println("error AddActivity", err2.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api ActivityAPI) AddMultiActivity(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var form model.AddMultiActivityForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(form)

	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _, _ := utils.GetTokenValue(c)
	path := ""
	file, header, err := c.Request.FormFile("image")
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
		m := resize.Resize(680, 0, img, resize.Lanczos3)

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
		//_, err = io.Copy(out, file)
	}
	fmt.Println(form)
	time1, err := time.Parse(time.RFC3339, form.ActivityDate)
	if err != nil {
		fmt.Println(err)
		time1 = time.Now()
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	form.UserID = userID
	form.ImageURL = path

	var arrEvent = form.EventCodes

	for index, each := range arrEvent {
		fmt.Printf("EventID value [%d] is [%s]\n", index, each)
		eventCode := each
		userObjectID, err := primitive.ObjectIDFromHex(userID)
		eventUser := eventCode + "." + userID

		if err != nil {
			log.Fatal(err)
		}

		activityInfo := model.ActivityInfo{
			Caption:      form.Caption,
			Distance:     form.Distance,
			ImageURL:     form.ImageURL,
			ActivityDate: time1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		activityModel := model.AddActivity{
			EventUser:    eventUser,
			UserID:       userObjectID,
			EventCode:    eventCode,
			ActivityInfo: activityInfo,
		}

		err2 := api.ActivityRepository.AddActivity(activityModel)
		if err2 != nil {
			log.Println("error AddActivity", err2.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
			return
		}

	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api ActivityAPI) GetHistoryDayByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var form model.HistoryDayFilter
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _, _ := utils.GetTokenValue(c)
	eventUser := form.EventID + "." + userID

	activity, err := api.ActivityRepository.GetHistoryDayByEvent(eventUser, form.Year, form.Month)

	if err != nil {
		log.Println("error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)

}

func (api ActivityAPI) GetHistoryMonthByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var form model.HistoryDayFilter
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _, _ := utils.GetTokenValue(c)
	eventUser := form.EventID + "." + userID

	activity, err := api.ActivityRepository.HistoryMonthByEvent(eventUser, form.Year)

	if err != nil {
		log.Println("error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)

}

func (api ActivityAPI) DeleteActivityEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	id := c.Param("id")
	activityID := c.Param("activityID")

	log.Printf("[info] id %s", id)
	log.Printf("[info] ActivityID %s", activityID)

	userID, _, _ := utils.GetTokenValue(c)
	eventUser := id + "." + userID

	err := api.ActivityRepository.DeleteActivity(eventUser, activityID)

	if err != nil {
		log.Println("error Delete Activity", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

//GetActivityAllInfo activity
func (api ActivityAPI) GetActivityAllInfo(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	eventID := c.Param("event_id")
	if eventID == "" {
		appG.Response(http.StatusBadRequest, http.StatusBadRequest, nil)
		return
	}
	datas, err := api.ActivityRepository.GetActivityAllEvent(eventID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, http.StatusInternalServerError, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, datas)
}
