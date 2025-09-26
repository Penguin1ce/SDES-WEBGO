package request

// DecryptRequest API 请求结构体
type DecryptRequest struct {
	Ciphertext string `json:"ciphertext" binding:"required"`
	Key        string `json:"key" binding:"required"`
}

// EncryptRequest API 请求结构体
type EncryptRequest struct {
	Plaintext string `json:"plaintext" binding:"required"`
	Key       string `json:"key" binding:"required"`
}
