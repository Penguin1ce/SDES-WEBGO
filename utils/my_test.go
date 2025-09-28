package utils

import (
	"reflect"
	"testing"
)

// 辅助函数：比较两个 []int 是否相等，不相等则报错
func assertEqual(t *testing.T, name string, got, want []int) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s: got %v, want %v", name, got, want)
	}
}

func TestSDES_AssignmentVersion(t *testing.T) {
	// 注意：非累计左移
	// k1 = P8(LS^1(P10(K)))，k2 = P8(LS^2(P10(K)))

	// 题设
	K := StringToBits("1010000010", 10)
	P := StringToBits("10110000", 8)

	// 1) 子密钥
	k1, k2 := KeyExpansion(K)
	assertEqual(t, "k1", k1, []int{1, 0, 1, 0, 0, 1, 0, 0}) // 10100100
	assertEqual(t, "k2", k2, []int{1, 0, 0, 1, 0, 0, 1, 0}) // 10010010

	// 2) 初始置换
	ip := Permute(P, IP[:])
	assertEqual(t, "IP(P)", ip, []int{0, 0, 1, 1, 1, 0, 0, 0}) // 00111000
	L0, R0 := ip[:4], ip[4:]
	assertEqual(t, "L0", L0, []int{0, 0, 1, 1}) // 0011
	assertEqual(t, "R0", R0, []int{1, 0, 0, 0}) // 1000

	// 3) 第一轮：f_k1 对 R0
	ep := Permute(R0, EP[:])
	assertEqual(t, "EP(R0)", ep, []int{0, 1, 0, 0, 0, 0, 0, 1}) // 01000001
	x1 := XOR(ep, k1)
	assertEqual(t, "EP(R0) XOR k1", x1, []int{1, 1, 1, 0, 0, 1, 0, 1}) // 11100101
	sout1 := SBoxSubstitution(x1)
	assertEqual(t, "Sout1 (S1||S2)", sout1, []int{1, 1, 0, 1}) // 1101
	p4_1 := Permute(sout1, SPBox[:])
	assertEqual(t, "P4_1", p4_1, []int{1, 1, 0, 1}) // 1101
	L1 := XOR(L0, p4_1)
	assertEqual(t, "L1", L1, []int{1, 1, 1, 0}) // 1110
	R1 := R0
	assertEqual(t, "R1", R1, []int{1, 0, 0, 0}) // 1000

	// 4) 交换
	L2, R2 := R1, L1
	assertEqual(t, "L2", L2, []int{1, 0, 0, 0}) // 1000
	assertEqual(t, "R2", R2, []int{1, 1, 1, 0}) // 1110

	// 5) 第二轮：f_k2 对 R2
	ep2 := Permute(R2, EP[:])
	assertEqual(t, "EP(R2)", ep2, []int{0, 1, 1, 1, 1, 1, 0, 1}) // 01111101
	x2 := XOR(ep2, k2)
	assertEqual(t, "EP(R2) XOR k2", x2, []int{1, 1, 1, 0, 1, 1, 1, 1}) // 11101111
	sout2 := SBoxSubstitution(x2)
	assertEqual(t, "Sout2", sout2, []int{1, 1, 1, 1}) // 1111
	p4_2 := Permute(sout2, SPBox[:])
	assertEqual(t, "P4_2", p4_2, []int{1, 1, 1, 1}) // 1111
	L3 := XOR(L2, p4_2)
	assertEqual(t, "L3", L3, []int{0, 1, 1, 1}) // 0111
	R3 := R2
	assertEqual(t, "R3", R3, []int{1, 1, 1, 0}) // 1110

	// 6) IP^-1
	out := Permute(append(L3, R3...), IPInverse[:])
	assertEqual(t, "Ciphertext", out, []int{1, 0, 1, 1, 1, 1, 0, 1}) // 10111101
}
