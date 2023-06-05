package model

type Configuration struct {
	Dbname   string `json:"dbname"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	User     string `json:"user"`
}

type Users struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
