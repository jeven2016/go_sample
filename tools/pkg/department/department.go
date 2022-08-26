package department

import (
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/fileutil"
	"io/ioutil"
	"time"
)

func ConvertDepartments(oaRootDepId *string, srcFilePath *string, destFileDepPath *string, writeFile bool) ([]*IamDepartment, map[int64]*IamDepartment) {
	var oaDeps = Import[OaDepartment](srcFilePath)
	var depList []*IamDepartment

	rootId, err := convertor.ToInt(*oaRootDepId)
	HandleError(err)

	var depMap = make(map[int64]*IamDepartment)

	for _, dep := range *oaDeps {
		iamDep := &IamDepartment{
			Id:          dep.Id,
			Name:        dep.Name,
			Priority:    dep.SortId,
			Enabled:     dep.Enabled,
			Description: dep.WholeName,
		}

		depMap[dep.Id] = iamDep

		if dep.Superior == rootId {
			depList = append(depList, iamDep)
			continue
		}

		//查找父节点并添加到子列表中去
		appended := appendChild(iamDep, &dep, depList)
		if !appended {
			depList = append(depList, iamDep)
		}
	}

	if writeFile {
		saveFile[*IamDepartmentRoot](&IamDepartmentRoot{Departments: depList}, destFileDepPath)
	}
	return depList, depMap
}

func ConvertUsers(srcUserFilePath *string, destUserFileDepPath *string, roles []string) []*IamUser {
	users := loadUsers(srcUserFilePath, roles)
	saveFile[*IamUserRoot](&IamUserRoot{Users: users}, destUserFileDepPath)
	return users
}

func ConvertAll(oaRootDepId *string, srcFilePath *string, destFileDepPath *string,
	srcUserFilePath *string, destUserFileDepPath *string, roles []string) {

	users := ConvertUsers(srcUserFilePath, destUserFileDepPath, roles)
	depList, depMap := ConvertDepartments(oaRootDepId, srcFilePath, destFileDepPath, false)

	for _, user := range users {
		dep := depMap[user.DepartmentId]
		if dep == nil {
			msg, err := fmt.Printf("No department found for user %v(%v)", user.Username, user.FirstName)
			HandleError(err)
			panic(msg)
		}
		departmentUsers := append(dep.Users, &IamDepartmentUser{Username: user.Username})
		dep.Users = departmentUsers
	}

	saveFile[*IamDepartmentRoot](&IamDepartmentRoot{Departments: depList}, destFileDepPath)
}

func loadUsers(srcUserFilePath *string, roles []string) []*IamUser {
	var users = Import[OaUser](srcUserFilePath)

	var userNameMap = make(map[int64]string)
	for _, u := range *users {
		userNameMap[u.Id] = u.Name
	}

	var iamUsers []*IamUser

	for _, user := range *users {
		iamUser := &IamUser{
			DepartmentId: user.DepartmentId,
			Username:     user.LoginName,
			FirstName:    user.Name,
			Enabled:      true,
			Email:        user.EmailAddress,
			Attributes: map[string]string{
				"手机号":  user.TelNumber,
				"拼音全名": user.Pinyin,
				"拼音简称": user.PinyinHead,
				"入职时间": time.Unix(user.HireDate, 0).Format("2006-01-02"),
				"职务":   user.OrgLevelName,
				"上级领导": userNameMap[user.Reporter],
			},
			RealmRoles: roles,
		}
		iamUsers = append(iamUsers, iamUser)
	}
	return iamUsers
}

func appendChild(iamDep *IamDepartment, oaDep *OaDepartment, depList []*IamDepartment) bool {
	for _, existingIamDep := range depList {
		if oaDep.Superior == existingIamDep.Id {
			var subDesps = append(existingIamDep.SubDepartments, iamDep)
			existingIamDep.SubDepartments = subDesps
			iamDep.ParentName = existingIamDep.Name
			return true
		}

		result := appendChild(iamDep, oaDep, existingIamDep.SubDepartments)
		if result {
			return result
		}
	}
	return false
}

func Import[T OaDepartment | OaUser](srcFilePath *string) *[]T {
	data, err := fileutil.ReadFileToString(*srcFilePath)
	if err != nil {
		panic(err)
	}

	var oaDeps []T
	err = json.Unmarshal([]byte(data), &oaDeps)
	if err != nil {
		panic(err)
	}
	return &oaDeps
}

func saveFile[T any](data T, destFilePath *string) {
	jsonData, err := convertor.ToJson(data)
	HandleError(err)

	err = ioutil.WriteFile(*destFilePath, []byte(jsonData), 0664)
	HandleError(err)
}
