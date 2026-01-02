package ir

import "time"

// IRConfig is the fully resolved, engine-consumable configuration.
// This is the final output of the config pipeline.
type IRConfig struct {
	Project   IRProject
	Paths     IRPaths
	Generator IRGenerator
	Runtime   IRRuntime
	Meta      IRMeta
}

// ------------------------------------------------------------
// Project
// ------------------------------------------------------------

type IRProject struct {
	Name        string
	Description string
	Version     string
	Author      string
	License     string
	Timezone    string
}

// ------------------------------------------------------------
// Paths (absolute, resolved)
// ------------------------------------------------------------

type IRPaths struct {
	Specifications string // absolute path
	Templates      string // absolute path
	Output         string // absolute path
}

// ------------------------------------------------------------
// Generator flags
// ------------------------------------------------------------

type IRGenerator struct {
	GenerateSchema  bool
	GenerateService bool
	DryRun          bool
	CleanBefore     bool
}

// ------------------------------------------------------------
// Runtime (single resolved runtime)
// ------------------------------------------------------------

type IRRuntime struct {
	Name     string // node
	Version  string // 18, 20, etc
	Scaffold bool

	Options IRRuntimeOptions
	Schema  IRRuntimeSchema
	Service IRRuntimeService
}

// ------------------------------------------------------------
// Runtime options
// ------------------------------------------------------------

type IRRuntimeOptions struct {
	PackageManager string // npm | pnpm | yarn
	Entry          string // src/server.ts
	DevNodemon     bool
}

// ------------------------------------------------------------
// Schema runtime (database layer)
// ------------------------------------------------------------

type IRRuntimeSchema struct {
	Framework string // mongoose
	Version   string

	Database IRDatabaseConfig
}

type IRDatabaseConfig struct {
	URI string // resolved value (no env refs)
	DB  string
}

// ------------------------------------------------------------
// Service runtime (API layer)
// ------------------------------------------------------------

type IRRuntimeService struct {
	Framework string // fastify
	Version   string

	Server IRServerConfig
}

type IRServerConfig struct {
	Logger bool
	Port   int
	Host   string
}

// ------------------------------------------------------------
// Meta
// ------------------------------------------------------------

type IRMeta struct {
	CreatedAt        time.Time
	GeneratorVersion string
}
