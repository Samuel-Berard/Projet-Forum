package repositories

import (
	"database/sql"
	"errors"
	"projet-forum/src/models"
)

type UtilisateurRepository struct {
	db *sql.DB
}

func NewUtilisateurRepository(db *sql.DB) *UtilisateurRepository {
	return &UtilisateurRepository{db: db}
}

func (r *UtilisateurRepository) Create(u *models.Utilisateur) error {
	query := `INSERT INTO utilisateurs (username, email, password_hash) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, u.Username, u.Email, u.PasswordHash)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

func (r *UtilisateurRepository) FindByUsernameOrEmail(identifiant string) (*models.Utilisateur, error) {
	query := `SELECT id, username, email, password_hash, role, banned, created_at FROM utilisateurs WHERE username = ? OR email = ?`
	row := r.db.QueryRow(query, identifiant, identifiant)

	u := &models.Utilisateur{}
	var createdAt []uint8 // Temporary handle for MySQL timestamp mapping to time.Time
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.Banned, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("utilisateur introuvable")
		}
		return nil, err
	}
	// Note: We ignore createdAt parsing for simplicity here if not strictly needed, 
	// or we can use time.Parse depending on MySQL driver config (parseTime=true).
	return u, nil
}

func (r *UtilisateurRepository) FindByID(id int) (*models.Utilisateur, error) {
	query := `SELECT id, username, email, password_hash, role, banned FROM utilisateurs WHERE id = ?`
	row := r.db.QueryRow(query, id)

	u := &models.Utilisateur{}
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.Banned)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("utilisateur introuvable")
		}
		return nil, err
	}
	return u, nil
}

func (r *UtilisateurRepository) BanUser(id int) error {
	query := `UPDATE utilisateurs SET banned = TRUE WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
