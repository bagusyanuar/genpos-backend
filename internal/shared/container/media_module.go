package container

import (
	mediaHttp "github.com/bagusyanuar/genpos-backend/internal/media/delivery/http"
	"github.com/bagusyanuar/genpos-backend/internal/shared/config"
	"github.com/bagusyanuar/genpos-backend/pkg/fileupload"
	"gorm.io/gorm"
)

func (c *Container) wireMediaModule(db *gorm.DB, conf *config.Config) {
	// Initialize the generic file uploader utility
	uploader := fileupload.NewLocalFileUploader("./public/uploads", "/public/uploads")

	// Initialize the media handler with the shared uploader
	handler := mediaHttp.NewMediaHandler(uploader, conf)

	c.MediaHandler = handler
}
