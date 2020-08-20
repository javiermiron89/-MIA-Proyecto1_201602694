package analizador

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"../funciones"
	"../metodos"
)

var palabras []string

var black = "\u001B[30m"
var red = "\033[1;31m"
var green = "\033[1;32m"
var yellow = "\033[1;33m"
var blue = "\033[1;34m"
var magenta = "\033[1;35m"
var cyan = "\033[1;36m"
var white = "\033[1;37m"
var reset = "\u001B[0m"
var whiteBACKGROUND = "\u001B[47m"
var blackBACKGROUND = "\u001B[40m"
var redBACKGROUND = "\u001B[41m"
var greenBACKGROUND = "\u001B[42m"
var yellowBACKGROUND = "\u001B[43m"
var blueBACKGROUND = "\u001B[44m"
var purpleBACKGROUND = "\u001B[45m"
var cyanBACKGROUND = "\u001B[46m"

//Global =
var Global string

//InitApp =
func InitApp() string {
	return Global
}

//Iniciar_Analisis =
func Iniciar_Analisis(cadena string) []string {
	cadenaDividida := strings.SplitN(cadena, " ", -1)

	//fmt.Println("Inicia la cadena separada: ")
	//Variable tipo bool que indica si alguna comilla ha sido encontrada
	var encontreComilla = false
	//Variable para almacenar los valores entre comillas temporalmente
	var temporal string
	//Cadena encargada de almacenar los nuevos valores dividos
	var cadenaReordenadaDividida []string
	/*
	*Ciclo for que se encarga de recorer la cadena principal de ingreso, para luego reordenarla
	*y poder verificar la exitencia de comillas para no separla por espacios
	 */
	//Este for se encarga de verificar el ultimo split y ver si tiene un salto de linea para eliminarlo
	for i := 0; i < len(cadenaDividida); i++ {
		if i == len(cadenaDividida)-1 {
			if (cadenaDividida[i][len(cadenaDividida[i])-1:]) == "\n" {
				r := []rune(cadenaDividida[i])
				cadenaDividida[i] = string(r[:len(r)-1])
			}
		}
	}
	/*
		fmt.Println(red + "-------" + reset)
		for i := 0; i < len(cadenaDividida); i++ {
			fmt.Print(cadenaDividida[i] + " -> " + red)
			fmt.Print(i)
			fmt.Println(reset)
		}
		fmt.Println(red + "-------" + reset)
	*/
	if strings.Contains(cadena, "#") {
		fmt.Println(cyan + "COMENTARIO: " + reset + cadena)
	} else {
		for i := 0; i < len(cadenaDividida); i++ {
			if strings.Contains(cadenaDividida[i], "\"") && encontreComilla == true {
				encontreComilla = false
				cadenaDividida[i] = strings.ReplaceAll(cadenaDividida[i], "\"", "")
				temporal += cadenaDividida[i]
				cadenaReordenadaDividida = append(cadenaReordenadaDividida, temporal)
				temporal = ""
			} else if strings.Contains(cadenaDividida[i], "\"") && encontreComilla == false {
				cadenaDividida[i] = strings.ReplaceAll(cadenaDividida[i], "\"", "")
				encontreComilla = true
			} else if encontreComilla == false {
				cadenaReordenadaDividida = append(cadenaReordenadaDividida, cadenaDividida[i])
			}
			if encontreComilla == true {
				temporal += cadenaDividida[i] + " "
			}
		}
		//cadenaReordenadaDividida = append(cadenaReordenadaDividida, temporal)
	}
	/*
		fmt.Println(yellow + temporal + reset)

		fmt.Println(green + "empieza la reordenada" + reset)

		for i := 0; i < len(cadenaReordenadaDividida); i++ {
			fmt.Println(cadenaReordenadaDividida[i])
		}
	*/
	return cadenaReordenadaDividida
}

//EjecutarLineas = separa las lineas obtenidas del archivo .mia y las recorre una por una
func EjecutarLineas(cadena string) {
	//fmt.Println(magenta + "********************" + reset)
	//fmt.Println(cadena)
	var cadenaDividida []string = strings.SplitN(cadena, "\n", -1)
	//fmt.Println(red + "Empezamos la separacion: " + reset)
	//fmt.Println()
	for i := 0; i < len(cadenaDividida); i++ {
		if cadenaDividida[i] == "" {
			//fmt.Println(red + "ERROR" + reset)
		} else {
			Verificar_Tipo(Iniciar_Analisis(cadenaDividida[i]))
		}
		//fmt.Print(cadenaDividida[i])
		//fmt.Print(yellow + " -> " + reset)
		//fmt.Println(i)
	}

	//fmt.Println(magenta + "********************" + reset)
}

//Metodo_EXEC = Este metodo obtiene la ruta del archivo y lo abre
func Metodo_EXEC(vector []string) string {
	var parametros [1]string //[0]PATH
	var vecAuxiliar []string = nil

	for i := 1; i < len(vector); i++ {
		vecAuxiliar = strings.SplitN(vector[i], "->", -1)
	}

	/*
		fmt.Println(green + "===================" + reset)
		for i := 0; i < len(vecAuxiliar); i++ {
			fmt.Print(vecAuxiliar[i])
			fmt.Print(yellow + " -> " + reset)
			fmt.Println(i)
		}
		fmt.Println(green + "===================" + reset)
	*/
	for j := 0; j < len(vecAuxiliar); j++ {
		if strings.ToLower(vecAuxiliar[j]) == "-path" {
			parametros[0] = vecAuxiliar[j+1]
		}
	}

	file, error := os.Open(parametros[0])
	if error != nil { //Valor nil quiere decir que no hay error
		fmt.Println(red + "ERROR CRITIO AL INTENTAR ABRIR EL ARCHIVO" + reset)
		log.Fatal(error)
	}
	defer file.Close() //defer es la encargada de asegurar que una funcion es llamada (obliga a funcionar)
	scanner := bufio.NewScanner(file)
	var contenido string
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "\\*") {
			reconstruido := strings.ReplaceAll(scanner.Text(), "\\*", "")
			//fmt.Println(yellow + reconstruido + reset)
			//fmt.Println(red + "SE ENCONTRO UN ERROR" + reset)
			contenido += reconstruido + " "
		} else {
			contenido += scanner.Text() + "\n"
		}
		//fmt.Println(scanner.Text())
	}
	return contenido
}

//SaltoLinea =
func SaltoLinea() {
	for i := 0; i < 3; i++ {
		fmt.Println()
	}
}

//Verificar_Tipo = Funcion encargada de verificar que comando se procede a ejecutar
func Verificar_Tipo(cadenaOrdenada []string) {
	if cadenaOrdenada == nil {
		fmt.Println(red + "VECTOR VACIO" + reset)
	} else {
		var mamarre string = cadenaOrdenada[0]
		//fmt.Println(cyan + mamarre + reset)
		if strings.ToLower(mamarre) == "exec" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "EXEC" + reset)
			var cadenaEXEC string = Metodo_EXEC(cadenaOrdenada)
			//var nuevoOrden []string = Iniciar_Analisis(cadenaEXEC)
			EjecutarLineas(cadenaEXEC)
		} else if strings.ToLower(mamarre) == "pause" {
			SaltoLinea()
			fmt.Println(red + "[PAUSE]" + reset + yellow + "Presione una tecla para ejecutar: " + reset)
			var tecla string
			fmt.Scanln(&tecla)
		} else if strings.ToLower(mamarre) == "mkdisk" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "MKDISK" + reset)
			funciones.FuncionMKDISK(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "rmdisk" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "RMDISK" + reset)
			funciones.FuncionRMDISK(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "fdisk" {
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "FDISK" + reset)
			funciones.FuncionFDISK(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "mount" {
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "MOUNT" + reset)
			funciones.FuncionMOUNT(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "unmount" {
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "UNMOUNT" + reset)
			funciones.FuncionUNMOUNT(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "leer" {
			metodos.ResumenMBR("/home/javier/Imágenes/Disco1.dsk")
		} else if mamarre == "" {
			SaltoLinea()
			fmt.Println(red + "NO SE INGRESO NINGUN PARAMETRO" + reset)
		} else {
			SaltoLinea()
			fmt.Println(red + "ERROR: EL PARAMETRO NO ES VALIDO" + reset)
		}
	}
}
