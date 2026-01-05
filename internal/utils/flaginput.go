package utils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
)

// InputType defines the type of prompt to show
// "input" for string, "bool" for yes/no, "select" for options
// You can extend this as needed

type InputType int

const (
	InputString InputType = iota
	InputBool
	InputSelect
)

// InputOption holds the configuration for a prompt

type InputOption struct {
	Message  string
	Help     string
	Default  interface{}
	Options  []string // for select
	Type     InputType
	Required bool // if true, must provide a value
}

// AskFlagInput prompts the user for input using survey, with type safety and redundancy
// dest must be a pointer to the correct type (string, bool, or int for select index)
func AskFlagInput(opt InputOption, dest interface{}) error {
	var validate survey.Validator
	if opt.Required {
		validate = survey.Required
	}
	// check if dest is already have a value (non-zero)
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}
	elem := v.Elem()
	zeroValue := reflect.Zero(elem.Type()).Interface()
	if !reflect.DeepEqual(elem.Interface(), zeroValue) {
		// already has a value, skip prompt
		return nil
	}

	var prompt survey.Prompt
	switch opt.Type {
	case InputString:
		p := &survey.Input{
			Message: opt.Message,
			Help:    opt.Help,
		}
		if def, ok := opt.Default.(string); ok && def != "" {
			p.Default = def
		}
		prompt = p
	case InputBool:
		p := &survey.Confirm{
			Message: opt.Message,
			Help:    opt.Help,
		}
		if def, ok := opt.Default.(bool); ok {
			p.Default = def
		}
		prompt = p
	case InputSelect:
		p := &survey.Select{
			Message: opt.Message,
			Options: opt.Options,
			Help:    opt.Help,
		}
		if def, ok := opt.Default.(string); ok && def != "" {
			p.Default = def
		}
		prompt = p
	default:
		return fmt.Errorf("unsupported input type")
	}
	var err error
	if validate != nil {
		err = survey.AskOne(prompt, dest, survey.WithValidator(validate))
	} else {
		err = survey.AskOne(prompt, dest)
	}
	if err != nil && err.Error() == "interrupt" {
		fmt.Println("\nInput cancelled by user.")
		os.Exit(1)
	}
	return err
}

// OsEnvCheck is a helper to check if an environment variable is set then update its value to the pointer
func OsEnvCheck(envVar string, ptr interface{}) {
	envValue := os.Getenv(envVar)

	// check if ptr is nil and envValue is not empty set the value
	if ptr != nil && envValue != "" {
		switch v := ptr.(type) {
		case *string:
			*v = envValue
		case *bool:
			if b, err := strconv.ParseBool(envValue); err == nil {
				*v = b
			}
		case *int:
			if i, err := strconv.Atoi(envValue); err == nil {
				*v = i
			}
		}
	}

}

func SetEnvVar(envVar string, value string) error {
	return os.Setenv(envVar, value)
}
