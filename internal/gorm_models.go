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

type User struct {
	gorm.Model
	Name         string `gorm:"column:name;not null"`
	DepartmentID uint   `gorm:"column:department_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	ADObjectID   string `gorm:"column:ad_object_id;type:varchar(255);uniqueIndex"`
	Role         string `gorm:"column:role;type:enum('student','professor');not null"`
	Username     string `gorm:"column:username"`
	StudentID    string `gorm:"column:student_id"`
	Department   Department
}

type Enrollment struct {
	gorm.Model
	UserID           uint `gorm:"column:user_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	SemesterCourseID uint `gorm:"column:semester_course_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	User             User
	SemesterCourse   SemesterCourse
}

type Semester struct {
	gorm.Model
	Name      string     `gorm:"column:name;type:varchar(255);not null;uniqueIndex"`
	IsCurrent bool       `gorm:"column:is_current;not null;default:false"`
	StartDate *time.Time `gorm:"column:start_date"`
	EndDate   *time.Time `gorm:"column:end_date"`
}

type SemesterCourseSchedule struct {
	gorm.Model
	SemesterCourseID uint      `gorm:"column:semester_course_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	Weekday          string    `gorm:"column:weekday;type:varchar(10);not null"`
	StartTime        time.Time `gorm:"column:start_time;type:time;not null"`
	EndTime          time.Time `gorm:"column:end_time;type:time;not null"`
	SemesterCourse   SemesterCourse
}

type SemesterCourse struct {
	gorm.Model
	SemesterID uint   `gorm:"column:semester_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	CourseID   uint   `gorm:"column:course_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	Section    string `gorm:"column:section;not null"`
	TeachBy    uint   `gorm:"column:taught_by;type:bigint unsigned;not null;constraint:OnDelete:CASCADE"`
	Semester   Semester
	Course     Course
	Professor  User `gorm:"foreignKey:TeachBy"`
}

type Token struct {
	gorm.Model
	UserID    uint      `gorm:"column:user_id;type:bigint unsigned;not null"`
	Token     string    `gorm:"column:token;type:varchar(255);not null,uniqueIndex"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
}

type ApiKey struct {
	gorm.Model
	Name  string `gorm:"column:name;type:varchar(255);not null"`
	Email string `gorm:"column:email;type:varchar(255);not null"`
	Key   string `gorm:"column:key;type:varchar(255);not null,uniqueIndex"`
	Valid bool   `gorm:"column:valid;not null;default:true"`
}

type Department struct {
	gorm.Model
	Name string `gorm:"column:name;type:varchar(255);not null,uniqueIndex"`
}

type RecommendedBook struct {
	gorm.Model
	SemesterCourseID uint   `gorm:"column:semester_course_id;type:bigint unsigned;not null;constraint:OnDelete:CASCADE;uniqueIndex:idx_semester_book"`
	Title            string `gorm:"column:title;not null"`
	Author           string `gorm:"column:author"`
	Isbn             string `gorm:"column:isbn"`
	CoverI           int    `gorm:"column:cover_i"`
	OpenLibraryKey   string `gorm:"column:open_library_key;type:varchar(255);uniqueIndex:idx_semester_book"`
	SemesterCourse   SemesterCourse
}

type SupportTicket struct {
	gorm.Model
	UserID  uint   `gorm:"column:user_id;type:bigint unsigned;not null"`
	Subject string `gorm:"column:subject;not null"`
	Status  string `gorm:"column:status;not null;default:'open'"`
}

var Models = []interface{}{
	&Course{},
	&User{},
	&Semester{},
	&Enrollment{},
	&SemesterCourseSchedule{},
	&SemesterCourse{},
	&Token{},
	&ApiKey{},
	&Department{},
	&SupportTicket{},
	&RecommendedBook{},
}
