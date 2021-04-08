package util

// 脱敏银行卡号
func MaskBankcard(bankcard string) string {
	if n := len(bankcard); n >= 7 {
		return bankcard[:n-7] + "****" + bankcard[n-4:]
	}
	return bankcard
}
