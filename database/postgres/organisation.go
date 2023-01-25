package postgres

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/frain-dev/convoy/datastore"
	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

var (
	ErrOrganizationNotCreated = errors.New("organization could not be created")
	ErrOrganizationNotUpdated = errors.New("organization could not be updated")
	ErrOrganizationNotDeleted = errors.New("organization could not be deleted")
)

type orgRepo struct {
	db *sqlx.DB
}

func NewOrgRepo(db *sqlx.DB) datastore.OrganisationRepository {
	return &orgRepo{db: db}
}

func (o *orgRepo) CreateOrganisation(ctx context.Context, org *datastore.Organisation) error {
	q := sq.Insert("convoy.organisations").
		Columns("name", "owner_id").
		Values(org.Name, org.OwnerID).
		PlaceholderFormat(sq.Dollar)

	sql, vals, err := q.ToSql()
	if err != nil {
		return err
	}

	result, err := o.db.ExecContext(ctx, sql, vals...)
	if err != nil {
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows < 1 {
		return ErrOrganizationNotCreated
	}

	return nil
}

func (o *orgRepo) LoadOrganisationsPaged(ctx context.Context, pageable datastore.Pageable) ([]datastore.Organisation, datastore.PaginationData, error) {
	skip := (pageable.Page - 1) * pageable.PerPage

	// TODO(raymond,daniel) implement cursor based pagination
	q := sq.Select("*").
		From("convoy.organisations").
		Where("deleted_at IS NULL").
		OrderBy("id").
		Limit(uint64(pageable.PerPage)).
		Offset(uint64(skip)).
		PlaceholderFormat(sq.Dollar)

	sql, vals, err := q.ToSql()
	if err != nil {
		return nil, datastore.PaginationData{}, err
	}

	rows, err := o.db.QueryxContext(ctx, sql, vals...)
	if err != nil {
		return nil, datastore.PaginationData{}, err
	}

	var organizations []datastore.Organisation
	for rows.Next() {
		var org datastore.Organisation

		err = rows.StructScan(&org)
		if err != nil {
			return nil, datastore.PaginationData{}, err
		}

		organizations = append(organizations, org)
	}

	var count int
	err = o.db.GetContext(ctx, &count, "SELECT COUNT(id) FROM convoy.organisations WHERE deleted_at IS NULL")
	if err != nil {
		return nil, datastore.PaginationData{}, err
	}

	pagination := datastore.PaginationData{
		Total:     int64(count),
		Page:      int64(pageable.Page),
		PerPage:   int64(pageable.PerPage),
		Prev:      int64(getPrevPage(pageable.Page)),
		Next:      int64(pageable.Page + 1),
		TotalPage: int64(math.Ceil(float64(count) / float64(pageable.PerPage))),
	}

	return organizations, pagination, nil
}

func (o *orgRepo) UpdateOrganisation(ctx context.Context, org *datastore.Organisation) error {
	q := sq.Update("convoy.organisations").
		Set("name", org.Name).
		Set("owner_id", org.OwnerID).
		Set("custom_domain", org.CustomDomain).
		Set("assigned_domain", org.AssignedDomain).
		Set("updated_at", time.Now()).
		Where("id = ?", org.UID).
		Where("deleted_at IS NULL").
		PlaceholderFormat(sq.Dollar)

	sql, vals, err := q.ToSql()
	if err != nil {
		return err
	}

	result, err := o.db.ExecContext(ctx, sql, vals...)
	if err != nil {
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows < 1 {
		return ErrOrganizationNotUpdated
	}

	return nil
}

func (o *orgRepo) DeleteOrganisation(ctx context.Context, uid string) error {
	q := sq.Update("convoy.organisations").
		Set("deleted_at", time.Now()).
		Where("deleted_at IS NULL").
		Where("id = ?", uid).
		PlaceholderFormat(sq.Dollar)

	sql, vals, err := q.ToSql()
	if err != nil {
		return err
	}

	result, err := o.db.ExecContext(ctx, sql, vals...)
	if err != nil {
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if nRows < 1 {
		return ErrOrganizationNotDeleted
	}

	return nil
}

func (o *orgRepo) FetchOrganisationByID(ctx context.Context, id string) (*datastore.Organisation, error) {
	q := sq.Select("*").
		From("convoy.organisations").
		Where("deleted_at IS NULL").
		Where("id = ?", id).
		PlaceholderFormat(sq.Dollar)

	sql, vals, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var org *datastore.Organisation
	err = o.db.QueryRowxContext(ctx, sql, vals...).StructScan(&org)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, datastore.ErrOrgNotFound
	}

	return org, nil
}

func (o *orgRepo) FetchOrganisationByAssignedDomain(ctx context.Context, domain string) (*datastore.Organisation, error) {
	q := sq.Select("*").
		From("convoy.organisations").
		Where("deleted_at IS NULL").
		Where("assigned_domain = ?", domain).
		PlaceholderFormat(sq.Dollar)

	sql, vals, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var org *datastore.Organisation
	err = o.db.QueryRowxContext(ctx, sql, vals...).StructScan(&org)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, datastore.ErrOrgNotFound
	}

	return org, nil
}

func (o *orgRepo) FetchOrganisationByCustomDomain(ctx context.Context, domain string) (*datastore.Organisation, error) {
	q := sq.Select("*").
		From("convoy.organisations").
		Where("deleted_at IS NULL").
		Where("custom_domain = ?", domain).
		PlaceholderFormat(sq.Dollar)

	sql, vals, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var org *datastore.Organisation
	err = o.db.QueryRowxContext(ctx, sql, vals...).StructScan(&org)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, datastore.ErrOrgNotFound
	}

	return org, nil
}

// getPrevPage returns calculated value for the prev page
func getPrevPage(page int) int {
	if page == 0 {
		return 1
	}

	prev := 0
	if page-1 <= 0 {
		prev = page
	} else {
		prev = page - 1
	}

	return prev
}
