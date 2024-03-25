package utills

import (
	"fmt"

	promptpayqr "github.com/kazekim/promptpay-qr-go"
)

func GeneratePromptpayQr(phoneNumber string, amount float64) ([]byte, error) {
	qr, err := promptpayqr.QRForTargetWithAmount(phoneNumber, fmt.Sprintf("%v", amount))
	if err != nil {
		return nil, err
	}
	return *qr, nil
}
