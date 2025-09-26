package utils

// S-DES 算法实现
// 分组长度：8-bit
// 密钥长度：10-bit

// 置换表定义
var (
	// IP 初始置换 IP
	IP = [8]int{2, 6, 3, 1, 4, 8, 5, 7}

	// IPInverse 初始置换逆 IP^-1
	IPInverse = [8]int{4, 1, 3, 5, 7, 2, 8, 6}

	// P10 10位密钥置换 P10
	P10 = [10]int{3, 5, 2, 7, 4, 10, 1, 9, 8, 6}

	// P8 8位密钥置换 P8
	P8 = [8]int{6, 3, 7, 4, 8, 5, 10, 9}

	// EP 扩展置换 EP
	EP = [8]int{4, 1, 2, 3, 2, 3, 4, 1}

	// P4 P4置换
	P4 = [4]int{2, 4, 3, 1}

	// S0 S盒
	S0 = [4][4]int{
		{1, 0, 3, 2},
		{3, 2, 1, 0},
		{0, 2, 1, 3},
		{3, 1, 3, 2},
	}

	S1 = [4][4]int{
		{0, 1, 2, 3},
		{2, 0, 1, 3},
		{3, 0, 1, 0},
		{2, 1, 0, 3},
	}
)

// Permute 置换函数
func Permute(input []int, table []int) []int {
	result := make([]int, len(table))
	for i, pos := range table {
		result[i] = input[pos-1] // 表中索引从1开始，数组索引从0开始
	}
	return result
}

// LeftShift 左循环移位
func LeftShift(bits []int, positions int) []int {
	n := len(bits)
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = bits[(i+positions)%n]
	}
	return result
}

// XOR 异或运算
func XOR(a, b []int) []int {
	result := make([]int, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}

// KeyExpansion 密钥扩展算法
func KeyExpansion(key []int) ([]int, []int) {
	// P10置换
	p10Result := Permute(key, P10[:])

	// 分成两个5位的部分
	left5 := p10Result[0:5]
	right5 := p10Result[5:10]

	// 生成k1：左移1位
	left51 := LeftShift(left5, 1)
	right51 := LeftShift(right5, 1)
	combined1 := append(left51, right51...)
	k1 := Permute(combined1, P8[:])

	// 生成k2：再左移1位（总共左移2位）
	left52 := LeftShift(left51, 1)
	right52 := LeftShift(right51, 1)
	combined2 := append(left52, right52...)
	k2 := Permute(combined2, P8[:])

	return k1, k2
}

// SBoxSubstitution S盒替换
func SBoxSubstitution(input []int) []int {
	// 分成两个4位部分
	left4 := input[0:4]
	right4 := input[4:8]

	// S0盒替换
	row0 := left4[0]*2 + left4[3]
	col0 := left4[1]*2 + left4[2]
	s0Output := S0[row0][col0]

	// S1盒替换
	row1 := right4[0]*2 + right4[3]
	col1 := right4[1]*2 + right4[2]
	s1Output := S1[row1][col1]

	// 将结果转换为2位二进制
	result := make([]int, 4)
	result[0] = s0Output / 2
	result[1] = s0Output % 2
	result[2] = s1Output / 2
	result[3] = s1Output % 2

	return result
}

// FFunction f函数
func FFunction(right4 []int, subkey []int) []int {
	// 扩展置换
	expanded := Permute(right4, EP[:])

	// 与子密钥异或
	xorResult := XOR(expanded, subkey)

	// S盒替换
	sBoxResult := SBoxSubstitution(xorResult)

	// P4置换
	p4Result := Permute(sBoxResult, P4[:])

	return p4Result
}

// Swap 交换函数 SW
func Swap(input []int) []int {
	result := make([]int, 8)
	copy(result[0:4], input[4:8]) // 右半部分移到左边
	copy(result[4:8], input[0:4]) // 左半部分移到右边
	return result
}

// Encrypt 加密算法
func Encrypt(plaintext []int, key []int) []int {
	// 密钥扩展
	k1, k2 := KeyExpansion(key)

	// 初始置换 IP
	ipResult := Permute(plaintext, IP[:])

	// 第一轮
	left4 := ipResult[0:4]
	right4 := ipResult[4:8]
	fResult1 := FFunction(right4, k1)
	newLeft := XOR(left4, fResult1)
	round1Result := append(newLeft, right4...)

	// 交换 SW
	swapped := Swap(round1Result)

	// 第二轮
	left4 = swapped[0:4]
	right4 = swapped[4:8]
	fResult2 := FFunction(right4, k2)
	newLeft = XOR(left4, fResult2)
	round2Result := append(newLeft, right4...)

	// 逆初始置换 IP^-1
	ciphertext := Permute(round2Result, IPInverse[:])

	return ciphertext
}

// Decrypt 解密算法
func Decrypt(ciphertext []int, key []int) []int {
	// 密钥扩展（注意解密时密钥顺序相反）
	k1, k2 := KeyExpansion(key)

	// 初始置换 IP
	ipResult := Permute(ciphertext, IP[:])

	// 第一轮（使用k2）
	left4 := ipResult[0:4]
	right4 := ipResult[4:8]
	fResult1 := FFunction(right4, k2)
	newLeft := XOR(left4, fResult1)
	round1Result := append(newLeft, right4...)

	// 交换 SW
	swapped := Swap(round1Result)

	// 第二轮（使用k1）
	left4 = swapped[0:4]
	right4 = swapped[4:8]
	fResult2 := FFunction(right4, k1)
	newLeft = XOR(left4, fResult2)
	round2Result := append(newLeft, right4...)

	// 逆初始置换 IP^-1
	plaintext := Permute(round2Result, IPInverse[:])

	return plaintext
}

// StringToBits 辅助函数：将字符串转换为位数组
func StringToBits(s string, length int) []int {
	bits := make([]int, length)
	for i, char := range s {
		if i >= length {
			break
		}
		if char == '1' {
			bits[i] = 1
		} else {
			bits[i] = 0
		}
	}
	return bits
}

// BitsToString 辅助函数：将位数组转换为字符串
func BitsToString(bits []int) string {
	result := ""
	for _, bit := range bits {
		if bit == 1 {
			result += "1"
		} else {
			result += "0"
		}
	}
	return result
}

// IsValidBinary 验证输入是否为有效的二进制字符串
func IsValidBinary(s string, expectedLength int) bool {
	if len(s) != expectedLength {
		return false
	}
	for _, char := range s {
		if char != '0' && char != '1' {
			return false
		}
	}
	return true
}
