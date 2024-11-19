-- 更新主菜单的中文名称
UPDATE menus 
SET name = '数据面板' 
WHERE path = '/dashboard' AND parent_id = 0;

UPDATE menus 
SET name = '任务管理' 
WHERE path = '/task' AND parent_id = 0;

UPDATE menus 
SET name = '短信管理' 
WHERE path = '/sms' AND parent_id = 0;

UPDATE menus 
SET name = '流水线管理' 
WHERE path = '/pipeline' AND parent_id = 0;

UPDATE menus 
SET name = '服务器管理' 
WHERE path = '/server' AND parent_id = 0;

UPDATE menus 
SET name = '数据库管理' 
WHERE path = '/database' AND parent_id = 0;

UPDATE menus 
SET name = '系统管理' 
WHERE path = '/system' AND parent_id = 0;

-- 更新任务管理子菜单
UPDATE menus 
SET name = '任务列表' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/task' AND parent_id = 0) AS t) 
AND path = '';

UPDATE menus 
SET name = '任务生成' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/task' AND parent_id = 0) AS t) 
AND path = 'generate';

UPDATE menus 
SET name = '任务日志' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/task' AND parent_id = 0) AS t) 
AND path = 'log';

-- 更新任务生成子菜单
UPDATE menus 
SET name = 'DataX任务' 
WHERE path = 'datax';

UPDATE menus 
SET name = 'Shell任务' 
WHERE path = 'shell';

UPDATE menus 
SET name = 'HTTP任务' 
WHERE path = 'http';

-- 更新短信管理子菜单
UPDATE menus 
SET name = '短信列表' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/sms' AND parent_id = 0) AS t) 
AND path = 'index';

UPDATE menus 
SET name = '短信模板' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/sms' AND parent_id = 0) AS t) 
AND path = 'template';

UPDATE menus 
SET name = '收件人管理' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/sms' AND parent_id = 0) AS t) 
AND path = 'recipient';

-- 更新流水线管理子菜单
UPDATE menus 
SET name = '容器构建' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/pipeline' AND parent_id = 0) AS t) 
AND path = 'build';

UPDATE menus 
SET name = '容器部署' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/pipeline' AND parent_id = 0) AS t) 
AND path = 'deploy';

-- 更新服务器管理子菜单
UPDATE menus 
SET name = '服务器列表' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/server' AND parent_id = 0) AS t) 
AND path = 'list';

UPDATE menus 
SET name = '服务器新增' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/server' AND parent_id = 0) AS t) 
AND path = 'add';

UPDATE menus 
SET name = '服务器日志' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/server' AND parent_id = 0) AS t) 
AND path = 'log';

-- 更新数据库管理子菜单
UPDATE menus 
SET name = '数据库连接' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/database' AND parent_id = 0) AS t) 
AND path = 'connection';

UPDATE menus 
SET name = '数据库操作' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/database' AND parent_id = 0) AS t) 
AND path = 'operation';

-- 更新系统管理子菜单
UPDATE menus 
SET name = '用户管理' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/system' AND parent_id = 0) AS t) 
AND path = 'user';

UPDATE menus 
SET name = '菜单管理' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/system' AND parent_id = 0) AS t) 
AND path = 'menu';

UPDATE menus 
SET name = '角色管理' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/system' AND parent_id = 0) AS t) 
AND path = 'role';

UPDATE menus 
SET name = '操作日志' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/system' AND parent_id = 0) AS t) 
AND path = 'operation-log';

UPDATE menus 
SET name = '系统日志' 
WHERE parent_id = (SELECT id FROM (SELECT id FROM menus WHERE path = '/system' AND parent_id = 0) AS t) 
AND path = 'system-log';
