package container

import (
	mediaHttp "github.com/bagusyanuar/genpos-backend/internal/media/delivery/http"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
)

func (c *Container) wireMediaModule(conf *config.Config) {
	// Initialize the media handler with the shared uploader
	handler := mediaHttp.NewMediaHandler(c.Uploader, conf)

	c.MediaHandler = handler
}
