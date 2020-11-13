package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/repository/v2"
)

//ConfigAPI struct for user repository
type ConfigAPI struct {
	ConfigRepo repository.ConfigRepository
}

// GetConfig godoc
// @Summary Get config app
// @Description get config info API calls
// @Tags config
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.ConfigModel}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /config [get]
func (api ConfigAPI) GetConfig(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	config := api.ConfigRepo.GetConfig()
	res.Response(http.StatusOK, "success", config)
	return
}

func (api ConfigAPI) GetTest(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	res.Response(http.StatusOK, "success", gin.H{"msg": "succes"})
	return
}
