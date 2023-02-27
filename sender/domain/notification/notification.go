package notification

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	id uuid.UUID

	message string
	status  NotificationStatus

	sender      string
	recipientId uuid.UUID

	creationDate time.Time
}

func NewNotification(message, sender string, recipientId uuid.UUID) (*Notification, error) {
	if message == "" {
		return nil, errors.New("empty message")
	}
	if sender == "" {
		return nil, errors.New("empty sender")
	}
	return &Notification{
		id: uuid.New(),

		message: message,
		status:  CREATED,

		sender:      sender,
		recipientId: recipientId,

		creationDate: time.Now(),
	}, nil
}
func (n *Notification) ID() uuid.UUID {
	return n.id
}
func (n *Notification) Message() string {
	return n.message
}
func (n *Notification) CreationDate() time.Time {
	return n.creationDate
}
func (n *Notification) Sender() string {
	return n.sender
}
func (n *Notification) RecipientId() uuid.UUID {
	return n.recipientId
}
func (n *Notification) Status() NotificationStatus {
	return n.status
}

func (n *Notification) IsCreated() bool {
	return n.status == CREATED
}
func (n *Notification) IsSend() bool {
	return n.status == SEND
}
func (n *Notification) IsReceived() bool {
	return n.status == RECEIVED
}
func (n *Notification) IsSeen() bool {
	return n.status == SEEN
}
func (n *Notification) Send() error {
	if n.IsSend() {
		return errors.New("notification already sended")
	}
	n.status = SEND
	return nil
}
func (n *Notification) Receive() error {
	if n.IsReceived() {
		return errors.New("notification already received")
	}
	n.status = RECEIVED
	return nil
}
func (n *Notification) Seen() error {
	if n.IsSeen() {
		return errors.New("notification already received")
	}
	n.status = SEEN
	return nil
}
func (n *Notification) UpdateStatus(status NotificationStatus) {
	n.status = status
}

const DEFAULT_DATE_FORMAT = "2006-01-02 15:04:05.999999999 -0700 MST"

type MapperConfig struct {
	DateFormat string
}

func (mc MapperConfig) Validate() error {
	var err error
	if mc.DateFormat == "" {
		err = errors.New("invalid date format")
	}
	return err
}

type Mapper struct {
	mc MapperConfig
}

func (m Mapper) IsZero() bool {
	return m == Mapper{}
}
func (m Mapper) UnmarshalNotificationFromDatabase(
	id string,
	message string,
	creationDate string,
	sender string,
	recipient string,
	status string,
) (*Notification, error) {
	if message == "" {
		return nil, errors.New("empty message")
	}
	if sender == "" {
		return nil, errors.New("empty sender")
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	recipientId, err := uuid.Parse(recipient)
	if err != nil {
		return nil, err
	}
	notificationStatus, err := NewNotificationStatusFromString(status)
	if err != nil {
		return nil, err
	}
	date, err := m.parseTime(creationDate)
	if err != nil {
		return nil, err
	}
	return &Notification{
		id: uid,

		message: message,
		status:  notificationStatus,

		sender:      sender,
		recipientId: recipientId,

		creationDate: date,
	}, nil
}

func (m Mapper) parseTime(creationDate string) (time.Time, error) {
	if date, err := time.Parse(m.mc.DateFormat, creationDate); err != nil {
		return time.Time{}, err
	} else {
		return date, nil
	}
}
func (m Mapper) CreationDateToString(creationdate time.Time) string {
	return creationdate.Format(m.mc.DateFormat)
}
func NewMapper(mc MapperConfig) (Mapper, error) {
	if err := mc.Validate(); err != nil {
		return Mapper{}, err
	}
	return Mapper{mc}, nil
}
func NewDefaultMapper() Mapper {
	return Mapper{
		mc: MapperConfig{DEFAULT_DATE_FORMAT},
	}
}
