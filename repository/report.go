package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"thinkdev.app/think/runex/runexapi/model"
)

type ReportRepository interface {
	GetDashboardByEvent(eventID string) (model.ReportDashboard, error)
}

type ReportRepositoryMongo struct {
	ConnectionDB *mongo.Database
}

func (reportMongo ReportRepositoryMongo) GetDashboardByEvent(eventID string) (model.ReportDashboard, error) {

	var dashboard model.ReportDashboard

	return dashboard, nil
}
