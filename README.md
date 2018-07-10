# woleet-cli

Woleet command line interface.
Currently, the tool only supports the 'anchor' command, which allows anchoring all files of a given directory recursively.

## Fonctionnalities

The tool scans a folder recursively and anchors all files available. It also gathers proof receipts and store them beside anchored files (in Chainpoint files named 'filename'-'anchorID'.receipt.json).

Since anchoring is not a realtime operation, the tool is supposed to be run on a regular basis (or at least until all proof receipts are ready to download). Obviously, the 
 files that were already anchored are not re-anchored.

If the option --strict is provided, for each file that already have a proof receipt, the tool checks that the hash of the file still matches the hash in the receipt (to detect file changes). If they differ, the file is re-anchored and the old receipt is kept, except if --strict--prune is used instead.

To sum up, this tool can be used to generate and maintain the set of timestamped proofs of existence for all files of a given directory.

Note: tags are added to the anchors according to the name of subfolders  

## Limitations

- All files and folders begining by '.' or finished by '.receipt|pending.json' are ignored
- Symlinks are not followed  
- Scanned subfolders cannot have a space in their name  
- The maximum length of the subfolder path (without delimiters) is 128 characters  

## Usage

```
Usage:
  woleet-cli anchor [flags]

Flags:
  -d, --directory string   source directory containing files to anchor
  -e, --exitonerror        exit the app with an error code if something goes wrong
  -h, --help               help for anchor
  -p, --private            create anchors as private anchors (not discoverable by the file hash)
      --strict             re-anchor any file that has changed since last anchoring
      --strict-prune       same as --strict, plus delete the previous anchoring receipt

Global Flags:
  -t, --token string   JWT token (required)
  -u, --url string     custom API url (default: "https://api.woleet.io/v1")
```

## Update models

```bash
# Update definition file
curl -s https://api.woleet.io/swagger.json > api/swagger.json

# Update models
swagger-codegen generate -i api/swagger.json -l go -o pkg/models -Dmodels -DmodelDocs=false -DpackageName=models
```

## Build from sources

```bash
# You will need 2 go libraries:
go get -u github.com/spf13/cobra
go get -u gopkg.in/resty.v1

# Once done at the root of this folder type:
go install
```
