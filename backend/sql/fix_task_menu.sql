-- 检查任务管理的父菜单
SELECT * FROM menus WHERE id = 30;

-- 添加任务列表子菜单
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
VALUES 
(30, 'TaskList', 'index', '@/views/task/index.vue', 2, 'List', 1, 1, 1, 1);

-- 检查任务生成子菜单是否存在
SELECT * FROM menus WHERE parent_id = 30 AND path LIKE 'generate%';

-- 如果任务生成子菜单不存在，添加它们
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
VALUES 
(30, 'TaskGenerate', 'generate', '', 1, 'Plus', 2, 1, 1, 1);

-- 获取任务生成菜单的ID
SET @generate_id = LAST_INSERT_ID();

-- 添加任务生成的子菜单
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
VALUES 
(@generate_id, 'DataxTask', 'datax', '@/views/task/generate/datax.vue', 2, 'DataAnalysis', 1, 1, 1, 1),
(@generate_id, 'ShellTask', 'shell', '@/views/task/generate/shell.vue', 2, 'Terminal', 2, 1, 1, 1),
(@generate_id, 'HttpTask', 'http', '@/views/task/generate/http.vue', 2, 'Link', 3, 1, 1, 1);

-- 检查任务日志菜单是否存在
SELECT * FROM menus WHERE parent_id = 30 AND path = 'log';

-- 如果任务日志菜单不存在，添加它
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
VALUES 
(30, 'TaskLog', 'log', '@/views/task/log/index.vue', 2, 'Document', 3, 1, 1, 1);
