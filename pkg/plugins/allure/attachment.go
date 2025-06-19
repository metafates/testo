package allure

import (
	"bytes"
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

// Image types for attachments supported by Allure report.
//
// See also [screenshot attachments].
//
// [screenshot attachments]: https://allurereport.org/docs/attachments/#screenshots
const (
	ImagePNG  MediaType = "image/png"
	ImageJPEG MediaType = "image/jpeg"
	ImageWEBP MediaType = "image/webp"
	ImageGIF  MediaType = "image/gif"
	ImageSVG  MediaType = "image/svg+xml"
	ImageTIFF MediaType = "image/tiff"
	ImageBMP  MediaType = "image/bmp"
)

// Video types for attachments supported by Allure report.
//
// See also [video attachments].
//
// [video attachments]: https://allurereport.org/docs/attachments/#videos
const (
	VideoMP4  MediaType = "video/mp4"
	VideoOGG  MediaType = "video/ogg"
	VideoWebM MediaType = "video/webm"
)

// Text types for attachments supported by Allure report.
//
// See also [text attachments].
//
// [text attachments]: https://allurereport.org/docs/attachments/#text
const (
	TextPlain MediaType = "text/plain"
	TextHTML  MediaType = "text/html"
)

// Table types for attachments supported by Allure report.
//
// See also [table attachments].
//
// [table attachments]: https://allurereport.org/docs/attachments/#tables
const (
	TableCSV MediaType = "text/csv"
	TableTSV MediaType = "text/tab-separated-values"
)

// URIList is uri list type for attachments.
//
// See also [uri lists attachments].
//
// [uri lists attachments]: https://allurereport.org/docs/attachments/#uri-lists
const URIList MediaType = "text/uri-list"

// Document types for attachments supported by Allure report.
//
// See also [document attachments].
//
// [document attachments]: https://allurereport.org/docs/attachments/#documents
const (
	DocumentXML  MediaType = "text/xml"
	DocumentJSON MediaType = "application/json"
	DocumentYAML MediaType = "application/yaml"
)

// Attachment to add into report.
//
// See [Allure attachments] for more information.
//
// [Allure attachments]: https://allurereport.org/docs/attachments/
type Attachment interface {
	// Open attachment for reading.
	Open() (io.ReadCloser, error)

	// UUID returns the unique id of this attachment.
	UUID() UUID

	// Type returns the media type of the content.
	Type() MediaType
}

// AttachmentBytes is an attachment which stores its contents in-memory.
// Consider using [AttachmentPath] for large files.
type AttachmentBytes struct {
	id        UUID
	data      []byte
	mediaType MediaType
}

// NewAttachmentBytes creates a new bytes attachment from the given bytes.
//
// See [AttachmentBytes.As] to specify media type to enable preview in Allure report.
func NewAttachmentBytes(data []byte) AttachmentBytes {
	return AttachmentBytes{
		id:   uuid.NewString(),
		data: data,
	}
}

// As returns new attachment with the given media type.
func (b AttachmentBytes) As(mediaType MediaType) AttachmentBytes {
	b.mediaType = mediaType

	return b
}

// Open attachment for reading.
func (b AttachmentBytes) Open() (io.ReadCloser, error) {
	// We clone data because NewBuffer takes ownership of passed bytes,
	// which may result unexpected behavior when attachment is shared
	// between results.
	//
	// TODO: avoid initial cloning with something like [io.Pipe].
	buf := bytes.NewBuffer(bytes.Clone(b.data))

	return io.NopCloser(buf), nil
}

// UUID returns the unique id of this attachment.
func (b AttachmentBytes) UUID() UUID { return b.id }

// Type returns the media type of the content.
func (b AttachmentBytes) Type() MediaType { return b.mediaType }

// AttachmentPath is an attachment
// stored in the file located at path.
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

// UUID returns the unique id of this attachment.
func (p AttachmentPath) UUID() UUID { return p.id }

// Open attachment for reading.
func (p AttachmentPath) Open() (io.ReadCloser, error) {
	return os.OpenFile(p.path, os.O_RDONLY, 0o600)
}

// Type returns the media type of the content.
func (p AttachmentPath) Type() MediaType {
	return MediaType(mime.TypeByExtension(filepath.Ext(p.path)))
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

	return NewAttachmentBytes(data).As(p.Type()), nil
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
