package firebase_service

import (
	"context"
	"encoding/json"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"git.samberi.com/dois/delivery_api/config"
	"git.samberi.com/dois/delivery_api/logger"
	"google.golang.org/api/option"
)

var App *firebase.App

func FirebaseApp(ctx context.Context) (*firebase.App, error) {
	var err error
	if App == nil {
		opt := option.WithCredentialsFile(config.Config.Firebase.CertPath)
		App, err = firebase.NewApp(ctx, nil, opt)
		if err != nil {
			logger.Logger.Error(err)
		} else {
			_, err := App.Auth(ctx)
			logger.HandleError(err)
		}
	}
	return App, err
}

func FirebaseSendNotification(data map[string]string, token, topic string) (err error) {
	ctx := context.Background()
	app, err := FirebaseApp(ctx)
	if err == nil {
		messageClient, err := app.Messaging(ctx)
		if err == nil {
			message := messaging.Message{
				Data: data,
			}
			if token != "" {
				message.Token = token
			} else {
				message.Topic = topic
			}
			ans, err := messageClient.Send(ctx, &message)
			logger.Logger.Info("firebase send message", ans)
			logger.HandleError(err)
		}
	}
	return err
}

func SendEvent(event string, message map[string]interface{}, token, topic string) error {
	data, _ := json.Marshal(message)
	return FirebaseSendNotification(map[string]string{
		"event":   "test",
		"type":    "EVENT",
		"message": string(data),
	}, token, topic)
}
