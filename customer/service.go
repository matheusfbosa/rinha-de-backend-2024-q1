package customer

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) MakeTransaction(tr *Transaction) (*AccountBalance, error) {
	balance, err := s.r.GetAccountBalance(tr.CustomerID)
	if err != nil {
		return nil, err
	}

	newTr, err := tr.withLastBalance(balance.Total, balance.Limit)
	if err != nil {
		return nil, err
	}

	err = s.r.CreateTransaction(newTr)
	if err != nil {
		return nil, err
	}

	return &AccountBalance{
		Limit:   balance.Limit,
		Balance: newTr.LastBalance,
	}, nil
}

func (s *Service) GetBankStatement(customerID string) (*BankStatement, error) {
	statement, err := s.r.GetBankStatement(customerID)
	if err != nil {
		return nil, err
	}

	return statement, nil
}
