package main

import (
	flag "github.com/spf13/pflag"
	"move-repository/pkg/department"
	"runtime"
)

var oaRootDepId *string = flag.StringP("rootId", "r", "670869647114347", "The root id of departments defined in OA json file")
var importType *string = flag.StringP("importType", "t", "all", "The type of json file, valid values: department, user, all")
var srcFilePath *string = flag.StringP("source-department-file", "s", "", "The source json file")
var destFileDepPath *string = flag.StringP("destination-department-file", "d", "", "The path to save the converted json file")
var srcUserFilePath *string = flag.StringP("source-user-file", "u", "", "The source json file")
var destUserFileDepPath *string = flag.StringP("destination-user-file", "o", "", "The path to save the converted json file")

func main() {
	flag.Parse()

	// check all parameters should be specified
	var params = &[]*string{importType, srcFilePath, destFileDepPath}
	validateParams(params)
	validateImportType(importType)

	runtime.GOMAXPROCS(5)

	switch *importType {
	case "department":
		department.ConvertDepartments(oaRootDepId, srcFilePath, destFileDepPath, true)
	case "user":
		department.ConvertUsers(srcUserFilePath, destUserFileDepPath)
	case "all":
		department.ConvertAll(oaRootDepId, srcFilePath, destFileDepPath, srcUserFilePath, destUserFileDepPath)
	}
}

func validateImportType(t *string) {
	var importType = *t
	if importType != "department" && importType != "user" && importType != "all" {
		panic("invalid importType , it should be department or user")
	}
}

func validateParams(array *[]*string) {
	for _, param := range *array {
		if len(*param) == 0 {
			panic("you should specify the cmd to run: handler or upload")
		}
	}
}
