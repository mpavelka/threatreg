-- Create "threat_assignment_resolutions" table
CREATE TABLE `threat_assignment_resolutions` (
  `id` uuid NULL,
  `threat_assignment_id` integer NOT NULL,
  `instance_id` uuid NULL,
  `product_id` uuid NULL,
  `status` varchar NOT NULL,
  `description` text NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_threat_assignment_resolutions_instance` FOREIGN KEY (`instance_id`) REFERENCES `instances` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `fk_threat_assignment_resolutions_threat_assignment` FOREIGN KEY (`threat_assignment_id`) REFERENCES `threat_assignments` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `fk_threat_assignment_resolutions_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_threat_assignment_resolution" to table: "threat_assignment_resolutions"
CREATE UNIQUE INDEX `idx_threat_assignment_resolution` ON `threat_assignment_resolutions` (`threat_assignment_id`, `instance_id`, `product_id`);
