package allure

import (
	"fmt"
)

func ExampleNewAttachmentBytes() {
	a := NewAttachmentBytes(
		[]byte(`{"name": "value"}`),
		AttachmentTypeJSON,
	)

	fmt.Println(a.Type())
	// Output: application/json
}
