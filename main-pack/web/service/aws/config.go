package aws

type Config struct {
	Address string `json:"address"`
	Region  string `json:"region"`
	Profile string `json:"profile"`
	ID      string `json:"ID"`
	Secret  string `json:"secret"`
}
