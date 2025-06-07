package allure

import (
	"bytes"
	"cmp"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Attachment interface {
	// Open attachment for reading.
	Open() (io.ReadCloser, error)

	// ID is the unique ID of this attachment.
	ID() uuid.UUID

	// Type is the media type of the content.
	Type() string
}

func NewAttachmentBytes(data []byte, mediaType string) Attachment {
	return attachmentBytes{
		id:        uuid.New(),
		data:      data,
		mediaType: mediaType,
	}
}

type attachmentBytes struct {
	id        uuid.UUID
	data      []byte
	mediaType string
}

func (b attachmentBytes) Open() (io.ReadCloser, error) {
	// We clone data because NewBuffer takes ownership of passed bytes,
	// which may result unexpected behavior when attachment is shared
	// between results.
	buf := bytes.NewBuffer(bytes.Clone(b.data))

	return io.NopCloser(buf), nil
}

func (b attachmentBytes) ID() uuid.UUID { return b.id }

func (b attachmentBytes) Type() string { return cmp.Or(b.mediaType, "text/plain") }

func NewAttachmentPath(path string) Attachment {
	return attachmentPath{
		id:   uuid.New(),
		path: path,
	}
}

type attachmentPath struct {
	id   uuid.UUID
	path string
}

func (p attachmentPath) ID() uuid.UUID { return p.id }

func (p attachmentPath) Open() (io.ReadCloser, error) {
	return os.OpenFile(p.path, os.O_RDONLY, 0o600)
}

func (p attachmentPath) Type() string {
	byExtension := mime.TypeByExtension(filepath.Ext(p.path))

	return cmp.Or(byExtension, "text/plain")
}
