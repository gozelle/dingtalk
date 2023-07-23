package dingtalk

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

type UserClient struct {
	client *Client
}

type ListIDsResult struct {
	UserIdList []string `json:"userid_list"`
}

type User struct {
	UserId           string          `json:"userid"`
	UnionId          string          `json:"unionid"`
	Name             string          `json:"name"`
	Avatar           string          `json:"avatar"`
	StateCode        string          `json:"state_code"`
	ManagerUserId    string          `json:"manager_userid"`
	Mobile           string          `json:"mobile"`
	HideMobile       bool            `json:"hide_mobile"`
	Telephone        string          `json:"telephone"`
	JobNumber        string          `json:"job_number"`
	Title            string          `json:"title"`
	Email            string          `json:"email"`
	WorkPlace        string          `json:"work_place"`
	Remark           string          `json:"remark"`
	ExclusiveAccount bool            `json:"exclusive_account"`
	OrgEmail         string          `json:"org_email"`
	DeptIdList       []int64         `json:"dept_id_list"`
	DeptOrderList    []UserDeptOrder `json:"dept_order_list"`
	Extension        string          `json:"extension"`
	HiredDate        *int64          `json:"hired_date"`
	Active           bool            `json:"active"`
	RealAuthed       bool            `json:"real_authed"`
	Senior           bool            `json:"senior"`
	Admin            bool            `json:"admin"`
	Boss             bool            `json:"boss"`
	LeaderInDept     []UserLeader    `json:"leader_in_dept"`
	RoleList         []UserRole      `json:"role_list"`
	UnionEmpExt      *UserUnion      `json:"union_emp_ext"`
}

type UserDeptOrder struct {
	DeptId int64 `json:"dept_id"`
	Order  int64 `json:"order"`
}

type UserLeader struct {
	DeptId int64 `json:"dept_id"`
	Leader bool  `json:"leader"`
}

type UserRole struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	GroupName string `json:"group_name"`
}

type UserUnion struct {
	UserId          string          `json:"userid"`
	UnionEmpMapList []UserUnionItem `json:"union_emp_map_list"`
	CropId          string          `json:"crop_id"`
}

type UserUnionItem struct {
	UserId string `json:"userid"`
	CropId string `json:"crop_id"`
}

type UserGetParams struct {
	UserId   string `json:"userid"`   // 必填
	Language string `json:"language"` // 选填
}

// ListIDs 获取部门用户userid列表
// https://open.dingtalk.com/document/orgapp-server/query-the-list-of-department-userids
// 部门 ID，如果是根部门，该参数传 1
func (d *UserClient) ListIDs(deptId int64) (list []string, err error) {
	req := d.client.newRestyRequest()
	req.SetBody(map[string]int64{
		"dept_id": deptId,
	})
	err = d.client.wrapAccessToken(req)
	if err != nil {
		return
	}

	resp, err := req.Post(d.client.url("/topapi/user/listid"))
	if err != nil {
		return
	}

	data, err := extractResult(resp)
	if err != nil {
		return
	}

	var r = &ListIDsResult{}
	err = json.Unmarshal(data, r)
	if err != nil {
		return
	}
	list = r.UserIdList

	return
}

// UserGet 查询用户详情
// https://open.dingtalk.com/document/orgapp-server/query-user-details
func (d *UserClient) UserGet(params UserGetParams) (user *User, err error) {
	req := d.client.newRestyRequest()
	req.SetBody(params)
	err = d.client.wrapAccessToken(req)
	if err != nil {
		return
	}

	resp, err := req.Post(d.client.url("/topapi/v2/user/get"))
	if err != nil {
		return
	}
	data, err := extractResult(resp)
	if err != nil {
		return
	}
	user = new(User)
	err = json.Unmarshal(data, user)
	if err != nil {
		return
	}
	return
}

// GetUserIdByUnionId 根据 unionid 获取用户 userid
// https://open.dingtalk.com/document/orgapp-server/you-can-call-this-operation-to-retrieve-the-userids-of
func (d *UserClient) GetUserIdByUnionId(unionId string) (userId string, err error) {
	req := d.client.newRestyRequest()
	req.QueryParam.Add("unionid", unionId)
	err = d.client.wrapAccessToken(req)
	if err != nil {
		return
	}

	resp, err := req.Get(d.client.url("/user/getUseridByUnionid"))
	if err != nil {
		return
	}

	userId = gjson.Parse(resp.String()).Get("userid").String()
	return
}
