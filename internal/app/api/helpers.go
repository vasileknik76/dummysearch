package api

func errorResponseWithText(msg string, code int) response {
	return response{
		responseData{
			Status: false,
			Error: struct {
				Message string `json:"message"`
			}{msg},
		},
		code,
	}
}

func successResponse(data responseData) response {
	return response{
		data,
		200,
	}
}
