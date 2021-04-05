package model

// ResponseResult => For response structure
type ResponseResult struct {
	Error  string      `json:"error"`
	Result interface{} `json:"result"`
}
