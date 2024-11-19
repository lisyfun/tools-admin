-- 修复任务列表菜单的路径
UPDATE menus 
SET path = '', component = '@/views/task/index.vue'
WHERE parent_id = 30 AND name = 'TaskList';

-- 如果上面的更新没有影响任何行（说明记录不存在），则插入新记录
INSERT INTO menus (parent_id, name, path, component, type, icon, sort, status, visible, keep_alive)
SELECT 30, 'TaskList', '', '@/views/task/index.vue', 2, 'List', 1, 1, 1, 1
WHERE NOT EXISTS (
    SELECT 1 FROM menus WHERE parent_id = 30 AND name = 'TaskList'
);
