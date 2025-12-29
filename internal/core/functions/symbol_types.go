package functions

// EnvKind represents the kind of environment reference (env, var, secret)
type EnvKind string

const (
	EnvKindEnv    EnvKind = "env"
	EnvKindVar    EnvKind = "var"
	EnvKindSecret EnvKind = "secret"
)

// EnvRef represents a reference to an environment variable, variable, or secret
type EnvRef struct {
	Name string  `hcl:"name,attr" cty:"name"`
	Kind EnvKind `hcl:"kind,attr" cty:"kind"`
}
