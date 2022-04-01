package phoneswitches

type TwilioPayload struct {
	Payload TilioData `json:"payload"`
	Action  string    `json:"type"`
}

type TilioData struct {
	CallID string `json:"callsid"`
}
