# Woleet Command Line Interface

woleet-cli is an open source command line tool on top of the Woleet API.
The tool is written in Go and has been tested on Windows, macOS and Linux.

Currently, the tool only supports:

* the `timestamp` command, allowing to recursively timestamp all files in a given directory (legacy `anchor` command)
* the `seal` command, allowing to recursively seal all files in a given directory (using Woleet.ID Server: <https://github.com/woleet/woleet.id-server>) (legacy `sign` command)
* the `export` command, allowing to download all your proof receipts in a given directory

## Timestamp / Seal

### Functionalities

The tool scans a folder recursively and timestamps or seal all files found. It also gathers proof receipts and stores them beside timestamped or sealed files (in a Chainpoint file named 'filename'-'proofID'.(timestamp|seal)-receipt.json).

Since timestamping is not a realtime operation, the tool is supposed to be run on a regular basis (or at least a second time when all proof receipts are ready to download). Obviously, the files that were already timestamp are not re-timestamped.

If the option --strict is provided, for each file that already have a proof receipt, the tool checks that the hash of the file still matches the hash in the receipt (to detect file changes), in addition when sealing the public key is checked as well. If they differ, the file is re-timestamped and the old receipt is kept, except if --prune is set in that case the old receipt is deleted.  
If the original file is no longer present and the option --prune is provided, the old receipt/pending file will be deleted.

If you want to timestamp a subset of the files present in a folder or a subfolder, you can use the --filter option which will limit the scope of this tool to the files that matches the provided regex, you can test the regex here: <https://regex-golang.appspot.com/assets/html/index.html>, for example.

To sum up, this tool can be used to generate and maintain the set of proofs of timestamp or proof of seal for all the files of a set of directories.

### S3 support

When filling --s3AccessKeyID, --s3SecretAccessKey, --s3Bucket and --s3Endpoint you will not have to specify --directory.  

In that configuration, woleet-cli will timestamp/seal all files in the input bucket (regex still works), that process can be long because files will be downloaded to calculate their hashes.  
Proof receipts and pending files will be stored along original files in the S3 bucket.  

When using an S3-like directory, we advise to not use the --strict parameter as it will download all files at each run.  

### Limitations

* All files and folders beginning by '.' or finished by '.(timestamp|seal)-(receipt|pending).json' are ignored
* Symlinks are not followed

## Export

### Functionalities

The tool dumps all your proof receipts into a folder.  
You can specify a limit date to get all receipt created from this date.  

### Limitations

* Each receipt will be named: 'timestamp name'-'proofID'.(timestamp|seal)-receipt.json
* Exporting can be quite long, as each receipt is downloaded individually

## Install woleet-cli

### Get latest binaries

The latest binaries can be found [here](https://github.com/woleet/woleet-cli/releases)

Just download it, decompress it and execute it (grant execution permissions if necessary)

#### For Linux or MacOS

You can use this command to install the latest binaries in /usr/local/bin

```bash
CLI_URL=$(curl --silent https://api.github.com/repos/woleet/woleet-cli/releases/latest | grep 'browser_download_url' | grep -ioE "https://.*$(uname -s)_x86_64.tar.gz") && \
sudo curl -L "$CLI_URL" | sudo tar -xz -C /usr/local/bin woleet-cli && \
sudo chmod +x /usr/local/bin/woleet-cli
```

## Configure woleet-cli

The tool behavior can be configured using command line arguments, environment variables or a configuration file. When several configuration means are used, the following priorities are applied:

* command line arguments
* environment variables
* config file
* default value (if any)

There is also a special environment variable or config path to disable environment configuration and configuration file:

```bash
woleet-cli --config "DISABLED" ...
# or
export WCLI_CONFIG="DISABLED"
```

### Usage

```
Usage:
  woleet-cli timestamp [flags]

Aliases:
  timestamp, anchor

Flags:
  -d, --directory string           source directory containing files to timestamp (required)
      --dryRun                     print information about files to timetamp without timetamping
  -e, --exitOnError                exit with an error code if anything goes wrong
  -f, --filter string              timestamp only files matching this regex
      --fixReceipts                Check the format and fix (if necessary) every existing receipts,
                                   also rename legacy receipts ending by signature-receipt.json to seal-receipt.json
  -h, --help                       help for timestamp
  -p, --private                    create non discoverable proofs
      --prune                      delete receipts that are not along the original file,
                                   with --strict it checks the hash of the original file and deletes the receipt if they do not match
  -r, --recursive                  explore sub-folders recursively
      --s3AccessKeyID string       your AccessKeyID
      --s3Bucket string            bucket name that contains files to timestamp
      --s3Endpoint string          specify an alternative S3 endpoint: ex: storage.googleapis.com,
                                   don't specify the transport (https://), https will be used by default if you want to use http see --s3NoSSL param (default "s3.amazonaws.com")
      --s3NoSSL                    use S3 without SSL (strongly discouraged)
      --s3SecretAccessKey string   your SecretAccessKey
      --strict                     re-timetamp any file that has changed since last timetamping


Usage:
  woleet-cli seal [flags]

Aliases:
  seal, sign

Flags:
  -d, --directory string           source directory containing files to seal (required)
      --dryRun                     print information about files to seal without sealing
  -e, --exitOnError                exit with an error code if anything goes wrong
  -i, --filter string              seal only files matching this regex
      --fixReceipts                Check the format and fix (if necessary) every existing receipts,
                                   also rename legacy receipts ending by signature-receipt.json to seal-receipt.json
  -h, --help                       help for seal
  -p, --private                    create non discoverable proofs
      --prune                      delete receipts that are not along the original file,
                                   with --strict it checks the hash of the original file and deletes the receipt if they do not match or if the pubkey has changed
  -r, --recursive                  explore sub-folders recursively
      --s3AccessKeyID string       your AccessKeyID
      --s3Bucket string            bucket name that contains files to seal
      --s3Endpoint string          specify an alternative S3 endpoint: ex: storage.googleapis.com,
                                   don't specify the transport (https://), https will be used by default if you want to use http see --s3NoSSL param (default "s3.amazonaws.com")
      --s3NoSSL                    use S3 without SSL (strongly discouraged)
      --s3SecretAccessKey string   your SecretAccessKey
      --strict                     re-seal any file that has changed since last sealing or if the pubkey was changed
      --widsPubKey string          public key (ie. bitcoin address) to use to seal (required)
      --widsSignURL string         Woleet.ID Server sign URL ex: "https://idserver.com:3002" (required)
      --widsToken string           Woleet.ID Server API token (required)
      --widsUnsecureSSL            do not check Woleet.ID Server's SSL certificate validity (only for development)

Usage:
  woleet-cli export [flags]

Flags:
  -d, --directory string   directory where to store the proofs (required)
  -e, --exitOnError        exit with an error code if anything goes wrong
      --fixReceipts        Rename legacy receipts ending by anchor/signature-receipt.json to timestamp/seal-receipt.json
  -h, --help               help for export
  -l, --limitDate string   get only proofs created after the provided date (format: yyyy-MM-dd)

Global Flags:
  -c, --config string     config file (default is $HOME/.woleet-cli.yaml)
  -h, --help              display help for woleet-cli
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
  directory: /home/folder/to/proof
  filter: '.*\.json'
  fixReceipts: true
  strict: true
  prune: true
  exitOnError: true
  recursive: true
  dryRun: false
seal:
  widsSignURL: https://idserver.com:3002
  widsToken: insert-your-idserver-token-here
  widsPubKey: insert-your-idserver-pubkey-here
  widsUnsecureSSL: false
s3:
  bucket: bucket-name
  endpoint: storage.googleapis.com
  accessKeyID: insert-your-accessKeyID-here
  secretAccessKey: insert-your-secretAccessKey-here
  noSSL: true
export:
  directory: /home/folder/to/proof
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
    "private": true
  },
  "app": {
    "directory": "/home/folder/to/proof",
    "filter": ".*\.json",
    "fixReceipts": true,
    "exitOnError": true,
    "strict": true,
    "prune": true,
    "recursive": true,
    "dryRun": true
  },
  "seal": {
    "widsSignURL": "https://idserver.com:3002",
    "widsToken": "insert-your-idserver-token-here",
    "widsPubKey": "insert-your-idserver-pubkey-here",
    "widsUnsecureSSL": false
  },
  "s3": {
    "bucket": "bucket-name",
    "endpoint": "storage.googleapis.com",
    "accessKeyID": "insert-your-accessKeyID-here",
    "secretAccessKey": "insert-your-secretAccessKey-here",
    "noSSL": true
  },
  "export": {
    "directory": "/home/folder/to/proof",
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
export WCLI_APP_DIRECTORY="/home/folder/to/proof"
export WCLI_APP_FILTER='.*\.json'
export WCLI_APP_FIXRECEIPTS="true"
export WCLI_APP_EXITONERROR="true"
export WCLI_APP_STRICT="true"
export WCLI_APP_PRUNE="true"
export WCLI_APP_RECURSIVE="true"
export WCLI_APP_DRYRUN="true"
export WCLI_SEAL_WIDSSIGNURL="https://idserver.com:3002"
export WCLI_SEAL_WIDSTOKEN="insert-your-idserver-token-here"
export WCLI_SEAL_WIDSPUBKEY="insert-your-idserver-pubkey-here"
export WCLI_SEAL_WIDSUNSECURESSL="false"
export S3_BUCKET="bucket-name"
export S3_ENDPOINT="storage.googleapis.com"
export S3_ACCESSKEYID="insert-your-accessKeyID-here"
export S3_SECRETACCESSKEY="insert-your-secretAccessKey-here"
export S3_NOSSL="true"
export WCLI_EXPORT_DIRECTORY="/home/folder/to/proof"
export WCLI_EXPORT_LIMITDATE="2018-01-21"
export WCLI_EXPORT_EXITONERROR="true"
export WCLI_LOG_JSON="true"
export WCLI_LOG_LEVEL="info"
```

## Build woleet-cli

### Clone and build in your $GOPATH folder

```bash
$ GO111MODULE=on go get github.com/woleet/woleet-cli
```

After this step the created binary will be installed in your $GOBIN folder (traditionally $GOPATH/bin).

### Clone and build in any folder

```bash
$ git clone git@github.com:woleet/woleet-cli.git
$ cd woleet-cli
$ go build
```

After this step you can install the binary in your $GOBIN folder by doing:

```bash
$ go install
```

### Generate models from OpenAPI/Swagger specifications

The tool calls Woleet API and Woleet.ID Server API using model classes generated from their OpenAPI/Swagger specification.
If this specification were to be changed, model classes can be updated using the following commands:

```bash
# Update definition files
curl -s https://api.woleet.io/v1/openapi.json > api/swagger.json
curl -s https://raw.githubusercontent.com/woleet/woleet.id-server/master/swagger.yaml > api/swaggerIDServer.yaml


# Update models
rm -rf pkg/models/woleetapi pkg/models/idserver && \
JAVA_TOOL_OPTIONS='-Dmodels=anchor,anchors -DmodelDocs=false -DmodelTests=false' openapi-generator generate -i api/swagger.json -g go -o pkg/models/woleetapi -p packageName=woleetapi -p enumClassPrefix=true -p generateAliasAsModel=false && \
JAVA_TOOL_OPTIONS='-Dmodels=UserModeEnum,UserStatusEnum,UserRoleEnum,KeyStatusEnum,KeyTypeEnum,KeyHolderEnum,KeyDeviceEnum,SignatureResult,UserDisco,KeyGet,FullIdentity,ConfigDisco -DmodelDocs=false -DmodelTests=false' openapi-generator generate -i api/swaggerIDServer.yaml -g go -o pkg/models/idserver -p packageName=idserver -p enumClassPrefix=true -p generateAliasAsModel=false
```
