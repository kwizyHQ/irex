config {
  name               = "IRS"
  description        = "Intermediate Representation Specification"
  version            = "1.0.0"
  output_root        = "./src"
  specification_root = "./irs"
  generator {
    schema  = true
    service = true
    dry_run = false
  }
  runtime {
    name       = "node-ts"
    scaffold   = true
    output_dir = "."
    options {
      package_manager = "npm"
      entry           = "src/app.ts"
      dev_nodemon     = false
    }
  }
  modules {
    schema {
      framework  = "mongoose"
      output_dir = "vendor/models"
      options {
        uri = "$${env.MONGO_URI}"
        db  = "$${env.MONGO_DB}"
      }
    }
    service {
      framework  = "fastify"
      output_dir = "src/routes"
      options {
        logger = true
        port   = 8080
        host   = "localhost"
      }
    }
  }
  env {
    file    = "./.env"
    require = false
  }
  meta {
    created_at        = "2025-12-10"
    generator_version = "0.1.0"
  }
}
