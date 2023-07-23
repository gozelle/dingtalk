package dingtalk_test

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/gozelle/dingtalk"
	"github.com/stretchr/testify/require"
)

func TestDepartment(t *testing.T) {

	// 测试
	client := NewTestClient()

	subs, err := ListAllSub(client, nil)
	require.NoError(t, err)

	subs = append([]*dingtalk.Department{{DeptID: 1}}, subs...) // 加入根部门

	var users []*dingtalk.DepartmentUser
	for _, sub := range subs {
		r, e := client.DepartmentClient().ListUsers(dingtalk.ListUsersParams{
			DeptId: sub.DeptID,
			Cursor: 0,
			Size:   100,
		})
		require.NoError(t, e)
		users = append(users, r.List...)
	}

	var _users []*dingtalk.DepartmentUser
	check := map[string]bool{}

	for _, v := range users {
		if _, ok := check[v.UserID]; !ok {
			_users = append(_users, v)
			check[v.UserID] = true
		}
	}

	d, err := json.MarshalIndent(_users, "", "\t")
	require.NoError(t, err)
	//t.Log(string(d))

	f, err := os.Create("users.json")
	require.NoError(t, err)

	defer func() {
		_ = f.Close()
	}()
	_, err = f.Write(d)
	require.NoError(t, err)

	_ = writeCSV(_users)
}

func writeCSV(users []*dingtalk.DepartmentUser) (err error) {
	f, err := os.Create("users.csv") //创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	data := [][]string{}
	for i, v := range users {
		data = append(data, []string{fmt.Sprintf("%d", i+1), v.Name, v.Avatar})
	}

	w.WriteAll(data) //写入数据
	w.Flush()
	return
}

func ListAllSub(client *dingtalk.Client, params *dingtalk.DepartmentListSubParams) (r []*dingtalk.Department, err error) {

	list, err := client.DepartmentClient().ListSub(params)
	if err != nil {
		return
	}
	for _, v := range list {
		r = append(r, v)
		rr, e := ListAllSub(client, &dingtalk.DepartmentListSubParams{
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
