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

//IniciarAnalisis =
func IniciarAnalisis(cadena string) []string {
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
			//fmt.Println("Mamarre: ", cadenaDividida[i][len(cadenaDividida[i])-1:])
			if (cadenaDividida[i][len(cadenaDividida[i])-1:]) == "\n" {
				//fmt.Println("uts")
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
		SaltoLinea()
		fmt.Println(cyan + "COMENTARIO: " + reset + cadena)
	} else {
		for i := 0; i < len(cadenaDividida); i++ {
			contadorComillas := strings.Count(cadenaDividida[i], "\"")
			if contadorComillas == 2 {
				cadenaDividida[i] = strings.ReplaceAll(cadenaDividida[i], "\"", "")
				temporal += cadenaDividida[i]
				cadenaReordenadaDividida = append(cadenaReordenadaDividida, temporal)
			} else if strings.Contains(cadenaDividida[i], "\"") && encontreComilla == true {
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
			VerificarTipo(IniciarAnalisis(cadenaDividida[i]))
		}
		//fmt.Print(cadenaDividida[i])
		//fmt.Print(yellow + " -> " + reset)
		//fmt.Println(i)
	}

	//fmt.Println(magenta + "********************" + reset)
}

//MetodoEXEC = Este metodo obtiene la ruta del archivo y lo abre
func MetodoEXEC(vector []string) string {
	var parametros [1]string //[0]PATH
	var vecAuxiliar []string = nil

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-path" {
			parametros[0] = vecAuxiliar[1]
		}
	}

	fmt.Println("Ruta: " + parametros[0])
	//file, error := os.OpenFile(parametros[0], os.O_RDWR, 0755)
	file, error := os.Open(parametros[0])
	if error != nil { //Valor nil quiere decir que no hay error
		fmt.Println(red + "ERROR CRITICO AL INTENTAR ABRIR EL ARCHIVO" + reset)
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
			contenido += reconstruido + ""
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

//VerificarTipo = Funcion encargada de verificar que comando se procede a ejecutar
func VerificarTipo(cadenaOrdenada []string) {
	if cadenaOrdenada == nil {
		//fmt.Println(red + "VECTOR VACIO" + reset)
	} else {
		var mamarre string = cadenaOrdenada[0]
		//fmt.Println(cyan + mamarre + reset)
		if strings.ToLower(mamarre) == "exec" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "EXEC" + reset)
			var cadenaEXEC string = MetodoEXEC(cadenaOrdenada)
			//fmt.Println(cadenaEXEC)
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
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "FDISK" + reset)
			funciones.FuncionFDISK(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "mount" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "MOUNT" + reset)
			funciones.FuncionMOUNT(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "unmount" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "UNMOUNT" + reset)
			funciones.FuncionUNMOUNT(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "mkfs" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "MKFS" + reset)
			funciones.FuncionMKFS(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "login" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "LOGIN" + reset)
			funciones.FuncionLOGIN(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "logout" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "LOGOUT" + reset)
			metodos.FuncionLOGOUT()
		} else if strings.ToLower(mamarre) == "mkgrp" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "MKGRP" + reset)
			funciones.FuncionMKGRP(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "rmgrp" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "RMGRP" + reset)
			funciones.FuncionRMGRP(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "mkusr" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "MKUSR" + reset)
			funciones.FuncionMKUSR(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "rmusr" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "RMUSR" + reset)
			funciones.FuncionRMUSR(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "mkfile" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "MKFILE" + reset)
			funciones.FuncionMKFILE(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "mkdir" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "MKDIR" + reset)
			funciones.FuncionMKDIR(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "rep" {
			SaltoLinea()
			fmt.Println(yellow + "Palabra reservada a ejecutar: " + magenta + "REP" + reset)
			funciones.FuncionREP(cadenaOrdenada)
		} else if strings.ToLower(mamarre) == "leer" {
			SaltoLinea()
			metodos.ResumenMBR("/home/javier/Imágenes/Disco1.dsk")
		} else if strings.ToLower(mamarre) == "pruebita" {
			funciones.FuncionPRUEBITA()
		} else if mamarre == "" {
			SaltoLinea()
			fmt.Println(red + "NO SE INGRESO NINGUN PARAMETRO" + reset)
		} else {
			SaltoLinea()
			fmt.Println(red + "[ERROR]" + reset + "EL PARAMETRO " + cyan + mamarre + reset + " NO ES VALIDO" + reset)
		}
	}
}
