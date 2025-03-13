package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EBOOST-LTD/eboost-cms-partner-BE/cronjobs"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/pkg/log"
	"github.com/spf13/viper"
)

var (
	ErrUnknownCommand  = errors.New("unknown command")
	ErrUsage           = errors.New("usage")
	ErrJobNameRequired = errors.New("job name is required")
)

const (
	notifierEngineRollbar = "rollbar"
	cronJobsConfigFile    = "assets/cron_jobs.yaml"
)

//nolint:gomnd
func main() {
	viper.AutomaticEnv()

	m := NewMain(viper.GetString("ENV"))
	if err := m.Run(os.Args[1:]...); err != nil {
		if errors.Is(err, ErrUsage) {
			os.Exit(2)
		}

		log.WithError(err).Errorln()
		os.Exit(1)
	}
}

type Main struct {
	env string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NewMain(env string) *Main {
	return &Main{
		env:    env,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run is entry point
func (m *Main) Run(args ...string) error {
	// Require a command at the beginning.
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		fmt.Fprintln(m.Stderr, m.Usage())

		return ErrUsage
	}

	switch args[0] {
	case "help":
		fmt.Fprintln(m.Stderr, m.Usage())

		return ErrUsage
	case "start":
		return newStartCommand(m).Run(args[1:]...)
	case "run":
		return newRunCommand(m).Run(args[1:]...)
	case "list":
		return newListCommand(m).Run(args[1:]...)
	default:
		return ErrUnknownCommand
	}
}

func (m *Main) Usage() string {
	return strings.TrimLeft(`
eboost-cms-partner-BE_cron for running cron jobs in eboost-cms-partner-BE
Usage:
	eboost-cms-partner-BE_cron command [arguments]
The commands are:
- start starts the scheduler
- run   runs a specific cronjob
- list  lists all the scheduled jobs

Use "eboost-cms-partner-BE_cron [command] -h" for more information about a command.
`, "\n")
}

func loadConfig() {
	// load config
	if err := cronjobs.LoadConfig(); err != nil {
		panic(fmt.Errorf("err loading env config: %w", err))
	}

	if err := log.SetLevel(cronjobs.GetConfig().LogLevel); err != nil {
		panic(err)
	}
}
