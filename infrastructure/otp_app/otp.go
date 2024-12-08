package otpapp

import (
	"time"

	"github.com/edutav/licentia-usoris/internal/utils"
	"github.com/pquerna/otp/totp"
)

func GenerateOTP() (string, string, error) {
	// TODO: Criar uma chave e adicionar nas variaveis de ambiente
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Auctoritas",
		AccountName: "auctoritas@localhost.com",
		Digits:      6,
	})
	if err != nil {
		return "", "", err
	}

	otpCode, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", "", utils.ErrGenerateOTP
	}
	return key.Secret(), otpCode, nil
}

func ValidateOTP(otpCode string, secret string) bool {
	return totp.Validate(otpCode, secret)
}
