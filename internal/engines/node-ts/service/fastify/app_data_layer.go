package fastify

import (
	"github.com/kwizyHQ/irex/internal/ir"
)

type AppDataLayer struct {
	EnvPort int
	EnvHost string
}

func BuildAppDataLayer(irb *ir.IRBundle) *AppDataLayer {
	dl := &AppDataLayer{
		EnvPort: irb.Config.Runtime.Service.Server.Port,
		EnvHost: irb.Config.Runtime.Service.Server.Host,
	}
	return dl
}
