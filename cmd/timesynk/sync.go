// SPDX-License-Identifier: GPL-3.0-or-later
package main

import (
	"flag"
	"strings"
	"time"

	assert "github.com/cobratbq/goutils/assert"
	osutils "github.com/cobratbq/goutils/std/os"
	"github.com/cobratbq/httptime/internal/timefile"
	"github.com/cobratbq/httptime/internal/timesync"
)

// 🗹 checkpoint file mod-time
// 🗹 http
// 🗹 https
// ☐ deb metadata (freshness 'ValidUntil' based on "authentic" source)
func cmdSync(handle TimeProcessor, args []string) {
	cfg := configureSync(args)
	if cfg.checkpointSync {
		timestamp, err := timefile.ReadTime(cfg.checkpointPath)
		if err != nil {
			osutils.ExitWithError(1, "Failed to sync with checkpoint file: "+err.Error())
		}
		if err = handle(timestamp); err != nil {
			osutils.ExitWithError(1, "Failed to process time from checkpoint file: "+err.Error())
		}
	}
	if cfg.webSync {
		synctime, err := timesync.SyncHttpsTime(cfg.httpURL, cfg.httpsURL)
		if err != nil {
			osutils.ExitWithError(1, "Failed to query timestamp from HTTPS server: "+err.Error())
		}
		if err := handle(synctime.Remote.Add(time.Since(synctime.Local))); err != nil {
			osutils.ExitWithError(1, "Failed to process time from synchronization: "+err.Error())
		}
	}
}

func configureSync(args []string) syncOptions {
	cfg := syncOptions{}
	syncset := flag.NewFlagSet("sync", flag.ExitOnError)
	syncset.BoolVar(&cfg.checkpointSync, "checkpoint", false, "Restore time snapshot from checkpoint file.")
	syncset.StringVar(&cfg.checkpointPath, "path", CHECKPOINT_FILE, "Synchronize with previous snapshot stored in checkpoint-file.")
	syncset.BoolVar(&cfg.webSync, "web", false, "Enable web-sync (using the urls as provided in arguments)")
	syncset.StringVar(&cfg.httpURL, "http-url", "http://deb.debian.org/debian/", "HTTP (insecure) URL for time-synchronization")
	syncset.StringVar(&cfg.httpsURL, "https-url", "https://deb.debian.org/debian/", "HTTPS URL for time-synchronization")
	syncset.Parse(args)
	assert.True(strings.HasPrefix(cfg.httpURL, "http://"))
	assert.True(strings.HasPrefix(cfg.httpsURL, "https://"))
	return cfg
}

type syncOptions struct {
	checkpointSync bool
	checkpointPath string
	webSync        bool
	httpURL        string
	httpsURL       string
}
