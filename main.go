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
	patronNumeroBoleto := regexp.MustCompile(`([0-9]{10})`)

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
			if len(matches) > 1 {
				numeroBoleto := matches[1]

				// Verificar si es un "agent coupon"
				esAgentCoupon := false
				if contenido != nil && bytesContainsIgnoreCase(contenido, []byte("AGENT COUPON")) {
					esAgentCoupon = true
				}

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

// Función auxiliar para buscar una subcadena en un slice de bytes sin importar mayúsculas o minúsculas
func bytesContainsIgnoreCase(s []byte, substr []byte) bool {
	n := len(substr)
	if n == 0 {
		return true
	}
	if len(s) < n {
		return false
	}
	lowerSubstr := bytesToLower(substr)
	for i := 0; i <= len(s)-n; i++ {
		if bytesEqualFold(s[i:i+n], lowerSubstr) {
			return true
		}
	}
	return false
}

// Función auxiliar para convertir un slice de bytes a minúsculas
func bytesToLower(b []byte) []byte {
	res := make([]byte, len(b))
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			res[i] = c + 'a' - 'A'
		} else {
			res[i] = c
		}
	}
	return res
}

// Función auxiliar para comparar dos slices de bytes sin importar mayúsculas o minúsculas
func bytesEqualFold(b1 []byte, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}
	for i := 0; i < len(b1); i++ {
		c1 := b1[i]
		c2 := b2[i]
		if c1 != c2 && (c1|0x20) != (c2|0x20) {
			return false
		}
	}
	return true
}
