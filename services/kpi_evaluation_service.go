package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type kpiEvaluationRepoService struct {
	config            config.Config
	kpiRepo           ports.KPIRepository
	userRepo          ports.UserRepository
	kpiEvaluationRepo ports.KPIEvaluationRepository
	taskRepo          ports.TaskRepository
	departmentRepo    ports.DepartmentRepository
	projectRepo       ports.ProjectRepository
	signJobRepo       ports.SignJobRepository
}

func NewKPIEvaluationService(cfg config.Config, kpiRepo ports.KPIRepository, userRepo ports.UserRepository, kpiEvaluationRepo ports.KPIEvaluationRepository, taskRepo ports.TaskRepository, departmentRepo ports.DepartmentRepository, projectRepo ports.ProjectRepository, signJobRepo ports.SignJobRepository) ports.KPIEvaluationService {
	return &kpiEvaluationRepoService{config: cfg, kpiRepo: kpiRepo, userRepo: userRepo, kpiEvaluationRepo: kpiEvaluationRepo, taskRepo: taskRepo, departmentRepo: departmentRepo, projectRepo: projectRepo, signJobRepo: signJobRepo}
}

// func (s *kpiService) GetKPITemplateByID(ctx context.Context, kpiID string) (*dto.KPITemplateDTO, error) {

// 	filter := bson.M{"kpi_id": kpiID, "deleted_at": nil}
// 	projection := bson.M{}

// 	m, err := s.kpiRepo.GetOneKPIByFilter(ctx, filter, projection)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if m == nil {
// 		return nil, nil
// 	}

// 	ItemsDTO := make([]dto.KPITemplateItemDTO, 0, len(m.Items))
// 	for _, st := range m.Items {
// 		ItemsDTO = append(ItemsDTO, dto.KPITemplateItemDTO{
// 			ItemID:      st.ItemID,
// 			Name:        st.Name,
// 			Description: st.Description,
// 			Category:    st.Category,
// 			MaxScore:    st.MaxScore,
// 			Weight:      st.Weight,
// 			CreatedAt:   st.CreatedAt,
// 			UpdatedAt:   st.UpdatedAt,
// 		})
// 	}

// 	dtoObj := &dto.KPITemplateDTO{
// 		KPIID:       m.KPIID,
// 		KPIName:     m.KPIName,
// 		Department:  m.Department,
// 		TotalWeight: m.TotalWeight,
// 		Items:       ItemsDTO,
// 		Version:     m.Version,
// 		CreatedBy:   m.CreatedBy,
// 		CreatedAt:   m.CreatedAt,
// 		UpdatedAt:   m.UpdatedAt,
// 	}
// 	return dtoObj, nil
// }

func (s *kpiEvaluationRepoService) UpdateKPIEvaluation(ctx context.Context, evaluationID string, req dto.UpdateKPIEvaluationRequest, claims *dto.JWTClaims) error {

	now := time.Now()
	// ดึงข้อมูลเดิม
	filter := bson.M{"evaluation_id": evaluationID, "deleted_at": nil}
	existing, err := s.kpiEvaluationRepo.GetOneKPIEvaluationByFilter(ctx, filter, bson.M{})
	if err != nil {
		return err
	}

	if existing == nil {
		return mongo.ErrNoDocuments
	}

	if req.Scores != nil {

		reqMap := make(map[string]dto.KPIScoreRequest, len(req.Scores))
		for _, r := range req.Scores {
			id := strings.TrimSpace(r.ItemID)
			if id == "" {
				continue
			}
			reqMap[id] = r
		}

		updatedAny := false
		total := 0
		for i, sc := range existing.Scores {
			if r, ok := reqMap[sc.ItemID]; ok {
				score := r.Score
				if score < 0 {
					score = 0
				}
				if score > sc.MaxScore {
					score = sc.MaxScore
				}
				existing.Scores[i].Score = score
				existing.Scores[i].Notes = strings.TrimSpace(r.Notes)
				updatedAny = true
			}
			total += existing.Scores[i].Score
		}

		if updatedAny {
			existing.TotalScore = total
			existing.IsEvaluated = true
		}
	}

	existing.Version += 1    // เพิ่มเวอร์ชัน
	existing.UpdatedAt = now // อัปเดตเวลา

	updated, err := s.kpiEvaluationRepo.UpdateKPIEvaluationByID(ctx, evaluationID, *existing)
	if err != nil {
		return err
	}
	if updated == nil {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *kpiEvaluationRepoService) ListKPIEvaluation(ctx context.Context, claims *dto.JWTClaims, page, size int, search string, department string, sortBy string, sortOrder string) (dto.Pagination, error) {
	skip := int64((page - 1) * size)
	limit := int64(size)

	filter := bson.M{
		"deleted_at": nil,
	}

	department = strings.TrimSpace(department)
	if department != "" {
		filter["department_id"] = department
	}

	search = strings.TrimSpace(search)
	if search != "" {
		safe := regexp.QuoteMeta(search)
		re := primitive.Regex{Pattern: safe, Options: "i"}
		filter["$or"] = []bson.M{
			{"project_name": re},
			{"job_name": re},
		}
	}

	projection := bson.M{}

	allowedSortFields := map[string]string{
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"project_name": "project_name",
		"job_name":     "job_name",
	}
	field, ok := allowedSortFields[sortBy]
	if !ok || field == "" {
		field = "created_at"
	}
	order := int32(-1)
	if strings.EqualFold(sortOrder, "asc") {
		order = 1
	}

	sort := bson.D{
		{Key: field, Value: order},
		{Key: "_id", Value: -1},
	}

	items, total, err := s.kpiEvaluationRepo.GetListKPIEvaluationByFilter(ctx, filter, projection, sort, skip, limit)
	if err != nil {
		return dto.Pagination{}, fmt.Errorf("list tasks: %w", err)
	}

	list := make([]interface{}, 0, len(items))
	for _, m := range items {

		departmentsName := "ไม่พบแผนก"
		assigneeName := "ไม่พบผู้รับผิดชอบ"
		jobName := "ไม่พบใบงาน"
		projectName := "ไม่พบโครงการ"

		if dept, _ := s.departmentRepo.GetOneDepartmentByFilter(ctx, bson.M{"department_id": m.Department, "deleted_at": nil}, bson.M{"_id": 0, "department_name": 1}); dept != nil {
			departmentsName = dept.DepartmentName
		}
		if assignee, _ := s.userRepo.GetByID(ctx, m.EvaluateeID); assignee != nil {
			assigneeName = fmt.Sprintf("%s %s %s", assignee.TitleTH, assignee.FirstNameTH, assignee.LastNameTH)
		}

		if job, _ := s.signJobRepo.GetOneSignJobByFilter(ctx, bson.M{"job_id": m.JobID, "deleted_at": nil}, bson.M{"_id": 0, "job_name": 1, "project_id": 1}); job != nil {
			jobName = job.JobName
			if project, _ := s.projectRepo.GetOneProjectByFilter(ctx, bson.M{"project_id": job.ProjectID, "deleted_at": nil}, bson.M{"_id": 0, "project_name": 1}); project != nil {
				projectName = project.ProjectName
			}
		}

		scores := make([]dto.KPIScoreResponse, 0, len(m.Scores))
		for _, score := range m.Scores {
			scores = append(scores, dto.KPIScoreResponse{
				ItemID:   score.ItemID,
				Name:     score.Name,
				Category: score.Category,
				Weight:   score.Weight,
				MaxScore: score.MaxScore,
				Score:    score.Score,
				Notes:    score.Notes,
			})
		}

		list = append(list, dto.KPIEvaluationResponse{
			EvaluationID:   m.EvaluationID,
			JobID:          m.JobID,
			JobName:        jobName,
			ProjectID:      m.ProjectID,
			ProjectName:    projectName,
			TaskID:         m.TaskID,
			KPIID:          m.KPIID,
			KPIName:        "",
			Version:        1,
			EvaluatorID:    "",
			EvaluatorName:  "ไม่มีผู้ประเมิน",
			EvaluateeID:    m.EvaluateeID,
			EvaluateeName:  assigneeName,
			Department:     m.Department,
			DepartmentName: departmentsName,
			KPIScore:       "",
			Scores:         scores,
			TotalScore:     m.TotalScore,
			IsEvaluated:    m.IsEvaluated,
			Feedback:       m.Feedback,
			FinishedAt:     m.UpdatedAt,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
		})
	}

	totalPages := 0
	if total > 0 && size > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
	}

	return dto.Pagination{
		Page:       page,
		Size:       size,
		TotalCount: int(total),
		TotalPages: totalPages,
		List:       list,
	}, nil
}
