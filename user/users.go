package users

import (
	"fmt"
	"sync"

	"github.com/gozelle/dingtalk"
)

type IUser interface {
	GetUserOrNil(id string) (user *dingtalk.DepartmentUser, err error)
}

func NewUsers(api *dingtalk.Client) *Users {
	return &Users{api: api}
}

var _ IUser = (*Users)(nil)

type Users struct {
	lock    sync.Mutex
	api     *dingtalk.Client
	users   []*dingtalk.DepartmentUser
	mapping map[string]*dingtalk.DepartmentUser
}

func (s *Users) GetUserOrNil(id string) (user *dingtalk.DepartmentUser, err error) {

	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()

	if s.mapping == nil {
		s.users, err = s.getUserIds()
		if err != nil {
			return
		}
		s.mapping = map[string]*dingtalk.DepartmentUser{}
		for _, v := range s.users {
			s.mapping[v.UserID] = v
		}
	}

	user = s.mapping[id]

	return
}

func (s *Users) listAllSub(params *dingtalk.DepartmentListSubParams) (r []*dingtalk.Department, err error) {
	list, err := s.api.DepartmentClient().ListSub(params)
	if err != nil {
		return
	}
	for _, v := range list {
		r = append(r, v)
		rr, e := s.listAllSub(&dingtalk.DepartmentListSubParams{
			DeptID: v.DeptID,
		})
		if e != nil {
			err = e
			return
		}
		if len(rr) > 0 {
			r = append(r, rr...)
		}
	}
	return
}

func (s *Users) getUserIds() (users []*dingtalk.DepartmentUser, err error) {
	subs, err := s.listAllSub(nil)
	if err != nil {
		err = fmt.Errorf("查询部门列表错误: %s", err)
		return
	}
	// 统计根部门
	subs = append(subs, &dingtalk.Department{
		DeptID: 1,
	})
	var dUsers []*dingtalk.DepartmentUser
	for _, sub := range subs {
		r, e := s.api.DepartmentClient().ListUsers(dingtalk.ListUsersParams{
			DeptId: sub.DeptID,
			Cursor: 0,
			Size:   100,
		})
		if err = e; err != nil {
			err = fmt.Errorf("查询部门用户错误: %s", err)
			return
		}
		dUsers = append(dUsers, r.List...)
	}

	check := map[string]bool{}

	for _, v := range dUsers {
		if _, ok := check[v.UserID]; !ok {
			users = append(users, v)
			check[v.UserID] = true
		}
	}

	return
}
