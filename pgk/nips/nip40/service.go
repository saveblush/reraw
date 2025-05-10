package nip40

import (
	"errors"
	"strconv"

	"github.com/saveblush/reraw/core/cctx"
	"github.com/saveblush/reraw/core/config"
	"github.com/saveblush/reraw/core/utils"
	"github.com/saveblush/reraw/models"
)

// Service service interface
type Service interface {
	Expiration(c *cctx.Context, evt *models.Event) (*models.Timestamp, error)
}

type service struct {
	config *config.Configs
}

func NewService() Service {
	return &service{
		config: config.CF,
	}
}

// Expiration expiration
func (s *service) Expiration(c *cctx.Context, evt *models.Event) (*models.Timestamp, error) {
	var expiration int64
	expirationTag := evt.Tags.FindFirst("expiration")
	if expirationTag != nil && len(*expirationTag) >= 2 {
		expiration, _ = strconv.ParseInt(expirationTag.Value(), 10, 64)
		if expiration < 100 {
			return utils.Pointer(models.Timestamp(expiration)), errors.New("invalid: expiration")
		}
	}

	return utils.Pointer(models.Timestamp(expiration)), nil
}
