package tradfi

type SwitchPayload struct {
	Action      string `json:"action"`
	Batter      uint16 `json:"batter"`
	LinkQuality uint16 `json:"linkquality"`
}
