// Package dingtalk /*
package dingtalk

import (
	"encoding/json"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

type DepartmentClient struct {
	client *Client
}

func (d *DepartmentClient) ListSub(params *DepartmentListSubParams) (list []*Department, err error) {
	req := d.client.newRestyRequest()
	if params != nil {
		req.SetBody(params)
	}
	err = d.client.wrapAccessToken(req)
	if err != nil {
		return
	}
	resp, err := req.Post(d.client.url("/topapi/v2/department/listsub"))
	if err != nil {
		return
	}
	data, err := extractResult(resp)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &list)
	if err != nil {
		return
	}

	return
}

// ListSubID 获取子部门列表
func (d *DepartmentClient) ListSubID(deptID int64) (list []int64, err error) {
	req := d.client.newRestyRequest()
	req.SetBody(map[string]int64{
		"dept_id": deptID,
	})
	err = d.client.wrapAccessToken(req)
	if err != nil {
		return
	}
	resp, err := req.Post(d.client.url("/topapi/v2/department/listsubid"))
	if err != nil {
		return
	}

	data, err := extractResult(resp)
	if err != nil {
		return
	}

	g := gjson.ParseBytes(data)
	err = json.Unmarshal([]byte(g.Get("dept_id_list").Raw), &list)
	if err != nil {
		return
	}
	return
}

// Get 获取部门详情详情
// https://open.dingtalk.com/document/orgapp-server/query-department-details0-v2
func (d *DepartmentClient) Get(deptID int64) (r *DepartmentDetail, err error) {
	req := d.client.newRestyRequest()
	req.SetBody(map[string]interface{}{
		"dept_id":  deptID,
		"language": "zh_CN",
	})
	err = d.client.wrapAccessToken(req)
	if err != nil {
		return
	}
	resp, err := req.Post(d.client.url("/topapi/v2/department/get"))
	if err != nil {
		return
	}
	data, err := extractResult(resp)
	if err != nil {
		return
	}

	r = new(DepartmentDetail)
	err = json.Unmarshal(data, r)
	if err != nil {
		return
	}

	return
}

// ListUserIDs 获取部门人员 ID 列表
func (d *DepartmentClient) ListUserIDs(deptID int64) (userList []string, err error) {
	req := d.client.newRestyRequest()
	req.SetBody(map[string]int64{
		"dept_id": deptID,
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
	g := gjson.ParseBytes(data)
	err = json.Unmarshal([]byte(g.Get("userid_list").Raw), &userList)
	if err != nil {
		return
	}

	return
}

// ListUsers 获取部门人员详情列表
func (d *DepartmentClient) ListUsers(params ListUsersParams) (r *ListUsersResult, err error) {
	req := d.client.newRestyRequest()
	req.SetBody(params)
	err = d.client.wrapAccessToken(req)
	if err != nil {
		return
	}
	resp, err := req.Post(d.client.url("/topapi/v2/user/list"))
	if err != nil {
		return
	}

	data, err := extractResult(resp)
	if err != nil {
		return
	}

	r = new(ListUsersResult)
	err = json.Unmarshal(data, r)
	if err != nil {
		return
	}

	return
}

type DepartmentListSubParams struct {
	DeptID   int64  `json:"dept_id"`  // 父部门 ID
	Language string `json:"language"` // 通讯录语言，zh_CN: 中文（默认），en_US 英文
}

type Department struct {
	DeptID          int64  `json:"dept_id"`
	Name            string `json:"name"`
	ParentID        int64  `json:"parent_id"`
	CreateDeptGroup bool   `json:"create_dept_group"`
	AutoAddUser     bool   `json:"auto_add_user"`
}

type DepartmentDetail struct {
	DeptID                int64    `json:"dept_id"`
	Name                  string   `json:"name"`
	ParentID              int64    `json:"parent_id"`
	SourceIdentifier      string   `json:"source_identifier"`
	CreateDeptGroup       bool     `json:"create_dept_group"`
	AutoAddUser           bool     `json:"auto_add_user"`
	AutoApproveApply      bool     `json:"auto_approve_apply"`
	FromUnionOrg          bool     `json:"from_union_org"`
	Tags                  string   `json:"tags"`
	Order                 int64    `json:"order"`
	DeptGroupChatID       string   `json:"dept_group_chat_id"`
	GroupContainSubDept   bool     `json:"group_contain_sub_dept"`
	OrgDeptOwner          string   `json:"org_dept_owner"`
	DeptManagerUseridList []string `json:"dept_manager_userid_list"`
	OuterDept             bool     `json:"outer_dept"`
	OuterPermitDepts      []int64  `json:"outer_permit_depts"`
	OuterPermitUsers      []string `json:"outer_permit_users"`
	HideDept              bool     `json:"hide_dept"`
	UserPermits           []string `json:"user_permits"`
	DeptPermits           []int64  `json:"dept_permits"`
}

type ListUsersParams struct {
	DeptId             int64  `json:"dept_id"`
	Cursor             int64  `json:"cursor"`
	Size               int64  `json:"size"`
	OrderField         string `json:"order_field"`
	ContainAccessLimit bool   `json:"contain_access_limit"`
	Language           string `json:"language"`
}

type ListUsersResult struct {
	HasMore    bool              `json:"has_more"`
	NextCursor int64             `json:"next_cursor"`
	List       []*DepartmentUser `json:"list"`
}

type DepartmentUser struct {
	UserID           string          `json:"userid"`
	UnionID          string          `json:"unionid"`
	Name             string          `json:"name"`
	Avatar           string          `json:"avatar"`
	StateCode        string          `json:"state_code"`
	Mobile           string          `json:"mobile"`
	HideMobile       bool            `json:"hide_mobile"`
	Telephone        string          `json:"telephone"`
	JobNumber        string          `json:"jobNumber"`
	Title            string          `json:"title"`
	Email            string          `json:"email"`
	OrgEmail         string          `json:"org_email"`
	WorkPlace        string          `json:"work_place"`
	Remark           string          `json:"remark"`
	DeptIDList       pq.Int64Array   `json:"dept_id_list"`
	DeptOrder        decimal.Decimal `json:"dept_order"`
	Extension        string          `json:"extension"`
	HiredDate        int64           `json:"hired_date"`
	Active           bool            `json:"active"`
	Admin            bool            `json:"admin"`
	Boss             bool            `json:"boss"`
	Leader           bool            `json:"leader"`
	ExclusiveAccount bool            `json:"exclusive_account"`
}
