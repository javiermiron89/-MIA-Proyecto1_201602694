package funciones

import (
	"fmt"
	"strconv"
	"strings"

	"../metodos"
)

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
//FUNCIONES-----FUNCIONES-----FUNCIONES-----FUNCIONES-----FUNCIONES-----FUNCIONES-----FUNCIONES-----FUNCIONES-----FUNCIONES-----FUNCIONES----
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//InitApp =
func InitApp() {

}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKDISK-----MKDISK-----FUNCIONES-----FUNCIONES-----MKDISK-----MKDISK-----FUNCIONES-----FUNCIONES-----MKDISK-----MKDISK-----FUNCIONES--------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionMKDISK =
func FuncionMKDISK(vector []string) {
	var parametros [4]string //[0]SIZE   [1]PATH   [2]NAME   [3]UNIT
	var vecAuxiliar []string = nil
	sizeObligatorio := false
	pathObligatorio := false
	nameObligatorio := false
	unitOpcional := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-size" {
			//fmt.Println("SOY SIZE" + "    " + vecAuxiliar[1])
			tam, _ := strconv.Atoi(vecAuxiliar[1])
			if tam > 0 {
				parametros[0] = vecAuxiliar[1]
				sizeObligatorio = true
			} else {
				fmt.Println(red + "[ERROR]" + reset + "El valor de size no es mayor a cero")
			}
		} else if strings.ToLower(vecAuxiliar[0]) == "-path" {
			//fmt.Println("SOY PATH" + "    " + vecAuxiliar[1])
			parametros[1] = vecAuxiliar[1]
			pathObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-name" {
			//fmt.Println("SOY NAME" + "    " + vecAuxiliar[1])
			parametros[2] = vecAuxiliar[1]
			nameObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-unit" {
			//fmt.Println("SOY UNIT" + "    " + vecAuxiliar[1])
			if strings.ToLower(vecAuxiliar[1]) == "k" {
				parametros[3] = "k"
				unitOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "m" {
				parametros[3] = "m"
				unitOpcional = true
			}
		}
	}

	if sizeObligatorio == true && pathObligatorio == true && nameObligatorio == true {
		if unitOpcional == true {
			//ya fue ingresado
		} else {
			parametros[3] = "m"
		}
		metodos.CrearDisco(parametros[0], parametros[1], parametros[2], parametros[3])
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "MKDISK" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//RMDISK-----RMDISK-----FUNCIONES-----FUNCIONES-----RMDISK-----RMDISK-----FUNCIONES-----FUNCIONES-----RMDISK-----RMDISK-----FUNCIONES--------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionRMDISK =
func FuncionRMDISK(vector []string) {
	var parametros [1]string //[0]PATH
	var vecAuxiliar []string = nil
	pathObligatorio := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-path" {
			parametros[0] = vecAuxiliar[1]
			pathObligatorio = true
		}
	}

	if pathObligatorio == true {
		metodos.EliminarDisco(parametros[0])
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "RMDISK" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//FDISK-----FDISK-----FUNCIONES-----FUNCIONES-----FDISK-----FDISK-----FUNCIONES-----FUNCIONES-----FDISK-----FDISK-----FUNCIONES--------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionFDISK =
func FuncionFDISK(vector []string) {
	var parametros [8]string //[0]SIZE   [1]UNIT   [2]PATH   [3]TYPE	[4]FIT	[5]DELETE	[6]NAME		[7]ADD
	var vecAuxiliar []string = nil
	sizeObligatorio := false
	pathObligatorio := false
	nameObligatorio := false
	unitOpcional := false
	typeOpcional := false
	fitOpcional := false
	deleteOpcional := false
	addOpcional := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-size" {
			//fmt.Println("SOY SIZE" + "    " + vecAuxiliar[1])
			tam, _ := strconv.Atoi(vecAuxiliar[1])
			if tam > 0 {
				parametros[0] = vecAuxiliar[1]
				sizeObligatorio = true
			} else {
				fmt.Println(red + "[ERROR]" + reset + "El valor de size no es mayor a cero")
				break
			}
		} else if strings.ToLower(vecAuxiliar[0]) == "-unit" {
			//fmt.Println("SOY UNIT" + "    " + vecAuxiliar[1])
			if strings.ToLower(vecAuxiliar[1]) == "b" {
				parametros[1] = "B"
				unitOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "k" {
				parametros[1] = "K"
				unitOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "m" {
				parametros[1] = "M"
				unitOpcional = true
			} else {
				fmt.Println(red + "[ERROR]" + reset + "El parametro ingresado no es una opcion valida")
			}
		} else if strings.ToLower(vecAuxiliar[0]) == "-path" {
			//fmt.Println("SOY PATH" + "    " + vecAuxiliar[1])
			if vecAuxiliar[1][len(vecAuxiliar[1])-4:] == ".dsk" {
				if metodos.ExisteCarpeta(vecAuxiliar[1]) {
					parametros[2] = vecAuxiliar[1]
					pathObligatorio = true
				} else {
					fmt.Println(red + "[ERROR]" + reset + "La carpeta que ha especificado no existe")
					break
				}
			} else {
				fmt.Println(red + "[ERROR]" + reset + "La extension en el parametro main no coincide con " + cyan + "[.dsk]" + reset)
				break
			}
		} else if strings.ToLower(vecAuxiliar[0]) == "-type" {
			if strings.ToLower(vecAuxiliar[1]) == "p" {
				parametros[3] = "P"
				typeOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "e" {
				parametros[3] = "E"
				typeOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "l" {
				parametros[3] = "L"
				typeOpcional = true
			} else {
				fmt.Println(red + "[ERROR]" + reset + "El parametro ingresado no es una opcion valida")
			}
		} else if strings.ToLower(vecAuxiliar[0]) == "-fit" {
			if strings.ToLower(vecAuxiliar[1]) == "bf" {
				parametros[4] = "B"
				fitOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "ff" {
				parametros[4] = "F"
				fitOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "wf" {
				parametros[4] = "W"
				fitOpcional = true
			} else {
				fmt.Println(red + "[ERROR]" + reset + "El parametro ingresado no es una opcion valida")
			}
		} else if strings.ToLower(vecAuxiliar[0]) == "-delete" {
			if strings.ToLower(vecAuxiliar[1]) == "fast" {
				parametros[5] = "FAST"
				deleteOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "full" {
				parametros[5] = "FULL"
				deleteOpcional = true
			} else {
				fmt.Println(red + "[ERROR]" + reset + "El parametro ingresado no es una opcion valida")
			}
		} else if strings.ToLower(vecAuxiliar[0]) == "-name" {
			//fmt.Println("SOY NAME" + "    " + vecAuxiliar[1])
			parametros[6] = vecAuxiliar[1]
			nameObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-add" {
			addOpcional = true
			parametros[7] = vecAuxiliar[1]
		}
	}

	//[0]SIZE   [1]UNIT   [2]PATH   [3]TYPE	[4]FIT	[5]DELETE	[6]NAME		[7]ADD

	if sizeObligatorio == true && pathObligatorio == true && nameObligatorio == true {
		if unitOpcional == false {
			parametros[1] = "K"
		}
		if typeOpcional == false {
			parametros[3] = "P"
		}
		if fitOpcional == false {
			parametros[4] = "WF"
		}
		if deleteOpcional == true && addOpcional == true {
			fmt.Println(red + "[ERROR]" + reset + "Los parametros " + cyan + "[-delete]" + reset + " y " + cyan + "[-add]" + reset + " no son compatibles")
		} else if deleteOpcional == true && addOpcional == false {
			metodos.EliminarParticion()
		} else if deleteOpcional == false && addOpcional == true {
			metodos.ModificarParticion()
		} else if deleteOpcional == false && addOpcional == false {
			//size, unit, path, type2, fit, name string
			metodos.CrearParticion(parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[6])
			fmt.Println(green + "[EXITO]" + reset + "La particion ha sido creada con exito")
		}

	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "FDISK" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MOUNT-----MOUNT-----FUNCIONES-----FUNCIONES-----MOUNT-----MOUNT-----FUNCIONES-----FUNCIONES-----MOUNT-----MOUNT-----FUNCIONES--------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionMOUNT =
func FuncionMOUNT(vector []string) {
	var parametros [2]string //[0]PATH   [1]NAME
	var vecAuxiliar []string = nil
	mountVieneSolo := true
	pathObligatorio := false
	nameObligatorio := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-path" {
			parametros[0] = vecAuxiliar[1]
			pathObligatorio = true
			mountVieneSolo = false
		} else if strings.ToLower(vecAuxiliar[0]) == "-name" {
			parametros[1] = vecAuxiliar[1]
			nameObligatorio = true
			mountVieneSolo = false
		}
	}

	if mountVieneSolo == true {
		metodos.ResumenParticionesMontadas()
	} else {
		if pathObligatorio == true && nameObligatorio == true {
			metodos.MontarParticion(parametros[0], parametros[1])
			fmt.Println(green + "[EXITO]" + reset + "La particion " + cyan + parametros[1] + reset + " fue montado con exito")
		} else {
			fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "MOUNT" + reset + " obligatorios no han sido completamente ingresados")
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MOUNT-----MOUNT-----FUNCIONES-----FUNCIONES-----MOUNT-----MOUNT-----FUNCIONES-----FUNCIONES-----MOUNT-----MOUNT-----FUNCIONES--------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionUNMOUNT =
func FuncionUNMOUNT(vector []string) {
	var parametros [2]string //[0]PATH   [1]NAME
	var listadoDesmontar []string
	var vecAuxiliar []string = nil
	idObligatorio := false

	contador := 0
	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		vecAuxiliar[0] = vecAuxiliar[0][:len(vecAuxiliar[0])-1]
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			listadoDesmontar = append(listadoDesmontar, parametros[0])
			idObligatorio = true
			contador++
		}
	}

	if idObligatorio == true {
		for i := 0; i < len(listadoDesmontar); i++ {
			metodos.DesmontarParticion(listadoDesmontar[i])
		}
	} else {

	}
}
