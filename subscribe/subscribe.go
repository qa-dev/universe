package subscribe

type Subscribe struct {
	EventName   string `json:"event_name"`
	WebHookPath string `json:"webhook"`
}