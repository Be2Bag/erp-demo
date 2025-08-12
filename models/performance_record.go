package models

import "time"

const CollectionPerformanceRecords = "performance_records"

type PerformanceRecord struct {
	ID          string    `bson:"_id" json:"id"`
	UserID      string    `bson:"user_id"`      // รหัสพนักงาน
	PeriodMonth int       `bson:"period_month"` // เดือนที่วัด (1-12)
	PeriodYear  int       `bson:"period_year"`  // ปีที่วัด
	ReviewerID  string    `bson:"reviewer_id"`  // คนประเมิน (หัวหน้าโดยตรง)
	Scores      []KPIItem `bson:"scores"`       // รายการ KPI แต่ละหัวข้อ
	Comment     string    `bson:"comment"`      // ความคิดเห็นเพิ่มเติม
	TotalScore  float64   `bson:"total_score"`  // คะแนนรวม
	CreatedAt   time.Time `bson:"created_at"`
}

type KPIItem struct {
	Category string  `bson:"category"`  // หมวด: design, install, teamwork, punctuality, quality
	Title    string  `bson:"title"`     // หัวข้อเช่น "ออกแบบตรงเวลา", "งานติดตั้งไม่ผิดพลาด"
	Score    float64 `bson:"score"`     // คะแนนเต็ม (0-5 หรือ 0-10)
	Weight   float64 `bson:"weight"`    // น้ำหนักคะแนนของหัวข้อนี้ (เช่น 0.2 หมายถึง 20%)
	MaxScore float64 `bson:"max_score"` // คะแนนเต็มจริง
	Comment  string  `bson:"comment"`   // คำอธิบายเพิ่มเติม (optional)
}
