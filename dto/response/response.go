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

type BlastingResponse struct {
	Plaintext   string   `json:"plaintext,omitempty"`
	Ciphertext  string   `json:"ciphertext,omitempty"`
	Keys        []string `json:"keys,omitempty"`
	KeysDecimal []int    `json:"keys_decimal,omitempty"`
	KeyCount    int      `json:"key_count,omitempty"`
	Success     bool     `json:"success"`
	Message     string   `json:"message,omitempty"`
	Time        string   `json:"time,omitempty"`
}
