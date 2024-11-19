package service

import (
	"fmt"
	"time"

	"tools-admin/backend/model"
	"tools-admin/backend/pkg/cronutil"
	"tools-admin/backend/pkg/db"
	"tools-admin/backend/pkg/log"
)

type TaskService struct{}

// List 获取任务列表
func (s *TaskService) List(page, pageSize int, name string, taskType int) ([]*model.TaskResponse, int64, error) {
	var tasks []*model.Task
	var total int64

	// 构建查询
	dbQuery := db.Db.Model(&model.Task{})

	// 添加查询条件
	if name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+name+"%")
	}
	if taskType != 0 {
		dbQuery = dbQuery.Where("type = ?", taskType)
	}

	// 获取总数
	if err := dbQuery.Count(&total).Error; err != nil {
		log.Error(fmt.Sprintf("获取任务总数失败: %v", err))
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := dbQuery.Offset(offset).Limit(pageSize).Order("id desc").Find(&tasks).Error; err != nil {
		log.Error(fmt.Sprintf("获取任务列表失败: %v", err))
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []*model.TaskResponse
	for _, task := range tasks {
		responses = append(responses, task.ToResponse())
	}

	return responses, total, nil
}

// Get 获取单个任务
func (s *TaskService) Get(id uint) (*model.TaskResponse, error) {
	var task model.Task
	if err := db.Db.First(&task, id).Error; err != nil {
		log.Error(fmt.Sprintf("获取任务失败, ID: %d, 错误: %v", id, err))
		return nil, err
	}
	return task.ToResponse(), nil
}

// GetByID 根据ID获取任务
func (s *TaskService) GetByID(id uint) (*model.Task, error) {
	var task model.Task
	if err := db.Db.First(&task, id).Error; err != nil {
		log.Error(fmt.Sprintf("获取任务失败, ID: %d, 错误: %v", id, err))
		return nil, err
	}
	return &task, nil
}

// Create 创建任务
func (s *TaskService) Create(task *model.Task) error {
	// 验证cron表达式
	if task.CronExpr != "" {
		if err := cronutil.ValidateCronExpr(task.CronExpr); err != nil {
			log.Error(fmt.Sprintf("无效的cron表达式: %v", err))
			return err
		}

		// 计算下次执行时间
		nextRunTime, err := cronutil.GetNextRunTime(task.CronExpr)
		if err != nil {
			log.Error(fmt.Sprintf("计算下次执行时间失败: %v", err))
			return err
		}
		task.NextRunTime = nextRunTime
	}

	if err := db.Db.Create(task).Error; err != nil {
		log.Error(fmt.Sprintf("创建任务失败: %v", err))
		return err
	}
	return nil
}

// Update 更新任务
func (s *TaskService) Update(task *model.Task) error {
	// 如果有cron表达式，重新计算下次执行时间
	if task.CronExpr != "" {
		nextRunTime, err := cronutil.GetNextRunTime(task.CronExpr)
		if err != nil {
			log.Error(fmt.Sprintf("计算下次执行时间失败: %v", err))
			return err
		}
		task.NextRunTime = nextRunTime
	}

	if err := db.Db.Save(task).Error; err != nil {
		log.Error(fmt.Sprintf("更新任务失败: %v", err))
		return err
	}
	return nil
}

// Delete 删除任务
func (s *TaskService) Delete(id uint) error {
	if err := db.Db.Delete(&model.Task{}, id).Error; err != nil {
		log.Error(fmt.Sprintf("删除任务失败, ID: %d, 错误: %v", id, err))
		return err
	}
	return nil
}

// BatchDelete 批量删除任务
func (s *TaskService) BatchDelete(ids []uint) error {
	if err := db.Db.Delete(&model.Task{}, ids).Error; err != nil {
		log.Error(fmt.Sprintf("批量删除任务失败, IDs: %v, 错误: %v", ids, err))
		return err
	}
	return nil
}

// GetLogs 获取任务日志
func (s *TaskService) GetLogs(taskID uint) ([]*model.TaskLogResponse, error) {
	var logs []*model.TaskLog
	if err := db.Db.Where("task_id = ?", taskID).Order("id desc").Find(&logs).Error; err != nil {
		log.Error(fmt.Sprintf("获取任务日志失败, TaskID: %d, 错误: %v", taskID, err))
		return nil, err
	}

	// 转换为响应格式
	var responses []*model.TaskLogResponse
	for _, log := range logs {
		responses = append(responses, log.ToResponse())
	}
	return responses, nil
}

// RunTask 运行任务
func (s *TaskService) RunTask(id uint) error {
	// 获取任务
	task, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// 检查任务状态
	if task.Status != model.TaskStatusStarted {
		return fmt.Errorf("任务未启动，无法执行")
	}

	// 更新执行状态为执行中
	task.ExecStatus = model.TaskExecStatusRunning
	if err := s.Update(task); err != nil {
		return err
	}

	// TODO: 实际执行任务的逻辑

	return nil
}

// UpdateTaskStatus 更新任务状态
func (s *TaskService) UpdateTaskStatus(taskID int64, status model.TaskStatus) error {
	// 验证状态转换的合法性
	if status != model.TaskStatusStarted && status != model.TaskStatusStopped {
		return fmt.Errorf("无效的状态值")
	}

	// 更新状态
	result := db.Db.Model(&model.Task{}).Where("id = ?", taskID).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("更新任务状态失败")
	}

	return nil
}

// BatchUpdateTaskStatus 批量更新任务状态
func (s *TaskService) BatchUpdateTaskStatus(ids []uint, status model.TaskStatus) error {
	// 验证状态转换的合法性
	if status != model.TaskStatusStarted && status != model.TaskStatusStopped {
		return fmt.Errorf("无效的状态值")
	}

	// 批量更新状态
	result := db.Db.Model(&model.Task{}).Where("id IN ?", ids).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("批量更新任务状态失败")
	}

	return nil
}

// GetNextRunTimes 根据cron表达式获取未来4次执行时间
func (s *TaskService) GetNextRunTimes(cronExpr string) ([]string, error) {
	if err := cronutil.ValidateCronExpr(cronExpr); err != nil {
		return nil, err
	}

	times := make([]string, 0, 4)
	current := time.Now()

	for i := 0; i < 4; i++ {
		nextTime, err := cronutil.GetNextRunTimeFrom(cronExpr, current)
		if err != nil {
			return nil, err
		}
		if nextTime == nil {
			break
		}
		times = append(times, nextTime.Format("2006-01-02 15:04:05"))
		current = *nextTime
	}

	return times, nil
}
