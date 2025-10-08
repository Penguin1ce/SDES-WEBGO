package controller

import (
	"SDES/dto/request"
	"SDES/dto/response"
	"SDES/utils"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EncryptHandler API 处理函数
func EncryptHandler(c *gin.Context) {
	var req request.EncryptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.EncryptResponse{
			Success: false,
			Message: "无效的请求格式",
		})
		return
	}

	if req.Plaintext == "" && req.PlaintextASCII == nil {
		c.JSON(http.StatusBadRequest, response.EncryptResponse{
			Success: false,
			Message: "必须提供二进制明文或 ASCII 明文",
		})
		return
	}

	if !utils.IsValidBinary(req.Key, 10) {
		c.JSON(http.StatusBadRequest, response.EncryptResponse{
			Success: false,
			Message: "密钥必须是10位二进制字符串（只包含0和1）",
		})
		return
	}

	keyBits := utils.StringToBits(req.Key, 10)

	if req.PlaintextASCII != nil {
		plaintextBytes, err := utils.ASCIIStringToBytes(*req.PlaintextASCII)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.EncryptResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		if len(plaintextBytes) == 0 {
			c.JSON(http.StatusBadRequest, response.EncryptResponse{
				Success: false,
				Message: "ASCII 明文不能为空",
			})
			return
		}

		ciphertextBytes := utils.EncryptBytes(plaintextBytes, keyBits)
		ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertextBytes)

		c.JSON(http.StatusOK, response.EncryptResponse{
			CiphertextBase64: ciphertextBase64,
			Success:          true,
		})
		return
	}

	// 验证二进制输入
	if !utils.IsValidBinary(req.Plaintext, 8) {
		c.JSON(http.StatusBadRequest, response.EncryptResponse{
			Success: false,
			Message: "明文必须是8位二进制字符串（只包含0和1）",
		})
		return
	}

	// 转换为位数组
	plaintextBits := utils.StringToBits(req.Plaintext, 8)

	// 加密
	ciphertextBits := utils.Encrypt(plaintextBits, keyBits)
	ciphertext := utils.BitsToString(ciphertextBits)

	c.JSON(http.StatusOK, response.EncryptResponse{
		CiphertextBinary: ciphertext,
		Success:          true,
	})
}
