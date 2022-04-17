// SPDX-License-Identifier: GPL-3.0-or-later
package timesync

import (
	"time"

	timeutils "github.com/cobratbq/goutils/std/time"
	"github.com/cobratbq/httptime/internal/transport"
	"github.com/pkg/errors"
)

// SyncHttpTime reads a timestamp from the 'Date' header of the http-server response.
func SyncHttpTime(url string) (time.Time, error) {
	return transport.QueryHttpHeader(url)
}

// SyncHttpsTime starts with synchronizing with http server, then using that time reference,
// attempts to establish a TLS-session with an https-server (in the process confirming that the
// time is sufficiently accurate) and acquire a new date/time reference from the response headers.
func SyncHttpsTime(httpURL, httpsURL string) (LocalRemoteTime, error) {
	// Get rough date-time from HTTP request to have a reference time for successful HTTPS (TLS)
	// request.
	systemRequestTime := time.Now()
	serverTimestamp, err := transport.QueryHttpHeader(httpURL)
	if err != nil {
		return LocalRemoteTime{}, errors.Wrap(err, "Unable to parse date-time information from 'Date' header.")
	}
	// Execute HTTPS request, which requires at least approximately correct time.
	preHttpsRequestTimestamp := time.Now()
	// The https request is affected by DNS resolution speed, TLS session establishment, etc., so
	// time difference is truly just a rough approximation.
	httpsRequestDate, err := transport.QueryHttpsHeader(&systemRequestTime, &serverTimestamp, httpsURL)
	postHttpsRequestTimestamp := time.Now()
	if err != nil {
		return LocalRemoteTime{}, err
	}
	// Approximating local system's time around time that https response gets constructed, by
	// assuming equal parts for sending request and receiving response.
	systemTimeDuringRequest := timeutils.MiddleTimestamps(preHttpsRequestTimestamp, postHttpsRequestTimestamp)
	return LocalRemoteTime{
		Local:  systemTimeDuringRequest,
		Remote: httpsRequestDate,
	}, nil
}
