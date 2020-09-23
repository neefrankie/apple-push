package message

import (
	"github.com/neefrankie/apple-push/pkg/config"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
	"runtime"
	"sync"
	"time"
)

type ClientBuilder struct {
	token      *token.Token
	production bool
	topics     config.APNTopic
}

func NewClientBuilder(opt config.APNOption, prod bool) (ClientBuilder, error) {
	authKey, err := token.AuthKeyFromBytes([]byte(opt.AuthKey))
	if err != nil {
		return ClientBuilder{}, err
	}

	t := &token.Token{
		AuthKey: authKey,
		KeyID:   opt.KeyID,
		TeamID:  opt.TeamID,
	}

	_, err = t.Generate()
	if err != nil {
		return ClientBuilder{}, err
	}

	return ClientBuilder{
		token:      t,
		production: prod,
		topics:     config.MustAPNTopic(),
	}, nil
}

func (b ClientBuilder) CreateClient() *apns2.Client {
	b.token.GenerateIfExpired()
	client := apns2.NewTokenClient(b.token)

	if b.production {
		client.Production()
	}

	return client
}

func (b ClientBuilder) ComposeNotification(m *Message, d Device) *apns2.Notification {

	// Create a notification
	noti := &apns2.Notification{}

	switch d.DeviceType {
	case DeviceTypePhone:
		noti.Topic = b.topics.Phone

	case DeviceTypePad:
		noti.Topic = b.topics.Pad

	default:
		return nil
	}

	noti.DeviceToken = d.Token
	noti.Expiration = time.Now().Add(30 * time.Minute)

	noti.ApnsID = m.ApnsID
	noti.Payload = m.Payload()

	return noti
}

func (b ClientBuilder) CreatePool() chan *apns2.Client {
	clients := make(chan *apns2.Client, 50)

	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)

	for i := 0; i < runtime.NumCPU(); i++ {
		client := b.CreateClient()
		clients <- client
	}

	return clients
}

func (b ClientBuilder) Push(m *Message, devices chan Device) {

	clients := b.CreatePool()

	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < 100; i++ {
		client := <-clients
		clients <- client

		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			for device := range devices {
				noti := b.ComposeNotification(m, device)
				resp, err := client.Push(noti)
				if err != nil {
					return
				}

				// We do not handle status code above 410. Only 400, 410 deserve attention
				if resp.StatusCode == 200 {
					continue
				}

				if resp.StatusCode > 410 {
					continue
				}

				// For all those responses from APNs incdicating push failure, not all reasons mean the device token is invalid.
				// We only take the following reasons as indicating the device is no longer valid.
				switch resp.Reason {
				case apns2.ReasonBadDeviceToken:
				case apns2.ReasonDeviceTokenNotForTopic:
				case apns2.ReasonTopicDisallowed:
				case apns2.ReasonUnregistered:
					// Invalid device +1

				}
			}
		}(i)
	}

	wg.Wait()

	m.TimeElapsed = int(time.Since(start).Seconds())
}
