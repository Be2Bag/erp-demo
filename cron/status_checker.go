package cron

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Be2Bag/erp-demo/pkg/util"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// StatusUpdateResult เก็บผลลัพธ์การอัปเดตสถานะแต่ละรายการ
type StatusUpdateResult struct {
	Type      string  `json:"type"`       // "payable" หรือ "receivable"
	ID        string  `json:"id"`         // ID ของรายการ
	InvoiceNo string  `json:"invoice_no"` // เลขที่ใบแจ้งหนี้
	OldStatus string  `json:"old_status"` // สถานะเดิม
	NewStatus string  `json:"new_status"` // สถานะใหม่
	Amount    float64 `json:"amount"`     // จำนวนเงินทั้งหมด
	Balance   float64 `json:"balance"`    // ยอดคงเหลือ
	DueDate   string  `json:"due_date"`   // วันครบกำหนด
}

// CronRunSummary สรุปผลการรัน cronjob
type CronRunSummary struct {
	RunAt              string               `json:"run_at"`              // เวลาที่รัน
	TotalPayables      int                  `json:"total_payables"`      // จำนวน Payable ที่ตรวจสอบ
	UpdatedPayables    int                  `json:"updated_payables"`    // จำนวน Payable ที่อัปเดต
	TotalReceivables   int                  `json:"total_receivables"`   // จำนวน Receivable ที่ตรวจสอบ
	UpdatedReceivables int                  `json:"updated_receivables"` // จำนวน Receivable ที่อัปเดต
	UpdatedItems       []StatusUpdateResult `json:"updated_items"`       // รายการที่อัปเดต
}

// StatusChecker รับผิดชอบในการตรวจสอบและอัปเดตสถานะของ Payable และ Receivable
type StatusChecker struct {
	payableRepo    ports.PayableRepository
	receivableRepo ports.ReceivableRepository
	cron           *cron.Cron
	lastRunSummary *CronRunSummary // เก็บผลลัพธ์การรันล่าสุด
}

// NewStatusChecker สร้าง StatusChecker ใหม่
func NewStatusChecker(
	payableRepo ports.PayableRepository,
	receivableRepo ports.ReceivableRepository,
) *StatusChecker {
	// ใช้ timezone ไทย (Asia/Bangkok, GMT+7)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		// fallback ถ้าโหลด timezone ไม่ได้
		loc = time.FixedZone("Asia/Bangkok", 7*60*60)
	}

	return &StatusChecker{
		payableRepo:    payableRepo,
		receivableRepo: receivableRepo,
		cron:           cron.New(cron.WithLocation(loc)),
	}
}

// Start เริ่มต้น cronjob
// จะรันทุกวันเวลา 00:00 น. (เที่ยงคืน) ตามเวลาไทย (GMT+7)
func (sc *StatusChecker) Start() error {
	// รันทุกวันเวลา 00:00 น. ตามเวลาไทย (Asia/Bangkok)
	// รูปแบบ: minute hour day month weekday (5 ฟิลด์)
	_, err := sc.cron.AddFunc("0 0 * * *", func() {
		log.Println("[CRON] เริ่มตรวจสอบสถานะ Payable และ Receivable...")

		summary, err := sc.runStatusCheck()
		if err != nil {
			log.Printf("[CRON ERROR] ตรวจสอบสถานะไม่สำเร็จ: %v", err)
			return
		}

		sc.lastRunSummary = summary
		sc.printSummary(summary, "[CRON]")

		log.Println("[CRON] ตรวจสอบสถานะเสร็จสิ้น")
	})

	if err != nil {
		return err
	}

	sc.cron.Start()
	log.Println("[CRON] Status Checker เริ่มทำงานแล้ว (รันทุกวันเวลา 00:00 น. ตามเวลาไทย)")

	return nil
}

// Stop หยุด cronjob
func (sc *StatusChecker) Stop() {
	log.Println("[CRON] หยุด Status Checker...")
	sc.cron.Stop()
}

// checkPayableStatus ตรวจสอบและอัปเดตสถานะของ Payable
func (sc *StatusChecker) checkPayableStatus(ctx context.Context) ([]StatusUpdateResult, int, error) {
	now := time.Now()
	var updatedItems []StatusUpdateResult

	// ดึงข้อมูล Payable ทั้งหมดที่ยังไม่ถูกลบ และยังไม่ได้ชำระครบ
	filter := bson.M{
		"deleted_at": nil,
		"status": bson.M{
			"$in": []string{"pending", "partial"},
		},
		"balance": bson.M{
			"$gte": 0.01, // มากกว่า 1 สตางค์
		},
	}

	payables, err := sc.payableRepo.GetAllPayablesByFilter(ctx, filter, nil)
	if err != nil {
		return nil, 0, err
	}

	for _, payable := range payables {
		needUpdate := false
		oldStatus := payable.Status

		// ตรวจสอบว่าเลยกำหนดชำระหรือไม่
		if !payable.DueDate.IsZero() && payable.DueDate.Before(now) {
			// เลยกำหนดและยังมียอดคงเหลือ (มากกว่า 1 สตางค์)
			if util.IsPositiveAmount(payable.Balance) {
				payable.Status = "overdue"
				needUpdate = true
			}
		} else {
			// ยังไม่เลยกำหนด
			if util.IsPositiveAmount(payable.Balance) && payable.Balance < payable.Amount {
				// จ่ายบางส่วน
				payable.Status = "partial"
				needUpdate = true
			} else if payable.Balance == payable.Amount {
				// ยังไม่จ่ายเลย
				if payable.Status != "pending" {
					payable.Status = "pending"
					needUpdate = true
				}
			}
		}

		// อัปเดตถ้ามีการเปลี่ยนแปลง
		if needUpdate && oldStatus != payable.Status {
			payable.UpdatedAt = now
			if _, err := sc.payableRepo.UpdatePayableByID(ctx, payable.IDPayable, *payable); err != nil {
				log.Printf("[ERROR] อัปเดต Payable %s ไม่สำเร็จ: %v", payable.IDPayable, err)
				continue
			}

			// เพิ่มข้อมูลลงในรายการที่อัปเดต
			dueDateStr := ""
			if !payable.DueDate.IsZero() {
				dueDateStr = payable.DueDate.Format("02/01/2006")
			}

			updatedItems = append(updatedItems, StatusUpdateResult{
				Type:      "payable",
				ID:        payable.IDPayable,
				InvoiceNo: payable.InvoiceNo,
				OldStatus: oldStatus,
				NewStatus: payable.Status,
				Amount:    payable.Amount,
				Balance:   payable.Balance,
				DueDate:   dueDateStr,
			})

			log.Printf("[UPDATE] Payable %s: %s → %s", payable.InvoiceNo, oldStatus, payable.Status)
		}
	}

	return updatedItems, len(payables), nil
}

// checkReceivableStatus ตรวจสอบและอัปเดตสถานะของ Receivable
func (sc *StatusChecker) checkReceivableStatus(ctx context.Context) ([]StatusUpdateResult, int, error) {
	now := time.Now()
	var updatedItems []StatusUpdateResult

	// ดึงข้อมูล Receivable ทั้งหมดที่ยังไม่ถูกลบ และยังไม่ได้รับชำระครบ
	filter := bson.M{
		"deleted_at": nil,
		"status": bson.M{
			"$in": []string{"pending", "partial"},
		},
		"balance": bson.M{
			"$gte": 0.01, // มากกว่า 1 สตางค์
		},
	}

	receivables, err := sc.receivableRepo.GetAllReceivablesByFilter(ctx, filter, nil)
	if err != nil {
		return nil, 0, err
	}

	for _, receivable := range receivables {
		needUpdate := false
		oldStatus := receivable.Status

		// ตรวจสอบว่าเลยกำหนดรับชำระหรือไม่
		if !receivable.DueDate.IsZero() && receivable.DueDate.Before(now) {
			// เลยกำหนดและยังมียอดคงเหลือ (มากกว่า 1 สตางค์)
			if util.IsPositiveAmount(receivable.Balance) {
				receivable.Status = "overdue"
				needUpdate = true
			}
		} else {
			// ยังไม่เลยกำหนด
			if util.IsPositiveAmount(receivable.Balance) && receivable.Balance < receivable.Amount {
				// รับชำระบางส่วน
				receivable.Status = "partial"
				needUpdate = true
			} else if receivable.Balance == receivable.Amount {
				// ยังไม่รับชำระเลย
				if receivable.Status != "pending" {
					receivable.Status = "pending"
					needUpdate = true
				}
			}
		}

		// อัปเดตถ้ามีการเปลี่ยนแปลง
		if needUpdate && oldStatus != receivable.Status {
			receivable.UpdatedAt = now
			if _, err := sc.receivableRepo.UpdateReceivableByID(ctx, receivable.IDReceivable, *receivable); err != nil {
				log.Printf("[ERROR] อัปเดต Receivable %s ไม่สำเร็จ: %v", receivable.IDReceivable, err)
				continue
			}

			// เพิ่มข้อมูลลงในรายการที่อัปเดต
			dueDateStr := ""
			if !receivable.DueDate.IsZero() {
				dueDateStr = receivable.DueDate.Format("02/01/2006")
			}

			updatedItems = append(updatedItems, StatusUpdateResult{
				Type:      "receivable",
				ID:        receivable.IDReceivable,
				InvoiceNo: receivable.InvoiceNo,
				OldStatus: oldStatus,
				NewStatus: receivable.Status,
				Amount:    receivable.Amount,
				Balance:   receivable.Balance,
				DueDate:   dueDateStr,
			})

			log.Printf("[UPDATE] Receivable %s: %s → %s", receivable.InvoiceNo, oldStatus, receivable.Status)
		}
	}

	return updatedItems, len(receivables), nil
}

// runStatusCheck รันการตรวจสอบสถานะและส่งคืนผลสรุป
func (sc *StatusChecker) runStatusCheck() (*CronRunSummary, error) {
	ctx := context.Background()
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)

	summary := &CronRunSummary{
		RunAt:        now.Format("02/01/2006 15:04:05"),
		UpdatedItems: []StatusUpdateResult{},
	}

	// ตรวจสอบ Payables
	payableUpdates, totalPayables, err := sc.checkPayableStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("ตรวจสอบ Payable ไม่สำเร็จ: %v", err)
	}
	summary.TotalPayables = totalPayables
	summary.UpdatedPayables = len(payableUpdates)
	summary.UpdatedItems = append(summary.UpdatedItems, payableUpdates...)

	// ตรวจสอบ Receivables
	receivableUpdates, totalReceivables, err := sc.checkReceivableStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("ตรวจสอบ Receivable ไม่สำเร็จ: %v", err)
	}
	summary.TotalReceivables = totalReceivables
	summary.UpdatedReceivables = len(receivableUpdates)
	summary.UpdatedItems = append(summary.UpdatedItems, receivableUpdates...)

	return summary, nil
}

// printSummary แสดงผลสรุปการรัน
func (sc *StatusChecker) printSummary(summary *CronRunSummary, prefix string) {
	log.Printf("%s ========== สรุปผลการตรวจสอบ ==========", prefix)
	log.Printf("%s เวลาที่รัน: %s", prefix, summary.RunAt)
	log.Printf("%s Payable: อัปเดต %d จาก %d รายการ", prefix, summary.UpdatedPayables, summary.TotalPayables)
	log.Printf("%s Receivable: อัปเดต %d จาก %d รายการ", prefix, summary.UpdatedReceivables, summary.TotalReceivables)

	if len(summary.UpdatedItems) > 0 {
		log.Printf("%s ---------- รายการที่อัปเดต ----------", prefix)
		for i, item := range summary.UpdatedItems {
			log.Printf("%s %d. [%s] %s: %s → %s (ยอด: %.2f, คงเหลือ: %.2f, ครบกำหนด: %s)",
				prefix, i+1, item.Type, item.InvoiceNo, item.OldStatus, item.NewStatus,
				item.Amount, item.Balance, item.DueDate)
		}
	} else {
		log.Printf("%s ไม่มีรายการที่ต้องอัปเดต", prefix)
	}
	log.Printf("%s ======================================", prefix)
}

// GetLastRunSummary คืนค่าผลสรุปการรันล่าสุด
func (sc *StatusChecker) GetLastRunSummary() *CronRunSummary {
	return sc.lastRunSummary
}

// RunNow รันการตรวจสอบทันที (สำหรับ testing หรือรันด้วยตนเอง)
func (sc *StatusChecker) RunNow() (*CronRunSummary, error) {
	log.Println("[MANUAL] เริ่มตรวจสอบสถานะ Payable และ Receivable...")

	summary, err := sc.runStatusCheck()
	if err != nil {
		log.Printf("[MANUAL ERROR] ตรวจสอบสถานะไม่สำเร็จ: %v", err)
		return nil, err
	}

	sc.lastRunSummary = summary
	sc.printSummary(summary, "[MANUAL]")

	log.Println("[MANUAL] ตรวจสอบสถานะเสร็จสิ้น")

	return summary, nil
}
