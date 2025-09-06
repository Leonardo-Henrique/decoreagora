package usecases

import (
	"errors"
	"log"
	"time"

	"github.com/Leonardo-Henrique/decoreagora/app/core/models"
	"github.com/Leonardo-Henrique/decoreagora/app/core/ports"
)

type LoginUsecase struct {
	db    ports.Database
	token ports.TokenHandler
}

func NewLoginUsecase(db ports.Database, token ports.TokenHandler) *LoginUsecase {
	return &LoginUsecase{
		db:    db,
		token: token,
	}
}

func (uc *LoginUsecase) Login(login models.LoginRequest) (string, error) {
	if err := login.ValidateRequiredFields(); err != nil {
		return "", err
	}

	exist, err := uc.db.CheckIfEmailIsRegistered(login.Email)
	if err != nil {
		return "", err
	}

	if !exist {
		return "", errors.New("no user found")
	}

	user, err := uc.db.GetUserByEmail(login.Email)
	if err != nil {
		return "", err
	}

	code := models.AccessCode{}
	code.Generate(user.ID)

	if err := uc.db.CreateAccessCode(code); err != nil {
		return "", err
	}

	return code.Code, nil
}

func (uc *LoginUsecase) AuthenticateCode(req models.AutheticateCodeRequest) (*models.LoginResponse, error) {
	if err := req.ValidateRequiredFields(); err != nil {
		return nil, errors.New("missing required fields")
	}

	user, err := uc.db.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	storedCodes, err := uc.db.GetAccessCodeByUserID(user.ID)
	if err != nil {
		log.Println(err)
		return nil, errors.New("authentication failed")
	}

	var foundValidCode bool

	for _, storedCode := range storedCodes {
		if storedCode.Code == req.Code && !*storedCode.IsUsed && time.Now().Before(storedCode.ExpireAt) {
			foundValidCode = true
			break
		}
	}

	if !foundValidCode {
		return nil, errors.New("no valid code code")
	}

	/* if storedCode.Code != req.Code {
		return nil, errors.New("invalid code")
	}

	if time.Now().After(storedCode.ExpireAt) {
		return nil, errors.New("code expired")
	}

	if storedCode.IsUsed != nil && *storedCode.IsUsed {
		return nil, errors.New("code already used")
	} */

	token, err := uc.token.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("could not generate token")
	}

	if err := uc.db.DeleteAccessCode(user.ID, req.Code); err != nil {
		log.Println("error when deleting code", err)
	}

	return &models.LoginResponse{Token: token}, nil
}
