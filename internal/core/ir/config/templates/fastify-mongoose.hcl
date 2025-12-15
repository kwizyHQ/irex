# ============================================================================
# IREX PROJECT CONFIGURATION
# ============================================================================
# This file defines how IREX should:
# - locate specifications
# - generate code
# - target a runtime
# - configure frameworks (schema & service)
# ============================================================================

project {

  # --------------------------------------------------------------------------
  # Project metadata (informational only)
  # --------------------------------------------------------------------------
  # Name of the project / specification set
  name        = "IRS"

  # Short description of what this project represents
  description = "Intermediate Representation Specification"

  # Project version (used for tracking, not runtime behavior)
  version     = "1.0.0"

  # Project author or organization
  author      = "IRS Team"

  # License identifier (MIT, Apache-2.0, etc.)
  license     = "MIT"

  # Default timezone used for generated metadata, logs, timestamps
  timezone    = "UTC"


  # --------------------------------------------------------------------------
  # Paths configuration
  # --------------------------------------------------------------------------
  paths {

    # Root directory containing IR specifications (.hcl files).
    # This directory MUST contain schema/ and service/ subfolders.
    specifications = "./irs"

    # Directory containing generator templates (.tpl files)
    templates      = "./irs/templates"

    # Output directory where generated code will be written
    output         = "./src/generated"
  }


  # --------------------------------------------------------------------------
  # Generator behavior flags
  # --------------------------------------------------------------------------
  generator {

    # Enable / disable schema generation (models, entities, DB layer)
    schema  = true

    # Enable / disable service generation (routes, controllers, APIs)
    service = true

    # Dry-run mode:
    # - Parses and validates specs
    # - Builds IR
    # - Does NOT write any files
    # Useful for validation, CI checks, and compatibility testing
    dry_run = false

    # If true, previously generated files are removed before regeneration.
    # Only generator-owned files should be deleted (never user code).
    clean_before = true
  }


  # --------------------------------------------------------------------------
  # Runtime configuration
  # --------------------------------------------------------------------------
  # Defines the currently active runtime.
  # Only ONE runtime is active at a time.
  runtime {

    # Runtime identifier (e.g. node-ts, go, rust, php)
    name = "node-ts"

    # If true, scaffold base runtime files (package.json, tsconfig, etc.)
    scaffold = true

    # Target runtime version (minimum supported unless stated otherwise)
    version  = "18.0.0"


    # ------------------------------------------------------------------------
    # Runtime-level options (tooling & entry point)
    # ------------------------------------------------------------------------
    options {

      # Package manager to use for dependency installation
      # Supported: npm | yarn | pnpm
      package_manager = "npm"

      # Application entry file
      entry = "src/app.ts"

      # Enable nodemon for development (Node runtime only)
      dev_nodemon = false
    }


    # ------------------------------------------------------------------------
    # Schema generation configuration (database / models)
    # ------------------------------------------------------------------------
    schema {

      # Schema framework / ORM / ODM
      framework = "mongoose"

      # Framework version used for template selection and compatibility
      version   = "6.0.0"

      options {

        # Database connection URI
        # Value is resolved from the environment at runtime
        uri = env("MONGO_URI")

        # Database name
        # Also resolved from the environment
        db  = env("MONGO_DB")
      }
    }


    # ------------------------------------------------------------------------
    # Service generation configuration (API / routes)
    # ------------------------------------------------------------------------
    service {

      # HTTP framework used to expose services
      framework = "fastify"

      # Framework version for compatibility-aware generation
      version   = "4.0.0"

      options {

        # Enable framework logger
        logger = true

        # Server port
        port = 8080

        # Server host
        host = "localhost"
      }
    }
  }


  # --------------------------------------------------------------------------
  # Metadata (non-functional, informational)
  # --------------------------------------------------------------------------
  meta {

    # Timestamp when this config was created
    created_at = "2025-12-10"

    # Version of IREX generator used to create this project
    generator_version = "0.1.0"
  }
}
