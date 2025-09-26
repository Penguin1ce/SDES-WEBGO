package request

// EncryptRequest API 请求和响应结构体
type EncryptRequest struct {
	Plaintext string `json:"plaintext" binding:"required"`
	Key       string `json:"key" binding:"required"`
}
