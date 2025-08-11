-- 创建数据库
CREATE DATABASE IF NOT EXISTS student_management;
USE student_management;

-- 创建院系表
CREATE TABLE IF NOT EXISTS department (
    dept_name VARCHAR(20) PRIMARY KEY,
    building VARCHAR(15),
    budget DECIMAL(12,2)
);

-- 创建学生表
CREATE TABLE IF NOT EXISTS student (
    ID VARCHAR(5) PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    dept_name VARCHAR(20),
    tot_cred DECIMAL(3,0) DEFAULT 0,
    password VARCHAR(100),
    salt VARCHAR(50),
    FOREIGN KEY (dept_name) REFERENCES department(dept_name)
);

-- 创建教师表
CREATE TABLE IF NOT EXISTS instructor (
    ID VARCHAR(5) PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    dept_name VARCHAR(20),
    salary DECIMAL(8,2),
    FOREIGN KEY (dept_name) REFERENCES department(dept_name)
);

-- 创建课程表
CREATE TABLE IF NOT EXISTS course (
    course_id VARCHAR(8) PRIMARY KEY,
    title VARCHAR(50),
    dept_name VARCHAR(20),
    credits DECIMAL(2,0),
    FOREIGN KEY (dept_name) REFERENCES department(dept_name)
);

-- 创建教室表
CREATE TABLE IF NOT EXISTS classroom (
    building VARCHAR(15),
    room_number VARCHAR(7),
    capacity DECIMAL(4,0),
    PRIMARY KEY (building, room_number)
);

-- 创建时间段表
CREATE TABLE IF NOT EXISTS time_slot (
    time_slot_id VARCHAR(4) PRIMARY KEY,
    day VARCHAR(1),
    start_hr DECIMAL(2),
    start_min DECIMAL(2),
    end_hr DECIMAL(2),
    end_min DECIMAL(2)
);

-- 创建课程段表
CREATE TABLE IF NOT EXISTS section (
    course_id VARCHAR(8),
    sec_id VARCHAR(8),
    semester VARCHAR(6),
    year DECIMAL(4,0),
    building VARCHAR(15),
    room_number VARCHAR(7),
    time_slot_id VARCHAR(4),
    PRIMARY KEY (course_id, sec_id, semester, year),
    FOREIGN KEY (course_id) REFERENCES course(course_id),
    FOREIGN KEY (building, room_number) REFERENCES classroom(building, room_number),
    FOREIGN KEY (time_slot_id) REFERENCES time_slot(time_slot_id)
);

-- 创建学生选课表
CREATE TABLE IF NOT EXISTS takes (
    ID VARCHAR(5),
    course_id VARCHAR(8),
    sec_id VARCHAR(8),
    semester VARCHAR(6),
    year DECIMAL(4,0),
    grade VARCHAR(2),
    PRIMARY KEY (ID, course_id, sec_id, semester, year),
    FOREIGN KEY (ID) REFERENCES student(ID),
    FOREIGN KEY (course_id, sec_id, semester, year) REFERENCES section(course_id, sec_id, semester, year)
);

-- 创建教师授课表
CREATE TABLE IF NOT EXISTS teaches (
    ID VARCHAR(5),
    course_id VARCHAR(8),
    sec_id VARCHAR(8),
    semester VARCHAR(6),
    year DECIMAL(4,0),
    PRIMARY KEY (ID, course_id, sec_id, semester, year),
    FOREIGN KEY (ID) REFERENCES instructor(ID),
    FOREIGN KEY (course_id, sec_id, semester, year) REFERENCES section(course_id, sec_id, semester, year)
);

-- 创建导师关系表
CREATE TABLE IF NOT EXISTS advisor (
    s_ID VARCHAR(5),
    i_ID VARCHAR(5),
    PRIMARY KEY (s_ID),
    FOREIGN KEY (s_ID) REFERENCES student(ID),
    FOREIGN KEY (i_ID) REFERENCES instructor(ID)
);

-- 创建先修课程表
CREATE TABLE IF NOT EXISTS prereq (
    course_id VARCHAR(8),
    prereq_id VARCHAR(8),
    PRIMARY KEY (course_id, prereq_id),
    FOREIGN KEY (course_id) REFERENCES course(course_id),
    FOREIGN KEY (prereq_id) REFERENCES course(course_id)
);

-- 插入示例数据
INSERT IGNORE INTO department VALUES ('计算机科学', '工程楼', 100000.00);
INSERT IGNORE INTO department VALUES ('数学', '科学楼', 80000.00);
INSERT IGNORE INTO department VALUES ('物理', '科学楼', 90000.00);

-- 插入学生数据（包含密码和盐值字段，示例密码统一为 '123456'）
INSERT IGNORE INTO student VALUES
('S001', '张三', '计算机科学', 30, '$2a$12$s94zBo0.hs6z5qLIQVTueuP/U8Zm0rDYzGq/n.2Mm2pNRcZNWgJ6u', 's3cr3ts4lt'),
('S002', '李四', '数学', 25, '$2a$10$xJwL5v5z3V1lO6B9QYqZNuYbU1wYk7Xe7n6jKJc8bLm0v1aG2sD1C', 's4ltv4lu3'),
('S003', '王五', '计算机科学', 35, '$2a$10$xJwL5v5z3V1lO6B9QYqZNuYbU1wYk7Xe7n6jKJc8bLm0v1aG2sD1C', 's0m3s4lt');

INSERT IGNORE INTO instructor VALUES ('I001', '陈教授', '计算机科学', 80000.00);
INSERT IGNORE INTO instructor VALUES ('I002', '刘教授', '数学', 75000.00);
INSERT IGNORE INTO instructor VALUES ('I003', '赵教授', '物理', 85000.00);

INSERT IGNORE INTO course VALUES ('CS101', '计算机科学导论', '计算机科学', 4);
INSERT IGNORE INTO course VALUES ('CS102', '数据结构', '计算机科学', 4);
INSERT IGNORE INTO course VALUES ('MATH101', '微积分', '数学', 3);

INSERT IGNORE INTO classroom VALUES ('工程楼', '101', 50);
INSERT IGNORE INTO classroom VALUES ('工程楼', '102', 40);
INSERT IGNORE INTO classroom VALUES ('科学楼', '201', 60);

INSERT IGNORE INTO time_slot VALUES ('A', 'M', 8, 0, 8, 50);
INSERT IGNORE INTO time_slot VALUES ('B', 'M', 9, 0, 9, 50);
INSERT IGNORE INTO time_slot VALUES ('C', 'T', 10, 0, 10, 50);

INSERT IGNORE INTO section VALUES ('CS101', '1', 'Fall', 2024, '工程楼', '101', 'A');
INSERT IGNORE INTO section VALUES ('CS102', '1', 'Fall', 2024, '工程楼', '102', 'B');
INSERT IGNORE INTO section VALUES ('MATH101', '1', 'Fall', 2024, '科学楼', '201', 'C');

INSERT IGNORE INTO teaches VALUES ('I001', 'CS101', '1', 'Fall', 2024);
INSERT IGNORE INTO teaches VALUES ('I001', 'CS102', '1', 'Fall', 2024);
INSERT IGNORE INTO teaches VALUES ('I002', 'MATH101', '1', 'Fall', 2024);

INSERT IGNORE INTO advisor VALUES ('S001', 'I001');
INSERT IGNORE INTO advisor VALUES ('S002', 'I002');
INSERT IGNORE INTO advisor VALUES ('S003', 'I001');

INSERT IGNORE INTO prereq VALUES ('CS102', 'CS101');