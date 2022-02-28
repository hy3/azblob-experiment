package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/go-autorest/autorest"
	"github.com/hy3/azblob-experiment/internal"
)

var config internal.Config

func main() {
	c, err := internal.InitConfig()
	if err != nil {
		internal.Logger.Fatal().Msg(err.Error())
	}
	config = c

	os.Exit(realMain())
}

func realMain() int {
	if err := internal.SetLogLevel(config.LogLevel); err != nil {
		internal.Logger.Error().Err(err).
			Msg("can't set log level")
		return 1
	}

	ctx := context.Background()
	logger := internal.Logger

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create azure credential")
		return 1
	}

	clientOptions := &azblob.ClientOptions{
		Retry: policy.RetryOptions{
			MaxRetries:    3,
			TryTimeout:    1 * time.Minute,
			RetryDelay:    500 * time.Millisecond,
			MaxRetryDelay: 500 * time.Millisecond,
			StatusCodes:   autorest.StatusCodesForRetry,
		},
	}
	serviceClient, err := azblob.NewServiceClient(
		fmt.Sprintf("https://%s.blob.core.windows.net", config.AccountName),
		cred,
		clientOptions,
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create service client")
		return 1
	}

	blobClient := serviceClient.
		NewContainerClient(config.ContainerName).
		NewBlobClient(config.BlobName)

	if len(os.Args) > 1 {
		versionID := os.Args[1]
		logger.Info().Msgf("use version id: %s", versionID)

		// NOTE:
		// `blobClient.WithVersionID()` does not work.
		// https://github.com/Azure/azure-sdk-for-go/issues/17188

		u, err := url.Parse(blobClient.URL())
		if err != nil {
			panic(err)
		}
		q := u.Query()
		q.Set("versionId", versionID)
		u.RawQuery = q.Encode()

		cli, err := azblob.NewBlobClient(u.String(), cred, clientOptions)
		if err != nil {
			logger.Error().Err(err).Msg("failed to re-create service client")
			return 1
		}
		blobClient = cli
	}

	resp, err := blobClient.Download(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to download data")
		return 1
	}
	defer resp.RawResponse.Body.Close()

	logger.Debug().Msg("***headers***")
	for key, vals := range resp.RawResponse.Header {
		for i, val := range vals {
			logger.Debug().Msgf("%s.%d: %s", key, i, val)
		}
	}

	body, err := io.ReadAll(resp.RawResponse.Body)
	if err != nil {
		logger.Error().Err(err).Msg("can't read response body")
		return 1
	}

	logger.Debug().Msg("***body***")
	logger.Debug().Msgf("%s", body)

	return 0
}
