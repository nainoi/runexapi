package activityV2

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

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
	"thinkdev.app/think/runex/runexapi/utils"
)

type ActivityV2API struct {
	ActivityV2Repository repository.ActivityV2Repository
}

type ActivityEvent struct {
	EventID string `json:"event_id" bson:"event_id" binding:"required"`
}

func (api ActivityV2API) AddActivity(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)
}

func (api ActivityV2API) GetActivityByEvent2(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	eventID := c.Param("event")
	//userID := "5d772660c8a56133c2d7c5ba"
	userID, _, _ := utils.GetTokenValue(c)

	activity, err := api.ActivityV2Repository.GetActivityByEvent2(eventID, userID)

	if err != nil {
		log.Println("error AddEvent Get Event info2", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, activity)
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
