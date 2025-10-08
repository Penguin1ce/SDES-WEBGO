package request

// DecryptRequest API 请求结构体
type DecryptRequest struct {
	Ciphertext       string  `json:"ciphertext"`
	CiphertextBase64 *string `json:"ciphertext_base64"`
	Key              string  `json:"key" binding:"required"`
}

// EncryptRequest API 请求结构体
type EncryptRequest struct {
	Plaintext      string  `json:"plaintext"`
	PlaintextASCII *string `json:"plaintext_ascii"`
	Key            string  `json:"key" binding:"required"`
}

type BlastingRequest struct {
	Plaintext  string `json:"plaintext"`
	Ciphertext string `json:"ciphertext"`
}
