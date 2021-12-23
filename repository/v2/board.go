package repository

import (
	"context"
	"log"
	"sort"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/config/db"
	"thinkdev.app/think/runex/runexapi/model"
)

// BoardRepository interface
// type BoardRepository interface {
// 	GetBoardByEvent(req model.RankingRequest, userID string) (model.Event, int64, []model.Ranking, []model.Ranking, error)
// }

// // BoardRepositoryMongo db struct
// type BoardRepositoryMongo struct {
// 	ConnectionDB *mongo.Database
// }

const (
	boardCollection = "activityV2"
)

//GetBoardByEvent board ranking
func GetBoardByEvent(req model.RankingRequest, userID string) (model.Event, int64, []model.Ranking, []model.Ranking, error) {
	//var activity model.Activity
	//var activityInfo []model.ActivityInfo
	var event = model.Event{}
	var activities = []model.Ranking{}
	myActivities := []model.Ranking{}

	// objectEventID, err := primitive.ObjectIDFromHex(eventID)
	// if err != nil {
	// 	log.Println(err)
	// 	return event, 0, activities, myActivities, err
	// }
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return event, 0, activities, myActivities, err
	}

	var temps = []model.Ranking{}

	if req.ParentRegID.IsZero() {
		//filter := bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
		filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "ticket.id", Value: req.TicketID}}

		option := options.Find()
		log.Println(userID)
		//filterEvent := bson.D{primitive.E{Key: "_id", Value: objectEventID}}

		event, err = DetailEventByCode(req.EventCode)
		if err != nil {
			return event, 0, activities, myActivities, err
		}
		//option.SetLimit(10)
		count, err := db.DB.Collection(activityCollection).CountDocuments(context.TODO(), filter)
		if count == 0 {
			return event, count, activities, myActivities, err
		}
		//filter = bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
		option.SetSort(bson.D{primitive.E{Key: "total_distance", Value: -1}})
		cur, err := db.DB.Collection(activityCollection).Find(context.TODO(), filter, option)

		if err != nil {
			log.Println(err)
			return event, count, activities, myActivities, err
		}

		n := 0

		for cur.Next(context.TODO()) {
			// if n >= 10 {
			// 	break
			// }

			var activity model.ActivityV2
			var a model.Ranking
			// decode the document
			if err := cur.Decode(&activity); err != nil {
				log.Println(err)
			}
			if err == nil {
				a = model.Ranking{
					UserID:        activity.UserID,
					ActivityInfo:  []model.ActivityInfo{},
					ToTalDistance: activity.ToTalDistance,
					UserInfo:      activity.UserInfo,
					EventCode:     req.EventCode,
					ID:            activity.ID,
					RegID:         activity.RegID,
				}
				if n < 10 {

					a.RankNo = n + 1
					//a.UserInfo = user
					url, _ := GetUserAvatar(a.UserID)
					a.ImageURL = url
					activities = append(activities, a)
				}

				temps = append(temps, a)

				n++
			}
		}
		index := -1
		for n, s := range temps {
			if s.UserID == objectUserID {
				index = n
				break
			}
		}
		//index := SliceIndex(len(activities), func(i int) bool { return activities[i].UserID == objectUserID })
		if index != -1 && index > 9 {

			// if (index - 2) >= 0 {
			// 	// log.Printf("[info] userID %s", userID)
			// 	ranking := temps[index-2]
			// 	if err != nil {
			// 		log.Println(err)
			// 	}
			// 	ranking.RankNo = index - 1
			// 	myActivities = append(myActivities, ranking)
			// }
			if (index - 1) >= 0 {
				ranking := temps[index-1]
				if err != nil {
					log.Println(err)
				}
				ranking.RankNo = index
				url, _ := GetUserAvatar(ranking.UserID)
				ranking.ImageURL = url
				myActivities = append(myActivities, ranking)
			}

			ranking := temps[index]

			ranking.RankNo = index + 1
			url, _ := GetUserAvatar(ranking.UserID)
			ranking.ImageURL = url
			myActivities = append(myActivities, ranking)
			if (index + 1) < len(temps) {
				ranking := temps[index+1]
				ranking.RankNo = index + 2
				url, _ := GetUserAvatar(ranking.UserID)
				ranking.ImageURL = url
				myActivities = append(myActivities, ranking)
			}
			// if (index + 2) < len(temps) {
			// 	ranking.RankNo = index + 3
			// 	myActivities = append(myActivities, ranking)
			// }
		}

		// //activityInfo = activity.ActivityInfo

		// filter = bson.D{{"event_id", objectEventID}}
		// option := options.Find()
		// option.SetLimit(10)
		// option.SetSort(bson.D{{"total_distance", -1}})
		// cur, err := boardMongo.ConnectionDB.Collection(activityCollection).Find(context.TODO(), filter, option)

		return event, count, activities, myActivities, err
	} else {
		filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "ticket.id", Value: req.TicketID}}
		//filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "activities.user_id", Value: objectUserID}}

		event, err = DetailEventByCode(req.EventCode)
		if err != nil {
			return event, 0, activities, myActivities, err
		}
		//option.SetLimit(10)
		count, err := db.DB.Collection(activityCollection).CountDocuments(context.TODO(), filter)
		if count == 0 {
			return event, count, activities, myActivities, err
		}
		//filter = bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
		matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"$and": []interface{}{bson.M{"event_code": req.EventCode}, bson.M{"ticket.id": req.TicketID}}}}}
		//groupStage := bson.D{primitive.E{Key: "$group", Value: bson.A{"_id", "$parent_reg_id"}}}
		groupStage := bson.D{primitive.E{Key: "$group", Value: bson.M{"_id": bson.M{"parent_reg_id": "$parent_reg_id", "event_code": "$event_code", "ticket_id": "$ticket_id"},
			"total_distance": bson.M{"$sum": "$total_distance"},
			"reg_id":         bson.M{"$first": "$reg_id"},
			"user_id":        bson.M{"$first": "$user_id"},
			"user_info":      bson.M{"$first": "$user_info"},
			"parent_reg_id":  bson.M{"$first": "$parent_reg_id"},
			"ticket":         bson.M{"$first": "$ticket"},
		}}}
		sortStage := bson.D{primitive.E{Key: "$sort", Value: bson.M{"total_distance": -1}}}
		type ActivityTemp struct {
			EventCode     string             `json:"event_code" bson:"event_code"`
			Ticket        model.Tickets      `json:"ticket" bson:"ticket"`
			OrderID       string             `json:"order_id" bson:"order_id"`
			RegID         primitive.ObjectID `json:"reg_id" bson:"reg_id"`
			ParentRegID   primitive.ObjectID `json:"parent_reg_id" bson:"parent_reg_id"`
			UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
			ToTalDistance float64            `json:"total_distance" bson:"total_distance"`
			UserInfo      model.UserOption   `json:"user_info" bson:"user_info"`
		}
		cur, err := db.DB.Collection(activityCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage, sortStage})
		if err != nil {
			log.Println(err)
			return event, count, activities, myActivities, err
		}

		n := 0
		for cur.Next(context.TODO()) {
			// if n >= 10 {
			// 	break
			// }
			// log.Println(cur)
			//var activity model.ActivityV2
			var t ActivityTemp
			var a model.Ranking
			// decode the document
			if err := cur.Decode(&t); err != nil {
				log.Println(err)
			}
			// activity = t.ID

			// activity.RegID = activity.RegID
			// activity.UserID = activity.UserID
			// activity.UserInfo = activity.UserInfo
			a = model.Ranking{
				UserID:        t.UserID,
				ActivityInfo:  []model.ActivityInfo{},
				ToTalDistance: t.ToTalDistance,
				UserInfo:      t.UserInfo,
				EventCode:     req.EventCode,
				ParentRegID:   t.ParentRegID,
				TicketID:      req.TicketID,
				RegID:         t.ParentRegID,
			}
			a.RankNo = n + 1
			if n < 10 {
				if err == nil {
					//a.UserInfo = user
					url, _ := GetUserAvatar(t.UserID)
					teamModel := model.TeamIcon{
						RegID:   a.ParentRegID,
						IconURL: url,
					}
					teamModel, err = GetTeamIcon(teamModel)
					//
					a.ImageURL = teamModel.IconURL
					u, err := GetTeamInfo(req.EventCode, t.ParentRegID)
					if err == nil {
						a.UserInfo = u
					}
					activities = append(activities, a)
				}
			}

			temps = append(temps, a)
			n++
		}

		index := -1
		for n, s := range temps {
			if s.ParentRegID == req.ParentRegID {
				index = n
				break
			}
		}

		//index := SliceIndex(len(activities), func(i int) bool { return activities[i].UserID == objectUserID })
		if index != -1 && index > 9 {
			// if (index - 2) >= 0 {
			// 	ranking := temps[index-2]
			// 	if err != nil {
			// 		log.Println(err)
			// 	}
			// 	ranking.RankNo = index - 1
			// 	myActivities = append(myActivities, ranking)
			// }
			if (index - 1) >= 0 {
				ranking := temps[index-1]
				ranking.RankNo = index
				url, _ := GetUserAvatar(ranking.UserID)
				teamModel := model.TeamIcon{
					RegID:   ranking.ParentRegID,
					IconURL: url,
				}
				teamModel, err = GetTeamIcon(teamModel)
				ranking.ImageURL = teamModel.IconURL
				u, err := GetTeamInfo(req.EventCode, ranking.ParentRegID)
				if err == nil {
					ranking.UserInfo = u
				}
				myActivities = append(myActivities, ranking)
			}
			ranking := temps[index]
			ranking.RankNo = index + 1
			url, _ := GetUserAvatar(ranking.UserID)
			teamModel := model.TeamIcon{
				RegID:   ranking.ParentRegID,
				IconURL: url,
			}
			teamModel, err = GetTeamIcon(teamModel)
			ranking.ImageURL = teamModel.IconURL
			u, err := GetTeamInfo(req.EventCode, ranking.ParentRegID)
			if err == nil {
				ranking.UserInfo = u
			}
			myActivities = append(myActivities, ranking)
			if (index + 1) < len(temps) {
				ranking := temps[index+1]
				url, _ := GetUserAvatar(ranking.UserID)
				teamModel := model.TeamIcon{
					RegID:   ranking.ParentRegID,
					IconURL: url,
				}
				teamModel, err = GetTeamIcon(teamModel)
				//url, _ := GetUserAvatar(ranking.UserID)
				ranking.ImageURL = teamModel.IconURL
				ranking.RankNo = index + 2
				u, err := GetTeamInfo(req.EventCode, ranking.ParentRegID)
				if err == nil {
					ranking.UserInfo = u
				}
				myActivities = append(myActivities, ranking)
			}
			// if (index + 2) < len(temps) {
			// 	ranking := temps[index+2]
			// 	ranking.RankNo = index + 3
			// 	myActivities = append(myActivities, ranking)
			// }
		}

		count = int64(len(temps))
		// //activityInfo = activity.ActivityInfo

		// filter = bson.D{{"event_id", objectEventID}}
		// option := options.Find()
		// option.SetLimit(10)
		// option.SetSort(bson.D{{"total_distance", -1}})
		// cur, err := boardMongo.ConnectionDB.Collection(activityCollection).Find(context.TODO(), filter, option)

		return event, count, activities, myActivities, err
	}
}

//GetAllBoardByEvent board ranking
func GetAllBoardByEvent(req model.AllRankingRequest) (model.Event, int64, []model.Ranking, error) {
	//var activity model.Activity
	//var activityInfo []model.ActivityInfo
	var event = model.Event{}
	var activities = []model.Ranking{}

	// objectEventID, err := primitive.ObjectIDFromHex(eventID)
	// if err != nil {
	// 	log.Println(err)
	// 	return event, 0, activities, myActivities, err
	// }
	event, err := DetailEventByCode(req.EventCode)
	if err != nil {
		return event, 0, activities, err
	}
	isTeam := false
	for _, s := range event.Tickets {
		if s.ID == req.TicketID {
			if s.Category == "team" {
				isTeam = true
			}
			break
		}
	}

	var temps = []model.Ranking{}

	if !isTeam {
		//filter := bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
		filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "ticket.id", Value: req.TicketID}}

		option := options.Find()
		//option.SetLimit(10)
		count, err := db.DB.Collection(activityCollection).CountDocuments(context.TODO(), filter)
		if count == 0 {
			return event, count, activities, err
		}
		//filter = bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
		option.SetSort(bson.D{primitive.E{Key: "total_distance", Value: -1}})
		cur, err := db.DB.Collection(activityCollection).Find(context.TODO(), filter, option)

		if err != nil {
			log.Println(err)
			return event, count, activities, err
		}

		n := 0

		for cur.Next(context.TODO()) {
			// if n >= 10 {
			// 	break
			// }

			var activity model.ActivityV2
			var a model.Ranking
			// decode the document
			if err := cur.Decode(&activity); err != nil {
				log.Println(err)
			}
			regs, err := GetRegisterNumber(activity.EventCode, activity.RegID)

			if err == nil {
				if len(regs.TicketOptions) > 0 {
					a = model.Ranking{
						UserID:        activity.UserID,
						ActivityInfo:  []model.ActivityInfo{},
						ToTalDistance: activity.ToTalDistance,
						UserInfo:      activity.UserInfo,
						EventCode:     req.EventCode,
						ID:            activity.ID,
						RegID:         activity.RegID,
						BibNo:         regs.TicketOptions[0].RegisterNumber,
						Teams:         []model.UserOption{},
					}
					a.RankNo = n + 1
					//a.UserInfo = user

					url, _ := GetUserAvatar(a.UserID)
					a.ImageURL = url
					activities = append(activities, a)

					temps = append(temps, a)
				} else {
					a = model.Ranking{
						UserID:        activity.UserID,
						ActivityInfo:  []model.ActivityInfo{},
						ToTalDistance: activity.ToTalDistance,
						UserInfo:      activity.UserInfo,
						EventCode:     req.EventCode,
						ID:            activity.ID,
						RegID:         activity.RegID,
						BibNo:         "",
					}
					a.RankNo = n + 1
					//a.UserInfo = user

					url, _ := GetUserAvatar(a.UserID)
					a.ImageURL = url
					activities = append(activities, a)

					temps = append(temps, a)
				}
				n++
			}
		}

		// //activityInfo = activity.ActivityInfo

		// filter = bson.D{{"event_id", objectEventID}}
		// option := options.Find()
		// option.SetLimit(10)
		// option.SetSort(bson.D{{"total_distance", -1}})
		// cur, err := boardMongo.ConnectionDB.Collection(activityCollection).Find(context.TODO(), filter, option)

		return event, count, activities, err
	} else {
		filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "ticket.id", Value: req.TicketID}}
		//filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "activities.user_id", Value: objectUserID}}

		event, err = DetailEventByCode(req.EventCode)
		if err != nil {
			return event, 0, activities, err
		}
		//option.SetLimit(10)
		count, err := db.DB.Collection(activityCollection).CountDocuments(context.TODO(), filter)
		if count == 0 {
			return event, count, activities, err
		}
		//filter = bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
		matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"$and": []interface{}{bson.M{"event_code": req.EventCode}, bson.M{"ticket.id": req.TicketID}}}}}
		//groupStage := bson.D{primitive.E{Key: "$group", Value: bson.A{"_id", "$parent_reg_id"}}}
		groupStage := bson.D{primitive.E{Key: "$group", Value: bson.M{"_id": bson.M{"parent_reg_id": "$parent_reg_id", "event_code": "$event_code", "ticket_id": "$ticket_id"},
			"total_distance": bson.M{"$sum": "$total_distance"},
			"reg_id":         bson.M{"$first": "$reg_id"},
			"user_id":        bson.M{"$first": "$user_id"},
			"user_info":      bson.M{"$first": "$user_info"},
			"teams":          bson.M{"$push": "$user_info"},
			"parent_reg_id":  bson.M{"$first": "$parent_reg_id"},
			"ticket":         bson.M{"$first": "$ticket"},
		}}}
		sortStage := bson.D{primitive.E{Key: "$sort", Value: bson.M{"total_distance": -1}}}
		type ActivityTemp struct {
			EventCode     string             `json:"event_code" bson:"event_code"`
			Ticket        model.Tickets      `json:"ticket" bson:"ticket"`
			OrderID       string             `json:"order_id" bson:"order_id"`
			RegID         primitive.ObjectID `json:"reg_id" bson:"reg_id"`
			ParentRegID   primitive.ObjectID `json:"parent_reg_id" bson:"parent_reg_id"`
			UserID        primitive.ObjectID `json:"user_id" bson:"user_id"`
			ToTalDistance float64            `json:"total_distance" bson:"total_distance"`
			UserInfo      model.UserOption   `json:"user_info" bson:"user_info"`
			Teams         []model.UserOption `json:"teams" bson:"teams"`
		}
		cur, err := db.DB.Collection(activityCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage, sortStage})
		if err != nil {
			log.Println(err)
			return event, count, activities, err
		}

		n := 0
		for cur.Next(context.TODO()) {
			// if n >= 10 {
			// 	break
			// }
			// log.Println(cur)
			//var activity model.ActivityV2
			var t ActivityTemp
			var a model.Ranking
			// decode the document
			if err := cur.Decode(&t); err != nil {
				log.Println(err)
			}
			// activity = t.ID

			// activity.RegID = activity.RegID
			// activity.UserID = activity.UserID
			// activity.UserInfo = activity.UserInfo
			a = model.Ranking{
				UserID:        t.UserID,
				ActivityInfo:  []model.ActivityInfo{},
				ToTalDistance: t.ToTalDistance,
				UserInfo:      t.UserInfo,
				EventCode:     req.EventCode,
				ParentRegID:   t.ParentRegID,
				TicketID:      req.TicketID,
				RegID:         t.ParentRegID,
				Teams:         t.Teams,
			}
			a.RankNo = n + 1
			if err == nil {
				//a.UserInfo = user
				url, _ := GetUserAvatar(t.UserID)
				teamModel := model.TeamIcon{
					RegID:   a.ParentRegID,
					IconURL: url,
				}
				teamModel, err = GetTeamIcon(teamModel)
				//
				a.ImageURL = teamModel.IconURL
				u, err := GetTeamInfo(req.EventCode, t.ParentRegID)
				if err == nil {
					a.UserInfo = u
				}
				activities = append(activities, a)
			}

			temps = append(temps, a)
			n++
		}

		count = int64(len(temps))
		// //activityInfo = activity.ActivityInfo

		// filter = bson.D{{"event_id", objectEventID}}
		// option := options.Find()
		// option.SetLimit(10)
		// option.SetSort(bson.D{{"total_distance", -1}})
		// cur, err := boardMongo.ConnectionDB.Collection(activityCollection).Find(context.TODO(), filter, option)

		return event, count, activities, err
	}
}

func GetRankFinish(req model.AllRankingRequest) (model.Event, int64, []model.RankingFinish, error) {
	//var activity model.Activity
	//var activityInfo []model.ActivityInfo
	var event = model.Event{}
	var activities = []model.RankingFinish{}
	//var ticket model.Tickets

	// objectEventID, err := primitive.ObjectIDFromHex(eventID)
	// if err != nil {
	// 	log.Println(err)
	// 	return event, 0, activities, myActivities, err
	// }
	event, err := DetailEventByCode(req.EventCode)
	if err != nil {
		return event, 0, activities, err
	}
	isTeam := false
	var ticketDistance = 0.0
	for _, s := range event.Tickets {
		if s.ID == req.TicketID {
			dist, err := strconv.ParseFloat(s.Distance, 64)
			if err == nil {
				ticketDistance = dist
			}
			if s.Category == "team" {
				isTeam = true
			}
			break
		}
	}

	var temps = []model.RankingFinish{}

	if !isTeam {
		//filter := bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
		filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "ticket.id", Value: req.TicketID}}

		option := options.Find()
		//option.SetLimit(10)
		count, err := db.DB.Collection(activityCollection).CountDocuments(context.TODO(), filter)
		if count == 0 {
			return event, count, activities, err
		}
		//filter = bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
		option.SetSort(bson.D{primitive.E{Key: "total_distance", Value: -1}})
		cur, err := db.DB.Collection(activityCollection).Find(context.TODO(), filter, option)

		if err != nil {
			log.Println(err)
			return event, count, activities, err
		}

		//n := 0

		for cur.Next(context.TODO()) {
			// if n >= 10 {
			// 	break
			// }

			var activity model.ActivityV2
			var a model.RankingFinish
			// decode the document
			if err := cur.Decode(&activity); err != nil {
				log.Println(err)
			}
			var sumTotal = 0.0
			var dateFinished time.Time
			for _, s := range activity.ActivityInfo {
				if s.Status == config.ACTIVITY_STATUS_APPROVE {
					sumTotal += float64(s.Distance)
					if sumTotal >= ticketDistance {
						dateFinished = s.ActivityDate
						//break
					}
				}
			}
			regs, err := GetRegisterNumber(activity.EventCode, activity.RegID)
			if err == nil {
				if !dateFinished.IsZero() {
					a = model.RankingFinish{
						UserID:        activity.UserID,
						ActivityInfo:  []model.ActivityInfo{},
						ToTalDistance: sumTotal,
						UserInfo:      activity.UserInfo,
						EventCode:     req.EventCode,
						ID:            activity.ID,
						RegID:         activity.RegID,
						DateFinished:  dateFinished,
						BibNo:         regs.TicketOptions[0].RegisterNumber,
						Teams:         []model.UserOption{},
					}
					//a.RankNo = n + 1
					//a.UserInfo = user
					url, _ := GetUserAvatar(a.UserID)
					a.ImageURL = url
					activities = append(activities, a)

					temps = append(temps, a)
				}
				//n++
			}
		}

		sort.SliceStable(activities, func(i, j int) bool {
			return activities[i].DateFinished.Before(activities[j].DateFinished)
		})

		count = int64(len(activities))
		// //activityInfo = activity.ActivityInfo

		// filter = bson.D{{"event_id", objectEventID}}
		// option := options.Find()
		// option.SetLimit(10)
		// option.SetSort(bson.D{{"total_distance", -1}})
		// cur, err := boardMongo.ConnectionDB.Collection(activityCollection).Find(context.TODO(), filter, option)

		return event, count, activities, err
	} else {
		filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "ticket.id", Value: req.TicketID}}
		//filter := bson.D{primitive.E{Key: "event_code", Value: req.EventCode}, primitive.E{Key: "activities.user_id", Value: objectUserID}}

		event, err = DetailEventByCode(req.EventCode)
		if err != nil {
			return event, 0, activities, err
		}
		//option.SetLimit(10)
		count, err := db.DB.Collection(activityCollection).CountDocuments(context.TODO(), filter)
		if count == 0 {
			return event, count, activities, err
		}
		//filter = bson.D{primitive.E{Key: "event_id", Value: objectEventID}, primitive.E{Key: "activities.user_id", Value: objectUserID}}
		matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"$and": []interface{}{bson.M{"event_code": req.EventCode}, bson.M{"ticket.id": req.TicketID}}}}}
		unwindStage := bson.D{primitive.E{Key: "$unwind", Value: "$activity_info"}}
		groupStage := bson.D{primitive.E{Key: "$group", Value: bson.M{"_id": bson.M{"parent_reg_id": "$parent_reg_id", "event_code": "$event_code", "ticket_id": "$ticket_id"},
			// "total_distance": bson.M{"$sum": "$total_distance"},
			"reg_id":        bson.M{"$first": "$reg_id"},
			"user_id":       bson.M{"$first": "$user_id"},
			"user_info":     bson.M{"$first": "$user_info"},
			"teams":         bson.M{"$push": "$user_info"},
			"parent_reg_id": bson.M{"$first": "$parent_reg_id"},
			"ticket":        bson.M{"$first": "$ticket"},
			"activities":    bson.M{"$push": "$activity_info"},
		}}}
		// sortStage := bson.D{primitive.E{Key: "$sort", Value: bson.M{"date": -1}}}
		type ActivityTemp struct {
			EventCode     string               `json:"event_code" bson:"event_code"`
			Ticket        model.Tickets        `json:"ticket" bson:"ticket"`
			OrderID       string               `json:"order_id" bson:"order_id"`
			RegID         primitive.ObjectID   `json:"reg_id" bson:"reg_id"`
			ParentRegID   primitive.ObjectID   `json:"parent_reg_id" bson:"parent_reg_id"`
			UserID        primitive.ObjectID   `json:"user_id" bson:"user_id"`
			ToTalDistance float64              `json:"total_distance" bson:"total_distance"`
			UserInfo      model.UserOption     `json:"user_info" bson:"user_info"`
			Teams         []model.UserOption   `json:"teams" bson:"teams"`
			Activities    []model.ActivityInfo `json:"activities" bson:"activities"`
		}
		cur, err := db.DB.Collection(activityCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, unwindStage, groupStage})
		if err != nil {
			log.Println(err)
			return event, count, activities, err
		}

		//n := 0
		for cur.Next(context.TODO()) {
			// if n >= 10 {
			// 	break
			// }
			// log.Println(cur)
			//var activity model.ActivityV2
			var t ActivityTemp
			var a model.RankingFinish
			// decode the document
			if err := cur.Decode(&t); err != nil {
				log.Println(err)
			}

			sort.SliceStable(t.Activities, func(i, j int) bool {
				return t.Activities[i].ActivityDate.Before(t.Activities[j].ActivityDate)
			})

			var sumTotal = 0.0
			var dateFinished time.Time
			for _, s := range t.Activities {
				if s.Status == config.ACTIVITY_STATUS_APPROVE {
					sumTotal += float64(s.Distance)
					if sumTotal >= ticketDistance {
						dateFinished = s.ActivityDate
						//break
					}
				}
			}
			regs, err := GetRegisterNumber(t.EventCode, t.RegID)
			var bib = ""
			if l := len(regs.TicketOptions); l > 0 {
				bib = regs.TicketOptions[0].RegisterNumber
			}
			if err == nil {
				if !dateFinished.IsZero() {
					a = model.RankingFinish{
						UserID:        t.UserID,
						ActivityInfo:  []model.ActivityInfo{},
						ToTalDistance: sumTotal,
						UserInfo:      t.UserInfo,
						EventCode:     req.EventCode,
						RegID:         t.RegID,
						ParentRegID:   t.ParentRegID,
						DateFinished:  dateFinished,
						BibNo:         bib,
						Teams:         t.Teams,
					}
					//a.RankNo = n + 1
					//a.UserInfo = user
					url, _ := GetUserAvatar(t.UserID)
					teamModel := model.TeamIcon{
						RegID:   a.ParentRegID,
						IconURL: url,
					}
					teamModel, err = GetTeamIcon(teamModel)
					//
					a.ImageURL = teamModel.IconURL
					u, err := GetTeamInfo(req.EventCode, t.ParentRegID)
					if err == nil {
						a.UserInfo = u
					}
					activities = append(activities, a)

					temps = append(temps, a)
				}
				//n++
			}
			// activity = t.ID

			// activity.RegID = activity.RegID
			// activity.UserID = activity.UserID
			// activity.UserInfo = activity.UserInfo
			// a = model.RankingFinish{
			// 	UserID:        t.UserID,
			// 	ActivityInfo:  []model.ActivityInfo{},
			// 	ToTalDistance: t.ToTalDistance,
			// 	UserInfo:      t.UserInfo,
			// 	EventCode:     req.EventCode,
			// 	ParentRegID:   t.ParentRegID,
			// 	TicketID:      req.TicketID,
			// 	RegID:         t.ParentRegID,
			// }
			//a.RankNo = n + 1
			// if err == nil {
			// 	//a.UserInfo = user
			// 	url, _ := GetUserAvatar(t.UserID)
			// 	teamModel := model.TeamIcon{
			// 		RegID:   a.ParentRegID,
			// 		IconURL: url,
			// 	}
			// 	teamModel, err = GetTeamIcon(teamModel)
			// 	//
			// 	a.ImageURL = teamModel.IconURL
			// 	u, err := GetTeamInfo(req.EventCode, t.ParentRegID)
			// 	if err == nil {
			// 		a.UserInfo = u
			// 	}
			// 	activities = append(activities, a)
			// }

			// temps = append(temps, a)
			//n++
		}

		sort.SliceStable(activities, func(i, j int) bool {
			return activities[i].DateFinished.Before(activities[j].DateFinished)
		})

		count = int64(len(activities))
		// //activityInfo = activity.ActivityInfo

		// filter = bson.D{{"event_id", objectEventID}}
		// option := options.Find()
		// option.SetLimit(10)
		// option.SetSort(bson.D{{"total_distance", -1}})
		// cur, err := boardMongo.ConnectionDB.Collection(activityCollection).Find(context.TODO(), filter, option)

		return event, count, activities, err
	}
}

//SliceIndex get index array object
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func GetTeamInfo(eventCode string, regID primitive.ObjectID) (model.UserOption, error) {
	var r model.UserOption

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": eventCode}}}
	// projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$and": []interface{}{bson.M{"$$regs._id": req.RegID}, bson.M{"$$regs.user_id": req.UserID}}}}}, "event_code": 1, "ref2": 1, "event": 1}}}
	// if req.ParentRegID.IsZero() {
	// 	log.Println("zero")
	// 	log.Println(req.RegID)
	// 	projectStage = bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.M{"$$regs._id": req.RegID, "$$regs.user_id": req.UserID}}, "event_code": 1, "ref2": 1, "event": 1}}}}}
	// }
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs._id", regID}}}}, "event_code": 1, "ref2": 1, "event": 1}}}
	unwindStage := bson.D{primitive.E{Key: "$unwind", Value: "$regs"}}
	cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, unwindStage})
	//cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	//u.Regs = []model.Regs{}
	for cur.Next(context.TODO()) {
		var u Reg
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Print(err)
		}
		return u.Regs.TicketOptions[0].UserOption, err
	}
	return r, err
}

func GetUserInfo(req model.ActivityV2) (model.UserOption, error) {
	var r model.UserOption

	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.M{"event_code": req.EventCode}}}
	// projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$and": []interface{}{bson.M{"$$regs._id": req.RegID}, bson.M{"$$regs.user_id": req.UserID}}}}}, "event_code": 1, "ref2": 1, "event": 1}}}
	// if req.ParentRegID.IsZero() {
	// 	log.Println("zero")
	// 	log.Println(req.RegID)
	// 	projectStage = bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.M{"$$regs._id": req.RegID, "$$regs.user_id": req.UserID}}, "event_code": 1, "ref2": 1, "event": 1}}}}}
	// }
	projectStage := bson.D{bson.E{Key: "$project", Value: bson.M{"regs": bson.M{"$filter": bson.M{"input": "$regs", "as": "regs", "cond": bson.M{"$eq": bson.A{"$$regs._id", req.RegID}}}}, "event_code": 1, "ref2": 1, "event": 1}}}
	unwindStage := bson.D{primitive.E{Key: "$unwind", Value: "$regs"}}
	cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage, unwindStage})
	//cur, err := registerMongo.ConnectionDB.Collection(registerCollection).Find(context.TODO(), filter)
	//log.Printf("[info] cur %s", cur)
	if err != nil {
		log.Println(err)
	}

	//u.Regs = []model.Regs{}
	for cur.Next(context.TODO()) {
		var u Reg
		// decode the document
		if err := cur.Decode(&u); err != nil {
			log.Print(err)
		}
		return u.Regs.TicketOptions[0].UserOption, err
		// event, err := DetailEventOwnerByCode(u.EventCode)
		// if err == nil {
		// 	u.Event = event
		// }
		// r := model.RegisterV2{
		// 	UserCode:  u.UserCode,
		// 	OwnerID:   u.OwnerID,
		// 	EventCode: u.EventCode,
		// 	Event:     event,
		// 	Ref2:      u.Ref2,
		// 	Regs:      append(regs, u.Regs),
		// }
		//register = append(register, r)
	}
	// cur, err := db.DB.Collection(registerCollection).Aggregate(context.TODO(), mongo.Pipeline{matchStage, projectStage})
	// if err != nil {
	// 	return r, err
	// }

	// for cur.Next(context.TODO()) {
	// 	// decode the document
	// 	if err := cur.Decode(&r); err != nil {
	// 		log.Println(err)
	// 		return r, err
	// 	}
	// 	//fmt.Printf("post: %+v\n", p)
	// }
	// log.Println(r.Regs)
	return r, err
}

func UpdateActivityV3() {
	cur, err := db.DB.Collection(activityCollection).Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err)
	}
	for cur.Next(context.TODO()) {
		var activity model.ActivityV2
		// decode the document
		if err := cur.Decode(&activity); err != nil {
			log.Println(err)
		}
		user, err := GetUserInfo(activity)
		if err == nil {
			// if len(user.Regs) > 0 {
			// 	log.Println(user.Regs[0].TicketOptions[0].UserOption)
			filter := bson.M{"_id": activity.ID}
			update := bson.M{"$set": bson.M{"user_info": user}}
			res := db.DB.Collection(activityCollection).FindOneAndUpdate(context.TODO(), filter, update)
			log.Println(res.Err())
			// }
		}

	}
}
