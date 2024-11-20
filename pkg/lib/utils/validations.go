package utils

import (
	"errors"
	emailverifier "github.com/AfterShip/email-verifier"
	"log/slog"
)

var (
	verifier = emailverifier.NewVerifier().EnableSMTPCheck()
)

func VerifyEmail(email string) error {
	ret, err := verifier.Verify(email)
	if err != nil {
		slog.Error("verify email address failed, error is: ", err.Error())
		return err
	}
	if !ret.Syntax.Valid {
		return errors.New("невалидный email")
	}
	if ret.Disposable {
		return errors.New("используйте действительную почту")
	}
	return nil
}
