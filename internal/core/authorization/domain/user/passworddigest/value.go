package passworddigest

import "errors"

type Crypt interface {
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	Cost(hashedPassword []byte) (int, error)
}

const (
	hashCost = 10
	minLen   = 8
	maxLen   = 64
)

var (
	ErrPasswordEmpty       = errors.New("password is empty")
	ErrPasswordTooShort    = errors.New("password is too short")
	ErrPasswordTooLong     = errors.New("password is too long")
	ErrPasswordDoesntMatch = errors.New("password and password confirmation don't match")
	ErrHashInvalid         = errors.New("hash is invalid")
	ErrHashCostInvalid     = errors.New("hash cost is invalid")

	ErrInvalidPassword = errors.New("invalid password")
)

type Value struct {
	hash                 []byte
	password             string
	passwordConfirmation string
	crypt                Crypt
}

func New(hash string, crypt Crypt) (Value, error) {
	var errs error

	cost, err := crypt.Cost([]byte(hash))
	if err != nil {
		errs = errors.Join(errs, ErrHashInvalid, err)
	}

	if cost != hashCost {
		errs = errors.Join(errs, ErrHashCostInvalid)
	}

	return Value{
		hash:                 []byte(hash),
		password:             "",
		passwordConfirmation: "",
		crypt:                crypt,
	}, errs
}

func NewFromPassword(password string, passwordConfirmation string, crypt Crypt) (Value, error) {
	var errs error

	if password == "" {
		errs = errors.Join(errs, ErrPasswordEmpty)
	}

	if len(password) < minLen {
		errs = errors.Join(errs, ErrPasswordTooShort)
	}

	if len(password) > maxLen {
		errs = errors.Join(errs, ErrPasswordTooLong)
	}

	if password != passwordConfirmation {
		errs = errors.Join(errs, ErrPasswordDoesntMatch)
	}

	hash, err := crypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		errs = errors.Join(errs, ErrHashInvalid, err)
	}

	return Value{
		hash:                 hash,
		password:             password,
		passwordConfirmation: passwordConfirmation,
		crypt:                crypt,
	}, errs
}

func (v Value) String() string {
	return string(v.hash)
}

func (v Value) Compare(password string) error {
	err := v.crypt.CompareHashAndPassword(v.hash, []byte(password))
	if err != nil {
		return errors.Join(err, ErrInvalidPassword)
	}

	return nil
}
