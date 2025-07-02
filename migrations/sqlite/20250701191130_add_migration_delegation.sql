-- Create "threat_assignment_resolution_delegations" table
CREATE TABLE `threat_assignment_resolution_delegations` (
  `id` uuid NULL,
  `delegated_by` uuid NOT NULL,
  `delegated_to` uuid NOT NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_delegation_unique" to table: "threat_assignment_resolution_delegations"
CREATE UNIQUE INDEX `idx_delegation_unique` ON `threat_assignment_resolution_delegations` (`delegated_by`, `delegated_to`);
