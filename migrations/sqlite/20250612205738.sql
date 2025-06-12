-- Create "domains" table
CREATE TABLE `domains` (
  `id` uuid NOT NULL,
  `name` varchar NULL,
  `description` text NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_domains_name" to table: "domains"
CREATE INDEX `idx_domains_name` ON `domains` (`name`);
-- Create "domain_instances" table
CREATE TABLE `domain_instances` (
  `instance_id` uuid NOT NULL,
  `domain_id` uuid NOT NULL,
  PRIMARY KEY (`instance_id`, `domain_id`),
  CONSTRAINT `fk_domain_instances_domain` FOREIGN KEY (`domain_id`) REFERENCES `domains` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_domain_instances_instance` FOREIGN KEY (`instance_id`) REFERENCES `instances` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
