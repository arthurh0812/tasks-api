package api

type Response struct {
	Status int32 `json:"status"`
	Message string `json:"message"`
	Count int64 `json:"count,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}
