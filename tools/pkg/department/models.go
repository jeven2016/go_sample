package department

type OaDepartment struct {
	Id        int64  `json:"id"`
	SortId    int    `json:"sortId"`
	Enabled   bool   `json:"enabled"`
	Name      string `json:"name"`
	Superior  int64  `json:"superior"`
	WholeName string `json:"wholeName"`
}

type OaUser struct {
	Id           int64  `json:"id"`
	LoginName    string `json:"loginName"`
	Name         string `json:"name"`
	DepartmentId int64  `json:"orgDepartmentId,omitempty"`
	Pinyin       string `json:"pinyin,omitempty"`
	PinyinHead   string `json:"pinyinhead,omitempty"`
	EmailAddress string `json:"emailAddress"`
	HireDate     int64  `json:"hiredate,omitempty"`
	TelNumber    string `json:"telnumber,omitempty"`
	Reporter     int64  `json:"reporter,omitempty"`
	OrgLevelName string `json:"orgLevelName,omitempty"`
}

type IamDepartmentUser struct {
	Username string `json:"username"`
}

type IamDepartment struct {
	Id             int64                `json:"-"`
	Priority       int                  `json:"priority"`
	Enabled        bool                 `json:"enabled"`
	Description    string               `json:"description"`
	SubDepartments []*IamDepartment     `json:"subDepartments"`
	Users          []*IamDepartmentUser `json:"users"`
}

type IamUser struct {
	Username     string            `json:"username"`
	FirstName    string            `json:"firstName"`
	Enabled      bool              `json:"enabled"`
	Email        string            `json:"email"`
	DepartmentId int64             `json:"-"`
	Attributes   map[string]string `json:"attributes"`
}

type IamDepartmentRoot struct {
	Departments []*IamDepartment `json:"departments"`
}

type IamUserRoot struct {
	Users []*IamUser `json:"users"`
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}

}
