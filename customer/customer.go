package customer

import (
	"errors"
	"strconv"
)

var (
	ErrInvalidTransaction = errors.New("invalid transaction")
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
)

const (
	DebitType  string = "d"
	CreditType        = "c"
)

type (
	Transaction struct {
		CustomerID  string
		Value       int    `json:"valor"`
		Type        string `json:"tipo"`
		Description string `json:"descricao"`
		LastBalance int
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
		Total int    `json:"total"`
		Date  string `json:"data_extrato"`
		Limit int    `json:"limite"`
	}

	TransactionBankStatement struct {
		Value       int    `json:"valor"`
		Type        string `json:"tipo"`
		Description string `json:"descricao"`
		CreatedAt   string `json:"realizada_em"`
	}

	UseCase interface {
		MakeTransaction(tr *Transaction) (*AccountBalance, error)
		GetBankStatement(customerID string) (*BankStatement, error)
	}

	Repository interface {
		CreateTransaction(tr *Transaction) error
		GetAccountBalance(customerID string) (*BalanceBankStatement, error)
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

func (tr Transaction) withLastBalance(balance, limit int) (*Transaction, error) {
	newTr := tr

	switch tr.Type {
	case CreditType:
		newTr.LastBalance = balance + tr.Value
	case DebitType:
		if (balance - tr.Value) < -limit {
			return nil, ErrInsufficientFunds
		}
		newTr.LastBalance = balance - tr.Value
	default:
		return nil, ErrInvalidTransaction
	}

	return &newTr, nil
}
