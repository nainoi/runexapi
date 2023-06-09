package firebase

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

//SendToToken send msg once token
func SendToToken(app *firebase.App, token string, msg string) {
	// [START send_to_token_golang]
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Println(err)
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"msg": msg,
		},
		Token: token,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	// [END send_to_token_golang]
}

func sendToTopic(ctx context.Context, client *messaging.Client) {
	// [START send_to_topic_golang]
	// The topic name can be optionally prefixed with "/topics/".
	topic := "highScores"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Topic: topic,
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	// [END send_to_topic_golang]
}

func sendToCondition(ctx context.Context, client *messaging.Client) {
	// [START send_to_condition_golang]
	// Define a condition which will send to devices which are subscribed
	// to either the Google stock or the tech industry topics.
	condition := "'stock-GOOG' in topics || 'industry-tech' in topics"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Condition: condition,
	}

	// Send a message to devices subscribed to the combination of topics
	// specified by the provided condition.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	// [END send_to_condition_golang]
}

func sendAll(ctx context.Context, client *messaging.Client) {
	// This registration token comes from the client FCM SDKs.
	registrationToken := "YOUR_REGISTRATION_TOKEN"

	// [START send_all]
	// Create a list containing up to 100 messages.
	messages := []*messaging.Message{
		{
			Notification: &messaging.Notification{
				Title: "Price drop",
				Body:  "5% off all electronics",
			},
			Token: registrationToken,
		},
		{
			Notification: &messaging.Notification{
				Title: "Price drop",
				Body:  "2% off all books",
			},
			Topic: "readers-club",
		},
	}

	br, err := client.SendAll(context.Background(), messages)
	if err != nil {
		log.Println(err)
	}

	// See the BatchResponse reference documentation
	// for the contents of response.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	// [END send_all]
}

func sendMulticast(ctx context.Context, client *messaging.Client) {
	// [START send_multicast]
	// Create a list containing up to 100 registration tokens.
	// This registration tokens come from the client FCM SDKs.
	registrationTokens := []string{
		"YOUR_REGISTRATION_TOKEN_1",
		// ...
		"YOUR_REGISTRATION_TOKEN_n",
	}
	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		
		Tokens: registrationTokens,
	}

	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Println(err)
	}

	// See the BatchResponse reference documentation
	// for the contents of response.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	// [END send_multicast]
}

// SendMulticastAndHandleErrors send multicast notification to multi device
func SendMulticastAndHandleErrors(ctx context.Context, client *messaging.Client, tokens []string, title string, body string) {
	// [START send_multicast_error]
	// Create a list containing up to 100 registration tokens.
	// This registration tokens come from the client FCM SDKs.
	oneHour := time.Duration(1) * time.Hour
	badge := 42
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Android: &messaging.AndroidConfig{
			TTL: &oneHour,
			Notification: &messaging.AndroidNotification{
				Icon:  "stock_ticker_update",
				Color: "#d45704",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Badge: &badge,
				},
			},
		},
		Tokens: tokens,
	}

	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Println(err)
	}

	if br.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range br.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, tokens[idx])
			}
		}

		fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
	}
	// [END send_multicast_error]
}

func sendDryRun(ctx context.Context, client *messaging.Client) {
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: "token",
	}

	// [START send_dry_run_golang]
	// Send a message in the dry run mode.
	response, err := client.SendDryRun(ctx, message)
	if err != nil {
		log.Println(err)
	}
	// Response is a message ID string.
	fmt.Println("Dry run successful:", response)
	// [END send_dry_run_golang]
}

func androidMessage(title string, body string) *messaging.Message {
	// [START android_message_golang]
	oneHour := time.Duration(1) * time.Hour
	message := &messaging.Message{
		Android: &messaging.AndroidConfig{
			TTL:      &oneHour,
			Priority: "normal",
			Notification: &messaging.AndroidNotification{
				Title: title,
				Body:  body,
				Icon:  "stock_ticker_update",
				Color: "#d45704",
			},
		},
		Topic: "runex",
	}
	// [END android_message_golang]
	return message
}

func apnsMessage(title string, body string) *messaging.Message {
	// [START apns_message_golang]
	badge := 42
	message := &messaging.Message{
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: title,
						Body:  body,
					},
					Badge: &badge,
				},
			},
		},
		Topic: "runex",
	}
	// [END apns_message_golang]
	return message
}

func webpushMessage(title string, body string) *messaging.Message {
	// [START webpush_message_golang]
	message := &messaging.Message{
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: title,
				Body:  body,
				Icon:  "https://runex.co/icon.png",
			},
		},
		Topic: "runex",
	}
	// [END webpush_message_golang]
	return message
}

func allPlatformsMessage(title string, body string) *messaging.Message {
	// [START multi_platforms_message_golang]
	oneHour := time.Duration(1) * time.Hour
	badge := 42
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Android: &messaging.AndroidConfig{
			TTL: &oneHour,
			Notification: &messaging.AndroidNotification{
				Icon:  "stock_ticker_update",
				Color: "#d45704",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Badge: &badge,
				},
			},
		},
		Topic: "runex",
	}
	// [END multi_platforms_message_golang]
	return message
}

func subscribeToTopic(ctx context.Context, client *messaging.Client) {
	topic := "highScores"

	// [START subscribe_golang]
	// These registration tokens come from the client FCM SDKs.
	registrationTokens := []string{
		"YOUR_REGISTRATION_TOKEN_1",
		// ...
		"YOUR_REGISTRATION_TOKEN_n",
	}

	// Subscribe the devices corresponding to the registration tokens to the
	// topic.
	response, err := client.SubscribeToTopic(ctx, registrationTokens, topic)
	if err != nil {
		log.Println(err)
	}
	// See the TopicManagementResponse reference documentation
	// for the contents of response.
	fmt.Println(response.SuccessCount, "tokens were subscribed successfully")
	// [END subscribe_golang]
}

func unsubscribeFromTopic(ctx context.Context, client *messaging.Client) {
	topic := "highScores"

	// [START unsubscribe_golang]
	// These registration tokens come from the client FCM SDKs.
	registrationTokens := []string{
		"YOUR_REGISTRATION_TOKEN_1",
		// ...
		"YOUR_REGISTRATION_TOKEN_n",
	}

	// Unsubscribe the devices corresponding to the registration tokens from
	// the topic.
	response, err := client.UnsubscribeFromTopic(ctx, registrationTokens, topic)
	if err != nil {
		log.Println(err)
	}
	// See the TopicManagementResponse reference documentation
	// for the contents of response.
	fmt.Println(response.SuccessCount, "tokens were unsubscribed successfully")
	// [END unsubscribe_golang]
}
