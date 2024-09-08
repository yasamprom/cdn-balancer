package usecases

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type uriTest struct {
	host      string
	incoming  string
	outcoming string
}

func TestConstructCDNUri(t *testing.T) {

	testCases := []uriTest{
		{
			"host.com",
			"https://s1.origin-cluster/video/123/xcg2djHckad.m3u8",
			"https://host.com/s1/video/123/xcg2djHckad.m3u8",
		},
		{
			"host.com",
			"https://s23.complex.cluster.addres/video/123/abcab.m3u8",
			"https://host.com/s23/video/123/abcab.m3u8",
		},
		{
			"host.subdomain.com",
			"https://s1.complex.cluster.addres/video/123/xcg2djHckad.other_extansion",
			"https://host.subdomain.com/s1/video/123/xcg2djHckad.other_extansion",
		},
		{
			"host.subdomain.com",
			"http://s1.complex.cluster.addres/video/123/xcg2djHckad.m3u8",
			"http://host.subdomain.com/s1/video/123/xcg2djHckad.m3u8",
		},
	}

	for _, tc := range testCases {
		res := constructCDNUri(tc.incoming, tc.host)
		assert.Equal(t, res, tc.outcoming)
	}

}

func TestGetLink(t *testing.T) {
	uc := New()
	os.Setenv("CDN_HOST", "cdnhost")
	ctx := context.Background()
	testLink := "http://s1.complex.cluster.addres/video/123/xcg2djHckad.m3u8"
	cdnLink := "http://cdnhost.com/s1/video/123/xcg2djHckad.m3u8"

	viper.Set("originalSourceRoutePercent", 0.1)
	viper.Set("cdnSourceHost", "cdnhost.com")
	t.Run("simple case", func(t *testing.T) {
		b, err := uc.getBoarer()
		assert.NoError(t, err)
		assert.True(t, b > 0)

		for range b - 1 {
			uri, err := uc.GetLink(ctx, testLink)
			assert.NoError(t, err)
			assert.Equal(t, uri, cdnLink)
		}
		uri, err := uc.GetLink(ctx, testLink)
		assert.NoError(t, err)
		assert.Equal(t, uri, testLink)
	})

	t.Run("concurrent case", func(t *testing.T) {
		b, err := uc.getBoarer()
		assert.NoError(t, err)
		assert.True(t, b > 0)

		runs := 1000
		originalRoutes := atomic.Uint32{}
		wg := sync.WaitGroup{}
		for range runs {
			wg.Add(1)
			go func() {
				uri, err := uc.GetLink(ctx, testLink)
				assert.NoError(t, err)
				if uri == testLink {
					originalRoutes.Add(1)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		// check original host redirect count
		assert.Equal(t, uint32(runs)/b, originalRoutes.Load())
	})

}
