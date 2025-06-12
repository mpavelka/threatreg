-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Drop "applications" table
DROP TABLE `applications`;
-- Create "new_threat_assignments" table
CREATE TABLE `new_threat_assignments` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `threat_id` uuid NULL,
  `product_id` uuid NULL,
  `instance_id` uuid NULL,
  CONSTRAINT `fk_threats_threat_assignments` FOREIGN KEY (`threat_id`) REFERENCES `threats` (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT `fk_products_threat_assignments` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT `fk_instances_threat_assignments` FOREIGN KEY (`instance_id`) REFERENCES `instances` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);
-- Copy rows from old table "threat_assignments" to new temporary table "new_threat_assignments"
INSERT INTO `new_threat_assignments` (`id`, `threat_id`, `product_id`) SELECT `id`, `threat_id`, `product_id` FROM `threat_assignments`;
-- Drop "threat_assignments" table after copying rows
DROP TABLE `threat_assignments`;
-- Rename temporary table "new_threat_assignments" to "threat_assignments"
ALTER TABLE `new_threat_assignments` RENAME TO `threat_assignments`;
-- Create index "threat_assignments_id" to table: "threat_assignments"
CREATE UNIQUE INDEX `threat_assignments_id` ON `threat_assignments` (`id`);
-- Create "instances" table
CREATE TABLE `instances` (
  `id` uuid NOT NULL,
  `name` varchar NULL,
  `instance_of` uuid NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_products_instances` FOREIGN KEY (`instance_of`) REFERENCES `products` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_instances_name" to table: "instances"
CREATE INDEX `idx_instances_name` ON `instances` (`name`);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
