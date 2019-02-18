# Woleet command line interface

woleet-cli is an open source command line tool on top of the Woleet API.
The tool is written in Go and has been tested on Windows, macOS and Linux.

Currently, the tool only supports:

* the `anchor` command, allowing to recursively anchor all files in a given directory
* the `sign` command, allowing to recursively sign all files in a given directory (using Woleet.ID Server: <https://github.com/woleet/woleet.id-server>)
* the `export` command, allowing to download all your receipts in a given directory

## Anchor / Sign

### Functionalities

The tool scans a folder recursively and anchors or sign all files found. It also gathers proof receipts and stores them beside anchored or signed files (in a Chainpoint file named 'filename'-'anchorID'.(anchor|signature)-receipt.json).

Since anchoring is not a realtime operation, the tool is supposed to be run on a regular basis (or at least a second time when all proof receipts are ready to download). Obviously, the files that were already anchored are not re-anchored.

If the option --strict is provided, for each file that already have a proof receipt, the tool checks that the hash of the file still matches the hash in the receipt (to detect file changes). If they differ, the file is re-anchored and the old receipt is kept, except if --prune is set in that case the old receipt is deleted.  
If the original file is no longer present and the option --prune is provided, the old receipt/pending file will be deleted.

To sum up, this tool can be used to generate and maintain the set of timestamped proofs of existence for all files in a given directory.

Note: tags are added to the anchors according to the name of sub-folders  

### Limitations

* All files and folders beginning by '.' or finished by '.(anchor|signature)-(receipt|pending).json' are ignored
* Symlinks are not followed
* Scanned sub-folders cannot have a space in their name
* The maximum length of the subfolder path (without delimiters) is 128 characters

## Export

### Functionalities

The tool dumps all your receipts into a folder.  
You can specify a limit date to get all receipt created from this date.  

### Limitations

* Each receipt will be named: 'anchor name'-'anchor ID'.(anchor|signature)-receipt.json
* Exporting can be quite long, as each recepit is downloaded individually

## Configuration

The tool behavior can be configured using command line arguments, environment variables or a configuration file. When several configuration means are used, the following priorities are applied:

- command line arguments
- environment variables
- config file
- default value (if any)

There is also a special environnement variable or config path to disable environnement configuration and configuration file:

```bash
woleet-cli --config "DISABLED" ...
# or
export WCLI_CONFIG="DISABLED"
```

### Usage

```
Usage:
woleet-cli anchor [flags]
  -d, --directory string   source directory containing files to anchor (required)
      --dryRun             print information about files to anchor without anchoring
  -e, --exitOnError        exit with an error code if anything goes wrong
  -h, --help               help for anchor
  -p, --private            create non discoverable proofs
      --prune              delete receipts that are not along the original file,
                           with --strict it checks the hash of the original file and deletes the receipt if they do not match
  -r, --recursive          explore sub-folders recursively
      --strict             re-anchor any file that has changed since last anchoring


woleet-cli sign [flags]
  -d, --directory string         source directory containing files to sign (required)
      --dryRun                   print information about files to sign without signing
  -e, --exitOnError              exit with an error code if anything goes wrong
  -h, --help                     help for sign
      --widsPubKey string    public key (ie. bitcoin address) to use to sign
      --widsSignURL string   Woleet.ID Server sign URL ex: "https://idserver.com:4443/sign" (required)
      --widsToken string     Woleet.ID Server API token (required)
      --widsUnsecureSSL      do not check Woleet.ID Server's SSL certificate validity (only for developpement)
  -p, --private                  create non discoverable proofs
      --prune                    delete receipts that are not along the original file,
                                 with --strict it checks the hash of the original file and deletes the receipt if they do not match
  -r, --recursive                explore sub-folders recursively
      --strict                   re-sign any file that has changed since last signature

woleet-cli export [flags]
  -d, --directory string   directory where to store the proofs (required)
  -e, --exitOnError        exit with an error code if anything goes wrong
  -h, --help               help for export
  -l, --limitDate string   get only proofs created after the provided date (format: yyyy-MM-dd)

Global Flags:
  -c, --config string     config file (default is $HOME/.woleet-cli.yaml)
  -h, --help              help for woleet-cli
      --json              use JSON as log output format
      --logLevel string   select log level info|warn|error|fatal (default "info")
  -t, --token string      Woleet API token (required)
  -u, --url string        Woleet API URL (default "https://api.woleet.io/v1")
      --version           version for woleet-cli

```

### Configuration file format

YAML:

```yaml
api:
  url: https://api.woleet.io/v1
  token: insert-your-token-here
  private: true
app:
  directory: /home/folder/to/anchor
  strict: true
  prune: true
  exitOnError: true
  recursive: true
  dryRun: false
sign:
  widsSignURL: https://idserver.com:4443/sign
  widsToken: insert-your-idserver-token-here
  widsPubKey: insert-your-idserver-pubkey-here
  widsUnsecureSSL: false
export:
  directory: /home/folder/to/anchor
  limitDate: 2018-01-21
  exitOnError: true
log:
  json: true
  level: info
```

JSON:

```json
{
  "api": {
    "url": "https://api.woleet.io/v1",
    "token": "insert-your-token-here",
    "private": true,
  },
  "app": {
    "directory": "/home/folder/to/anchor",
    "exitOnError": true,
    "strict": true,
    "prune": true,
    "recursive": true,
    "dryRun": true
  },
  "sign": {
    "widsSignURL": "https://idserver.com:4443/sign",
    "widsToken": "insert-your-idserver-token-here",
    "widsPubKey": "insert-your-idserver-pubkey-here",
    "widsUnsecureSSL": false
  },
  "export": {
    "directory": "/home/folder/to/anchor",
    "limitDate": "2018-01-21",
    "exitOnError": true
  },
  "log": {
    "json": true,
    "level": "info"
  }
}
```

ENV:

```bash
export WCLI_CONFIG="$HOME/.woleet-cli.json"
export WCLI_API_URL="https://api.woleet.io/v1"
export WCLI_API_TOKEN="insert-your-token-here"
export WCLI_API_PRIVATE="true"
export WCLI_APP_DIRECTORY="/home/folder/to/anchor"
export WCLI_APP_EXITONERROR="true"
export WCLI_APP_STRICT="true"
export WCLI_APP_PRUNE="true"
export WCLI_APP_RECURSIVE="true"
export WCLI_APP_DRYRUN="true"
export WCLI_SIGN_WIDSSIGNURL="https://idserver.com:4443/sign"
export WCLI_SIGN_WIDSTOKEN="insert-your-idserver-token-here"
export WCLI_SIGN_WIDSPUBKEY="insert-your-idserver-pubkey-here"
export WCLI_SIGN_WIDSUNSECURESSL="false"
export WCLI_EXPORT_DIRECTORY="/home/folder/to/anchor"
export WCLI_EXPORT_LIMITDATE="2018-01-21"
export WCLI_EXPORT_EXITONERROR="true"
export WCLI_LOG_JSON="true"
export WCLI_LOG_LEVEL="info"
```

## Install

see: [INSTALL](INSTALL.md)

## Build

see: [BUILD](BUILD.md)
