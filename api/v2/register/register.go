package registerV2

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/mail"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository/v2"
	"thinkdev.app/think/runex/runexapi/utils"
)

const (
	//TEST pkey_test_5i6ivm4cotoab601bfr
	//skey_test_5h9vov2hqe9iv55o8tu
	// Read these from environment variables or configuration files!
	//PD pkey_5i1p3nkjgq6vrrrfhkp
	// skey_5i6wunbg5thk3eqd5kk
	//product
	OmisePublicKey = "pkey_5i1p3nkjgq6vrrrfhkp"
	OmiseSecretKey = "skey_5i6wunbg5thk3eqd5kk"

	//test
	// OmisePublicKey = "pkey_test_5i6ivm4cotoab601bfr"
	// OmiseSecretKey = "skey_test_5h9vov2hqe9iv55o8tu"
)

type RegisterAPI struct {
	RegisterRepository repository.RegisterRepository
}

func (api RegisterAPI) GetByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	eventID := c.Param("eventID")
	register, err := api.RegisterRepository.GetRegisterByEvent(eventID)
	if err != nil {
		log.Println("error GetAll", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, register)
}

func (api RegisterAPI) GetByUserID(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	userID, _, _ := utils.GetTokenValue(c)
	register, err := api.RegisterRepository.GetRegisterByUserID(userID)
	if err != nil {
		log.Println("error GetAll", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, register)
}

func (api RegisterAPI) ChargeRegEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	userID, _, _ := utils.GetTokenValue(c)
	token := c.PostForm("token")
	price := c.PostForm("price")
	eventID := c.PostForm("event_id")
	regID := c.PostForm("reg_id")
	// log.Printf("created charge: %s\n", token)
	// log.Printf("created charge: %s\n", price)
	// log.Printf("created eventid: %s\n", eventID)
	if amount, err := strconv.ParseInt(price, 10, 64); err == nil {
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Token is null"})
			return
		}

		client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Token is null"})
			return
		}
		charge, create := &omise.Charge{}, &operations.CreateCharge{
			Amount:   amount,
			Currency: "thb",
			Card:     token,
		}

		if err := client.Do(charge, create); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Token is null"})
			return
		}

		err = api.RegisterRepository.AddMerChant(userID, eventID, *charge)
		if err != nil {
			log.Println("error add merchant", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if regID != "" {
			register, err := api.RegisterRepository.GetRegEventByID(regID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			register.TotalPrice = float64(charge.Amount / 100)
			register.PaymentType = config.PAYMENT_CREDIT_CARD
			register.Status = config.PAYMENT_SUCCESS
			register.OrderID = charge.ID
			register.UpdatedAt = time.Now()

			err = api.RegisterRepository.EditRegister(regID, register)
			if err != nil {
				log.Println("error update register payment success", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}

		appG.Response(http.StatusOK, e.SUCCESS, charge)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token or price is null"})
	}

}

func (api RegisterAPI) GetRegEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	userID, _, _ := utils.GetTokenValue(c)
	regID := c.Param("regID")
	register, err := api.RegisterRepository.GetRegisterByUserAndEvent(userID, regID)
	if err != nil {
		log.Println("error GetAll", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, register)
}

func (api RegisterAPI) GetAll(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	register, err := api.RegisterRepository.GetRegisterAll()
	if err != nil {
		log.Println("error GetAll", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, register)
}

func (api RegisterAPI) AddRegister(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	userID, _ := oauth.GetValuesToken(c)
	ownerObjectID, _ := primitive.ObjectIDFromHex(userID)

	var json model.RegisterAdd

	json.Regs.UserID = ownerObjectID

	//categoryObjectID, _ := primitive.ObjectIDFromHex(json.Category.)

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	json.Regs.TotalPrice = utils.ToFixed(json.Regs.TotalPrice, 2)

	registerID, err := api.RegisterRepository.AddRegister(json)
	if err != nil {
		log.Println("error AddRegister", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, registerID)
}

func (api RegisterAPI) AddRaceRegister(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	userID, _, _ := utils.GetTokenValue(c)
	ownerObjectID, _ := primitive.ObjectIDFromHex(userID)

	var json model.Register

	json.UserID = ownerObjectID

	//categoryObjectID, _ := primitive.ObjectIDFromHex(json.Category.)

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	json.TotalPrice = utils.ToFixed(json.TotalPrice, 2)

	registerID, err := api.RegisterRepository.AddRaceRegister(json)
	if err != nil {
		log.Println("error AddRegister", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, registerID)
}

func (api RegisterAPI) EditRegister(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	id := c.Param("id")
	log.Printf("[info] id %s", id)

	var json model.Register

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.RegisterRepository.EditRegister(id, json)
	if err != nil {
		log.Println("error AddEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func (api RegisterAPI) SendSlipTransfer(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	id := c.Param("id")
	log.Printf("[info] id %s", id)

	var json model.SlipTransfer

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.RegisterRepository.NotifySlipRegister(id, json)
	if err != nil {
		log.Println("error SendSlipTransfer", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api RegisterAPI) AdminUpSlip(c *gin.Context) {
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

	var json model.SlipTransfer

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := api.RegisterRepository.AdminNotifySlipRegister(id, json)
	if err != nil {
		log.Println("error SendSlipTransfer", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api RegisterAPI) SendMailRegister(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	var json model.EmailTemplateData
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mail.SendRegEventMail(json)

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api RegisterAPI) CountRegisterEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	eventID := c.Param("eventID")
	count, err := api.RegisterRepository.CountByEvent(eventID)
	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, count)
}

func (api RegisterAPI) CheckUserRegisterEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	userID, _, _ := utils.GetTokenValue(c)
	eventID := c.Param("eventID")
	check, err := api.RegisterRepository.CheckUserRegisterEvent(eventID, userID)
	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, check)
}

func (api RegisterAPI) SendMailRegister2(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	var json model.EmailTemplateData
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mail.SendRegEventMail(json)

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api RegisterAPI) GetMyRegEventActivate(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	userID, _, _ := utils.GetTokenValue(c)
	events, err := api.RegisterRepository.GetRegisterActivateEvent(userID)
	if err != nil {
		log.Println("error GetAll", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, events)
}

func (api RegisterAPI) GetReport(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var form model.DataRegisterRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := api.RegisterRepository.GetRegisterReport(form)

	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, report)
}

func (api RegisterAPI) GetReportAll(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var form model.DataRegisterRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := api.RegisterRepository.GetRegisterReportAll(form)

	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, report)
}

func (api RegisterAPI) FindPersonRegEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var form model.DataRegisterRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := api.RegisterRepository.FindPersonRegEvent(form)

	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, report)
}

func (api RegisterAPI) UpdateStatus(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	userID, _, _ := utils.GetTokenValue(c)

	isAdmin := utils.ISAdmin(c)

	if isAdmin != true {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you do not permission access."})
		return
	}

	var form model.UpdayeRegisterStatusRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registerID := form.RegisterID
	status := form.Status
	err := api.RegisterRepository.UpdateStatusRegister(registerID, status, userID)

	if err != nil {
		log.Println("error UpdateRegisterStatus", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
