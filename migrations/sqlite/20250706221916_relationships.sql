-- Create "relationships" table
CREATE TABLE `relationships` (
  `id` uuid NOT NULL,
  `type` varchar NULL,
  `from_instance_id` uuid NULL,
  `from_product_id` uuid NULL,
  `to_instance_id` uuid NULL,
  `to_product_id` uuid NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_relationships_from_product` FOREIGN KEY (`from_product_id`) REFERENCES `products` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `fk_relationships_from_instance` FOREIGN KEY (`from_instance_id`) REFERENCES `instances` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `fk_relationships_to_product` FOREIGN KEY (`to_product_id`) REFERENCES `products` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `fk_relationships_to_instance` FOREIGN KEY (`to_instance_id`) REFERENCES `instances` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_relationships_to_product_id" to table: "relationships"
CREATE INDEX `idx_relationships_to_product_id` ON `relationships` (`to_product_id`);
-- Create index "idx_relationships_to_instance_id" to table: "relationships"
CREATE INDEX `idx_relationships_to_instance_id` ON `relationships` (`to_instance_id`);
-- Create index "idx_relationships_from_product_id" to table: "relationships"
CREATE INDEX `idx_relationships_from_product_id` ON `relationships` (`from_product_id`);
-- Create index "idx_relationships_from_instance_id" to table: "relationships"
CREATE INDEX `idx_relationships_from_instance_id` ON `relationships` (`from_instance_id`);
-- Create index "idx_relationships_type" to table: "relationships"
CREATE INDEX `idx_relationships_type` ON `relationships` (`type`);
