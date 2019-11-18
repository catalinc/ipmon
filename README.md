# ipmon

`ipmon` is small command line utility to monitor changes in the hostname and IP addresses of the local host.

# Prerequisites

[Golang](https://golang.org/doc/install)

Tested with Go 1.11.

# Installation

1. Install the utility with `go get github.com/catalinc/ipmon` 
2. Edit mailer configuration according to your setup. You can use the `mail.json` file as template.
3. Run it:
```bash
ipmon -help # to view options OR

ipmon -interval <networkCheckIntervalInSeconds> \
      -netConfig </path/to/save/currentNetworkConfiguration.json> \
      -mailConfig </path/to/mailerConfiguration.json> \
      -runOnce 
``` 

The utility checks network configuration for changes at regular intervals.

When changes are detected (i.e. a new IP address) the utility sends an email alert to a preconfigured list of recipients.

# Warning

The SMTP password is stored in clear in `mail.json` configuration file.

So consider yourself warned ;-).
