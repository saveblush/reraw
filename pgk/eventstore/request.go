package eventstore

import "github.com/saveblush/reraw/models"

type Request struct {
	NostrFilter *models.Filter
	DoCount     bool
	NoLimit     bool
}
