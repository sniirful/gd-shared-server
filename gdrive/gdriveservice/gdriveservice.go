package gdriveservice

import (
	"app/oauth"
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	driveService *drive.Service = nil
)

func GetService() (*drive.Service, error) {
	if driveService != nil {
		return driveService, nil
	}

	client, err := oauth.GetClient(func(url string) {
		fmt.Printf("Grant permission to your account using the following link:\n%v\n", url)
	}, drive.DriveScope)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return driveService, nil
}
