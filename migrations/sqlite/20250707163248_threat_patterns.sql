-- Create "pattern_conditions" table
CREATE TABLE `pattern_conditions` (
  `id` uuid NOT NULL,
  `pattern_id` uuid NOT NULL,
  `condition_type` varchar NOT NULL,
  `operator` varchar NOT NULL,
  `value` varchar NULL,
  `relationship_type` varchar NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_threat_patterns_conditions` FOREIGN KEY (`pattern_id`) REFERENCES `threat_patterns` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "threat_patterns" table
CREATE TABLE `threat_patterns` (
  `id` uuid NOT NULL,
  `name` varchar NOT NULL,
  `description` text NULL,
  `threat_id` uuid NOT NULL,
  `is_active` numeric NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `fk_threat_patterns_threat` FOREIGN KEY (`threat_id`) REFERENCES `threats` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
