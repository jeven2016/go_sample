package main

import (
	flag "github.com/spf13/pflag"
	"move-repository/pkg/department"
	"runtime"
	"strings"
)

var oaRootDepId *string = flag.StringP("rootId", "i", "670869647114347", "The root id of departments defined in OA json file")
var importType *string = flag.StringP("importType", "t", "all", "The type of json file, valid values: department, user, all")
var srcFilePath *string = flag.StringP("source-department-file", "s", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/oa-dep.json", "The source json file")
var destFileDepPath *string = flag.StringP("destination-department-file", "d", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/iam-dep.json", "The path to save the converted json file")
var srcUserFilePath *string = flag.StringP("source-user-file", "u", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/oa-user.json", "The source json file")
var destUserFileDepPath *string = flag.StringP("destination-user-file", "o", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/iam-user.json", "The path to save the converted json file")
var realRoles *string = flag.StringP("realm-roles", "r", "member", "the realm roles to configure for this use, multiple roles separate with comma. format: realm-role,member")

func main() {
	flag.Parse()

	// check all parameters should be specified
	var params = &[]*string{importType, srcFilePath, destFileDepPath, realRoles}
	validateParams(params)
	validateImportType(importType)

	var roles = strings.Split(*realRoles, ",")

	runtime.GOMAXPROCS(5)

	switch *importType {
	case "department":
		department.ConvertDepartments(oaRootDepId, srcFilePath, destFileDepPath, true)
	case "user":
		department.ConvertUsers(srcUserFilePath, destUserFileDepPath, roles)
	case "all":
		department.ConvertAll(oaRootDepId, srcFilePath, destFileDepPath, srcUserFilePath, destUserFileDepPath, roles)
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
			panic("Some parameters are required")
		}
	}
}
