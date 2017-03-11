package web

type SubscribeData struct {
	EventName string `json:"event_name"`
	Url       string `json:"url"`
}

type UnsubscribeData struct {
	EventName string `json:"event_name"`
	Url       string `json:"url"`
}
