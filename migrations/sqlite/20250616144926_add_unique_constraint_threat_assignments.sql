-- Create index "idx_threat_assignment" to table: "threat_assignments"
CREATE UNIQUE INDEX `idx_threat_assignment` ON `threat_assignments` (`threat_id`, `product_id`, `instance_id`);
