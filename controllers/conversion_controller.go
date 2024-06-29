package controllers

import (
	"Alabbi_BTEX/config"
	"Alabbi_BTEX/models"
	"encoding/base64"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConversionController struct {
}

func NewConversionController() *ConversionController {
	return &ConversionController{}
}

func (r *ConversionController) Upload(file string, filename string) map[string]any {

	// Extraer la parte base64 de la cadena (omitir el prefijo "data:image/jpeg;base64,")
	parts := strings.Split(file, ",")
	if len(parts) != 2 {
		panic("Invalid base64 string")
	}
	base64Str := parts[1]

	// Decodificar la cadena base64 a datos binarios
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		panic("Error decoding base64 string: " + err.Error())
	}

	currentTime := time.Now()
	/*tz, _ := time.LoadLocation("America/Havana")*/
	folder := fmt.Sprintf("%s_%s-%s-%s", currentTime.Format(time.DateOnly), strconv.Itoa(currentTime.Hour()), strconv.Itoa(currentTime.Minute()), strconv.Itoa(currentTime.Second()))

	err = os.Mkdir("./storage/files/"+folder, 0755)
	if err != nil {
		panic("Error make dir: " + err.Error())
	}

	// Escribir los datos binarios en un archivo
	if err := ioutil.WriteFile("./storage/files/"+folder+"/"+filename, data, 0755); err != nil {
		panic("Error writing file: " + err.Error())
	}

	ruta, _ := os.Getwd()

	process := models.Process{
		File:      filename,
		Extension: ".mp4",
		Folder:    folder,
		Route:     fmt.Sprintf("%s", ruta+"\\files\\"),
	}
	config.Db.Create(&process)

	runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "Operacion exitosa",
		Message:       "Archivo importado listo para procesar.",
		DefaultButton: "Aceptar",
	})

	return map[string]any{
		"error":     nil,
		"message":   "Archivo importado listo para procesar.",
		"validator": nil,
		"filename":  process.File,
		"extension": ".mp4",
		"folder":    folder,
		"outputDir": fmt.Sprintf("%s", ruta+"\\files\\"+folder),
		"ruta":      fmt.Sprintf("%s", ruta+"\\files\\"+folder+"\\"+process.File),
	}

}
