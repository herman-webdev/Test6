package user

import (
	"awesomeProject/internal/user/model"
	"awesomeProject/internal/user/storage"
	"awesomeProject/pkg/client/postgresql"
	"awesomeProject/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) Create(ctx context.Context, user *model.User) error {
	q := `
			INSERT INTO users 
			    (name, surname, patronymic, age, gender, country_id) 
			VALUES 
			    ($1, $2, $3, $4, $5, $6) 
			RETURNING id
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	if err := r.client.QueryRow(
		ctx,
		q,
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Age,
		user.Gender,
		user.CountryId).Scan(&user.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf(
				"SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message,
				pgErr.Detail,
				pgErr.Where,
				pgErr.Code,
				pgErr.SQLState()))

			r.logger.Error(newErr)

			return newErr
		}

		return err
	}

	return nil
}

func (r *repository) FindAll(ctx context.Context, sortOptions storage.SortOptions) ([]model.User, error) {
	orderBy := sortOptions.GetOrderBy()

	q := `
        SELECT id, name, surname, patronymic, age, gender, country_id, created_at
        FROM public.users
        ORDER BY ` + orderBy + `;
    `

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	var users []model.User

	for rows.Next() {
		var usr model.User

		err = rows.Scan(&usr.ID, &usr.Name, &usr.Surname, &usr.Patronymic, &usr.Age, &usr.Gender, &usr.CountryId, &usr.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, usr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (model.User, error) {
	q := `
		SELECT id, name, surname, patronymic, age, gender, country_id, created_at, updated_at FROM public.users WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var usr model.User
	err := r.client.QueryRow(ctx, q, id).Scan(&usr.ID, &usr.Name, &usr.Surname, &usr.Patronymic, &usr.Age, &usr.Gender, &usr.CountryId, &usr.CreatedAt, &usr.UpdatedAt)
	if err != nil {
		return model.User{}, err
	}

	return usr, nil
}

func (r *repository) Update(ctx context.Context, user *model.User, id string) error {
	q := `UPDATE public.users SET updated_at = now() AT TIME ZONE 'utc'`

	var params []interface{}
	var paramIndex = 1

	if user.Name != "" {
		q += fmt.Sprintf(", name = $%d", paramIndex)
		params = append(params, user.Name)
		paramIndex++
	}

	if user.Surname != "" {
		q += fmt.Sprintf(", surname = $%d", paramIndex)
		params = append(params, user.Surname)
		paramIndex++
	}

	if user.Patronymic != "" {
		q += fmt.Sprintf(", patronymic = $%d", paramIndex)
		params = append(params, user.Patronymic)
		paramIndex++
	}

	if user.Age != 0 {
		q += fmt.Sprintf(", age = $%d", paramIndex)
		params = append(params, user.Age)
		paramIndex++
	}

	if user.Gender != "" {
		q += fmt.Sprintf(", gender = $%d", paramIndex)
		params = append(params, user.Gender)
		paramIndex++
	}

	if user.CountryId != "" {
		q += fmt.Sprintf(", country_id = $%d", paramIndex)
		params = append(params, user.CountryId)
		paramIndex++
	}

	q += fmt.Sprintf(" WHERE id = $%d RETURNING id, name, surname, patronymic, age, gender, country_id, created_at, updated_at", paramIndex)
	params = append(params, id)

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var usr model.User
	err := r.client.QueryRow(ctx, q, params...).
		Scan(&usr.ID, &usr.Name, &usr.Surname, &usr.Patronymic, &usr.Age, &usr.Gender, &usr.CountryId, &usr.CreatedAt, &usr.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	q := `
		DELETE FROM public.users WHERE id = $1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) storage.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
