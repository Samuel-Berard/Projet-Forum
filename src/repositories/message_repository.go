package repositories

import (
	"database/sql"
	"errors"
	"projet-forum/src/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(m *models.Message) error {
	query := `INSERT INTO messages (contenu, fil_id, auteur_id) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, m.Contenu, m.FilID, m.AuteurID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = int(id)
	return nil
}

func (r *MessageRepository) FindByID(id int) (*models.Message, error) {
	query := `SELECT id, contenu, fil_id, auteur_id, score_popularite FROM messages WHERE id = ?`
	row := r.db.QueryRow(query, id)

	m := &models.Message{}
	err := row.Scan(&m.ID, &m.Contenu, &m.FilID, &m.AuteurID, &m.ScorePopularite)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("message introuvable")
		}
		return nil, err
	}

	return m, nil
}

func (r *MessageRepository) FindByFilID(filID int, page, limit int, sortBy string) ([]models.Message, error) {
	offset := (page - 1) * limit
	orderBy := "created_at DESC"
	if sortBy == "popularite" {
		orderBy = "score_popularite DESC, created_at DESC"
	} else if sortBy == "chronologique" {
		orderBy = "created_at ASC"
	}

	query := `SELECT id, contenu, fil_id, auteur_id, score_popularite, created_at FROM messages WHERE fil_id = ? ORDER BY ` + orderBy + ` LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, filID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		var createdAt []uint8
		if err := rows.Scan(&m.ID, &m.Contenu, &m.FilID, &m.AuteurID, &m.ScorePopularite, &createdAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *MessageRepository) Update(m *models.Message) error {
	query := `UPDATE messages SET contenu = ? WHERE id = ?`
	_, err := r.db.Exec(query, m.Contenu, m.ID)
	return err
}

func (r *MessageRepository) Delete(id int) error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MessageRepository) React(userID, messageID int, typeReaction string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Chercher s'il y a déjà une réaction
	var existType string
	err = tx.QueryRow(`SELECT type FROM reactions WHERE utilisateur_id = ? AND message_id = ?`, userID, messageID).Scan(&existType)

	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}

	scoreChange := 0
	if err == sql.ErrNoRows {
		// Nouvelle réaction
		_, err = tx.Exec(`INSERT INTO reactions (utilisateur_id, message_id, type) VALUES (?, ?, ?)`, userID, messageID, typeReaction)
		if err != nil {
			tx.Rollback()
			return err
		}
		if typeReaction == "like" {
			scoreChange = 1
		} else {
			scoreChange = -1
		}
	} else {
		if existType == typeReaction {
			// L'utilisateur refait la même action : on pourrait annuler la réaction ou renvoyer une erreur
			tx.Rollback()
			return errors.New("vous avez déjà réagi de cette manière à ce message")
		} else {
			// Change le like en dislike ou dislike en like
			_, err = tx.Exec(`UPDATE reactions SET type = ? WHERE utilisateur_id = ? AND message_id = ?`, typeReaction, userID, messageID)
			if err != nil {
				tx.Rollback()
				return err
			}
			if typeReaction == "like" {
				scoreChange = 2 // -1 -> +1
			} else {
				scoreChange = -2 // +1 -> -1
			}
		}
	}

	if scoreChange != 0 {
		_, err = tx.Exec(`UPDATE messages SET score_popularite = score_popularite + ? WHERE id = ?`, scoreChange, messageID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
