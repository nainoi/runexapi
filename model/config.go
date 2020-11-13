package model

// ConfigModel struct
type ConfigModel struct {
	LeaderBoardURL string `json:"leader_board_url" bson:"leader_board_url"`
	AuthenURL      string `json:"authen_url" bson:"authen_url"`
	AuthenToken    string `json:"authen_token" bson:"authen_token"`
	PreviewURL     string `json:"preview_url" bson:"preview_url"`
}
