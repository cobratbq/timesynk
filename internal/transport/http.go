// SPDX-License-Identifier: GPL-3.0-or-later
package transport

import (
	"net/http"
	"time"

	httputils "github.com/cobratbq/goutils/std/net/http"
	timeutils "github.com/cobratbq/goutils/std/time"
	"github.com/pkg/errors"
)

func QueryHttpsHeader(systemTimeRef, timestamp *time.Time, url string) (time.Time, error) {
	transp := http.DefaultTransport.(*http.Transport).Clone()
	transp.TLSClientConfig.Time = timeutils.TimeDeltaCorrectionFunc(systemTimeRef, timestamp)
	client := http.Client{Transport: transp}
	resp, err := client.Head(url)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to query https server")
	}
	return httputils.ExtractResponseHeaderDate(resp)
}

func QueryHttpHeader(url string) (time.Time, error) {
	resp, err := http.Head(url)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to query http server for 'Date' header")
	}
	return httputils.ExtractResponseHeaderDate(resp)
}
