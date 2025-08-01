// Atlas configuration file
// See: https://atlasgo.io/getting-started

data "external_schema" "gorm_sqlite" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",
    "--dialect", "sqlite",
  ]
}

data "external_schema" "gorm_postgres" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/models",
    "--dialect", "postgres",
  ]
}


env "sqlite" {
  src = data.external_schema.gorm_sqlite.url
  dev = "sqlite://app.dev.db"
  url = "sqlite://app.db"
  migration {
    dir = "file://migrations/sqlite"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "postgres" {
  src = data.external_schema.gorm_postgres.url
  dev = getenv("APP_DEV_DATABASE_URL")
  url = getenv("APP_DATABASE_URL")
  migration {
    dir = "file://migrations/postgres"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
