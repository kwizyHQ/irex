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

	cfg := rootBody.AppendNewBlock("config", nil)
	cfgBody := cfg.Body()
	cfgBody.SetAttributeValue("name", cty.StringVal(name))
	cfgBody.SetAttributeValue("description", cty.StringVal("Intermediate Representation Specification"))
	cfgBody.SetAttributeValue("version", cty.StringVal("1.0.0"))

	gen := cfgBody.AppendNewBlock("generate", nil)
	genBody := gen.Body()
	genBody.SetAttributeValue("schema", cty.BoolVal(true))
	genBody.SetAttributeValue("service", cty.BoolVal(true))
	genBody.SetAttributeValue("output_root", cty.StringVal("./myapp"))
	genBody.SetAttributeValue("force_overwrite", cty.BoolVal(false))
	genBody.SetAttributeValue("verbose", cty.BoolVal(false))

	rt := cfgBody.AppendNewBlock("runtime", nil)
	rtBody := rt.Body()
	rtBody.SetAttributeValue("name", cty.StringVal("node-ts"))
	rtBody.SetAttributeValue("scaffold", cty.BoolVal(true))
	rtBody.SetAttributeValue("output_dir", cty.StringVal("."))
	opts := rtBody.AppendNewBlock("options", nil)
	opts.Body().SetAttributeValue("package_manager", cty.StringVal(pm))
	opts.Body().SetAttributeValue("entry", cty.StringVal("src/app.ts"))
	opts.Body().SetAttributeValue("dev_nodemon", cty.BoolVal(useNodemon))

	mods := cfgBody.AppendNewBlock("modules", nil)
	modsBody := mods.Body()

	schema := modsBody.AppendNewBlock("schema", nil)
	schemaBody := schema.Body()
	schemaBody.SetAttributeValue("framework", cty.StringVal("mongoose"))
	schemaBody.SetAttributeValue("output_dir", cty.StringVal("vendor/models"))
	schOpts := schemaBody.AppendNewBlock("options", nil)
	schOpts.Body().SetAttributeValue("uri", cty.StringVal("${env.MONGO_URI}"))
	schOpts.Body().SetAttributeValue("db", cty.StringVal("${env.MONGO_DB}"))

	service := modsBody.AppendNewBlock("service", nil)
	serviceBody := service.Body()
	serviceBody.SetAttributeValue("framework", cty.StringVal("fastify"))
	serviceBody.SetAttributeValue("output_dir", cty.StringVal("src/routes"))
	svcOpts := serviceBody.AppendNewBlock("options", nil)
	svcOpts.Body().SetAttributeValue("logger", cty.BoolVal(true))
	svcOpts.Body().SetAttributeValue("port", cty.NumberIntVal(8080))
	svcOpts.Body().SetAttributeValue("host", cty.StringVal("localhost"))

	env := cfgBody.AppendNewBlock("env", nil)
	env.Body().SetAttributeValue("file", cty.StringVal("./.env"))
	env.Body().SetAttributeValue("require", cty.BoolVal(false))

	meta := cfgBody.AppendNewBlock("meta", nil)
	meta.Body().SetAttributeValue("created_at", cty.StringVal(time.Now().Format("2006-01-02")))
	meta.Body().SetAttributeValue("generator_version", cty.StringVal("0.1.0"))

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
				Message: "Target path (must not already exist)",
				Help:    "Path to the target directory for your IREX project",
				Default: "",
				Type:    utils.InputString,
			}, &target)
			if _, err := os.Stat(target); err == nil {
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
