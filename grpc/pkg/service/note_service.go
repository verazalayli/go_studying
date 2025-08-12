package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Доменная модель
type Note struct {
	ID        string
	Title     string
	Content   string
	CreatedAt time.Time
}

// Порт хранилища
type NoteRepository interface {
	Save(n Note) error
	GetByID(id string) (Note, error)
	List() ([]Note, error)
}

// Ошибки прикладного слоя
var (
	ErrNotFound   = errors.New("note not found")
	ErrBadRequest = errors.New("bad request")
)

// Входной порт прикладного слоя (то, что вызывает handler)
type NoteService interface {
	Create(ctx context.Context, title, content string) (Note, error)
	Get(ctx context.Context, id string) (Note, error)
	List(ctx context.Context) ([]Note, error)
}

type noteService struct {
	repo NoteRepository
}

func NewNoteService(repo NoteRepository) NoteService {
	return &noteService{repo: repo}
}

func (s *noteService) Create(ctx context.Context, title, content string) (Note, error) {
	if title == "" {
		return Note{}, ErrBadRequest
	}
	n := Note{
		ID:        uuid.NewString(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Save(n); err != nil {
		return Note{}, err
	}
	return n, nil
}

func (s *noteService) Get(ctx context.Context, id string) (Note, error) {
	n, err := s.repo.GetByID(id)
	if err != nil {
		return Note{}, ErrNotFound
	}
	return n, nil
}

func (s *noteService) List(ctx context.Context) ([]Note, error) {
	return s.repo.List()
}
