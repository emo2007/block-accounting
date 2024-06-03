package chain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/emochka2007/block-accounting/internal/infrastructure/repository/transactions"
	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/logger"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/organizations"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type ChainInteractor interface {
	PubKey(ctx context.Context, user *models.User) ([]byte, error)

	NewMultisig(ctx context.Context, params NewMultisigParams) error
	ListMultisigs(ctx context.Context, params ListMultisigsParams) ([]models.Multisig, error)

	PayrollDeploy(ctx context.Context, params PayrollDeployParams) error
	ListPayrolls(ctx context.Context, params ListPayrollsParams) ([]models.Payroll, error)
	SetSalary(ctx context.Context, params SetSalaryParams) error
}

type chainInteractor struct {
	log                     *slog.Logger
	config                  config.Config
	txRepository            transactions.Repository
	organizationsInteractor organizations.OrganizationsInteractor
}

func NewChainInteractor(
	log *slog.Logger,
	config config.Config,
	txRepository transactions.Repository,
	organizationsInteractor organizations.OrganizationsInteractor,
) ChainInteractor {
	return &chainInteractor{
		log:                     log,
		config:                  config,
		txRepository:            txRepository,
		organizationsInteractor: organizationsInteractor,
	}
}

type NewMultisigParams struct {
	Title         string
	Owners        []models.OrganizationParticipant
	Confirmations int
}

type newMultisigChainResponse struct {
	Address string `json:"address"`
}

func (i *chainInteractor) NewMultisig(ctx context.Context, params NewMultisigParams) error {
	endpoint := i.config.ChainAPI.Host + "/multi-sig/deploy"

	i.log.Debug(
		"deploy multisig",
		slog.String("endpoint", endpoint),
		slog.Any("params", params),
	)

	pks := make([]string, len(params.Owners))

	for i, owner := range params.Owners {
		if owner.GetUser() == nil {
			return fmt.Errorf("error invalis owners set")
		}

		pks[i] = "0x" + common.Bytes2Hex(owner.GetUser().PublicKey())
	}

	requestBody, err := json.Marshal(map[string]any{
		"owners":        pks,
		"confirmations": params.Confirmations,
	})
	if err != nil {
		return fmt.Errorf("error marshal request body. %w", err)
	}

	user, err := ctxmeta.User(ctx)
	if err != nil {
		return fmt.Errorf("error fetch user from context. %w", err)
	}

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return fmt.Errorf("error fetch organization id from context. %w", err)
	}

	go func() { // TODO remove this subroutine shit and replace it with worker pools
		pid := uuid.Must(uuid.NewV7()).String()
		startTime := time.Now()

		i.log.Info(
			"new multisig worker started",
			slog.String("pid", pid),
		)

		doneCh := make(chan struct{}, 1)

		defer func() {
			if err := recover(); err != nil {
				i.log.Error("worker paniced!", slog.Any("panic", err))
			}

			doneCh <- struct{}{}
			close(doneCh)
		}()

		go func() {
			warn := time.After(1 * time.Minute)
			select {
			case <-doneCh:
				i.log.Info(
					"new multisig worker done",
					slog.String("pid", pid),
					slog.Time("started at", startTime),
					slog.Time("done at", time.Now()),
					slog.Duration("work time", time.Since(startTime)),
				)
			case <-warn:
				i.log.Warn(
					"new multisig worker seems sleeping",
					slog.String("pid", pid),
					slog.Duration("work time", time.Since(startTime)),
				)
			}
		}()

		requestContext, cancel := context.WithTimeout(context.TODO(), time.Minute*15)
		defer cancel()

		body := bytes.NewBuffer(requestBody)

		req, err := http.NewRequestWithContext(requestContext, http.MethodPost, endpoint, body)
		if err != nil {
			i.log.Error(
				"error build request",
				logger.Err(err),
			)

			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Seed", common.Bytes2Hex(user.Seed()))

		// TODO replace http.DefaultClient with custom ChainAPi client from infrastructure/ethapi mod
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			i.log.Error(
				"error send deploy multisig request",
				slog.String("endpoint", endpoint),
				slog.Any("params", params),
			)

			return
		}

		defer resp.Body.Close()

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			i.log.Error(
				"error read body",
				logger.Err(err),
			)

			return
		}

		respObject := new(newMultisigChainResponse)

		if err := json.Unmarshal(raw, &respObject); err != nil {
			i.log.Error(
				"error parse chain-api response body",
				logger.Err(err),
			)

			return
		}

		if respObject.Address == "" {
			i.log.Error(
				"error multisig address is empty",
			)

			return
		}

		multisigAddress := common.Hex2Bytes(respObject.Address[2:])

		createdAt := time.Now()

		msg := models.Multisig{
			ID:                    uuid.Must(uuid.NewV7()),
			Title:                 params.Title,
			Address:               multisigAddress,
			OrganizationID:        organizationID,
			Owners:                params.Owners,
			ConfirmationsRequired: params.Confirmations,
			CreatedAt:             createdAt,
			UpdatedAt:             createdAt,
		}

		i.log.Debug(
			"deploy multisig response",
			slog.Int("code", resp.StatusCode),
			slog.String("body", string(raw)),
			slog.Any("parsed", respObject),
			slog.Any("multisig object", msg),
		)

		if err := i.txRepository.AddMultisig(requestContext, msg); err != nil {
			i.log.Error(
				"error add new multisig",
				logger.Err(err),
			)

			return
		}
	}()

	return nil
}

func (i *chainInteractor) PubKey(ctx context.Context, user *models.User) ([]byte, error) {
	pubAddr := i.config.ChainAPI.Host + "/address-from-seed"

	requestBody, err := json.Marshal(map[string]any{
		"seedPhrase": user.Mnemonic,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshal request body. %w", err)
	}

	body := bytes.NewBuffer(requestBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pubAddr, body)
	if err != nil {
		return nil, fmt.Errorf("error build request. %w", err)
	}

	req.Header.Add("X-Seed", common.Bytes2Hex(user.Seed()))
	req.Header.Add("Content-Type", "application/json")

	// TODO replace http.DefaultClient with custom ChainAPi client from infrastructure/ethapi mod
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetch pub address. %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read resp body. %w", err)
	}

	pubKeyStr := string(respBody)[2:]

	if pubKeyStr == "" {
		return nil, fmt.Errorf("error empty public key")
	}

	return common.Hex2Bytes(pubKeyStr), nil
}

type PayrollDeployParams struct {
	FirstAdmin models.OrganizationParticipant
	MultisigID uuid.UUID
	Title      string
}

type newPayrollContractChainResponse struct {
	Address string `json:"address"`
}

func (i *chainInteractor) PayrollDeploy(
	ctx context.Context,
	params PayrollDeployParams,
) error {
	user, err := ctxmeta.User(ctx)
	if err != nil {
		return fmt.Errorf("error fetch user from context. %w", err)
	}

	if user.Id() != params.FirstAdmin.Id() || params.FirstAdmin.GetUser() == nil {
		return fmt.Errorf("error unauthorized access")
	}

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return fmt.Errorf("error fetch organization id from context. %w", err)
	}

	multisigs, err := i.ListMultisigs(ctx, ListMultisigsParams{
		OrganizationID: organizationID,
		IDs:            uuid.UUIDs{params.MultisigID},
	})
	if err != nil {
		return fmt.Errorf("error fetch multisigs by id. %w", err)
	}

	if len(multisigs) == 0 {
		return fmt.Errorf("error empty multisigs set")
	}

	i.log.Debug(
		"PayrollDeploy",
		slog.String("organization id", organizationID.String()),
		slog.String("multisig id", params.MultisigID.String()),
		slog.String("multisig address", common.Bytes2Hex(multisigs[0].Address)),
		slog.String("X-Seed header data", common.Bytes2Hex(user.Seed())),
	)

	maddr := common.Bytes2Hex(multisigs[0].Address)

	if maddr == "" {
		return fmt.Errorf("empty multisig address")
	}

	if maddr[0] != 0 && maddr[1] != 'x' {
		maddr = "0x" + maddr
	}

	requestBody, err := json.Marshal(map[string]any{
		"authorizedWallet": maddr,
	})
	if err != nil {
		return fmt.Errorf("error marshal request body. %w", err)
	}

	go func() { // TODO remove this subroutine shit and replace it with worker pools
		pid := uuid.Must(uuid.NewV7()).String()
		startTime := time.Now()

		i.log.Info(
			"new paroll worker started",
			slog.String("pid", pid),
		)

		doneCh := make(chan struct{}, 1)

		defer func() {
			if err := recover(); err != nil {
				i.log.Error("worker paniced!", slog.Any("panic", err))
			}

			doneCh <- struct{}{}
			close(doneCh)
		}()

		go func() {
			warn := time.After(2 * time.Minute)
			select {
			case <-doneCh:
				i.log.Info(
					"new payroll worker done",
					slog.String("pid", pid),
					slog.Time("started at", startTime),
					slog.Time("done at", time.Now()),
					slog.Duration("work time", time.Since(startTime)),
				)
			case <-warn:
				i.log.Warn(
					"new paroll worker seems sleeping",
					slog.String("pid", pid),
					slog.Duration("work time", time.Since(startTime)),
				)
			}
		}()

		requestContext, cancel := context.WithTimeout(context.TODO(), time.Minute*20)
		defer cancel()

		body := bytes.NewBuffer(requestBody)

		endpoint := i.config.ChainAPI.Host + "/salaries/deploy"

		i.log.Debug(
			"request",
			slog.String("body", string(requestBody)),
			slog.String("endpoint", endpoint),
		)

		req, err := http.NewRequestWithContext(requestContext, http.MethodPost, endpoint, body)
		if err != nil {
			i.log.Error(
				"error build request",
				logger.Err(fmt.Errorf("error build request. %w", err)),
			)
			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Seed", common.Bytes2Hex(user.Seed()))

		// TODO replace http.DefaultClient with custom ChainAPi client from infrastructure/ethapi mod
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			i.log.Error(
				"error fetch deploy salary contract",
				logger.Err(err),
			)

			return
		}

		defer resp.Body.Close()

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			i.log.Error(
				"error read body",
				logger.Err(err),
			)

			return
		}

		respObject := new(newPayrollContractChainResponse)

		if err := json.Unmarshal(raw, &respObject); err != nil {
			i.log.Error(
				"error parse chain-api response body",
				logger.Err(err),
			)

			return
		}

		i.log.Debug(
			"payroll deploy",
			slog.Any("response", respObject),
		)

		if respObject.Address == "" {
			i.log.Error(
				"error multisig address is empty",
			)

			return
		}

		addr := common.Hex2Bytes(respObject.Address[2:])

		createdAt := time.Now()

		if err := i.txRepository.AddPayrollContract(requestContext, transactions.AddPayrollContract{
			ID:             uuid.Must(uuid.NewV7()),
			Title:          params.Title,
			Address:        addr,
			OrganizationID: organizationID,
			MultisigID:     params.MultisigID,
			CreatedAt:      createdAt,
		}); err != nil {
			i.log.Error(
				"error add new payroll contract",
				logger.Err(err),
			)

			return
		}
	}()

	return nil
}

type ListMultisigsParams struct {
	IDs            uuid.UUIDs
	OrganizationID uuid.UUID
}

func (i *chainInteractor) ListMultisigs(
	ctx context.Context,
	params ListMultisigsParams,
) ([]models.Multisig, error) {
	multisigs, err := i.txRepository.ListMultisig(ctx, transactions.ListMultisigsParams{
		IDs:            params.IDs,
		OrganizationID: params.OrganizationID,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch multisigs. %w", err)
	}

	return multisigs, nil
}

type ListPayrollsParams struct {
	IDs            []uuid.UUID
	Limit          int
	OrganizationID uuid.UUID
}

func (i *chainInteractor) ListPayrolls(
	ctx context.Context,
	params ListPayrollsParams,
) ([]models.Payroll, error) {
	payrolls, err := i.txRepository.ListPayrolls(ctx, transactions.ListPayrollsParams{
		IDs:            params.IDs,
		Limit:          int64(params.Limit),
		OrganizationID: params.OrganizationID,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch payrolls from repository. %w", err)
	}

	return payrolls, nil
}

type SetSalaryParams struct {
	PayrollID  uuid.UUID
	EmployeeID uuid.UUID
	Salary     float64
}

func (i *chainInteractor) SetSalary(
	ctx context.Context,
	params SetSalaryParams,
) error {
	user, err := ctxmeta.User(ctx)
	if err != nil {
		return fmt.Errorf("error fetch user from context. %w", err)
	}

	organizationID, err := ctxmeta.OrganizationId(ctx)
	if err != nil {
		return fmt.Errorf("error fetch organization id from context. %w", err)
	}

	i.log.Debug(
		"SetSalary",
		slog.String("org id", organizationID.String()),
		slog.Any("user", user),
	)

	payrolls, err := i.ListPayrolls(ctx, ListPayrollsParams{
		IDs:            []uuid.UUID{params.PayrollID},
		OrganizationID: organizationID,
	})
	if err != nil {
		return fmt.Errorf("error fetch payroll. %w", err)
	}

	if len(payrolls) == 0 {
		return fmt.Errorf("error payroll not found. %w", err)
	}

	payroll := payrolls[0]

	multisigs, err := i.ListMultisigs(ctx, ListMultisigsParams{
		IDs:            uuid.UUIDs{payroll.MultisigID},
		OrganizationID: organizationID,
	})
	if err != nil {
		return fmt.Errorf("error fetch multisig. %w", err)
	}

	if len(multisigs) == 0 {
		return fmt.Errorf("error multisig not found. %w", err)
	}

	multisig := multisigs[0]

	employee, err := i.organizationsInteractor.Participant(ctx, organizations.ParticipantParams{
		ID:             params.EmployeeID,
		OrganizationID: organizationID,
		EmployeesOnly:  true,
	})
	if err != nil {
		return fmt.Errorf("error fetch employee from repository. %w", err)
	}

	if employee.GetEmployee() == nil {
		return fmt.Errorf("error employee is nil")
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				i.log.Error("worker paniced!", slog.Any("panic", err))
			}

			// doneCh <- struct{}{}
			// close(doneCh)
		}()

		ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Minute)
		defer cancel()

		maddr := common.Bytes2Hex(multisig.Address)
		if maddr[0] != 0 && maddr[1] != 'x' {
			maddr = "0x" + maddr
		}

		caddr := common.Bytes2Hex(payroll.Address)
		if caddr[0] != 0 && caddr[1] != 'x' {
			caddr = "0x" + caddr
		}

		eaddr := common.Bytes2Hex(employee.GetEmployee().WalletAddress)
		if caddr[0] != 0 && caddr[1] != 'x' {
			caddr = "0x" + caddr
		}

		bodyMap := map[string]any{
			"multiSigWallet":  maddr,
			"contractAddress": caddr,
			"employeeAddress": eaddr,
			"salary":          params.Salary,
		}

		bodyRaw, err := json.Marshal(&bodyMap)
		if err != nil {
			i.log.Error(
				"error marshal request body",
				logger.Err(err),
			)

			return
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			i.config.ChainAPI.Host+"/salaries/set-salary",
			bytes.NewBuffer(bodyRaw),
		)
		if err != nil {
			i.log.Error(
				"error build request",
				logger.Err(err),
			)

			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Seed", common.Bytes2Hex(user.Seed()))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			i.log.Error(
				"error do request",
				logger.Err(err),
			)

			return
		}

		defer resp.Body.Close()

		// todo parse body
	}()

	return nil
}
