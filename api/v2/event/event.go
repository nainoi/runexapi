package eventV2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	//handle_user "thinkdev.app/think/runex/runexapi/api/v1/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/repository/v2"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
)

// EventAPI reference
type EventAPI struct {
	EventRepository repository.EventRepository
}

// EventStatus struct
type EventStatus struct {
	Status string `json:"status" bson:"status" binding:"required"`
}

// AddEvent api godoc
// @Summary add new event
// @Description add new event API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags event
// @Accept  application/json
// @Produce application/json
// @Param payload body model.EventV2 true "payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /event [post]
func (api EventAPI) AddEvent(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)

	userID, _ := oauth.GetValuesToken(c)
	ownerObjectID, _ := primitive.ObjectIDFromHex(userID)

	var json model.EventV2

	//categoryObjectID, _ := primitive.ObjectIDFromHex(json.Category.)

	if err := c.ShouldBindJSON(&json); err != nil {
		appG.Response(http.StatusBadRequest, err.Error(), gin.H{"error": err.Error()})
		return
	}

	json.OwnerID = ownerObjectID

	exists, err2 := api.EventRepository.ExistByName(json.Name)

	if err2 != nil {
		appG.Response(http.StatusInternalServerError, err2.Error(), gin.H{"message": err2.Error()})
		return
	}

	if exists {
		appG.Response(http.StatusBadRequest, "event not exits", nil)
		return
	}

	eventID, err := api.EventRepository.AddEvent(json)
	if err != nil {
		log.Println("error AddEvent", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, "success", eventID)

}

// MyEvent api godoc
// @Summary Get my event
// @Description get event owner API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags event
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.EventV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /event/myEvent [get]
func (api EventAPI) MyEvent(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	event, err := api.EventRepository.GetEventByUser(userID)
	if err != nil {
		log.Println("error get my event", err.Error())
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	res.Response(http.StatusOK, "success", event)

}

// GetAll api godoc
// @Summary Get event all
// @Description get Get event all API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags event
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.EventList}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /event/all [get]
func GetAll(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	urlS := fmt.Sprintf("https://events-api.thinkdev.app/events")
	var bearer = "Bearer olcgZVpqDXQikRDG"
	//reqURL, _ := url.Parse(urlS)
	req, err := http.NewRequest("GET", urlS, nil)
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//log.Println(err)
			appG.Response(http.StatusOK, err.Error(), model.EventList{})
			c.Abort()
			return
		}
		var events model.EventLists
		err = json.Unmarshal(body, &events)
		if err != nil {
			//log.Println(err)
			appG.Response(http.StatusOK, err.Error(), model.EventList{})
			c.Abort()
			return
		}
		if events.Events != nil {
			appG.Response(http.StatusOK, "success", events.Events)
			c.Abort()
			return
		}
		appG.Response(http.StatusOK, "success", model.EventList{})
		c.Abort()
		return
	}

	appG.Response(http.StatusInternalServerError, err.Error(), nil)

}

// GetAllActive api godoc
// @Summary Get event active status
// @Description get Get event active status API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags event
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.EventV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /event/active [get]
func (api EventAPI) GetAllActive(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	urlS := fmt.Sprintf("https://events-api.thinkdev.app/events")
	var bearer = "Bearer olcgZVpqDXQikRDG"
	//reqURL, _ := url.Parse(urlS)
	req, err := http.NewRequest("GET", urlS, nil)
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		var koaObject model.KaoObject
		err = json.Unmarshal(body, &koaObject)
		if err != nil {
			log.Println(err)
		}
		appG.Response(http.StatusOK, "success", koaObject)
		c.Abort()
		return
	}

	appG.Response(http.StatusInternalServerError, err.Error(), nil)
	event, err := api.EventRepository.GetEventActive()
	if err != nil {
		log.Println("error AddEvent", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, "success", event)

}

// GetDetail api godoc
// @Summary Get event detail
// @Description get Get event detail API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags event
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.EventList}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /event/detail/{code} [get]
func GetDetail(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	code := c.Param("code")
	if code == "" {
		appG.Response(http.StatusInternalServerError, "code not found", nil)
		c.Abort()
		return
	}
	event, err := DetailEventByCode(code)

	if err != nil {
		appG.Response(http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return

	}

	appG.Response(http.StatusOK, "success", event)
	c.Abort()
	return

}

//DetailEventByCode go doc
//Description get Get event detail API calls to event runex
func DetailEventByCode(code string) (model.EventData, error) {
	urlS := fmt.Sprintf("https://events-api.thinkdev.app/event/%s", code)
	var bearer = "Bearer olcgZVpqDXQikRDG"
	//reqURL, _ := url.Parse(urlS)
	req, err := http.NewRequest("GET", urlS, nil)
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	var event model.EventData

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//log.Println(err)
			return event, err
		}

		err = json.Unmarshal(body, &event)
		if err != nil {
			//log.Println(err)
			return event, err
		}
		return event, err
	}

	return event, err

}

//GetByStatus go doc
func (api EventAPI) GetByStatus(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	status := c.Param("status")
	log.Printf("[info] status %s", status)

	event, err := api.EventRepository.GetEventByStatus(status)
	if err != nil {
		log.Println("error AddEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, event)

}

/*
func (api EventAPI) GetByID(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	log.Printf("[info] id %s", id)

	event, err := api.EventRepository.GetEventByID(id)
	if err != nil {
		log.Println("error AddEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user, err2 := api.EventRepository.GetUserEvent(event.OwnerID.Hex())
	if err2 != nil {
		log.Println("error owner event", err2.Error())
	}
	//data := EventRes{event, user}
	appG.Response(http.StatusOK, e.SUCCESS, gin.H{"event": event, "user": user})

}

func (api EventAPI) EditEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	log.Printf("[info] id %s", id)
	var json model.EventV2
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err2 := api.EventRepository.ExistByNameForEdit(json.Name, id)

	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	if exists {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_EVENT, nil)
		return
	}

	err := api.EventRepository.EditEvent(id, json)
	if err != nil {
		log.Println("error AddEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func (api EventAPI) DeleteEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	log.Printf("[info] id %s", id)

	err := api.EventRepository.DeleteEventByID(id)
	if err != nil {
		log.Println("error DeleteEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func (api EventAPI) UploadImage(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	id := c.Param("id")
	log.Printf("[info] id %s", id)

	file, header, err := c.Request.FormFile("upload")
	filename := header.Filename
	fmt.Println(filename)

	uniqidFilename := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", uniqidFilename.String())

	pathDir := "./upload/image/event/"
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		os.MkdirAll(pathDir, os.ModePerm)
	}

	out, err := os.Create(pathDir + uniqidFilename.String() + ".png")

	path := "/upload/image/event/" + uniqidFilename.String() + ".png"
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
		//log.Fatal(err)
	}

	err2 := api.EventRepository.UploadCoverEvent(id, path)
	if err2 != nil {
		log.Println("error UploadImage", err2.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, out)
}

func (api EventAPI) AddProduct(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	var json model.ProduceEvent

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err2 := api.EventRepository.AddProductEvent(id, json)
	if err2 != nil {
		log.Println("error AddProductEvent", err2.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, productID)
}

func (api EventAPI) EditProduct(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	var json model.ProduceEvent

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err2 := api.EventRepository.EditProductEvent(id, json)
	if err2 != nil {
		log.Println("error EditProduct", err2.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api EventAPI) GetProductEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")

	product, err := api.EventRepository.GetProductByEventID(id)
	if err != nil {
		log.Println("error GetProductEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, product)
}

func (api EventAPI) DeleteProductEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	productID := c.Param("productID")

	log.Printf("[info] id %s", id)
	log.Printf("[info] productID %s", productID)

	err := api.EventRepository.DeleteProductEvent(id, productID)
	if err != nil {
		log.Println("error GetProductEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api EventAPI) AddTicket(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	var json model.TicketEventV2

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketID, err2 := api.EventRepository.AddTicketEvent(id, json)
	if err2 != nil {
		log.Println("error AddTicket", err2.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, ticketID)
}

func (api EventAPI) EditTicket(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	var json model.TicketEventV2

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err2 := api.EventRepository.EditTicketEvent(id, json)
	if err2 != nil {
		log.Println("error EditProduct", err2.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api EventAPI) DeleteTicketEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	productID := c.Param("ticketID")

	log.Printf("[info] id %s", id)
	log.Printf("[info] ticketID %s", productID)

	err := api.EventRepository.DeleteTicketEvent(id, productID)
	if err != nil {
		log.Println("error DeleteTicketEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
*/

func (api EventAPI) SearchEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	var json model.SearchEvent
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[info] Term %s", json.Term)
	event, err := api.EventRepository.SearchEvent(json.Term)
	if err != nil {
		log.Println("error SearchEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, event)

}

func (api EventAPI) ValidateSlug(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var json model.Slug

	if err := c.ShouldBindJSON(&json); err != nil {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}

	slugText := slug.MakeLang(json.Slug, "th")
	//slug.Make(json.Slug)
	//slug.MakeLang("Diese & Dass", "de")
	//slug.Make(json.Slug)
	fmt.Println(slugText) // Will print: "hello-world-khello-vorld"
	existsSlug, err2 := api.EventRepository.ValidateBySlug(json.Slug)

	if err2 != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, gin.H{"message": err2.Error()})
		return
	}

	if !existsSlug {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_EVENT_SLUG, json.Slug)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, json.Slug)

}

// GetBySlug api godoc
// @Summary Get my event
// @Description get event owner API calls
// @Consume application/x-www-form-urlencoded
// @Security bearerAuth
// @Tags event
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.EventV2}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /event/GetBySlug/{slug} [get]
func (api EventAPI) GetBySlug(c *gin.Context) {
	var (
		appG = response.Gin{C: c}
	)
	slug := c.Param("slug")
	log.Printf("[info] id %s", slug)

	existsSlug, err2 := api.EventRepository.ValidateBySlug(slug)
	log.Printf("[ValidateBySlug] id %s", slug)
	if existsSlug {
		appG.Response(http.StatusBadRequest, "slug event exit", slug)
		return
	}

	event, err := api.EventRepository.GetEventBySlug(slug)
	if err != nil {
		log.Println("error get Event", err.Error())
		appG.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	user, err2 := api.EventRepository.GetUserEvent(event.OwnerID.Hex())
	if err2 != nil {
		log.Println("error owner event", err2.Error())
	}
	//data := EventRes{event, user}
	appG.Response(http.StatusOK, "success", gin.H{"event": event, "user": user})

}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
	return nil
}
