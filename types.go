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

type TtsAsyncRep struct {
	Reqid      string `json:"reqid"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	TaskId     string `json:"task_id"`
	TaskStatus int    `json:"task_status"`
	TextLength int    `json:"text_length"`
}

type TtsAsyncQueryRep struct {
	Reqid         string `json:"reqid"`
	Code          int    `json:"code"`
	Message       string `json:"message"`
	AudioUrl      string `json:"audio_url"`
	TaskId        string `json:"task_id"`
	TaskStatus    int    `json:"task_status"`
	TextLength    int    `json:"text_length"`
	UrlExpireTime int    `json:"url_expire_time"`
}

type ChanJoinVoice struct {
	Index int
	Audio []byte
	Err   error
}
