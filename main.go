package main

import (
	"bufio"
	"fmt"
	"os"

	"./analizador"
)

var palabras []string

//Colores
var red = "\033[1;31m"
var green = "\033[1;32m"
var yellow = "\033[1;33m"
var blue = "\033[1;34m"
var magenta = "\033[1;35m"
var cyan = "\033[1;36m"
var white = "\033[1;37m"
var reset = "\u001B[0m"

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN-----MAIN
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

func main() {
	Iniciar()
	var respuesta string

	for {
		SaltoLinea()
		fmt.Println(green + "DESEA REALIZAR OTRA OPERACION?[y/n]" + reset)
		fmt.Scanf("%s", &respuesta)
		if respuesta == "y" {
			Iniciar()
		} else if respuesta == "n" {
			break
		} else {
			fmt.Println(red + "-" + respuesta + "-" + " no es una opcion valida, ingresa alguna de las opciones indicadas" + reset)
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//METODOS-----METODOS-----METODOS-----METODOS-----METODOS-----METODOS-----METODOS-----METODOS-----METODOS-----METODOS-----METODOS-----METODOS
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//Iniciar = sirve para empezar el proceso del analisis lexico
func Iniciar() {
	//lectorLinea es un Reader (vector binario) que almacena en formato ASCII cada uno de los caracteres
	lectorLinea := bufio.NewReader(os.Stdin)
	fmt.Println(green + "Ingrese un comando:" + reset)
	//entrada se encarga de convertir
	entrada, _ := lectorLinea.ReadString('\n') //_ sirve para no darle uso al segundo parametro

	var mamarre string = entrada

	//fmt.Println("Cadena Ingresada :=" + mamarre)

	palabras = nil
	palabras = analizador.Iniciar_Analisis(mamarre)
	analizador.Verificar_Tipo(palabras)

	//RecorrerVectorPalabras()
	//fmt.Println(analizador.Global)
}

//RecorrerVectorPalabras =
func RecorrerVectorPalabras() {
	fmt.Println(red + "################" + reset)
	for i := 0; i < len(palabras); i++ {
		if palabras[i] != "" {
			fmt.Println(palabras[i])
		}
	}
	fmt.Println(red + "################" + reset)
}

//SaltoLinea =
func SaltoLinea() {
	for i := 0; i < 3; i++ {
		fmt.Println()
	}
}
