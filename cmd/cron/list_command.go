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

type ListCommand struct {
	env string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func newListCommand(m *Main) *ListCommand {
	return &ListCommand{
		env:    m.env,
		Stdin:  m.Stdin,
		Stdout: m.Stdout,
		Stderr: m.Stderr,
	}
}

func (cmd *ListCommand) Run(args ...string) error {
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

	var (
		out          strings.Builder
		scheduledJob = s.ScheduledJobs()
	)

	out.WriteString("Scheduled jobs: \n")

	for name, schedules := range scheduledJob {
		out.WriteString(fmt.Sprintf("- %v on [ %v ] \n", name, strings.Join(schedules, " , ")))
	}

	_, _ = fmt.Fprint(cmd.Stdout, out.String())

	return nil
}

func (cmd *ListCommand) Usage() string {
	return strings.TrimLeft(`
usage: list - lists the jobs and their schedules 
`, "\n")
}
