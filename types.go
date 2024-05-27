package go_byte_tts

type App struct {
	Appid   string `json:"appid"`
	Token   string `json:"token"`
	Cluster string `json:"cluster"`
}

type Rep struct {
	ReqID     string `json:"reqid"`
	Code      int    `json:"code"`
	Message   string `json:"Message"`
	Operation string `json:"operation"`
	Sequence  int    `json:"sequence"`
	Data      string `json:"data"`
}
