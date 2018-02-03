package services

type Config struct {
	CollectionPath string
	Sessions       *SessionsConfig
}

type Provider struct {
	ClientCollection *ClientCollection
	Sessions         *Sessions
}

// NewProvider returns the ServiceProvider
func NewProvider(conf *Config) *Provider {
	var provider = &Provider{}

	provider.ClientCollection = NewClientCollection(conf.CollectionPath)
	provider.Sessions = NewSessions(conf.Sessions)

	return provider
}
