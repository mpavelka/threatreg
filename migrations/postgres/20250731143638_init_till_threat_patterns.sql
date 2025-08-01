-- Create "controls" table
CREATE TABLE "public"."controls" (
  "id" uuid NOT NULL,
  "title" character varying(255) NULL,
  "description" text NULL,
  CONSTRAINT "uni_controls_id" PRIMARY KEY ("id")
);
-- Create "threat_assignment_resolution_delegations" table
CREATE TABLE "public"."threat_assignment_resolution_delegations" (
  "id" uuid NOT NULL,
  "delegated_by" uuid NOT NULL,
  "delegated_to" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_delegation_unique" to table: "threat_assignment_resolution_delegations"
CREATE UNIQUE INDEX "idx_delegation_unique" ON "public"."threat_assignment_resolution_delegations" ("delegated_by", "delegated_to");
-- Create "products" table
CREATE TABLE "public"."products" (
  "id" uuid NOT NULL,
  "name" character varying(255) NULL,
  "description" text NULL,
  CONSTRAINT "uni_products_id" PRIMARY KEY ("id")
);
-- Create index "idx_products_name" to table: "products"
CREATE INDEX "idx_products_name" ON "public"."products" ("name");
-- Create "instances" table
CREATE TABLE "public"."instances" (
  "id" uuid NOT NULL,
  "name" character varying(255) NULL,
  "instance_of" uuid NULL,
  CONSTRAINT "uni_instances_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_instances" FOREIGN KEY ("instance_of") REFERENCES "public"."products" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_instances_name" to table: "instances"
CREATE INDEX "idx_instances_name" ON "public"."instances" ("name");
-- Create "threats" table
CREATE TABLE "public"."threats" (
  "id" uuid NOT NULL,
  "title" character varying(255) NULL,
  "description" text NULL,
  CONSTRAINT "uni_threats_id" PRIMARY KEY ("id")
);
-- Create "threat_assignments" table
CREATE TABLE "public"."threat_assignments" (
  "id" bigserial NOT NULL,
  "threat_id" uuid NULL,
  "product_id" uuid NULL,
  "instance_id" uuid NULL,
  CONSTRAINT "uni_threat_assignments_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_instances_threat_assignments" FOREIGN KEY ("instance_id") REFERENCES "public"."instances" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_products_threat_assignments" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_threats_threat_assignments" FOREIGN KEY ("threat_id") REFERENCES "public"."threats" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_threat_assignment" to table: "threat_assignments"
CREATE UNIQUE INDEX "idx_threat_assignment" ON "public"."threat_assignments" ("threat_id", "product_id", "instance_id");
-- Create "control_assignments" table
CREATE TABLE "public"."control_assignments" (
  "id" bigserial NOT NULL,
  "threat_assignment_id" bigint NOT NULL,
  "control_id" uuid NULL,
  CONSTRAINT "uni_control_assignments_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_controls_control_assignments" FOREIGN KEY ("control_id") REFERENCES "public"."controls" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_threat_assignments_control_assignments" FOREIGN KEY ("threat_assignment_id") REFERENCES "public"."threat_assignments" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create "domains" table
CREATE TABLE "public"."domains" (
  "id" uuid NOT NULL,
  "name" character varying(255) NULL,
  "description" text NULL,
  CONSTRAINT "uni_domains_id" PRIMARY KEY ("id")
);
-- Create index "idx_domains_name" to table: "domains"
CREATE INDEX "idx_domains_name" ON "public"."domains" ("name");
-- Create "domain_instances" table
CREATE TABLE "public"."domain_instances" (
  "instance_id" uuid NOT NULL,
  "domain_id" uuid NOT NULL,
  PRIMARY KEY ("instance_id", "domain_id"),
  CONSTRAINT "fk_domain_instances_domain" FOREIGN KEY ("domain_id") REFERENCES "public"."domains" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_domain_instances_instance" FOREIGN KEY ("instance_id") REFERENCES "public"."instances" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "tags" table
CREATE TABLE "public"."tags" (
  "id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "color" character varying(7) NULL,
  CONSTRAINT "uni_tags_id" PRIMARY KEY ("id")
);
-- Create index "idx_tags_name" to table: "tags"
CREATE UNIQUE INDEX "idx_tags_name" ON "public"."tags" ("name");
-- Create "instance_tags" table
CREATE TABLE "public"."instance_tags" (
  "tag_id" uuid NOT NULL,
  "instance_id" uuid NOT NULL,
  PRIMARY KEY ("tag_id", "instance_id"),
  CONSTRAINT "fk_instance_tags_instance" FOREIGN KEY ("instance_id") REFERENCES "public"."instances" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_instance_tags_tag" FOREIGN KEY ("tag_id") REFERENCES "public"."tags" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "threat_patterns" table
CREATE TABLE "public"."threat_patterns" (
  "id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "threat_id" uuid NOT NULL,
  "is_active" boolean NOT NULL,
  CONSTRAINT "uni_threat_patterns_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_threat_patterns_threat" FOREIGN KEY ("threat_id") REFERENCES "public"."threats" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "pattern_conditions" table
CREATE TABLE "public"."pattern_conditions" (
  "id" uuid NOT NULL,
  "pattern_id" uuid NOT NULL,
  "condition_type" character varying(50) NOT NULL,
  "operator" character varying(20) NOT NULL,
  "value" character varying(255) NULL,
  "relationship_type" character varying(100) NULL,
  CONSTRAINT "uni_pattern_conditions_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_threat_patterns_conditions" FOREIGN KEY ("pattern_id") REFERENCES "public"."threat_patterns" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "product_tags" table
CREATE TABLE "public"."product_tags" (
  "tag_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  PRIMARY KEY ("tag_id", "product_id"),
  CONSTRAINT "fk_product_tags_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_product_tags_tag" FOREIGN KEY ("tag_id") REFERENCES "public"."tags" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "relationships" table
CREATE TABLE "public"."relationships" (
  "id" uuid NOT NULL,
  "type" character varying(255) NULL,
  "from_instance_id" uuid NULL,
  "from_product_id" uuid NULL,
  "to_instance_id" uuid NULL,
  "to_product_id" uuid NULL,
  CONSTRAINT "uni_relationships_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_relationships_from_instance" FOREIGN KEY ("from_instance_id") REFERENCES "public"."instances" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_relationships_from_product" FOREIGN KEY ("from_product_id") REFERENCES "public"."products" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_relationships_to_instance" FOREIGN KEY ("to_instance_id") REFERENCES "public"."instances" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_relationships_to_product" FOREIGN KEY ("to_product_id") REFERENCES "public"."products" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_relationships_from_instance_id" to table: "relationships"
CREATE INDEX "idx_relationships_from_instance_id" ON "public"."relationships" ("from_instance_id");
-- Create index "idx_relationships_from_product_id" to table: "relationships"
CREATE INDEX "idx_relationships_from_product_id" ON "public"."relationships" ("from_product_id");
-- Create index "idx_relationships_to_instance_id" to table: "relationships"
CREATE INDEX "idx_relationships_to_instance_id" ON "public"."relationships" ("to_instance_id");
-- Create index "idx_relationships_to_product_id" to table: "relationships"
CREATE INDEX "idx_relationships_to_product_id" ON "public"."relationships" ("to_product_id");
-- Create index "idx_relationships_type" to table: "relationships"
CREATE INDEX "idx_relationships_type" ON "public"."relationships" ("type");
-- Create "threat_assignment_resolutions" table
CREATE TABLE "public"."threat_assignment_resolutions" (
  "id" uuid NOT NULL,
  "threat_assignment_id" bigint NOT NULL,
  "instance_id" uuid NULL,
  "product_id" uuid NULL,
  "status" character varying(20) NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_threat_assignment_resolutions_instance" FOREIGN KEY ("instance_id") REFERENCES "public"."instances" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_threat_assignment_resolutions_product" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_threat_assignment_resolutions_threat_assignment" FOREIGN KEY ("threat_assignment_id") REFERENCES "public"."threat_assignments" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_threat_assignment_resolution" to table: "threat_assignment_resolutions"
CREATE UNIQUE INDEX "idx_threat_assignment_resolution" ON "public"."threat_assignment_resolutions" ("threat_assignment_id", "instance_id", "product_id");
-- Create "threat_controls" table
CREATE TABLE "public"."threat_controls" (
  "id" bigserial NOT NULL,
  "threat_id" uuid NULL,
  "control_id" uuid NULL,
  CONSTRAINT "uni_threat_controls_id" PRIMARY KEY ("id"),
  CONSTRAINT "fk_controls_threat_controls" FOREIGN KEY ("control_id") REFERENCES "public"."controls" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_threats_threat_controls" FOREIGN KEY ("threat_id") REFERENCES "public"."threats" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
