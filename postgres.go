package ipam

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// import for sqlx to use postgres driver
	_ "github.com/lib/pq"
)

const postgresSchema = `
CREATE TABLE IF NOT EXISTS prefixes (
	cidr   text NOT NULL,
	tenantid text NOT NULL,
	prefix JSONB,
	PRIMARY KEY (cidr,tenantid)
);

CREATE INDEX IF NOT EXISTS prefix_idx ON prefixes USING GIN(prefix);
CREATE INDEX IF NOT EXISTS tenant_prefix_idx ON prefixes (tenantid);
CREATE INDEX IF NOT EXISTS tenant_prefix_idx ON prefixes (cidr);
`

// SSLMode specifies how to configure ssl encryption to the database
type SSLMode string

func (s SSLMode) String() string {
	return "sslmode=" + string(s)
}

const (
	// SSLModeAllow I don't care about security
	// but I will pay the overhead of encryption if the server insists on it
	SSLModeAllow = SSLMode("allow")
	// SSLModeDisable I don't care about security
	// and I don't want to pay the overhead of encryption.
	SSLModeDisable = SSLMode("disable")
	// SSLModePrefer I don't care about encryption
	// but I wish to pay the overhead of encryption if the server supports it.
	SSLModePrefer = SSLMode("prefer")
	// SSLModeRequire I want my data to be encrypted and I accept the overhead.
	// I trust that the network will make sure I always connect to the server I want.
	SSLModeRequire = SSLMode("require")
	// SSLModeVerifyCA I want my data encrypted and I accept the overhead.
	// I want to be sure that I connect to a server that I trust.
	SSLModeVerifyCA = SSLMode("verify-ca")
	// SSLModeVerifyFull I want my data encrypted and I accept the overhead.
	// I want to be sure that I connect to a server I trust, and that it's the one I specify.
	SSLModeVerifyFull = SSLMode("verify-full")
)

// NewPostgresStorage creates a new Storage which uses postgres.
func NewPostgresStorage(host, port, user, password, dbname string, sslmode SSLMode) (*sql, error) {
	dsn := dataSource(host, port, user, password, dbname, sslmode)
	return NewPostgresStorageFromDsn(dsn)
}

func NewPostgresStorageFromDsn(dsn string) (*sql, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database:%v", err)
	}
	db.MustExec(postgresSchema)
	return &sql{
		db: db,
	}, nil
}

func dataSource(host, port, user, password, dbname string, sslmode SSLMode) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", user, password, host, port, dbname, sslmode)
}
