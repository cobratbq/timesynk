// SPDX-License-Identifier: GPL-3.0-or-later
package timesync

import "time"

// LocalRemoteTime contains both the system time and the reference time for the same instant.
type LocalRemoteTime struct {
	Local  time.Time
	Remote time.Time
}
