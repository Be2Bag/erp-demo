# Cron Jobs - Status Checker

## คำอธิบาย

Cronjob สำหรับตรวจสอบและอัปเดตสถานะของ **Payable** (เจ้าหนี้) และ **Receivable** (ลูกหนี้) อัตโนมัติ

## สถานะที่ระบบจัดการ

ระบบจะตรวจสอบและอัปเดตสถานะดังนี้:

### 1. **pending** (รอดำเนินการ)
- ยังไม่มีการชำระเงินเลย
- Balance = Amount (ยอดคงเหลือเท่ากับยอดทั้งหมด)
- ยังไม่เกินกำหนดชำระ

### 2. **partial** (ชำระบางส่วน)
- มีการชำระเงินบางส่วนแล้ว
- 0 < Balance < Amount
- ยังไม่เกินกำหนดชำระ

### 3. **paid** (ชำระครบแล้ว)
- ชำระเงินครบถ้วนแล้ว
- Balance = 0
- สถานะนี้จะถูกตั้งค่าเมื่อมีการบันทึกการชำระเงินในระบบ

### 4. **overdue** (เกินกำหนด)
- เลยกำหนดชำระ (DueDate < วันปัจจุบัน)
- ยังมียอดคงเหลือ (Balance > 0)
- ระบบจะเปลี่ยนสถานะอัตโนมัติเมื่อเลยกำหนด

## กำหนดการทำงาน

Cronjob จะทำงานอัตโนมัติ **ทุกวันเวลา 00:00 น. (เที่ยงคืน)**

รูปแบบ Cron Expression: `0 0 0 * * *`
- Second: 0
- Minute: 0
- Hour: 0 (เที่ยงคืน)
- Day: * (ทุกวัน)
- Month: * (ทุกเดือน)
- Weekday: * (ทุกวัน)

## โลจิกการตรวจสอบ

### สำหรับ Payable และ Receivable:

1. **ดึงข้อมูล**: ดึงรายการที่มีสถานะ `pending` หรือ `partial` และยังมียอดคงเหลือ (Balance > 0)

2. **ตรวจสอบวันครบกำหนด**:
   - ถ้า DueDate < วันปัจจุบัน และ Balance > 0 → เปลี่ยนเป็น **overdue**
   - ถ้ายังไม่เลยกำหนด:
     - Balance = Amount → สถานะ **pending**
     - 0 < Balance < Amount → สถานะ **partial**

3. **บันทึกการเปลี่ยนแปลง**: อัปเดตสถานะและเวลาที่แก้ไขล่าสุด (UpdatedAt)

## การใช้งาน

### 1. เริ่มต้น Cronjob อัตโนมัติ (ใน main.go)

```go
// เริ่มต้น Cronjob
statusChecker := cron.NewStatusChecker(payableRepo, receivableRepo)
if err := statusChecker.Start(); err != nil {
    log.Printf("เริ่ม cronjob ไม่สำเร็จ: %v", err)
} else {
    log.Println("Cronjob เริ่มทำงานแล้ว")
}
```

### 2. รันด้วยตนเอง (Manual Trigger)

หากต้องการรันการตรวจสอบทันที (ไม่รอถึงเวลา 00:00 น.):

```go
statusChecker := cron.NewStatusChecker(payableRepo, receivableRepo)
if err := statusChecker.RunNow(); err != nil {
    log.Printf("รัน cronjob ไม่สำเร็จ: %v", err)
}
```

### 3. หยุดการทำงาน

```go
statusChecker.Stop()
```

## Log Messages

Cronjob จะแสดง log ดังนี้:

### เมื่อเริ่มต้น:
```
[CRON] Status Checker เริ่มทำงานแล้ว (รันทุกวันเวลา 00:00 น.)
```

### เมื่อรันตรวจสอบ:
```
[CRON] เริ่มตรวจสอบสถานะ Payable และ Receivable...
[CRON] อัปเดต Payable INV-2024-001: partial → overdue
[CRON] ตรวจสอบ Payable เสร็จสิ้น: อัปเดต 3 รายการ จากทั้งหมด 10 รายการ
[CRON] อัปเดต Receivable REC-2024-005: pending → overdue
[CRON] ตรวจสอบ Receivable เสร็จสิ้น: อัปเดต 2 รายการ จากทั้งหมด 8 รายการ
[CRON] ตรวจสอบสถานะเสร็จสิ้น
```

### เมื่อเกิด Error:
```
[CRON ERROR] ตรวจสอบ Payable ไม่สำเร็จ: <error message>
[CRON ERROR] อัปเดต Receivable <id> ไม่สำเร็จ: <error message>
```

## ตัวอย่างการทำงาน

### Scenario 1: ใบแจ้งหนี้เลยกำหนด
- **วันที่สร้าง**: 2024-01-01
- **วันครบกำหนด**: 2024-01-15
- **จำนวนเงิน**: 10,000 บาท
- **ยอดคงเหลือ**: 10,000 บาท
- **สถานะเริ่มต้น**: pending

เมื่อถึงวันที่ 2024-01-16 (00:00 น.) → Cronjob จะเปลี่ยนสถานะเป็น **overdue**

### Scenario 2: ชำระบางส่วนและเลยกำหนด
- **จำนวนเงิน**: 10,000 บาท
- **ยอดคงเหลือ**: 5,000 บาท (ชำระไป 5,000)
- **สถานะปัจจุบัน**: partial
- **วันครบกำหนด**: 2024-01-15

เมื่อถึงวันที่ 2024-01-16 (00:00 น.) → Cronjob จะเปลี่ยนสถานะเป็น **overdue**

### Scenario 3: ชำระครบถ้วน
- **จำนวนเงิน**: 10,000 บาท
- **ยอดคงเหลือ**: 0 บาท (ชำระครบ)
- **สถานะ**: paid

→ Cronjob จะไม่ดึงข้อมูลนี้มาตรวจสอบ (เพราะ Balance = 0)

## Dependencies

- `github.com/robfig/cron/v3` - สำหรับจัดการ cron schedule

## โครงสร้างไฟล์

```
cron/
  └── status_checker.go    # Logic สำหรับตรวจสอบและอัปเดตสถานะ
```

## หมายเหตุ

1. **Performance**: ระบบจะดึงเฉพาะรายการที่จำเป็นต้องตรวจสอบ (สถานะ pending/partial และมียอดคงเหลือ)
2. **Error Handling**: หากอัปเดตรายการใดไม่สำเร็จ จะ skip และดำเนินการกับรายการถัดไป
3. **Timezone**: ใช้ timezone ของ server ในการเปรียบเทียบวันที่
4. **Manual Trigger**: สามารถสร้าง API endpoint เพื่อให้สามารถเรียกใช้ `RunNow()` ได้จากภายนอก

## การปรับแต่งเวลารัน

หากต้องการเปลี่ยนเวลาในการรัน cronjob สามารถแก้ไขได้ที่ไฟล์ `cron/status_checker.go`:

```go
// รันทุกวันเวลา 01:00 น.
_, err := sc.cron.AddFunc("0 0 1 * * *", func() { ... })

// รันทุก 6 ชั่วโมง
_, err := sc.cron.AddFunc("0 0 */6 * * *", func() { ... })

// รันทุกชั่วโมง
_, err := sc.cron.AddFunc("0 0 * * * *", func() { ... })
```

## การทดสอบ

สามารถทดสอบการทำงานได้โดย:

1. สร้างข้อมูล Payable/Receivable ที่มี DueDate ในอดีต
2. รัน `statusChecker.RunNow()` เพื่อทดสอบ
3. ตรวจสอบ log และสถานะในฐานข้อมูล
