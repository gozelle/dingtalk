package process

import "sort"

type Instance struct {
	AttachedProcessInstanceIds []interface{}                 `json:"attached_process_instance_ids"`
	BizAction                  string                        `json:"biz_action"`
	BusinessId                 string                        `json:"business_id"`
	CreateTime                 string                        `json:"create_time"`
	FormComponentValues        []*InstanceFormComponentValue `json:"form_component_values"`
	OperationRecords           InstanceOperationRecords      `json:"operation_records"`
	OriginatorDeptId           string                        `json:"originator_dept_id"`
	OriginatorDeptName         string                        `json:"originator_dept_name"`
	OriginatorUserid           string                        `json:"originator_userid"`
	Result                     string                        `json:"result"`
	Status                     string                        `json:"status"`
	Tasks                      InstanceTasks                 `json:"tasks"`
	Title                      string                        `json:"title"`
}

type InstanceFormComponentValue struct {
	ComponentType string `json:"component_type"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	ExtValue      string `json:"ext_value,omitempty"`
}

var _ sort.Interface = (*InstanceTasks)(nil)

type InstanceTasks []*InstanceTask

func (p InstanceTasks) Len() int {
	return len(p)
}

func (p InstanceTasks) Less(i, j int) bool {
	return p[i].CreateTime < p[j].CreateTime
}

func (p InstanceTasks) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type InstanceTask struct {
	ActivityId       string                   `json:"activity_id"`
	CreateTime       string                   `json:"create_time"`
	FinishTime       string                   `json:"finish_time,omitempty"`
	TaskResult       string                   `json:"task_result"`
	TaskStatus       string                   `json:"task_status"`
	Taskid           string                   `json:"taskid"`
	Url              string                   `json:"url"`
	Userid           string                   `json:"userid"`
	Username         string                   `json:"username"`
	OperationRecords InstanceOperationRecords `json:"operation_records"`
}

var _ sort.Interface = (*InstanceOperationRecords)(nil)

type InstanceOperationRecords []*InstanceOperationRecord

func (p InstanceOperationRecords) Len() int {
	return len(p)
}

func (p InstanceOperationRecords) Less(i, j int) bool {
	return p[i].Date < p[j].Date
}

func (p InstanceOperationRecords) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type InstanceOperationRecord struct {
	Date            string `json:"date"`
	OperationResult string `json:"operation_result"`
	OperationType   string `json:"operation_type"`
	Userid          string `json:"userid"`
	Username        string `json:"username"`
	Remark          string `json:"remark,omitempty"`
}
