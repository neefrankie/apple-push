package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/neefrankie/apple-push/pkg/message"
	"go.uber.org/zap"
)

type Repo struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func New(db *sqlx.DB, logger *zap.Logger) Repo {
	return Repo{
		db:     db,
		logger: logger,
	}
}

func (r Repo) LoadDevices(devices chan<- message.Device) error {
	defer close(devices)

	rows, err := r.db.Queryx(stmtAllDevice)
	if err != nil {
		return err
	}

	d := message.Device{}

	for rows.Next() {
		err := rows.StructScan(&d)
		if err != nil {
			continue
		}

		devices <- d
	}

	return nil
}
