package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/liebeSonne/shortlink/internal/model"
	"github.com/liebeSonne/shortlink/internal/repository"
)

type shortLinkStorageData struct {
	ID          int        `json:"uuid"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	UserID      *uuid.UUID `json:"user_id"`
}

func NewFileShortLinkRepository(filePath string) (repository.ShortLinkRepositoryWithCloser, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	storage := &fileShortLinkRepository{
		filePath: filePath,
		file:     file,
		lastID:   0,
	}

	err = storage.init()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

type fileShortLinkRepository struct {
	filePath string
	file     *os.File
	lastID   int
	mu       sync.Mutex
}

func (s *fileShortLinkRepository) Find(_ context.Context, shortID string) (*model.ShortLink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	itemPtr, err := s.findItem(shortID)
	if err != nil {
		return nil, fmt.Errorf("failed find item: %w", err)
	}

	if itemPtr != nil {
		return &model.ShortLink{ID: itemPtr.ShortURL, URL: itemPtr.OriginalURL}, nil
	}

	return nil, nil
}

func (s *fileShortLinkRepository) FindByURL(_ context.Context, url string) (*model.ShortLink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	itemPtr, err := s.findItemByURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed find item: %w", err)
	}

	if itemPtr != nil {
		return &model.ShortLink{ID: itemPtr.ShortURL, URL: itemPtr.OriginalURL}, nil
	}

	return nil, nil
}

func (s *fileShortLinkRepository) FindByUserID(_ context.Context, userID uuid.UUID) ([]model.ShortLink, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items, err := s.findItemsByUserUD(userID)
	if err != nil {
		return nil, fmt.Errorf("failed find item: %w", err)
	}

	result := make([]model.ShortLink, 0, len(items))
	for _, item := range items {
		shortLink := model.ShortLink{ID: item.ShortURL, URL: item.OriginalURL}
		result = append(result, shortLink)
	}

	return result, nil
}

func (s *fileShortLinkRepository) Store(_ context.Context, shortLink model.ShortLink, userID *uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nextID := s.nextID()
	item := shortLinkStorageData{
		ID:          nextID,
		ShortURL:    shortLink.ID,
		OriginalURL: shortLink.URL,
		UserID:      userID,
	}

	items := []shortLinkStorageData{item}

	err := s.save(items)
	if err != nil {
		return fmt.Errorf("failed save items: %w", err)
	}
	return nil
}

func (s *fileShortLinkRepository) StoreAll(_ context.Context, shortLinks []model.ShortLink, userID *uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := make([]shortLinkStorageData, len(shortLinks))

	for _, shortLink := range shortLinks {
		nextID := s.nextID()
		item := shortLinkStorageData{
			ID:          nextID,
			ShortURL:    shortLink.ID,
			OriginalURL: shortLink.URL,
			UserID:      userID,
		}
		items = append(items, item)
	}

	err := s.save(items)
	if err != nil {
		return fmt.Errorf("failed save items: %w", err)
	}

	return nil
}

func (s *fileShortLinkRepository) Close() error {
	return s.file.Close()
}

func (s *fileShortLinkRepository) init() error {
	stat, err := s.file.Stat()
	if err != nil {
		return err
	}
	if stat.Size() == 0 {
		writer := bufio.NewWriter(s.file)
		_, err = writer.WriteString("[\n]")
		if err != nil {
			return fmt.Errorf("failed write to file: %w", err)
		}
		err = writer.Flush()
		if err != nil {
			return fmt.Errorf("failed flush: %w", err)
		}
		return nil
	}

	err = s.initLastID()
	if err != nil {
		return err
	}

	return nil
}

func (s *fileShortLinkRepository) initLastID() error {
	lastID, err := s.findLastID()
	if err != nil {
		return fmt.Errorf("failed find last id: %w", err)
	}
	s.lastID = lastID

	return nil
}

func (s *fileShortLinkRepository) findLastID() (int, error) {
	lastID := 0

	_, err := s.file.Seek(0, 0)
	if err != nil {
		return lastID, fmt.Errorf("failed seek file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		b := scanner.Bytes()

		itemPtr, err := s.parseItem(b)
		if err != nil {
			return 0, fmt.Errorf("failed parse item: %w", err)
		}

		if itemPtr != nil {
			if itemPtr.ID > lastID {
				lastID = itemPtr.ID
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		return 0, fmt.Errorf("failed scan file: %w", err)
	}

	return lastID, nil
}

func (s *fileShortLinkRepository) nextID() int {
	nextID := s.lastID + 1
	s.lastID = nextID
	return nextID
}

func (s *fileShortLinkRepository) findItem(id string) (*shortLinkStorageData, error) {
	_, err := s.file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed seek file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		b := scanner.Bytes()

		itemPtr, err := s.parseItem(b)
		if err != nil {
			return nil, fmt.Errorf("failed parse item: %w", err)
		}

		if itemPtr != nil && itemPtr.ShortURL == id {
			return itemPtr, nil
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("failed scan file: %w", err)
	}

	return nil, nil
}

func (s *fileShortLinkRepository) findItemByURL(url string) (*shortLinkStorageData, error) {
	_, err := s.file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed seek file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		b := scanner.Bytes()

		itemPtr, err := s.parseItem(b)
		if err != nil {
			return nil, fmt.Errorf("failed parse item: %w", err)
		}

		if itemPtr != nil && itemPtr.OriginalURL == url {
			return itemPtr, nil
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("failed scan file: %w", err)
	}

	return nil, nil
}

func (s *fileShortLinkRepository) findItemsByUserUD(userID uuid.UUID) ([]shortLinkStorageData, error) {
	_, err := s.file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed seek file: %w", err)
	}

	result := make([]shortLinkStorageData, 0)

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		b := scanner.Bytes()

		itemPtr, err := s.parseItem(b)
		if err != nil {
			return nil, fmt.Errorf("failed parse item: %w", err)
		}

		if itemPtr != nil && itemPtr.UserID != nil && *(itemPtr.UserID) == userID {
			result = append(result, *itemPtr)
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("failed scan file: %w", err)
	}

	return result, nil
}

func (s *fileShortLinkRepository) parseItem(b []byte) (*shortLinkStorageData, error) {
	str := string(b)

	if str == "[" || str == "]" {
		return nil, nil
	}

	str = strings.TrimPrefix(str, ",")
	str = strings.TrimSuffix(str, ",")

	var item shortLinkStorageData
	err := json.Unmarshal([]byte(str), &item)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal item: %w", err)
	}

	return &item, nil
}

func (s *fileShortLinkRepository) save(items []shortLinkStorageData) error {
	tmpDir, err := os.MkdirTemp("", "tmp-*")
	defer func() {
		err = os.RemoveAll(tmpDir)
		if err != nil {
			fmt.Printf("failed remove tmp dir: %v", err)
		}
	}()

	tmpFile, err := os.CreateTemp(tmpDir, "tmp.storage.*.json")
	if err != nil {
		return fmt.Errorf("failed create tmp file: %w", err)
	}

	writer := bufio.NewWriter(tmpFile)

	err = s.savaItemsToWriter(writer, items)
	if err != nil {
		return fmt.Errorf("failed save to file: %w", err)
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("failed flush writer: %w", err)
	}

	err = tmpFile.Close()
	if err != nil {
		return fmt.Errorf("failed close tmp file: %w", err)
	}

	err = s.file.Close()
	if err != nil {
		return fmt.Errorf("failed close file: %w", err)
	}

	err = os.Rename(tmpFile.Name(), s.filePath)
	if err != nil {
		return fmt.Errorf("failed rename file: %w", err)
	}

	err = s.reopenFile()
	if err != nil {
		return fmt.Errorf("failed to reopen file: %w", err)
	}

	return nil
}

func (s *fileShortLinkRepository) reopenFile() error {
	file, err := os.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed reopen file: %w", err)
	}

	s.file = file
	return nil
}

func (s *fileShortLinkRepository) savaItemsToWriter(writer *bufio.Writer, items []shortLinkStorageData) error {
	if len(items) == 0 {
		return nil
	}

	_, err := s.file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed seek file: %w", err)
	}

	scanner := bufio.NewScanner(s.file)

	for scanner.Scan() {
		b := scanner.Bytes()
		str := string(b)

		if str != "]" {
			if str != "[" && !strings.HasSuffix(str, ",") {
				str = str + ","
			}
			str += "\n"

			_, err = writer.WriteString(str)
			if err != nil {
				return fmt.Errorf("failed write to file: %w", err)
			}
		} else {
			for i, item := range items {
				data, err := json.Marshal(&item)
				if err != nil {
					return fmt.Errorf("failed encode item: %w", err)
				}

				newStr := string(data)
				if i < len(items)-1 {
					newStr += ","
				}
				newStr += "\n"

				_, err = writer.WriteString(newStr)
				if err != nil {
					return fmt.Errorf("failed write new item: %w", err)
				}
			}

			_, err = writer.WriteString("]")
			if err != nil {
				return fmt.Errorf("failed write end array: %w", err)
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("failed scan file: %w", err)
	}

	return nil
}
