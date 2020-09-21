package api

// ResponseData represent base response structure. All request returns data like this structure.
type responseData struct {
	Status  bool        `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Response it's contains both payload and http entities like status code
type response struct {
	ResponseData responseData
	StatusCode   int
}
