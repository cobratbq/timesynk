# timesynk

Time synchronization using http and https protocols. Approximate synchronization of time from a source that is accurate enough to enable TLS communication on the internet.

First, synchronizes time with an 'http'-server. Next, use that time to establish a TLS-secured connection with an https-server. If we can establish secure communication, the acquired time is reasonably accurate. As we test the secure connection, we acquire a new timestamp. Next, output time in `hwclock` format and/or set time directly using a syscall.

## Usage

`timesynk`:

- `-help` print global parameters
- `-print` to output `hwclock`-formatted time on stdout
- `-set` to directly update time on the system (requires elevated privileges, e.g. `sudo`)

Subcommands:

Follow subcommand with `-help` to print options.

- `checkpoint` write current system time to a "checkpoint" file.
  - `-path` path to the checkpoint-file.

- `sync` synchronize system time with the specified source.
  - `-checkpoint` enable 'checkpoint' synchronization  
    - `-path` path to checkpoint-file from which to take the modification time
  - `-web` enable 'web' synchronization: _Acquire a timestamp from the http-request, immediately test acquired time with an https-request and acquires a new timestamp in the process. If successful, use the date/time._
    - `-http` url to http-server
    - `-https` url to https-server

For web sources, any http and https server will do. Defaults in the application use the debian package mirrors as these are expected to be fast, reliable and as basic as to serve mere static files.

## Use cases

Use cases may vary and depends on circumstances. Small board computers and system-on-chips do not always retain time on powerloss. Many servers, services and tools require reasonably accurate time to even work, because of their security properties. This is a convenient way to get time from a ubiquitous source that is practical.

## Building

Build with: `go build ./cmd/timesynk`
