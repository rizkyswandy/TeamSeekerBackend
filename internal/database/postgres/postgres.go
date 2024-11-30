package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rizkyswandy/TeamSeekerBackend/api"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(connString string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

// Creating profile
func (p *PostgresDB) CreateProfile(profile *api.StudentProfile) error {
    query := `
        INSERT INTO student_profiles 
        (name, email, faculty, field_of_study, semester, skills, focus, is_available)
        VALUES ($1, $2, $3, $4, $5, $6::text[], $7::text[], $8)
        RETURNING id, created_at, updated_at`

    err := p.db.QueryRow(
        query,
        profile.Name,
        profile.Email,
        profile.Faculty,
        profile.FieldOfStudy,
        profile.Semester,
        pq.Array(profile.Skills),
        pq.Array(profile.Focus),
        profile.IsAvailable,
    ).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)

    if err != nil {
        log.Printf("Database error: %v", err)
        return err
    }

    return nil
}

func (p *PostgresDB) GetProfile(id string) (api.StudentProfile, error) {
	var profile api.StudentProfile

	query := `
		SELECT id, name, email, faculty, field_of_study, semester, skills, focus, 
			is_available, created_at, updated_at
		FROM student_profiles WHERE id = $1`

	profileID, err := uuid.Parse(id)
	if err != nil {
		return api.StudentProfile{}, fmt.Errorf("invalid ID format")
	}

	err = p.db.QueryRow(query, profileID).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Email,
		&profile.Faculty,
		&profile.FieldOfStudy,
		&profile.Semester,
		pq.Array(&profile.Skills),
		pq.Array(&profile.Focus),
		&profile.IsAvailable,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No profile found with id: %s", id)
			return api.StudentProfile{}, fmt.Errorf("profile not found")
		}
		log.Printf("Database error getting profile %s: %v", id, err)
		return api.StudentProfile{}, err
	}

	return profile, nil
}

func (p *PostgresDB) UpdateProfile(id string, profile *api.StudentProfile) error {
	profileID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	query := `
		UPDATE student_profiles 
		SET name = $1,
			email = $2,
			faculty = $3,
			field_of_study = $4,
			semester = $5,
			skills = $6,
			focus = $7,
			is_available = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $9`

	result, err := p.db.Exec(
		query,
		profile.Name,
		profile.Email,
		profile.Faculty,
		profile.FieldOfStudy,
		profile.Semester,
		pq.Array(profile.Skills),
		pq.Array(profile.Focus),
		profile.IsAvailable,
		profileID,
	)

	if err != nil {
		log.Printf("Database error updating profile %s: %v", id, err)
		return err
	}

	// Check if any row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("profile not found")
	}

	return nil
}

func (p *PostgresDB) DeleteProfile(id string) error {
	query := `
		DELETE FROM student_profiles WHERE id = $1`

	profileID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	result, err := p.db.Exec(query, profileID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("profile with id %s not found", id)
	}

	return nil
}

func (p *PostgresDB) GetAllProfiles() ([]api.StudentProfile, error) {
	query := `
		SELECT *
		FROM student_profiles`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []api.StudentProfile

	for rows.Next() {
		var profile api.StudentProfile

		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Email,
			&profile.Faculty,
			&profile.FieldOfStudy,
			&profile.Semester,
			pq.Array(&profile.Skills),
			pq.Array(&profile.Focus),
			&profile.IsAvailable,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (p *PostgresDB) SearchProfiles(filter api.SearchFilters) ([]api.StudentProfile, error) {
	query := `
		SELECT id, name, email, faculty, field_of_study, semester, skills, focus, is_available, 
			   created_at, updated_at
		FROM student_profiles
		WHERE 1=1`

	var params []interface{}
	paramCount := 1

	if filter.Faculty != "" {
		query += fmt.Sprintf(" AND faculty = $%d", paramCount)
		params = append(params, filter.Faculty)
		paramCount++
	}

	if len(filter.Skills) > 0 {
		query += fmt.Sprintf(" AND skills && $%d", paramCount)
		params = append(params, pq.Array(filter.Skills))
		paramCount++
	}

	if len(filter.Focus) > 0 {
		query += fmt.Sprintf(" AND focus && $%d", paramCount)
		params = append(params, pq.Array(filter.Focus))
		paramCount++
	}

	query += fmt.Sprintf(" AND is_available = $%d", paramCount)
	params = append(params, filter.Availability)

	rows, err := p.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []api.StudentProfile

	for rows.Next() {
		var profile api.StudentProfile
		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Email,
			&profile.Faculty,
			&profile.FieldOfStudy,
			&profile.Semester,
			pq.Array(&profile.Skills),
			pq.Array(&profile.Focus),
			&profile.IsAvailable,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}