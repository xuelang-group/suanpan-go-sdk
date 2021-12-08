package logkit

type Data struct {
	Node string `json:"node"`
}

type EventLog struct {
	Title string `json:"title"`
	Level string `json:"level"`
	Time  string `json:"time"`
	Data  Data   `json:"data"`
}

type Event struct {
	Name  string
	AppID string
	Log   EventLog
}
