package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
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

func (r *PostgreSQL) CreateTransaction(tr *customer.Transaction) error {
	query := `
		INSERT INTO transactions (
			type,
			value,
			description,
			customer_id,
			last_balance
		) VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.dbpool.Exec(context.Background(),
		query,
		tr.Type,
		tr.Value,
		tr.Description,
		tr.CustomerID,
		tr.LastBalance,
	)
	pqErr, isPQError := err.(*pq.Error)
	if isPQError && pqErr.Code == foreignKeyViolationCode {
		return customer.ErrCustomerNotFound
	}

	return err
}

func (r *PostgreSQL) GetBankStatement(customerID string) (*customer.BankStatement, error) {
	balance, err := r.GetAccountBalance(customerID)
	if err != nil {
		return nil, err
	}

	lastTransactions, err := r.getLastTransactions(customerID)
	if err != nil {
		return nil, err
	}

	return &customer.BankStatement{
		Balance:          *balance,
		LastTransactions: lastTransactions,
	}, nil
}

func (r *PostgreSQL) GetAccountBalance(customerID string) (*customer.BalanceBankStatement, error) {
	query := `
		SELECT
			COALESCE(t.last_balance, 0) AS total,
			c.account_limit AS limit,
			NOW() AS date
		FROM customers c
		LEFT JOIN (
			SELECT customer_id, last_balance
			FROM transactions
			WHERE customer_id = $1
			ORDER BY created_at DESC
			LIMIT 1
		) t ON c.customer_id = t.customer_id
		WHERE c.customer_id = $1 AND c.customer_id IS NOT NULL;
	`
	var ab customer.BalanceBankStatement
	err := r.dbpool.QueryRow(context.Background(), query, customerID).
		Scan(&ab.Total, &ab.Limit, &ab.Date)
	if err == sql.ErrNoRows {
		return nil, customer.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}

	return &ab, nil
}

func (r *PostgreSQL) getLastTransactions(customerID string) ([]*customer.TransactionBankStatement, error) {
	query := `
		SELECT
			value,
			type,
			description,
			created_at
		FROM transactions
		WHERE customer_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.dbpool.Query(context.Background(), query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]*customer.TransactionBankStatement, 0)
	for rows.Next() {
		var tr customer.TransactionBankStatement
		err := rows.Scan(&tr.Value, &tr.Type, &tr.Description, &tr.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &tr)
	}

	return transactions, nil
}
