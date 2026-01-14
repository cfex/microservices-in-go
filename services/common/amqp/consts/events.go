package amqp

const (
	EmailWelcome = "email.welcome"
)

const (
	UserFollow = "user.follow"
)

const (
	ProjectCreated = "project.created"
	ProjectUpdated = "project.updated"
)

type NotificationEvent struct {
	Type any
	Body any
}
