package contracts

import "context"

// Mailer defines the email sending interface
type Mailer interface {
	Send(ctx context.Context, mail Mail) error
	Queue(ctx context.Context, mail Mail) error
}

// Mail represents an email message
type Mail interface {
	To() []string
	Subject() string
	Body() string
	IsHTML() bool
}

// Mailable is implemented by types that can be sent as email
type Mailable interface {
	Build() Mail
}
