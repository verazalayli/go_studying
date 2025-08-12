package memory

import (
	"errors"
	"github.com/verazalayli/go_studying/grpc/pkg/service"
	"sync"
)

// Простое in-memory хранилище (адаптер к порту service.NoteRepository)
type NoteRepo struct {
	mu    sync.RWMutex
	items map[string]service.Note
}

func NewNoteRepo() *NoteRepo {
	return &NoteRepo{items: make(map[string]service.Note)}
}

func (r *NoteRepo) Save(n service.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[n.ID] = n
	return nil
}

func (r *NoteRepo) GetByID(id string) (service.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, ok := r.items[id]
	if !ok {
		return service.Note{}, errors.New("not found")
	}
	return n, nil
}

func (r *NoteRepo) List() ([]service.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]service.Note, 0, len(r.items))
	for _, n := range r.items {
		out = append(out, n)
	}
	return out, nil
}
