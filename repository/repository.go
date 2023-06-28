package repository

import (
	"BIMSupportBot/models"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repository interface {
	GetById(ctx context.Context, input string) (*models.ResponseModel, error)
}

type pgRepository struct {
	db *sqlx.DB
}

func NewPgRepository(db *sqlx.DB) Repository {
	return &pgRepository{db: db}
}

func (p *pgRepository) GetById(ctx context.Context, input string) (*models.ResponseModel, error) {
	var a = &models.ResponseModel{}
	if err := p.db.GetContext(
		ctx, a,
		findById,
		input,
	); err != nil {
		return nil, errors.Wrap(err, "pgRepository.GetById.GetContext")
	}
	return a, nil
}
