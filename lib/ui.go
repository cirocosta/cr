package lib

import (
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

type ActivityType int

const (
	ActivityUnknown ActivityType = iota
	ActivityStarted
	ActivityErrored
	ActivitySuccess
	ActivityAborted
)

type Activity struct {
	Type ActivityType
	Time time.Time
	Job  *Job
}

var (
	ActivityMapping = map[ActivityType]string{
		ActivityAborted: "ABORTED",
		ActivityStarted: "STARTED",
		ActivityErrored: "ERRORED",
		ActivitySuccess: "SUCCESS",
		ActivityUnknown: "UNKNOWN",
	}
	WriterMapping = map[ActivityType]*color.Color{
		ActivityAborted: color.New(color.FgYellow),
		ActivityStarted: color.New(color.FgBlue),
		ActivityErrored: color.New(color.FgRed),
		ActivitySuccess: color.New(color.FgGreen),
		ActivityUnknown: color.New(color.FgCyan),
	}
)

type Ui struct {
	writer *tabwriter.Writer
	sync.Mutex
}

func NewUi() (u Ui) {
	u.writer = new(tabwriter.Writer)
	u.writer.Init(os.Stdout, 10, 8, 2, '\t', 0)

	return
}

func (u *Ui) WriteActivity(a *Activity) (err error) {
	u.Lock()
	defer u.Unlock()

	switch a.Type {
	case ActivityStarted:
		WriterMapping[a.Type].
			Fprintf(u.writer, "%s\tstatus=%s\tstart=%s\n",
				a.Job.Id,
				ActivityMapping[a.Type],
				time.Now().Format("15:04:05"))
	case ActivityErrored, ActivitySuccess, ActivityAborted:
		WriterMapping[a.Type].
			Fprintf(u.writer, "%s\tstatus=%s\tstart=%s\telapsed=%s\n",
				a.Job.Id,
				ActivityMapping[a.Type],
				a.Job.StartTime.Format("15:04:05"),
				a.Job.EndTime.Sub(*a.Job.StartTime).String())
	default:
		err = errors.Errorf(
			"unknown activity type %+v", a)
		return
	}

	u.writer.Flush()
	return
}
