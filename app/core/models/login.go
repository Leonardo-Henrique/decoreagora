package models

import "errors"

type (
	LoginRequest struct {
		Email string `json:"email"`
	}

	AutheticateCodeRequest struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	LoginResponse struct {
		Token string `json:"auth_token"`
	}
)

func (a *AutheticateCodeRequest) ValidateRequiredFields() error {
	if a.Email == "" || a.Code == "" {
		return errors.New("required fields are missing")
	}
	return nil
}

func (l *LoginRequest) ValidateRequiredFields() error {
	if l.Email == "" {
		return errors.New("required fields are missing")
	}
	return nil
}
