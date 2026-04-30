package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zzl/go-winmd/apimodel"
	"github.com/zzl/go-winmd/mdmodel"
)

func main() {
	mdFilePath := `c:\Users\Admin\Downloads\winmd\go-winapi-gen\assets\Windows.Wdk.winmd`

	mdModelParser := mdmodel.NewModelParser()
	mdModel, err := mdModelParser.Parse(mdFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer mdModel.Close()

	apiModelParser := apimodel.NewModelParser(nil)
	apiModel := apiModelParser.Parse(mdModel)

	outputFile, err := os.Create("wdk_analysis.txt")
	if err != nil {
		log.Panic(err)
	}
	defer outputFile.Close()

	for _, ns := range apiModel.AllNamespaces {
		if ns.Name == "" || len(ns.Types) == 0 {
			continue
		}
		fmt.Fprintln(outputFile, ns.FullName)
		fmt.Fprintln(outputFile, "  Types:", len(ns.Types))
		for _, typ := range ns.Types {
			fmt.Fprintf(outputFile, "\t%s (Kind: %d)\n", typ.Name, typ.Kind)
			if typ.Struct {
				fmt.Fprintln(outputFile, "\t\tStruct fields:")
				for _, f := range typ.StructDef.Fields {
					fmt.Fprintf(outputFile, "\t\t\t%s %s\n", f.Name, f.Type.Name)
				}
			} else if typ.Interface {
				fmt.Fprintln(outputFile, "\t\tInterface methods:")
				for _, m := range typ.InterfaceDef.Methods {
					fmt.Fprintf(outputFile, "\t\t\tfunc %s", m.Name)
					fmt.Fprint(outputFile, "(")
					for i, p := range m.Params {
						if i > 0 {
							fmt.Fprint(outputFile, ", ")
						}
						fmt.Fprintf(outputFile, "%s %s", p.Name, p.Type.Name)
					}
					fmt.Fprintf(outputFile, ") %s\n", m.ReturnType.Name)
				}
			}
		}
	}

	println("Analysis completed. Output written to wdk_analysis.txt")
}
