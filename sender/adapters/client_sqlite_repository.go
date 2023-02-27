package adapters

import (
	"context"
	"database/sql"
	"errors"
	"sender/domain/client"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteClient struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}
type SqliteClientRepository struct {
	db *sql.DB
}

func (s SqliteClientRepository) GetAllClients(ctx context.Context) ([]client.Client, error) {
	query := "SELECT * FROM 'client'"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var clients []client.Client
	for rows.Next() {
		clt := sqliteClient{}
		err = rows.Scan(&clt.Id, &clt.Name)
		if err != nil {
			continue
		}
		uid, err := uuid.Parse(clt.Id)
		if err != nil {
			continue
		}
		c, err := client.NewClientWithId(uid, clt.Name)
		if err != nil {
			continue
		}

		clients = append(clients, *c)
	}
	return clients, nil
}
func (s SqliteClientRepository) GetClient(ctx context.Context, clientId uuid.UUID) (*client.Client, error) {
	return s.getOrCreateClient(ctx, clientId.String(), s.db)
}
func (s SqliteClientRepository) GetClientByName(ctx context.Context, clientName string) (*client.Client, error) {
	return s.getOrCreateClient(ctx, clientName, s.db)
}
func (s SqliteClientRepository) getOrCreateClient(
	ctx context.Context,
	clt string,
	db sqliteDb,
) (*client.Client, error) {
	query := "SELECT * FROM 'client' WHERE id = ? OR name = ?"
	row := db.QueryRowContext(ctx, query, clt, clt)
	var sClient sqliteClient
	err := row.Scan(&sClient.Id, &sClient.Name)
	if errors.Is(err, sql.ErrNoRows) {
		//If client doesnt exists it means that clt is clientName
		sC, err := s.createClient(ctx, clt, db)
		if err != nil {
			return nil, err
		}
		sClient = *sC
	} else if err != nil {
		return nil, err
	}
	//TODO REFACTOR TO CLIENT MAPPER REMOVE UUID PARSING INSIDE REPOSITORY
	uid, err := uuid.Parse(sClient.Id)
	if err != nil {
		return nil, err
	}
	c, err := client.NewClientWithId(uid, sClient.Name)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (s SqliteClientRepository) createClient(
	ctx context.Context,
	clientName string,
	db sqliteDb,
) (*sqliteClient, error) {
	query := "INSERT INTO 'client' VALUES (?, ?)"
	newUuid := uuid.New().String()
	_, err := db.ExecContext(ctx, query, newUuid, clientName)
	if err != nil {
		return nil, err
	}
	return &sqliteClient{
		Id:   newUuid,
		Name: clientName,
	}, nil
}
func (s SqliteClientRepository) UpdateClient(
	ctx context.Context,
	clientId uuid.UUID,
	updateFn func(c *client.Client) (*client.Client, error),
) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return errors.New("unable to start transaction")
	}
	defer func() {
		err = s.finishTransaction(err, tx)
	}()
	clt, err := s.getOrCreateClient(ctx, clientId.String(), tx)
	if err != nil {
		return err
	}
	updatedClient, err := updateFn(clt)
	if err != nil {
		return err
	}
	return s.updateClient(ctx, updatedClient, tx)
}
func (s SqliteClientRepository) updateClient(
	ctx context.Context,
	clt *client.Client,
	db sqliteDb,
) error {
	query := "UPDATE 'client' SET name = ? WHERE id = ?"
	_, err := db.ExecContext(ctx, query, clt.ID(), clt.Name())
	return err
}
func (s SqliteClientRepository) SaveClient(ctx context.Context, client *client.Client) (*client.Client, error) {
	return s.saveClient(ctx, client, s.db)
}
func (s SqliteClientRepository) saveClient(ctx context.Context, client *client.Client, db sqliteDb) (*client.Client, error) {
	query := "INSERT INTO 'client' VALUES (?, ?)"
	_, err := db.ExecContext(ctx, query, client.ID(), client.Name())
	return client, err
}
func (s SqliteClientRepository) finishTransaction(err error, tx *sql.Tx) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return commitErr
		}
		return nil
	}
}
func NewSqliteClientRepository(db *sql.DB) *SqliteClientRepository {
	if db == nil {
		panic("missing db")
	}
	return &SqliteClientRepository{
		db: db,
	}
}
