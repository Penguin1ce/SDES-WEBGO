package controller

import (
	"SDES/dto/request"
	"SDES/dto/response"
	"SDES/utils"
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

	// 验证输入
	if !utils.IsValidBinary(req.Plaintext, 8) {
		c.JSON(http.StatusBadRequest, response.EncryptResponse{
			Success: false,
			Message: "明文必须是8位二进制字符串（只包含0和1）",
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

	// 转换为位数组
	plaintextBits := utils.StringToBits(req.Plaintext, 8)
	keyBits := utils.StringToBits(req.Key, 10)

	// 加密
	ciphertextBits := utils.Encrypt(plaintextBits, keyBits)
	ciphertext := utils.BitsToString(ciphertextBits)

	c.JSON(http.StatusOK, response.EncryptResponse{
		Ciphertext: ciphertext,
		Success:    true,
	})
}
