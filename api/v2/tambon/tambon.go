package tambon

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

// SearchZipcode api godoc
// @Summary get tambon datas list by zipcode
// @Description get tambon datas list by zipcode
// @Consume application/x-www-form-urlencoded
// @Tags tambon
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.Tambon}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /tambon/:{zipcode}  [get]
func SearchZipcode(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	zipcode := c.Param("zipcode")
	tambons := repository.TambonByZipcode(zipcode)
	res.Response(http.StatusOK, "success", tambons)
	return
}

// SearchProvince api godoc
// @Summary get tambon datas list by zipcode
// @Description get tambon datas list by zipcode
// @Consume application/x-www-form-urlencoded
// @Tags tambon
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.Tambon}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /province/:{province}  [get]
func SearchProvince(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	province := c.Param("province")
	tambons := repository.Province(province)
	res.Response(http.StatusOK, "success", tambons)
	return
}

// SearchAmphoe api godoc
// @Summary get tambon datas list by zipcode
// @Description get tambon datas list by zipcode
// @Consume application/x-www-form-urlencoded
// @Tags tambon
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.Tambon}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /amphoe/:{amphoe}  [get]
func SearchAmphoe(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	amphoe := c.Param("amphoe")
	tambons := repository.Amphoe(amphoe)
	res.Response(http.StatusOK, "success", tambons)
	return
}

// SearchDistrict api godoc
// @Summary get tambon datas list by zipcode
// @Description get tambon datas list by zipcode
// @Consume application/x-www-form-urlencoded
// @Tags tambon
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.Tambon}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /district/:{district}  [get]
func SearchDistrict(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	district := c.Param("district")
	tambons := repository.District(district)
	res.Response(http.StatusOK, "success", tambons)
	return
}

// All api godoc
// @Summary get tambon datas list all
// @Description get tambon datas list all
// @Consume application/x-www-form-urlencoded
// @Tags tambon
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=[]model.Tambon}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /tambon  [get]
func All(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	tambons := repository.TambonAll()
	res.Response(http.StatusOK, "success", tambons)
	return
}
