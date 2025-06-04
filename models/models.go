// package models

// import (
// 	"time"
// )

// // AcademicGroups соответствует таблице Academic_Groups
// type AcademicGroup struct {
// 	AcademicGroupID int32     `gorm:"primaryKey"`
// 	Name            string    `gorm:"type:varchar(255);not null"`
// 	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// }

// // Groups соответствует таблице Groups
// type Group struct {
// 	GroupID         int32     `gorm:"primaryKey"`
// 	Name            string    `gorm:"type:varchar(255);not null"`
// 	AcademicGroupID int32     `gorm:"foreignKey:AcademicGroupID;references:AcademicGroups(AcademicGroupID)"`
// 	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	// AcademicGroup   AcademicGroup
// 	AdminID User `gorm:"foreignKey:AdminID;references:UserID"`
// }

// // Модераторы группы
// type GroupModer struct {
// 	GroupID int32 `gorm:"primaryKey"`
// 	UserId  int32 `gorm:"primaryKey"`
// }

// // Users соответствует таблице Users
// type User struct {
// 	UserID       int32     `gorm:"primaryKey"`
// 	Username     string    `gorm:"type:varchar(255);not null"`
// 	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	Email        string    `gorm:"type:varchar(255); not null"`
// 	HashPassword string    `gorm:"type:varchar(255); not null"`
// }

// // GroupUsers соответствует таблице Group_Users
// type GroupUser struct {
// 	GroupID  int32     `gorm:"primaryKey;foreignKey:GroupID;references:Groups(GroupID)"`
// 	UserID   int32     `gorm:"primaryKey;foreignKey:UserID;references:Users(UserID)"`
// 	Role     string    `gorm:"type:varchar(50);default:'member'"`
// 	JoinedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	Group    Group
// 	User     User
// }

// // Subjects соответствует таблице Subjects
// type Subject struct {
// 	SubjectID   int32     `gorm:"primaryKey"`
// 	GroupID     int32     `gorm:"foreignKey:GroupID;references:Groups(GroupID)"`
// 	Name        string    `gorm:"type:varchar(255);not null"`
// 	Description string    `gorm:"type:text"`
// 	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	Group       Group
// }

// // Tasks соответствует таблице Tasks
// type Task struct {
// 	TaskID      int32     `gorm:"primaryKey"`
// 	SubjectID   int32     `gorm:"foreignKey:SubjectID;references:Subjects(SubjectID)"`
// 	Title       string    `gorm:"type:varchar(255);not null"`
// 	Description string    `gorm:"type:text"`
// 	CreatedBy   int32     `gorm:"foreignKey:CreatedBy;references:Users(UserID)"`
// 	IsActive    bool      `gorm:"default:true"`
// 	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	Subject     Subject
// 	Creator     User
// }

// // TaskUsers соответствует таблице Task_Users
// type TaskUser struct {
// 	TaskID      int32      `gorm:"primaryKey;foreignKey:TaskID;references:Tasks(TaskID)"`
// 	UserID      int32      `gorm:"primaryKey;foreignKey:UserID;references:Users(UserID)"`
// 	IsCompleted bool       `gorm:"default:false"`
// 	CompletedAt *time.Time // Используем указатель для nullable поля
// 	Task        Task
// 	User        User
// }

// // Materials соответствует таблице Materials
// type Material struct {
// 	MaterialID int32     `gorm:"primaryKey"`
// 	SubjectID  int32     `gorm:"foreignKey:SubjectID;references:Subjects(SubjectID)"`
// 	Title      string    `gorm:"type:varchar(255);not null"`
// 	Content    string    `gorm:"type:text"`
// 	CreatedBy  int32     `gorm:"foreignKey:CreatedBy;references:Users(UserID)"`
// 	IsActive   bool      `gorm:"default:true"`
// 	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	Subject    Subject
// 	Creator    User
// }

// // TimeSlots соответствует таблице Time_Slots
// type TimeSlot struct {
// 	SlotID     int32     `gorm:"primaryKey"`
// 	SlotNumber int32     `gorm:"unique;check:slot_number BETWEEN 1 AND 9"`
// 	StartTime  time.Time `gorm:"type:time;not null"`
// 	EndTime    time.Time `gorm:"type:time;not null"`
// }

// // Schedules соответствует таблице Schedules
// type Schedule struct {
// 	ScheduleID      int32     `gorm:"primaryKey"`
// 	GroupID         int32     `gorm:"foreignKey:GroupID;references:Groups(GroupID)"`
// 	SubjectID       int32     `gorm:"foreignKey:SubjectID;references:Subjects(SubjectID)"`
// 	TeacherInitials string    `gorm:"type:varchar(50)"`
// 	Classroom       string    `gorm:"type:varchar(50)"`
// 	SlotID          int32     `gorm:"foreignKey:SlotID;references:TimeSlots(SlotID)"`
// 	Date            time.Time `gorm:"type:date;not null"`
// 	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	Group           Group
// 	Subject         Subject
// 	TimeSlot        TimeSlot
// }

package models

import (
	"time"
)

// AcademicGroup
type AcademicGroup struct {
	AcademicGroupID int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"type:varchar(255);not null"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Group
type OldGroup struct {
	// GroupID         int32         `gorm:"primaryKey"`
	// Name            string        `gorm:"type:varchar(255);not null"`
	// AcademicGroupID int32         `gorm:"foreignKey:AcademicGroupID;references:AcademicGroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	// CreatedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
	// AdminID         int32         `gorm:"foreignKey:AdminID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	// AcademicGroup   AcademicGroup `gorm:"foreignKey:AcademicGroupID"`
	// Admin           User          `gorm:"foreignKey:AdminID"`
	ID              int32 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string
	CreatedAt       time.Time
	AcademicGroupID int32
	AdminID         int32
	AcademicGroup   AcademicGroup
	Admin           User
}

type Group struct {
	ID              int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"type:varchar(255)" json:"name"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	AcademicGroupID int32     `gorm:"index" json:"academic_group_id"` // Foreign key to AcademicGroup
	AdminID         int32     `gorm:"index"`                          // Foreign key to User (admin)

	AcademicGroup AcademicGroup

	Admin User
}

// GroupModer
type OldGroupModer struct {
	GroupID int32 `gorm:"primaryKey;foreignKey:GroupID;references:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID  int32 `gorm:"primaryKey;foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type GroupModer struct {
	GroupID   int32 `gorm:"primaryKey;foreignKey:GroupID"`
	UserID    int32 `gorm:"primaryKey;foreignKey:UserID"`
	Group     Group `gorm:"foreignKey:GroupID"`
	User      User  `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
}

// Users
type User struct {
	UserID       int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Email        string    `gorm:"type:varchar(255);not null"`
	HashPassword string    `gorm:"type:varchar(255);not null"`
}

// GroupUsers
type GroupUser struct {
	GroupID  int32     `gorm:"primaryKey;foreignKey:GroupID;references:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID   int32     `gorm:"primaryKey;foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role     string    `gorm:"type:varchar(50);default:'member'"`
	JoinedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Group    Group     `gorm:"foreignKey:GroupID"`
	User     User      `gorm:"foreignKey:UserID"`
}

// Subjects
type Subject struct {
	SubjectID       int32     `gorm:"primaryKey"`
	AcademicGroupID int32     `gorm:"foreignKey:AcademicGroupID;references:AcademicGroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Name            string    `gorm:"type:varchar(255);not null"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	AcademicGroup   AcademicGroup
}
type OldSubject struct {
	SubjectID int32  `gorm:"primaryKey"`
	GroupID   int32  `gorm:"foreignKey:GroupID;references:GroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Name      string `gorm:"type:varchar(255);not null"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Group     Group     `gorm:"foreignKey:GroupID"`
}

// Tasks
//
//		type Task struct {
//			TaskID      int32  `gorm:"primaryKey"`
//			SubjectID   int32  `gorm:"foreignKey:SubjectID;references:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
//			Title       string `gorm:"type:varchar(255);not null"`
//			Description string `gorm:"type:text"`
//			// CreatedBy   int32     `gorm:"foreignKey:CreatedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
//			CreaterID int32
//			IsActive  bool      `gorm:"default:true"`
//			CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
//			UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
//			Subject   Subject   `gorm:"foreignKey:SubjectID"`
//			// Creator   User      `gorm:"foreignKey:CreatedBy"`
//			Creator User
//	}
type OldTask struct {
	TaskID      int32     `gorm:"column:task_id;primaryKey;autoIncrement"`
	SubjectID   int32     `gorm:"foreignKey:SubjectID;references:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Title       string    `gorm:"column:title;type:varchar(255);not null"`
	Description string    `gorm:"column:description;type:text"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	CreatedBy   int       `gorm:"column:created_by"`
	IsActive    bool      `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
	// Subject     Subject   `gorm:"foreignKey:SubjectID"`
	Subject Subject `json:"-"`

	User User `gorm:"foreignKey:CreatedBy;references:user_id" json:"-"`
}
type Task struct {
	ID          int32     `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID     int32     `gorm:"not null" json:"group_id"`
	UserID      int32     `gorm:"not null" json:"user_id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	IsVerified  bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"-"`
	Group       Group     `json:"-"`
}

// Materials
type Material struct {
	MaterialID int32     `gorm:"primaryKey"`
	SubjectID  int32     `gorm:"foreignKey:SubjectID;references:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Title      string    `gorm:"type:varchar(255);not null"`
	Content    string    `gorm:"type:text"`
	CreatedBy  int32     `gorm:"foreignKey:CreatedBy;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Subject    Subject
	Creator    User `gorm:"foreignKey:CreatedBy"`
	// Subject    Subject   `gorm:"foreignKey:SubjectID"`
}

// TimeSlots
type TimeSlot struct {
	SlotID     int32     `gorm:"primaryKey"`
	SlotNumber int32     `gorm:"unique;check:slot_number BETWEEN 1 AND 9"`
	StartTime  time.Time `gorm:"type:time;not null"`
	EndTime    time.Time `gorm:"type:time;not null"`
}

// Schedules
type Schedule struct {
	ScheduleID int32 `gorm:"primaryKey"`
	GroupID    int32 `gorm:"foreignKey:GroupID;references:GroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	// SubjectID       int32  `gorm:"foreignKey:SubjectID;references:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	SubjectID int32 `gorm:"foreignKey:SubjectID;references:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	TeacherInitials string `gorm:"type:varchar(50)"`
	Classroom       string `gorm:"type:varchar(50)"`
	TimeSlotID      int32

	Date      time.Time `gorm:"type:date;not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Group Group `gorm:"foreignKey:GroupID"`
	// Subject Subject `gorm:"foreignKey:SubjectID"`
	Subject Subject

	// Subject  Subject `gorm:"foreignKey:SubjectID"`
	TimeSlot TimeSlot
}
type GroupApplication struct {
	ApplicationID int32     `gorm:"primaryKey;autoIncrement" json:"application_id"`
	GroupID       int32     `gorm:"not null" json:"group_id"`
	UserID        int32     `gorm:"not null" json:"user_id"`
	Message       string    `gorm:"type:text" json:"message"`
	Status        string    `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, approved, rejected
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Group Group
	User  User
}
