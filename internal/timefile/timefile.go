// SPDX-License-Identifier: GPL-3.0-or-later
package timefile

import (
	"os"
	"time"
)

// UpdateTime updates the times for the file located at filePath.
func UpdateTime(filePath string, timestamp time.Time) error {
	return os.Chtimes(filePath, timestamp, timestamp)
}

// ReadTIme reads the modification time from the file located at filePath.
func ReadTime(filePath string) (time.Time, error) {
	info, err := os.Lstat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}
