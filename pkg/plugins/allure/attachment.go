package allure

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// MediaType stands for MIME type strings.
//
// See also [list of all official MIME types].
//
// [list of all official MIME types]: https://www.iana.org/assignments/media-types/media-types.xhtml
type MediaType string

func (m MediaType) String() string { return string(m) }

// Common media types for allure attachments.
const (
	AttachmentTypePNG  MediaType = "image/png"
	AttachmentTypeJPEG MediaType = "image/jpeg"
	AttachmentTypeWEBP MediaType = "image/webp"
	AttachmentTypeGIF  MediaType = "image/gif"
	AttachmentTypeSVG  MediaType = "image/svg+xml"

	AttachmentTypeMP4 MediaType = "video/mp4"

	AttachmentTypeCSV  MediaType = "text/csv"
	AttachmentTypeText MediaType = "text/plain"
	AttachmentTypeHTML MediaType = "text/html"

	AttachmentTypeMP3 MediaType = "audio/mp3"

	AttachmentTypePDF  MediaType = "application/pdf"
	AttachmentTypeJSON MediaType = "application/json"
	AttachmentTypeXML  MediaType = "application/xml"
	AttachmentTypeDocX MediaType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	AttachmentTypeDoc  MediaType = "application/msword"
	AttachmentTypeXLS  MediaType = "application/vnd.ms-excel"
	AttachmentTypeZIP  MediaType = "application/zip"
	AttachmentTypeTAR  MediaType = "application/x-tar"
)

// Attachment to add into report.
//
// See [Allure attachments] for more information.
//
// [Allure attachments]: https://allurereport.org/docs/attachments/
type Attachment interface {
	// Open attachment for reading.
	Open() (io.ReadCloser, error)

	// UUID is the unique id of this attachment.
	UUID() UUID

	// Type is the media type of the content.
	Type() MediaType
}

// AttachmentBytes is an attachment which stores its contents in-memory.
// Consider using [AttachmentPath] for large files.
type AttachmentBytes struct {
	id        UUID
	data      []byte
	mediaType MediaType
}

// NewAttachmentBytes creates a new bytes attachment from the given bytes and media type.
// If media type is empty text/plain will be used instead.
func NewAttachmentBytes(data []byte, mediaType MediaType) AttachmentBytes {
	return AttachmentBytes{
		id:        uuid.NewString(),
		data:      data,
		mediaType: mediaType,
	}
}

func (b AttachmentBytes) Open() (io.ReadCloser, error) {
	// We clone data because NewBuffer takes ownership of passed bytes,
	// which may result unexpected behavior when attachment is shared
	// between results.
	//
	// TODO: avoid initial cloning with something like [io.Pipe].
	buf := bytes.NewBuffer(bytes.Clone(b.data))

	return io.NopCloser(buf), nil
}

func (b AttachmentBytes) UUID() UUID { return b.id }

func (b AttachmentBytes) Type() MediaType { return cmp.Or(b.mediaType, "text/plain") }

type AttachmentPath struct {
	id   UUID
	path string
}

// NewAttachmentPath creates a new attachment from the given file path.
// Note that file at path won't be read until all suite tests
// are finished (in AfterAll hook).
// Use [AttachmentPath.Read] method to convert it to [AttachmentBytes].
func NewAttachmentPath(path string) AttachmentPath {
	return AttachmentPath{
		id:   uuid.NewString(),
		path: path,
	}
}

func (p AttachmentPath) UUID() UUID { return p.id }

func (p AttachmentPath) Open() (io.ReadCloser, error) {
	return os.OpenFile(p.path, os.O_RDONLY, 0o600)
}

func (p AttachmentPath) Type() MediaType {
	byExtension := MediaType(mime.TypeByExtension(filepath.Ext(p.path)))

	return cmp.Or(byExtension, AttachmentTypeText)
}

// Read the file at path and return it as [AttachmentBytes].
// Error is returned if reading a file failed.
//
// See also [AttachmentPath.MustRead].
func (p AttachmentPath) Read() (AttachmentBytes, error) {
	data, err := os.ReadFile(p.path)
	if err != nil {
		return AttachmentBytes{}, err
	}

	return NewAttachmentBytes(data, p.Type()), nil
}

// MustRead reads the file at path and return it as [AttachmentBytes].
// Panics on error.
//
// See also [AttachmentPath.Read] for non-panicking version.
func (p AttachmentPath) MustRead() AttachmentBytes {
	b, err := p.Read()
	if err != nil {
		panic(fmt.Errorf("error in AttachmentPath.MustRead: %w", err))
	}

	return b
}
