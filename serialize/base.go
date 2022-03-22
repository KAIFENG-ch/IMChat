package serialize

type Base struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

type Datalist struct {
	Item  interface{} `json:"item"`
	Total int         `json:"total"`
}
