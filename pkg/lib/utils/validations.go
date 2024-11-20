package utils

import (
	"errors"
	emailverifier "github.com/AfterShip/email-verifier"
	"log/slog"
)

var (
	verifier = emailverifier.NewVerifier()
)

func VerifyEmail(email string) error {
	ret, err := verifier.Verify(email)
	if err != nil {
		slog.Error("verify email address failed, error is: ", err.Error())
		return err
	}
	if !ret.Syntax.Valid {
		return errors.New("невалтдный email")
	}

	if ret.SMTP == nil {
		return errors.New("Используйте действительную почту")
	}
	return nil
}
