package controller

import (
	"SDES/dto/request"
	"SDES/dto/response"
	"SDES/utils"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func BlastingHandler(c *gin.Context) {
	var req request.BlastingRequest
	var startTime = time.Now()
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BlastingResponse{
			Success: false,
			Message: "无效的请求格式",
		})
		return
	}
	if req.Plaintext == "" || req.Ciphertext == "" {
		c.JSON(http.StatusBadRequest, response.BlastingResponse{
			Success: false,
			Message: "plaintext和ciphertext不能为空",
		})
		return
	}
	if !utils.IsValidBinary(req.Plaintext, 8) {
		c.JSON(http.StatusBadRequest, response.BlastingResponse{
			Success: false,
			Message: "plaintext必须是8位二进制字符串（只包含0和1）",
		})
		return
	}
	if !utils.IsValidBinary(req.Ciphertext, 8) {
		c.JSON(http.StatusBadRequest, response.BlastingResponse{
			Success: false,
			Message: "ciphertext必须是8位二进制字符串（只包含0和1）",
		})
		return
	}
	// 将输入的明文和密文转换为位数组
	plaintextBits := utils.StringToBits(req.Plaintext, 8)
	ciphertextBits := utils.StringToBits(req.Ciphertext, 8)
	log.Printf("plaintextBits: %s", utils.BitsToString(plaintextBits))
	log.Printf("ciphertextBits: %s", utils.BitsToString(ciphertextBits))
	log.Println("开始暴力破解...")

	// 用于收集所有匹配的密钥
	var (
		foundKeys        []string
		foundKeysDecimal []int
	)
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	// 用两个线程进行破解
	ranges := [][2]int{{0, 512}, {512, 1024}}

	for _, r := range ranges {
		start, end := r[0], r[1]
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			localKeys := make([]string, 0, 4)
			localKeysDecimal := make([]int, 0, 4)

			for i := start; i < end; i++ {
				keyBits := utils.IntTo10BitKey(i)
				encryptedBits := utils.Encrypt(plaintextBits, keyBits)

				match := true
				for j := 0; j < 8; j++ {
					if encryptedBits[j] != ciphertextBits[j] {
						match = false
						break
					}
				}

				if match {
					keyString := utils.BitsToString(keyBits)
					localKeys = append(localKeys, keyString)
					localKeysDecimal = append(localKeysDecimal, i)
					log.Printf("找到匹配密钥：%s（十进制：%d）", keyString, i)
				}
			}

			mu.Lock()
			foundKeys = append(foundKeys, localKeys...)
			foundKeysDecimal = append(foundKeysDecimal, localKeysDecimal...)
			mu.Unlock()
		}(start, end)
	}

	wg.Wait()
	var endTime = time.Now()
	var duration = endTime.Sub(startTime)
	var timeString = fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1000000)
	// 根据找到的密钥数量返回相应结果
	if len(foundKeys) > 0 {
		var message string
		if len(foundKeys) == 1 {
			message = fmt.Sprintf("成功破解！找到1个密钥：%s（十进制：%d）", foundKeys[0], foundKeysDecimal[0])
		} else {
			message = fmt.Sprintf("成功破解！找到%d个可能的密钥", len(foundKeys))
		}
		log.Printf("暴力破解完成！总共找到%d个匹配的密钥", len(foundKeys))

		c.JSON(http.StatusOK, response.BlastingResponse{
			Success:     true,
			Message:     message,
			Plaintext:   req.Plaintext,
			Ciphertext:  req.Ciphertext,
			Keys:        foundKeys,
			KeysDecimal: foundKeysDecimal,
			KeyCount:    len(foundKeys),
			Time:        timeString,
		})
	} else {
		log.Println("暴力破解完成，未找到匹配的密钥")
		c.JSON(http.StatusOK, response.BlastingResponse{
			Success: false,
			Message: "暴力破解失败：未找到匹配的密钥",
			Time:    timeString,
		})
	}
}
