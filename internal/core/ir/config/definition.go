package config

import "github.com/kwizyHQ/irex/internal/core/ir"

// ConfigDefinition is the root struct for the config HCL file, matching fastify-mongoose.hcl
type ConfigDefinition struct {
	Project *ProjectBlock `hcl:"project,block"`
}

type ProjectBlock struct {
	Name        string `hcl:"name,optional"`
	Description string `hcl:"description,optional"`
	Version     string `hcl:"version,optional"`
	Author      string `hcl:"author,optional"`
	License     string `hcl:"license,optional"`
	Timezone    string `hcl:"timezone,optional"`

	Paths     *PathsBlock     `hcl:"paths,block"`
	Generator *GeneratorBlock `hcl:"generator,block"`
	Runtime   *RuntimeBlock   `hcl:"runtime,block"`
	Meta      *MetaBlock      `hcl:"meta,block"`
}

type PathsBlock struct {
	Specifications string `hcl:"specifications,optional"`
	Templates      string `hcl:"templates,optional"`
	Output         string `hcl:"output,optional"`
}

type GeneratorBlock struct {
	Schema      bool `hcl:"schema,optional"`
	Service     bool `hcl:"service,optional"`
	DryRun      bool `hcl:"dry_run,optional"`
	CleanBefore bool `hcl:"clean_before,optional"`
}

type RuntimeBlock struct {
	Name     string               `hcl:"name,optional"`
	Scaffold bool                 `hcl:"scaffold,optional"`
	Version  string               `hcl:"version,optional"`
	Options  *RuntimeOptions      `hcl:"options,block"`
	Schema   *RuntimeSchemaBlock  `hcl:"schema,block"`
	Service  *RuntimeServiceBlock `hcl:"service,block"`
}

type RuntimeOptions struct {
	PackageManager string `hcl:"package_manager,optional"`
	Entry          string `hcl:"entry,optional"`
	DevNodemon     bool   `hcl:"dev_nodemon,optional"`
}

type RuntimeSchemaBlock struct {
	Framework string                `hcl:"framework,optional"`
	Version   string                `hcl:"version,optional"`
	Options   *RuntimeSchemaOptions `hcl:"options,block"`
}

type RuntimeSchemaOptions struct {
	URI ir.EnvRef `hcl:"uri,optional"`
	DB  ir.EnvRef `hcl:"db,optional"`
}

type RuntimeServiceBlock struct {
	Framework string                 `hcl:"framework,optional"`
	Version   string                 `hcl:"version,optional"`
	Options   *RuntimeServiceOptions `hcl:"options,block"`
}

type RuntimeServiceOptions struct {
	Logger bool   `hcl:"logger,optional"`
	Port   int    `hcl:"port,optional"`
	Host   string `hcl:"host,optional"`
}

type MetaBlock struct {
	CreatedAt        string `hcl:"created_at,optional"`
	GeneratorVersion string `hcl:"generator_version,optional"`
}
