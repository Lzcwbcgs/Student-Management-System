package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql" // 使用MySQL驱动
	"github.com/yourusername/student-management-system/internal/api/handler"
	"github.com/yourusername/student-management-system/internal/api/middleware"
	"github.com/yourusername/student-management-system/internal/repository"
	"github.com/yourusername/student-management-system/internal/service"
	"github.com/yourusername/student-management-system/pkg/config"
)

func main() {
	// 加载配置
	cfg, err := config.Load("../config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := sql.Open("mysql", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 设置数据库连接池参数
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	// 检查数据库连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 执行数据库初始化脚本
	if err := executeSQLScript(db, "../scripts/init.sql"); err != nil {
		log.Fatalf("Failed to execute init.sql: %v", err)
	}

	// 初始化仓库层
	studentRepo := repository.NewStudentRepository(db)
	instructorRepo := repository.NewInstructorRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	sectionRepo := repository.NewSectionRepository(db)
	takesRepo := repository.NewTakesRepository(db)
	advisorRepo := repository.NewAdvisorRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)
	classroomRepo := repository.NewClassroomRepository(db)
	timeSlotRepo := repository.NewTimeSlotRepository(db)
	teachesRepo := repository.NewTeachesRepository(db)
	prereqRepo := repository.NewPrereqRepository(db)

	// 初始化服务层
	studentService := service.NewStudentService(studentRepo, takesRepo, prereqRepo, sectionRepo, advisorRepo)
	instructorService := service.NewInstructorService(instructorRepo, teachesRepo, takesRepo, advisorRepo, sectionRepo, studentRepo)
	courseService := service.NewCourseService(courseRepo, prereqRepo, departmentRepo)
	sectionService := service.NewSectionService(sectionRepo, courseRepo, classroomRepo, timeSlotRepo)
	enrollmentService := service.NewEnrollmentService(takesRepo, studentRepo, sectionRepo, courseRepo, prereqRepo, timeSlotRepo, teachesRepo)
	adminService := service.NewAdminService(studentRepo, instructorRepo, courseRepo, sectionRepo, departmentRepo, classroomRepo, timeSlotRepo, teachesRepo, advisorRepo, prereqRepo)

	// 初始化认证中间件
	authMiddleware := middleware.NewAuthMiddleware()

	// 初始化处理器
	studentHandler := handler.NewStudentHandler(studentService)
	instructorHandler := handler.NewInstructorHandler(instructorService)
	courseHandler := handler.NewCourseHandler(courseService)
	sectionHandler := handler.NewSectionHandler(sectionService)
	registrationHandler := handler.NewRegistrationHandler(enrollmentService) // 注意这里改为enrollmentService
	adminHandler := handler.NewAdminHandler(adminService)
	authHandler := handler.NewAuthHandler(studentService, instructorService)

	// 创建路由
	mux := http.NewServeMux()

	// 添加 CORS 中间件
	corsMux := middleware.CORSMiddleware(mux)

	// 认证路由
	mux.HandleFunc("/api/login", authHandler.Login)

	// 学生路由
	mux.HandleFunc("/api/students/profile", authMiddleware.Authenticate(authMiddleware.AuthorizeStudent(studentHandler.GetProfile)))
	mux.HandleFunc("/api/students/profile/update", authMiddleware.Authenticate(authMiddleware.AuthorizeStudent(studentHandler.UpdateProfile)))
	mux.HandleFunc("/api/students/advisor", authMiddleware.Authenticate(authMiddleware.AuthorizeStudent(studentHandler.GetAdvisor)))
	mux.HandleFunc("/api/students/courses", authMiddleware.Authenticate(authMiddleware.AuthorizeStudent(studentHandler.GetCourses)))
	mux.HandleFunc("/api/students/transcript", authMiddleware.Authenticate(authMiddleware.AuthorizeStudent(studentHandler.GetTranscript)))

	// 课程和选课路由
	mux.HandleFunc("/api/courses", authMiddleware.Authenticate(courseHandler.GetCourses))
	mux.HandleFunc("/api/sections", authMiddleware.Authenticate(sectionHandler.GetSections))
	mux.HandleFunc("/api/registration/register", authMiddleware.Authenticate(authMiddleware.AuthorizeStudent(registrationHandler.RegisterCourse)))
	mux.HandleFunc("/api/registration/drop", authMiddleware.Authenticate(authMiddleware.AuthorizeStudent(registrationHandler.DropCourse)))

	// 教师路由
	mux.HandleFunc("/api/instructors/profile", authMiddleware.Authenticate(authMiddleware.AuthorizeInstructor(instructorHandler.GetProfile)))
	mux.HandleFunc("/api/instructors/profile/update", authMiddleware.Authenticate(authMiddleware.AuthorizeInstructor(instructorHandler.UpdateProfile)))
	mux.HandleFunc("/api/instructors/sections", authMiddleware.Authenticate(authMiddleware.AuthorizeInstructor(instructorHandler.GetSections)))
	mux.HandleFunc("/api/instructors/sections/students", authMiddleware.Authenticate(authMiddleware.AuthorizeInstructor(instructorHandler.GetSectionStudents)))
	mux.HandleFunc("/api/instructors/grade/update", authMiddleware.Authenticate(authMiddleware.AuthorizeInstructor(instructorHandler.UpdateGrade)))
	mux.HandleFunc("/api/instructors/advisees", authMiddleware.Authenticate(authMiddleware.AuthorizeInstructor(instructorHandler.GetAdvisees)))
	mux.HandleFunc("/api/instructors/advisees/info", authMiddleware.Authenticate(authMiddleware.AuthorizeInstructor(instructorHandler.GetAdviseeInfo)))

	// 管理员路由
	mux.HandleFunc("/api/admin/students", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetStudents)))
	mux.HandleFunc("/api/admin/students/create", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.CreateStudent)))
	mux.HandleFunc("/api/admin/students/update", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.UpdateStudent)))
	mux.HandleFunc("/api/admin/students/delete", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.DeleteStudent)))
	mux.HandleFunc("/api/admin/instructors", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetInstructors)))
	mux.HandleFunc("/api/admin/instructors/create", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.CreateInstructor)))
	mux.HandleFunc("/api/admin/instructors/update", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.UpdateInstructor)))
	mux.HandleFunc("/api/admin/instructors/delete", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.DeleteInstructor)))
	mux.HandleFunc("/api/admin/departments", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetDepartments)))
	mux.HandleFunc("/api/admin/departments/create", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.CreateDepartment)))
	mux.HandleFunc("/api/admin/departments/update", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.UpdateDepartment)))
	mux.HandleFunc("/api/admin/departments/delete", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.DeleteDepartment)))
	mux.HandleFunc("/api/admin/courses", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetCourses)))
	mux.HandleFunc("/api/admin/courses/create", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.CreateCourse)))
	mux.HandleFunc("/api/admin/courses/update", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.UpdateCourse)))
	mux.HandleFunc("/api/admin/courses/delete", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.DeleteCourse)))
	mux.HandleFunc("/api/admin/prereqs", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetPrereqs)))
	mux.HandleFunc("/api/admin/prereqs/create", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.CreatePrereq)))
	mux.HandleFunc("/api/admin/prereqs/delete", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.DeletePrereq)))
	mux.HandleFunc("/api/admin/classrooms", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetClassrooms)))
	mux.HandleFunc("/api/admin/classrooms/create", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.CreateClassroom)))
	mux.HandleFunc("/api/admin/classrooms/update", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.UpdateClassroom)))
	mux.HandleFunc("/api/admin/classrooms/delete", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.DeleteClassroom)))
	mux.HandleFunc("/api/admin/sections", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetSections)))
	mux.HandleFunc("/api/admin/sections/create", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.CreateSection)))
	mux.HandleFunc("/api/admin/sections/update", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.UpdateSection)))
	mux.HandleFunc("/api/admin/sections/delete", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.DeleteSection)))
	mux.HandleFunc("/api/admin/teaches", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetTeaches)))
	mux.HandleFunc("/api/admin/teaches/create", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.CreateTeaches)))
	mux.HandleFunc("/api/admin/teaches/delete", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.DeleteTeaches)))
	mux.HandleFunc("/api/admin/advisors", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetAdvisors)))
	mux.HandleFunc("/api/admin/advisors/create", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.CreateAdvisor)))
	mux.HandleFunc("/api/admin/advisors/delete", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.DeleteAdvisor)))
	mux.HandleFunc("/api/admin/stats", authMiddleware.Authenticate(authMiddleware.AuthorizeAdmin(adminHandler.GetStats)))

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: corsMux,
	}

	// 启动服务器
	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// executeSQLScript 读取并执行 SQL 脚本文件
func executeSQLScript(db *sql.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL script file %s: %w", filePath, err)
	}

	sqls := string(content)
	// 分割 SQL 语句，这里简单地按分号分割，实际项目中可能需要更复杂的解析
	statements := splitSQLStatements(sqls)

	for _, stmt := range statements {
		stmt = trimSpaceAndNewline(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute SQL statement: %s\nError: %w", stmt, err)
		}
	}
	log.Printf("Successfully executed SQL script: %s", filePath)
	return nil
}

// splitSQLStatements 辅助函数：分割 SQL 语句
func splitSQLStatements(sqlContent string) []string {
	// 简单的分割逻辑，不处理复杂的 SQL 语法，例如包含分号的字符串字面量
	return strings.Split(sqlContent, ";")
}

// trimSpaceAndNewline 辅助函数：去除字符串首尾的空格和换行符
func trimSpaceAndNewline(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\n\r\t "))
}
