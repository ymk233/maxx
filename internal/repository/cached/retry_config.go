package cached

import (
    "sync"

    "github.com/Bowl42/maxx-next/internal/domain"
    "github.com/Bowl42/maxx-next/internal/repository"
)

type RetryConfigRepository struct {
    repo         repository.RetryConfigRepository
    cache        map[uint64]*domain.RetryConfig
    defaultCache *domain.RetryConfig
    mu           sync.RWMutex
}

func NewRetryConfigRepository(repo repository.RetryConfigRepository) *RetryConfigRepository {
    return &RetryConfigRepository{
        repo:  repo,
        cache: make(map[uint64]*domain.RetryConfig),
    }
}

func (r *RetryConfigRepository) Load() error {
    list, err := r.repo.List()
    if err != nil {
        return err
    }
    r.mu.Lock()
    defer r.mu.Unlock()
    for _, c := range list {
        r.cache[c.ID] = c
        if c.IsDefault {
            r.defaultCache = c
        }
    }
    return nil
}

func (r *RetryConfigRepository) Create(c *domain.RetryConfig) error {
    if err := r.repo.Create(c); err != nil {
        return err
    }
    r.mu.Lock()
    r.cache[c.ID] = c
    if c.IsDefault {
        r.defaultCache = c
    }
    r.mu.Unlock()
    return nil
}

func (r *RetryConfigRepository) Update(c *domain.RetryConfig) error {
    if err := r.repo.Update(c); err != nil {
        return err
    }
    r.mu.Lock()
    r.cache[c.ID] = c
    if c.IsDefault {
        r.defaultCache = c
    }
    r.mu.Unlock()
    return nil
}

func (r *RetryConfigRepository) Delete(id uint64) error {
    if err := r.repo.Delete(id); err != nil {
        return err
    }
    r.mu.Lock()
    if r.defaultCache != nil && r.defaultCache.ID == id {
        r.defaultCache = nil
    }
    delete(r.cache, id)
    r.mu.Unlock()
    return nil
}

func (r *RetryConfigRepository) GetByID(id uint64) (*domain.RetryConfig, error) {
    r.mu.RLock()
    if c, ok := r.cache[id]; ok {
        r.mu.RUnlock()
        return c, nil
    }
    r.mu.RUnlock()
    return r.repo.GetByID(id)
}

func (r *RetryConfigRepository) GetDefault() (*domain.RetryConfig, error) {
    r.mu.RLock()
    if r.defaultCache != nil {
        r.mu.RUnlock()
        return r.defaultCache, nil
    }
    r.mu.RUnlock()
    return r.repo.GetDefault()
}

func (r *RetryConfigRepository) List() ([]*domain.RetryConfig, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    list := make([]*domain.RetryConfig, 0, len(r.cache))
    for _, c := range r.cache {
        list = append(list, c)
    }
    return list, nil
}
