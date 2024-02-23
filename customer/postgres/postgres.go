package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matheusfbosa/rinha-de-backend-2024-q1/customer"
)

const foreignKeyViolationCode = "23503"

type PostgreSQL struct {
	dbpool *pgxpool.Pool
}

func NewPostgreSQL(dbpool *pgxpool.Pool) *PostgreSQL {
	return &PostgreSQL{
		dbpool: dbpool,
	}
}

func (r *PostgreSQL) MakeTransaction(tr *customer.Transaction) (int, error) {
	query := `
		select make_transaction($1, $2, $3, $4, $5);
	`
	var balance int
	err := r.dbpool.QueryRow(context.Background(),
		query,
		tr.CustomerID,
		tr.Type,
		tr.Value,
		tr.Description,
		tr.AccountLimit,
	).Scan(&balance)
	if err != nil {
		if isAccountLimitError(err) {
			return 0, customer.ErrInsufficientFunds
		}

		return 0, err
	}

	return balance, nil
}

func (r *PostgreSQL) GetBankStatement(customerID string) (*customer.BankStatement, error) {
	query := `
		select
			last_balance,
			type,
			value,
			description,
			created_at
		from transactions
		where customer_id = $1
		order by created_at desc
		limit 10
    `
	rows, err := r.dbpool.Query(context.Background(), query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var idx int
	var lastBalance int
	var statement customer.BankStatement
	statement.LastTransactions = make([]*customer.TransactionBankStatement, 0)
	for rows.Next() {
		var tr customer.TransactionBankStatement
		err := rows.Scan(
			&lastBalance,
			&tr.Type,
			&tr.Value,
			&tr.Description,
			&tr.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if idx == 0 {
			statement.Balance.Total = lastBalance
		}

		statement.LastTransactions = append(statement.LastTransactions, &tr)
		idx++
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &statement, nil
}

func isAccountLimitError(err error) bool {
	return err.Error() == "ERROR: insufficient funds (SQLSTATE P0001)"
}
