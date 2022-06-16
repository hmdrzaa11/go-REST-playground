package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{Name: "Mocato", Price: 2.99, SKU: "abcd-abcd-abcd"}
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
