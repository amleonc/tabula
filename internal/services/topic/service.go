package topic

import (
	"context"
	"strings"

	ev "github.com/amleonc/evalrunes"
	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/dto"
	"github.com/amleonc/tabula/internal/services/media"
	"github.com/amleonc/tabula/internal/services/thread"
	"github.com/amleonc/tabula/internal/services/topic/internal"
	"github.com/gofrs/uuid"
)

// ------------ Types ------------ //

type Service interface {
	Create(context.Context, *dto.Topic) (*dto.Topic, error)
	ReadOneView(context.Context, uuid.UUID) (*dto.Topic, error)
}

type serviceStruct struct {
	repo internal.Repository
}

type ServiceError struct {
	msg string
}

// ------------ Constants ------------ //

const (
	MinTitleLength      = 5
	MaxTitleLength      = 30
	MinShortTitleLength = 2
	MaxShortTitleLength = 4
	MinThreads          = 12
	MaxThreads          = 128
)

// ------------ Variables ------------ //

var (
	service = newService()
	ms      = media.GetService()
	ts      = thread.GetService()

	errInvalidShortTitle = "error: invalid short title, its length must be between 2 and 4 alphanumeric characters"
	errInvalidTitle      = "error: invalid title, its length must be between 5 and 30 alphanumeric or space characters"
	errThreads           = "error: thread capacity must be between 12 and 128"
)

// ------------ Methods ------------ //

func (s *serviceStruct) Create(ctx context.Context, t *dto.Topic) (*dto.Topic, error) {
	var err error
	if err = sanitizeTopic(t); err != nil {
		return nil, err
	}
	t.Media, err = ms.Create(ctx, t.Media)
	if err != nil {
		return nil, err
	}
	daoT := &dao.Topic{
		Media:      t.Media.ID,
		User:       t.User,
		Title:      t.Title,
		ShortTitle: t.ShortTitle,
		NSFW:       t.NSFW,
		MaxThreads: t.MaxThreads,
	}
	err = s.repo.Create(ctx, daoT)
	if err != nil {
		return nil, err
	}
	t.ID = daoT.ID
	t.CreatedAt = daoT.CreatedAt
	t.UpdatedAt = daoT.UpdatedAt
	return t, nil
}

func (s *serviceStruct) ReadOneView(ctx context.Context, id uuid.UUID) (*dto.Topic, error) {
	daot, err := s.repo.SelectOneByID(ctx, id)
	if err != nil {
		return nil, err
	}
	m, err := ms.GetByID(ctx, daot.Media)
	if err != nil {
		return nil, err
	}
	threads, err := ts.GetByTopicID(ctx, id)
	if err != nil {
		return nil, err
	}
	t := daoToDto(daot)
	t.Threads = threads
	t.Media = m
	return t, nil
}

func (se ServiceError) Error() string {
	return se.msg
}

// ------------ Functions ------------ //

func GetService() Service {
	return service
}

func newService() *serviceStruct {
	return &serviceStruct{internal.Repository{}}
}

func newServiceError(msg string) *ServiceError {
	return &ServiceError{msg}
}

func sanitizeTopic(t *dto.Topic) error {
	var ok bool
	t.ShortTitle = strings.TrimSpace(t.ShortTitle)
	if ok = validateTopicShortTitle(t.ShortTitle); !ok {
		return newServiceError(errInvalidShortTitle)
	}
	t.Title = strings.TrimSpace(t.Title)
	if ok = validateTopicTitle(t.Title); !ok {
		return newServiceError(errInvalidTitle)
	}
	if ok = validateThreadCapacity(t.MaxThreads); !ok {
		return newServiceError(errThreads)
	}
	return nil
}

func validateTopicTitle(n string) bool {
	r := []rune(n)
	return ev.CompareAgainstValidators(r, ev.IsAlphanumeric, ev.IsSpace, ev.IsEmoji) &&
		ev.ValidateRuneArrayLength(r, MinTitleLength, MaxTitleLength)
}

func validateTopicShortTitle(s string) bool {
	r := []rune(s)
	return ev.CompareAgainstValidators(r, ev.IsAlphanumeric) &&
		ev.ValidateRuneArrayLength(r, MinShortTitleLength, MaxShortTitleLength)
}

func validateThreadCapacity(c int64) bool {
	return MinThreads <= c && c <= MaxThreads
}

func daoToDto(daot *dao.Topic) *dto.Topic {
	t := &dto.Topic{
		ID:         daot.ID,
		Title:      daot.Title,
		ShortTitle: daot.ShortTitle,
		NSFW:       daot.NSFW,
		MaxThreads: daot.MaxThreads,
		CreatedAt:  daot.CreatedAt,
		UpdatedAt:  daot.UpdatedAt,
	}
	return t
}
