package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/EBOOST-LTD/eboost-cms-partner-BE/cronjobs"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/pkg/cronjob"
	"github.com/EBOOST-LTD/eboost-cms-partner-BE/pkg/cronjob/scheduler"
)

type StartCommand struct {
	env string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newStartCommand(m *Main) *StartCommand {
	return &StartCommand{
		env:    m.env,
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}

func (cmd *StartCommand) Run(args ...string) error {
	// Parse flags.
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	help := fs.Bool("h", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	} else if *help {
		_, _ = fmt.Fprintln(cmd.Stderr, cmd.Usage())

		return ErrUsage
	}

	loadConfig()

	// load available jobs
	jobConfigs, err := cronjob.LoadCronJobConfig(cronJobsConfigFile)
	if err != nil {
		return fmt.Errorf("err loading job config: %w", err)
	}

	// set up scheduler
	s, err := scheduler.Setup(cmd.env, jobConfigs, cronjobs.ProvideJobs())
	if err != nil {
		return fmt.Errorf("err setting up scheduler: %w", err)
	}

	s.Start()

	return nil
}

func (cmd *StartCommand) Usage() string {
	return strings.TrimLeft(`
usage: start predefined jobs, with schedules in cron_jobs.yaml 
`, "\n")
}
