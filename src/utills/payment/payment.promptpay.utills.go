package payment

import (
	"strconv"

	promptpayqr "github.com/kazekim/promptpay-qr-go"
)

func GeneratePromptpayQr(phoneNumber string, amount int) (*[]byte, error) {
	qr, err := promptpayqr.QRForTargetWithAmount(phoneNumber, strconv.Itoa(amount))
	if err != nil {
		return nil, err
	}
	return qr, nil
}
