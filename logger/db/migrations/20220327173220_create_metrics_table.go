package migrations

import (
	"context"
	"database/sql"
)

type CreateMetricsTable struct {
	DB *sql.DB
}

func (m *CreateMetricsTable) Up() error {
	sql := `
  CREATE TABLE IF NOT EXISTS metrics (
    time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    device text NULL,
    channel text NULL,
    value double PRECISION NULL
  );
  `

	ctx := context.Background()
	_, err := m.DB.ExecContext(ctx, sql)
	if err != nil {
		return err
	}

	sql = `
  SELECT create_hypertable('metrics', 'time') WHERE NOT EXISTS (
    SELECT * FROM information_schema.tables WHERE table_name = 'metrics'
  );
  `
	_, err = m.DB.ExecContext(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

func (m *CreateMetricsTable) Down() error {
	sql := `DROP TABLE IF EXISTS metrics;`

	ctx := context.Background()
	_, err := m.DB.ExecContext(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}
