package data

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"` //here we are adding our custom validator
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New() //create a validator
	//we want to add a custom validation
	validate.RegisterValidation("sku", validateSku)
	return validate.Struct(p) //validate our struct
}

func validateSku(f validator.FieldLevel) bool {
	//we want this format for our SKU "abc-abcd-abcs"
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	//apply it to the field
	m := re.FindAllString(f.Field().String(), -1) //getting the field and casting it into "string"
	return len(m) == 1
}

func (p *Product) FromJson(i io.Reader) error { //we call this method on an empty "product"
	de := json.NewDecoder(i)
	return de.Decode(p)
}

type Products []*Product //make an alias because we want to add some utility methods to it

func (p *Products) ToJson(w io.Writer) error {
	en := json.NewEncoder(w)
	return en.Encode(p)
}

func GetProducts() Products {
	return productList
}

func AddProducts(p *Product) {
	p.ID = GetNextId()
	productList = append(productList, p)
}

func UpdateProducts(id int, p *Product) error {
	prod, index, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = prod.ID //because we already set all other fields we only need the ID
	p.UpdatedOn = time.Now().String()

	productList[index] = p
	return nil
}

func findProduct(id int) (*Product, int, error) {
	for i, prod := range productList {
		if prod.ID == id {
			return prod, i, nil
		}
	}

	return nil, -1, errors.New("product not found")
}

func GetNextId() int {
	lastProduct := productList[len(productList)-1]
	return lastProduct.ID + 1
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Normal milky coffee",
		Price:       2.45,
		SKU:         "abc12",
		CreatedOn:   time.Now().String(),
		UpdatedOn:   time.Now().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Strong coffee",
		Price:       1.99,
		SKU:         "fjf123",
		CreatedOn:   time.Now().String(),
		UpdatedOn:   time.Now().String(),
	},
}
