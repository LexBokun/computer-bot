package multimonitor

import (
	"context"
	"encoding/csv"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/LexBokun/ControlAgent/internal/domain/display"
	"github.com/LexBokun/ControlAgent/internal/infrastructure/multi-monitor/model"
)

func New(path string) *Config {
	return &Config{Path: path}
}

// List — список дисплеев.
func (t *Config) List(ctx context.Context) ([]display.Display, error) {
	tmpFile, err := os.CreateTemp("", "monitors-*.csv")
	if err != nil {
		slog.Error("error creating temporary file:", "error", err.Error())
		return nil, err
	}

	tmpFilePath := tmpFile.Name()
	err = tmpFile.Close()
	if err != nil {
		return nil, err
	}
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			slog.Error("error deleting file:", "error", err.Error())
			return
		}
	}(tmpFilePath)

	//nolint:gosec // MultiMonitorTool trust
	// Запускаем MultiMonitorTool
	cmd := exec.CommandContext(ctx, t.Path, "/scomma", tmpFilePath)
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	// Парсим CSV-файл
	file, err := os.Open(tmpFilePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Error("error closing file:", "error", err.Error())
		}
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var dtoList []model.Display
	for i, record := range records {
		if i == 0 {
			continue // Пропустить заголовок
		}
		if len(record) <= 20 { // Минимальная длинна файла
			continue
		}

		id := strings.TrimSpace(record[12])             // Monitor ID
		name := strings.TrimSpace(record[20])           // Monitor Name
		enable := strings.TrimSpace(record[3]) == "Yes" // Active

		dtoList = append(dtoList, model.Display{
			ID:     id,
			Name:   name,
			Enable: enable,
		})
	}

	var result []display.Display
	for _, dto := range dtoList {
		result = append(result, model.ToDomain(dto))
	}

	return result, nil
}

// SetEnabled — включить или выключить дисплей.
func (t *Config) SetEnabled(ctx context.Context, id string, enabled bool) error {
	action := "/enable"
	if !enabled {
		action = "/disable"
	}

	//nolint:gosec // MultiMonitorTool trust
	cmd := exec.CommandContext(ctx, t.Path, action, id)
	return cmd.Run()
}
