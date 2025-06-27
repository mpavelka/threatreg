package forms

import (
	"fmt"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func CreateNewInstanceForm(domainID uuid.UUID, onClose func()) tview.Primitive {
	nameField := ""
	var selectedProductID uuid.UUID
	createProduct := false

	// Get list of products for dropdown
	products, err := service.ListProducts()
	if err != nil {
		// If we can't load products, show error and return
		errorView := tview.NewTextView().SetText(fmt.Sprintf("Error loading products: %v", err))
		return errorView
	}

	// Create product options for dropdown
	productOptions := make([]string, len(products))
	productMap := make(map[string]uuid.UUID)
	for i, product := range products {
		productOptions[i] = product.Name
		productMap[product.Name] = product.ID
	}

	// Set initial selected product if we have products
	if len(products) > 0 {
		selectedProductID = products[0].ID
	}

	// Create a Pages container to switch between form states
	pages := tview.NewPages()

	// Function to create the form with current state
	var createForm func() *tview.Form
	createForm = func() *tview.Form {
		form := tview.NewForm()

		form.AddInputField("Instance Name", nameField, 50, nil, func(text string) {
			nameField = text
		})

		form.AddCheckbox("Create Product", createProduct, func(checked bool) {
			createProduct = checked
			// Switch to new form state
			newForm := createForm()
			pages.RemovePage("form")
			pages.AddPage("form", newForm, true, true)
		})

		// Only show product dropdown if not creating a new product
		if !createProduct {
			form.AddDropDown("Product", productOptions, 0, func(option string, optionIndex int) {
				selectedProductID = productMap[option]
			})
		}

		form.AddButton("Create & Add", func() {
			if nameField == "" {
				// TODO: Show validation error in the future
				return
			}

			var productID uuid.UUID
			if createProduct {
				// Create new product with the same name as the instance
				product, err := service.CreateProduct(nameField, "")
				if err != nil {
					// TODO: Show error message in the future
					return
				}
				productID = product.ID
			} else {
				productID = selectedProductID
			}

			// Create the instance
			instance, err := service.CreateInstance(nameField, productID)
			if err != nil {
				// TODO: Show error message in the future
				return
			}

			// Associate instance with domain
			err = service.AddInstanceToDomain(domainID, instance.ID)
			if err != nil {
				// TODO: Show error message in the future
				return
			}

			onClose()
		})

		form.AddButton("Cancel", func() {
			onClose()
		})

		return form
	}

	// Add initial form
	initialForm := createForm()
	pages.AddPage("form", initialForm, true, true)

	return pages
}
