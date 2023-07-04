package comment

import (
	"context"
	"math/rand"
	"time"

	"github.com/amleonc/tabula/config"
	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/dto"
	"github.com/amleonc/tabula/internal/helpers"
	"github.com/amleonc/tabula/internal/services/comment/internal"
	"github.com/amleonc/tabula/internal/services/media"
	"github.com/gofrs/uuid"
)

// ------------ Types ------------ //

type Service interface {
	Create(context.Context, *dto.Comment) (*dto.Comment, error)
	GetThreadComments(context.Context, uuid.UUID) ([]*dto.Comment, error)
}

type serviceStruct struct {
	repo internal.Repository
}

type ServiceError struct {
	msg string
}

// ------------ Constants ------------ //

const (
	GripLength = 7

	probGroup1 = 4
	probGroup2 = 1
	probGroup3 = 1
	probGroup4 = 1

	limitResults = 20
)

// ------------ Variables ------------ //

var (
	ms = media.GetService()

	userID = config.UserIdKey()

	colorRandomizer = rand.New(rand.NewSource(time.Now().UnixNano()))

	service = newService()
)

// ------------ Methods ------------ //

func (s *serviceStruct) Create(ctx context.Context, c *dto.Comment) (*dto.Comment, error) {
	var err error
	c.Media, err = ms.Create(ctx, c.Media)
	if err != nil {
		return nil, err
	}
	c.Grip = generateGrip()
	c.Color = randColor()
	daoC := &dao.Comment{
		Media:  c.Media.ID,
		Thread: c.Thread,
		User:   c.User,
		Grip:   c.Grip,
		Body:   c.Body,
		Color:  c.Color,
		IsOP:   c.IsOP,
	}
	err = s.repo.Create(ctx, daoC)
	if err != nil {
		return nil, err
	}
	c.ID = daoC.ID
	c.CreatedAt = daoC.CreatedAt
	c.UpdatedAt = daoC.UpdatedAt
	return c, nil
}

func (s *serviceStruct) GetThreadComments(ctx context.Context, id uuid.UUID) ([]*dto.Comment, error) {
	daoc, err := s.repo.SelectByThreadID(ctx, id, limitResults)
	if err != nil {
		return nil, err
	}
	dtoc := toDtoComment(daoc...)
	isOP(ctx.Value(userID).(uuid.UUID), dtoc...)
	return dtoc, nil
}

func (s ServiceError) Error() string {
	return s.msg
}

// ------------ Functions ------------ //

func GetService() Service {
	return service
}

func newService() *serviceStruct {
	return &serviceStruct{internal.Repository{}}
}

func generateGrip() string {
	return helpers.GenStringWithLength(GripLength)
}

// Colors:
//
// 0 = Red
// 1 = Blue
// 2 = Yellow
// 3 = Green
// 4 = Purple
// 5 = Orange
// 6 = Multicolor
// 7 = Negative multicolor
// 8 = White
// 9 = Black

func randColor() uint8 {
	// Generate a random number between 0 and the sum of probabilities
	random := colorRandomizer.Intn(probGroup1 + probGroup2 + probGroup3 + probGroup4)
	// Check the range and return the corresponding number
	if random < probGroup1 {
		return uint8(random % 4) // Group 1: 0 to 3
	} else if random < probGroup1+probGroup2 {
		return uint8(4 + (random % 2)) // Group 2: 4 to 5
	} else if random < probGroup1+probGroup2+probGroup3 {
		return uint8(6 + (random % 2)) // Group 3: 6 to 7
	} else {
		return uint8(8 + (random % 2)) // Group 4: 8 to 9
	}
}

func toDtoComment(daoC ...*dao.Comment) []*dto.Comment {
	dtoc := make([]*dto.Comment, len(daoC))
	for i, daoc := range daoC {
		dtoc[i].ID = daoc.ID
		dtoc[i].Thread = daoc.Thread
		dtoc[i].User = daoc.User
		dtoc[i].Grip = daoc.Grip
		dtoc[i].Body = daoc.Body
		dtoc[i].Color = daoc.Color
		dtoc[i].CreatedAt = daoc.CreatedAt
		dtoc[i].UpdatedAt = daoc.UpdatedAt
		dtoc[i].DeletedAt = daoc.DeletedAt
	}
	return dtoc
}

func isOP(opID uuid.UUID, comments ...*dto.Comment) {
	for _, c := range comments {
		c.IsOP = opID == c.User
	}
}