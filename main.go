package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	// Carpeta actual del programa
	carpetaActual, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Obtener la lista de archivos en la carpeta actual
	archivos, err := ioutil.ReadDir(carpetaActual)
	if err != nil {
		log.Fatal(err)
	}

	// Expresión regular para encontrar el número de boleto
	patronNumeroBoleto := regexp.MustCompile(`\b[0-9]{10}\b`)

	// Recorrer cada archivo en la carpeta actual
	for _, archivo := range archivos {
		// Verificar si el archivo es un archivo de texto
		if !archivo.IsDir() && filepath.Ext(archivo.Name()) == ".txt" {
			// Leer el contenido del archivo
			contenido, err := ioutil.ReadFile(filepath.Join(carpetaActual, archivo.Name()))
			if err != nil {
				log.Printf("Error al leer el archivo %s: %v\n", archivo.Name(), err)
				continue
			}

			// Buscar el número de boleto en el contenido
			matches := patronNumeroBoleto.FindStringSubmatch(string(contenido))
			if len(matches) > 0 {
				numeroBoleto := matches[0]

				// Verificar si es un "agent coupon"
				esAgentCoupon := regexp.MustCompile(`(?i)AGENT COUPON`).MatchString(string(contenido))

				// Nuevo nombre del archivo
				nuevoNombre := numeroBoleto
				if esAgentCoupon {
					nuevoNombre += " C"
				}
				nuevoNombre += ".txt"

				// Renombrar el archivo
				err = os.Rename(filepath.Join(carpetaActual, archivo.Name()), filepath.Join(carpetaActual, nuevoNombre))
				if err != nil {
					log.Printf("Error al renombrar el archivo %s: %v\n", archivo.Name(), err)
					continue
				}

				fmt.Printf("Archivo %s renombrado a %s\n", archivo.Name(), nuevoNombre)
			} else {
				log.Printf("No se encontró un número de boleto en el archivo %s\n", archivo.Name())
			}
		}
	}

	fmt.Println("Proceso completado.")
}
