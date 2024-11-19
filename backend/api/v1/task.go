package v1

import (
	"fmt"
	"strconv"
	"tools-admin/backend/model"
	"tools-admin/backend/pkg/cronutil"
	"tools-admin/backend/pkg/log"
	"tools-admin/backend/service"

	"github.com/gin-gonic/gin"
)

var taskService = &service.TaskService{}

// GetTasks 获取任务列表
func GetTasks(c *gin.Context) {
	// 如果是OPTIONS请求，直接返回
	if c.Request.Method == "OPTIONS" {
		c.Status(204)
		return
	}

	// 从查询参数获取数据
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	taskType, _ := strconv.Atoi(c.DefaultQuery("type", "0"))

	tasks, total, err := taskService.List(page, pageSize, name, taskType)
	if err != nil {
		log.Error(fmt.Sprintf("获取任务列表失败: %v", err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取任务列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "获取任务列表成功",
		"data": gin.H{
			"list":  tasks,
			"total": total,
		},
	})
}

// CreateTask 创建任务
func CreateTask(c *gin.Context) {
	var requestData struct {
		Name        string          `json:"name"`
		Type        interface{}     `json:"type"`
		Description string          `json:"description"`
		Status      interface{}     `json:"status"`
		Priority    string          `json:"priority"`
		CronExpr    string          `json:"cronExpr"`
		TaskContent string          `json:"taskContent"`
		TaskParams  string          `json:"taskParams"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 创建任务对象
	task := &model.Task{
		Name:        requestData.Name,
		Description: requestData.Description,
		Priority:    requestData.Priority,
		CronExpr:    requestData.CronExpr,
		TaskContent: requestData.TaskContent,
		TaskParams:  requestData.TaskParams,
		Status:      model.TaskStatusStopped,  // 默认为停止状态
		ExecStatus:  model.TaskExecStatusPending,  // 默认为待执行状态
	}

	// 处理类型字段
	switch v := requestData.Type.(type) {
	case float64:
		task.Type = model.TaskType(int8(v))
	case string:
		typeMap := map[string]int8{
			"shell":     1,
			"http":      2,
			"datax":     3,
			"regular":   4,
			"urgent":    5,
			"longterm":  6,
			"recurring": 7,
		}
		if typeNum, exists := typeMap[v]; exists {
			task.Type = model.TaskType(typeNum)
		} else {
			task.Type = model.TaskTypeShell // 默认为shell类型
		}
	default:
		task.Type = model.TaskTypeShell
	}

	// 处理状态字段
	if requestData.Status != nil {
		switch v := requestData.Status.(type) {
		case float64:
			task.Status = model.TaskStatus(int8(v))
		case string:
			statusMap := map[string]int8{
				"started": 1,
				"stopped": 2,
			}
			if statusNum, exists := statusMap[v]; exists {
				task.Status = model.TaskStatus(statusNum)
			}
		}
	}

	// 验证cron表达式
	if task.CronExpr != "" {
		if err := cronutil.ValidateCronExpr(task.CronExpr); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "无效的cron表达式",
				"error":   err.Error(),
			})
			return
		}
	}

	// 创建任务
	err := taskService.Create(task)
	if err != nil {
		log.Error(fmt.Sprintf("创建任务失败: %v", err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": "创建任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "创建任务成功",
		"data":    task.ToResponse(),
	})
}

// UpdateTask 更新任务
func UpdateTask(c *gin.Context) {
	// 如果是OPTIONS请求，直接返回
	if c.Request.Method == "OPTIONS" {
		c.Status(204)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "无效的任务ID",
		})
		return
	}

	// 获取当前任务
	task, err := taskService.GetByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "任务不存在",
			"error":   err.Error(),
		})
		return
	}

	// 绑定更新数据
	var updates model.Task
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "无效的请求参数",
			"error":   err.Error(),
		})
		return
	}

	// 更新任务字段
	task.Name = updates.Name
	task.Description = updates.Description
	task.Type = updates.Type
	task.Priority = updates.Priority
	task.Status = updates.Status
	task.CronExpr = updates.CronExpr
	task.TaskContent = updates.TaskContent
	task.TaskParams = updates.TaskParams

	// 保存更新
	if err := taskService.Update(task); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "更新任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "更新任务成功",
		"data":    task.ToResponse(),
	})
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	// 如果是OPTIONS请求，直接返回
	if c.Request.Method == "OPTIONS" {
		c.Status(204)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "无效的任务ID",
		})
		return
	}

	if err := taskService.Delete(uint(id)); err != nil {
		log.Error(fmt.Sprintf("删除任务失败: %v", err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": fmt.Sprintf("删除任务失败: %v", err),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "删除任务成功",
	})
}

// BatchDeleteTasks 批量删除任务
func BatchDeleteTasks(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 1,
			"error": "无效的请求参数",
		})
		return
	}

	if err := taskService.BatchDelete(req.IDs); err != nil {
		c.JSON(500, gin.H{
			"code": 1,
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"message": "删除成功",
	})
}

// GetTaskLogs 获取任务日志
func GetTaskLogs(c *gin.Context) {
	// 如果是OPTIONS请求，直接返回
	if c.Request.Method == "OPTIONS" {
		c.Status(204)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Error(fmt.Sprintf("获取任务日志ID参数错误: %v", err))
		c.JSON(400, gin.H{
			"code":    400,
			"message": "无效的任务ID",
			"error":   err.Error(),
		})
		return
	}

	logs, err := taskService.GetLogs(uint(id))
	if err != nil {
		log.Error(fmt.Sprintf("获取任务日志失败, ID: %d, 错误: %v", id, err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取任务日志失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "获取任务日志成功",
		"data":    logs,
	})
}

// RunTask 运行任务
func RunTask(c *gin.Context) {
	// 如果是OPTIONS请求，直接返回
	if c.Request.Method == "OPTIONS" {
		c.Status(204)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Error(fmt.Sprintf("运行任务ID参数错误: %v", err))
		c.JSON(400, gin.H{
			"code":    400,
			"message": "无效的任务ID",
			"error":   err.Error(),
		})
		return
	}

	if err := taskService.RunTask(uint(id)); err != nil {
		log.Error(fmt.Sprintf("运行任务失败, ID: %d, 错误: %v", id, err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": "运行任务失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "任务运行成功",
	})
}

// GetCommonCronPatterns 获取常用的cron表达式
func GetCommonCronPatterns(c *gin.Context) {
	// 如果是OPTIONS请求，直接返回
	if c.Request.Method == "OPTIONS" {
		c.Status(204)
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "获取常用cron表达式成功",
		"data":    cronutil.CommonPatterns,
	})
}

// UpdateTaskStatus 更新任务状态
func UpdateTaskStatus(c *gin.Context) {
	var requestData struct {
		Status int8 `json:"status"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "无效的任务ID",
			"error":   err.Error(),
		})
		return
	}

	// 获取当前任务
	task, err := taskService.GetByID(uint(taskID))
	if err != nil {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "任务不存在",
			"error":   err.Error(),
		})
		return
	}

	// 更新状态
	oldStatus := task.Status
	task.Status = model.TaskStatus(requestData.Status)

	// 如果任务从停止状态变为启动状态，将执行状态设置为待执行
	if oldStatus == model.TaskStatusStopped && task.Status == model.TaskStatusStarted {
		task.ExecStatus = model.TaskExecStatusPending
	}

	// 如果任务从启动状态变为停止状态，将执行状态设置为待执行
	if oldStatus == model.TaskStatusStarted && task.Status == model.TaskStatusStopped {
		task.ExecStatus = model.TaskExecStatusPending
	}

	if err := taskService.Update(task); err != nil {
		log.Error(fmt.Sprintf("更新任务状态失败: %v", err))
		c.JSON(500, gin.H{
			"code":    500,
			"message": "更新任务状态失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "更新任务状态成功",
		"data":    task.ToResponse(),
	})
}

// GetTaskById 获取单个任务详情
func GetTaskById(c *gin.Context) {
    // 从路径参数获取任务ID
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(400, gin.H{
            "code":    400,
            "message": "无效的任务ID",
            "error":   err.Error(),
        })
        return
    }

    // 获取任务详情
    task, err := taskService.GetByID(uint(id))
    if err != nil {
        log.Error(fmt.Sprintf("获取任务详情失败: %v", err))
        c.JSON(500, gin.H{
            "code":    500,
            "message": "获取任务详情失败",
            "error":   err.Error(),
        })
        return
    }

    if task == nil {
        c.JSON(404, gin.H{
            "code":    404,
            "message": "任务不存在",
        })
        return
    }

    c.JSON(200, gin.H{
        "code":    0,
        "message": "获取任务详情成功",
        "data":    task.ToResponse(),
    })
}

// GetNextRunTimes 获取下次执行时间
// @Summary 获取下次执行时间
// @Description 根据cron表达式获取未来4次执行时间
// @Tags 任务管理
// @Accept json
// @Produce json
// @Param cronExpr query string true "cron表达式"
// @Success 200 {array} string "执行时间列表"
// @Router /api/v1/tasks/next-run-times [get]
func GetNextRunTimes(c *gin.Context) {
	cronExpr := c.Query("cronExpr")
	if cronExpr == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "cron表达式不能为空",
		})
		return
	}

	times, err := taskService.GetNextRunTimes(cronExpr)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取下次执行时间失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "获取下次执行时间成功",
		"data":    times,
	})
}
