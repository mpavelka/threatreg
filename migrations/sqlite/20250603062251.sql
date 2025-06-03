-- Create "products" table
CREATE TABLE `products` (
  `id` uuid NOT NULL,
  `name` varchar NULL,
  `description` text NULL,
  PRIMARY KEY (`id`)
);
-- Create "applications" table
CREATE TABLE `applications` (
  `id` uuid NOT NULL,
  `instance_of` uuid NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_products_applications` FOREIGN KEY (`instance_of`) REFERENCES `products` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create "controls" table
CREATE TABLE `controls` (
  `id` uuid NOT NULL,
  `title` varchar NULL,
  `description` text NULL,
  PRIMARY KEY (`id`)
);
-- Create "threats" table
CREATE TABLE `threats` (
  `id` uuid NOT NULL,
  `title` varchar NULL,
  `description` text NULL,
  PRIMARY KEY (`id`)
);
-- Create "threat_assignments" table
CREATE TABLE `threat_assignments` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `threat_id` uuid NULL,
  `product_id` uuid NULL,
  `application_id` uuid NULL,
  CONSTRAINT `fk_threats_threat_assignments` FOREIGN KEY (`threat_id`) REFERENCES `threats` (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT `fk_applications_threat_assignments` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT `fk_products_threat_assignments` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "threat_assignments_id" to table: "threat_assignments"
CREATE UNIQUE INDEX `threat_assignments_id` ON `threat_assignments` (`id`);
-- Create "control_assignments" table
CREATE TABLE `control_assignments` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `threat_assignment_id` integer NOT NULL,
  `control_id` uuid NULL,
  CONSTRAINT `fk_threat_assignments_control_assignments` FOREIGN KEY (`threat_assignment_id`) REFERENCES `threat_assignments` (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT `fk_controls_control_assignments` FOREIGN KEY (`control_id`) REFERENCES `controls` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "control_assignments_id" to table: "control_assignments"
CREATE UNIQUE INDEX `control_assignments_id` ON `control_assignments` (`id`);
-- Create "threat_controls" table
CREATE TABLE `threat_controls` (
  `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  `threat_id` uuid NULL,
  `control_id` uuid NULL,
  CONSTRAINT `fk_threats_threat_controls` FOREIGN KEY (`threat_id`) REFERENCES `threats` (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT `fk_controls_threat_controls` FOREIGN KEY (`control_id`) REFERENCES `controls` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "threat_controls_id" to table: "threat_controls"
CREATE UNIQUE INDEX `threat_controls_id` ON `threat_controls` (`id`);
