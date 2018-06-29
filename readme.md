# woleet-cli

A command-line tool to (for now) anchors all files of a specific folder  

## Fonctionnalities

The tool will scan a folder, anchoring all the files available and creating a receipt named 'filename'-'anchorID'.recept.json next to the original file if this one is available  

The files that were already anchored will not be reanchored, to have this behavior, at the end of a post we will create a json which will be a subset of a chainpoint receipt with the name 'filename'-'anchorID'.pending.json:
```json
{
  "target_hash": "sha256here"
}
```

If the option --strict is passed and file and its receipt are present, we will check the hash in the receipt to see if it is the same as the file, in case of these differs, the new file will be anchored and the old recipt wille be kept, if --strict--prune is passe instead, ols receipts will be deleted.  

To sum up, the first call of this tool will anchors all the files and create pending jsons, and the second call will gather all the receipt for these files if they are available and if they are it wil deletes the pending.json file

Tags will be added to the anchors according the names of the subfolders  

By default the tool will use your current folder if you do not specify any folder on the cli  

## Limitations

- All files and folders begining by '.' or finished by '.<receipt|pending>.json' will be ignored by this tool  
- Symlinks will not be followed  
- Scanned subfolders cannot have a space in their name  
- The max subfolder path is 128 characters  

## Usage

```
Usage:
  woleet-cli anchor [flags]
Flags:
  -d, --directory       string  Source directory to read from (default ".")
      --strict          boolean Rescan the directory for file changes (default "false")
      --strict-prune    boolean Rescan the directory for file changes and delete old receipts (default "false")
      --token           string  JWT token
  -u  --url             string  url of the platform (default: "https://api.woleet.io/v1")
  -h, --help                    help for anchor (Display this)
```



## Update

```bash
# Update definition file
curl -s https://api.woleet.io/swagger.json > api/swagger.json

# Update models
swagger-codegen generate -i api/swagger.json -l go -o pkg/models -Dmodels -DmodelDocs=false -DpackageName=models
```

## Build from sources