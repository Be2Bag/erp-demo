```text
.
+--- .gitignore
+--- README.MD
+--- app
|   \--- app.go
+--- cmd
|   \--- main.go
+--- config
|   \--- config.go
+--- docs
|   +--- docs.go
|   +--- swagger.json
|   \--- swagger.yaml
+--- dto
|   +--- admin_dto.go
|   +--- auth_dto.go
|   +--- dropdown_dto.go
|   +--- jwt_claims.go
|   +--- kpi_dto.go
|   +--- pagination_dto.go
|   +--- response_dto.go
|   +--- swagger_wrapper.go
|   +--- upload_dto.go
|   +--- user_dto.go
|   \--- workflow_dto.go
+--- go.mod
+--- go.sum
+--- handlers
|   +--- admin_handler.go
|   +--- auth_handler.go
|   +--- dropdown_handler.go
|   +--- kpi_handler.go
|   +--- task_handler.go
|   +--- upload_handler.go
|   +--- user_handler.go
|   \--- workflow_handler.go
+--- middleware
|   \--- middleware.go
+--- models
|   +--- department.go
|   +--- district.go
|   +--- employee_history.go
|   +--- kpi.go
|   +--- performance_record.go
|   +--- position.go
|   +--- province.go
|   +--- sub_district.go
|   +--- user.go
|   \--- workflow.go
+--- pkg
|   +--- db
|   |   \--- mongo.go
|   +--- storage
|   |   +--- cloudflare.go
|   |   \--- supabase.go
|   \--- util
|       +--- cookie.go
|       +--- crypto.go
|       +--- idcard.go
|       \--- mail.go
+--- ports
|   +--- admin_port.go
|   +--- auth_port.go
|   +--- dropdown_port.go
|   +--- kpi_port.go
|   +--- task_port.go
|   +--- upload_port.go
|   +--- user_port.go
|   \--- workflow_port.go
+--- repositories
|   +--- admin_repository.go
|   +--- auth_repository.go
|   +--- dropdown_repository.go
|   +--- kpi_repository.go
|   +--- task_repository.go
|   +--- upload_repository.go
|   +--- user_repository.go
|   \--- workflow_repository.go
\--- services
    +--- admin_service.go
    +--- auth_service.go
    +--- dropdown_service.go
    +--- kpi_service.go
    +--- task_service.go
    +--- upload_service.go
    +--- user_service.go
    \--- workflow_service.go
```

## คำอธิบายโครงสร้าง (ภาษาไทย)

- .gitignore: ระบุไฟล์/โฟลเดอร์ที่ไม่ให้ Git ติดตาม
- README.MD: คำอธิบายโปรเจกต์และวิธีเริ่มต้นใช้งาน
- app/: ประกอบและบูตแอป เช่น ตั้งค่า Router, Middleware, DI
  - app.go: สร้างอินสแตนซ์แอปและผูกเส้นทาง
- cmd/:
  - main.go: จุดเริ่มรันโปรแกรม (entry point)
- config/:
  - config.go: โหลด/จัดการค่าคอนฟิก (เช่น ENV)
- docs/: เอกสาร API (Swagger/OpenAPI)
  - docs.go, swagger.json, swagger.yaml: สคีมาและเอกสาร API
- dto/: โครงสร้างข้อมูลแลกเปลี่ยนระหว่างเลเยอร์ (request/response/payload)
- handlers/: เลเยอร์ HTTP (Controller) รับคำขอและเรียก Service
- middleware/: มิดเดิลแวร์ (เช่น Auth, Logging, CORS)
- models/: โดเมนโมเดลของระบบ (เช่น User, KPI, Workflow)
- pkg/:
  - db/mongo.go: เชื่อมต่อ/จัดการฐานข้อมูล MongoDB
  - storage/: การจัดเก็บไฟล์ภายนอก (Cloudflare, Supabase)
  - util/: เครื่องมือช่วย (cookie, crypto, idcard, mail)
- ports/: สัญญา/อินเทอร์เฟซระหว่างเลเยอร์ (Hexagonal/Clean Architecture)
- repositories/: เข้าถึงข้อมูล (CRUD) ติดต่อฐานข้อมูล/สตอเรจ
- services/: ธุรกิจลอจิก ประสานงาน repository และกติกาโดเมน

### การไหลของคำขอ (โดยย่อ)
Client → handlers → services → repositories → pkg/db
             ↳ middleware (ตรวจสอบ/แทรกกลางคำขอ)
             ↳ pkg/storage (อัปโหลด/ดึงไฟล์)