Example usages:

```go
package main

import (
	"fmt"
	"log"

	"github.com/267H/orderedform"
)

func main() {
	testOrderedForm()
	fmt.Println("\nOrderedForm tests passed successfully!")
}

func testOrderedForm() {
	fmt.Println("\nTesting OrderedForm:")

	form := orderedform.NewForm(6)
	form.Set("username", "john_doe")
	form.Set("email", "john@example.com")
	form.Set("password", "p@ssw0rd!")
	form.Set("confirm_password", "p@ssw0rd!")
	form.Set("preferences[theme]", "dark")
	form.Set("preferences[notifications]", "email,sms")

	formData := form.URLEncode()
	fmt.Printf("Form data: %s\n", formData)
	expectedForm := "username=john_doe&email=john%40example.com&password=p%40ssw0rd%21&confirm_password=p%40ssw0rd%21&preferences%5Btheme%5D=dark&preferences%5Bnotifications%5D=email%2Csms"
	assertEqual(formData, expectedForm, "Complex Form")
	fmt.Println("✓ Form encoding verified")
}

func assertEqual(got, expected, testName string) {
	if got != expected {
		log.Fatalf("\n%s test failed.\nExpected: %s\nGot: %s", testName, expected, got)
	}
}

