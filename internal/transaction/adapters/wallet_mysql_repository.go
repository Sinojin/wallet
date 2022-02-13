package adapters

import (
	"errors"

	"github.com/sinojin/wallet/internal/transaction/domain/wallet"
	"gorm.io/gorm"
)

type ActionType int

const (
	None    ActionType = iota
	Deposit            // Deposit
	Credit             // Credit
)

//Against duplications it is truested mysql transactions

type MysqlModel struct {
	gorm.Model
	UserID string     `json:"user_id" gorm:"column:user_id"`
	Amount string     `json:"amount" gorm:"column:amount"`
	Action ActionType `json:"action" gorm:"column:action_type"`
}

func (MysqlModel) TableName() string {
	return "transactions"
}

type walletRepository struct {
	client *gorm.DB
}

//init func for mysql implementation
func NewMysqlWalletRepository(db *gorm.DB) *walletRepository {
	db.AutoMigrate(&MysqlModel{})
	w := walletRepository{db}
	_, err := w.GetByUserID("1234")
	if err != nil {
		db.Create(&MysqlModel{
			UserID: "1234",
			Action: Deposit,
			Amount: "100",
		})
	}

	return &w
}

// sql limit is not important right now...
//there is a limit to get but it is not consern
func (w *walletRepository) GetByUserID(userID string) (*wallet.Wallet, error) {

	tx := w.client.Begin()
	uWallet, err := w.getByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return uWallet, err
	}
	tx.Commit()
	return uWallet, err
}

func (w *walletRepository) Deposit(userID string, amount string) (*wallet.Wallet, error) {
	tx := w.client.Begin()
	userWallet, err := w.getByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = userWallet.Deposit(amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	result := tx.Create(&MysqlModel{
		UserID: userWallet.UserID,
		Action: Deposit,
		Amount: amount,
	})

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	updatedWallet, err := w.getByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return updatedWallet, nil
}

//
// Deposit(userID string, amount float64) (*Wallet, error)
func (w *walletRepository) Credit(userID string, amount string) (*wallet.Wallet, error) {
	tx := w.client.Begin()
	userWallet, err := w.getByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = userWallet.Credit(amount)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	result := tx.Create(&MysqlModel{
		UserID: userWallet.UserID,
		Action: Credit,
		Amount: amount,
	})

	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	updatedWallet, err := w.getByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return updatedWallet, nil
}

func (w *walletRepository) getByUserID(tx *gorm.DB, userID string) (*wallet.Wallet, error) {
	var transactions []MysqlModel
	result := tx.Where("user_id = ? ", userID).Find(&transactions)
	if (result.Error != nil && result.Error == gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("User not found !!")
	} else if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	userWallet := wallet.NewWallet(userID)
	for _, trns := range transactions {
		if trns.Action == Deposit {
			userWallet.Deposit(trns.Amount)
		} else if trns.Action == Credit {
			userWallet.Credit(trns.Amount)
		}
	}

	return userWallet, nil
}
