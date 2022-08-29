package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"move-repository/pkg/department"
	"strings"
)

var oaRootDepId *string = flag.StringP("rootId", "i", "670869647114347", "The root id of departments defined in OA json file")
var importType *string = flag.StringP("importType", "t", "all", "The type of json file, valid values: department, user, all")
var srcFilePath *string = flag.StringP("source-department-file", "s", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/oa-dep.json", "The source json file")
var destFileDepPath *string = flag.StringP("destination-department-file", "d", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/iam-dep.json", "The path to save the converted json file")
var srcUserFilePath *string = flag.StringP("source-user-file", "u", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/oa-user.json", "The source json file")
var destUserFileDepPath *string = flag.StringP("destination-user-file", "o", "/home/jujucom/Desktop/workspace/projects/go_samples/tools/conf/iam-user.json", "The path to save the converted json file")
var realRoles *string = flag.StringP("realm-roles", "r", "member", "the realm roles to configure for this use, multiple roles separate with comma. format: realm-role,member")
var defaultPassword *string = flag.StringP("default-password", "p", "NextGen#2021", "initial password for each user")

func main() {
	flag.Parse()

	// check all parameters should be specified
	var params = &[]*string{srcFilePath, destFileDepPath}
	validateParams(params)

	if len(*realRoles) == 0 {
		panic("realm-roles is required")
	}
	var roles = strings.Split(*realRoles, ",")

	switch *importType {
	case "department":
		department.ConvertDepartments(oaRootDepId, srcFilePath, destFileDepPath, true)
	case "user":
		department.ConvertUsers(srcUserFilePath, destUserFileDepPath, roles, defaultPassword)
	case "all":
		department.ConvertAll(oaRootDepId, srcFilePath, destFileDepPath, srcUserFilePath, destUserFileDepPath, roles, defaultPassword)
	default:
		panic("invalid importType")
	}
}

func validateParams(array *[]*string) {
	for _, param := range *array {
		if len(*param) == 0 {
			panic("The parameters are required: source-department-file, destination-department-file, source-user-file, destination-user-file")
		}
		if !strings.HasSuffix(*param, ".json") && !strings.HasSuffix(*param, ".JSON") {
			msg, _ := fmt.Printf("invalid file path, it should end with '.json' or '.JSON' (%v)", *param)
			panic(msg)
		}
	}
}
