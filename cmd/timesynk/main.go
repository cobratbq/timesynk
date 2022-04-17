// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"syscall"
	"time"

	osutils "github.com/cobratbq/goutils/std/os"
)

const HWCLOCK_FORMAT = "2006-01-02 15:04:05.999999-07:00"

// TODO handle user input via commandline args in nicer way than panicking
func main() {
	options := configureGlobal()
	processor := processAll(options.SetProcessor, options.PrintProcessor)
	cmd := flag.Arg(0)
	switch cmd {
	case "":
		flag.CommandLine.Usage()
		osutils.ExitWithError(1, "\nSpecify command for operation to perform.")
	case "checkpoint":
		cmdSnapshot(flag.Args()[1:])
	case "sync":
		cmdSync(processor, flag.Args()[1:])
	default:
		osutils.ExitWithError(1, "Unknown command specified: "+cmd)
	}
}

func configureGlobal() globalOptions {
	options := globalOptions{}
	flag.BoolVar(&options.Set, "set", false, "Immediately set hardware clock by the time value.")
	flag.BoolVar(&options.Print, "print", false, "Print formatted output of synchronized time.")
	// TODO add format specifier flag (`hwclock`, ...)?
	flag.CommandLine.Usage = func() {
		flag.CommandLine.Output().Write([]byte("Global options:\n"))
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.Output().Write([]byte(`
Subcommands:
  checkpoint
	record current system time in the modification-time of a file
  sync
	acquire date/time from a file or a webserver
`))
	}
	flag.Parse()
	if options.Set {
		options.SetProcessor = handleSetHardwareClock
	}
	if options.Print {
		options.PrintProcessor = prepareProcessPrintDateTime(HWCLOCK_FORMAT)
	}
	return options
}

type globalOptions struct {
	Set            bool
	SetProcessor   TimeProcessor
	Print          bool
	PrintProcessor TimeProcessor
}

// TimeProcessor defines the format for functions that can handle date/time
// values acquired through synchronization.
type TimeProcessor func(time.Time) error

func processAll(proc1 TimeProcessor, proc2 TimeProcessor) TimeProcessor {
	if proc1 == nil && proc2 == nil {
		return func(time.Time) error { return ErrNoProcessors }
	}
	if proc1 == nil {
		return proc2
	}
	if proc2 == nil {
		return proc1
	}
	return func(value time.Time) error {
		if err := proc1(value); err != nil {
			log.Printf("%v: %s\n", proc1, err.Error())
		}
		if err := proc2(value); err != nil {
			log.Printf("%v: %s\n", proc2, err.Error())
		}
		return nil
	}
}

// prepareProcessPrintDateTime defines a TimeProcessor function that
// outputs/prints the time according to the specified format.
func prepareProcessPrintDateTime(format string) TimeProcessor {
	return func(value time.Time) error {
		os.Stdout.WriteString(value.Format(format))
		return nil
	}
}

// handleSetHardwareClock directly sets the hardware clock to the provided
// value for time.
func handleSetHardwareClock(value time.Time) error {
	timeval := syscall.NsecToTimeval(value.UnixNano())
	return syscall.Settimeofday(&timeval)
}

// ErrNoProcessors indicates that there are no time processors configured.
var ErrNoProcessors = errors.New("no processors configured")
