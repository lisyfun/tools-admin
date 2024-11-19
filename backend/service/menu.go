package service

import (
	"tools-admin/backend/model"
	"tools-admin/backend/pkg/db"
)

type MenuService struct{}

// GetMenuTree 获取菜单树
func (s *MenuService) GetMenuTree() ([]model.MenuResponse, error) {
	var allMenus []model.Menu
	if err := db.Db.Order("sort").Find(&allMenus).Error; err != nil {
		return nil, err
	}

	// 初始化所有菜单的Children为空数组
	menuResponses := make([]model.MenuResponse, len(allMenus))
	for i := range allMenus {
		menuResponses[i] = allMenus[i].ToResponse()
	}

	// 构建菜单树
	return buildMenuTree(menuResponses, 0), nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []model.MenuResponse, parentID uint) []model.MenuResponse {
    var tree []model.MenuResponse
    for _, menu := range menus {
        if menu.ParentID == parentID {
            // 如果是目录类型，则递归构建子菜单
            if menu.Type == 1 {
                menu.Children = buildMenuTree(menus, menu.ID)
            } else {
                // 如果是菜单或按钮类型，则不包含子菜单
                menu.Children = nil
            }
            tree = append(tree, menu)
        }
    }
    return tree
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	// 确保新菜单的Children字段为空数组而不是nil
	menu.Children = make([]model.Menu, 0)
	return db.Db.Create(menu).Error
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(menu *model.Menu) error {
	// 确保更新时Children字段为空数组而不是nil
	menu.Children = make([]model.Menu, 0)
	return db.Db.Save(menu).Error
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(id uint) error {
	return db.Db.Delete(&model.Menu{}, id).Error
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := db.Db.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	// 确保返回的菜单Children字段为空数组而不是nil
	menu.Children = make([]model.Menu, 0)
	return &menu, err
}
