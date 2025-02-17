package pixelita

import (
	"github.com/keshon/pixelita/internal/config"
)

type Pixelita struct {
	config *config.Config
}

func NewPixelita(config *config.Config) *Pixelita {
	return &Pixelita{
		config: config,
	}
}
