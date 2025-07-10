package attacker

import (
	"encoding/hex"
	"errors"
	"github.com/sirupsen/logrus"
	attackclient "github.com/tsinghua-cel/attacker-client-go/client"
	"os"
	"strings"
	"sync"
)

var (
	serviceUrl string
	client     *attackclient.Client
	flagCache  sync.Map
)

func FromHex(s string) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	if len(s)%2 != 0 {
		return nil, errors.New("invalid hex string")
	}
	return hex.DecodeString(s)
}

func InitAttacker(url string) error {
	env := os.Getenv("ATTACKER_SERVICE_URL")
	if url != "" {
		env = url
	}
	serviceUrl = env
	logrus.WithField("url", serviceUrl).Info("Attacker service init")
	return nil
}

func GetAttacker() *attackclient.Client {
	if client != nil {
		return client
	}
	if serviceUrl == "" {
		return nil
	}

	c, err := attackclient.Dial(serviceUrl, 0)
	if err != nil {
		logrus.WithField("url", serviceUrl).WithError(err).Error("connect to attacker failed")
		return nil
	}
	client = c
	return client
}

func GetBoolFlag(key string) bool {
	if v, ok := flagCache.Load(key); ok {
		if b, ok := v.(bool); ok {
			return b
		}
		return false
	}
	return false
}

func SetFlag(key string, value interface{}) {
	flagCache.Store(key, value)
}
