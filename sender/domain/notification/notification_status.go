package notification

import "fmt"

var (
	CREATED   = NotificationStatus{"created"}
	SEND      = NotificationStatus{"send"}
	RECEIVED  = NotificationStatus{"received"}
	SEEN      = NotificationStatus{"seen"}
	BROADCAST = NotificationStatus{"broadcast"}
)

var notificationStatusValues = []NotificationStatus{
	CREATED,
	SEND,
	RECEIVED,
	SEEN,
	BROADCAST,
}

type NotificationStatus struct {
	s string
}

func (s NotificationStatus) String() string {
	return s.s
}
func (s NotificationStatus) IsZero() bool {
	return s == NotificationStatus{}
}
func NewNotificationStatusFromString(status string) (NotificationStatus, error) {
	for _, notificationStatus := range notificationStatusValues {
		if notificationStatus.String() == status {
			return notificationStatus, nil
		}
	}
	return NotificationStatus{}, fmt.Errorf("unknown '%s' notificationStatus", status)
}
