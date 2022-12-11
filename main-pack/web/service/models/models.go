package models

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserInfo struct {
	Credentials Credentials `json:"credentials"`
	HomeDir     string      `json:"homeDir"`
}

type Note struct {
	Name string `json:"name"`
	Text string `json:"text"`
}
