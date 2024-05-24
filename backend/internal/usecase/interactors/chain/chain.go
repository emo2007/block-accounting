package chain

import (
	"context"

	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/transactions"
)

type ChainInteractor interface {
}

type chainInteractor struct {
	txRepository    transactions.Repository
	usersInteractor users.UsersInteractor
}

type NewMultisigParams struct {
	OwnersPKs []string
}

func (i *chainInteractor) NewMultisig(ctx context.Context) {

}
