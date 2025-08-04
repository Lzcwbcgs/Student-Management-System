-- 修改 student 表添加认证相关字段
ALTER TABLE student
ADD COLUMN password VARCHAR(100) NOT NULL DEFAULT '',
ADD COLUMN salt VARCHAR(32) NOT NULL DEFAULT '';

-- 修改 instructor 表添加认证相关字段
ALTER TABLE instructor
ADD COLUMN password VARCHAR(100) NOT NULL DEFAULT '',
ADD COLUMN salt VARCHAR(32) NOT NULL DEFAULT '';

-- 创建一个初始管理员账号
INSERT INTO users (id, username, password, salt, role)
VALUES ('admin', 'admin', 'hashed_password', 'salt', 'admin');
