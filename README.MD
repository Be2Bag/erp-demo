```text
.
├── .gitignore
├── README.MD
├── app
│   └── app.go
├── cmd
│   └── main.go
├── config
│   └── config.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── dto
│   ├── admin_dto.go
│   ├── auth_dto.go
│   ├── dropdown_dto.go
│   ├── jwt_claims.go
│   ├── kpi_dto.go
│   ├── pagination_dto.go
│   ├── project_dto.go
│   ├── response_dto.go
│   ├── sign_job_dto.go
│   ├── swagger_wrapper.go
│   ├── task_dto.go
│   ├── upload_dto.go
│   ├── user_dto.go
│   └── workflow_dto.go
├── go.mod
├── go.sum
├── handlers
│   ├── admin_handler.go
│   ├── auth_handler.go
│   ├── dropdown_handler.go
│   ├── kpi_handler.go
│   ├── project_handler.go
│   ├── sign_job_handler.go
│   ├── task_handler.go
│   ├── upload_handler.go
│   ├── user_handler.go
│   └── workflow_handler.go
├── middleware
│   └── middleware.go
├── models
│   ├── customer_types.go
│   ├── department.go
│   ├── district.go
│   ├── employee_history.go
│   ├── kpi.go
│   ├── performance_record.go
│   ├── position.go
│   ├── project.go
│   ├── province.go
│   ├── sign_job.go
│   ├── sign_types.go
│   ├── sub_district.go
│   ├── task.go
│   ├── user.go
│   └── workflow.go
├── pkg
│   ├── db
│   │   └── mongo.go
│   ├── storage
│   │   ├── cloudflare.go
│   │   └── supabase.go
│   └── util
│       ├── cookie.go
│       ├── crypto.go
│       ├── idcard.go
│       └── mail.go
├── ports
│   ├── admin_port.go
│   ├── auth_port.go
│   ├── dropdown_port.go
│   ├── kpi_port.go
│   ├── project_port.go
│   ├── sign_job_port.go
│   ├── task_port.go
│   ├── upload_port.go
│   ├── user_port.go
│   └── workflow_port.go
├── repositories
│   ├── admin_repository.go
│   ├── auth_repository.go
│   ├── dropdown_repository.go
│   ├── kpi_repository.go
│   ├── project_repository.go
│   ├── sign_job_repository.go
│   ├── task_repository.go
│   ├── upload_repository.go
│   ├── user_repository.go
│   └── workflow_repository.go
└── services
  ├── admin_service.go
  ├── auth_service.go
  ├── dropdown_service.go
  ├── kpi_service.go
  ├── project_service.go
  ├── sign_job_service.go
  ├── task_service.go
  ├── upload_service.go
  ├── user_service.go
  └── workflow_service.go
```

## คำอธิบายโครงสร้าง (ภาษาไทย)

- .gitignore: ระบุไฟล์/โฟลเดอร์ที่ไม่ให้ Git ติดตาม
- README.MD: คำอธิบายโปรเจกต์และวิธีเริ่มต้นใช้งาน
- app/: บูตแอป ตั้งค่า Router, Middleware, DI
- cmd/main.go: จุดเริ่มรันโปรแกรม
- config/config.go: โหลด/จัดการคอนฟิก (ENV)
- docs/: เอกสาร API (Swagger/OpenAPI)
- dto/: โครงสร้างแลกเปลี่ยนข้อมูล (request/response) เพิ่ม project_dto, sign_job_dto, task_dto
- handlers/: เลเยอร์ HTTP (Controller) เพิ่ม project, sign_job, task
- middleware/: มิดเดิลแวร์ (Auth, Logging, CORS)
- models/: โดเมนโมเดล เพิ่ม project, task, sign_job, sign_types, customer_types
- pkg/:
  - db/: การเชื่อมต่อฐานข้อมูล
  - storage/: จัดเก็บไฟล์ภายนอก
  - util/: ฟังก์ชันช่วยเหลือทั่วไป
- ports/: อินเทอร์เฟซ (สัญญา) สำหรับ service/repository (เพิ่ม project, sign_job, task)
- repositories/: เข้าถึง/จัดการข้อมูลถาวร (เพิ่ม project, sign_job, task)
- services/: ธุรกิจลอจิก (เพิ่ม project, sign_job, task)

### การไหลของคำขอ (โดยย่อ)
Client → handlers → services → repositories → pkg/db
      ↳ middleware (ตรวจสอบ/แทรกกลางคำขอ)
      ↳ pkg/storage (อัปโหลด/ดึงไฟล์)

swag init -g cmd/main.go 