-- Create "tags" table
CREATE TABLE `tags` (
  `id` uuid NOT NULL,
  `name` varchar NOT NULL,
  `description` text NULL,
  `color` varchar NULL,
  PRIMARY KEY (`id`)
);
-- Create index "idx_tags_name" to table: "tags"
CREATE UNIQUE INDEX `idx_tags_name` ON `tags` (`name`);
-- Create "instance_tags" table
CREATE TABLE `instance_tags` (
  `tag_id` uuid NOT NULL,
  `instance_id` uuid NOT NULL,
  PRIMARY KEY (`tag_id`, `instance_id`),
  CONSTRAINT `fk_instance_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_instance_tags_instance` FOREIGN KEY (`instance_id`) REFERENCES `instances` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "product_tags" table
CREATE TABLE `product_tags` (
  `tag_id` uuid NOT NULL,
  `product_id` uuid NOT NULL,
  PRIMARY KEY (`tag_id`, `product_id`),
  CONSTRAINT `fk_product_tags_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_product_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
