package controllers

import (
	"Alabbi_BTEX/config"
	"Alabbi_BTEX/models"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
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

func (r *ConversionController) Delete(id int) (models.Process, error) {
	var process models.Process
	config.Db.First(&process, id)
	if process.ID == 0 {
		runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Operación errónea",
			Message:       "Este Producto no existe",
			DefaultButton: "Aceptar",
		})
		return process, fmt.Errorf("Este Producto no existe")
	}
	config.Db.Delete(&process, id)
	// Eliminar un archivo
	if err := os.Remove("./storage/files/" + process.Folder + "/" + process.File); err != nil {
		runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
			Type:          runtime.WarningDialog,
			Title:         "Advertencia",
			Message:       "No se ha encontrado el archivo multimedia: '" + process.File + "' en el sistema, se eliminará el registro en la base de datos",
			DefaultButton: "Aceptar",
		})
		os.Remove("./storage/files/" + process.Folder)
	}
	runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "Operación exitosa",
		Message:       "Archivo eliminado con éxito",
		DefaultButton: "Aceptar",
	})
	return process, nil
}

func (r *ConversionController) DeleteHistory() (bool, error) {
	var processes []models.Process
	config.Db.Find(&processes)

	for _, el := range processes {
		config.Db.Delete(&el, el.ID)
		// Eliminar un archivo
		if err := os.RemoveAll("./storage/files/" + el.Folder); err != nil {
			fmt.Println(err.Error())
		}
	}

	runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "Operación exitosa",
		Message:       "Historial borrado con éxito",
		DefaultButton: "Aceptar",
	})
	return true, nil
}
