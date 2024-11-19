package model

import "time"

type MenuMeta struct {
	Title     string   `json:"title"`
	Icon      string   `json:"icon,omitempty"`
	NoCache   bool     `json:"noCache,omitempty"`
	Breadcrumb bool    `json:"breadcrumb,omitempty"`
	ActiveMenu string  `json:"activeMenu,omitempty"`
}

type Menu struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	ParentID   uint      `json:"parentId" gorm:"column:parent_id;default:0;comment:父菜单ID"`
	Name       string    `json:"name" gorm:"comment:菜单名称"`
	Path       string    `json:"path" gorm:"comment:路由路径"`
	Component  string    `json:"component" gorm:"comment:组件路径"`
	Permission string    `json:"permission" gorm:"comment:权限标识"`
	Type       int8      `json:"type" gorm:"comment:类型(1:目录,2:菜单,3:按钮)"`
	Icon       string    `json:"icon" gorm:"comment:图标"`
	Sort       int       `json:"sort" gorm:"default:0;comment:排序"`
	Status     int8      `json:"status" gorm:"default:1;comment:状态(1:启用,0:禁用)"`
	Visible    int8      `json:"visible" gorm:"default:1;comment:是否可见(1:是,0:否)"`
	KeepAlive  int8      `json:"keepAlive" gorm:"column:keep_alive;default:1;comment:是否缓存(1:是,0:否)"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt  time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	Children   []Menu    `json:"children" gorm:"-"`
	Meta       MenuMeta  `json:"meta" gorm:"-"`
}

type MenuResponse struct {
	ID       uint           `json:"id"`
	ParentID uint           `json:"parentId"`
	Name     string         `json:"name"`
	Path     string         `json:"path"`
	Type     int8           `json:"type"`
	Icon     string         `json:"icon"`
	Children []MenuResponse `json:"children"`
	Meta     MenuMeta       `json:"meta"`
}

// TableName 设置表名
func (Menu) TableName() string {
	return "menus"
}

// AfterFind GORM的钩子函数，在从数据库加载后填充Meta字段
func (m *Menu) AfterFind() error {
	m.Meta = MenuMeta{
		Title:      m.Name,
		Icon:       m.Icon,
		NoCache:    m.KeepAlive == 0,
		Breadcrumb: true,
	}
	return nil
}

// ToResponse 将Menu转换为MenuResponse
func (m *Menu) ToResponse() MenuResponse {
	return MenuResponse{
		ID:       m.ID,
		ParentID: m.ParentID,
		Name:     m.Name,
		Path:     m.Path,
		Type:     m.Type,
		Icon:     m.Icon,
		Meta:     m.Meta,
		Children: make([]MenuResponse, 0),
	}
}
