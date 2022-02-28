package main

import (
	"context"
	"fmt"
	"io"
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

	serviceClient, err := azblob.NewServiceClient(
		fmt.Sprintf("https://%s.blob.core.windows.net", config.AccountName),
		cred,
		&azblob.ClientOptions{
			Retry: policy.RetryOptions{
				MaxRetries:    3,
				TryTimeout:    1 * time.Minute,
				RetryDelay:    500 * time.Millisecond,
				MaxRetryDelay: 500 * time.Millisecond,
				StatusCodes:   autorest.StatusCodesForRetry,
			},
		},
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create service client")
		return 1
	}

	blobClient := serviceClient.
		NewContainerClient(config.ContainerName).
		NewBlobClient(config.BlobName)

	resp, err := blobClient.Delete(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed to upload data")
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
