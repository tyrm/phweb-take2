package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID uint `json:"-"`

	Subject           string  `json:"sub"`
	Email             *string `json:"email"`
	Name              *string `json:"name"`
	PreferredUsername *string `json:"preferred_username"`

	// unsaved values
	Audience        string   `json:"aud"`
	AuthnCtxClsRef  string   `json:"acr"`
	AuthorizedParty string   `json:"azp"`
	AuthTime        int64    `json:"auth_time"`
	EmailVerified   bool     `json:"email_verified"`
	ExpiresAt       int64    `json:"exp"`
	Groups          []string `json:"groups"`
	IssuedAt        int64    `json:"iat"`
	Issuer          string   `json:"iss"`
	SessionState    string   `json:"session_state"`
	TokenID         string   `json:"jti"`
	Type            string   `json:"typ"`

	// metadata
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

const sqlCreateUser = `
INSERT INTO "public"."users" (subject, email, name, preferred_username)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;`

const sqlReadUser = `
SELECT id, subject, email, name, preferred_username, created_at, updated_at
FROM "public"."users"
WHERE subject = $1;`

const sqlUpdateUser = `
UPDATE "public"."users"
SET email = $2, name = $3, preferred_username = $4,updated_at = now()
WHERE id = $1
RETURNING updated_at;`

func (u *User) Upsert() error {
	var newID uint
	var newSubject string
	var newEmail string
	var newName string
	var newPreferredUsername string
	var newCreatedAt time.Time
	var newUpdatedAt time.Time

	// try to get user
	err := db.QueryRow(sqlReadUser, u.Subject).Scan(&newID, &newSubject, &newEmail, &newName,
		&newPreferredUsername, &newCreatedAt, &newUpdatedAt)

	if err == sql.ErrNoRows {
		// Create user
		err = db.QueryRow(sqlCreateUser, u.Subject, u.Email, u.Name, u.PreferredUsername).
			Scan(&newID, &newCreatedAt, &newUpdatedAt)
		if err != nil {
			logger.Errorf("Could not create OktaUser: %s", err.Error())
			return err
		}

		// update object
		u.ID = newID
		u.CreatedAt = newCreatedAt
		u.UpdatedAt = newUpdatedAt

		return nil
	} else if err != nil {
		logger.Errorf("Could not get OktaUser from database: %s", err.Error())
		return err
	}

	// Update
	err = db.QueryRow(sqlUpdateUser, newID, u.Email, u.Name, u.PreferredUsername).
		Scan(&newUpdatedAt)
	if err != nil {
		logger.Errorf("Could not update OktaUser: %s", err.Error())
		return err
	}

	// update object
	u.ID = newID
	u.CreatedAt = newCreatedAt
	u.UpdatedAt = newUpdatedAt
	return nil
}
