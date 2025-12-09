package adapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/dto"
	"github.com/I-Van-Radkov/corporate-messenger/directory-service/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DirectoryRepo struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewDirectoryRepo(db *pgxpool.Pool) *DirectoryRepo {
	return &DirectoryRepo{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *DirectoryRepo) CreateDepartment(ctx context.Context, dep *models.Department) (uuid.UUID, error) {
	query := r.builder.Insert("departments").
		Columns("department_id", "name", "parent_id", "created_at", "updated_at").
		Values(dep.DepartmentID.String(), dep.Name, nullString(dep.ParentID), dep.CreatedAt, dep.UpdatedAt).
		Suffix("RETURNING department_id")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var id string
	err = r.db.QueryRow(ctx, sqlQuery, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(id)
}

func (r *DirectoryRepo) GetDepartmentByID(ctx context.Context, id uuid.UUID) (*models.Department, error) {
	query := r.builder.Select("department_id, name, parent_id, created_at, updated_at").
		From("departments").
		Where(squirrel.Eq{"department_id": id.String()})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sqlQuery, args...)

	var dep models.Department
	var parentID sql.NullString
	err = row.Scan(&dep.DepartmentID, &dep.Name, &parentID, &dep.CreatedAt, &dep.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if parentID.Valid {
		pid, _ := uuid.Parse(parentID.String)
		dep.ParentID = &pid
	}

	return &dep, nil
}

func (r *DirectoryRepo) GetDepartments(ctx context.Context, limit, offset int) ([]*models.Department, int, error) {
	// Count
	countQuery := r.builder.Select("COUNT(*)").From("departments")
	sqlCount, _, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	err = r.db.QueryRow(ctx, sqlCount).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// List
	query := r.builder.Select("department_id, name, parent_id, created_at, updated_at").
		From("departments").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var deps []*models.Department
	for rows.Next() {
		var dep models.Department
		var parentID sql.NullString
		err = rows.Scan(&dep.DepartmentID, &dep.Name, &parentID, &dep.CreatedAt, &dep.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		if parentID.Valid {
			pid, _ := uuid.Parse(parentID.String)
			dep.ParentID = &pid
		}
		deps = append(deps, &dep)
	}

	return deps, total, nil
}

func (r *DirectoryRepo) GetAllDepartments(ctx context.Context) ([]*models.Department, error) {
	query := r.builder.Select("department_id, name, parent_id, created_at, updated_at").
		From("departments")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deps []*models.Department
	for rows.Next() {
		var dep models.Department
		var parentID sql.NullString
		err = rows.Scan(&dep.DepartmentID, &dep.Name, &parentID, &dep.CreatedAt, &dep.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if parentID.Valid {
			pid, _ := uuid.Parse(parentID.String)
			dep.ParentID = &pid
		}
		deps = append(deps, &dep)
	}

	return deps, nil
}

func (r *DirectoryRepo) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	query := r.builder.Delete("departments").
		Where(squirrel.Eq{"department_id": id.String()})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sqlQuery, args...)
	return err
}

func (r *DirectoryRepo) GetDepartmentMembers(ctx context.Context, depID uuid.UUID, limit, offset int) ([]*models.User, int, error) {
	countQuery := r.builder.Select("COUNT(*)").From("users").
		Where(squirrel.Eq{"department_id": depID.String()})
	sqlCount, argsCount, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	err = r.db.QueryRow(ctx, sqlCount, argsCount...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := r.builder.Select("user_id, email, first_name, last_name, position, department_id, avatar_url, is_active, created_at, updated_at").
		From("users").
		Where(squirrel.Eq{"department_id": depID.String()}).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *DirectoryRepo) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	query := r.builder.Insert("users").
		Columns("user_id, email, first_name, last_name, position, department_id, avatar_url, is_active, created_at, updated_at").
		Values(user.UserID.String(), user.Email, user.FirstName, user.LastName, nullString(user.Position), nullString(user.DepartmentID), nullString(user.AvatarURL), user.IsActive, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING user_id")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var id string
	err = r.db.QueryRow(ctx, sqlQuery, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(id)
}

func (r *DirectoryRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := r.builder.Select("user_id, email, first_name, last_name, position, department_id, avatar_url, is_active, created_at, updated_at").
		From("users").
		Where(squirrel.Eq{"user_id": id.String()})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sqlQuery, args...)

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *DirectoryRepo) GetUsers(ctx context.Context, filter *dto.GetUsersRequest) ([]*models.User, int, error) {
	qb := r.builder.Select("user_id, email, first_name, last_name, position, department_id, avatar_url, is_active, created_at, updated_at").
		From("users")

	if filter.DepartmentID != nil {
		qb = qb.Where(squirrel.Eq{"department_id": filter.DepartmentID.String()})
	}
	if filter.IsActive != nil {
		qb = qb.Where(squirrel.Eq{"is_active": *filter.IsActive})
	}

	countQuery := r.builder.Select("COUNT(*)").FromSelect(qb, "t")
	sqlCount, argsCount, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	err = r.db.QueryRow(ctx, sqlCount, argsCount...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	qb = qb.Limit(uint64(filter.Limit)).Offset(uint64(filter.Offset))
	sqlQuery, args, err := qb.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *DirectoryRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := r.builder.Delete("users").
		Where(squirrel.Eq{"user_id": id.String()})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sqlQuery, args...)
	return err
}

func nullString(v interface{}) interface{} {
	switch val := v.(type) {
	case *uuid.UUID:
		if val == nil {
			return nil
		}
		return val.String()
	case *string:
		if val == nil {
			return nil
		}
		return *val
	default:
		return v
	}
}

func scanUser(scanner pgx.Row) (*models.User, error) {
	var user models.User
	var position, departmentID, avatarURL sql.NullString
	err := scanner.Scan(
		&user.UserID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&position,
		&departmentID,
		&avatarURL,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if position.Valid {
		user.Position = &position.String
	}
	if departmentID.Valid {
		did, _ := uuid.Parse(departmentID.String)
		user.DepartmentID = &did
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}

	return &user, nil
}
