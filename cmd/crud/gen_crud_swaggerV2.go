package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	modelDir   = "internal/app/model"
	outputFile = "internal/crud/crud_docs.go"
)

var actionMap = map[string]string{
	"创建":   "Create",
	"更新":   "Update",
	"删除":   "Delete",
	"分页查询": "Page",
}

type ModelStub struct {
	Name   string
	Prefix string
	Flags  map[string]bool
}

func main() {
	models := scanModels()

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("创建输出文件失败: %v", err)
	}
	defer f.Close()

	write := func(s string) {
		f.WriteString(s + "\n")
	}

	write("package crud\n")
	write("")
	// write("import \"github.com/gin-gonic/gin\"")
	// write("import \"tier-up/internal/app/model\"")
	write("import (")
	write(`  "github.com/gin-gonic/gin"`)
	write(`"tier-up/internal/app/model"`)
	write(")")
	for _, m := range models {
		generateStubs(write, m)
	}
}

// 扫描Struct
func scanModels() []ModelStub {
	var models []ModelStub
	err := filepath.Walk(modelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("解析失败 %s: %v", path, err)
			return nil
		}
		for _, decl := range node.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				model := ModelStub{
					Name:  typeSpec.Name.Name,
					Flags: make(map[string]bool),
				}
				for _, field := range structType.Fields.List {
					if field.Tag != nil {
						tag := strings.Trim(field.Tag.Value, "`")
						crudTag := reflect.StructTag(tag).Get("crud")
						if crudTag == "" {
							continue
						}
						parts := strings.Split(crudTag, ",")
						for _, part := range parts {
							part = strings.TrimSpace(part)
							if strings.HasPrefix(part, "prefix:") {
								model.Prefix = strings.TrimPrefix(part, "prefix:")
							} else {
								model.Flags[part] = true
							}
						}

					}
				}

				if len(model.Flags) > 0 {
					if model.Prefix == "" {
						model.Prefix = "/" + strings.ToLower(model.Name)
					}
					models = append(models, model)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("读取model失败")
		panic(err)
	}
	return models
}

func generateStubs(write func(string), m ModelStub) {
	write("")
	write(fmt.Sprintf("// ===== Auto-generated stub for %s =====", m.Name))

	if m.Flags["create"] {
		printFunc(write, m.Name, m.Prefix+"/create", "post", "创建")
	}
	if m.Flags["delete"] {
		printFunc(write, m.Name, m.Prefix+"/delete/:id", "delete", "删除")
	}
	if m.Flags["update"] {
		printFunc(write, m.Name, m.Prefix+"/update/:id", "put", "更新")
	}

	if m.Flags["page"] {
		printFunc(write, m.Name, m.Prefix+"/page", "get", "分页查询")
	}
	printResponseType(write, m.Name)
}

func printFunc(write func(string), model, path, method, action string) {
	write("")
	write(fmt.Sprintf("// @Summary %s %s", action, model))
	write(fmt.Sprintf("// @Description %s %s", action, model))
	write(fmt.Sprintf("// @Tags %s", model))
	write("// @Accept json")
	write("// @Produce json")
	if method == "post" || method == "put" {
		write(fmt.Sprintf("// @Param data body model.%s true \"%s 数据\"", model, model))
	}
	resp := model
	if action == "分页查询" {
		resp += "PageResponse"
	} else {
		resp += "Response"
	}
	write(fmt.Sprintf("// @Success 200 {object} %s", resp))
	write(fmt.Sprintf("// @Router %s [%s]", path, method))
	write(fmt.Sprintf("func %s%sDoc(ctx *gin.Context) {}\n", model, actionMap[action]))
}

func printResponseType(write func(string), model string) {
	write("")
	write(fmt.Sprintf("type %sResponse struct {", model))
	write("	Code    int    `json:\"code\"`")
	write("	Message string `json:\"message\"`")
	write(fmt.Sprintf("	Data   model. %s `json:\"data\"`", model))
	write("}")

	write("")
	write(fmt.Sprintf("type %sPageResponse struct {", model))
	write("	Code    int    `json:\"code\"`")
	write("	Message string `json:\"message\"`")
	write("	Data struct {")
	write("		Page  int     `json:\"page\"`")
	write("		Limit int     `json:\"limit\"`")
	write("		Total int64   `json:\"total\"`")
	write(fmt.Sprintf("		Data  []model.%s `json:\"data\"`", model))
	write("	} `json:\"data\"`")
	write("}")
}
