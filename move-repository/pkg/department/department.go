package department

import (
	"container/list"
	"encoding/json"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/fileutil"
)

func ConvertDepartments(oaRootDepId *string, srcFilePath *string, destFileDepPath *string) {
	var oaDeps = ImportOaDepartments(srcFilePath)
	var depList = make([]IamDepartment, 40)

	rootId, err := convertor.ToInt(oaRootDepId)
	HandleError(err)

	for _, dep := range *oaDeps {
		iamDep := IamDepartment{
			Id:          dep.Id,
			Priority:    dep.SortId,
			Enabled:     dep.Enabled,
			Description: dep.WholeName,
		}

		if dep.Id == rootId {
			depList = append(depList, iamDep)
			continue
		}

		//查找父节点并添加到子列表中去
		appended := appendChild(iamDep, &dep, depList)
		if !appended {
			panic("Failed to process the data: " + dep.WholeName + ", Id=" + convertor.ToString(dep.Id))
		}

	}
}

func appendChild(iamDep *IamDepartment, oaDep *OaDepartment, depList *list.List) bool {
	for dep := depList.Front(); dep != nil; dep = dep.Next() {
		existingIamDep := dep.Value.(IamDepartment)
		if oaDep.Superior == existingIamDep.Id {
			var subDesps = append(*existingIamDep.SubDepartments, *iamDep)
			existingIamDep.SubDepartments = &subDesps
			return true
		}

		for sub := range *existingIamDep.SubDepartments {
			return appendChild(iamDep, oaDep)
		}
	}
}

func ImportOaDepartments(srcFilePath *string) *[]OaDepartment {
	data, err := fileutil.ReadFileToString(*srcFilePath)
	if err != nil {
		panic(err)
	}

	var oaDeps []OaDepartment
	err = json.Unmarshal([]byte(data), &oaDeps)
	if err != nil {
		panic(err)
	}
	return &oaDeps
}
