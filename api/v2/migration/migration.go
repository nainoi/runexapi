package migration

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
)

type MigrationAPI struct {
	MigrationRepository repository.MigrationRepository
}

func (api MigrationAPI) MigrateWorkout(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	newCollection := c.Param("newCollection")

	// integer, fraction := math.Modf(24.3488998413086)
	// var modTimeMin int = int(integer * 60)
	// var modTimeSec int = int(fraction * 100)
	// fmt.Println(modTimeMin)
	// fmt.Println(modTimeSec)
	// durTime := time.Duration(modTimeMin+modTimeSec) * time.Second
	// modTime := time.Now().Round(0).Add(-(durTime))
	// since := time.Since(modTime)
	// fmt.Println(since)
	// durStr := fmtDuration(since)
	// fmt.Println(durStr)

	err2 := api.MigrationRepository.MigrateWorkout(newCollection)
	if err2 != nil {
		log.Println("error MigrateWorkout", err2.Error())
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
