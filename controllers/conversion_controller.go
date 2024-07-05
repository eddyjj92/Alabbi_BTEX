package controllers

import (
	"Alabbi_BTEX/config"
	"Alabbi_BTEX/models"
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type ConversionController struct {
}

func NewConversionController() *ConversionController {
	return &ConversionController{}
}

func (r *ConversionController) Upload(file string, filename string, extension string) (map[string]any, error) {
	// Extraer la parte base64 de la cadena (omitir el prefijo "data:image/jpeg;base64,")
	parts := strings.Split(file, ",")
	if len(parts) != 2 {
		runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Operacion abortada",
			Message:       "Ha ocurrido un error: Archivo no procesable",
			DefaultButton: "Aceptar",
		})
		return nil, fmt.Errorf("Archivo no procesable")
	}
	base64Str := parts[1]
	// Decodificar la cadena base64 a datos binarios
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Operacion abortada",
			Message:       "Ha ocurrido un error: " + err.Error(),
			DefaultButton: "Aceptar",
		})
		return nil, err
	}
	currentTime := time.Now()
	/*tz, _ := time.LoadLocation("America/Havana")*/
	folder := fmt.Sprintf("%s_%s-%s-%s", currentTime.Format(time.DateOnly), strconv.Itoa(currentTime.Hour()), strconv.Itoa(currentTime.Minute()), strconv.Itoa(currentTime.Second()))
	err = os.Mkdir("./storage/files/"+folder, 0755)
	if err != nil {
		runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Operacion abortada",
			Message:       "Ha ocurrido un error: " + err.Error(),
			DefaultButton: "Aceptar",
		})
		return nil, err
	}
	// Escribir los datos binarios en un archivo
	if err := ioutil.WriteFile("./storage/files/"+folder+"/"+filename, data, 0755); err != nil {
		panic("Error writing file: " + err.Error())
	}
	ruta, _ := os.Getwd()
	process := models.Process{
		File:      filename,
		Extension: extension,
		Folder:    folder,
		Route:     fmt.Sprintf("%s", ruta+"\\storage\\files\\"),
		Status:    0,
	}
	config.Db.Create(&process)
	runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
		Type:          runtime.InfoDialog,
		Title:         "Operacion exitosa",
		Message:       "Archivo importado listo para procesar.",
		DefaultButton: "Aceptar",
	})
	return map[string]any{
		"id":         process.ID,
		"filename":   process.File,
		"extension":  extension,
		"folder":     folder,
		"folderPath": fmt.Sprintf("%s", ruta+"\\storage\\files\\"+folder),
	}, nil
}

func (r *ConversionController) OpenVideo(videoPath string) bool {
	// Ruta al ejecutable de KMPlayer
	basePath, _ := os.Getwd()
	args := []string{
		fmt.Sprintf("%s", videoPath),
		fmt.Sprintf("%s", "--video-on-top"),
		fmt.Sprintf("--width=%d", 500),
		fmt.Sprintf("--heigh=%d", 500),
	}
	potplayerPath := basePath + `\storage\auxiliar\VLCPortable\VLCPortable.exe`
	// Crear el comando para abrir KMPlayer con el archivo de video
	cmd := exec.Command(potplayerPath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	go func() {
		// Espera unos segundos para que VLC se inicie y luego centra la ventana
		exec.Command("powershell", "-Command", fmt.Sprintf(`
                Start-Sleep -s 4;
                $window = (Get-Process -Name vlc).MainWindowHandle;
                Add-Type -TypeDefinition @"
                using System;
                using System.Runtime.InteropServices;
                public class User32 {
                    [DllImport("user32.dll", SetLastError = true)]
                    public static extern bool MoveWindow(IntPtr hWnd, int X, int Y, int nWidth, int nHeight, bool bRepaint);
                }
                "@;
                [User32]::MoveWindow($window, %d, %d, %d, %d, $true);
            `, 0, 0, 300, 300)).Run()
	}()
	// Ejecutar el comando
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error" + err.Error())
		panic(err)
	}
	return true
}

func (r *ConversionController) Start(id int, model string, input string, outputDir string, threads int) (bool, error) {
	args := []string{
		fmt.Sprintf("--model=%s", model),
		fmt.Sprintf("%s", input),
		/*fmt.Sprintf("--language=%s", "English"),*/
		fmt.Sprintf("%s", "-pp"),
		/*fmt.Sprintf("--beam_size=%d", 1),
		fmt.Sprintf("--best_of=%d", 1),*/
		fmt.Sprintf("--thread=%d", threads),
		fmt.Sprintf("--output_format=%s", "all"),
		fmt.Sprintf("--output_dir=%s", outputDir),
	}
	ruta, _ := os.Getwd()
	cmd := exec.Command(ruta+`\storage\auxiliar\whisper\whisper-faster.exe`, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := os.Remove(ruta + `\storage\logs\logs.log`)
	if err != nil {
		runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Operacion abortada1",
			Message:       "Ha ocurrido un error: " + err.Error(),
			DefaultButton: "Aceptar",
		})
	}
	// Create a new file for logging requests.
	logFile, err := os.Create(ruta + `\storage\logs\logs.log`)
	if err != nil {
		runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Operacion abortada2",
			Message:       "Ha ocurrido un error: " + err.Error(),
			DefaultButton: "Aceptar",
		})
	}
	defer logFile.Close()
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	stopCh := make(chan struct{})
	go worker(stopCh, ruta)
	// Execute the command
	if err2 := cmd.Run(); err2 != nil {
		runtime.MessageDialog(config.Ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "Operacion abortada3",
			Message:       "Ha ocurrido un error: ",
			DefaultButton: "Aceptar",
		})
	}
	close(stopCh)
	var process models.Process
	config.Db.Find(&process, id)
	process.Status = 1
	config.Db.Save(&process)

	return true, nil
}

func worker(stopCh chan struct{}, ruta string) {
	for {
		select {
		default:
			time.Sleep(time.Second * 1)
			// Abre el archivo
			file, err := os.Open(ruta + `\storage\logs\logs.log`)
			if err != nil {
				fmt.Println("Error al abrir el archivo:", err.Error())
			}
			defer file.Close()
			// Crea un nuevo scanner
			scanner := bufio.NewScanner(file)
			// Lee línea por línea
			currentLine := 1
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(line)

				if currentLine == 10 {
					number := getLastNumber(line)
					runtime.EventsEmit(config.Ctx, "porciento", number)
				}
				currentLine++
			}
			// Verifica si hubo algún error durante la lectura
			if err := scanner.Err(); err != nil {
				fmt.Println("Error al leer el archivo:", err)
			}
		case <-stopCh:
			// Recibimos la señal de detención
			fmt.Println("Deteniendo la goroutine...")
			return
		}
	}
}

func getLastNumber(text string) string {
	fmt.Println(text)
	// Definir la expresión regular para encontrar números entre un espacio y %
	re := regexp.MustCompile(`(\d+(\.\d+)?)%`)

	// Encontrar todas las coincidencias
	matches := re.FindAllStringSubmatch(text, -1)

	// Obtener el último número entre espacio y %
	if len(matches) > 0 {
		lastMatch := matches[len(matches)-1]
		if len(lastMatch) > 1 {
			fmt.Println("Último número encontrado:", lastMatch[1])
			return lastMatch[1]
		}
	} else {
		fmt.Println("No se encontraron números entre espacio y %.")
	}
	return "0"
}

func (r *ConversionController) OutputFile(file string) (bool, error) {
	cmd := exec.Command(
		"notepad.exe",
		[]string{
			fmt.Sprintf("%s", file),
		}...)
	cmd.Run()
	return true, nil
}
