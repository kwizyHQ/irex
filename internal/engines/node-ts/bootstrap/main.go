package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/hcl/v2/hclwrite"
	utils "github.com/kwizyHQ/irex/internal/utils"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

func createHCLFile(name, target, pm string, useNodemon bool) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// project block
	project := rootBody.AppendNewBlock("project", nil)
	projectBody := project.Body()

	// Metadata
	projectBody.SetAttributeValue("name", cty.StringVal(name))
	projectBody.SetAttributeValue("description", cty.StringVal("Intermediate Representation Specification"))
	projectBody.SetAttributeValue("version", cty.StringVal("1.0.0"))
	projectBody.SetAttributeValue("author", cty.StringVal("IRS Team"))
	projectBody.SetAttributeValue("license", cty.StringVal("MIT"))
	projectBody.SetAttributeValue("timezone", cty.StringVal("UTC"))

	// paths block
	paths := projectBody.AppendNewBlock("paths", nil)
	pathsBody := paths.Body()
	pathsBody.SetAttributeValue("specifications", cty.StringVal("./spec"))
	pathsBody.SetAttributeValue("templates", cty.StringVal("./spec/templates"))
	pathsBody.SetAttributeValue("output", cty.StringVal("./src/generated"))

	// generator block
	generator := projectBody.AppendNewBlock("generator", nil)
	generatorBody := generator.Body()
	generatorBody.SetAttributeValue("schema", cty.BoolVal(true))
	generatorBody.SetAttributeValue("service", cty.BoolVal(true))
	generatorBody.SetAttributeValue("dry_run", cty.BoolVal(false))
	generatorBody.SetAttributeValue("clean_before", cty.BoolVal(true))

	// runtime block
	runtime := projectBody.AppendNewBlock("runtime", nil)
	runtimeBody := runtime.Body()
	runtimeBody.SetAttributeValue("name", cty.StringVal("node-ts"))
	runtimeBody.SetAttributeValue("scaffold", cty.BoolVal(true))
	runtimeBody.SetAttributeValue("version", cty.StringVal("18.0.0"))

	// runtime.options block
	rtOpts := runtimeBody.AppendNewBlock("options", nil)
	rtOptsBody := rtOpts.Body()
	rtOptsBody.SetAttributeValue("package_manager", cty.StringVal(pm))
	rtOptsBody.SetAttributeValue("entry", cty.StringVal("src/app.ts"))
	rtOptsBody.SetAttributeValue("dev_nodemon", cty.BoolVal(useNodemon))

	// runtime.schema block
	rtSchema := runtimeBody.AppendNewBlock("schema", nil)
	rtSchemaBody := rtSchema.Body()
	rtSchemaBody.SetAttributeValue("framework", cty.StringVal("mongoose"))
	rtSchemaBody.SetAttributeValue("version", cty.StringVal("6.0.0"))
	rtSchemaOpts := rtSchemaBody.AppendNewBlock("options", nil)
	rtSchemaOptsBody := rtSchemaOpts.Body()
	rtSchemaOptsBody.SetAttributeRaw("uri", hclwrite.TokensForFunctionCall("env", hclwrite.TokensForValue(cty.StringVal("MONGO_URI"))))
	rtSchemaOptsBody.SetAttributeRaw("db", hclwrite.TokensForFunctionCall("env", hclwrite.TokensForValue(cty.StringVal("MONGO_DB"))))

	// runtime.service block
	rtService := runtimeBody.AppendNewBlock("service", nil)
	rtServiceBody := rtService.Body()
	rtServiceBody.SetAttributeValue("framework", cty.StringVal("fastify"))
	rtServiceBody.SetAttributeValue("version", cty.StringVal("4.0.0"))
	rtServiceOpts := rtServiceBody.AppendNewBlock("options", nil)
	rtServiceOptsBody := rtServiceOpts.Body()
	rtServiceOptsBody.SetAttributeValue("logger", cty.BoolVal(true))
	rtServiceOptsBody.SetAttributeValue("port", cty.NumberIntVal(8080))
	rtServiceOptsBody.SetAttributeValue("host", cty.StringVal("localhost"))

	// meta block
	meta := projectBody.AppendNewBlock("meta", nil)
	metaBody := meta.Body()
	metaBody.SetAttributeValue("created_at", cty.StringVal(time.Now().Format("2006-01-02")))
	metaBody.SetAttributeValue("generator_version", cty.StringVal("0.1.0"))

	outPath := filepath.Join(target, "irex.hcl")
	fmt.Printf("generating %s...\n", outPath)
	if err := os.WriteFile(outPath, f.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write irex.hcl: %w", err)
	}
	fmt.Printf("created %s (package manager=%s, nodemon=%t)\n", outPath, pm, useNodemon)
	return nil
}

// NewNodeTsCmd returns a Cobra command for node-ts bootstrapping
func Run() *cobra.Command {
	var name string
	var target string
	var pm string
	var useNodemon bool
	var schemaFramework string
	var serviceFramework string

	cmd := &cobra.Command{
		Use:   "node-ts",
		Short: "Bootstrap a node-ts project",
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.OsEnvCheck("IREX_NAME", &name)
			utils.OsEnvCheck("IREX_TARGET", &target)
			utils.OsEnvCheck("IREX_NODE_TS_PKG_MANAGER", &pm)
			utils.OsEnvCheck("IREX_NODE_TS_SCHEMA_FRAMEWORK", &schemaFramework)
			utils.OsEnvCheck("IREX_NODE_TS_SERVICE_FRAMEWORK", &serviceFramework)
			utils.OsEnvCheck("IREX_NODE_TS_DEV_NODEMON", &useNodemon)

			utils.AskFlagInput(utils.InputOption{
				Message: "Project name",
				Help:    "Name of your IREX project",
				Default: "IREX",
				Type:    utils.InputString,
			}, &name)

			utils.AskFlagInput(utils.InputOption{
				Message:  "Target path (must not already exist)",
				Help:     "Path to the target directory for your IREX project",
				Default:  "",
				Type:     utils.InputString,
				Required: true,
			}, &target)
			if target == "." {
				entries, err := os.ReadDir(target)
				if err != nil {
					if !(len(entries) == 1 && entries[0].Name() == ".env") {
						return fmt.Errorf("failed to read current directory: %w", err)
					}
				}
				if len(entries) > 1 {
					return fmt.Errorf("current directory must be empty")
				}
			} else if _, err := os.Stat(target); err == nil {
				return fmt.Errorf("target path %s already exists", target)
			}

			utils.AskFlagInput(utils.InputOption{
				Message: "Schema framework",
				Help:    "Choose a schema framework",
				Default: "mongoose",
				Options: []string{"mongoose", "sequelize"},
				Type:    utils.InputSelect,
			}, &schemaFramework)

			utils.AskFlagInput(utils.InputOption{
				Message: "Service framework",
				Help:    "Choose a service framework",
				Default: "fastify",
				Options: []string{"fastify", "express"},
				Type:    utils.InputSelect,
			}, &serviceFramework)

			utils.AskFlagInput(utils.InputOption{
				Message: "Package manager to use (npm/yarn)",
				Help:    "Choose your preferred package manager",
				Default: "npm",
				Options: []string{"npm", "yarn"},
				Type:    utils.InputSelect,
			}, &pm)

			utils.AskFlagInput(utils.InputOption{
				Message: "Add nodemon as a devDependency?",
				Help:    "Should nodemon be installed for development?",
				Default: false,
				Type:    utils.InputBool,
			}, &useNodemon)

			utils.SetEnvVar("IREX_NAME", name)
			utils.SetEnvVar("IREX_TARGET", target)
			utils.SetEnvVar("IREX_RUNTIME", "node-ts")
			utils.SetEnvVar("IREX_NODE_TS_SCHEMA_FRAMEWORK", schemaFramework)
			utils.SetEnvVar("IREX_NODE_TS_SERVICE_FRAMEWORK", serviceFramework)
			utils.SetEnvVar("IREX_NODE_TS_PKG_MANAGER", pm)

			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create target: %w", err)
			}

			if err := createHCLFile(name, target, pm, useNodemon); err != nil {
				return err
			}

			if err := Scaffold(); err != nil {
				return fmt.Errorf("scaffold failed: %w", err)
			}
			fmt.Println("scaffold completed")
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Project name")
	cmd.Flags().StringVarP(&target, "target", "t", "", "Target path (must not already exist)")
	cmd.Flags().StringVarP(&pm, "pkg", "p", "npm", "Package manager to use (npm/yarn)")
	cmd.Flags().BoolVarP(&useNodemon, "nodemon", "d", false, "Add nodemon as a devDependency?")
	cmd.Flags().StringVarP(&schemaFramework, "schema", "s", "", "Schema framework")
	cmd.Flags().StringVarP(&serviceFramework, "service", "v", "", "Service framework")

	return cmd
}
