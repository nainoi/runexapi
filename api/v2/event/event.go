package eventV2

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	guuid "github.com/google/uuid"

	//handle_user "thinkdev.app/think/runex/runexapi/api/v1/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"thinkdev.app/think/runex/runexapi/middleware/oauth"
	"thinkdev.app/think/runex/runexapi/repository/v2"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
)

type EventAPI struct {
	EventRepository repository.EventRepository
}

type EventStatus struct {
	Status string `json:"status" bson:"status" binding:"required"`
}

func (api EventAPI) AddEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	userID, _ := oauth.GetValuesToken(c)
	ownerObjectID, _ := primitive.ObjectIDFromHex(userID)

	var json model.EventV2

	json.OwnerID = ownerObjectID

	//categoryObjectID, _ := primitive.ObjectIDFromHex(json.Category.)

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err2 := api.EventRepository.ExistByName(json.Name)

	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	if exists {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_EVENT, nil)
		return
	}

	eventID, err := api.EventRepository.AddEvent(json)
	if err != nil {
		log.Println("error AddEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, eventID)

}

func (api EventAPI) MyEvent(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	userID, _ := oauth.GetValuesToken(c)
	event, err := api.EventRepository.GetEventByUser(userID)
	if err != nil {
		log.Println("error AddEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, event)

}

func (api EventAPI) GetAll(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	event, err := api.EventRepository.GetEventAll()
	if err != nil {
		log.Println("error AddEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, event)

}

func (api EventAPI) GetAllActive(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	event, err := api.EventRepository.GetEventActive()
	if err != nil {
		log.Println("error AddEvent", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, event)

}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slugText := slug.MakeLang(json.Slug, "th")
	//slug.Make(json.Slug)
	//slug.MakeLang("Diese & Dass", "de")
	//slug.Make(json.Slug)
	fmt.Println(slugText) // Will print: "hello-world-khello-vorld"
	existsSlug, err2 := api.EventRepository.ValidateBySlug(json.Slug)

	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	if !existsSlug {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_EVENT_SLUG, json.Slug)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, json.Slug)

}

func (api EventAPI) GetBySlug(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	slug := c.Param("slug")
	log.Printf("[info] id %s", slug)

	existsSlug, err2 := api.EventRepository.ValidateBySlug(slug)
	log.Printf("[ValidateBySlug] id %s", slug)
	if existsSlug {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_EVENT_SLUG, slug)
		return
	}

	event, err := api.EventRepository.GetEventBySlug(slug)
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
