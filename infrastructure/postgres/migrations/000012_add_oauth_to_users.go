package migrations

import (
	"gorm.io/gorm"
)

// AddOAuthToUsers adds OAuth support fields to users table
func AddOAuthToUsers(db *gorm.DB) error {
	// Add new columns for OAuth support
	if err := db.Exec(`
		ALTER TABLE users
		ADD COLUMN IF NOT EXISTS oauth_provider VARCHAR(50),
		ADD COLUMN IF NOT EXISTS oauth_id VARCHAR(255),
		ADD COLUMN IF NOT EXISTS is_oauth_user BOOLEAN DEFAULT FALSE;
	`).Error; err != nil {
		return err
	}

	// Add indexes for OAuth fields
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_oauth_provider ON users(oauth_provider);
		CREATE INDEX IF NOT EXISTS idx_users_oauth_id ON users(oauth_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_oauth_provider_id ON users(oauth_provider, oauth_id)
		WHERE oauth_provider IS NOT NULL AND oauth_id IS NOT NULL;
	`).Error; err != nil {
		return err
	}

	// Make password nullable for OAuth users
	if err := db.Exec(`
		ALTER TABLE users ALTER COLUMN password DROP NOT NULL;
	`).Error; err != nil {
		return err
	}

	return nil
}

// RollbackOAuthFromUsers removes OAuth support from users table
func RollbackOAuthFromUsers(db *gorm.DB) error {
	if err := db.Exec(`
		DROP INDEX IF EXISTS idx_users_oauth_provider_id;
		DROP INDEX IF EXISTS idx_users_oauth_id;
		DROP INDEX IF EXISTS idx_users_oauth_provider;
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		ALTER TABLE users
		DROP COLUMN IF EXISTS is_oauth_user,
		DROP COLUMN IF EXISTS oauth_id,
		DROP COLUMN IF EXISTS oauth_provider;
	`).Error; err != nil {
		return err
	}

	// Make password NOT NULL again
	if err := db.Exec(`
		ALTER TABLE users ALTER COLUMN password SET NOT NULL;
	`).Error; err != nil {
		return err
	}

	return nil
}
