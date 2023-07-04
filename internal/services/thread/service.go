package thread

import (
	"context"
	"strings"

	ev "github.com/amleonc/evalrunes"
	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/dto"
	"github.com/amleonc/tabula/internal/services/comment"
	"github.com/amleonc/tabula/internal/services/media"
	"github.com/amleonc/tabula/internal/services/thread/internal"
	"github.com/gofrs/uuid"
)

// ------------ Types ------------ //

type Service interface {
	Create(context.Context, *dto.Thread) (*dto.Thread, error)
	GetOneByID(context.Context, uuid.UUID) (*dto.Thread, error)
	GetByTopicID(context.Context, uuid.UUID) ([]*dto.Thread, error)
}

type serviceStruct struct {
	repo internal.Repostiroy
}

type ServiceError struct {
	msg string
}

// ------------ Constants ------------ //

const (
	MinThreadTitleLen = 1
	MaxThreadTitleLen = 100
	MinThreadBodyLen  = 1
	MaxThreadBodyLen  = 3000

	limitQuery = 40
)

// ------------ Variables ------------ //

var (
	cs = comment.GetService()
	ms = media.GetService()

	service = newService()

	errInvalidTitle = "error: title length must be between 1 and 100 alphanumeric characters, emojis and/or spaces"
	errInvalidBody  = "error: body length must be between 1 and 3000 alphanumeric characters, emojis and/or spaces"
)

// ------------ Methods ------------ //

func (s *serviceStruct) Create(ctx context.Context, t *dto.Thread) (*dto.Thread, error) {
	var err error
	if err = sanitizeThread(t); err != nil {
		return nil, err
	}
	t.Media, err = ms.Create(ctx, t.Media)
	if err != nil {
		return nil, err
	}
	daoT := &dao.Thread{
		Media: t.Media.ID,
		Topic: t.Topic,
		User:  t.User,
		Title: t.Title,
		Body:  t.Body,
		NSFW:  t.NSFW,
	}
	if err = s.repo.Create(ctx, daoT); err != nil {
		return nil, err
	}
	t = daoToDto(daoT)
	return t, nil
}

func (s *serviceStruct) GetOneByID(ctx context.Context, id uuid.UUID) (*dto.Thread, error) {
	// get comments concurrently
	fn := func(cc chan<- []*dto.Comment, errc chan<- error) {
		defer close(cc)
		defer close(errc)
		comments, err := cs.GetThreadComments(ctx, id)
		if err != nil {
			cc <- nil
			errc <- err
		}
		cc <- comments
		errc <- nil
	}
	cc := make(chan []*dto.Comment, 1)
	errc := make(chan error, 1)
	go fn(cc, errc)
	fn1 := func(tc chan<- *dao.Thread, errc chan<- error) {
		defer close(tc)
		defer close(errc)
		t, err := s.repo.SelectByID(ctx, id)
		if err != nil {
			tc <- nil
			errc <- err
		}
		tc <- t
		errc <- nil
	}
	tc := make(chan *dao.Thread, 1)
	errc1 := make(chan error, 1)
	go fn1(tc, errc1)
	comments, err1 := <-cc, <-errc
	daoT, err2 := <-tc, <-errc1
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	t := daoToDto(daoT)
	t.Comments = comments
	return t, nil
}

func (s *serviceStruct) GetByTopicID(ctx context.Context, id uuid.UUID) ([]*dto.Thread, error) {
	daoThreads, err := s.repo.SelectByTopicID(ctx, id, limitQuery)
	if err != nil {
		return nil, err
	}
	for _, v := range daoThreads {
		v.User = uuid.UUID{}
	}
	dtoThreads := make([]*dto.Thread, len(daoThreads))
	for i, daot := range daoThreads {
		dtoThreads[i] = daoToDto(daot)
	}
	return dtoThreads, nil
}

func (s ServiceError) Error() string {
	return s.msg
}

// ------------ Functions ------------ //

func GetService() Service {
	return service
}

func newService() *serviceStruct {
	return &serviceStruct{internal.Repostiroy{}}
}

func newServiceError(msg string) ServiceError {
	return ServiceError{msg}
}

func sanitizeThread(t *dto.Thread) error {
	t.Title = strings.TrimSpace(t.Title)
	var ok bool
	if ok = validateThreadTitle(t.Title); !ok {
		return newServiceError(errInvalidTitle)
	}
	t.Body = strings.TrimSpace(t.Body)
	if ok = validateThreadBodyLen(t.Body); !ok {
		return newServiceError(errInvalidBody)
	}
	t.Comments = nil // flush it, just in case
	return nil
}

func validateThreadTitle(t string) bool {
	r := []rune(t)
	return ev.CompareAgainstValidators(
		r,
		ev.IsAlphanumeric,
		ev.IsSpace,
		ev.IsEmoji,
		ev.IsPunctuation,
	) &&
		ev.ValidateRuneArrayLength(r, MinThreadTitleLen, MaxThreadTitleLen)
}

func validateThreadBodyLen(b string) bool {
	r := []rune(b)
	return ev.CompareAgainstValidators(
		r,
		ev.IsAlphanumeric,
		ev.IsSpace,
		ev.IsEmoji,
		ev.IsPunctuation,
	) && ev.ValidateRuneArrayLength(r, MinThreadBodyLen, MaxThreadBodyLen)
}

func daoToDto(daoT *dao.Thread) *dto.Thread {
	dtoT := &dto.Thread{
		ID:        daoT.ID,
		Topic:     daoT.Topic,
		User:      daoT.User,
		Comments:  nil,
		CreatedAt: daoT.CreatedAt,
		UpdatedAt: daoT.UpdatedAt,
	}
	return dtoT
}
