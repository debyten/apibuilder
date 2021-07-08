package registrar

import "github.com/debyten/apibuilder"

type Registrar interface {
	Register(apis []apibuilder.API)
}
