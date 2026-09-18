package transport

type CreateItemTypeRequest struct {
	Name string `json:"name"`
}

type ItemTypeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
