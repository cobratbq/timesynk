// SPDX-License-Identifier: GPL-3.0-or-later
package timesync

import (
	"time"

	time_ "github.com/cobratbq/goutils/std/time"
	"github.com/cobratbq/httptime/internal/transport"
	"github.com/pkg/errors"
)

// SyncHttpsTime starts with synchronizing with http server, then using that time reference,
// attempts to establish a TLS-session with an https-server (in the process confirming that the
// time is sufficiently accurate) and acquire a new date/time reference from the response headers.
func SyncHttpsTime(httpURL, httpsURL string) (LocalRemoteTime, error) {
	delta, err := SyncHttpTime(httpURL)
	if err != nil {
		return LocalRemoteTime{}, errors.Wrap(err, "Unable to parse date-time information from 'Date' header.")
	}
	return syncHttpsTime(httpsURL, delta)
}

// SyncHttpTime reads a timestamp from the 'Date' header of the http-server response.
func SyncHttpTime(httpURL string) (LocalRemoteTime, error) {
	if err := transport.LookupURLHost(httpURL); err != nil {
		return LocalRemoteTime{}, err
	}
	systemRequestTime := time.Now()
	remoteTime, err := transport.QueryHttpHeader(httpURL)
	if err != nil {
		return LocalRemoteTime{}, err
	}
	return LocalRemoteTime{Local: systemRequestTime, Remote: remoteTime}, nil
}

// syncHttpsTime acquires the date/time from an https request, using provided local-remote time
// delta to correct for (possibly) bad system time to establish a secure session.
func syncHttpsTime(httpsURL string, delta LocalRemoteTime) (LocalRemoteTime, error) {
	if err := transport.LookupURLHost(httpsURL); err != nil {
		return LocalRemoteTime{}, err
	}
	// Execute HTTPS request, which requires at least approximately correct time. Furthermore,
	// the https request is affected by TLS session establishment, etc., so time difference is
	// truly just a rough approximation.
	preHttpsRequestTimestamp := time.Now()
	httpsRequestDate, err := transport.QueryHttpsHeader(&delta.Local, &delta.Remote, httpsURL)
	postHttpsRequestTimestamp := time.Now()
	if err != nil {
		return LocalRemoteTime{}, err
	}
	// Approximating local system's time around time that https response gets constructed, by
	// assuming equal parts for sending request and receiving response.
	systemTimeDuringRequest := time_.MiddleTimestamps(preHttpsRequestTimestamp, postHttpsRequestTimestamp)
	return LocalRemoteTime{
		Local:  systemTimeDuringRequest,
		Remote: httpsRequestDate,
	}, nil
}
