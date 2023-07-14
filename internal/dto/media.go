package dto

import (
	"github.com/lestrrat-go/option"
	"github.com/uptrace/bun"
	"io"
	"time"
)

// Types

type Media struct {
	bun.BaseModel `bun:"table:media"`
	Id            string            `json:"id,omitempty" bun:",pk,type:varchar(64)"`
	Url           string            `json:"url,omitempty" bun:",nullzero,notnull,unique"`
	ThumbnailUrl  string            `json:"thumbnail_url,omitempty" bun:",nullzero,notnull,unique"`
	TypeId        uint8             `json:"-" bun:",nullzero,notnull"`
	Type          *Type             `json:"type,omitempty" bun:",nullzero,notnull,rel:has-one,join:type_id=id"`
	FormatId      uint8             `json:"-" bun:",nullzero,notnull"`
	Format        *Format           `json:"format,omitempty" bun:",nullzero,notnull,rel:has-one,join:format_id=id"`
	Bytes         io.ReadSeekCloser `json:"file,omitempty" bun:"-"`
	Blacklist     bool              `json:"blacklist,omitempty" bun:",nullzero,notnull,default:false"`
	CreatedAt     *time.Time        `json:"created_at,omitempty" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time        `json:"updated_at,omitempty" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     *time.Time        `json:"deleted_at,omitempty" bun:",nullzero,soft_delete"`
}

type Type struct {
	bun.BaseModel `bun:"table:media_types"`
	Id            uint8      `json:"id,omitempty" bun:",pk,autoincrement"`
	Type          string     `json:"type,omitempty" bun:",nullzero,notnull,unique"`
	CreatedAt     *time.Time `json:"created_at,omitempty" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" bun:",nullzero,soft_delete"`
}

type Format struct {
	bun.BaseModel `bun:"table:media_formats"`
	Id            uint8      `json:"id,omitempty" bun:",pk,autoincrement"`
	Format        string     `json:"format,omitempty" bun:",nullzero,notnull,unique"`
	CreatedAt     *time.Time `json:"created_at,omitempty" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" bun:",nullzero,soft_delete"`
}

type mediaType string

type MediaOption interface {
	Option
	mediaMethod() mediaType
}

type mediaStr struct {
	Option
}

type identMediaId struct{}
type identMediaUrl struct{}
type identMediaThumbUrl struct{}
type identMediaTypeId struct{}
type identMediaType struct{}
type identMediaFormatId struct{}
type identMediaFormat struct{}
type identMediaBytes struct{}
type identMediaBlacklist struct{}

// Methods

func (mediaStr) mediaMethod() mediaType { return "" }

// Constructor

func NewMedia(opts ...MediaOption) Media {
	m := Media{
		BaseModel:    bun.BaseModel{},
		Id:           "",
		Url:          "",
		ThumbnailUrl: "",
		TypeId:       0,
		Type:         nil,
		FormatId:     0,
		Format:       nil,
		Bytes:        nil,
		Blacklist:    false,
		CreatedAt:    nil,
		UpdatedAt:    nil,
		DeletedAt:    nil,
	}

	for _, o := range opts {
		switch o.Ident() {
		case &identMediaId{}:
			m.Id = o.Value().(string)
		case &identMediaUrl{}:
			m.Url = o.Value().(string)
		case &identMediaThumbUrl{}:
			m.ThumbnailUrl = o.Value().(string)
		case &identMediaTypeId{}:
			m.TypeId = o.Value().(uint8)
		case &identMediaType{}:
			m.Type = o.Value().(*Type)
		case &identMediaFormatId{}:
			m.FormatId = o.Value().(uint8)
		case &identMediaFormat{}:
			m.Format = o.Value().(*Format)
		case &identMediaBytes{}:
			m.Bytes = o.Value().(io.ReadSeekCloser)
		case &identMediaBlacklist{}:
			m.Blacklist = o.Value().(bool)
		}
	}

	return m
}

func MediaId(id string) MediaOption {
	return &mediaStr{option.New(&identMediaId{}, id)}
}

func MediaUrl(u string) MediaOption {
	return &mediaStr{option.New(&identMediaUrl{}, u)}
}

func MediaThumbUrl(u string) MediaOption {
	return &mediaStr{option.New(&identMediaThumbUrl{}, u)}
}

func MediaTypeId(id uint8) MediaOption {
	return &mediaStr{option.New(&identMediaTypeId{}, id)}
}

func MediaType(t Type) MediaOption {
	return &mediaStr{option.New(&identMediaType{}, &t)}
}

func MediaFormatId(id uint8) MediaOption {
	return &mediaStr{option.New(&identMediaFormatId{}, id)}
}

func MediaFormat(f Format) MediaOption {
	return &mediaStr{option.New(&identMediaFormat{}, &f)}
}

func MediaBytes(b io.ReadSeekCloser) MediaOption {
	return &mediaStr{option.New(&identMediaBytes{}, b)}
}

func MediaBlacklist(b bool) MediaOption {
	return &mediaStr{option.New(&identMediaBlacklist{}, b)}
}
