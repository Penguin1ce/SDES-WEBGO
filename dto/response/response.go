package response

type EncryptResponse struct {
	Ciphertext      string `json:"ciphertext,omitempty"`
	CiphertextASCII string `json:"ciphertext_ascii,omitempty"`
	Success         bool   `json:"success"`
	Message         string `json:"message,omitempty"`
}

type DecryptResponse struct {
	Plaintext      string `json:"plaintext,omitempty"`
	PlaintextASCII string `json:"plaintext_ascii,omitempty"`
	Success        bool   `json:"success"`
	Message        string `json:"message,omitempty"`
}
