package infra

import (
	"context"
	"database/sql"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	smpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/francisco-alonso/stock-portfolio/portfolio-service/domain"
	_ "github.com/lib/pq"
)

type PortfolioRepository interface {
    AddPosition(userID, asset string, quantity int, price float64) error
    GetPositions(userID string) ([]domain.Position, error)
}

type PostgresDB struct {
    db *sql.DB
}

func NewPostgresDB(host, user, password, dbname string) (PortfolioRepository, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", 
        host, user, password, dbname)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("error connecting to database: %w", err)
    }

    return &PostgresDB{db: db}, nil
}

func GetSecret(ctx context.Context, secretName string) (string, error) {
    client, err := secretmanager.NewClient(ctx)
    if err != nil {
        return "", err
    }
    defer client.Close()

    req := &smpb.AccessSecretVersionRequest{
        Name: fmt.Sprintf("%s/versions/latest", secretName),
    }

    result, err := client.AccessSecretVersion(ctx, req)
    if err != nil {
        return "", err
    }

    return string(result.Payload.Data), nil
}

func (p *PostgresDB) AddPosition(userID, asset string, quantity int, price float64) error {
    query := `
        INSERT INTO positions (user_id, asset, quantity, price) 
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, asset) 
        DO UPDATE SET quantity = positions.quantity + $3, price = $4
    `
    _, err := p.db.Exec(query, userID, asset, quantity, price)
    return err
}

func (p *PostgresDB) GetPositions(userID string) ([]domain.Position, error) {
    query := "SELECT user_id, asset, quantity, price FROM positions WHERE user_id = $1"
    rows, err := p.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var positions []domain.Position
    for rows.Next() {
        var pos domain.Position
        if err := rows.Scan(&pos.UserID, &pos.Asset, &pos.Quantity, &pos.Price); err != nil {
            return nil, err
        }
        positions = append(positions, pos)
    }
    return positions, nil
}
