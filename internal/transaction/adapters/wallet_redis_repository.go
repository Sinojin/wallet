package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sinojin/wallet/internal/transaction/domain/wallet"
)

type walletRedisRepository struct {
	walletRepo  wallet.WalletRepository
	redisClient *redis.Client
}

func NewRedisRepository(walletRepo wallet.WalletRepository, redisClient *redis.Client) *walletRedisRepository {
	return &walletRedisRepository{walletRepo: walletRepo, redisClient: redisClient}
}

//should be long duration but its okey for now

const DEFAULT_TIMEOUT_DURATION = time.Second * 60

//Prefix for keys should be good but it's not concerned because of this is test assigment
func (r *walletRedisRepository) GetByUserID(userID string) (*wallet.Wallet, error) {
	str, err := r.redisClient.Get(context.Background(), userID).Result()
	if err == redis.Nil {
		userWallet, err := r.walletRepo.GetByUserID(userID)
		if err != nil {
			return nil, err
		}
		// cache errors should'nt block normal repositories
		r.setCache(userID, userWallet)
		return userWallet, nil
	} else if err != nil {
		return nil, err
	}
	var UserWallet *wallet.Wallet
	err = json.Unmarshal([]byte(str), &UserWallet)
	if err != nil {
		return nil, err
	}
	return UserWallet, nil

}

func (r *walletRedisRepository) Deposit(userID string, amount string) (*wallet.Wallet, error) {
	userWallet, err := r.walletRepo.Deposit(userID, amount)
	if err != nil {
		return userWallet, err
	}
	r.setCache(userID, userWallet)
	return userWallet, err
}

func (r *walletRedisRepository) Credit(userID string, amount string) (*wallet.Wallet, error) {
	userWallet, err := r.walletRepo.Credit(userID, amount)
	if err != nil {
		return userWallet, err
	}
	r.setCache(userID, userWallet)
	return userWallet, err
}

func (r *walletRedisRepository) setCache(userID string, userWallet *wallet.Wallet) error {
	walletByteA, err := json.Marshal(userWallet)
	if err != nil {
		return err
	}
	_, err = r.redisClient.SetEX(context.Background(), userID, fmt.Sprintf("%s", walletByteA), DEFAULT_TIMEOUT_DURATION).Result()
	return err

}
