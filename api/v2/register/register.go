package registerV2

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/mail"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository/v2"
	repoV1 "thinkdev.app/think/runex/runexapi/repository"
	"thinkdev.app/think/runex/runexapi/utils"
)

// const (

// 	OmisePublicKey = "pkey_test_5i6ivm4cotoab601bfr"
// 	OmiseSecretKey = "skey_test_5h9vov2hqe9iv55o8tu"
// )

//RegisterAPI repo struct
type RegisterAPI struct {
	RegisterRepository repository.RegisterRepository
}

// GetByEvent api godoc
// @Summary Get register by event id
// @Description get register API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.Register}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/:{eventID} [get]
func (api RegisterAPI) GetByEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	eventID := c.Param("eventID")
	register, err := api.RegisterRepository.GetRegisterByEvent(eventID)
	if err != nil {
		log.Println("error GetAll", err.Error())
		appG.Response(http.StatusInternalServerError, http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, register)
}

// GetByUserID api godoc
// @Summary Get register by user id
// @Description get register API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.RegisterV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/myRegEvent [get]
func (api RegisterAPI) GetByUserID(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	register, err := api.RegisterRepository.GetRegisterByUserID(userID)
	if err != nil {
		log.Println("error GetAll", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", register)
}

// ChargeRegEvent api doc
// @Summary payment charge register event by register id
// @Description payment charge register event API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.RegisterChargeRequest true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/payment [post]
func (api RegisterAPI) ChargeRegEvent(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)

	var json model.RegisterChargeRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		appG.Response(http.StatusBadRequest, err.Error(), err.Error())
		return
	}
	var amount int64 = int64(json.Price) * 100
	if json.TokenOmise == "" {
		appG.Response(http.StatusBadRequest, "Token invalid", "Token invalid")
		return
	}

	OmisePublicKey := viper.GetString("omise.OmisePublicKey")
	OmiseSecretKey := viper.GetString("omise.OmiseSecretKey")

	client, err := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if err != nil {
		appG.Response(http.StatusInternalServerError, err.Error(), nil)
		return
	}
	charge, create := &omise.Charge{}, &operations.CreateCharge{
		Amount:   amount,
		Currency: "thb",
		Card:     json.TokenOmise,
	}

	if err := client.Do(charge, create); err != nil {
		log.Println(err)
		appG.Response(http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	err = api.RegisterRepository.AddMerChant(userID, json.EventCode, json.RegID, *charge, json.OrderID)
	if err != nil {
		log.Println("error add merchant", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	appG.Response(http.StatusOK, "success", charge)

}

// SendSlip api doc
// @Summary payment charge register event by register id
// @Description payment charge register event API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.RegisterAttachSlipRequest true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/sendSlip [post]
func (api RegisterAPI) SendSlip(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)

	var json model.RegisterAttachSlipRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		appG.Response(http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	objectID, err := primitive.ObjectIDFromHex(userID)

	success := repository.SlipUpdatePaymentStatus(json.RegID, json.EventCode, objectID, config.PAYMENT_WAITING_APPROVE, json.PaymentType, json.Image)
	if !success {
		log.Println("error add merchant", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), err.Error())
		return
	}
	appG.Response(http.StatusOK, "success", nil)
}

// AdminSendSlip api doc
// @Summary payment charge register event by register id
// @Description payment charge register event API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.AdminAttachSlipRequest true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/adminUpdatePayment [post]
func (api RegisterAPI) AdminSendSlip(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var json model.AdminAttachSlipRequest
	token := c.GetHeader("token")
	key := viper.GetString("public.token")
	if err := c.ShouldBindJSON(&json); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		return
	}
	if token != key {
		res.Response(http.StatusNotFound, "", nil)
		return
	}
	log.Println(json)
	success := repository.AdminUpdateSlipPaymentStatus(json.RegID, json.EventCode, json.UserID, json.Status, json.PaymentType, json.Image)
	if !success {
		res.Response(http.StatusInternalServerError, "Update failed!", nil)
		return
	}
	res.Response(http.StatusOK, "success", nil)
}

// GetRegEvent api doc
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

// GetAll api godoc
// @Summary Get register all
// @Description get register all API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.RegisterV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/myRegEvent [get]
func (api RegisterAPI) GetAll(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	register, err := api.RegisterRepository.GetRegisterAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, "success", register)
}

// AddRegister api godoc
// @Summary add register event
// @Description save register API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.RegisterRequest true "payload"
// @Success 200 {object} response.Response{data=model.RegisterV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/add [post]
func (api RegisterAPI) AddRegister(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)

	userID, _ := oauth.GetValuesToken(c)
	ownerObjectID, _ := primitive.ObjectIDFromHex(userID)

	var json model.RegisterRequest

	json.Regs.UserID = ownerObjectID

	//categoryObjectID, _ := primitive.ObjectIDFromHex(json.Category.)

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err)
		res.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}

	json.Regs.TotalPrice = utils.ToFixed(json.Regs.TotalPrice, 2)

	registerID, err := api.RegisterRepository.AddRegister(json)
	if err != nil {
		log.Println("error AddRegister", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "ลงทะเบียนสำเร็จ", registerID)
}

// AddTeamRegister api godoc
// @Summary add register event
// @Description save register API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.AddTeamRequest true "payload"
// @Success 200 {object} response.Response{data=model.RegisterV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/team [post]
func (api RegisterAPI) AddTeamRegister(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)

	userID, _ := oauth.GetValuesToken(c)
	ownerObjectID, _ := primitive.ObjectIDFromHex(userID)

	var json model.AddTeamRequest

	//categoryObjectID, _ := primitive.ObjectIDFromHex(json.Category.)

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err)
		res.Response(http.StatusBadRequest, err.Error(), nil)
		return
	}

	if check, _ := repository.CheckTeamLead(json, ownerObjectID); !check {
		res.Response(http.StatusBadRequest, "Your not team lead.", nil)
		return
	}

	if _, check, _ := repository.CheckRunerInTeam(json); !check {
		res.Response(http.StatusBadRequest, "Team is full.", nil)
		return
	}

	check, err := repository.CheckSameUser(json)
	if err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		return
	}
	if check {
		res.Response(http.StatusBadRequest, "User in team.", "")
		return
	}

	json.Regs.TotalPrice = 0
	json.Regs.UserID = json.TeamUserID
	register, err := repository.AddTeamRegister(json)
	if err != nil {
		log.Println("error AddRegister", err.Error())
		res.Response(http.StatusInternalServerError, "Cannot add to team.", gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", register)
}

//AddRaceRegister api doc
func (api RegisterAPI) AddRaceRegister(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	userID, _ := oauth.GetValuesToken(c)
	ownerObjectID, _ := primitive.ObjectIDFromHex(userID)

	var json model.RegisterRequest

	json.Regs.UserID = ownerObjectID

	//categoryObjectID, _ := primitive.ObjectIDFromHex(json.Category.)

	if err := c.ShouldBindJSON(&json); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	json.Regs.TotalPrice = utils.ToFixed(json.Regs.TotalPrice, 2)

	registerID, err := api.RegisterRepository.AddRaceRegister(json)
	if err != nil {
		log.Println("error AddRegister", err.Error())
		appG.Response(http.StatusInternalServerError, http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, registerID)
}

//EditRegister api doc
func (api RegisterAPI) EditRegister(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	id := c.Param("id")
	log.Printf("[info] id %s", id)

	var json model.RegisterRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := api.RegisterRepository.EditRegister(id, json)
	if err != nil {
		log.Println("error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

//SendSlipTransfer api doc
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
		appG = response.Gin{C: c}
	)
	eventID := c.Param("eventID")
	count, err := api.RegisterRepository.CountByEvent(eventID)
	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, "success", count)
}

// CheckUserRegisterEvent api godoc
// @Summary check register by user id
// @Description check register API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/checkUserRegisterEvent/:{eventID} [get]
func (api RegisterAPI) CheckUserRegisterEvent(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	eventID := c.Param("eventID")
	check, err := api.RegisterRepository.CheckUserRegisterEvent(eventID, userID)
	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", gin.H{"is_reg": check})
}

// CheckUserRegisterEventCode api godoc
// @Summary check register by user id and event code
// @Description check register API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/checkRegEventCode/:{code} [get]
func (api RegisterAPI) CheckUserRegisterEventCode(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	code := c.Param("code")
	check, err := api.RegisterRepository.CheckUserRegisterEventCode(code, userID)
	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", gin.H{"is_reg": check})
}

// CheckUserRegisterER api godoc
// @Summary check register by user id and event code
// @Description check register API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/checkRegER [post]
func (api RegisterAPI) CheckUserRegisterER(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	var json model.CheckRegisterERRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}
	check, err := repository.CheckUserRegisterER(json.EventCode, userID, json.TicketID, json.CitycenID)
	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", gin.H{"is_reg": check})
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

// GetMyRegEventActivate api godoc
// @Summary get register payment success by user id
// @Description get register payment success API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.RegisterV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/myRegEventActivate [get]
func (api RegisterAPI) GetMyRegEventActivate(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	events, err := api.RegisterRepository.GetRegisterActivateEvent(userID)
	if err != nil {
		log.Println("error GetAll activate", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, "success", events)
}

// GetRegEventDashboard api godoc
// @Summary get register payment success by user id
// @Description get register payment success API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.RegEventDashboardRequest true "payload"
// @Success 200 {object} response.Response{data=model.RegisterV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/myRegEventActivate [get]
func (api RegisterAPI) GetRegEventDashboard(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var json model.RegEventDashboardRequest

	//categoryObjectID, _ := primitive.ObjectIDFromHex(json.Category.)

	if err := c.ShouldBindJSON(&json); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}

	reg, err := repository.GetRegEventDashboard(json)
	if err != nil {
		log.Println("error GetAll activate", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", reg)
}

// GetRegEventFromEventer api godoc
// @Summary get register datas 's eventer and admin
// @Description get register datas 's eventer and admin
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.RegisterV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/regsEvent/:{eventID}  [get]
func (api RegisterAPI) GetRegEventFromEventer(c *gin.Context) {
	var res = response.Gin{C: c}
	userID, role := oauth.GetValuesToken(c)
	eventID := c.Param("eventID")

	if !repository.IsOwner(eventID, userID) && role != "ADMIN" {
		res.Response(http.StatusUnauthorized, "You do not have access to the information.", nil)
		return
	}

	datas, err := api.RegisterRepository.GetRegisterByEvent(eventID)
	if err != nil {
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res.Response(http.StatusOK, "success", datas)
}

// GetRegEventFromOwner api godoc
// @Summary get register datas 's eventer and admin
// @Description get register datas 's eventer and admin
// @Consume application/x-www-form-urlencoded
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.OwnerRequest true "payload"
// @Success 200 {object} response.Response{data=[]model.RegisterV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/regsEvent  [post]
func (api RegisterAPI) GetRegEventFromOwner(c *gin.Context) {
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
	// if !repository.IsOwner(form.EventCode, form.OwnerID) {
	// 	res.Response(http.StatusUnauthorized, "You do not have access to the information.", nil)
	// 	return
	// }

	datas, err := api.RegisterRepository.GetRegisterByEvent(form.EventCode)
	if err != nil {
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res.Response(http.StatusOK, "success", datas)
}

func (api RegisterAPI) GetReport(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var form model.DataRegisterRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		appG.Response(http.StatusBadRequest, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := api.RegisterRepository.GetRegisterReport(form)

	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		appG.Response(http.StatusInternalServerError, http.StatusInternalServerError, gin.H{"message": err.Error()})
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
		appG.Response(http.StatusBadRequest, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := api.RegisterRepository.GetRegisterReportAll(form)

	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		appG.Response(http.StatusInternalServerError, http.StatusInternalServerError, gin.H{"message": err.Error()})
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
		appG.Response(http.StatusBadRequest, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := api.RegisterRepository.FindPersonRegEvent(form)

	if err != nil {
		log.Println("error CountRegisterEvent", err.Error())
		appG.Response(http.StatusInternalServerError, http.StatusBadRequest, gin.H{"message": err.Error()})
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

// PaymentHook api godoc
// @Summary payment register hook
// @Description payment register hook data API calls
// @Consume application/x-www-form-urlencoded
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.SCBPayment true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/payment_hook [post]
func (api RegisterAPI) PaymentHook(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var form model.SCBPayment
	token := c.GetHeader("token")
	key := viper.GetString("public.token")
	// body, _ := ioutil.ReadAll(c.Request.Body)
	// println(string(body))

	if token != key {
		res.Response(http.StatusNotFound, "", nil)
		return
	}
	if err := c.ShouldBind(&form); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := repository.PaymentWithSCB(form)
	if err != nil {
		log.Println("error hook payment", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res.Response(http.StatusOK, "success", nil)
}

// RegEventUpdateUserInfo api godoc
// @Summary get register payment success by user id
// @Description get register payment success API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.RegUpdateUserInfoRequest true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/regUpdateUserInfo [post]
func (api RegisterAPI) RegEventUpdateUserInfo(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var json model.RegUpdateUserInfoRequest
	userID, _ := oauth.GetValuesToken(c)
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	if err := c.ShouldBindJSON(&json); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}

	err := repository.UpdateRegisterUserInfo(userObjectID, json)
	if err != nil {
		log.Println("error update user info register", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", nil)
}

// EditUserInfo api godoc
// @Summary get register payment success by user id
// @Description get register payment success API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags register
// @Accept  application/json
// @Produce application/json
// @Param payload body model.RegUpdateUserInfoRequest true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /register/editUserInfo [post]
func (api RegisterAPI) EditUserInfo(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var json model.RegUpdateUserInfoRequest
	userID, _ := oauth.GetValuesToken(c)
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	if err := c.ShouldBindJSON(&json); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}

	err := repository.EditUserOption(userObjectID, json)
	if err != nil {
		log.Println("error update user info register", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	err = repoV1.UpdateUserInfoActivity(json.TicketOption.UserOption, json.EventCode, userObjectID, json.RegID)
	if err != nil {
		log.Println("Update userinfo activity fail")
	}

	res.Response(http.StatusOK, "success", nil)
}
