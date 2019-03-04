# Build the tool from sources

## Standard way

```bash
go get -u github.com/woleet/woleet-cli
# After this step the created binary will be in your $GOBIN folder, traditionnaly $GOPATH/bin
```

### For go >= 1.11

```bash
# Clone this project wherever you want
git clone git@github.com:woleet/woleet-cli.git

# Generating the actual binary
go build -o $<desired_path>/woleet-cli

# or

# The created binary will be in your $GOBIN folder
go install
```

### For go < 1.11

```bash
# Clone this project in $GOPATH/src/github.com/woleet
# get mandatory libraries:
go get -u gopkg.in/resty.v1
go get -u github.com/spf13/cobra
go get -u github.com/spf13/viper
go get -u github.com/mitchellh/go-homedir
go get -u github.com/kennygrant/sanitize
go get -u github.com/sirupsen/logrus
# For windows only:
go get -u github.com/inconshreveable/mousetrap

# Generating the actual binary
go build -o $<desired_path>/woleet-cli
```

## Generate models from OpenAPI/Swagger specifications

The tool calls Woleet API and Woleet.ID Server API using model classes generated from their OpenAPI/Swagger specification.
If this specification were to be changed, model classes can be updated using the following commands:

```bash
# Update definition files
curl -s https://api.woleet.io/swagger.json > api/swagger.json
curl -s https://raw.githubusercontent.com/woleet/woleet.id-server/master/swagger.yaml > api/swaggerIDServer.yaml


# Update models
openapi-generator generate -i api/swagger.json -g go -o pkg/models/woleetapi -Dmodels -DmodelDocs=false -DpackageName=woleetapi --type-mappings boolean=*bool && \
ANCHOR_FILE=$(cat pkg/models/woleetapi/model_anchor.go) && echo "$ANCHOR_FILE" | sed 's/`json:"hash"`/`json:"hash,omitempty"`/' > pkg/models/woleetapi/model_anchor.go
openapi-generator generate -i api/swaggerIDServer.yaml -g go -o pkg/models/idserver -Dmodels -DmodelDocs=false -DpackageName=idserver --type-mappings boolean=*bool && rm pkg/models/idserver/model_api_* pkg/models/idserver/model_key_* pkg/models/idserver/model_mnemonics* pkg/models/idserver/model_server_* pkg/models/idserver/model_user_*
```