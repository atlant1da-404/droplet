package database

import (
	"context"
	"fmt"
	"time"

	"github.com/a631807682/zerofield"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgreSQLConfig - represents PostgreSQL service config.
type PostgreSQLConfig struct {
	User     string
	Password string
	Host     string
	Database string
}

// PostgreSQL - represents postgresql service.
type PostgreSQL struct {
	DB *gorm.DB
}

var _ Database = (*PostgreSQL)(nil)

// PostgreSQLModel provides base fields for database models (like gorm.PostgreSQLModel).
type PostgreSQLModel struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" swaggerignore:"true"`
}

// Option - represents PostgreSQL service option.
type Option func(*PostgreSQL)

// NewPostgreSQL - creates new instance of PostgreSQL service.
func NewPostgreSQL(cfg PostgreSQLConfig, opts ...Option) (*PostgreSQL, error) {
	// create instance of mysql
	sql := &PostgreSQL{}

	// apply custom options
	for _, opt := range opts {
		opt(sql)
	}

	// connect to database
	var err error
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s",
		cfg.User, cfg.Password, cfg.Database, cfg.Host,
	)
	sql.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
	}

	// https://github.com/a631807682/zerofield
	// allow update zero value field
	err = sql.DB.Use(zerofield.NewPlugin())
	if err != nil {
		return nil, fmt.Errorf("failed to use zerofield plugin: %w", err)
	}

	// create UUID extension.
	err = sql.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	return sql, nil
}

func (p *PostgreSQL) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	pgSql, err := p.DB.WithContext(ctx).DB()
	if err != nil {
		return fmt.Errorf("failed to connect")
	}
	err = pgSql.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect")
	}
	return nil
}

// Close - closes mysql service database connection.
func (p *PostgreSQL) Close() error {
	if p.DB != nil {
		db, _ := p.DB.DB()
		err := db.Close()
		if err != nil {
			return fmt.Errorf("failed to close postgresql connection: %w", err)
		}
	}
	return nil
}
