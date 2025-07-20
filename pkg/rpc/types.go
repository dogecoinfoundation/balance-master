package rpc

type PostTrackersRequest struct {
	Address string `json:"address"`
}

type PostTrackersResponse struct {
	ID int `json:"id"`
}
