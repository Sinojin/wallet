package wallet

type WalletRepository interface {
	GetByUserID(userID string) (*Wallet, error)
	Deposit(userID string, amount string) (*Wallet, error)
	Credit(userID string, amount string) (*Wallet, error)
}
