package notification

import (
	"context"
	"sync"
)

// Notifiable represents an entity that can receive notifications
type Notifiable interface {
	// RouteNotificationFor returns the routing information for a channel
	RouteNotificationFor(channel string) string

	// GetID returns the notifiable's ID
	GetID() interface{}
}

// Notification represents a notification
type Notification interface {
	// Via returns the channels to send the notification through
	Via(notifiable Notifiable) []string
}

// MailNotification can be sent via mail
type MailNotification interface {
	Notification
	ToMail(notifiable Notifiable) *MailMessage
}

// SMSNotification can be sent via SMS
type SMSNotification interface {
	Notification
	ToSMS(notifiable Notifiable) *SMSMessage
}

// DatabaseNotification can be stored in database
type DatabaseNotification interface {
	Notification
	ToDatabase(notifiable Notifiable) map[string]interface{}
}

// SlackNotification can be sent to Slack
type SlackNotification interface {
	Notification
	ToSlack(notifiable Notifiable) *SlackMessage
}

// Channel defines a notification channel
type Channel interface {
	// Send sends the notification through this channel
	Send(ctx context.Context, notifiable Notifiable, notification Notification) error
}

// --- Mail Message ---

// MailMessage represents an email notification
type MailMessage struct {
	Subject     string
	From        string
	ReplyTo     string
	To          []string
	CC          []string
	BCC         []string
	Body        string
	HTMLBody    string
	Attachments []Attachment
	Headers     map[string]string
}

// Attachment represents an email attachment
type Attachment struct {
	Name    string
	Content []byte
	Type    string
}

// NewMailMessage creates a new mail message
func NewMailMessage() *MailMessage {
	return &MailMessage{
		Headers: make(map[string]string),
	}
}

// SetSubject sets the email subject
func (m *MailMessage) SetSubject(subject string) *MailMessage {
	m.Subject = subject
	return m
}

// SetFrom sets the sender
func (m *MailMessage) SetFrom(from string) *MailMessage {
	m.From = from
	return m
}

// SetTo sets the recipients
func (m *MailMessage) SetTo(to ...string) *MailMessage {
	m.To = to
	return m
}

// AddCC adds CC recipients
func (m *MailMessage) AddCC(cc ...string) *MailMessage {
	m.CC = append(m.CC, cc...)
	return m
}

// AddBCC adds BCC recipients
func (m *MailMessage) AddBCC(bcc ...string) *MailMessage {
	m.BCC = append(m.BCC, bcc...)
	return m
}

// SetBody sets the plain text body
func (m *MailMessage) SetBody(body string) *MailMessage {
	m.Body = body
	return m
}

// SetHTMLBody sets the HTML body
func (m *MailMessage) SetHTMLBody(html string) *MailMessage {
	m.HTMLBody = html
	return m
}

// Attach adds an attachment
func (m *MailMessage) Attach(name string, content []byte, contentType string) *MailMessage {
	m.Attachments = append(m.Attachments, Attachment{
		Name:    name,
		Content: content,
		Type:    contentType,
	})
	return m
}

// SetHeader sets a header
func (m *MailMessage) SetHeader(key, value string) *MailMessage {
	m.Headers[key] = value
	return m
}

// --- SMS Message ---

// SMSMessage represents an SMS notification
type SMSMessage struct {
	To      string
	From    string
	Content string
}

// NewSMSMessage creates a new SMS message
func NewSMSMessage(content string) *SMSMessage {
	return &SMSMessage{Content: content}
}

// SetTo sets the recipient phone number
func (m *SMSMessage) SetTo(to string) *SMSMessage {
	m.To = to
	return m
}

// SetFrom sets the sender ID
func (m *SMSMessage) SetFrom(from string) *SMSMessage {
	m.From = from
	return m
}

// --- Slack Message ---

// SlackMessage represents a Slack notification
type SlackMessage struct {
	Channel     string
	Text        string
	Username    string
	IconEmoji   string
	Attachments []SlackAttachment
}

// SlackAttachment represents a Slack message attachment
type SlackAttachment struct {
	Color      string
	Title      string
	TitleLink  string
	Text       string
	Fields     []SlackField
	Footer     string
	FooterIcon string
	Timestamp  int64
}

// SlackField represents a field in Slack attachment
type SlackField struct {
	Title string
	Value string
	Short bool
}

// NewSlackMessage creates a new Slack message
func NewSlackMessage(text string) *SlackMessage {
	return &SlackMessage{Text: text}
}

// SetChannel sets the Slack channel
func (m *SlackMessage) SetChannel(channel string) *SlackMessage {
	m.Channel = channel
	return m
}

// SetUsername sets the bot username
func (m *SlackMessage) SetUsername(username string) *SlackMessage {
	m.Username = username
	return m
}

// SetIcon sets the icon emoji
func (m *SlackMessage) SetIcon(emoji string) *SlackMessage {
	m.IconEmoji = emoji
	return m
}

// AddAttachment adds an attachment
func (m *SlackMessage) AddAttachment(attachment SlackAttachment) *SlackMessage {
	m.Attachments = append(m.Attachments, attachment)
	return m
}

// --- Notification Manager ---

// Manager manages notification channels and sending
type Manager struct {
	mu       sync.RWMutex
	channels map[string]Channel
}

var (
	manager *Manager
	once    sync.Once
)

// Global returns the global notification manager
func Global() *Manager {
	once.Do(func() {
		manager = New()
	})
	return manager
}

// New creates a new notification manager
func New() *Manager {
	return &Manager{
		channels: make(map[string]Channel),
	}
}

// RegisterChannel registers a notification channel
func (m *Manager) RegisterChannel(name string, channel Channel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.channels[name] = channel
}

// Channel returns a channel by name
func (m *Manager) Channel(name string) Channel {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.channels[name]
}

// Send sends a notification to a notifiable
func (m *Manager) Send(ctx context.Context, notifiable Notifiable, notification Notification) error {
	channels := notification.Via(notifiable)

	for _, channelName := range channels {
		channel := m.Channel(channelName)
		if channel == nil {
			continue
		}

		if err := channel.Send(ctx, notifiable, notification); err != nil {
			return err
		}
	}

	return nil
}

// SendNow sends a notification immediately (alias for Send)
func (m *Manager) SendNow(ctx context.Context, notifiable Notifiable, notification Notification) error {
	return m.Send(ctx, notifiable, notification)
}

// --- Convenience Functions ---

// RegisterChannel registers a channel with the global manager
func RegisterChannel(name string, channel Channel) {
	Global().RegisterChannel(name, channel)
}

// Send sends a notification using the global manager
func Send(ctx context.Context, notifiable Notifiable, notification Notification) error {
	return Global().Send(ctx, notifiable, notification)
}

// SendNow sends a notification immediately using the global manager
func SendNow(ctx context.Context, notifiable Notifiable, notification Notification) error {
	return Global().SendNow(ctx, notifiable, notification)
}

// --- Log Channel (for testing/development) ---

// LogChannel logs notifications instead of sending them
type LogChannel struct {
	mu       sync.Mutex
	Messages []LoggedNotification
}

// LoggedNotification represents a logged notification
type LoggedNotification struct {
	Notifiable   Notifiable
	Notification Notification
	Channel      string
}

// NewLogChannel creates a new log channel
func NewLogChannel() *LogChannel {
	return &LogChannel{
		Messages: make([]LoggedNotification, 0),
	}
}

// Send logs the notification
func (c *LogChannel) Send(ctx context.Context, notifiable Notifiable, notification Notification) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Messages = append(c.Messages, LoggedNotification{
		Notifiable:   notifiable,
		Notification: notification,
		Channel:      "log",
	})

	return nil
}

// Clear clears all logged messages
func (c *LogChannel) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Messages = make([]LoggedNotification, 0)
}

// Count returns the number of logged messages
func (c *LogChannel) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.Messages)
}

// --- Simple Notifiable Implementation ---

// SimpleNotifiable is a simple implementation of Notifiable
type SimpleNotifiable struct {
	ID    interface{}
	Email string
	Phone string
	Slack string
}

// RouteNotificationFor returns the routing information
func (n *SimpleNotifiable) RouteNotificationFor(channel string) string {
	switch channel {
	case "mail":
		return n.Email
	case "sms":
		return n.Phone
	case "slack":
		return n.Slack
	default:
		return ""
	}
}

// GetID returns the ID
func (n *SimpleNotifiable) GetID() interface{} {
	return n.ID
}

// --- Anonymous Notification ---

// AnonymousNotification allows creating notifications inline
type AnonymousNotification struct {
	channels []string
	mail     func(Notifiable) *MailMessage
	sms      func(Notifiable) *SMSMessage
	slack    func(Notifiable) *SlackMessage
	database func(Notifiable) map[string]interface{}
}

// NewAnonymousNotification creates a new anonymous notification
func NewAnonymousNotification(channels ...string) *AnonymousNotification {
	return &AnonymousNotification{
		channels: channels,
	}
}

// Via returns the channels
func (n *AnonymousNotification) Via(notifiable Notifiable) []string {
	return n.channels
}

// WithMail sets the mail builder
func (n *AnonymousNotification) WithMail(fn func(Notifiable) *MailMessage) *AnonymousNotification {
	n.mail = fn
	return n
}

// WithSMS sets the SMS builder
func (n *AnonymousNotification) WithSMS(fn func(Notifiable) *SMSMessage) *AnonymousNotification {
	n.sms = fn
	return n
}

// WithSlack sets the Slack builder
func (n *AnonymousNotification) WithSlack(fn func(Notifiable) *SlackMessage) *AnonymousNotification {
	n.slack = fn
	return n
}

// WithDatabase sets the database builder
func (n *AnonymousNotification) WithDatabase(fn func(Notifiable) map[string]interface{}) *AnonymousNotification {
	n.database = fn
	return n
}

// ToMail returns the mail message
func (n *AnonymousNotification) ToMail(notifiable Notifiable) *MailMessage {
	if n.mail != nil {
		return n.mail(notifiable)
	}
	return nil
}

// ToSMS returns the SMS message
func (n *AnonymousNotification) ToSMS(notifiable Notifiable) *SMSMessage {
	if n.sms != nil {
		return n.sms(notifiable)
	}
	return nil
}

// ToSlack returns the Slack message
func (n *AnonymousNotification) ToSlack(notifiable Notifiable) *SlackMessage {
	if n.slack != nil {
		return n.slack(notifiable)
	}
	return nil
}

// ToDatabase returns the database data
func (n *AnonymousNotification) ToDatabase(notifiable Notifiable) map[string]interface{} {
	if n.database != nil {
		return n.database(notifiable)
	}
	return nil
}
