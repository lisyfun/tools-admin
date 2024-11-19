package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tools-admin/backend/model"
	"tools-admin/backend/service"
)

// GetMenus 获取菜单列表
func GetMenus(c *gin.Context) {
	menuService := service.MenuService{}
	menus, err := menuService.GetMenuTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取菜单列表成功",
		"data":    menus,
	})
}

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	menuService := service.MenuService{}
	if err := menuService.CreateMenu(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建菜单成功",
		"data":    menu,
	})
}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	menuService := service.MenuService{}
	if err := menuService.UpdateMenu(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新菜单成功",
		"data":    menu,
	})
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的菜单ID",
		})
		return
	}

	menuService := service.MenuService{}
	if err := menuService.DeleteMenu(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除菜单成功",
	})
}
