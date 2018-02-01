package services

type Config struct {
	DB       *DBConfig
	Sessions *SessionsConfig
}

type Provider struct {
	DB       *DB
	Sessions *Sessions
}

// NewProvider returns the ServiceProvider
func NewProvider(conf *Config) *Provider {
	var provider = &Provider{}

	provider.DB = NewDB(conf.DB)
	provider.Sessions = NewSessions(conf.Sessions)

	return provider
}
