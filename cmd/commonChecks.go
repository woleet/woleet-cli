package cmd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/woleet/woleet-cli/internal/app"
)

func checkToken(cmd *cobra.Command) string {
	if strings.EqualFold(viper.GetString("api.token"), "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a token")
	}
	return viper.GetString("api.token")
}

func checkExportDirectory(cmd *cobra.Command) string {
	viper.Set("app.directory", viper.GetString("export.directory"))
	return checkDirectory(cmd)
}

func checkDirectory(cmd *cobra.Command) string {
	if strings.EqualFold(viper.GetString("app.directory"), "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a directory")
	}

	absDirectory, errAbs := filepath.Abs(viper.GetString("app.directory"))

	if errAbs != nil {
		log.Fatalln("Unable to get Absolute directory from --directory")
	}

	info, err := os.Stat(absDirectory)
	if err != nil {
		log.Fatalln("The provided directory does not exists")
	} else {
		if !info.IsDir() {
			log.Fatalln("The provided path is not a directory")
		}
	}
	return absDirectory
}

func checkFilter(cmd *cobra.Command) *regexp.Regexp {
	if strings.EqualFold(viper.GetString("app.filter"), "") {
		return nil
	}
	filter, errFilter := regexp.Compile(viper.GetString("app.filter"))
	if errFilter != nil {
		log.Fatalf("Unable parse the regexp specified by the --filter: \n%s\n", errFilter)
	}
	return filter
}

func checkWidSignURL(cmd *cobra.Command) string {
	widsSignURL := ""
	if viper.IsSet("sign.widsSignURL") {
		widsSignURL = viper.GetString("sign.widsSignURL")
	}
	if viper.IsSet("seal.widsSignURL") {
		widsSignURL = viper.GetString("seal.widsSignURL")
	}
	if strings.EqualFold(widsSignURL, "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a widsSignURL")
	}
	return widsSignURL
}

func checkWidToken(cmd *cobra.Command) string {
	widsToken := ""
	if viper.IsSet("sign.widsToken") {
		widsToken = viper.GetString("sign.widsToken")
	}
	if viper.IsSet("seal.widsToken") {
		widsToken = viper.GetString("seal.widsToken")
	}
	if strings.EqualFold(widsToken, "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a widsToken")
	}
	return widsToken
}

func checkWidPubKey(cmd *cobra.Command) string {
	widsPubKey := ""
	if viper.IsSet("sign.widsPubKey") {
		widsPubKey = viper.GetString("sign.widsPubKey")
	}
	if viper.IsSet("seal.widsPubKey") {
		widsPubKey = viper.GetString("seal.widsPubKey")
	}
	if strings.EqualFold(widsPubKey, "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a widsPubKey")
	}
	return widsPubKey
}

func checkS3(cmd *cobra.Command) *minio.Client {
	if strings.EqualFold(viper.GetString("s3.bucket"), "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a bucket for your S3 connection")
	}

	if strings.EqualFold(viper.GetString("s3.endpoint"), "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a endpoint for your S3 connection")
	}

	if strings.EqualFold(viper.GetString("s3.accessKeyID"), "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a accessKeyID for your S3 connection")
	}

	if strings.EqualFold(viper.GetString("s3.secretAccessKey"), "") {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set a secretAccessKey for your S3 connection")
	}

	minioClient, errMinioClient := minio.New(viper.GetString("s3.endpoint"), viper.GetString("s3.accessKeyID"), viper.GetString("s3.secretAccessKey"), !viper.GetBool("s3.noSSL"))
	if errMinioClient != nil {
		log.Fatalln(errMinioClient)
	}

	_, errBucket := minioClient.BucketExists(viper.GetString("s3.bucket"))
	if errBucket != nil {
		log.Fatalln(errBucket)
	}

	return minioClient
}

func checkFolderType(cmd *cobra.Command, runParameters *app.RunParameters) {
	runParameters.IsFS = false
	runParameters.IsS3 = false
	if !strings.EqualFold(viper.GetString("app.directory"), "") {
		runParameters.IsFS = true
	}

	if !strings.EqualFold(viper.GetString("s3.bucket"), "") {
		runParameters.IsS3 = true
	}

	if runParameters.IsFS && runParameters.IsS3 {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("directory and bucket cannot be set simultaneously")
	}

	if !(runParameters.IsFS || runParameters.IsS3) {
		if !viper.GetBool("log.json") {
			cmd.Help()
		}
		log.Fatalln("Please set directory or bucket")
	}
}
