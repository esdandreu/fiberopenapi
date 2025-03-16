package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strings"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"golang.org/x/tools/go/packages"
)

func GenerateRouter(packagePath, specPath, outputPath, typeName string) error {
	g := &Generated{}

	// Get the package name and generate file header.
	packageName, err := LoadPackageName(packagePath)
	if err != nil {
		return fmt.Errorf("cannot load package name: %w", err)
	}
	g.Printf("package %s", packageName)
	g.Printf("\n")
	g.Printf("// Code generated by \"genrouter %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")

	// Import required packages.
	g.Printf("import (\n")
	g.Printf("\t\"github.com/gofiber/fiber/v2\"\n")
	g.Printf(")\n\n")

	// Load OpenAPI document.
	spec, err := LoadOpenAPIDocument(specPath)
	if err != nil {
		return err
	}
	routes := ExtractRoutes(spec)

	// Generate the Handlers interface.
	g.Printf("type %s interface {\n", typeName)

	for _, route := range routes {
		g.Printf("\t%s(c *fiber.Ctx) error\n", ToPascalCase(route.OperationId))
	}
	g.Printf("}\n\n")

	// Generate AddHandlers function.
	g.Printf("func Add%s(app *fiber.App, h %s) {\n", typeName, typeName)
	for _, route := range routes {
		g.Printf("\tapp.%s(\"%s\", h.%s)\n",
			route.Method,
			route.Path,
			ToPascalCase(route.OperationId))
	}
	g.Printf("}\n")

	// Write the generated code back to main.go
	if err := g.WriteFile(outputPath); err != nil {
		return fmt.Errorf("cannot write generated code: %w", err)
	}
	return nil
}

// Contains the generated code. One can append content to it using Printf. Once
// done, one can use MustWriteFile to write formatted content to a file.
type Generated struct {
	content bytes.Buffer
}

func (g *Generated) Printf(format string, args ...any) {
	fmt.Fprintf(&g.content, format, args...)
}

func (g *Generated) WriteFile(path string) error {
	src, err := format.Source(g.content.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}
	err = os.WriteFile(path, src, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write generated code: %w", err)
	}
	return nil
}

// Loads the package name of a directory.
func LoadPackageName(dir string) (string, error) {
	cfg := &packages.Config{Mode: packages.NeedName}
	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return "", fmt.Errorf("cannot load package info: %w", err)
	}
	if len(pkgs) == 0 {
		return "", fmt.Errorf("no packages found")
	}
	return pkgs[0].Name, nil
}

// Read and parse an OpenAPI specification file.
func LoadOpenAPIDocument(specPath string) (*libopenapi.DocumentModel[v3.Document], error) {
	specData, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read specification file at %s: %w", specPath, err)
	}
	document, err := libopenapi.NewDocument(specData)
	if err != nil {
		return nil, fmt.Errorf("cannot create new document: %w", err)
	}
	model, errors := document.BuildV3Model()
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %v\n", errors[i])
		}
		return nil, fmt.Errorf("cannot create v3 model from document: %d errors reported", len(errors))
	}
	return model, nil
}

// Converts a string from kebab-case or snake_case to PascalCase
func ToPascalCase(s string) string {
	// First replace hyphens with underscores to handle both kebab-case and snake_case
	s = strings.ReplaceAll(s, "-", "_")

	// Split the string by underscores
	parts := strings.Split(s, "_")

	// Convert each part to title case and join them
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}

	return strings.Join(parts, "")
}

// Convert an OpenAPI path with {param} style parameters to Fiber's :param
// style https://docs.gofiber.io/guide/routing/#parameters.
func ToFiberPath(path string) string {
	re := regexp.MustCompile(`{([^}]+)}`)
	return re.ReplaceAllString(path, ":$1")
}

// Route represents an API endpoint with its method, path and operation ID
type Route struct {
	Method      string
	Path        string
	OperationId string
}

// ExtractRoutes returns a list of routes from an OpenAPI specification
func ExtractRoutes(spec *libopenapi.DocumentModel[v3.Document]) []Route {
	var routes []Route

	for pair := spec.Model.Paths.PathItems.First(); pair != nil; pair = pair.Next() {
		path := ToFiberPath(pair.Key())
		pathItem := pair.Value()

		if pathItem.Get != nil && pathItem.Get.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Get",
				Path:        path,
				OperationId: pathItem.Get.OperationId,
			})
		}
		if pathItem.Put != nil && pathItem.Put.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Put",
				Path:        path,
				OperationId: pathItem.Put.OperationId,
			})
		}
		if pathItem.Post != nil && pathItem.Post.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Post",
				Path:        path,
				OperationId: pathItem.Post.OperationId,
			})
		}
		if pathItem.Delete != nil && pathItem.Delete.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Delete",
				Path:        path,
				OperationId: pathItem.Delete.OperationId,
			})
		}
		if pathItem.Options != nil && pathItem.Options.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Options",
				Path:        path,
				OperationId: pathItem.Options.OperationId,
			})
		}
		if pathItem.Head != nil && pathItem.Head.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Head",
				Path:        path,
				OperationId: pathItem.Head.OperationId,
			})
		}
		if pathItem.Patch != nil && pathItem.Patch.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Patch",
				Path:        path,
				OperationId: pathItem.Patch.OperationId,
			})
		}
		if pathItem.Trace != nil && pathItem.Trace.OperationId != "" {
			routes = append(routes, Route{
				Method:      "Trace",
				Path:        path,
				OperationId: pathItem.Trace.OperationId,
			})
		}
	}

	return routes
}

func main() {
	var packagePath, outputPath, specPath, typeName string
	flag.StringVar(&packagePath, "path", ".", "path to the package to generate the router for; defaults to current directory")
	flag.StringVar(&outputPath, "output", "router.go", "output file name; defaults to router.go")
	flag.StringVar(&specPath, "spec", "", "path to the OpenAPI specification file; must be set")
	flag.StringVar(&typeName, "type-name", "Handlers", "name of the interface to generate; defaults to Handlers")
	flag.Parse()
	if specPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	if err := GenerateRouter(packagePath, specPath, outputPath, typeName); err != nil {
		panic(err)
	}
}
