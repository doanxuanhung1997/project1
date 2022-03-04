package model

type Message struct {
	Title   string `json:"notification"`
	Content string `json:"content"`
	Icon    string `json:"icon"`
	To      string `json:"to"`
}

type InputToken struct {
	Token string `json:"token"`
}
