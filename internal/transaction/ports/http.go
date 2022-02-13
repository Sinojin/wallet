package ports

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinojin/wallet/internal/transaction/services"
)

type Balance struct {
	Balance string `json:"balance"`
}

const URL_PARAM_USER_ID = "wallet_id"

type TransactionHandler struct {
	walletService services.WalletService
}

func NewTransactionHandler(walletService services.WalletService) TransactionHandler {
	return TransactionHandler{walletService: walletService}
}
func (h TransactionHandler) GetBalance(c *gin.Context) {
	userID := c.Param(URL_PARAM_USER_ID)
	userWallet, err := h.walletService.GetBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Balance": userWallet.Balance.String(),
	})
}

func (h TransactionHandler) CreditWallet(c *gin.Context) {
	userID := c.Param(URL_PARAM_USER_ID)
	var request Balance
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
		})
		return
	}
	userWallet, err := h.walletService.CreditWallet(userID, request.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Balance": userWallet.Balance.String(),
	})
}

func (h TransactionHandler) DepositWallet(c *gin.Context) {
	userID := c.Param(URL_PARAM_USER_ID)
	var request Balance
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
		})
		return
	}
	userWallet, err := h.walletService.DepositWallet(userID, request.Balance)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Balance": userWallet.Balance.String(),
	})
}

func (h *TransactionHandler) Handle(router *gin.RouterGroup) {
	wallets := router.Group("/wallets")
	//get balance of wallet
	wallets.GET(fmt.Sprintf("/:%v/balance", URL_PARAM_USER_ID), h.GetBalance)
	//give money to wallet
	wallets.POST(fmt.Sprintf("/:%v/credit", URL_PARAM_USER_ID), h.DepositWallet)
	//take money from wallet.
	wallets.POST(fmt.Sprintf("/:%v/debit", URL_PARAM_USER_ID), h.CreditWallet)
}
