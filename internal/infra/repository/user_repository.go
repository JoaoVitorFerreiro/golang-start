package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/JoaoVitorFerreiro/golang-start/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
    Save(user *entity.User) error
    FindByID(id string) (*entity.User, error)
    FindByEmail(email string) (*entity.User, error)
    FindAll() ([]*entity.User, error)
    Delete(id string) error
}

type InMemoryUserRepository struct {
    users map[string]*entity.User
    mutex sync.RWMutex
}

func NewInMemoryUserRepository() UserRepository {
    return &InMemoryUserRepository{
        users: make(map[string]*entity.User),
    }
}

func (r *InMemoryUserRepository) Save(u *entity.User) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    r.users[u.ID] = u
    return nil
}

func (r *InMemoryUserRepository) FindByID(id string) (*entity.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    u, exists := r.users[id]
    if !exists {
        return nil, nil
    }
    return u, nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*entity.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    for _, u := range r.users {
        if u.Email == email {
            return u, nil
        }
    }
    return nil, nil
}

func (r *InMemoryUserRepository) FindAll() ([]*entity.User, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    var users []*entity.User
    for _, u := range r.users {
        users = append(users, u)
    }
    return users, nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    if _, exists := r.users[id]; !exists {
        return errors.New("user not found")
    }
    delete(r.users, id)
    return nil
}

type PostgresUserRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) UserRepository {
    return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Save(u *entity.User) error {
    query := `
        INSERT INTO users (id, name, email, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE SET 
            name = EXCLUDED.name,
            email = EXCLUDED.email,
            updated_at = EXCLUDED.updated_at`
    
    _, err := r.pool.Exec(context.Background(), query, u.ID, u.Name, u.Email, u.CreatedAt, u.UpdatedAt)
    return err
}

func (r *PostgresUserRepository) FindByID(id string) (*entity.User, error) {
    query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
    
    var u entity.User
    err := r.pool.QueryRow(context.Background(), query, id).Scan(
        &u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt,
    )
    
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*entity.User, error) {
    query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = $1`
    
    var u entity.User
    err := r.pool.QueryRow(context.Background(), query, email).Scan(
        &u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt,
    )
    
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *PostgresUserRepository) FindAll() ([]*entity.User, error) {
    query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`
    
    rows, err := r.pool.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*entity.User
    for rows.Next() {
        var u entity.User
        err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, &u)
    }
    
    return users, rows.Err()
}

func (r *PostgresUserRepository) Delete(id string) error {
    query := `DELETE FROM users WHERE id = $1`
    
    result, err := r.pool.Exec(context.Background(), query, id)
    if err != nil {
        return err
    }
    
    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        return errors.New("user not found")
    }
    
    return nil
}


type RepositoryType string

const (
    InMemory RepositoryType = "memory"
    Postgres RepositoryType = "postgres"
)

func NewUserRepository(repoType RepositoryType, pool *pgxpool.Pool) UserRepository {
    switch repoType {
    case InMemory:
        return NewInMemoryUserRepository()
    case Postgres:
        if pool == nil {
            panic("pgxpool is required for postgres repository")
        }
        return NewPostgresUserRepository(pool)
    default:
        return NewInMemoryUserRepository()
    }
}