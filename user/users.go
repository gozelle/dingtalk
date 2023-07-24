package users

import (
	"fmt"
	"sync"
	
	"github.com/gozelle/dingtalk"
)

type IManger interface {
	GetUser(id string) (user *dingtalk.DepartmentUser, err error)
	Users() (users []*dingtalk.DepartmentUser, err error)
}

func NewUserManger(api *dingtalk.Client) *Manager {
	return &Manager{api: api}
}

var _ IManger = (*Manager)(nil)

type Manager struct {
	lock    sync.Mutex
	api     *dingtalk.Client
	mapping map[string]*dingtalk.DepartmentUser
}

func (s *Manager) GetUser(id string) (user *dingtalk.DepartmentUser, err error) {
	
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()
	
	err = s.initUsers()
	if err != nil {
		return
	}
	var ok bool
	user, ok = s.mapping[id]
	if !ok {
		user = nil
		return
	}
	
	return
}

func (s *Manager) initUsers() (err error) {
	if s.mapping != nil {
		return
	}
	
	users, err := s.getUserIds()
	if err != nil {
		return
	}
	s.mapping = map[string]*dingtalk.DepartmentUser{}
	for _, v := range users {
		s.mapping[v.UserID] = v
	}
	
	return
}

func (s *Manager) Users() (users []*dingtalk.DepartmentUser, err error) {
	
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()
	
	err = s.initUsers()
	if err != nil {
		return
	}
	
	for _, v := range s.mapping {
		users = append(users, v)
	}
	
	return
}

func (s *Manager) listAllSub(params *dingtalk.DepartmentListSubParams) (r []*dingtalk.Department, err error) {
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

func (s *Manager) getUserIds() (users []*dingtalk.DepartmentUser, err error) {
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
