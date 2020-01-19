package grpc

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/grpclog"
	"log"
	"projectionist/consts"
	"projectionist/models"
	projProto "projectionist/proto"
)

func (p *ProjectionistServer) Login(ctx context.Context, r *projProto.LoginRequest) (*projProto.LoginResponse, error) {
	var respond = &projProto.LoginResponse{Meta: &projProto.DefaultResponse{}}
	err := r.Validate()
	if err != nil {
		grpclog.Errorf("LoginRequest.Validate error: %v", err)
		return respond, errors.New(consts.InputDataInvalidResp)
	}

	var iUser models.Model
	var user = &models.User{}

	iUser, err = p.dbProvider.GetByName(user, r.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			grpclog.Error("dbProvider.GetByName with username %s error: %v", r.Username, err)
			return respond, errors.New(consts.NotExistResp)
		}
		grpclog.Errorf("dbProvider.GetByName error: %v", err)
		respond.Meta.Message = consts.SmtWhenWrongResp
		respond.Meta.Status = false
		return respond, errors.New(consts.SmtWhenWrongResp)
	}

	user, ok := iUser.(*models.User)
	if !ok {
		grpclog.Errorf("LoginApi() error: iUser is not User model")
		respond.Meta.Message = consts.SmtWhenWrongResp
		respond.Meta.Status = false
		return respond, errors.New(consts.SmtWhenWrongResp)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)); err != nil {
		log.Printf("LoginApi() bcrypt.CompareHashAndPassword() error: %v", err)
		return respond, errors.New("Not authorized")
	}

	var tokenM = &models.Token{UserId: uint64(user.ID)}
	var token = jwt.NewWithClaims(jwt.GetSigningMethod(jwt.SigningMethodHS256.Name), tokenM)
	tokenStr, err := token.SignedString([]byte(p.cfg.TokenSecretKey))
	if err != nil {
		grpclog.Errorf("LoginApi() token.SignedString() error: %v", err)
		return respond, errors.New("Authorization failed")
	}

	grpclog.Infof("Login successful for %v", user.Username)

	respond.Meta.Message = "Login successful"
	respond.User = &projProto.User{
		Id:       int64(user.ID),
		Username: user.Username,
		Password: "",
		Role:     projProto.UserRole(user.Role),
		Token:    tokenStr,
		Deleted:  projProto.Deleted(user.Deleted),
	}
	respond.Meta.Status = true
	return respond, nil
}
