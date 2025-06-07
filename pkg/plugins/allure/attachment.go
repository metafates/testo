package allure

import (
	"bytes"
	"cmp"
	"errors"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Attachment interface {
	// Open attachment for reading.
	//
	// If returned value also implements [io.Closer]
	// it will be closed appropriately.
	Open() (io.Reader, error)

	// ID is the unique ID of this attachment.
	ID() uuid.UUID

	// Validate attachment.
	// It must exist and must not be a directory.
	Validate() error

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

func (b attachmentBytes) Open() (io.Reader, error) {
	// We clone data because NewBuffer takes ownership of passed bytes,
	// which may result unexpected behavior when attachment is shared
	// between results.
	return bytes.NewBuffer(bytes.Clone(b.data)), nil
}

func (b attachmentBytes) ID() uuid.UUID { return b.id }

func (b attachmentBytes) Validate() error { return nil }

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

func (p attachmentPath) Validate() error {
	stat, err := os.Stat(p.path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.New("dir attachment")
	}

	return nil
}

func (p attachmentPath) Open() (io.Reader, error) {
	return os.OpenFile(p.path, os.O_RDONLY, 0o600)
}

func (p attachmentPath) Type() string {
	byExtension := mime.TypeByExtension(filepath.Ext(p.path))

	return cmp.Or(byExtension, "text/plain")
}
