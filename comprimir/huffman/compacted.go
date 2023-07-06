package huffman

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
)

/*
* Crea un string codificado en huffman segun los datos del arbol de frecuencias
* y el texto original, en otras palabras codifica a huffman el texto original.
 */
func Compacted(texto string, arbol *arbol) string {
	codigos := make(map[rune]string)
	obtenerCodigos(arbol, "", codigos) // Generar los códigos Huffman a partir del árbol

	var compactado string
	var buffer bytes.Buffer
	for _, ch := range texto {
		buffer.WriteString(codigos[ch])
	}
	compactado = buffer.String()
	return compactado
}

/* ----- NOTA PARA MI -------
*  Creo que es este el que cambie
*  para que siempre devuelva
*  el mismo arbol para los
*  mismos datos, copiar del
*  codigo paralelo
* --------------------------
 */
func obtenerCodigos(nodo *arbol, prefijo string, codigos map[rune]string) {
	if nodo == nil {
		return
	}

	if nodo.izq == nil && nodo.der == nil {
		// Es un nodo hoja, asignar el prefijo como código Huffman para el carácter
		codigos[nodo.c] = prefijo
		return
	}
	// Recursivamente obtener los códigos Huffman de los subárboles izquierdo y derecho
	obtenerCodigos(nodo.izq, prefijo+"0", codigos)
	obtenerCodigos(nodo.der, prefijo+"1", codigos)
}

func BinaryToBytes(binaryString string) []byte {
	// Convierte el string binario en un número entero sin signo
	number, _ := strconv.ParseUint(binaryString, 2, len(binaryString))

	// Convierte el número entero en un slice de bytes
	bytes := make([]byte, len(binaryString)/8)
	for i := range bytes {
		bytes[i] = byte(number >> (8 * (len(bytes) - 1 - i)))
	}

	return bytes
}

func SaveCompacted(compacted string, raiz *arbol) error {
	//Obtiene la longitud original de los datos
	originalLength := len(compacted)

	//Calcula la cantidad de ceros adicionales necesarios para que la longitud sea multiplo de 8
	extraZeros := 8 - (originalLength % 8)

	compacted += strings.Repeat("0", extraZeros)

	file, err := os.Create("./comprimir/resultados/comprimido.huf")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Escribe la longitud original de los datos en el archivo
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(originalLength))
	writer.Write(lengthBytes)

	for i := 0; i < len(compacted); i += 8 {
		byteStr := compacted[i : i+8]
		byteVal := byte(0)
		for j := 0; j < 8; j++ {
			if byteStr[j] == '1' {
				byteVal |= 1 << (7 - j)
			}
		}
		writer.WriteByte(byteVal)
	}

	err = guardarArbol(raiz, writer)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func guardarArbol(raiz *arbol, writer *bufio.Writer) error {
	if raiz == nil {
		writer.WriteByte(0)
		return nil
	}

	writer.WriteByte(1) //Escribe un byte 1 para indicar nodo valido
	writer.WriteByte(byte(raiz.freq))
	writer.WriteRune(raiz.c) //escribe el caracter del nodo

	err := guardarArbol(raiz.izq, writer) //Recursion para guardar subarbol izq
	if err != nil {
		return err
	}

	err = guardarArbol(raiz.der, writer)
	if err != nil {
		return err
	}

	return nil
}

func SaveCompactedAmbos(compacted string, raiz *arbol) error {
	//Obtiene la longitud original de los datos
	originalLength := len(compacted)

	//Calcula la cantidad de ceros adicionales necesarios para que la longitud sea multiplo de 8
	extraZeros := 8 - (originalLength % 8)

	compacted += strings.Repeat("0", extraZeros)

	file, err := os.Create("./ambos/resultados/comprimido.huf")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Escribe la longitud original de los datos en el archivo
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(originalLength))
	writer.Write(lengthBytes)

	for i := 0; i < len(compacted); i += 8 {
		byteStr := compacted[i : i+8]
		byteVal := byte(0)
		for j := 0; j < 8; j++ {
			if byteStr[j] == '1' {
				byteVal |= 1 << (7 - j)
			}
		}
		writer.WriteByte(byteVal)
	}

	err = guardarArbol(raiz, writer)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
