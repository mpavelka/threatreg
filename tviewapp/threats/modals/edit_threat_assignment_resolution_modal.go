package modals

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func CreateEditThreatAssignmentResolutionModal(
	assignment models.ThreatAssignment,
	existingResolution *models.ThreatAssignmentResolution,
	onSave func(),
	onClose func(),
) tview.Primitive {
	form := tview.NewForm()

	// Set title based on whether we're creating or editing
	if existingResolution == nil {
		form.SetBorder(true).SetTitle("Add Resolution")
	} else {
		form.SetBorder(true).SetTitle("Edit Resolution")
	}

	// Initialize values
	var status string
	var description string
	if existingResolution != nil {
		status = string(existingResolution.Status)
		description = existingResolution.Description
	} else {
		status = string(models.ThreatAssignmentResolutionStatusAwaiting) // Default status
		description = ""
	}

	// Status dropdown
	statusOptions := []string{
		string(models.ThreatAssignmentResolutionStatusAwaiting),
		string(models.ThreatAssignmentResolutionStatusAccepted),
		string(models.ThreatAssignmentResolutionStatusResolved),
	}

	// Find current status index
	statusIndex := 0
	for i, option := range statusOptions {
		if option == status {
			statusIndex = i
			break
		}
	}

	form.AddDropDown("Status", statusOptions, statusIndex, func(option string, optionIndex int) {
		status = option
	})

	form.AddInputField("Description", description, 50, nil, func(text string) {
		description = text
	})

	form.AddButton("Save", func() {
		// Determine if we're creating or updating
		if existingResolution == nil {
			// Create new resolution
			var instanceID *uuid.UUID
			var productID *uuid.UUID

			if assignment.InstanceID != uuid.Nil {
				instanceID = &assignment.InstanceID
			} else {
				productID = &assignment.ProductID
			}

			_, err := service.CreateThreatResolution(
				assignment.ID,
				instanceID,
				productID,
				models.ThreatAssignmentResolutionStatus(status),
				description,
			)
			if err == nil {
				onSave()
			}
		} else {
			// Update existing resolution
			newStatus := models.ThreatAssignmentResolutionStatus(status)
			_, err := service.UpdateThreatResolution(
				existingResolution.ID,
				&newStatus,
				&description,
			)
			if err == nil {
				onSave()
			}
		}
	})

	form.AddButton("Close", func() {
		onClose()
	})

	// Show assignment information in the modal
	infoText := fmt.Sprintf("Threat: %s\nType: %s",
		assignment.Threat.Title,
		func() string {
			if assignment.InstanceID != uuid.Nil {
				return fmt.Sprintf("Instance (%s)", assignment.Instance.Name)
			}
			return fmt.Sprintf("Product (%s)", assignment.Product.Name)
		}())

	infoView := tview.NewTextView()
	infoView.SetText(infoText)
	infoView.SetBorder(true).SetTitle("Assignment Info")

	// Container for form and info
	formContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	formContainer.AddItem(infoView, 6, 0, false)
	formContainer.AddItem(form, 0, 1, true)

	modalContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false)

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)
	centerFlex.AddItem(formContainer, 80, 0, true)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)

	modalContainer.AddItem(centerFlex, 20, 0, true)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false)

	return modalContainer
}
