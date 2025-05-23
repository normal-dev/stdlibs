// TODO: Read https://stackoverflow.com/questions/32132064/how-to-discover-all-package-types-at-runtime
package api

import (
	"fmt"
	"go/types"
	"log"
	"regexp"

	"golang.org/x/tools/go/packages"
)

func init() {
	log.SetFlags(0)
}

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
	pkgs := getAllPkgs()
	stripePkgs(pkgs)
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
					switch typ.Kind() {
					case types.UntypedBool:
						api.Type = types.Typ[types.Bool].String()

					case types.UntypedInt:
						api.Type = types.Typ[types.Int].String()

					case types.UntypedFloat:
						api.Type = types.Typ[types.Float64].String()

					case types.UntypedString:
						api.Type = types.Typ[types.String].String()

					case types.UntypedRune:
						api.Type = types.Typ[types.Rune].String()

					case types.UntypedNil:
						api.Type = types.Typ[types.UntypedNil].String()

					default:
						api.Type = typ.Name()
					}

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

			log.Printf("namespace: %s", api.Ns)
			log.Printf("name: %s", api.Name)
			log.Printf("type: %s", api.Type)

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
		log.Printf("pkg: %s", pkg.ID)
		for _, name := range pkg.Types.Scope().Names() {
			log.Printf("name: %s", name)
			pkgs[pkg.ID] = append(pkgs[pkg.ID], pkg.Types.Scope().Lookup(name))
		}
	}
	return pkgs
}

func stripePkgs(pkgs map[string][]types.Object) {
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
