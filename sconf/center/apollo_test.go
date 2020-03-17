package center

import (
	"context"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func assertStringEqual(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		t.Errorf("%s and %s should be equal", s1, s2)
	}
}

func TestConfFromEnv(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		conf := confFromEnv()

		assertStringEqual(t, conf.CacheDir, defaultCacheDir)
		assertStringEqual(t, conf.Cluster, defaultCluster)
		assertStringEqual(t, conf.IP, defaultHostPort)
	})
}

func TestInit(t *testing.T) {
	ass := assert.New(t)
	center,err := NewConfigCenter(ApolloConfigCenter)
	ass.Nil(err)
	err = center.Init(context.TODO(), "base/servmonitor", []string{"rpc.client", "rpc.test111"})
	ass.Nil(err)

	time.Sleep(5*time.Minute)
}