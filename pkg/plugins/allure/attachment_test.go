package allure

import (
	"fmt"
)

func ExampleNewAttachmentBytes() {
	a := NewAttachmentBytes(
		[]byte(`{"name": "value"}`),
	).As(DocumentJSON)

	fmt.Println(a.Type())
	// Output: application/json
}
