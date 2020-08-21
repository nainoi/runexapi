package importdata

import (
	// "fmt"
	// "log"
	// "net/http"

	// "thinkdev.app/think/runex/runexapi/model"
	// "thinkdev.app/think/runex/runexapi/pkg/app"
	// "thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/repository"
	// "github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/360EntSecGroup-Skylar/excelize"
)

type ImportDataAPI struct {
	ImportDataRepository repository.ImportDataRepository
}

// func (api ImportDataAPI) ImportExcel(c *gin.Context) {
// 	var (
// 		appG = app.Gin{C: c}
// 	)

// 	event := c.Param("event")
// 	fmt.Println(event)

// 	f, err := excelize.OpenFile("./WHJJ3_DATA.xlsx")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	rows, err := f.Rows("Sheet1")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for rows.Next() {
// 		row, err2 := rows.Columns()
// 		if err2 != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("%s\t%s\n", row[1], row[3]) // Print values in columns B and D

// 		//var address = []model.Address
// 		var userImport = model.ExcelUserForm{
// 			Email:     row[4],
// 			FullName:  row[1],
// 			FirstName: row[2],
// 			Phone:     row[3],
// 			Address: []model.Address{
// 				model.Address{
// 					Address: row[9],
// 				},
// 			},
// 		}

// 		userID, exists, err := api.ImportDataRepository.ExistUserByEmail(userImport.Email)
// 		if exists {

// 			existsRegister, err := api.ImportDataRepository.ExistRegisterByUserAndEvent(userID, event)
// 			if !existsRegister {

// 				ownerObjectID, _ := primitive.ObjectIDFromHex(userID)
// 				eventObjectID, _ := primitive.ObjectIDFromHex(event)

// 				var registerImport = model.Register{
// 					UserID:  ownerObjectID,
// 					EventID: eventObjectID,
// 				}

// 				err2 := api.ImportDataRepository.AddRegister(registerImport)

// 				if err != nil {
// 					log.Fatal(err)
// 				}

// 				if err2 != nil {

// 					log.Fatal(err2)
// 				}
// 			}
// 		} else {

// 			userID, err := api.ImportDataRepository.AddUser(userImport)

// 			ownerObjectID, _ := primitive.ObjectIDFromHex(userID)
// 			eventObjectID, _ := primitive.ObjectIDFromHex(event)

// 			var registerImport = model.Register{
// 				UserID:  ownerObjectID,
// 				EventID: eventObjectID,
// 			}

// 			err2 := api.ImportDataRepository.AddRegister(registerImport)

// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			if err2 != nil {
// 				log.Fatal(err2)
// 			}
// 		}

// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	appG.Response(http.StatusOK, e.SUCCESS, nil)
// }
