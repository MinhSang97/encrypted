package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	flags := []cli.Flag{
		// Common
		EnvFlag,
		AppNameFlag,
		AppVersionFlag,

		// MySQL
		MYSQLConnFlag,
		MYSQLHostFlag,
		MYSQLPortFlag,
		MYSQLMasterHostsFlag,
		MySQLSlaveHostsFlag,
		MySQLUserFlag,
		MySQLPasswordFlag,
		MySQLDatabaseFlag,
		MySQLMaxOpenConnsFlag,
		MySQLMaxIdleConnsFlag,
		MySQLConnMaxLifetimeFlag,

		// Redis
		RedisConnFlag,
		RedisHostFlag,
		RedisPortFlag,
		RedisUserFlag,
		RedisPasswordFlag,
		RedisDatabaseFlag,
		RedisPoolSizeFlag,
		RedisInsecureSkipVerifyFlag,
		RedisEnabledTLSFlag,

		HTTPPortFlag,

		LogLevelFlag,

		CORSAllowHostsFlag,

		SaltFlag,
		JWTSecretFlag,

		RealtimeURLFlag,
		RealtimeUsernameFlag,
		RealtimePasswordFlag,
		RealtimeAPITimeoutFlag,

		BasicAuthUsernameFlag,
		BasicAuthPasswordFlag,
	}

	app := &cli.App{
		Name:  "EBoost CMS PARTNER Service",
		Flags: flags,
		Action: func(ctx *cli.Context) error {
			srv := newService(ctx)

			return srv.start()
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
