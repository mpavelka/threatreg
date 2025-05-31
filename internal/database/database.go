package database

import (
	"database/sql"
	"fmt"
	"strings"
	"threatreg/internal/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type DB struct {
	*sqlx.DB
}

var db *DB

// User represents a user in the database
type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Post represents a blog post
type Post struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	UserID    int       `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Connect establishes a database connection
func Connect() error {
	dbURL := config.BuildDatabaseURL()

	// Parse the database URL to determine the driver
	var driverName string
	var dataSourceName string

	if strings.HasPrefix(dbURL, "sqlite3://") {
		driverName = "sqlite3"
		dataSourceName = strings.TrimPrefix(dbURL, "sqlite3://")
	} else if strings.HasPrefix(dbURL, "postgres://") || strings.HasPrefix(dbURL, "postgresql://") {
		driverName = "postgres"
		dataSourceName = dbURL
	} else {
		return fmt.Errorf("unsupported database URL format: %s", dbURL)
	}

	// Connect to database
	sqlxDB, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlxDB.SetMaxOpenConns(25)
	sqlxDB.SetMaxIdleConns(25)
	sqlxDB.SetConnMaxLifetime(5 * time.Minute)

	db = &DB{sqlxDB}

	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Database connected successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// CreateTables creates the database tables (for development only - use migrations in production)
func CreateTables() error {
	if db == nil {
		return fmt.Errorf("database not connected")
	}

	// SQLite-compatible table creation
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	postsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title VARCHAR(200) NOT NULL,
		content TEXT,
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Create indexes
	userIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);",
	}

	postIndexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);",
	}

	// Execute table creation
	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	if _, err := db.Exec(postsTable); err != nil {
		return fmt.Errorf("failed to create posts table: %w", err)
	}

	// Create indexes
	for _, index := range userIndexes {
		if _, err := db.Exec(index); err != nil {
			return fmt.Errorf("failed to create user index: %w", err)
		}
	}

	for _, index := range postIndexes {
		if _, err := db.Exec(index); err != nil {
			return fmt.Errorf("failed to create post index: %w", err)
		}
	}

	fmt.Println("✅ Tables created successfully")
	return nil
}

// User repository methods
func CreateUser(username, email string) (*User, error) {
	query := `INSERT INTO users (username, email) VALUES (?, ?) RETURNING id, created_at`

	// SQLite doesn't support RETURNING, so we use a different approach
	if strings.Contains(config.GetDatabaseURL(), "sqlite") {
		result, err := db.Exec(`INSERT INTO users (username, email) VALUES (?, ?)`, username, email)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to get last insert ID: %w", err)
		}

		// Get the created user
		return GetUserByID(int(id))
	}

	// PostgreSQL with RETURNING
	user := &User{}
	err := db.QueryRowx(query, username, email).StructScan(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.Username = username
	user.Email = email
	return user, nil
}

func GetUserByID(id int) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

func DeleteUser(id int) error {
	result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

func GetUserCount() (int, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users")
	return count, err
}

func GetPostCount() (int, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM posts")
	return count, err
}
