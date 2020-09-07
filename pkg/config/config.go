package config

import (
	glauth "github.com/owncloud/ocis-glauth/pkg/config"
	graphExplorer "github.com/owncloud/ocis-graph-explorer/pkg/config"
	graph "github.com/owncloud/ocis-graph/pkg/config"
	hello "github.com/owncloud/ocis-hello/pkg/config"
	konnectd "github.com/owncloud/ocis-konnectd/pkg/config"
	ocs "github.com/owncloud/ocis-ocs/pkg/config"
	phoenix "github.com/owncloud/ocis-phoenix/pkg/config"
	proxy "github.com/owncloud/ocis-proxy/pkg/config"
	reva "github.com/owncloud/ocis-reva/pkg/config"
	settings "github.com/owncloud/ocis-settings/pkg/config"
	store "github.com/owncloud/ocis-store/pkg/config"
	thumbnails "github.com/owncloud/ocis-thumbnails/pkg/config"
	webdav "github.com/owncloud/ocis-webdav/pkg/config"
	accounts "github.com/refs/ocis-mono/ocis-accounts/pkg/config"
	pman "github.com/refs/pman/pkg/config"
)

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr string
	Root string
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Config combines all available configuration parts.
type Config struct {
	File    string
	Log     Log
	Debug   Debug
	HTTP    HTTP
	GRPC    GRPC
	Tracing Tracing

	Accounts      *accounts.Config
	Graph         *graph.Config
	GraphExplorer *graphExplorer.Config
	GLAuth        *glauth.Config
	Hello         *hello.Config
	Konnectd      *konnectd.Config
	OCS           *ocs.Config
	Phoenix       *phoenix.Config
	Proxy         *proxy.Config
	Reva          *reva.Config
	Thumbnails    *thumbnails.Config
	WebDAV        *webdav.Config
	Settings      *settings.Config
	Store         *store.Config
	Runtime       *pman.Config
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{
		Accounts:      accounts.New(),
		Graph:         graph.New(),
		GraphExplorer: graphExplorer.New(),
		Hello:         hello.New(),
		Konnectd:      konnectd.New(),
		OCS:           ocs.New(),
		Phoenix:       phoenix.New(),
		WebDAV:        webdav.New(),
		Reva:          reva.New(),
		GLAuth:        glauth.New(),
		Proxy:         proxy.New(),
		Thumbnails:    thumbnails.New(),
		Settings:      settings.New(),
		Store:         store.New(),
		Runtime:       pman.NewConfig(),
	}
}
