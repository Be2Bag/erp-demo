package services

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/pkg/storage"
	"github.com/Be2Bag/erp-demo/pkg/util"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userService struct {
	config            config.Config
	userRepo          ports.UserRepository
	dropDownRepo      ports.DropDownRepository
	storageService    *storage.SupabaseStorage
	storageCloudflare *storage.CloudflareStorage
	taskRepo          ports.TaskRepository
}

func NewUserService(cfg config.Config, ur ports.UserRepository, dr ports.DropDownRepository, ss *storage.SupabaseStorage, sc *storage.CloudflareStorage, tr ports.TaskRepository) ports.UserService {
	return &userService{config: cfg, userRepo: ur, dropDownRepo: dr, storageService: ss, storageCloudflare: sc, taskRepo: tr}
}

func (s *userService) Create(ctx context.Context, req dto.RequestCreateUser) error {

	hashPassword := util.HashPassword(req.Password, s.config.Hash.Salt)

	user := &models.User{
		UserID:            uuid.New().String(),
		Email:             req.Email,
		Password:          hashPassword,
		TitleTH:           req.TitleTH,
		TitleEN:           req.TitleEN,
		FirstNameTH:       req.FirstNameTH,
		LastNameTH:        req.LastNameTH,
		FirstNameEN:       req.FirstNameEN,
		LastNameEN:        req.LastNameEN,
		IDCard:            req.IDCard,
		Role:              "user",
		Avatar:            req.Avatar,
		Phone:             req.Phone,
		Status:            "pending",
		EmployeeCode:      req.EmployeeCode,
		Gender:            req.Gender,
		BirthDate:         req.BirthDate,
		PositionID:        req.PositionID,
		DepartmentID:      req.DepartmentID,
		HireDate:          req.HireDate,
		EmploymentType:    req.EmploymentType,
		EmploymentHistory: []models.EmploymentHistory{},
		Address: models.Address{
			AddressLine1: req.Address.AddressLine1,
			AddressLine2: req.Address.AddressLine2,
			Subdistrict:  req.Address.Subdistrict,
			District:     req.Address.District,
			Province:     req.Address.Province,
			PostalCode:   req.Address.PostalCode,
			Country:      req.Address.Country,
		},
		BankInfo: models.BankInfo{
			BankName:    req.BankInfo.BankName,
			AccountNo:   req.BankInfo.AccountNo,
			AccountName: req.BankInfo.AccountName,
		},
		Documents: []models.Document{
			{Name: "", FileURL: "", Type: "idcards", CreatedAt: time.Now(), UploadedAt: time.Now(), DeletedAt: nil},
			{Name: "", FileURL: "", Type: "graduation", CreatedAt: time.Now(), UploadedAt: time.Now(), DeletedAt: nil},
			{Name: "", FileURL: "", Type: "transcript", CreatedAt: time.Now(), UploadedAt: time.Now(), DeletedAt: nil},
			{Name: "", FileURL: "", Type: "resume", CreatedAt: time.Now(), UploadedAt: time.Now(), DeletedAt: nil},
			{Name: "", FileURL: "", Type: "health", CreatedAt: time.Now(), UploadedAt: time.Now(), DeletedAt: nil},
			{Name: "", FileURL: "", Type: "military", CreatedAt: time.Now(), UploadedAt: time.Now(), DeletedAt: nil},
			{Name: "", FileURL: "", Type: "criminal", CreatedAt: time.Now(), UploadedAt: time.Now(), DeletedAt: nil},
			{Name: "", FileURL: "", Type: "other", CreatedAt: time.Now(), UploadedAt: time.Now(), DeletedAt: nil},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	filter := bson.M{"id_card": user.IDCard, "deleted_at": nil}
	projection := bson.M{"_id": 0, "id_card": 1}

	checkUser, errOnGetByIDCard := s.userRepo.GetUserByFilter(ctx, filter, projection)
	if errOnGetByIDCard != nil {
		return errOnGetByIDCard
	}

	if checkUser != nil {
		return fmt.Errorf("user with ID card %s already exists", user.IDCard)
	}

	_, errOnCreateUser := s.userRepo.Create(ctx, user)
	if errOnCreateUser != nil {
		return fmt.Errorf("failed to create user: %w", errOnCreateUser)
	}

	return nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*dto.ResponseGetUserByID, error) {

	filter := bson.M{"user_id": id, "deleted_at": nil}
	projection := bson.M{}

	users, errOnGetUser := s.userRepo.GetUserByFilter(ctx, filter, projection)

	if errOnGetUser != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", errOnGetUser)
	}

	if len(users) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	user := users[0]
	positionsName := "ไม่พบตำแหน่ง"
	departmentsName := "ไม่พบแผนก"
	provincesName := "ไม่พบจังหวัด"
	districtsName := "ไม่พบอำเภอ"
	subDistrictsName := "ไม่พบตำบล"
	positions, errOnGetPositions := s.dropDownRepo.GetPositions(ctx, bson.M{"position_id": user.PositionID}, bson.M{"_id": 0, "position_name": 1, "level": 1})
	if errOnGetPositions != nil {
		return nil, fmt.Errorf("failed to get position: %w", errOnGetPositions)
	}

	if len(positions) > 0 {
		positionsName = positions[0].PositionName + " (" + positions[0].Level + ")"
	}

	departments, errOnGetDepartments := s.dropDownRepo.GetDepartments(ctx, bson.M{"department_id": user.DepartmentID}, bson.M{"_id": 0, "department_name": 1})
	if errOnGetDepartments != nil {
		return nil, fmt.Errorf("failed to get department: %w", errOnGetDepartments)
	}
	if len(departments) > 0 {
		departmentsName = departments[0].DepartmentName
	}

	provinces, errOnGetProvinces := s.dropDownRepo.GetProvinces(ctx, bson.M{"id": user.Address.Province}, bson.M{"_id": 0, "name_th": 1})
	if errOnGetProvinces != nil {
		return nil, fmt.Errorf("failed to get province: %w", errOnGetProvinces)
	}
	if len(provinces) == 0 {
		return nil, fmt.Errorf("province not found")
	}
	provincesName = provinces[0].NameTH

	districts, errOnGetDistricts := s.dropDownRepo.GetDistricts(ctx, bson.M{"id": user.Address.District}, bson.M{"_id": 0, "name_th": 1})
	if errOnGetDistricts != nil {
		return nil, fmt.Errorf("failed to get district: %w", errOnGetDistricts)
	}
	if len(districts) == 0 {
		return nil, fmt.Errorf("district not found")
	}
	districtsName = districts[0].NameTH

	subDistricts, errOnGetSubDistricts := s.dropDownRepo.GetSubDistricts(ctx, bson.M{"id": user.Address.Subdistrict}, bson.M{"_id": 0, "name_th": 1})
	if errOnGetSubDistricts != nil {
		return nil, fmt.Errorf("failed to get subdistrict: %w", errOnGetSubDistricts)
	}
	if len(subDistricts) == 0 {
		return nil, fmt.Errorf("subdistrict not found")
	}
	subDistrictsName = subDistricts[0].NameTH

	var dtoDocuments []dto.Document
	for _, doc := range user.Documents {
		if doc.FileURL == "" {
			continue
		}
		dtoDocuments = append(dtoDocuments, dto.Document{
			Name:       doc.Name,
			FileURL:    doc.FileURL,
			Type:       doc.Type,
			CreatedAt:  doc.CreatedAt,
			UploadedAt: doc.UploadedAt,
			DeletedAt:  doc.DeletedAt,
		})
	}

	var dtoEmploymentHistory []dto.EmploymentHistory
	for _, eh := range user.EmploymentHistory {
		dtoEmploymentHistory = append(dtoEmploymentHistory, dto.EmploymentHistory{
			UserID:         eh.UserID,
			PositionID:     eh.PositionID,
			DepartmentID:   eh.DepartmentID,
			FromDate:       eh.FromDate,
			ToDate:         eh.ToDate,
			EmploymentType: eh.EmploymentType,
			Note:           eh.Note,
			CreatedAt:      eh.CreatedAt,
			UpdatedAt:      eh.UpdatedAt,
			DeletedAt:      eh.DeletedAt,
		})
	}

	dtoUser := &dto.ResponseGetUserByID{
		UserID:            user.UserID,
		Email:             user.Email,
		TitleTH:           user.TitleTH,
		TitleEN:           user.TitleEN,
		FirstNameTH:       user.FirstNameTH,
		LastNameTH:        user.LastNameTH,
		FirstNameEN:       user.FirstNameEN,
		LastNameEN:        user.LastNameEN,
		Phone:             user.Phone,
		Role:              user.Role,
		Avatar:            user.Avatar,
		IDCard:            user.IDCard,
		Status:            user.Status,
		EmployeeCode:      user.EmployeeCode,
		Gender:            user.Gender,
		BirthDate:         user.BirthDate,
		Position:          positionsName,
		Department:        departmentsName,
		HireDate:          user.HireDate,
		EmploymentType:    user.EmploymentType,
		EmploymentHistory: dtoEmploymentHistory,
		Address: dto.Address{
			AddressLine1: user.Address.AddressLine1,
			AddressLine2: user.Address.AddressLine2,
			Subdistrict:  subDistrictsName,
			District:     districtsName,
			Province:     provincesName,
			PostalCode:   user.Address.PostalCode,
			Country:      user.Address.Country,
		},
		BankInfo: dto.BankInfo{
			BankName:    user.BankInfo.BankName,
			AccountNo:   user.BankInfo.AccountNo,
			AccountName: user.BankInfo.AccountName,
		},
		Documents: dtoDocuments,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}

	return dtoUser, nil
}

func (s *userService) GetAll(ctx context.Context, req dto.RequestGetUserAll) (dto.Pagination, error) {

	filter := bson.M{
		"deleted_at": nil,
		"role":       bson.M{"$ne": "admin"},
	}

	if req.Status != "" {
		filter["status"] = req.Status
	}

	if req.Role != "" {
		filter["role"] = req.Role
	}

	if req.Search != "" {

		filter["$or"] = []bson.M{
			{"first_name_th": bson.M{"$regex": req.Search, "$options": "i"}},
			{"last_name_th": bson.M{"$regex": req.Search, "$options": "i"}},
			{"first_name_en": bson.M{"$regex": req.Search, "$options": "i"}},
			{"last_name_en": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
	}

	if req.SortBy != "" {
		order := 1
		if strings.ToLower(req.SortOrder) == "desc" {
			order = -1
		}
		pipeline = append(pipeline, bson.D{
			{Key: "$sort", Value: bson.D{{Key: req.SortBy, Value: order}}},
		})
	}

	skip := (req.Page - 1) * req.Limit
	pipeline = append(pipeline,
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: req.Limit}},
	)

	users, err := s.userRepo.AggregateUser(ctx, pipeline)
	if err != nil {
		return dto.Pagination{}, err
	}

	var dtoUsers []*dto.ResponseGetUserAll
	for _, u := range users {

		positionsName := "ไม่พบตำแหน่ง"
		departmentsName := "ไม่พบแผนก"
		tasksCompleted := "0"
		tasksTotal := "0"

		positions, _ := s.dropDownRepo.GetPositions(ctx, bson.M{"position_id": u.PositionID, "deleted_at": nil}, bson.M{"_id": 0, "position_name": 1, "level": 1})

		if len(positions) > 0 {
			positionsName = positions[0].PositionName + " (" + positions[0].Level + ")"
		}

		departments, _ := s.dropDownRepo.GetDepartments(ctx, bson.M{"department_id": u.DepartmentID, "deleted_at": nil}, bson.M{"_id": 0, "department_name": 1})

		if len(departments) > 0 {
			departmentsName = departments[0].DepartmentName
		}

		filterUserTaskStats := bson.M{"user_id": u.UserID, "deleted_at": nil}
		projectionUserTaskStats := bson.M{}

		existingStats, _ := s.taskRepo.GetOneUserTaskStatsByFilter(ctx, filterUserTaskStats, projectionUserTaskStats)

		if existingStats != nil {
			tasksCompleted = fmt.Sprintf("%d", existingStats.Totals.Completed)
			tasksTotal = fmt.Sprintf("%d", existingStats.Totals.Assigned)
		}

		dtoUsers = append(dtoUsers, &dto.ResponseGetUserAll{
			UserID:         u.UserID,
			TitleTH:        u.TitleTH,
			FirstNameTH:    u.FirstNameTH,
			LastNameTH:     u.LastNameTH,
			TitleEN:        u.TitleEN,
			FirstNameEN:    u.FirstNameEN,
			LastNameEN:     u.LastNameEN,
			Avatar:         u.Avatar,
			Email:          u.Email,
			Phone:          u.Phone,
			PositionID:     u.PositionID,
			PositionName:   positionsName,
			DepartmentID:   u.DepartmentID,
			DepartmentName: departmentsName,
			Role:           u.Role,
			Status:         u.Status,
			KPIScore:       "100%",
			TasksCompleted: tasksCompleted,
			TasksTotal:     tasksTotal,
			CreatedAt:      u.CreatedAt,
			UpdatedAt:      u.UpdatedAt,
			DeletedAt:      u.DeletedAt,
		})
	}

	totalCount, _ := s.userRepo.CountUsers(ctx, filter)

	list := make([]interface{}, len(dtoUsers))
	for i, v := range dtoUsers {
		list[i] = v
	}

	pagination := dto.Pagination{
		Page:       req.Page,
		Size:       len(dtoUsers),
		TotalCount: int(totalCount),
		TotalPages: int(math.Ceil(float64(totalCount) / float64(req.Limit))),
		List:       list,
	}

	return pagination, nil
}

func (s *userService) UpdateUserByID(ctx context.Context, id string, req dto.RequestUpdateUser) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, mongo.ErrNoDocuments
	}

	if req.EmploymentHistory != nil {
		var employmentHistory []models.EmploymentHistory
		for _, eh := range req.EmploymentHistory {
			employmentHistory = append(employmentHistory, models.EmploymentHistory{
				UserID:         eh.UserID,
				PositionID:     eh.PositionID,
				DepartmentID:   eh.DepartmentID,
				FromDate:       eh.FromDate,
				ToDate:         eh.ToDate,
				EmploymentType: eh.EmploymentType,
				Note:           eh.Note,
				CreatedAt:      eh.CreatedAt,
				UpdatedAt:      eh.UpdatedAt,
				DeletedAt:      eh.DeletedAt,
			})
		}
		user.EmploymentHistory = employmentHistory
	}

	if req.Documents != nil {
		var documents []models.Document
		for _, doc := range req.Documents {
			documents = append(documents, models.Document{
				Name:       doc.Name,
				FileURL:    doc.FileURL,
				Type:       doc.Type,
				CreatedAt:  doc.CreatedAt,
				UploadedAt: doc.UploadedAt,
				DeletedAt:  doc.DeletedAt,
			})
		}
		user.Documents = documents
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.TitleTH != "" {
		user.TitleTH = req.TitleTH
	}
	if req.TitleEN != "" {
		user.TitleEN = req.TitleEN
	}
	if req.FirstNameTH != "" {
		user.FirstNameTH = req.FirstNameTH
	}
	if req.LastNameTH != "" {
		user.LastNameTH = req.LastNameTH
	}
	if req.FirstNameEN != "" {
		user.FirstNameEN = req.FirstNameEN
	}
	if req.LastNameEN != "" {
		user.LastNameEN = req.LastNameEN
	}
	if req.IDCard != "" {
		user.IDCard = req.IDCard
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.EmployeeCode != "" {
		user.EmployeeCode = req.EmployeeCode
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for birth_date: %w", err)
		}
		user.BirthDate = parsedDate
	}
	if req.PositionID != "" {
		user.PositionID = req.PositionID
	}
	if req.DepartmentID != "" {
		user.DepartmentID = req.DepartmentID
	}
	if req.HireDate != "" {
		parsedHireDate, err := time.Parse("2006-01-02", req.HireDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for hire_date: %w", err)
		}
		user.HireDate = parsedHireDate
	}
	if req.EmploymentType != "" {
		user.EmploymentType = req.EmploymentType
	}
	if req.Address != (dto.Address{}) {
		user.Address = models.Address{
			AddressLine1: req.Address.AddressLine1,
			AddressLine2: req.Address.AddressLine2,
			Subdistrict:  req.Address.Subdistrict,
			District:     req.Address.District,
			Province:     req.Address.Province,
			PostalCode:   req.Address.PostalCode,
			Country:      req.Address.Country,
		}
	}
	if req.BankInfo != (dto.BankInfo{}) {
		user.BankInfo = models.BankInfo{
			BankName:    req.BankInfo.BankName,
			AccountNo:   req.BankInfo.AccountNo,
			AccountName: req.BankInfo.AccountName,
		}
	}
	user.UpdatedAt = time.Now()

	if user.Status == "rejected" {
		user.Status = "pending"
	}

	updateUser, errOnUpdateUserByID := s.userRepo.UpdateUserByID(ctx, id, user)
	if errOnUpdateUserByID != nil {
		if errOnUpdateUserByID == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, errOnUpdateUserByID
	}

	return updateUser, nil

}

func (s *userService) DeleteUserByID(ctx context.Context, id string) error {

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user == nil {
		return mongo.ErrNoDocuments
	}

	now := time.Now()
	user.DeletedAt = &now
	user.UpdatedAt = time.Now()

	result, err := s.userRepo.UpdateUserByID(ctx, id, user)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result == nil {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *userService) UpdateDocuments(ctx context.Context, req dto.RequestUpdateDocuments) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if user == nil {
		return nil, mongo.ErrNoDocuments
	}

	if req.Type == "avatars" {

		avatarURL := user.Avatar
		parts := strings.Split(avatarURL, "/")
		filename := ""
		if len(parts) > 0 {
			filename = parts[len(parts)-1]
		}
		errOnDelete := s.storageCloudflare.DeleteFile(req.Type, filename)
		if errOnDelete != nil {
			return nil, fmt.Errorf("failed to delete file: %w", errOnDelete)
		}

		user.Avatar = req.FileURL

	} else {

		var documents []models.Document
		for _, doc := range user.Documents {
			if doc.Type == req.Type {

				documentsURL := doc.FileURL
				parts := strings.Split(documentsURL, "/")
				filename := ""
				if len(parts) > 0 {
					filename = parts[len(parts)-1]
				}

				errOnDelete := s.storageCloudflare.DeleteFile(req.Type, filename)
				if errOnDelete != nil {
					return nil, fmt.Errorf("failed to delete file: %w", errOnDelete)
				}

				doc.Name = req.Name
				doc.FileURL = req.FileURL
				doc.UploadedAt = time.Now()
			}
			documents = append(documents, doc)
		}

		user.Documents = documents
		user.UpdatedAt = time.Now()
	}

	if user.Status == "rejected" {
		user.Status = "pending"
	}

	updatedUser, errOnUpdate := s.userRepo.UpdateUserByID(ctx, req.UserID, user)
	if errOnUpdate != nil {
		return nil, fmt.Errorf("failed to update user documents: %w", errOnUpdate)
	}

	return updatedUser, nil
}
