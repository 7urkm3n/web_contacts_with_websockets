package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type ContactModel struct {
	DB *sql.DB
}

type Contact struct {
	ID          uint64 `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	EditHistory []EditHistory `json:"edit_history,omitempty"`
}

type EditHistory struct {
	ID        uint   `json:"id"`
	ContactID uint   `json:"contact_id"`
	Changes   string `json:"changes"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m ContactModel) GetAll() ([]*Contact, error) {
	query := `SELECT id, first_name, last_name, email, phone_number FROM contacts ORDER BY id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*Contact

	for rows.Next() {

		var c Contact

		err := rows.Scan(
			&c.ID,
			&c.FirstName,
			&c.LastName,
			&c.Email,
			&c.PhoneNumber,
		)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (m ContactModel) Insert(contact *Contact) error {
	query := `
			INSERT INTO contacts (first_name, last_name, email, phone_number)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 21*time.Second)
	defer cancel()
	args := []any{contact.FirstName, contact.LastName, contact.Email, contact.PhoneNumber}
	time.Sleep(20 * time.Second) // sleep for task
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&contact.ID, &contact.CreatedAt)
}

func (m ContactModel) Get(id int64) (*Contact, error) {
	if id < 1 {
		fmt.Println("Record not found: ", ErrRecordNotFound)
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, first_name, last_name, email, phone_number, created_at, updated_at
		FROM contacts
		WHERE id = $1`

	var c Contact

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.FirstName,
		&c.LastName,
		&c.Email,
		&c.PhoneNumber,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			fmt.Println("Record not found: ", ErrRecordNotFound)
			return nil, ErrRecordNotFound
		default:
			fmt.Println("Record not found err: ", err)
			return nil, err
		}
	}

	return &c, nil
}

func (m ContactModel) Update(contact *Contact, changes map[string]map[string]string) (*Contact, error) {
	query := `
		UPDATE contacts
		SET first_name = $1, last_name = $2, email = $3, phone_number = $4
		WHERE id = $5 
		RETURNING id`

	args := []any{contact.FirstName, contact.LastName, contact.Email, contact.PhoneNumber, contact.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&contact.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrEditConflict
		default:
			return nil, err
		}
	}

	queryHistory := `INSERT INTO contact_edit_history (contact_id, changes) 
		VALUES ($1, $2)
		RETURNING id`

	jsonData, err := json.Marshal(changes)
	if err != nil {
		return nil, err
	}

	argsHistory := []any{contact.ID, string(jsonData)}
	edit := EditHistory{}
	err = m.DB.QueryRowContext(ctx, queryHistory, argsHistory...).Scan(&edit.ID)
	if err != nil {
		return nil, err
	}

	// edits, _ := m.GetContactHistoryByID(contact.ID)
	// for _, v := range edits {
	// 	contact.EditHistory = append(contact.EditHistory, *v)
	// }

	return contact, nil
}

func (m ContactModel) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM contacts WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m ContactModel) GetContactHistoryByID(id uint64) ([]*EditHistory, error) {
	query := `SELECT * FROM contact_edit_history WHERE contact_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edits []*EditHistory

	for rows.Next() {
		var edit EditHistory
		err := rows.Scan(
			&edit.ID,
			&edit.ContactID,
			&edit.Changes,
			&edit.CreatedAt,
			&edit.UpdatedAt,
		)
		if err != nil {
			fmt.Println("ITS HERER err", err)
			return nil, err
		}
		edits = append(edits, &edit)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return edits, nil
}
