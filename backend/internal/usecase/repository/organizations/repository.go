package organizations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	sqltools "github.com/emochka2007/block-accounting/internal/pkg/sqlutils"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/users"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

var (
	ErrorNotFound = errors.New("not found")
)

type GetParams struct {
	Ids    uuid.UUIDs
	UserId uuid.UUID

	OffsetDate time.Time
	CursorId   uuid.UUID
	Limit      int64
}

type ParticipantsParams struct {
	OrganizationId uuid.UUID
	Ids            uuid.UUIDs

	// Filters
	UsersOnly     bool
	ActiveOnly    bool
	EmployeesOnly bool
}

type AddParticipantParams struct {
	OrganizationId uuid.UUID
	UserId         uuid.UUID
	EmployeeId     uuid.UUID
	IsAdmin        bool
}

type DeleteParticipantParams struct {
	OrganizationId uuid.UUID
	UserId         uuid.UUID
	EmployeeId     uuid.UUID
}

type Repository interface {
	Create(ctx context.Context, org models.Organization) error
	Get(ctx context.Context, params GetParams) ([]*models.Organization, error)
	Update(ctx context.Context, org models.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddParticipant(ctx context.Context, params AddParticipantParams) error
	Participants(ctx context.Context, params ParticipantsParams) ([]models.OrganizationParticipant, error)
	CreateAndAdd(ctx context.Context, org models.Organization, user *models.User) error
	DeleteParticipant(ctx context.Context, params DeleteParticipantParams) error
}

type repositorySQL struct {
	db              *sql.DB
	usersRepository users.Repository
}

func NewRepository(
	db *sql.DB,
	usersRepository users.Repository,
) Repository {
	return &repositorySQL{
		db:              db,
		usersRepository: usersRepository,
	}
}

func (s *repositorySQL) Conn(ctx context.Context) sqltools.DBTX {
	if tx, ok := ctx.Value(sqltools.TxCtxKey).(*sql.Tx); ok {
		return tx
	}

	return s.db
}

func (r *repositorySQL) Create(ctx context.Context, org models.Organization) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Insert("organizations").Columns(
			"id",
			"name",
			"address",
			"wallet_seed",
			"created_at",
			"updated_at",
		).Values(
			org.ID,
			org.Name,
			org.Address,
			org.WalletSeed,
			org.CreatedAt,
			org.UpdatedAt,
		).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error insert new organization. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) Get(ctx context.Context, params GetParams) ([]*models.Organization, error) {
	organizations := make([]*models.Organization, 0, params.Limit)

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Select(
			"o.id",
			"o.name",
			"o.address",
			"o.wallet_seed",
			"o.created_at",
			"o.updated_at",
		).From("organizations as o").
			Limit(uint64(params.Limit)).
			PlaceholderFormat(sq.Dollar)

		if params.UserId != uuid.Nil {
			query = query.InnerJoin("organizations_users as ou on o.id = ou.organization_id").
				Where(sq.Eq{
					"ou.user_id": params.UserId,
				})
		}

		if params.CursorId != uuid.Nil {
			query = query.Where(sq.Gt{
				"o.id": params.CursorId,
			})
		}

		if params.Ids != nil {
			query = query.Where(sq.Eq{
				"o.id": params.Ids,
			})
		}

		if !params.OffsetDate.IsZero() {
			query = query.Where(sq.GtOrEq{
				"o.updated_at": params.OffsetDate,
			})
		}

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch organizations from database. %w", err)
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				err = errors.Join(fmt.Errorf("error close rows. %w", closeErr), err)
			}
		}()

		for rows.Next() {
			var (
				id         uuid.UUID
				name       string
				address    string
				walletSeed []byte
				createdAt  time.Time
				updatedAt  time.Time
			)

			if err = rows.Scan(
				&id,
				&name,
				&address,
				&walletSeed,
				&createdAt,
				&updatedAt,
			); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}

			organizations = append(organizations, &models.Organization{
				ID:         id,
				Name:       name,
				Address:    address,
				WalletSeed: walletSeed,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error execute transactional operation. %w", err)
	}

	return organizations, nil
}

func (r *repositorySQL) Update(ctx context.Context, org models.Organization) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Update("organizations as o").
			SetMap(sq.Eq{
				"o.name":        org.Name,
				"o.address":     org.Address,
				"o.wallet_seed": org.WalletSeed,
				"o.created_at":  org.CreatedAt,
				"o.updated_at":  org.UpdatedAt,
			}).
			Where(sq.Eq{
				"o.id": org.ID,
			}).
			PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error update organization. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) Delete(ctx context.Context, id uuid.UUID) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Delete("organizations as o").
			Where(sq.Eq{
				"o.id": id,
			}).
			PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error delete organization. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) CreateAndAdd(ctx context.Context, org models.Organization, user *models.User) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		if err := r.Create(ctx, org); err != nil {
			return fmt.Errorf("error create organization. %w", err)
		}

		if err := r.AddParticipant(ctx, AddParticipantParams{
			OrganizationId: org.ID,
			UserId:         user.Id(),
			IsAdmin:        true,
		}); err != nil {
			return fmt.Errorf("error add user to newly created organization. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) AddParticipant(ctx context.Context, params AddParticipantParams) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Insert("organizations_users").
			Columns(
				"organization_id",
				"user_id",
				"employee_id",
				"added_at",
				"updated_at",
				"is_admin",
			).
			Values(
				params.OrganizationId,
				params.UserId,
				params.EmployeeId,
				time.Now(),
				time.Now(),
				params.IsAdmin,
			).
			PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error add new participant to organization. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) DeleteParticipant(ctx context.Context, params DeleteParticipantParams) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		deletedAt := time.Now()

		query := sq.Update("organizations_users as ou").
			SetMap(sq.Eq{
				"updated_at": deletedAt,
				"deleted_at": deletedAt,
				"is_admin":   false,
			}).
			Where(sq.Eq{
				"ou.organization_id": params.OrganizationId,
			}).
			PlaceholderFormat(sq.Dollar)

		if params.EmployeeId != uuid.Nil {
			query = query.Where(sq.Eq{
				"ou.employee_id": params.EmployeeId,
			})
		}

		if params.UserId != uuid.Nil {
			query = query.Where(sq.Eq{
				"ou.user_id": params.UserId,
			})
		}

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error delete participant from organization. %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error execute transactional operation. %w", err)
	}

	return nil
}

func (r *repositorySQL) Participants(
	ctx context.Context,
	params ParticipantsParams,
) ([]models.OrganizationParticipant, error) {
	participants := make([]models.OrganizationParticipant, 0, len(params.Ids))

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		orgUsersModels, err := r.fetchOrganizationUsers(ctx, params)
		if err != nil {
			return fmt.Errorf("error fetch organization users raw models. %w", err)
		}

		eg, egCtx := errgroup.WithContext(ctx)

		var employees []*models.Employee = make([]*models.Employee, 0, len(orgUsersModels))
		if !params.UsersOnly {
			eg.Go(func() error {
				ids := make(uuid.UUIDs, 0, len(orgUsersModels))

				for _, m := range orgUsersModels {
					if m.employeeID != uuid.Nil {
						ids = append(ids, m.employeeID)
					}
				}

				employees, err = r.fetchEmployees(egCtx, fetchEmployeesParams{
					IDs:            ids,
					OrganizationId: params.OrganizationId,
				})
				if err != nil {
					return fmt.Errorf("error fetch employees. %w", err)
				}

				return nil
			})
		}

		var usrs []*models.User
		if !params.EmployeesOnly {
			eg.Go(func() error {
				ids := make(uuid.UUIDs, 0, len(orgUsersModels))

				for _, m := range orgUsersModels {
					if m.userID != uuid.Nil {
						ids = append(ids, m.userID)
					}
				}

				usrs, err = r.usersRepository.Get(egCtx, users.GetParams{
					Ids: ids,
				})
				if err != nil {
					return fmt.Errorf("error fetch users by ids. %w", err)
				}

				return nil
			})
		}

		if err = eg.Wait(); err != nil {
			return fmt.Errorf("error organizations users entitiels. %w", err)
		}

		for _, ou := range orgUsersModels {
			var employee *models.Employee

			if ou.employeeID != uuid.Nil {
				for _, e := range employees {
					if e.ID == ou.employeeID {
						employee = e

						break
					}
				}
			}

			if ou.userID == uuid.Nil && employee != nil {
				participants = append(participants, employee)
			}

			for _, u := range usrs {
				if u.Id() == ou.userID {
					participants = append(participants, &models.OrganizationUser{
						User:        *u,
						OrgPosition: ou.position,
						Admin:       ou.isAdmin,
						Employee:    employee,
					})

					break
				}
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error execute transactional operation. %w", err)
	}

	if len(participants) == 0 {
		return nil, ErrorNotFound
	}

	return participants, nil
}

type fetchOrganizationUsersModel struct {
	organizationID uuid.UUID
	userID         uuid.UUID
	employeeID     uuid.UUID
	position       string
	addedAt        time.Time
	updatedAt      time.Time
	deletedAt      time.Time
	isAdmin        bool
}

func (r *repositorySQL) fetchOrganizationUsers(
	ctx context.Context,
	params ParticipantsParams,
) ([]fetchOrganizationUsersModel, error) {
	participants := make([]fetchOrganizationUsersModel, 0, len(params.Ids))

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		ouQuery := sq.Select(
			"ou.organization_id",
			"ou.user_id",
			"ou.employee_id",
			"ou.position",
			"ou.added_at",
			"ou.updated_at",
			"ou.deleted_at",
			"ou.is_admin",
		).Where(sq.Eq{
			"ou.organization_id": params.OrganizationId,
		}).From("organizations_users as ou").
			PlaceholderFormat(sq.Dollar)

		if len(params.Ids) > 0 {
			ouQuery = ouQuery.Where(sq.Eq{
				"ou.user_id": params.Ids,
			})
		}

		fmt.Println(ouQuery.ToSql())

		rows, err := ouQuery.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch organization participants. %w", err)
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				err = errors.Join(fmt.Errorf("error close rows. %w", closeErr), err)
			}
		}()

		for rows.Next() {
			var (
				organizationID uuid.UUID
				userID         uuid.UUID
				employeeID     uuid.UUID
				position       sql.NullString
				addedAt        time.Time
				updatedAt      time.Time
				deletedAt      sql.NullTime
				isAdmin        bool
			)

			if err = rows.Scan(
				&organizationID,
				&userID,
				&employeeID,
				&position,
				&addedAt,
				&updatedAt,
				&deletedAt,
				&isAdmin,
			); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}

			if params.EmployeesOnly && employeeID == uuid.Nil {
				continue
			}

			if params.UsersOnly && userID == uuid.Nil {
				continue
			}

			if params.ActiveOnly && deletedAt.Valid {
				continue
			}

			participants = append(participants, fetchOrganizationUsersModel{
				organizationID: organizationID,
				userID:         userID,
				employeeID:     employeeID,
				position:       position.String,
				addedAt:        addedAt,
				updatedAt:      updatedAt,
				deletedAt:      deletedAt.Time,
				isAdmin:        isAdmin,
			})
		}

		return err
	}); err != nil {
		return nil, fmt.Errorf("error execute transactional operation. %w", err)
	}

	return participants, nil
}

type fetchEmployeesParams struct {
	IDs            uuid.UUIDs
	OrganizationId uuid.UUID
}

func (r *repositorySQL) fetchEmployees(
	ctx context.Context,
	params fetchEmployeesParams,
) ([]*models.Employee, error) {
	employees := make([]*models.Employee, 0, len(params.IDs))

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) (err error) {
		query := sq.Select(
			"e.id",
			"e.user_id",
			"e.organization_id",
			"e.wallet_address",
			"e.created_at",
			"e.updated_at",
		).Where(sq.Eq{
			"e.organization_id": params.OrganizationId,
		}).From("employees as e").
			PlaceholderFormat(sq.Dollar)

		if len(params.IDs) > 0 {
			query = query.Where(sq.Eq{
				"e.id": params.IDs,
			})
		}

		fmt.Println(query.ToSql())

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch employees from database. %w", err)
		}

		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				err = errors.Join(fmt.Errorf("error close rows. %w", closeErr), err)
			}
		}()

		for rows.Next() {
			var (
				id         uuid.UUID
				userID     uuid.UUID
				orgID      uuid.UUID
				walletAddr []byte
				createdAt  time.Time
				updatedAt  time.Time
			)

			if err = rows.Scan(
				&id,
				&userID,
				&orgID,
				&walletAddr,
				&createdAt,
				&updatedAt,
			); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}

			employees = append(employees, &models.Employee{
				ID:             id,
				UserID:         userID,
				OrganizationId: orgID,
				WalletAddress:  walletAddr,
				CreatedAt:      createdAt,
				UpdatedAt:      updatedAt,
			})
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error execute transactional operation. %w", err)
	}

	return employees, nil
}
