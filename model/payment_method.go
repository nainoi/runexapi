package model

//PaymentMethod model
type PaymentMethod struct {
	ID            string `json:"id" bson:"_id"`
	Name          string `json:"name" bson:"name"`
	PaymentType   string `json:"type" bson:"type"`
	Charge        int64  `json:"charge" bson:"charge"`
	ChargePercent int64  `json:"charge_percent" bson:"charge_percent"`
	IsActive      bool   `json:"is_active" bson:"is_active"`
	Icon          string `json:"icon" bson:"icon"`
}
