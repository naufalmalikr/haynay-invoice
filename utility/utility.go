package utility

import (
	"strconv"
)

func Rupiah(amount int64) string {
	amountStr := strconv.FormatInt(amount, 10)
	n := len(amountStr)
	if n <= 3 {
		return "Rp" + amountStr
	}

	rupiah := ""
	for i := n - 1; i >= 0; i-- {
		rupiah = string(amountStr[i]) + rupiah
		if (n-i)%3 == 0 && i != 0 {
			rupiah = "." + rupiah
		}
	}

	return "Rp" + rupiah
}
