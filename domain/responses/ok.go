package responses

type OkResponse struct {
	Ok bool `json:"ok"`
}

func NewOkResponse(ok bool) OkResponse {
	return OkResponse{Ok: ok}
}
