package adapter

import (
	"context"
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
		Values(dep.DepartmentID, dep.Name, dep.ParentID, dep.CreatedAt, dep.UpdatedAt).
		Suffix("RETURNING department_id")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = r.db.QueryRow(ctx, sqlQuery, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *DirectoryRepo) GetDepartmentByID(ctx context.Context, id uuid.UUID) (*models.Department, error) {
	query := r.builder.Select("department_id, name, parent_id, created_at, updated_at").
		From("departments").
		Where(squirrel.Eq{"department_id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sqlQuery, args...)

	var dep models.Department
	err = row.Scan(&dep.DepartmentID, &dep.Name, &dep.ParentID, &dep.CreatedAt, &dep.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &dep, nil
}

func (r *DirectoryRepo) GetDepartments(ctx context.Context, limit, offset int) ([]*models.Department, int, error) {
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
		err = rows.Scan(&dep.DepartmentID, &dep.Name, &dep.ParentID, &dep.CreatedAt, &dep.UpdatedAt)
		if err != nil {
			return nil, 0, err
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
		err = rows.Scan(&dep.DepartmentID, &dep.Name, &dep.ParentID, &dep.CreatedAt, &dep.UpdatedAt)
		if err != nil {
			return nil, err
		}
		deps = append(deps, &dep)
	}

	return deps, nil
}

func (r *DirectoryRepo) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	query := r.builder.Delete("departments").
		Where(squirrel.Eq{"department_id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sqlQuery, args...)
	return err
}

func (r *DirectoryRepo) GetDepartmentMembers(ctx context.Context, depID uuid.UUID, limit, offset int) ([]*models.User, int, error) {
	countQuery := r.builder.Select("COUNT(*)").From("users").
		Where(squirrel.Eq{"department_id": depID})
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
		Where(squirrel.Eq{"department_id": depID}).
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
		Values(user.UserID, user.Email, user.FirstName, user.LastName, user.Position, user.DepartmentID, user.AvatarURL, user.IsActive, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING user_id")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = r.db.QueryRow(ctx, sqlQuery, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *DirectoryRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := r.builder.Select("user_id, email, first_name, last_name, position, department_id, avatar_url, is_active, created_at, updated_at").
		From("users").
		Where(squirrel.Eq{"user_id": id})

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
		qb = qb.Where(squirrel.Eq{"department_id": *filter.DepartmentID})
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
		Where(squirrel.Eq{"user_id": id})

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, sqlQuery, args...)
	return err
}

func scanUser(scanner pgx.Row) (*models.User, error) {
	var user models.User
	err := scanner.Scan(
		&user.UserID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Position,
		&user.DepartmentID,
		&user.AvatarURL,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
