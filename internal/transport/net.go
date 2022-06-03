// SPDX-License-Identifier: GPL-3.0-or-later
package transport

import (
	"context"
	"net"
	"net/url"
)

// lookupURLHost parses and resolves the URL's hostname.
func LookupURLHost(rawurl string) error {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	_, err = net.DefaultResolver.LookupIPAddr(context.Background(), uri.Hostname())
	return err
}
