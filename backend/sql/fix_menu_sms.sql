-- 删除重复的菜单项
DELETE FROM menus WHERE id = 56;

-- 确保使用ID 36作为父菜单
UPDATE menus SET name = '短信管理', icon = 'Message' WHERE id = 36;

-- 添加子菜单，使用ID 36作为父ID
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
VALUES 
(36, 'SmsList', 'index', '@/views/sms/index.vue', 2, 'List', 1, 1, 1, 1),
(36, 'SmsTemplate', 'template', '@/views/sms/template/index.vue', 2, 'Document', 2, 1, 1, 1),
(36, 'SmsRecipient', 'recipient', '@/views/sms/recipient/index.vue', 2, 'User', 3, 1, 1, 1);
