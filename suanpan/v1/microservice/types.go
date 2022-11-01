package microservice

type LookupResp struct {
	Success bool       `json:"success"`
	Data    LookupData `json:"data"`
}

type LookupData struct {
	Address string `json:"address"`
}
