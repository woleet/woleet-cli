# Woleet command line interface

woleet-cli is an open source command line tool on top of the Woleet API.
The tool is written in Go and has been tested on Windows, macOS and Linux.

Currently, the tool only supports:
 * the `anchor` command, allowing to recursively anchor all files in a given directory
 * the `sign` command, allowing to recursively sign all files in a given directory (using the backend kit: <https://github.com/woleet/woleet-backendkit>)

## Functionalities

The tool scans a folder recursively and anchors or sign all files found. It also gathers proof receipts and stores them beside anchored or signed files (in a Chainpoint file named 'filename'-'anchorID'.(signature-)?receipt.json).

Since anchoring is not a realtime operation, the tool is supposed to be run on a regular basis (or at least a second time when all proof receipts are ready to download). Obviously, the files that were already anchored are not re-anchored.

If the option --strict is provided, for each file that already have a proof receipt, the tool checks that the hash of the file still matches the hash in the receipt (to detect file changes). If they differ, the file is re-anchored and the old receipt is kept, except if --strict--prune is used instead.

To sum up, this tool can be used to generate and maintain the set of timestamped proofs of existence for all files in a given directory.

Note: tags are added to the anchors according to the name of sub-folders  

The tool behavior can be configured using command line arguments, environment variables or a configuration file. When several configuration means are used, the following priorities are applied:
- command line arguments
- environment variables
- config file
- default value (if any)

## Limitations

- All files and folders beginning by '.' or finished by '.(signature-)?receipt|pending.json' are ignored
- Symlinks are not followed  
- Scanned sub-folders cannot have a space in their name  
- The maximum length of the subfolder path (without delimiters) is 128 characters  

## Usage

```
Usage:
woleet-cli anchor [flags]
  -d, --directory string   source directory containing files to anchor (required)
  -e, --exitonerror        exit the app with an error code if something goes wrong
  -h, --help               help for anchor
  -p, --private            create anchors with non-public access
      --strict             re-anchor any file that has changed since last anchoring
      --strict-prune       same as --strict, plus delete the previous anchoring receipt

woleet-cli sign [flags]
      --backendkitPubKey string    backendkit pubkey
      --backendkitSignURL string   backendkit sign url ex: "https://backendkit.com:4443/signature" (required)
      --backendkitToken string     backendkit token (required)
  -d, --directory string           source directory containing files to sign (required)
  -e, --exitonerror                exit the app with an error code if something goes wrong
  -h, --help                       help for sign
  -p, --private                    create signatues with non-public access
      --strict                     re-sign any file that has changed since last signature
      --strict-prune               same as --strict, plus delete the previous signature receipt
      --unsecureSSL                Do not check the ssl certificate validity for the backendkit (only use in developpement)

Global Flags:
  -c, --config string   config file (default is $HOME/.woleet-cli.yaml)
  -t, --token string    JWT token (required)
  -u, --url string      custom API url (default "https://api.woleet.io/v1")
```

## Configuration file format

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
sign:
  backendkitSignURL: https://backendkit.com:4443/signature
  backendkitToken: insert-your-backendkit-token-here
  unsecureSSL: false
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
    "strict-prune": true
  },
  "sign": {
    "backendkitSignURL": "https://backendkit.com:4443/signature",
    "backendkitToken": "insert-your-backendkit-token-here",
    "unsecureSSL": false
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
export WLT_SIGN_BACKENDKITSIGNURL="https://backendkit.com:4443/signature"
export WLT_SIGN_BACKENDKITTOKEN="insert-your-backendkit-token-here"
export WLT_SIGN_UNSECURESSL="false"
```

## Generate models from OpenAPI/Swagger specifications

The tool calls the Woleet API and the BackendKit API using model classes generated from their OpenAPI/Swagger specification.
If this specification were to be changed, model classes can be updated using the following commands:

```bash
# Update definition files
curl -s https://api.woleet.io/swagger.json > api/swagger.json
curl -s https://raw.githubusercontent.com/woleet/woleet-backendkit/master/swagger.yaml > api/swaggerBackendkit.yaml

# Update models
swagger-codegen generate -i api/swagger.json -l go -o pkg/models -Dmodels -DmodelDocs=false -DpackageName=models --type-mappings boolean=*bool && \
ANCHOR_FILE=$(cat pkg/models/anchor.go) && echo "$ANCHOR_FILE" | sed 's/`json:"hash"`/`json:"hash,omitempty"`/' > pkg/models/anchor.go
swagger-codegen generate -i api/swaggerBackendkit.yaml -l go -o pkg/modelsBackendkit -Dmodels -DmodelDocs=false -DpackageName=modelsBackendkit --type-mappings boolean=*bool
```

## Build the tool from sources

Standard way:

```bash
go get -u github.com/woleet/woleet-cli
# After this step the created binary will be in your $GOBIN folder, traditionnaly $GOPATH/bin
```

Complex way:

```bash
# Clone this project in $GOPATH/src/github.com/woleet
# get mandatory libraries:
go get -u gopkg.in/resty.v1
go get -u github.com/spf13/cobra
go get -u github.com/spf13/viper
go get -u github.com/mitchellh/go-homedir
# For windows only, untested:
go get -u github.com/inconshreveable/mousetrap

# Generating the actual binary
go build -o $<desired_path>/woleet-cli
```

Alternatively, binaries for Linux and MacOS are available on github release tab
