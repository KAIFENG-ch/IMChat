package serialize

type Login struct {
	Reply string `json:"reply"`
	Token string `json:"token"`
}

type Update struct {
	Reply string `json:"reply"`
	Url   string `json:"url"`
}
