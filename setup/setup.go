package setup

import (
	"os"

	"github.com/Sirupsen/logrus"

	que "github.com/bgentry/que-go"
	"github.com/jackc/pgx"
)

//SetUp A Global hook for setup
var (
	SetUp setUp
)

//SetUp ...
type setUp struct {
	dbURL string
}

func init() {
	logrus.Info("In Setup init")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logrus.Fatal("Unable to get DATABASE_URL url")
	}
	SetUp = setUp{dbURL: dbURL}
}

// GetPgxPool based on the provided database URL
func (s *setUp) GetPgxPool() (*pgx.ConnPool, error) {
	pgxcfg, err := pgx.ParseURI(s.dbURL)
	if err != nil {
		return nil, err
	}

	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   pgxcfg,
		AfterConnect: que.PrepareStatements,
	})

	if err != nil {
		return nil, err
	}

	return pgxpool, nil
}

//PoolAndQueueConnection ...
func (s *setUp) PoolAndQueueConnection() (*pgx.ConnPool, *que.Client, error) {
	pgxpool, err := s.GetPgxPool()
	if err != nil {
		return nil, nil, err
	}

	qc := que.NewClient(pgxpool)

	return pgxpool, qc, err
}
