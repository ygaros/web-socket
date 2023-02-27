package adapters

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sender/domain/notification"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteNotification struct {
	Id           string `db:"id"`
	Message      string `db:"message"`
	CreationDate string `db:"creation_date"`
	Sender       string `db:"sender"`
}
type sqliteClientNotification struct {
	ClientId       string `db:"client_id"`
	NotificationId string `db:"notification_id"`
	Status         string `db:"notification_status"`
}
type sqliteDb interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
type SqliteNotificationRepository struct {
	db *sql.DB
	m  notification.Mapper
}

func (s SqliteNotificationRepository) GetAllNotifications(
	ctx context.Context,
) ([]notification.Notification, error) {
	queryN := "SELECT * FROM 'notification'"
	queryCN := "SELECT * FROM 'client_notifications'"

	var sqlNs []sqliteNotification
	var sqlCNs []sqliteClientNotification

	rows, err := s.db.QueryContext(ctx, queryN)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		sqlN := sqliteNotification{}
		err := rows.Scan(&sqlN.Id, &sqlN.Message, &sqlN.CreationDate, &sqlN.Sender)
		if err != nil {
			continue
		}
		sqlNs = append(sqlNs, sqlN)
	}
	rows, err = s.db.QueryContext(ctx, queryCN)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		sqlCN := sqliteClientNotification{}
		err := rows.Scan(&sqlCN.ClientId, &sqlCN.NotificationId, &sqlCN.Status)
		if err != nil {
			continue
		}
		sqlCNs = append(sqlCNs, sqlCN)
	}
	var notifications []notification.Notification

	for _, sqlNot := range sqlNs {
		for _, sqlCNot := range sqlCNs {
			if sqlCNot.NotificationId == sqlNot.Id {
				if n, err := s.m.UnmarshalNotificationFromDatabase(
					sqlCNot.NotificationId,
					sqlNot.Message,
					sqlNot.CreationDate,
					sqlNot.Sender,
					sqlCNot.ClientId,
					sqlCNot.Status,
				); err == nil {
					notifications = append(notifications, *n)
				} else {
					log.Println(err)
				}
			}
		}
	}
	return notifications, nil
}

func (s SqliteNotificationRepository) GetAllNotificationByClientId(
	ctx context.Context,
	clientId uuid.UUID) ([]notification.Notification, error) {
	query := "SELECT * FROM 'client_notifications' WHERE client_id = ?"
	rows, err := s.db.QueryContext(ctx, query, clientId)
	if err != nil {
		return nil, err
	}
	var cns []sqliteClientNotification
	for rows.Next() {
		cn := sqliteClientNotification{}
		err := rows.Scan(&cn.ClientId, &cn.NotificationId, &cn.Status)
		if err == nil {
			cns = append(cns, cn)
		}
	}
	var notifications []notification.Notification
	for _, cn := range cns {
		uid, err := uuid.Parse(cn.NotificationId)
		if err != nil {
			continue
		}
		notification, err := s.GetNotification(ctx, uid, clientId)
		if err == nil {
			notifications = append(notifications, *notification)
		}
	}
	return notifications, nil
}

func (s SqliteNotificationRepository) GetNotification(
	ctx context.Context,
	notificationId uuid.UUID,
	clientId uuid.UUID,
) (*notification.Notification, error) {
	//TODO IMPLEMENT
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = s.finishTransaction(err, tx)
	}()
	query := "SELECT * FROM 'client_notifications' WHERE client_id = ? AND notification_id = ?"
	query2 := "SELECT * FROM 'notification' WHERE id = ?"
	row := tx.QueryRowContext(ctx, query, clientId.String(), notificationId.String())
	cn := sqliteClientNotification{}
	err = row.Scan(&cn.ClientId, &cn.NotificationId, &cn.Status)
	if err != nil {
		return nil, err
	}
	row = tx.QueryRowContext(ctx, query2, notificationId.String())
	n := sqliteNotification{}
	err = row.Scan(&n.Id, &n.Message, &n.CreationDate, &n.Sender)
	if err != nil {
		return nil, err
	}
	return s.m.UnmarshalNotificationFromDatabase(n.Id, n.Message, n.CreationDate, n.Sender, cn.ClientId, cn.Status)
}
func (s SqliteNotificationRepository) UpdateNotification(
	ctx context.Context,
	notificationId uuid.UUID,
	clientId uuid.UUID,
	updateFn func(n *notification.Notification) (*notification.Notification, error),
) (err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return errors.New("unable to start transaction")
	}
	defer func() {
		err = s.finishTransaction(err, tx)
	}()
	sqn, err := s.getSqliteNotification(ctx, notificationId, tx)
	sqnc, err := s.getSqliteClientNotification(ctx, notificationId, clientId, tx)
	marshalledN, err := s.m.UnmarshalNotificationFromDatabase(
		sqn.Id,
		sqn.Message,
		sqn.CreationDate,
		sqn.Sender,
		sqnc.ClientId,
		sqnc.Status,
	)
	updated, err := updateFn(marshalledN)
	if err != nil {
		return err
	}
	updatedSqn, updatedSqnc := MarshallToNotificationFromDatabase(updated)
	err = s.updateSqliteNotification(ctx, *updatedSqn, tx)
	if err != nil {
		return err
	}
	if sqnc.Status != updatedSqnc.Status {
		err = s.updateSqliteClientNotification(ctx, *updatedSqnc, tx)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s SqliteNotificationRepository) SaveNotification(ctx context.Context, notification *notification.Notification) (notf *notification.Notification, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		s.finishTransaction(err, tx)
	}()
	return s.saveNotification(ctx, notification, tx)
}
func (s SqliteNotificationRepository) saveNotification(ctx context.Context, notification *notification.Notification, db sqliteDb) (*notification.Notification, error) {
	query := "INSERT INTO 'notification' VALUES (?, ?, ?, ?)"
	_, err := db.ExecContext(ctx, query,
		notification.ID().String(),
		notification.Message(),
		s.m.CreationDateToString(notification.CreationDate()),
		notification.Sender(),
	)
	if err != nil {
		return nil, err
	}

	return notification, s.saveClientNotification(ctx, notification, db)

}
func (s SqliteNotificationRepository) saveClientNotification(ctx context.Context, notification *notification.Notification, db sqliteDb) error {
	query := "INSERT INTO 'client_notifications' VALUES (?, ?, ?)"
	_, err := db.ExecContext(ctx, query, notification.RecipientId().String(), notification.ID().String(), notification.Status().String())
	if err != nil {
		return err
	}
	return nil
}
func (s SqliteNotificationRepository) getSqliteNotification(
	ctx context.Context,
	notificationId uuid.UUID,
	db sqliteDb,
) (*sqliteNotification, error) {

	query := "SELECT * FROM 'notification' WHERE id = ?"
	row := db.QueryRowContext(ctx, query, notificationId)
	var sqn sqliteNotification
	err := row.Scan(&sqn.Id, &sqn.Message, &sqn.CreationDate, &sqn.Sender)
	if err != nil {
		return &sqn, err
	}
	return &sqn, nil
}
func (s SqliteNotificationRepository) getSqliteClientNotification(
	ctx context.Context,
	notificationId uuid.UUID,
	clientId uuid.UUID,
	db sqliteDb,
) (*sqliteClientNotification, error) {

	query := "SELECT * FROM 'client_notifications' WHERE client_id = ? AND notification_id = ?"
	row := db.QueryRowContext(ctx, query, clientId, notificationId)
	var sqnc sqliteClientNotification
	err := row.Scan(&sqnc.ClientId, &sqnc.NotificationId, &sqnc.Status)
	if err != nil {
		return &sqnc, err
	}
	return &sqnc, nil
}
func (s SqliteNotificationRepository) updateSqliteNotification(
	ctx context.Context,
	sqn sqliteNotification,
	db sqliteDb,
) error {

	query := "UPDATE 'notification' SET message = ?, creation_date = ?, sender = ? WHERE id = ?"
	_, err := db.ExecContext(ctx, query, sqn.Message, sqn.CreationDate, sqn.Sender, sqn.Id)
	if err != nil {
		return err
	}
	return nil
}
func (s SqliteNotificationRepository) updateSqliteClientNotification(
	ctx context.Context,
	sqnc sqliteClientNotification,
	db sqliteDb,
) error {

	query := "UPDATE 'client_notifications' SET notification_status = ? WHERE client_id = ? AND notification_id = ?"
	_, err := db.ExecContext(ctx, query, sqnc.Status, sqnc.ClientId, sqnc.NotificationId)
	if err != nil {
		return err
	}
	return nil
}
func (s SqliteNotificationRepository) finishTransaction(err error, tx *sql.Tx) error {
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
func MarshallToNotificationFromDatabase(n *notification.Notification) (*sqliteNotification, *sqliteClientNotification) {
	return &sqliteNotification{
			Id:           n.ID().String(),
			Message:      n.Message(),
			CreationDate: n.CreationDate().String(),
			Sender:       n.Sender(),
		}, &sqliteClientNotification{
			ClientId:       n.RecipientId().String(),
			NotificationId: n.ID().String(),
			Status:         n.Status().String(),
		}
}
func NewSqliteNotificationRepostory(db *sql.DB, m notification.Mapper) *SqliteNotificationRepository {
	if db == nil {
		panic("missing db")
	}
	if m.IsZero() {
		panic("missing mapper")
	}
	return &SqliteNotificationRepository{
		db: db,
		m:  m,
	}
}
