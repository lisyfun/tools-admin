-- 添加短信管理父菜单
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
VALUES (0, 'SMS', '/sms', 'Layout', 1, 'Message', 4, 1, 1, 1);

-- 获取刚插入的父菜单ID
SET @parent_id = LAST_INSERT_ID();

-- 添加短信列表子菜单
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
VALUES (@parent_id, 'SmsList', 'index', '@/views/sms/index.vue', 2, 'List', 1, 1, 1, 1);

-- 添加短信模板子菜单
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
VALUES (@parent_id, 'SmsTemplate', 'template', '@/views/sms/template/index.vue', 2, 'Document', 2, 1, 1, 1);

-- 添加收件人子菜单
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
VALUES (@parent_id, 'SmsRecipient', 'recipient', '@/views/sms/recipient/index.vue', 2, 'User', 3, 1, 1, 1);
