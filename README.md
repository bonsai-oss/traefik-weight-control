### traefik-weight-control

*A cli tool adjusting and listing the weight of weighted round robin Traefik services.*

#### Usage
```
usage: aectl --file=FILE [<flags>] <command> [<args> ...]

Traefik Weight Control

Flags:
  -h, --help       Show context-sensitive help (also try --help-long and --help-man).
      --version    Show application version.
  -v, --verbose    Enable debug mode
  -f, --file=FILE  Path to the Traefik configuration file

Commands:
  help [<command>...]
    Show help.


  list [<flags>]
    List all services and servers

    -o, --format=text  Output format

  set --service=SERVICE --server=SERVER --weight=WEIGHT [<flags>]
    Set the weight of a server

    -d, --dry-run          Dry run
    -s, --service=SERVICE  Service name
    -n, --server=SERVER    Server name
    -w, --weight=WEIGHT    Server weight
```