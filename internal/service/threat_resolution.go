package service

import (
	"fmt"
	"threatreg/internal/database"
	"threatreg/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getThreatAssignmentResolutionRepository() (*models.ThreatAssignmentResolutionRepository, error) {
	db, err := database.GetDBOrError()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %w", err)
	}
	return models.NewThreatAssignmentResolutionRepository(db), nil
}

func CreateThreatResolution(
	threatAssignmentID int,
	instanceID *uuid.UUID,
	productID *uuid.UUID,
	status models.ThreatAssignmentResolutionStatus,
	description string,
) (*models.ThreatAssignmentResolution, error) {

	resolution := &models.ThreatAssignmentResolution{
		ThreatAssignmentID: threatAssignmentID,
		Status:             status,
		Description:        description,
	}

	// Set exactly one of InstanceID or ProductID
	if instanceID != nil {
		resolution.InstanceID = *instanceID
		resolution.ProductID = uuid.Nil
	} else if productID != nil {
		resolution.ProductID = *productID
		resolution.InstanceID = uuid.Nil
	} else {
		return nil, fmt.Errorf("either instanceID or productID must be provided")
	}

	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	err = resolutionRepository.Create(nil, resolution)
	if err != nil {
		return nil, fmt.Errorf("error creating threat resolution: %w", err)
	}

	return resolution, nil
}

func UpdateThreatResolution(
	id uuid.UUID,
	status *models.ThreatAssignmentResolutionStatus,
	description *string,
) (*models.ThreatAssignmentResolution, error) {

	var updatedResolution *models.ThreatAssignmentResolution
	err := database.GetDB().Transaction(func(tx *gorm.DB) error {
		resolutionRepository, err := getThreatAssignmentResolutionRepository()
		if err != nil {
			return err
		}

		resolution, err := resolutionRepository.GetByID(tx, id)
		if err != nil {
			return err
		}

		// Update fields if provided
		if status != nil {
			resolution.Status = *status
		}
		if description != nil {
			resolution.Description = *description
		}

		err = resolutionRepository.Update(tx, resolution)
		if err != nil {
			return err
		}
		updatedResolution = resolution

		// Delete any existing delegation for this resolution
		delegationRepository := models.NewThreatAssignmentResolutionDelegationRepository(tx)
		err = delegationRepository.DeleteThreatAssignmentResolutionDelegationBySourceId(tx, resolution.ID)
		if err != nil {
			return fmt.Errorf("error deleting existing delegation: %w", err)
		}

		err = updateUpstreamResolutionsStatus(*resolution, resolution.Status, tx)
		if err != nil {
			return fmt.Errorf("error updating upstream resolutions: %w", err)
		}
		return nil
	})

	return updatedResolution, err
}

func GetThreatResolution(id uuid.UUID) (*models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.GetByID(nil, id)
}

func GetThreatResolutionByThreatAssignmentID(threatAssignmentID int) (*models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.GetOneByThreatAssignmentID(nil, threatAssignmentID)
}

func GetInstanceLevelThreatResolution(threatAssignmentID int, instanceID uuid.UUID) (*models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.GetOneByThreatAssignmentIDAndInstanceID(nil, threatAssignmentID, instanceID)
}

func GetInstanceLevelThreatResolutionWithDelegation(threatAssignmentID int, instanceID uuid.UUID) (*models.ThreatAssignmentResolutionWithDelegation, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	// First get the resolution
	resolution, err := resolutionRepository.GetOneByThreatAssignmentIDAndInstanceID(nil, threatAssignmentID, instanceID)
	if err != nil {
		return nil, err
	}

	// Then get the delegation if it exists
	delegationRepo := models.NewThreatAssignmentResolutionDelegationRepository(database.GetDB())
	delegations, err := delegationRepo.GetThreatAssignmentResolutionDelegations(nil, &resolution.ID, nil)
	if err != nil {
		return nil, err
	}

	result := &models.ThreatAssignmentResolutionWithDelegation{
		Resolution: *resolution,
		Delegation: nil,
	}

	if len(delegations) > 0 {
		result.Delegation = &delegations[0]
	}

	return result, nil
}

func GetProductLevelThreatResolution(threatAssignmentID int, productID uuid.UUID) (*models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.GetOneByAssignmentIDAndProductID(nil, threatAssignmentID, productID)
}

func ListThreatResolutionsByProductID(productID uuid.UUID) ([]models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.ListByProductID(nil, productID)
}

func ListThreatResolutionsByInstanceID(instanceID uuid.UUID) ([]models.ThreatAssignmentResolution, error) {
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	return resolutionRepository.ListByInstanceID(nil, instanceID)
}

func DeleteThreatResolution(id uuid.UUID) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {
		resolutionRepository, err := getThreatAssignmentResolutionRepository()
		if err != nil {
			return err
		}

		// Get the resolution to be deleted
		resolution, err := resolutionRepository.GetByID(tx, id)
		if err != nil {
			return nil // Not found, nothing to delete
		}

		// Update upstream resolutions to "awaiting" status before deletion
		err = updateUpstreamResolutionsStatus(*resolution, models.ThreatAssignmentResolutionStatusAwaiting, tx)
		if err != nil {
			return fmt.Errorf("error updating upstream resolutions: %w", err)
		}

		// Delete any delegations associated with this resolution
		delegationRepository := models.NewThreatAssignmentResolutionDelegationRepository(tx)
		err = delegationRepository.DeleteByDelegatedTo(tx, id)
		if err != nil {
			return fmt.Errorf("error deleting delegations for resolution: %w", err)
		}

		// Delete the resolution
		return resolutionRepository.Delete(tx, id)
	})
}

func ListByDomainWithUnresolvedByInstancesCount(domainID uuid.UUID) ([]models.ThreatWithUnresolvedByInstancesCount, error) {
	threatRepository, err := getThreatRepository()
	if err != nil {
		return nil, err
	}

	return threatRepository.ListByDomainWithUnresolvedByInstancesCount(nil, domainID)
}

func DelegateResolution(threatResolution models.ThreatAssignmentResolution, targetThreatResolution models.ThreatAssignmentResolution) error {
	return database.GetDB().Transaction(func(tx *gorm.DB) error {

		delegationRepository := models.NewThreatAssignmentResolutionDelegationRepository(database.GetDB())

		// Create the delegation record
		delegation := &models.ThreatAssignmentResolutionDelegation{
			DelegatedBy: threatResolution.ID,
			DelegatedTo: targetThreatResolution.ID,
		}

		// Remove existing delegation if it exists
		delegationRepository.DeleteThreatAssignmentResolutionDelegationBySourceId(tx, threatResolution.ID)

		// Create the new delegation
		err := delegationRepository.CreateThreatAssignmentResolutionDelegation(tx, delegation)
		if err != nil {
			return fmt.Errorf("error creating delegation: %w", err)
		}

		// Update the resolution delegation chain
		rootResolution, e := FindResolutionRoot(targetThreatResolution, tx)
		if e != nil {
			return fmt.Errorf("error finding root resolution for delegation: %w", e)
		}
		updateUpstreamResolutionsStatus(*rootResolution, targetThreatResolution.Status, tx)

		return nil
	})
}

func GetDelegatedToResolutionByDelegatedByID(delegatedByID uuid.UUID) (*models.ThreatAssignmentResolution, error) {
	delegationRepository := models.NewThreatAssignmentResolutionDelegationRepository(database.GetDB())

	// Get delegation by source resolution ID
	delegations, err := delegationRepository.GetThreatAssignmentResolutionDelegations(nil, &delegatedByID, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting delegation: %w", err)
	}

	if len(delegations) == 0 {
		return nil, nil // No delegation found
	}

	// Get the target resolution
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	targetResolution, err := resolutionRepository.GetByID(nil, delegations[0].DelegatedTo)
	if err != nil {
		return nil, fmt.Errorf("error getting target resolution: %w", err)
	}

	return targetResolution, nil
}

// updateUpstreamResolutionsStatus recursively updates all resolutions that delegate to the rootResolution
func updateUpstreamResolutionsStatus(rootResolution models.ThreatAssignmentResolution, status models.ThreatAssignmentResolutionStatus, tx *gorm.DB) error {
	delegationRepository := models.NewThreatAssignmentResolutionDelegationRepository(database.GetDB())
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return err
	}

	// Find all resolutions that delegate TO this rootResolution
	delegations, err := delegationRepository.GetThreatAssignmentResolutionDelegations(tx, nil, &rootResolution.ID)
	if err != nil {
		return fmt.Errorf("error getting upstream delegations: %w", err)
	}

	// For each resolution that delegates to this one, update its status and continue recursively
	for _, delegation := range delegations {
		// Get the upstream resolution (the one that delegated to rootResolution)
		upstreamResolution, err := resolutionRepository.GetByID(tx, delegation.DelegatedBy)
		if err != nil {
			return fmt.Errorf("error getting upstream resolution: %w", err)
		}

		// Update the upstream resolution's status to the provided status
		upstreamResolution.Status = status
		err = resolutionRepository.Update(tx, upstreamResolution)
		if err != nil {
			return fmt.Errorf("error updating upstream resolution: %w", err)
		}

		// Recursively update any resolutions that delegate to this upstream resolution
		err = updateUpstreamResolutionsStatus(*upstreamResolution, status, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

// FindResolutionRoot traverses the delegation chain to find the root resolution
func FindResolutionRoot(resolution models.ThreatAssignmentResolution, tx *gorm.DB) (*models.ThreatAssignmentResolution, error) {
	if tx == nil {
		tx = database.GetDB()
	}
	delegationRepository := models.NewThreatAssignmentResolutionDelegationRepository(tx)
	resolutionRepository, err := getThreatAssignmentResolutionRepository()
	if err != nil {
		return nil, err
	}

	currentResolution := resolution
	visited := make(map[uuid.UUID]bool) // To detect cycles

	for {
		// Check if we've seen this resolution before (cycle detection)
		if visited[currentResolution.ID] {
			return nil, fmt.Errorf("cycle detected in delegation chain")
		}
		visited[currentResolution.ID] = true

		// Find if this resolution delegates to another resolution
		delegations, err := delegationRepository.GetThreatAssignmentResolutionDelegations(nil, &currentResolution.ID, nil)
		if err != nil {
			return nil, fmt.Errorf("error getting delegations: %w", err)
		}

		// If no delegation found, this is the root
		if len(delegations) == 0 {
			return &currentResolution, nil
		}

		// Get the target resolution (the one this resolution delegates to)
		targetResolution, err := resolutionRepository.GetByID(nil, delegations[0].DelegatedTo)
		if err != nil {
			return nil, fmt.Errorf("error getting target resolution: %w", err)
		}

		// Continue with the target resolution
		currentResolution = *targetResolution
	}
}
