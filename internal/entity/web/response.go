package web

import "time"

type GetUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Avatar   string `json:"avatar"`
}

type ResponseLogin struct {
	KeyID  string `json:"key_id"`
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
}

type GetDetailMail struct {
	Message  string `json:"msg"`
	From     string `json:"from"`
	To       string `json:"to"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type GetMail struct {
	ID        uint      `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	CreatedAt time.Time `json:"created_at"`
}

type GetMailLabel struct {
	Category string `json:"category"`
	Message  string `json:"msg"`
	From     string `json:"from"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	To       string `json:"to"`
}

type Category struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PGPResponse struct {
	Private  string `json:"private_key"`
	Username string `json:"username"`
	HexKeyID string `json:"hex_key"`
}
