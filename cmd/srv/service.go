package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/EBOOST-LTD/eboost-cms-partner-BE/app/config"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/app/external/adapter/restclient/realtime"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/app/external/framework"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/app/interface/api/payload"
	pkgConfig "github.com/EBOOST-LTD/eboost-cms-partner-BE/pkg/config"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/pkg/log"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/pkg/profiler"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/pkg/redis"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/pkg/timeutil"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/pkg/tracer"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const (
	// Constant for notifier engine
	// notifierEngineRollbar = "rollbar"
	logLevelDebug     = "debug"
	payloadConfigFile = "assets/payload_field_config.yaml"
)

type service struct {
	cfg *config.Config
}

func newService(ctx *cli.Context) *service {
	s := &service{}

	s.loadConfig(ctx)

	if err := log.SetLevel(s.cfg.LogLevel); err != nil {
		panic(err)
	}

	if err := payload.LoadGlobalPayloadConfig(payloadConfigFile); err != nil {
		panic(errors.Wrapf(err, "failed to load payload config file at %v", payloadConfigFile))
	}

	realtime.Setup(
		s.cfg.RealtimeService.URL,
		s.cfg.Env,
		s.cfg.RealtimeService.Username,
		s.cfg.RealtimeService.Password,
		s.cfg.RealtimeService.Timeout,
		timeutil.NewTimeFactory(),
	)

	return s
}

func (s *service) loadConfig(ctx *cli.Context) {
	conf := &config.Config{
		Env: ctx.String(EnvFlag.Name),
		HTTPServer: config.ServerAddr{
			Port: ctx.String(HTTPPortFlag.Name),
		},
		MySQL: config.MySQL{
			ConnectionString: ctx.String(MYSQLConnFlag.Name),
			Host:             ctx.String(MYSQLHostFlag.Name),
			Port:             ctx.Int64(MYSQLPortFlag.Name),
			Masters:          ctx.String(MYSQLMasterHostsFlag.Name),
			Slaves:           ctx.String(MySQLSlaveHostsFlag.Name),
			User:             ctx.String(MySQLUserFlag.Name),
			Password:         ctx.String(MySQLPasswordFlag.Name),
			DB:               ctx.String(MySQLDatabaseFlag.Name),
			MaxOpenConns:     ctx.Int(MySQLMaxOpenConnsFlag.Name),
			MaxIdleConns:     ctx.Int(MySQLMaxIdleConnsFlag.Name),
			ConnMaxLifetime:  ctx.Int(MySQLConnMaxLifetimeFlag.Name),
			IsEnabledLog:     ctx.String(LogLevelFlag.Name) == logLevelDebug,
		},
		Redis: redis.Config{
			ConnectionString:   ctx.String(RedisConnFlag.Name),
			Host:               ctx.String(RedisHostFlag.Name),
			Port:               ctx.Int(RedisPortFlag.Name),
			User:               ctx.String(RedisUserFlag.Name),
			Pass:               ctx.String(RedisPasswordFlag.Name),
			Database:           ctx.Int(RedisDatabaseFlag.Name),
			PoolSize:           ctx.Int(RedisPoolSizeFlag.Name),
			InsecureSkipVerify: ctx.Bool(RedisInsecureSkipVerifyFlag.Name),
			TLS:                ctx.Bool(RedisEnabledTLSFlag.Name),
		},
		CORS: config.CORS{
			AllowHosts: strings.Split(ctx.String(CORSAllowHostsFlag.Name), ","),
		},
		LogLevel:  ctx.String(LogLevelFlag.Name),
		ProxyURL:  ctx.String(ProxyURLFlag.Name),
		Salt:      ctx.String(SaltFlag.Name),
		JWTSecret: ctx.String(JWTSecretFlag.Name),
		RealtimeService: config.RealtimeService{
			URL:      ctx.String(RealtimeURLFlag.Name),
			Username: ctx.String(RealtimeUsernameFlag.Name),
			Password: ctx.String(RealtimePasswordFlag.Name),
			Timeout:  ctx.Int(RealtimeAPITimeoutFlag.Name),
		},
		BasicAuth: config.BasicAuth{
			Username: ctx.String(BasicAuthUsernameFlag.Name),
			Password: ctx.String(BasicAuthPasswordFlag.Name),
		},
	}

	s.cfg = conf

	pkgConfig.SetConfig(conf)
}

func (s *service) start() error {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)

	s.initDBConnection(s.cfg.MySQL)
	s.initCookieAuthenticator(s.cfg.Env)
	s.initRedisConnection(s.cfg.Redis)

	// #nosec
	srv := &http.Server{
		Addr:    ":" + s.cfg.HTTPServer.Port,
		Handler: framework.Handler(s.cfg),
	}

	go func() {
		tracer.Start()

		if err := profiler.Start(); err != nil {
			log.WithError(err).Errorln("profiler fail")
		}

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Errorln("ListenAndServe fail")
			panic(err)
		}
	}()

	log.WithField("port", s.cfg.HTTPServer.Port).Debugln("server started")

	<-stopSignal

	tracer.Stop()
	profiler.Stop()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Info("server stopped")

	return nil
}
