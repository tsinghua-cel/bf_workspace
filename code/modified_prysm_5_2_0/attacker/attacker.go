package attacker

import (
	"github.com/sirupsen/logrus"
	attackclient "github.com/tsinghua-cel/attacker-client-go/client"
	"os"
	"sync"
)

var (
	initOnce   sync.Once
	serviceUrl string
	client     *attackclient.Client
)

func initAttacker() {
	env := os.Getenv("ATTACKER_SERVICE_URL")
	serviceUrl = env
	logrus.WithField("url", serviceUrl).Info("Attacker service init")
}

func GetAttacker() *attackclient.Client {
	initOnce.Do(initAttacker)
	if client != nil {
		return client
	}
	if serviceUrl == "" {
		return nil
	}

	c, err := attackclient.Dial(serviceUrl, 0)
	if err != nil {
		return nil
	}
	client = c
	return client
}
