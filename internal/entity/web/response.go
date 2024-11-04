package web

type GetUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

type GetDetailMail struct {
	Message  string `json:"msg"`
	From     string `json:"from"`
	To       string `json:"to"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type GetMail struct {
	Category string `json:"category"`
	Message  string `json:"msg"`
	From     string `json:"from"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}
