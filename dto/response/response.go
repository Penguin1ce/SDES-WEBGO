package response

type EncryptResponse struct {
	Ciphertext string `json:"ciphertext"`
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
}

type DecryptResponse struct {
	Plaintext string `json:"plaintext"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
}
