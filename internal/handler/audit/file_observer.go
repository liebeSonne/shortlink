package audit

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/liebeSonne/shortlink/internal/logger"
)

// NewFileObserver - создание экземпляра наблюдателя сохраняющего данные в файл.
func NewFileObserver(
	filePath string,
	logger logger.Logger,
) (*FileObserver, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileObserver{
		filePath: filePath,
		file:     file,
		encoder:  json.NewEncoder(file),
		logger:   logger,
	}, nil
}

// FileObserver - структура наблюдателя сохраняющего данные в файл.
type FileObserver struct {
	filePath string
	file     *os.File
	encoder  *json.Encoder
	mu       sync.Mutex
	logger   logger.Logger
}

func (o *FileObserver) Update(event Event) {
	o.mu.Lock()
	defer o.mu.Unlock()

	err := o.encoder.Encode(&event)
	if err != nil {
		o.logger.Errorf("failed to save event: %v", err)
	}
}

func (o *FileObserver) Close() error {
	return o.file.Close()
}
