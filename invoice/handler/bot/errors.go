package bot

import "errors"

var (
	errInvalidChoice = errors.New("pilihan tidak valid")
	errInvalidNumber = errors.New("angka tidak valid")
	errNegative      = errors.New("tidak boleh negatif")
	errUnexpected    = errors.New("wah, coba kasih tau malik")
)
