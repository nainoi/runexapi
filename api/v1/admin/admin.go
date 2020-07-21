package admin

import (
	"log"
	"net/http"

	"bitbucket.org/suthisakch/runex/model"
	"bitbucket.org/suthisakch/runex/pkg/app"
	"bitbucket.org/suthisakch/runex/pkg/e"
	"bitbucket.org/suthisakch/runex/repository"
	"bitbucket.org/suthisakch/runex/utils"
	"github.com/gin-gonic/gin"
)

type AdminAPI struct {
	AdminRepository repository.AdminRepository
}

func (api AdminAPI) ChangeTicket(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	isAdmin := utils.ISAdmin(c)

	if isAdmin != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you do not permission access."})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api AdminAPI) ChangeShppingAddress(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	isAdmin := utils.ISAdmin(c)

	if isAdmin != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you do not permission access."})
		return
	}

	id := c.Param("id")
	log.Printf("[info] id %s", id)

	var json model.ShipingAddressUpdateForm

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.AdminRepository.EditShppingAddress(id, json)
	if err != nil {
		log.Println("error ChangeShppingAddress", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api AdminAPI) UpdateSlip(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	isAdmin := utils.ISAdmin(c)

	if isAdmin != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you do not permission access."})
		return
	}

	userID, _, _ := utils.GetTokenValue(c)

	var json model.SlipUpdateForm

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.AdminRepository.UpdateSlip(json, userID)
	if err != nil {
		log.Println("error ChangeShppingAddress", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api AdminAPI) GetSlip(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	isAdmin := utils.ISAdmin(c)

	if isAdmin != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you do not permission access."})
		return
	}

	regID := c.Param("regID")

	slips, err := api.AdminRepository.GetSlipByReg(regID)
	if err != nil {
		log.Println("error get slips ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, slips)
}

func (api AdminAPI) GetRegEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	isAdmin := utils.ISAdmin(c)

	if isAdmin != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you do not permission access."})
		return
	}

	regID := c.Param("regID")
	register, err := api.AdminRepository.GetRegEventByID(regID)
	if err != nil {
		log.Println("error GetAll", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, register)
}