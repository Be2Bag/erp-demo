package repositories

import (
	"context"

	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/ports"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dropDownRepo struct {
	departmentsColl   *mongo.Collection
	positionsColl     *mongo.Collection
	provincesColl     *mongo.Collection
	districtsColl     *mongo.Collection
	subDistrictsColl  *mongo.Collection
	signTypesColl     *mongo.Collection
	customerTypesColl *mongo.Collection
	signJobsColl      *mongo.Collection
	projectsColl      *mongo.Collection
	usersColl         *mongo.Collection
}

func NewDropDownRepository(db *mongo.Database) ports.DropDownRepository {
	return &dropDownRepo{
		departmentsColl:   db.Collection(models.CollectionDepartments),
		positionsColl:     db.Collection(models.CollectionPositions),
		provincesColl:     db.Collection(models.CollectionProvinces),
		districtsColl:     db.Collection(models.CollectionDistricts),
		subDistrictsColl:  db.Collection(models.CollectionSubDistricts),
		signTypesColl:     db.Collection(models.CollectionSignTypes),
		customerTypesColl: db.Collection(models.CollectionCustomerTypes),
		signJobsColl:      db.Collection(models.CollectionSignJobs),
		projectsColl:      db.Collection(models.CollectionSProject),
		usersColl:         db.Collection(models.CollectionUsers),
	}
}
func (r *dropDownRepo) GetPositions(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Position, error) {
	var positions []*models.Position
	cursor, err := r.positionsColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var position models.Position
		if err := cursor.Decode(&position); err != nil {
			return nil, err
		}
		positions = append(positions, &position)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return positions, nil
}

func (r *dropDownRepo) GetDepartments(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Department, error) {
	var departments []*models.Department
	cursor, err := r.departmentsColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var department models.Department
		if err := cursor.Decode(&department); err != nil {
			return nil, err
		}
		departments = append(departments, &department)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return departments, nil
}

func (r *dropDownRepo) GetProvinces(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Province, error) {
	var provinces []*models.Province
	cursor, err := r.provincesColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var province models.Province
		if err := cursor.Decode(&province); err != nil {
			return nil, err
		}
		provinces = append(provinces, &province)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return provinces, nil
}
func (r *dropDownRepo) GetDistricts(ctx context.Context, filter interface{}, projection interface{}) ([]*models.District, error) {
	var districts []*models.District
	cursor, err := r.districtsColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var district models.District
		if err := cursor.Decode(&district); err != nil {
			return nil, err
		}
		districts = append(districts, &district)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return districts, nil
}
func (r *dropDownRepo) GetSubDistricts(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SubDistrict, error) {
	var subDistricts []*models.SubDistrict
	cursor, err := r.subDistrictsColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var subDistrict models.SubDistrict
		if err := cursor.Decode(&subDistrict); err != nil {
			return nil, err
		}
		subDistricts = append(subDistricts, &subDistrict)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return subDistricts, nil
}

func (r *dropDownRepo) GetSignTypes(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SignType, error) {
	var signTypes []*models.SignType
	cursor, err := r.signTypesColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var signType models.SignType
		if err := cursor.Decode(&signType); err != nil {
			return nil, err
		}
		signTypes = append(signTypes, &signType)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return signTypes, nil
}

func (r *dropDownRepo) GetCustomerTypes(ctx context.Context, filter interface{}, projection interface{}) ([]*models.CustomerType, error) {
	var customerTypes []*models.CustomerType
	cursor, err := r.customerTypesColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var customerType models.CustomerType
		if err := cursor.Decode(&customerType); err != nil {
			return nil, err
		}
		customerTypes = append(customerTypes, &customerType)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return customerTypes, nil
}

func (r *dropDownRepo) GetSignJobsList(ctx context.Context, filter interface{}, projection interface{}) ([]*models.SignJob, error) {
	var signJobs []*models.SignJob
	cursor, err := r.signJobsColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var signJob models.SignJob
		if err := cursor.Decode(&signJob); err != nil {
			return nil, err
		}
		signJobs = append(signJobs, &signJob)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return signJobs, nil
}

func (r *dropDownRepo) GetProjectsList(ctx context.Context, filter interface{}, projection interface{}) ([]*models.Project, error) {
	var projects []*models.Project

	cursor, err := r.projectsColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var project models.Project
		if err := cursor.Decode(&project); err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *dropDownRepo) GetUsersList(ctx context.Context, filter interface{}, projection interface{}) ([]*models.User, error) {
	var users []*models.User

	cursor, err := r.usersColl.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
