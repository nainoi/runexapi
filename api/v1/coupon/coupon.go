package coupon

import (
	"log"
	"net/http"
	"time"

	"bitbucket.org/suthisakch/runex/model"
	"bitbucket.org/suthisakch/runex/pkg/app"
	"bitbucket.org/suthisakch/runex/pkg/e"
	"bitbucket.org/suthisakch/runex/repository"
	"github.com/gin-gonic/gin"
)

type CouponAPI struct {
	CouponRepository repository.CouponRepository
}

func (api CouponAPI) CreateCoupon(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	var json model.Coupon

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err2 := api.CouponRepository.ExistByCode(json.CouponCode)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	if exists {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_COUPON, nil)
		return
	}

	couponID, err := api.CouponRepository.CreateCoupon(json)
	if err != nil {
		log.Println("error CreateCoupon", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, couponID)
}

func (api CouponAPI) EditCoupon(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	id := c.Param("id")
	log.Printf("[info] id %s", id)
	var json model.EditCouponForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err2 := api.CouponRepository.ExistByCodeForEdit(json.CouponCode, id)

	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	if exists {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_COUPON, nil)
		return
	}

	err := api.CouponRepository.EditCoupon(id, json)
	if err != nil {
		log.Println("error EditCoupon", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func (api CouponAPI) DeleteCoupon(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	id := c.Param("id")
	log.Printf("[info] id %s", id)

	err := api.CouponRepository.DeleteCouponByID(id)
	if err != nil {
		log.Println("error DeleteCoupony", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func (api CouponAPI) GetAll(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	coupon, err := api.CouponRepository.GetCouponAll()
	if err != nil {
		log.Println("error Coupon", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, coupon)

}

func (api CouponAPI) ValidateCode(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	var json model.ValidateCoupon

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err2 := api.CouponRepository.ExistByCode(json.CouponCode)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	if !exists {
		appG.Response(http.StatusBadRequest, e.ERROR_NOT_EXIST_COUPON, nil)
		return
	}
	coupon, err := api.CouponRepository.GetCouponByCode(json.CouponCode)
	if err != nil {
		log.Println("error Get Code", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	today := time.Now()
	startDate := coupon.StartDate
	endDate := coupon.EndDate
	if (startDate.Before(today) && endDate.After(today)) && coupon.Active {
		appG.Response(http.StatusOK, e.SUCCESS, coupon)
		return
	} else {
		appG.Response(http.StatusBadRequest, e.ERROR_COUPON_FAIL_EXPIRE, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, coupon)

}

func (api CouponAPI) GetByCode(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)
	code := c.Param("code")
	log.Printf("[info] code %s", code)

	coupon, err := api.CouponRepository.GetCouponByCode(code)
	if err != nil {
		log.Println("error Get Code", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, coupon)
}
