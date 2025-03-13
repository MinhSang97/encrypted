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

type RunCommand struct {
	env string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newRunCommand(m *Main) *RunCommand {
	return &RunCommand{
		env:    m.env,
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}

func (cmd *RunCommand) Run(args ...string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	help := fs.Bool("h", false, "")

	if err := fs.Parse(args); err != nil {
		return err
	} else if *help {
		_, _ = fmt.Fprintln(cmd.Stderr, cmd.Usage())

		return ErrUsage
	}

	jobName := fs.Arg(0)
	if jobName == "" {
		return ErrJobNameRequired
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

	if err := s.Run(jobName); err != nil {
		return fmt.Errorf("err running job %s : %w", jobName, err)
	}

	return nil
}

func (cmd *RunCommand) Usage() string {
	return strings.TrimLeft(`
usage: run the specific job manually
`, "\n")
}
