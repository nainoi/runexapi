package category

import (
	"log"
	"net/http"

	"bitbucket.org/suthisakch/runex/model"
	"bitbucket.org/suthisakch/runex/pkg/app"
	"bitbucket.org/suthisakch/runex/pkg/e"
	"bitbucket.org/suthisakch/runex/repository"
	"github.com/gin-gonic/gin"
)

type CategoryAPI struct {
	CategoryRepository repository.CategoryRepository
}

func (api CategoryAPI) AddCategory(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var json model.CategoryMaster

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categoryID, err := api.CategoryRepository.AddCategory(json)
	if err != nil {
		log.Println("error AddCategory", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, categoryID)
}

func (api CategoryAPI) GetAll(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	category, err := api.CategoryRepository.GetCategoryAll()
	if err != nil {
		log.Println("error Get Category", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, category)
}

func (api CategoryAPI) EditCategory(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	id := c.Param("id")
	log.Printf("[info] id %s", id)
	var json model.CategoryUpdateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.CategoryRepository.EditCategory(id, json)
	if err != nil {
		log.Println("error EditCategory", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api CategoryAPI) DeleteCategory(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	log.Printf("[info] id %s", id)

	err := api.CategoryRepository.DeleteCategoryByID(id)
	if err != nil {
		log.Println("error DeleteCategoryByID", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
