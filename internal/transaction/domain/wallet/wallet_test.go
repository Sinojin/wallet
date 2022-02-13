package wallet

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestWallet_Deposit(t *testing.T) {
	type fields struct {
		UserID  string
		Balance decimal.Decimal
	}
	type args struct {
		amount string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Deposit should be okey!",
			fields: fields{
				UserID:  "123",
				Balance: decimal.RequireFromString("100"),
			},
			args:    args{amount: "100"},
			wantErr: false,
		},
		{
			name: "Deposit should't be okey!",
			fields: fields{
				UserID:  "123",
				Balance: decimal.RequireFromString("100"),
			},
			args:    args{amount: "-100"},
			wantErr: true,
		},
		{
			name: "Deposit should be okey!",
			fields: fields{
				UserID:  "123",
				Balance: decimal.RequireFromString("100"),
			},
			args:    args{amount: "0"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Wallet{
				UserID:  tt.fields.UserID,
				Balance: tt.fields.Balance,
			}
			if err := b.Deposit(tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("Wallet.Deposit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWallet_Credit(t *testing.T) {
	type fields struct {
		UserID  string
		Balance decimal.Decimal
	}
	type args struct {
		amount string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Credit should'nt be okey!",
			fields: fields{
				UserID:  "123",
				Balance: decimal.RequireFromString("0"),
			},
			args:    args{amount: "100"},
			wantErr: true,
		},
		{
			name: "Negative Credit should'nt be okey!",
			fields: fields{
				UserID:  "123",
				Balance: decimal.RequireFromString("1000"),
			},
			args:    args{amount: "-100"},
			wantErr: true,
		},
		{
			name: "Deposit should be okey!",
			fields: fields{
				UserID:  "123",
				Balance: decimal.RequireFromString("101"),
			},
			args:    args{amount: "100"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Wallet{
				UserID:  tt.fields.UserID,
				Balance: tt.fields.Balance,
			}
			if err := b.Credit(tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("Wallet.Credit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
