package model

import "github.com/LexBokun/ControlAgent/internal/domain/display"

type Display struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

func ToDomain(dto Display) display.Display {
	return display.Display{
		ID:     dto.ID,
		Name:   dto.Name,
		Enable: dto.Enable,
	}
}

func ToDTO(domain display.Display) Display {
	return Display{
		ID:     domain.ID,
		Name:   domain.Name,
		Enable: domain.Enable,
	}
}
