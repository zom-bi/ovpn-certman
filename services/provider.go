package services

type Config struct {
	DB       *DBConfig
	Sessions *SessionsConfig
	Email    *EmailConfig
}

type Provider struct {
	DB       *DB
	Sessions *Sessions
	Email    *Email
}

// NewProvider returns the ServiceProvider
func NewProvider(conf *Config) *Provider {
	var provider = &Provider{}

	provider.DB = NewDB(conf.DB)
	provider.Sessions = NewSessions(conf.Sessions)
	provider.Email = NewEmail(conf.Email)

	return provider
}
