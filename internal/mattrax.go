package mattrax

import (
	"github.com/gorilla/mux"
	"github.com/mattrax/Mattrax/internal/authentication"
	"github.com/mattrax/Mattrax/internal/certificates"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/settings"
	"github.com/patrickmn/go-cache"
)

// Version is contains the build information. It is injected at build time
const Version = "v0.0.0-dev Commit: xxxxxx"

// Server contains the global server state
type Server struct {
	Args         Arguments
	GlobalRouter *mux.Router
	Router       *mux.Router // Subrouter which is only accessible via secure origins (configured by admin)

	DB    *db.Queries
	Cache *cache.Cache

	Cert     *certificates.Service
	Auth     *authentication.Service
	Settings *settings.Service
}

// Arguments are the command line flags
type Arguments struct {
	Domain  string `placeholder:"\"mdm.example.com\"" help:"The domain your server is accessible from"`
	DB      string `placeholder:"\"postgres://localhost/Mattrax\"" help:"The Postgres database connection url"`
	Addr    string `default:":443" placeholder:"\":443\"" help:"The listen address of the https server"`
	TLSCert string `default:"./certs/tls.crt" placeholder:"\"./certs/tls.crt\"" help:"The path for the tls certificate"`
	TLSKey  string `default:"./certs/tls.key" placeholder:"\"./certs/tls.key\"" help:"The path for the tls certificates key"`

	Debug bool `help:"Enabled development mode. PLEASE DO NOT USE IN PRODUCTION!"`
}

// Description is for alexflint/go-args
func (Arguments) Description() string {
	return "Mattrax MDM Server. Created by Oscar Beaumont!"
}

// Version is for alexflint/go-args
func (Arguments) Version() string {
	return "Version: " + Version
}
