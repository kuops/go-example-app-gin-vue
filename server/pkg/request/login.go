package request

type Login struct {
	Username  string `json:"username" example:"admin"`
	Password  string `json:"password" example:"123456"`
}