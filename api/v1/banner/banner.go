package banner

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"

	"thinkdev.app/think/runex/runexapi/api/mail"
)

type BannerAPI struct {
	BannerRepository repository.BannerRepository
}

func (api BannerAPI) AddBanner(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var json model.BannerAddForm

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err2 := api.BannerRepository.ExistByEventID(json.EventID)

	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	if exists {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_EVENT, nil)
		return
	}

	bannerID, err := api.BannerRepository.AddBanner(json)
	if err != nil {
		log.Println("error AddBanner", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, bannerID)
}

func (api BannerAPI) DeleteBanner(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	log.Printf("[info] id %s", id)

	err := api.BannerRepository.DeleteBannerByID(id)
	if err != nil {
		log.Println("error DeleteBanner", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api BannerAPI) GetAll(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	banner, err := api.BannerRepository.GetBannerAll()
	if err != nil {
		log.Println("error Get banner", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, banner)
}

func (api BannerAPI) Testmail(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	mail.TestMailTemplate()
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
