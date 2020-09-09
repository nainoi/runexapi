package workouts

import (
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
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
	"thinkdev.app/think/runex/runexapi/utils"
)

type WorkoutsAPI struct {
	WorkoutsRepository repository.WorkoutsRepository
}

func (api WorkoutsAPI) AddWorkout(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	var form model.AddWorkoutForm
	fmt.Println(form)
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

		pathDir := "./upload/image/workout/" + userID + "/" + strconv.Itoa(year) + "_" + strconv.Itoa(int(month))
		if _, err := os.Stat(pathDir); os.IsNotExist(err) {
			os.MkdirAll(pathDir, os.ModePerm)
		}

		out, err := os.Create(pathDir + "/" + uniqidFilename.String() + ".png")

		path = "/upload/image/workout/" + userID + "/" + strconv.Itoa(year) + "_" + strconv.Itoa(int(month)) + "/" + uniqidFilename.String() + ".png"
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		// write new image to file
		jpeg.Encode(out, m, nil)
		//_, err = io.Copy(out, file)
	}
	//fmt.Println(form)
	time1, err := time.Parse(time.RFC3339, form.ActivityDate)
	if err != nil {
		fmt.Println(err)
		time1 = time.Now()
		//c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	form.UserID = userID
	form.ImagePath = path

	userObjectID, err := primitive.ObjectIDFromHex(userID)

	workInfo := model.WorkoutActivityInfo{
		ActivityType: form.ActivityType,
		Calory:       form.Calory,
		Caption:      form.Caption,
		Distance:     form.Distance,
		Pace:         form.Pace,
		Time:         form.Time,
		ActivityDate: time1,
		ImagePath:    form.ImagePath,
		GpxData:      form.GpxData,
	}

	workoutModel := model.AddWorkout{
		UserID:              userObjectID,
		WorkoutActivityInfo: workInfo,
	}

	err2 := api.WorkoutsRepository.AddWorkout(workoutModel)
	if err2 != nil {
		log.Println("error AddWorkout", err2.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
