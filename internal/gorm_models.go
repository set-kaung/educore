package internal

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Name       string `gorm:"column:name;not null"`
	CourseCode string `gorm:"column:course_code;not null,uniqueIndex"`
}

type Enrollment struct {
	gorm.Model
	StudentID        uint `gorm:"column:student_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	SemesterCourseID uint `gorm:"column:semester_course_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	Student          Student
	SemesterCourse   SemesterCourse
}

type Professor struct {
	gorm.Model
	Name         string `gorm:"column:name;not null"`
	DepartmentID uint   `gorm:"column:department_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	ADObjectID   string `gorm:"column:ad_object_id;type:varchar(255);uniqueIndex"`
	Department   Department
}

type SemesterCourseSchedule struct {
	gorm.Model
	SemesterCourseID uint       `gorm:"column:semester_course_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	From             time.Time  `gorm:"column:from;not null"`
	To               *time.Time `gorm:"column:to;null"`
	SemesterCourse   SemesterCourse
}

type SemesterCourse struct {
	gorm.Model
	Semester  string `gorm:"column:semester;not null"`
	CourseID  uint   `gorm:"column:course_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	Section   string `gorm:"column:section;not null"`
	TeachBy   uint   `gorm:"column:taught_by;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	Course    Course
	Professor Professor `gorm:"foreignKey:TeachBy"`
}

type Student struct {
	gorm.Model
	Username     string `gorm:"column:username;not null"`
	StudentID    string `gorm:"column:student_id;not null"`
	DepartmentID uint   `gorm:"column:department_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	ADObjectID   string `gorm:"column:ad_object_id;type:varchar(255);uniqueIndex"`
	Department   Department
}

type Token struct {
	gorm.Model
	UserID    uint      `gorm:"column:user_id;type:bigint unsigned;not null"`
	Token     string    `gorm:"column:token;type:varchar(255);not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
}

type ApiKey struct {
	gorm.Model
	Name  string `gorm:"column:name;type:varchar(255);not null"`
	Email string `gorm:"column:email;type:varchar(255);not null"`
	Key   string `gorm:"column:key;type:varchar(255);not null;uniqueIndex"`
	Valid bool   `gorm:"column:valid;not null;default:true"`
}

type Department struct {
	gorm.Model
	Name string `gorm:"column:name;type:varchar(255);not null;uniqueIndex"`
}

type SupportTicket struct {
	gorm.Model
	UserID  uint   `gorm:"column:user_id;type:bigint unsigned;not null"`
	Subject string `gorm:"column:subject;not null"`
	Status  string `gorm:"column:status;not null;default:'open'"`
}

var Models = []interface{}{
	&Course{},
	&Enrollment{},
	&Professor{},
	&SemesterCourseSchedule{},
	&SemesterCourse{},
	&Student{},
	&Token{},
	&ApiKey{},
	&Department{},
	&SupportTicket{},
}
