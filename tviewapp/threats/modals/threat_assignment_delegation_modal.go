package modals

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// DelegationModalState holds the state for the delegation modal
type DelegationModalState struct {
	sourceResolution        models.ThreatAssignmentResolution
	availableComponents     []models.Component
	instanceOptions         []string
	filteredResolutions     []models.ThreatAssignmentResolution
	resolutionOptions       []string
	selectedComponentIndex  int
	selectedResolutionIndex int
	onSave                  func(targetResolution models.ThreatAssignmentResolution)
	onClose                 func()
	form                    *tview.Form
}

// NewDelegationModalState creates a new delegation modal state
func NewDelegationModalState(
	sourceResolution models.ThreatAssignmentResolution,
	onSave func(targetResolution models.ThreatAssignmentResolution),
	onClose func(),
) (*DelegationModalState, error) {
	// Load instances for the dropdown
	instances, err := service.ListComponents()
	if err != nil {
		return nil, fmt.Errorf("error loading instances: %v", err)
	}

	// Filter out the current instance
	var availableComponents []models.Component
	instanceOptions := []string{}
	for _, instance := range instances {
		if instance.ID != sourceResolution.ComponentID {
			availableComponents = append(availableComponents, instance)
			instanceOptions = append(instanceOptions, instance.Name)
		}
	}

	return &DelegationModalState{
		sourceResolution:        sourceResolution,
		availableComponents:     availableComponents,
		instanceOptions:         instanceOptions,
		selectedComponentIndex:  -1,
		selectedResolutionIndex: -1,
		onSave:                  onSave,
		onClose:                 onClose,
		form:                    tview.NewForm(),
	}, nil
}

// GetSelectedResolution returns the currently selected resolution
func (state *DelegationModalState) GetSelectedResolution() *models.ThreatAssignmentResolution {
	if state.selectedResolutionIndex >= 0 && state.selectedResolutionIndex < len(state.filteredResolutions) {
		return &state.filteredResolutions[state.selectedResolutionIndex]
	}
	return nil
}

// LoadResolutionsForComponent loads threat resolutions for the selected instance
func (state *DelegationModalState) LoadResolutionsForComponent(instanceIndex int) error {
	if instanceIndex < 0 || instanceIndex >= len(state.availableComponents) {
		state.filteredResolutions = []models.ThreatAssignmentResolution{}
		state.resolutionOptions = []string{}
		return nil
	}

	selectedComponentID := state.availableComponents[instanceIndex].ID
	resolutions, err := service.ListThreatResolutionsByComponentID(selectedComponentID)
	if err != nil {
		return fmt.Errorf("error loading resolutions: %v", err)
	}

	// Filter and format resolutions
	state.filteredResolutions = []models.ThreatAssignmentResolution{}
	state.resolutionOptions = []string{}

	for _, resolution := range resolutions {
		var assignmentType string
		if resolution.ThreatAssignment.ComponentID != uuid.Nil {
			assignmentType = "Component-level"
		} else {
			assignmentType = "Product-level"
		}
		state.filteredResolutions = append(state.filteredResolutions, resolution)
		state.resolutionOptions = append(
			state.resolutionOptions,
			fmt.Sprintf("(%s) %s", assignmentType, resolution.ThreatAssignment.Threat.Title),
		)
	}

	return nil
}

// PopulateForm populates the form with current state
func (state *DelegationModalState) PopulateForm() {
	state.form.Clear(false)
	state.form.SetBorder(true).SetTitle("Delegate To Component")

	// Component dropdown
	state.form.AddDropDown("Target Component", state.instanceOptions, state.selectedComponentIndex, func(option string, optionIndex int) {
		state.selectedComponentIndex = optionIndex
		state.selectedResolutionIndex = -1 // Reset resolution selection

		// Load resolutions for selected instance
		err := state.LoadResolutionsForComponent(optionIndex)
		if err != nil {
			// TODO: Handle error properly
			return
		}

		// Update only the resolution dropdown instead of repopulating entire form
		state.UpdateResolutionDropdown()
	})

	// Resolution dropdown
	state.form.AddDropDown("Target Resolution", state.resolutionOptions, state.selectedResolutionIndex, func(option string, optionIndex int) {
		state.selectedResolutionIndex = optionIndex
	})

	// Buttons
	state.form.AddButton("Delegate", func() {
		if selectedResolution := state.GetSelectedResolution(); selectedResolution != nil {
			state.onSave(*selectedResolution)
		}
	})

	state.form.AddButton("Cancel", func() {
		state.onClose()
	})
}

// UpdateResolutionDropdown updates only the resolution dropdown without recreating the entire form
func (state *DelegationModalState) UpdateResolutionDropdown() {
	// Get the resolution dropdown (should be at index 1)
	if state.form.GetFormItemCount() >= 2 {
		if dropDown, ok := state.form.GetFormItem(1).(*tview.DropDown); ok {
			dropDown.SetOptions(state.resolutionOptions, func(option string, optionIndex int) {
				state.selectedResolutionIndex = optionIndex
			})
			dropDown.SetCurrentOption(-1) // Reset selection
		}
	}
}

func CreateThreatAssignmentDelegationModal(
	sourceResolution models.ThreatAssignmentResolution,
	onSave func(targetResolution models.ThreatAssignmentResolution),
	onClose func(),
) tview.Primitive {
	// Validate that source resolution has ComponentID
	if sourceResolution.ComponentID == uuid.Nil {
		return createErrorModal("Error: Threat resolutions can only be delegated from an instance to another instance.\nThis resolution is not associated with an instance.", "Delegation Error", onClose)
	}

	// Create modal state
	state, err := NewDelegationModalState(sourceResolution, onSave, onClose)
	if err != nil {
		return createErrorModal(fmt.Sprintf("Error loading instances: %v", err), "Error", onClose)
	}

	// Check if there are available instances
	if len(state.availableComponents) == 0 {
		return createErrorModal("No other instances available for delegation.", "No Components Available", onClose)
	}

	// Populate the form initially
	state.PopulateForm()

	// Show source resolution information
	infoText := fmt.Sprintf("Threat: %s\nComponent: %s",
		sourceResolution.ThreatAssignment.Threat.Title,
		sourceResolution.Component.Name,
	)

	infoView := tview.NewTextView()
	infoView.SetText(infoText)
	infoView.SetBorder(true).SetTitle("Threat Assignment Info")

	// Container for form and info
	formContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	formContainer.AddItem(infoView, 8, 0, false)
	formContainer.AddItem(state.form, 0, 1, true)

	modalContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false)

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)
	centerFlex.AddItem(formContainer, 80, 0, true)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)

	modalContainer.AddItem(centerFlex, 25, 0, true)
	modalContainer.AddItem(tview.NewBox(), 0, 1, false)

	return modalContainer
}

// createErrorModal creates a simple error modal
func createErrorModal(message, title string, onClose func()) tview.Primitive {
	errorView := tview.NewTextView()
	errorView.SetText(message)
	errorView.SetBorder(true).SetTitle(title)
	errorView.SetTextAlign(tview.AlignCenter)

	closeButton := tview.NewButton("Close")
	closeButton.SetSelectedFunc(onClose)

	container := tview.NewFlex().SetDirection(tview.FlexRow)
	container.AddItem(tview.NewBox(), 0, 1, false)

	centerFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)

	errorContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	errorContainer.AddItem(errorView, 6, 0, false)
	errorContainer.AddItem(closeButton, 3, 0, true)

	centerFlex.AddItem(errorContainer, 60, 0, true)
	centerFlex.AddItem(tview.NewBox(), 0, 1, false)

	container.AddItem(centerFlex, 12, 0, true)
	container.AddItem(tview.NewBox(), 0, 1, false)

	return container
}
