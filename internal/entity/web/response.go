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
	Key      []byte `json:"key"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
