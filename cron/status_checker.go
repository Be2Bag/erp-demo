package cron

import (
	"context"
	"log"
	"time"

	"github.com/Be2Bag/erp-demo/ports"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// StatusChecker รับผิดชอบในการตรวจสอบและอัปเดตสถานะของ Payable และ Receivable
type StatusChecker struct {
	payableRepo    ports.PayableRepository
	receivableRepo ports.ReceivableRepository
	cron           *cron.Cron
}

// NewStatusChecker สร้าง StatusChecker ใหม่
func NewStatusChecker(
	payableRepo ports.PayableRepository,
	receivableRepo ports.ReceivableRepository,
) *StatusChecker {
	return &StatusChecker{
		payableRepo:    payableRepo,
		receivableRepo: receivableRepo,
		cron:           cron.New(),
	}
}

// Start เริ่มต้น cronjob
// จะรันทุกวันเวลา 00:00 น. (เที่ยงคืน)
func (sc *StatusChecker) Start() error {
	// รันทุกวันเวลา 00:00 น.
	// รูปแบบ: minute hour day month weekday (5 ฟิลด์)
	_, err := sc.cron.AddFunc("0 0 * * *", func() {
		log.Println("[CRON] เริ่มตรวจสอบสถานะ Payable และ Receivable...")

		ctx := context.Background()

		// ตรวจสอบ Payables
		if err := sc.checkPayableStatus(ctx); err != nil {
			log.Printf("[CRON ERROR] ตรวจสอบ Payable ไม่สำเร็จ: %v", err)
		}

		// ตรวจสอบ Receivables
		if err := sc.checkReceivableStatus(ctx); err != nil {
			log.Printf("[CRON ERROR] ตรวจสอบ Receivable ไม่สำเร็จ: %v", err)
		}

		log.Println("[CRON] ตรวจสอบสถานะเสร็จสิ้น")
	})

	if err != nil {
		return err
	}

	sc.cron.Start()
	log.Println("[CRON] Status Checker เริ่มทำงานแล้ว (รันทุกวันเวลา 00:00 น.)")

	return nil
}

// Stop หยุด cronjob
func (sc *StatusChecker) Stop() {
	log.Println("[CRON] หยุด Status Checker...")
	sc.cron.Stop()
}

// checkPayableStatus ตรวจสอบและอัปเดตสถานะของ Payable
func (sc *StatusChecker) checkPayableStatus(ctx context.Context) error {
	now := time.Now()

	// ดึงข้อมูล Payable ทั้งหมดที่ยังไม่ถูกลบ และยังไม่ได้ชำระครบ
	filter := bson.M{
		"deleted_at": nil,
		"status": bson.M{
			"$in": []string{"pending", "partial"},
		},
		"balance": bson.M{
			"$gt": 0,
		},
	}

	payables, err := sc.payableRepo.GetAllPayablesByFilter(ctx, filter, nil)
	if err != nil {
		return err
	}

	updatedCount := 0

	for _, payable := range payables {
		needUpdate := false
		oldStatus := payable.Status

		// ตรวจสอบว่าเลยกำหนดชำระหรือไม่
		if !payable.DueDate.IsZero() && payable.DueDate.Before(now) {
			// เลยกำหนดและยังมียอดคงเหลือ
			if payable.Balance > 0 {
				payable.Status = "overdue"
				needUpdate = true
			}
		} else {
			// ยังไม่เลยกำหนด
			if payable.Balance > 0 && payable.Balance < payable.Amount {
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
				log.Printf("[CRON ERROR] อัปเดต Payable %s ไม่สำเร็จ: %v", payable.IDPayable, err)
				continue
			}
			updatedCount++
			log.Printf("[CRON] อัปเดต Payable %s: %s → %s", payable.InvoiceNo, oldStatus, payable.Status)
		}
	}

	log.Printf("[CRON] ตรวจสอบ Payable เสร็จสิ้น: อัปเดต %d รายการ จากทั้งหมด %d รายการ", updatedCount, len(payables))

	return nil
}

// checkReceivableStatus ตรวจสอบและอัปเดตสถานะของ Receivable
func (sc *StatusChecker) checkReceivableStatus(ctx context.Context) error {
	now := time.Now()

	// ดึงข้อมูล Receivable ทั้งหมดที่ยังไม่ถูกลบ และยังไม่ได้รับชำระครบ
	filter := bson.M{
		"deleted_at": nil,
		"status": bson.M{
			"$in": []string{"pending", "partial"},
		},
		"balance": bson.M{
			"$gt": 0,
		},
	}

	receivables, err := sc.receivableRepo.GetAllReceivablesByFilter(ctx, filter, nil)
	if err != nil {
		return err
	}

	updatedCount := 0

	for _, receivable := range receivables {
		needUpdate := false
		oldStatus := receivable.Status

		// ตรวจสอบว่าเลยกำหนดรับชำระหรือไม่
		if !receivable.DueDate.IsZero() && receivable.DueDate.Before(now) {
			// เลยกำหนดและยังมียอดคงเหลือ
			if receivable.Balance > 0 {
				receivable.Status = "overdue"
				needUpdate = true
			}
		} else {
			// ยังไม่เลยกำหนด
			if receivable.Balance > 0 && receivable.Balance < receivable.Amount {
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
				log.Printf("[CRON ERROR] อัปเดต Receivable %s ไม่สำเร็จ: %v", receivable.IDReceivable, err)
				continue
			}
			updatedCount++
			log.Printf("[CRON] อัปเดต Receivable %s: %s → %s", receivable.InvoiceNo, oldStatus, receivable.Status)
		}
	}

	log.Printf("[CRON] ตรวจสอบ Receivable เสร็จสิ้น: อัปเดต %d รายการ จากทั้งหมด %d รายการ", updatedCount, len(receivables))

	return nil
}

// RunNow รันการตรวจสอบทันที (สำหรับ testing หรือรันด้วยตนเอง)
func (sc *StatusChecker) RunNow() error {
	log.Println("[MANUAL] เริ่มตรวจสอบสถานะ Payable และ Receivable...")

	ctx := context.Background()

	// ตรวจสอบ Payables
	if err := sc.checkPayableStatus(ctx); err != nil {
		log.Printf("[MANUAL ERROR] ตรวจสอบ Payable ไม่สำเร็จ: %v", err)
		return err
	}

	// ตรวจสอบ Receivables
	if err := sc.checkReceivableStatus(ctx); err != nil {
		log.Printf("[MANUAL ERROR] ตรวจสอบ Receivable ไม่สำเร็จ: %v", err)
		return err
	}

	log.Println("[MANUAL] ตรวจสอบสถานะเสร็จสิ้น")

	return nil
}
