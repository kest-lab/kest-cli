package logger

// Predefined channel names for various system components
const (
	ChannelStack     = "stack"
	ChannelSingle    = "single"
	ChannelDaily     = "daily"
	ChannelSlack     = "slack"
	ChannelStderr    = "stderr"
	ChannelSyslog    = "syslog"
	ChannelErrorlog  = "errorlog"
	ChannelNull      = "null"
	ChannelEmergency = "emergency"
)

// Common application channels
const (
	ChannelApp      = "app"
	ChannelHTTP     = "http"
	ChannelDatabase = "database"
	ChannelQueue    = "queue"
	ChannelMail     = "mail"
	ChannelAuth     = "auth"
	ChannelAPI      = "api"
)

// App returns a logger for general application logs
func App() *Logger {
	return Channel(ChannelApp)
}

// HTTP returns a logger for HTTP request/response logs
func HTTP() *Logger {
	return Channel(ChannelHTTP)
}

// Database returns a logger for database query logs
func Database() *Logger {
	return Channel(ChannelDatabase)
}

// Queue returns a logger for queue job logs
func Queue() *Logger {
	return Channel(ChannelQueue)
}

// Mail returns a logger for email logs
func Mail() *Logger {
	return Channel(ChannelMail)
}

// Auth returns a logger for authentication logs
func Auth() *Logger {
	return Channel(ChannelAuth)
}

// API returns a logger for API logs
func API() *Logger {
	return Channel(ChannelAPI)
}
