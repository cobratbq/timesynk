// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"flag"
	"time"

	osutils "github.com/cobratbq/goutils/std/os"
	"github.com/cobratbq/httptime/internal/timefile"
)

const CHECKPOINT_FILE = "./checkpoint"

func cmdSnapshot(args []string) {
	config := checkpointOptions{}
	configureCheckpoint(args)
	if err := snapshotCheckpoint(config.checkpointPath); err != nil {
		osutils.ExitWithError(1, "Failed to update mod-time: "+err.Error()+"\n")
	}
}

func snapshotCheckpoint(path string) error {
	// TODO check for existence of file, if not create?
	return timefile.UpdateTime(path, time.Now())
}

func configureCheckpoint(args []string) checkpointOptions {
	config := checkpointOptions{}
	flagset := flag.NewFlagSet("checkpoint", flag.ExitOnError)
	flagset.StringVar(&config.checkpointPath, "file", CHECKPOINT_FILE, "The filename used to record the last checkpoint timestamp.")
	flagset.Parse(args)
	return config
}

type checkpointOptions struct {
	checkpointPath string
}
