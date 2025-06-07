package allure

import (
	"fmt"
	"mime"
)

func ExampleNewAttachmentBytes() {
	a := NewAttachmentBytes(
		[]byte(`{"name": "value"}`),
		mime.TypeByExtension(".json"),
	)

	fmt.Println(a.Type())
	// Output: application/json
}
