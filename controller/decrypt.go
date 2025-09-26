package controller

import (
	"SDES/dto/request"
	"SDES/dto/response"
	"SDES/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DecryptHandler(c *gin.Context) {
	var req request.DecryptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.DecryptResponse{
			Success: false,
			Message: "无效的请求格式",
		})
		return
	}

	// 验证输入
	if !utils.IsValidBinary(req.Ciphertext, 8) {
		c.JSON(http.StatusBadRequest, response.DecryptResponse{
			Success: false,
			Message: "密文必须是8位二进制字符串（只包含0和1）",
		})
		return
	}

	if !utils.IsValidBinary(req.Key, 10) {
		c.JSON(http.StatusBadRequest, response.DecryptResponse{
			Success: false,
			Message: "密钥必须是10位二进制字符串（只包含0和1）",
		})
		return
	}

	// 转换为位数组
	ciphertextBits := utils.StringToBits(req.Ciphertext, 8)
	keyBits := utils.StringToBits(req.Key, 10)

	// 解密
	plaintextBits := utils.Decrypt(ciphertextBits, keyBits)
	plaintext := utils.BitsToString(plaintextBits)

	c.JSON(http.StatusOK, response.DecryptResponse{
		Plaintext: plaintext,
		Success:   true,
	})
}
