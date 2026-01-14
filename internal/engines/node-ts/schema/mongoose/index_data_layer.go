package mongoose

import "github.com/kwizyHQ/irex/internal/ir"

type IndexDataLayer struct {
	Models       []string
	URI          string
	DatabaseName string
}

func BuildIndexDataLayer(ir *ir.IRBundle) *IndexDataLayer {
	dl := &IndexDataLayer{
		URI:          "process.env." + ir.Config.Runtime.Schema.Database.URI,
		DatabaseName: "process.env." + ir.Config.Runtime.Schema.Database.DB,
	}
	for _, m := range ir.Models {
		dl.Models = append(dl.Models, m.Name)
	}
	return dl
}
