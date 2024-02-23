// TODO: Read https://stackoverflow.com/questions/32132064/how-to-discover-all-package-types-at-runtime
package api

import (
	"fmt"
	"go/types"
	"log"
	"regexp"

	"golang.org/x/tools/go/packages"
)

type API struct {
	Doc   string  `json:"doc" bson:"doc"`                         // ToUpper returns s with all Unicode letters mapped to their upper case.
	Name  string  `json:"name" bson:"name"`                       // Reader, Writer, Buffer
	Ns    string  `json:"ns" bson:"ns"`                           // compress/lzw, net, bytes
	Type  string  `json:"type" bson:"type"`                       // struct, error, int, map, func
	Value *string `json:"value,omitempty" bson:"value,omitempty"` // NewFlagSet(os.Args[0], ExitOnError), 512, errors.New("bytes.Buffer: too large")
}

func (api API) ID() string {
	return fmt.Sprintf("%s.%s", api.Ns, api.Name)
}

func Get() []API {
	log.Println("getting all pkgs...")
	pkgs := getAllPkgs()
	log.Println("filtering pkgs...")
	filterPkgs(pkgs)
	log.Println("getting apis from pkgs...")
	return getAPIs(pkgs)
}

func getAPIs(pkgs map[string][]types.Object) []API {
	apis := make([]API, 0)

	for pkg, objs := range pkgs {
		for _, obj := range objs {
			if !obj.Exported() {
				continue
			}

			api := API{
				Name: obj.Name(),
				Ns:   pkg,
			}

			if doc := getGoDoc(api.ID()); doc != "" {
				api.Doc = doc
			}

			switch o := obj.(type) {
			case *types.Var:
				switch o.Type().String() {
				case "error":
					api.Type = "error"

				default:
					switch typ := o.Type().Underlying().(type) {
					case *types.Struct:
						api.Type = "struct"

					case *types.Map:
						api.Type = "map"

					case *types.Interface:
						api.Type = "interface"

					case *types.Basic:
						api.Type = typ.Name()

					case *types.Slice:
						api.Type = "slice"

					// Function or method type
					case *types.Signature:
						api.Type = "type"

					case *types.Array:
						api.Type = "array"

					case *types.Pointer:
						api.Type = "pointer"

					default:
						api.Type = o.Type().String()
					}
				}

			case *types.Const:
				api.Value = toPtr(o.Val().ExactString())

				switch typ := o.Type().Underlying().(type) {
				case *types.Basic:
					api.Type = typ.Name()

				default:
					api.Type = o.Type().String()
				}

			case *types.Func:
				api.Type = "func"

			case *types.TypeName:
				switch typ := o.Type().Underlying().(type) {
				case *types.Struct:
					api.Type = "struct"

				case *types.Map:
					api.Type = "map"

				case *types.Interface:
					api.Type = "interface"

				case *types.Basic:
					api.Type = typ.Name()

				case *types.Slice:
					api.Type = "slice"

				case *types.Signature:
					api.Type = "type"

				case *types.Array:
					api.Type = "array"

				case *types.Pointer:
					api.Type = "pointer"

				default:
					api.Type = o.Type().String()
				}

			default:
				continue
			}

			apis = append(apis, api)
		}
	}

	return apis
}

func getAllPkgs() map[string][]types.Object {
	stdPackages := func() []*packages.Package {
		pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedTypes}, "std")
		checkErr(err)
		return pkgs
	}
	pkgs := make(map[string][]types.Object)
	for _, pkg := range stdPackages() {
		for _, name := range pkg.Types.Scope().Names() {
			pkgs[pkg.ID] = append(pkgs[pkg.ID], pkg.Types.Scope().Lookup(name))
		}
	}
	return pkgs
}

func filterPkgs(pkgs map[string][]types.Object) {
	regexp := regexp.MustCompile("(^vendor|/internal|internal/|/internal/)")
	for pkg, objs := range pkgs {
		if regexp.MatchString(pkg) {
			// Internal, vendor
			delete(pkgs, pkg)
			continue
		}

		for i := len(objs) - 1; i >= 0; i-- {
			if !objs[i].Exported() {
				objs = append(objs[:i], objs[i+1:]...)
			}
		}
	}
}

func toPtr(s string) *string { return &s }

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
