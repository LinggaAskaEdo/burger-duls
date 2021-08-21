package dto

// Request struct
type Request struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      uint8  `json:"age"`
	Address  string `json:"address"`
}