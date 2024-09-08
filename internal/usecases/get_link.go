package usecases

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	model "github.com/yasamprom/cdn-balancer/internal/model"
)

func (u *usecases) GetLink(ctx context.Context, uri string) (string, error) {

	// Get percent of queries routed to original host
	boarder, err := u.getBoarer()
	if err != nil {
		return "", model.ErrBadRoutePercent
	}

	// Check if must redirect to cdn
	if !u.shouldRedirectToCDN(boarder) {
		return uri, nil
	}

	host := getEnv(viper.GetString("cdnHostEnv"), viper.GetString("cdnSourceHost"))
	link := constructCDNUri(uri, host)

	if err := checkAvailability(link); err != nil {
		return "", err
	}

	return link, nil
}

// Gets host from env or uses default
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}

// Constructs uri to cdn
func constructCDNUri(uri string, host string) string {
	matches := model.RegexpLink.FindStringSubmatch(uri)
	if len(matches) < 4 {
		return uri
	}
	domainParts := strings.Split(matches[2], ".")
	if len(domainParts) < 2 {
		return uri
	}
	originalHost := domainParts[0]

	return fmt.Sprintf("%s%s/%s%s", matches[1], host, originalHost, strings.Join(matches[3:], ""))
}

func (u *usecases) getBoarer() (uint32, error) {
	percent := viper.GetFloat64("originalSourceRoutePercent")
	if percent <= 0 {
		return 0, model.ErrBadRoutePercent
	}
	return uint32(1 / percent), nil
}

func (u *usecases) shouldRedirectToCDN(boarder uint32) bool {
	current := u.counter.Add(1)
	return current%boarder != 0
}

// Check if content is available
// TODO: implement
func checkAvailability(_ string) error {
	return nil
}
