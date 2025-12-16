package config

import (
	"database/sql"

	"student-report/app/repository"
	"student-report/app/service"

	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceContainer struct {
	AuthService        *service.AuthService
	UserService        *service.UserService
	StudentService     *service.StudentService
	LecturerService    *service.LecturerService
	AchievementService *service.AchievementService
}

func InitializeServices(db *sql.DB, mongoDB *mongo.Database) *ServiceContainer {
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	lecturerRepo := repository.NewLecturerRepository(db)

	return &ServiceContainer{
		AuthService:     service.NewAuthService(authRepo),
		UserService:     service.NewUserService(userRepo),
		StudentService:  service.NewStudentService(studentRepo),
		LecturerService: service.NewLecturerService(lecturerRepo),
	}
}
