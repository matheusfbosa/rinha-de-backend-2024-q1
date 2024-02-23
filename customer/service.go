package customer

import "time"

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) MakeTransaction(tr *Transaction) (*AccountBalance, error) {
	limit, ok := CustomerAccountLimit[tr.CustomerID]
	if !ok {
		return nil, ErrCustomerNotFound
	}

	tr.AccountLimit = limit
	balance, err := s.r.MakeTransaction(tr)
	if err != nil {
		return nil, err
	}

	return &AccountBalance{
		Limit:   limit,
		Balance: balance,
	}, nil
}

func (s *Service) GetBankStatement(customerID string) (*BankStatement, error) {
	limit, ok := CustomerAccountLimit[customerID]
	if !ok {
		return nil, ErrCustomerNotFound
	}

	statement, err := s.r.GetBankStatement(customerID)
	if err != nil {
		return nil, err
	}

	statement.Balance.Limit = limit
	statement.Balance.Date = time.Now()

	return statement, nil
}
