package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

// APP env: should only put environment variable related to the service itself, e.g.: Application Name, version, running environment, ...
var (
	EnvFlag = &cli.StringFlag{
		Name:    "env",
		Usage:   "Application environment: development, staging, production",
		EnvVars: []string{"DD_ENV"},
		Value:   os.Getenv("DD_ENV"),
	}

	AppNameFlag = &cli.StringFlag{
		Name:    "app_name",
		Usage:   "Application name",
		EnvVars: []string{"DD_SERVICE"},
		Value:   "Education_backend",
	}

	AppVersionFlag = &cli.StringFlag{
		Name:    "app_version",
		Usage:   "Application version",
		EnvVars: []string{"DD_VERSION"},
		Value:   "v1",
	}

	HTTPPortFlag = &cli.StringFlag{
		Name:    "http_port",
		Usage:   "Port binding to application",
		EnvVars: []string{"HTTP_PORT"},
		Value:   os.Getenv("HTTP_PORT"),
	}
	ProxyURLFlag = &cli.StringFlag{
		Name:    "proxy_url",
		Usage:   "Proxy URL",
		EnvVars: []string{"PROXY_URL"},
		Value:   "",
	}
)

// MYSQL env
var (
	MYSQLConnFlag = &cli.StringFlag{
		Name:    "mysql_conn",
		Usage:   `specify MySQL connection string. If non-empty then other flags begin with "mysql_" will be ignore`,
		EnvVars: []string{"JAWSDB_URL"},
		Value:   "",
	}

	// MYSQLHostFlag used for migration, not for subscribe
	MYSQLHostFlag = &cli.StringFlag{
		Name:    "mysql_host",
		Usage:   "specify MySQL master host",
		EnvVars: []string{"MYSQL_HOST"},
	}

	// MYSQLPortFlag used for migration, not for subscribe
	MYSQLPortFlag = &cli.StringFlag{
		Name:    "mysql_port",
		Usage:   "specify MySQL port",
		EnvVars: []string{"MYSQL_PORT"},
	}

	MYSQLMasterHostsFlag = &cli.StringFlag{
		Name:    "mysql_master_hosts",
		Usage:   "specify MySQL master hosts with port",
		EnvVars: []string{"MYSQL_MASTER_HOSTS"},
	}

	MySQLSlaveHostsFlag = &cli.StringFlag{
		Name:    "mysql_slave_hosts",
		Usage:   "specify MySQL slave hosts with port",
		EnvVars: []string{"MYSQL_SLAVE_HOSTS"},
	}

	MySQLUserFlag = &cli.StringFlag{
		Name:    "mysql_user",
		Usage:   "specify MySQL user",
		EnvVars: []string{"MYSQL_USER"},
	}

	MySQLPasswordFlag = &cli.StringFlag{
		Name:    "mysql_password",
		Usage:   "password used for MySQL user",
		EnvVars: []string{"MYSQL_PASSWORD"},
	}

	MySQLDatabaseFlag = &cli.StringFlag{
		Name:    "mysql_db",
		Usage:   "MySQL database is using by subscribe",
		EnvVars: []string{"MYSQL_DB"},
	}

	MySQLMaxOpenConnsFlag = &cli.IntFlag{
		Name:    "mysql_max_open_conns",
		Usage:   "sets the maximum number of open connections to the database",
		EnvVars: []string{"MYSQL_MAX_OPEN_CONNS"},
		Value:   50,
	}

	MySQLMaxIdleConnsFlag = &cli.IntFlag{
		Name:    "mysql_max_idle_conns",
		Usage:   "sets the maximum number of connections in the idle connection pool",
		EnvVars: []string{"MYSQL_MAX_IDLE_CONNS"},
		Value:   5,
	}

	MySQLConnMaxLifetimeFlag = &cli.IntFlag{
		Name:    "mysql_conn_max_lifetime",
		Usage:   "sets the maximum amount of time in minutes a connection may be reused",
		EnvVars: []string{"MYSQL_CONN_MAX_LIFETIME"},
		Value:   60,
	}
)

// For Redis
var (
	RedisConnFlag = &cli.StringFlag{
		Name:    "redis_conn",
		Usage:   `specify Redis connection string. If non empty then other flags begin with "redis_" will be ignore`,
		EnvVars: []string{"REDIS_CONN"}, // support for heroku deployment
		Value:   "",
	}

	RedisHostFlag = &cli.StringFlag{
		Name:    "redis_host",
		Usage:   "specify Redis host",
		EnvVars: []string{"REDIS_HOST"},
	}

	RedisPortFlag = &cli.StringFlag{
		Name:    "redis_port",
		Usage:   "Redis port is using by application",
		EnvVars: []string{"REDIS_PORT"},
	}

	RedisUserFlag = &cli.StringFlag{
		Name:    "redis_user",
		Usage:   "specify Redis user",
		EnvVars: []string{"REDIS_USER"},
		Value:   "default",
	}

	RedisPasswordFlag = &cli.StringFlag{
		Name:    "redis_password",
		Usage:   "password used for Redis user",
		EnvVars: []string{"REDIS_PASSWORD"},
	}

	RedisEnabledTLSFlag = &cli.BoolFlag{
		Name:    "redis_enabled_tls",
		Usage:   "enable tls for Redis tls connection",
		EnvVars: []string{"REDIS_ENABLED_TLS"},
		Value:   false,
	}

	RedisInsecureSkipVerifyFlag = &cli.BoolFlag{
		Name:    "redis_insecure_skip_verify",
		Usage:   "insecure_skip_verify used for Redis tls verify",
		EnvVars: []string{"REDIS_INSECURE_SKIP_VERIFY"},
		Value:   true,
	}

	RedisDatabaseFlag = &cli.IntFlag{
		Name:    "redis_db",
		Usage:   "Redis database is using by application",
		EnvVars: []string{"REDIS_DB"},
		Value:   0,
	}

	RedisPoolSizeFlag = &cli.IntFlag{
		Name:    "redis_max_open_conns",
		Usage:   "sets the maximum number of open connections to the database",
		EnvVars: []string{"REDIS_POOL_SIZE"},
		Value:   10,
	}
)

// Log and notifier env
var (
	LogLevelFlag = &cli.StringFlag{
		Name:    "log_level",
		Usage:   "Level to log message to standard logger: panic, fatal, error, warn, warning, info, debug",
		EnvVars: []string{"LOG_LEVEL"},
		Value:   "debug",
	}
)

// Monitoring tool env
// var (
//	DatadogAgentHostFlag = &cli.StringFlag{
//		Name:    "dd_agent_host",
//		Usage:   "Define Datadog agent host",
//		EnvVars: []string{"DD_AGENT_HOST"},
//		Value:   "localhost",
//	}
//
//	DatadogAgentAPMPortFlag = &cli.StringFlag{
//		Name:    "dd_agent_apm_port",
//		Usage:   "Define Datadog agent port",
//		EnvVars: []string{"DD_AGENT_APM_PORT"},
//		Value:   "8126",
//	}
//
//	TracerEngineFlag = &cli.StringFlag{
//		Name:    "tracer_engine",
//		Usage:   "Define logger engine to use: datadog",
//		EnvVars: []string{"TRACER_ENGINE"},
//		Value:   "",
//	}
//
//	ProfilerEngineFlag = &cli.StringFlag{
//		Name:    "profiler_engine",
//		Usage:   "Define profiler engine to use: datadog",
//		EnvVars: []string{"PROFILER_ENGINE"},
//		Value:   "",
//	}
//)

//// EnabledProfilingFlag For pprof middleware
// var (
//	EnabledProfilingFlag = &cli.BoolFlag{
//		Name:    "enabled_pprof",
//		Usage:   "enable pprof middleware",
//		EnvVars: []string{"ENABLED_PPROF"},
//		Value:   false,
//	}
//)

// CORSAllowHostsFlag CORS env
var (
	CORSAllowHostsFlag = &cli.StringFlag{
		Name:    "cors_allow_hosts",
		Usage:   "List of origins a cross-domain request can be executed from, separated by comma",
		EnvVars: []string{"CORS_ALLOW_HOSTS"},
		Value:   "http://localhost:3000,https://car-uat-cms.eboost.vn,https://cms.eboost.vn,https://car-cms.eboost.vn",
	}
)

// Salt for password
var (
	SaltFlag = &cli.StringFlag{
		Name:    "salt",
		Usage:   "salt for password",
		EnvVars: []string{"SALT"},
	}
)

// JWTSecretFlag for jwt
var (
	JWTSecretFlag = &cli.StringFlag{
		Name:    "jwt_secret",
		Usage:   "jwt",
		EnvVars: []string{"JWT_SECRET"},
	}
)

// Realtime service flag
var (
	RealtimeURLFlag = &cli.StringFlag{
		Name:    "realtime_url",
		Usage:   "URL for realtime service",
		EnvVars: []string{"REALTIME_URL"},
	}

	RealtimeUsernameFlag = &cli.StringFlag{
		Name:    "realtime_username",
		Usage:   "Username for realtime service",
		EnvVars: []string{"REALTIME_USERNAME"},
	}

	RealtimePasswordFlag = &cli.StringFlag{
		Name:    "realtime_password",
		Usage:   "Password for realtime service",
		EnvVars: []string{"REALTIME_PASSWORD"},
	}

	RealtimeAPITimeoutFlag = &cli.IntFlag{
		Name:    "realtime_api_timeout",
		Usage:   "Timeout for realtime API",
		EnvVars: []string{"REALTIME_API_TIMEOUT"},
		Value:   120,
	}
)

// BasicAuth flags
var (
	BasicAuthUsernameFlag = &cli.StringFlag{
		Name:    "basic_auth_username",
		Usage:   "authenticate requests from internal services",
		EnvVars: []string{"BASIC_AUTHENTICATION_USERNAME"},
	}
	BasicAuthPasswordFlag = &cli.StringFlag{
		Name:    "basic_auth_password",
		Usage:   "authenticate requests from internal services",
		EnvVars: []string{"BASIC_AUTHENTICATION_PASSWORD"},
	}
)
