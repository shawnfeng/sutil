package center

import (
	"testing"
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
