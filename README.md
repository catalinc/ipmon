# ipmon

`ipmon` is small command line utility to monitor changes in the IP addresses of the local host.

When changes are detected (i.e. a new IP address) the utility will send an email alert to a preconfigured list of recipients.

# Installation

1. Install the CLI utility with `go install github.com/catalinc/ipmon`
2. Edit mail configuration according to your setup. Use the `mail.json` file as template.
3. Run the utility:
```bash
ipmon -help # to view options OR

ipmon -interval <networkCheckIntervalInSeconds> \
      -netConfig </path/to/save/currentNetworkConfiguration.json> \
      -mailConfig </path/to/mailerConfiguration.json>
``` 

# Warning

The utility requires the SMTP password to be stored in clear in `mail.json` configuration file.
So, consider yourself warned ;-).