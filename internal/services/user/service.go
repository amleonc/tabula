package user

import (
	"context"
	"log"
	"strings"

	ev "github.com/amleonc/evalrunes"
	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/dto"
	"github.com/amleonc/tabula/internal/services/user/internal"
	"golang.org/x/crypto/bcrypt"
)

const (
	DevRole = iota
	AdminRole
	ModRole
	AnonRole

	MinNameLength     = 2
	MaxNameLength     = 25
	MinPasswordLength = 8
	MaxPasswordLength = 25
)

type Service interface {
	Signup(context.Context, *dto.User) (*dto.User, error)
	Login(context.Context, *dto.User) (*dto.User, error)
}

func NewService() Service {
	return service
}

type CredentialsError struct {
	msg string
}

func (c CredentialsError) Error() string {
	return c.msg
}

func newCredentialsError(msg string) CredentialsError {
	return CredentialsError{msg}
}

const (
	invalidName        = "error with invalid characters or out-of-bound name length"
	invalidPassword    = "error with invalid characters or out-of-bound password length"
	invalidCredentials = "invalid credentials"
)

type serviceStruct struct {
	repo internal.Repository
}

var service = &serviceStruct{internal.Repository{}}

func (s *serviceStruct) Signup(ctx context.Context, u *dto.User) (*dto.User, error) {
	var err error
	if err = validateNewUser(u); err != nil {
		return nil, err
	}
	u.Password, err = saltPassword(u.Password)
	if err != nil {
		return nil, err
	}
	userFromDB := &dao.User{
		Name:     u.Name,
		Password: u.Password,
		Role:     AnonRole,
	}
	err = s.repo.Create(ctx, userFromDB)
	if err != nil {
		return nil, err
	}
	u.ID = userFromDB.ID
	u.Password = ""
	u.CreatedAt = userFromDB.CreatedAt
	u.UpdatedAt = userFromDB.UpdatedAt
	return u, nil
}

func (s *serviceStruct) Login(ctx context.Context, u *dto.User) (*dto.User, error) {
	var err error
	userFromDB, err := s.repo.SelectByName(ctx, u.Name)
	if err != nil {
		log.Println(err.Error())
		return nil, newCredentialsError(invalidCredentials)
	}
	if err = comparePasswords(userFromDB.Password, u.Password); err != nil {
		return nil, newCredentialsError(invalidCredentials)
	}
	u.ID = userFromDB.ID
	u.Password = ""
	u.Role = userFromDB.Role
	u.CreatedAt = userFromDB.CreatedAt
	u.UpdatedAt = userFromDB.UpdatedAt
	u.DeletedAt = nil
	return u, nil
}

func validateNewUser(u *dto.User) error {
	var ok bool
	u.Name = strings.TrimSpace(u.Name)
	if ok = validateName(u.Name); !ok {
		return newCredentialsError(invalidName)
	}
	u.Password = strings.TrimSpace(u.Password)
	if ok = validatePassword(u.Password); !ok {
		return newCredentialsError(invalidPassword)
	}
	return nil
}

func validateName(n string) bool {
	r := []rune(n)
	return ev.CompareAgainstValidators(r, ev.IsAlphanumeric) &&
		ev.ValidateRuneArrayLength(r, MinNameLength, MaxNameLength)
}

func validatePassword(p string) bool {
	r := []rune(p)
	fn := func(r rune) bool {
		var specialChars = "!#$%&/(),.-_"
		for _, c := range specialChars {
			if r == c {
				return true
			}
		}
		return false
	}
	return ev.CompareAgainstValidators(r, ev.IsAlphanumeric, fn) &&
		ev.ValidateRuneArrayLength(r, MinPasswordLength, MaxPasswordLength)
}

func saltPassword(p string) (string, error) {
	saltedBytes, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(saltedBytes), nil
}

func comparePasswords(p1, p2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(p1), []byte(p2))
}
