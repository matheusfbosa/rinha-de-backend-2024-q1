package customer

import (
	"errors"
	"strconv"
	"time"
)

var (
	ErrInvalidTransaction = errors.New("invalid transaction")
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")

	CustomerAccountLimit = map[string]int{
		"1": 100000,
		"2": 80000,
		"3": 1000000,
		"4": 10000000,
		"5": 500000,
	}
)

const (
	DebitType  string = "d"
	CreditType        = "c"
)

type (
	Transaction struct {
		Value        int    `json:"valor"`
		Type         string `json:"tipo"`
		Description  string `json:"descricao"`
		CustomerID   string
		AccountLimit int
	}

	AccountBalance struct {
		Limit   int `json:"limite"`
		Balance int `json:"saldo"`
	}

	BankStatement struct {
		Balance          BalanceBankStatement        `json:"saldo"`
		LastTransactions []*TransactionBankStatement `json:"ultimas_transacoes"`
	}

	BalanceBankStatement struct {
		Total int       `json:"total"`
		Date  time.Time `json:"data_extrato"`
		Limit int       `json:"limite"`
	}

	TransactionBankStatement struct {
		Value       int       `json:"valor"`
		Type        string    `json:"tipo"`
		Description string    `json:"descricao"`
		CreatedAt   time.Time `json:"realizada_em"`
	}

	UseCase interface {
		MakeTransaction(tr *Transaction) (*AccountBalance, error)
		GetBankStatement(customerID string) (*BankStatement, error)
	}

	Repository interface {
		MakeTransaction(tr *Transaction) (int, error)
		GetBankStatement(customerID string) (*BankStatement, error)
	}
)

func (tr *Transaction) Validate() error {
	if tr.Value <= 0 || (tr.Type != DebitType && tr.Type != CreditType) ||
		tr.Description == "" || tr.CustomerID == "" {
		return ErrInvalidTransaction
	}
	if customerID, err := strconv.Atoi(tr.CustomerID); err != nil || customerID <= 0 {
		return ErrInvalidTransaction
	}

	return nil
}
