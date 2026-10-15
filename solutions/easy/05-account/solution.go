package lab

import (
	"errors"
	"math"
)

type Account struct{ cents int64 }

func NewAccount(cents int64) (*Account, error) {
	if cents < 0 {
		return nil, errors.New("negative balance")
	}
	return &Account{cents}, nil
}
func (a *Account) Deposit(cents int64) error {
	if cents <= 0 || cents > math.MaxInt64-a.cents {
		return errors.New("invalid deposit")
	}
	a.cents += cents
	return nil
}
func (a *Account) Withdraw(cents int64) error {
	if cents <= 0 || cents > a.cents {
		return errors.New("invalid withdrawal")
	}
	a.cents -= cents
	return nil
}
func (a *Account) Balance() int64 { return a.cents }
