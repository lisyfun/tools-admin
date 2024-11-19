export interface MenuMeta {
  title?: string
  icon?: string
  noCache?: boolean
  breadcrumb?: boolean
  affix?: boolean
  activeMenu?: string
  roles?: string[]
}

export interface MenuItem {
  id?: number
  parentId?: number
  name: string
  path: string
  component?: string
  permission?: string
  type?: number
  icon?: string
  sort?: number
  status?: number
  visible?: number
  keepAlive?: number
  createdAt?: string
  updatedAt?: string
  deletedAt?: string
  children?: MenuItem[]
  meta?: MenuMeta
}

export interface MenuState {
  menuList: MenuItem[]
  routes: MenuItem[]
}
