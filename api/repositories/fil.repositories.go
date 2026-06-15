package repositories

import (
	"database/sql"
	"errors"

	"projet-forum/api/models"
)

type FilRepository struct {
	db *sql.DB
}

func InitFilRepository(db *sql.DB) *FilRepository {
	return &FilRepository{db: db}
}

func (r *FilRepository) Create(f *models.FilDeDiscussion, categoryIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO fils_de_discussion (titre, etat, auteur_id) VALUES (?, ?, ?)`
	result, err := tx.Exec(query, f.Titre, f.Etat, f.AuteurID)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	f.ID = int(id)

	if len(categoryIDs) > 0 {
		catQuery := `INSERT INTO fils_categories (fil_id, categorie_id) VALUES (?, ?)`
		for _, catID := range categoryIDs {
			_, err = tx.Exec(catQuery, f.ID, catID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *FilRepository) FindAll(page, limit int, search string, categoryID int) ([]models.FilDeDiscussion, error) {
	offset := (page - 1) * limit
	args := []interface{}{}

	query := `
		SELECT DISTINCT f.id, f.titre, f.etat, f.auteur_id, f.created_at, f.updated_at
		FROM fils_de_discussion f
	`

	if categoryID > 0 {
		query += ` INNER JOIN fils_categories fc ON f.id = fc.fil_id `
	}

	query += ` WHERE f.etat != 'archive' `

	if search != "" {
		query += ` AND f.titre LIKE ? `
		args = append(args, "%"+search+"%")
	}

	if categoryID > 0 {
		query += ` AND fc.categorie_id = ? `
		args = append(args, categoryID)
	}

	query += ` ORDER BY f.created_at DESC LIMIT ? OFFSET ? `
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fils []models.FilDeDiscussion
	for rows.Next() {
		var f models.FilDeDiscussion
		var createdAt, updatedAt []uint8
		if err := rows.Scan(&f.ID, &f.Titre, &f.Etat, &f.AuteurID, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		fils = append(fils, f)
	}

	return fils, nil
}

func (r *FilRepository) CountAll(search string, categoryID int) (int, error) {
	args := []interface{}{}
	query := `SELECT COUNT(DISTINCT f.id) FROM fils_de_discussion f`

	if categoryID > 0 {
		query += ` INNER JOIN fils_categories fc ON f.id = fc.fil_id `
	}

	query += ` WHERE f.etat != 'archive' `

	if search != "" {
		query += ` AND f.titre LIKE ? `
		args = append(args, "%"+search+"%")
	}

	if categoryID > 0 {
		query += ` AND fc.categorie_id = ? `
		args = append(args, categoryID)
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FilRepository) FindByID(id int) (*models.FilDeDiscussion, error) {
	query := `SELECT id, titre, etat, auteur_id, created_at, updated_at FROM fils_de_discussion WHERE id = ?`
	row := r.db.QueryRow(query, id)

	f := &models.FilDeDiscussion{}
	var createdAt, updatedAt []uint8
	err := row.Scan(&f.ID, &f.Titre, &f.Etat, &f.AuteurID, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("fil introuvable")
		}
		return nil, err
	}

	return f, nil
}

func (r *FilRepository) Update(f *models.FilDeDiscussion) error {
	query := `UPDATE fils_de_discussion SET titre = ?, etat = ? WHERE id = ?`
	_, err := r.db.Exec(query, f.Titre, f.Etat, f.ID)
	return err
}

func (r *FilRepository) Delete(id int) error {
	query := `DELETE FROM fils_de_discussion WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
