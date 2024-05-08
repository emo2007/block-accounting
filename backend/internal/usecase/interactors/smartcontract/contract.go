package smartcontract

import (
	"context"

	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/transactions"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

type SignTransactionParams struct {
	Signee         *models.User
	TxId           uuid.UUID
	OrganizationId uuid.UUID
}

type SmartContractInteractor interface {
	SignTransaction(ctx context.Context, params SignTransactionParams) error
}

type smartContractInteractor struct {
	client ethclient.Client

	scAddr string

	txRepository    transactions.Repository
	usersInteractor users.UsersInteractor
}

// todo
func (s *smartContractInteractor) SignTransaction(ctx context.Context, params SignTransactionParams) error {
	// s.client.CallContract(ctx, ethereum.CallMsg{}, nil)

	return nil
}
