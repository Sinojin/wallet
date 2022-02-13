package wallet

import (
	"errors"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	UserID string `json:"user_id"`

	Balance decimal.Decimal
}

func NewWallet(UserID string) *Wallet {
	emptyBalance, _ := decimal.NewFromString("0")

	return &Wallet{
		UserID:  UserID,
		Balance: emptyBalance,
	}
}

func (b *Wallet) Deposit(amount string) error {
	//Amount must be positve
	ADecimal, err := b.GetAmount(amount)
	if err != nil {
		return err
	}
	newBalance := b.Balance.Add(ADecimal)

	b.Balance = newBalance
	return nil
}

func (b *Wallet) Credit(amount string) error {
	//Amount must be positve
	ADecimal, err := b.GetAmount(amount)
	if err != nil {
		return err
	}
	newBalance := b.Balance.Sub(ADecimal)
	//Decide whether it has enough money or not
	if newBalance.IsNegative() {
		return errors.New("Not enough balance to apply process")
	}
	//update the balance
	b.Balance = newBalance
	return nil
}

func (b *Wallet) GetAmount(amount string) (decimal.Decimal, error) {
	AmountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Decimal{}, errors.New("invalid amount!")
	}

	if AmountDecimal.IsNegative() {
		return decimal.Decimal{}, errors.New("Amount can not be negative")
	}
	return AmountDecimal, nil
}
