package views

import (
	"fmt"
	"threatreg/internal/models"
	"threatreg/internal/service"
	"threatreg/tviewapp/common"
	"threatreg/tviewapp/domains/modals"
	instancesViews "threatreg/tviewapp/instances/views"

	"github.com/rivo/tview"
)

func NewDomainDetailView(domain models.Domain, contentContainer ContentContainer) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetTitle(fmt.Sprintf("Domain: %s", domain.Name)).SetBorder(true)

	info := tview.NewTextView()
	info.SetTitle("Domain Information").SetBorder(true)
	info.SetText(fmt.Sprintf("ID: %s\nName: %s\nDescription: %s",
		domain.ID.String(), domain.Name, domain.Description))

	actionBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	actionBar.SetTitle("Actions").SetBorder(true)

	editButton := tview.NewButton("Edit").
		SetSelectedFunc(func() {
			contentContainer.PushContent(modals.CreateEditDomainModal(
				domain.Name, domain.Description,
				func(name, description string) {
					if _, err := service.UpdateDomain(domain.ID, &name, &description); err == nil {
						contentContainer.PopContent()
					}
				},
				func() { contentContainer.PopContent() },
			))
		})

	addInstanceButton := tview.NewButton("Add Instance").
		SetSelectedFunc(func() {
			contentContainer.PushContent(modals.CreateSelectInstanceModal(domain.ID, func() {
				contentContainer.PopContent()
			}))
		})

	actionBar.AddItem(editButton, 0, 1, false)
	actionBar.AddItem(addInstanceButton, 0, 1, false)
	actionBar.AddItem(tview.NewBox(), 0, 3, false) // Spacer

	instancesTable := tview.NewTable().SetBorders(true)
	instancesTable.SetTitle("Instances in Domain").SetBorder(true)

	instancesTable.SetCell(0, 0, tview.NewTableCell("[::b]Name").SetSelectable(false))
	instancesTable.SetCell(0, 1, tview.NewTableCell("[::b]Product").SetSelectable(false))
	instancesTable.SetCell(0, 2, tview.NewTableCell("[::b]Unresolved Threats").SetSelectable(false))
	instancesTable.SetCell(0, 3, tview.NewTableCell("[::b]Actions").SetSelectable(false))

	instances, err := service.GetInstancesByDomainIdWithThreatStats(domain.ID)
	if err != nil {
		// If we can't load instances, show an error message in the table
		instancesTable.SetCell(1, 0, tview.NewTableCell(fmt.Sprintf("Error loading instances: %v", err)))
		instancesTable.SetCell(1, 1, tview.NewTableCell(""))
		instancesTable.SetCell(1, 2, tview.NewTableCell(""))
		instancesTable.SetCell(1, 3, tview.NewTableCell(""))
		instances = []models.InstanceWithThreatStats{} // Set to empty slice to avoid nil access
	} else {
		for i, instance := range instances {
			productName := ""
			if instance.Product.Name != "" {
				productName = instance.Product.Name
			}
			instancesTable.SetCell(i+1, 0, tview.NewTableCell(instance.Name))
			instancesTable.SetCell(i+1, 1, tview.NewTableCell(productName))
			instancesTable.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprintf("%d", instance.UnresolvedThreatCount)))

			// Add remove button in Actions column
			removeButton := "[red]Remove[-]"
			instancesTable.SetCell(i+1, 3, tview.NewTableCell(removeButton).SetSelectable(true))
		}
	}

	instancesTable.SetSelectable(true, true)
	instancesTable.SetSelectedFunc(func(row, column int) {
		if row > 0 && row-1 < len(instances) {
			instance := instances[row-1]

			// Handle Actions column (Remove button)
			if column == 3 {
				// Show confirmation modal
				contentContainer.PushContent(common.CreateConfirmationModal(
					"Remove Instance",
					fmt.Sprintf("Are you sure you want to remove instance '%s' from this domain?", instance.Name),
					func() {
						// onYes callback - remove instance from domain
						err := service.RemoveInstanceFromDomain(domain.ID, instance.ID)
						if err != nil {
							// TODO: Show error message in the future
							return
						}
						// Refresh the view to reflect changes
						contentContainer.PopContent()
					},
					func() {
						// onNo callback - just close modal, go back to domain detail view
						contentContainer.PopContent()
					},
				))
			} else {
				// Navigate to instance detail for other columns
				contentContainer.PushContentWithFactory(func() tview.Primitive {
					return instancesViews.NewInstanceThreatManager(instance.ID, contentContainer)
				})
			}
		}
	})

	backButton := tview.NewButton("Back to Domains").
		SetSelectedFunc(func() {
			contentContainer.PopContent()
		})

	flex.AddItem(info, 5, 1, false)
	flex.AddItem(actionBar, 3, 0, false)
	flex.AddItem(instancesTable, 0, 2, true)
	flex.AddItem(backButton, 1, 0, false)

	return flex
}
