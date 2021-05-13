package team

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/api/v2/upload"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

// UploadTeamIcon api godoc
// @Summary Add team icon
// @Description save Add Team Icon API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags team
// @Accept  application/json
// @Produce application/json
// @Param payload body model.TeamIcon true "payload"
// @Success 200 {object} response.Response{data=model.TeamIcon}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /teamIcon [post]
func UploadTeamIcon(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	var regID = c.PostForm("reg_id")
	if regID == "" {
		appG.Response(http.StatusBadRequest, "error", gin.H{"error": "Parameter not found"})
		return
	}
	objectID, err := primitive.ObjectIDFromHex(regID)
	if err != nil {
		appG.Response(http.StatusBadRequest, "error", gin.H{"error": "Parameter not found"})
		return
	}
	var form = model.TeamIcon {
		RegID: objectID,
	}
	// if err := c.ShouldBind(&form); err != nil {
	// 	appG.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
	// 	return
	// }
	//userID := "5d772660c8a56133c2d7c5ba"

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

		uniqidFilename := uuid.New()

		pathDir := "./upload/image/profile/"
		if _, err := os.Stat(pathDir); os.IsNotExist(err) {
			os.MkdirAll(pathDir, os.ModePerm)
		}

		out, err := os.Create(pathDir + "/" + uniqidFilename.String() + ".png")

		path = "/upload/image/profile/" + uniqidFilename.String() + ".png"
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		// write new image to file
		jpeg.Encode(out, m, nil)

		resObject = upload.UploadWithFolderToS3(pathDir+"/"+uniqidFilename.String()+".png", "profile", uniqidFilename.String()+".png")
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

	form.IconURL = resObject.URL

	err = repository.UpdateIcon(form)
	if err != nil {
		log.Println("error update icon", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), form)
		return
	}

	appG.Response(http.StatusOK, "success", form)
}

// GetTeamIcon api godoc
// @Summary Get team icon
// @Description Get Team Icon API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags team
// @Accept  application/json
// @Produce application/json
// @Param payload body model.TeamIcon true "payload"
// @Success 200 {object} response.Response{data=model.TeamIcon}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /getTeamIcon [post]
func GetTeamIcon(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)

	var form model.TeamIcon
	if err := c.ShouldBind(&form); err != nil {
		appG.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}

	form, err := repository.GetTeamIcon(form)
	if err != nil {
		log.Println("error update icon", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), form)
		return
	}

	appG.Response(http.StatusOK, "success", form)
}
