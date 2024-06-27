package controllers

import (
	"Alabbi_BTEX/config"
	"Alabbi_BTEX/models"
)

type ProcessController struct {
}

func NewProcessController() *ProcessController {
	return &ProcessController{}
}

func (r *ConversionController) GetAll() []models.Process {
	var processes []models.Process
	config.Db.Find(&processes)
	return processes
}
