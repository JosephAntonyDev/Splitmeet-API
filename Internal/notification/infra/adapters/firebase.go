package adapters

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/JosephAntonyDev/splitmeet-api/internal/core"
	"google.golang.org/api/option"
)

type FirebasePushSender struct {
	client *messaging.Client
}

func NewFirebasePushSender(ctx context.Context) (*FirebasePushSender, error) {
	credentialsFile := os.Getenv("FIREBASE_CREDENTIALS_FILE")
	if credentialsFile == "" {
		return nil, fmt.Errorf("FIREBASE_CREDENTIALS_FILE no está configurada")
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("error inicializando firebase: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creando cliente de mensajería: %v", err)
	}

	return &FirebasePushSender{client: client}, nil
}

func (s *FirebasePushSender) SendAndroidPush(ctx context.Context, req core.PushRequest) (bool, error) {
	msg := &messaging.Message{
		Token: req.Token,
		Data:  req.Data,
		Notification: &messaging.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID:             "splitmeet_alerts_high",
				Sound:                 "default",
				DefaultVibrateTimings: true,
				DefaultSound:          true,
				Visibility:            messaging.VisibilityPublic,
				Priority:              messaging.PriorityMax,
				Title:                 req.Title,
				Body:                  req.Body,
				Tag:                   "splitmeet_notification",
			},
		},
	}

	_, err := s.client.Send(ctx, msg)
	if err != nil {
		if messaging.IsUnregistered(err) || messaging.IsInvalidArgument(err) {
			return true, err
		}
		return false, err
	}

	return false, nil
}
