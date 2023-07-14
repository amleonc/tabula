package media

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/amleonc/tabula/config"
	"github.com/amleonc/tabula/internal/dao"
	"github.com/amleonc/tabula/internal/dto"
	"github.com/amleonc/tabula/internal/fs"
	"github.com/amleonc/tabula/internal/services/media/internal"
)

// Types

type Service interface {
	Create(context.Context, *dto.Media) (*dto.Media, error)
	GetByID(context.Context, string) (*dto.Media, error)
}

type serviceStruct struct {
	repo internal.Repo
}

type Error struct {
	msg string
}

// Constants

const (
	chunkSize = 4096 // buffer size used to calculate the hash.

	errHash          = "error: cannot generate file hash"
	errForbidden     = "error: forbidden file"
	errNoMIME        = "error: cannot extract file MIME type from headers"
	errInvalidFormat = "error: invalid format (valid formats: jpeg, jpg, png, gif, webp, mp4, webm)"

	thumbnailPrefix        = "_thumb_"
	defaultThumbnailFormat = "jpeg"
)

// Variables

var (
	service          = newService()
	thumbnailFormat  = config.ThumbnailFormat()
	thumbnailSize    = config.ThumbnailSize()
	uploadsDirectory = config.UploadsDir()
)

// Methods

func (s *serviceStruct) Create(ctx context.Context, m *dto.Media) (*dto.Media, error) {
	var err error
	f := m.Bytes
	id, err := generateSHA256HashFromFileHeader(f)
	if err != nil {
		return nil, newMediaError(errHash)
	}
	r, err := s.repo.Select(ctx, &dao.Media{ID: id})
	if err == nil {
		if err = isBlacklisted(r); err != nil {
			return nil, err
		}
		m = daoToDto(r)
		return m, nil
	} else {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	mimetype, err := getFileMimetype(m.Bytes)
	if err != nil {
		return nil, newMediaError(errNoMIME)
	}
	fileType, extension, found := strings.Cut(mimetype, "/")
	if !found {
		return nil, newMediaError(errNoMIME)
	}
	if !validateFileExtension(extension) {
		return nil, newMediaError(errInvalidFormat)
	}
	fileName := buildFileName(id, extension)
	_, _ = f.Seek(0, io.SeekStart)
	err = fs.SaveToDisk(f, fileName)
	if err != nil {
		return nil, err
	}
	thumbnailFilename := buildFileName(thumbnailPrefix+id, thumbnailFormat)
	err = generateThumbnail(fileName, thumbnailFilename)
	if err != nil {
		return nil, err
	}
	r = &dao.Media{ID: id, Type: fileType, Extension: extension}
	r, err = s.repo.Insert(ctx, r)
	if err != nil {
		return nil, err
	}
	m = daoToDto(r)
	return m, nil
}

func (s *serviceStruct) GetByID(ctx context.Context, id string) (*dto.Media, error) {
	r, err := s.repo.Select(ctx, &dao.Media{ID: id})
	if err != nil {
		return nil, err
	}
	m := daoToDto(r)
	return m, nil
}

func (m Error) Error() string {
	return m.msg
}

// Functions

func GetService() Service {
	return service
}

func newService() *serviceStruct {
	return &serviceStruct{internal.Repo{}}
}

func newMediaError(msg string) Error {
	return Error{msg}
}

func generateSHA256HashFromFileHeader(file io.ReadSeekCloser) (string, error) {
	hash := sha256.New()
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		_, err = hash.Write(buf[:n])
		if err != nil {
			return "", err
		}
	}
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func buildFileName(hash, extension string) string {
	return fmt.Sprintf("%s/%s.%s", uploadsDirectory, hash, extension)
}

func getFileMimetype(f io.ReadSeekCloser) (string, error) {
	buffer := make([]byte, 512)
	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}
	fileType := http.DetectContentType(buffer)
	if fileType == "" {
		return "", err
	}
	return fileType, nil
}

func validateFileExtension(ext string) bool {
	switch ext {
	case "jpeg", "jpg", "png", "gif", "webp", "mp4", "webm":
		return true
	default:
		return false
	}
}

func isBlacklisted(m *dao.Media) error { // -
	switch m.Blacklist { // -----------------
	case true: // ---------------------------
		return newMediaError(errForbidden) //
	default: // -----------------------------
		return nil // ------------------------
	} // -------------------------------------
} // -----------------------------------------

func daoToDto(daom *dao.Media) *dto.Media {
	m := &dto.Media{
		Id:           daom.ID,
		Url:          fmt.Sprintf("%s.%s", daom.ID, daom.Extension),
		ThumbnailUrl: fmt.Sprintf("%s%s.%s", thumbnailPrefix, daom.ID, defaultThumbnailFormat),
		Type:         daom.Type,
		Format:       daom.Extension,
		Blacklist:    daom.Blacklist,
		CreatedAt:    daom.CreatedAt,
		UpdatedAt:    daom.UpdatedAt,
	}
	return m
}

func generateThumbnail(name, thumb string) error {
	scale := thumbnailSize
	args := make([]string, 0)
	args = append(
		args,
		"-i", name,
		"-vf", "select=eq(n\\,0)",
		"-vf", fmt.Sprintf("scale=%s:force_original_aspect_ratio=decrease", scale),
		"-vframes", "1",
		thumb,
	)
	cmd := exec.Command("ffmpeg", args...)
	err := cmd.Run()
	if err != nil && err.Error() != "exit status 1" {
		return err
	}
	return nil
}
