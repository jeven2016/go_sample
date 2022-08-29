package department

import (
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/fileutil"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

const (
	userMaxNumber = 300
)

func ConvertDepartments(oaRootDepId *string, srcFilePath *string, destFileDepPath *string,
	writeFile bool) ([]*IamDepartment, map[int64]*IamDepartment) {

	var oaDeps = Import[OaDepartment](srcFilePath)
	var depList []*IamDepartment

	rootId, err := convertor.ToInt(*oaRootDepId)
	HandleError(err)

	var depMap = make(map[int64]*IamDepartment)

	for _, dep := range *oaDeps {
		iamDep := &IamDepartment{
			Id:             dep.Id,
			Name:           dep.Name,
			Priority:       dep.SortId,
			Enabled:        dep.Enabled,
			Description:    dep.WholeName,
			SubDepartments: []*IamDepartment{},
			Users:          []*IamDepartmentUser{},
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

func ConvertUsers(srcUserFilePath *string, destUserFileDepPath *string, roles []string, defaultPassword *string) []*IamUser {
	users := loadUsers(srcUserFilePath, roles, defaultPassword)
	userCount := len(users)

	//每300个用户保存到一个文件
	var fileCount int = (len(users) + userMaxNumber - 1) / userMaxNumber

	if fileCount == 1 {
		saveFile[*IamUserRoot](&IamUserRoot{Users: users}, destUserFileDepPath)
	} else {
		ext := path.Ext(*destUserFileDepPath)
		base := path.Dir(*destUserFileDepPath)
		file := path.Base(*destUserFileDepPath)

		for i := 0; i < fileCount; i++ {
			newFilePath := strings.ReplaceAll(file, ext, fmt.Sprintf("-%v%v", i, ext))
			newFilePath = path.Join(base, newFilePath)
			//0   300
			//1  1*300 -> 2*300
			if i == fileCount-1 {
				saveFile[*IamUserRoot](&IamUserRoot{Users: users[i*userMaxNumber : userCount]}, &newFilePath)
			} else {
				saveFile[*IamUserRoot](&IamUserRoot{Users: users[i*userMaxNumber : (i+1)*userMaxNumber]}, &newFilePath)
			}
		}
	}

	return users
}

func ConvertAll(oaRootDepId *string, srcFilePath *string, destFileDepPath *string,
	srcUserFilePath *string, destUserFileDepPath *string, roles []string, defaultPassword *string) {

	users := ConvertUsers(srcUserFilePath, destUserFileDepPath, roles, defaultPassword)
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

func loadUsers(srcUserFilePath *string, roles []string, defaultPassword *string) []*IamUser {
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
				"职务类别": user.OrgLevelName,
				"上级领导": userNameMap[user.Reporter],
			},
			RealmRoles:      roles,
			Credentials:     []*Credential{{"password", *defaultPassword}},
			RequiredActions: []string{"UPDATE_PASSWORD"},
		}
		addTime("入职时间", user.HireDate, iamUser)
		addTime("登记时间", user.CreateTime, iamUser)
		addTime("上次更新", user.UpdateTime, iamUser)
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

func addTime(key string, timeDate int64, iamUser *IamUser) {
	if timeDate == 0 {
		return
	}
	iamUser.Attributes[key] = time.UnixMilli(timeDate).Format("2006-01-02")
}
