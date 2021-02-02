package model

// Province struct
type Province struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	AdLevel      string `json:"ad_level" bson:"AD_LEVEL"`
	TaID         string `json:"ta_id" bson:"TA_ID"`
	TambonThai   string `json:"tambon_t" bson:"TAMBON_T"`
	TambonEn     string `json:"tambon_e" bson:"TAMBON_E"`
	AmID         string `json:"am_id" bson:"AM_ID"`
	AmphoeThai   string `json:"amphoe_t" bson:"AMPHOE_T"`
	AmphoeEn     string `json:"amphoe_e" bson:"AMPHOE_E"`
	ChID         string `json:"ch_id" bson:"CH_ID"`
	ChangwatThai string `json:"changwat_t" bson:"CHANGWAT_T"`
	ChangwatEn   string `json:"changwat_e" bson:"CHANGWAT_E"`
	Lat          string `json:"lat" bson:"LAT"`
	Long         string `json:"long" bson:"LONG"`
}

// Tambon struct
type Tambon struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	District     string `json:"district" bson:"district"`
	Amphoe       string `json:"amphoe" bson:"amphoe"`
	Province     string `json:"province" bson:"province"`
	Zipcode      int64  `json:"zipcode" bson:"zipcode"`
	DistrictCode int64  `json:"district_code" bson:"district_code"`
	AmphoeCode   int64  `json:"amphoe_code" bson:"amphoe_code"`
	ProvinceCode int64  `json:"province_code" bson:"province_code"`
}
