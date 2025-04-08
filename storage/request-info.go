package storage

import "io"

type UploadRequest struct {
	FilePath         string
	File             io.Reader
	ContentType      string
	DownloadFileName *string
}
