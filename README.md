# Woleet command line interface

woleet-cli is an open source command line tool on top of the Woleet API.
The tool is written in Go and has been tested on Windows, macOS and Linux.

Currently, the tool only supports:

* the `anchor` command, allowing to recursively anchor all files in a given directory
* the `sign` command, allowing to recursively sign all files in a given directory (using the backend kit: <https://github.com/woleet/woleet-backendkit>)
* the `export` command, allowing to download all your receipts in a given directory

## Anchor / Sign

### Functionalities

The tool scans a folder recursively and anchors or sign all files found. It also gathers proof receipts and stores them beside anchored or signed files (in a Chainpoint file named 'filename'-'anchorID'.(anchor|signature)-receipt.json).

Since anchoring is not a realtime operation, the tool is supposed to be run on a regular basis (or at least a second time when all proof receipts are ready to download). Obviously, the files that were already anchored are not re-anchored.

If the option --strict is provided, for each file that already have a proof receipt, the tool checks that the hash of the file still matches the hash in the receipt (to detect file changes). If they differ, the file is re-anchored and the old receipt is kept, except if --strict--prune is used instead.  
If the original file is no longer present and the option --strict-prune is provided, the old receipt/pending file will be deleted.

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
export WLT_CONFIG="DISABLED"
```

### Usage

```
Usage:
woleet-cli anchor [flags]
  -d, --directory string   source directory containing files to anchor (required)
      --dryrun             Print information about the files to anchor without anchoring
  -e, --exitonerror        exit the app with an error code if anything goes wrong
  -h, --help               help for anchor
  -p, --private            create anchors with non-public access
  -r, --recursive          Explore subfolders
      --strict             re-anchor any file that has changed since last anchoring
      --strict-prune       same as --strict, plus delete the previous anchoring receipt


woleet-cli sign [flags]
      --backendkitPubKey string    backendkit pubkey
      --backendkitSignURL string   backendkit sign url ex: "https://backendkit.com:4443/signature" (required)
      --backendkitToken string     backendkit token (required)
  -d, --directory string           source directory containing files to sign (required)
      --dryrun                     Print information about the files to sign without signing
  -e, --exitonerror                exit the app with an error code if anything goes wrong
  -h, --help                       help for sign
  -p, --private                    create signatues with non-public access
  -r, --recursive                  Explore subfolders
      --strict                     re-sign any file that has changed since last signature
      --strict-prune               same as --strict, plus delete the previous signature receipt
      --unsecureSSL                Do not check the ssl certificate validity for the backendkit (only use in developpement)

woleet-cli export [flags]
  -d, --directory string   directory where the receipts will be downloaded (required)
  -e, --exitonerror        exit the app with an error code if anything goes wrong
  -h, --help               help for export
  -l, --limitdate string   get all receipts generated from the provided date format:yyyy-MM-dd (default is no limit)

Global Flags:
  -c, --config string     config file (default is $HOME/.woleet-cli.yaml)
      --json              Switch output format to json
      --loglevel string   Select log level info|warn|error|fatal (default is info) (default "info")
  -t, --token string      JWT token (required)
  -u, --url string        custom API url (default "https://api.woleet.io/v1")
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
  strict-prune: true
  exitonerror: true
  recursive: true
  dryrun: false
sign:
  backendkitSignURL: https://backendkit.com:4443/signature
  backendkitToken: insert-your-backendkit-token-here
  unsecureSSL: false
export:
  directory: /home/folder/to/anchor
  limitdate: 2018-01-21
  exitonerror: true
log:
  json: true
  level: info
```

JSON:

```json
{
  "api": {
    "private": true,
    "url": "https://api.woleet.io/v1",
    "token": "insert-your-token-here"
  },
  "app": {
    "directory": "/home/folder/to/anchor",
    "exitonerror": true,
    "strict": true,
    "strict-prune": true,
    "recursive": true,
    "dryrun": true,
  },
  "sign": {
    "backendkitSignURL": "https://backendkit.com:4443/signature",
    "backendkitToken": "insert-your-backendkit-token-here",
    "unsecureSSL": false
  },
  "export": {
    "directory": "/home/folder/to/anchor",
    "limitdate": "2018-01-21",
    "exitonerror": true
  },
  "log": {
    "json": true,
    "level": "info"
  }
}
```

ENV:

```bash
export WLT_CONFIG="$HOME/.woleet-cli.json"
export WLT_API_URL="https://api.woleet.io/v1"
export WLT_API_TOKEN="insert-your-token-here"
export WLT_API_PRIVATE="true"
export WLT_APP_DIRECTORY="/home/folder/to/anchor"
export WLT_APP_EXITONERROR="true"
export WLT_APP_STRICT="true"
export WLT_APP_STRICT_PRUNE="true"
export WLT_APP_RECURSIVE="true"
export WLT_APP_DRYRUN="true"
export WLT_SIGN_BACKENDKITSIGNURL="https://backendkit.com:4443/signature"
export WLT_SIGN_BACKENDKITTOKEN="insert-your-backendkit-token-here"
export WLT_SIGN_UNSECURESSL="false"
export WLT_EXPORT_DIRECTORY="/home/folder/to/anchor"
export WLT_EXPORT_LIMITDATE="2018-01-21"
export WLT_EXPORT_EXITONERROR="true"
export WLT_LOG_JSON="true"
export WLT_LOG_LEVEL="info"
```

## Generate models from OpenAPI/Swagger specifications

The tool calls the Woleet API and the BackendKit API using model classes generated from their OpenAPI/Swagger specification.
If this specification were to be changed, model classes can be updated using the following commands:

```bash
# Update definition files
curl -s https://api.woleet.io/swagger.json > api/swagger.json
curl -s https://raw.githubusercontent.com/woleet/woleet-backendkit/master/swagger.yaml > api/swaggerBackendkit.yaml

# Update models
openapi-generator generate -i api/swagger.json -g go -o pkg/models/woleetapi -Dmodels -DmodelDocs=false -DpackageName=woleetapi --type-mappings boolean=*bool && \
ANCHOR_FILE=$(cat pkg/models/woleetapi/model_anchor.go) && echo "$ANCHOR_FILE" | sed 's/`json:"hash"`/`json:"hash,omitempty"`/' > pkg/models/woleetapi/model_anchor.go
openapi-generator generate -i api/swaggerBackendkit.yaml -g go -o pkg/models/backendkit -Dmodels -DmodelDocs=false -DpackageName=backendkit --type-mappings boolean=*bool
```

## Build the tool from sources

### Standard way

```bash
go get -u github.com/woleet/woleet-cli
# After this step the created binary will be in your $GOBIN folder, traditionnaly $GOPATH/bin
```

Complex way:

### For go < 1.11

```bash
# Clone this project in $GOPATH/src/github.com/woleet
# get mandatory libraries:
go get -u gopkg.in/resty.v1
go get -u github.com/spf13/cobra
go get -u github.com/spf13/viper
go get -u github.com/mitchellh/go-homedir
# For windows only:
go get -u github.com/inconshreveable/mousetrap

# Generating the actual binary
go build -o $<desired_path>/woleet-cli
```

### For go >= 1.11

```bash
# Clone this project wherever you want
# get libraries:
go get

# Generating the actual binary
go build -o $<desired_path>/woleet-cli
```

Alternatively, binaries for Linux and MacOS are available on github release tab
