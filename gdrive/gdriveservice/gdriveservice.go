package gdriveservice

import (
	"app/oauth"
	"app/screen"
	"context"

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
		screen.Println("Grant permission to your account using the following link:")
		screen.Println("%v", url)
	}, func() {
		screen.Println("Reconnecting to Google Drive...")
	}, drive.DriveScope)
	if err != nil {
		screen.Println("%v", err)
		return nil, err
	}

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return driveService, nil
}
