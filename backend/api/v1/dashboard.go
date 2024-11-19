package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tools-admin/backend/service"
)

type DashboardApi struct {
	dashboardService service.DashboardService
}

// GetOverview 获取仪表盘概览数据
// @Summary 获取仪表盘概览数据
// @Description 获取今日统计数据和同比趋势
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} model.DashboardOverview
// @Router /api/v1/dashboard/overview [get]
func (api *DashboardApi) GetOverview(c *gin.Context) {
	data, code := api.dashboardService.GetOverview()
	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  "获取概览数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// GetTaskChart 获取任务图表数据
// @Summary 获取任务图表数据
// @Description 获取任务成功失败趋势数据
// @Tags dashboard
// @Accept json
// @Produce json
// @Param period query string true "时间周期(week/month)"
// @Success 200 {array} model.ChartData
// @Router /api/v1/dashboard/task-chart [get]
func (api *DashboardApi) GetTaskChart(c *gin.Context) {
	period := c.DefaultQuery("period", "week")
	data, code := api.dashboardService.GetTaskChart(period)
	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  "获取任务图表数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// GetSmsChart 获取短信图表数据
// @Summary 获取短信图表数据
// @Description 获取短信发送成功率趋势数据
// @Tags dashboard
// @Accept json
// @Produce json
// @Param period query string true "时间周期(week/month)"
// @Success 200 {array} model.ChartData
// @Router /api/v1/dashboard/sms-chart [get]
func (api *DashboardApi) GetSmsChart(c *gin.Context) {
	period := c.DefaultQuery("period", "week")
	data, code := api.dashboardService.GetSmsChart(period)
	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  "获取短信图表数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}
