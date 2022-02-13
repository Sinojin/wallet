package services

import "github.com/sinojin/wallet/internal/transaction/domain/wallet"

//This interface for business rules
//we can implement this for metrics and logs too
type WalletService interface {
	GetBalance(userID string) (*wallet.Wallet, error)
	CreditWallet(userID string, amount string) (*wallet.Wallet, error)
	DepositWallet(userID string, amount string) (*wallet.Wallet, error)
}

type walletService struct {
	WalletRepository wallet.WalletRepository
}

func NewWalletService(WalletRepository wallet.WalletRepository) WalletService {
	return &walletService{WalletRepository: WalletRepository}
}

//Notes: some business rules is solved on domain level thats why it is not checked again
//for example below zero balance rule

func (s *walletService) GetBalance(userID string) (*wallet.Wallet, error) {
	return s.WalletRepository.GetByUserID(userID)
}

func (s *walletService) CreditWallet(userID string, amount string) (*wallet.Wallet, error) {
	return s.WalletRepository.Credit(userID, amount)
}
func (s *walletService) DepositWallet(userID string, amount string) (*wallet.Wallet, error) {
	return s.WalletRepository.Deposit(userID, amount)
}
