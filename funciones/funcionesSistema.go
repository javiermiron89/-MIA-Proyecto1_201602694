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
	//fmt.Println(vector)
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
					fmt.Println(red + "[ERROR]" + reset + "El disco que ha especificado no existe")
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
	if deleteOpcional == true {
		sizeObligatorio = true
	}
	if addOpcional == true {
		sizeObligatorio = true
	}

	if sizeObligatorio == true && pathObligatorio == true && nameObligatorio == true {
		if unitOpcional == false {
			parametros[1] = "K"
		}
		if typeOpcional == false {
			parametros[3] = "P"
		}
		if fitOpcional == false {
			parametros[4] = "W"
		}
		if deleteOpcional == true && addOpcional == true {
			fmt.Println(red + "[ERROR]" + reset + "Los parametros " + cyan + "[-delete]" + reset + " y " + cyan + "[-add]" + reset + " no son compatibles")
		} else if deleteOpcional == true && addOpcional == false {
			//fmt.Println("VOY A BORRAR UNA PARTICION")
			metodos.EliminarParticion(parametros[2], parametros[6], parametros[5])
		} else if deleteOpcional == false && addOpcional == true {
			metodos.ModificarParticion()
		} else if deleteOpcional == false && addOpcional == false {
			//size, unit, path, type2, fit, name string
			if parametros[3] == "L" {
				metodos.InsertarParticionLogica(parametros[2], parametros[6], parametros[0], parametros[4])
				//fmt.Println(red + "****************" + reset)
				//metodos.ResumenEBR(parametros[2], parametros[6])
			} else {
				metodos.CrearParticion(parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[6])
			}
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
		} else {
			fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "MOUNT" + reset + " obligatorios no han sido completamente ingresados")
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//UNMOUNT-----UNMOUNT-----FUNCIONES-----FUNCIONES-----UNMOUNT-----UNMOUNT-----FUNCIONES-----FUNCIONES-----UNMOUNT-----UNMOUNT-----FUNCIONES--
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
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "UNMOUNT" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKFS-----MKFS-----FUNCIONES-----FUNCIONES-----MKFS-----MKFS-----FUNCIONES-----FUNCIONES-----MKFS-----MKFS-----FUNCIONES-----FUNCIONES------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionMKFS =
func FuncionMKFS(vector []string) {
	var parametros [4]string //[0]ID   [1]TIPO	[2]ADD	[3]UNIT
	var vecAuxiliar []string = nil
	idObligatorio := false
	tipoOpcional := false
	addOpcional := false
	unitOpcional := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-tipo" {
			if strings.ToLower(vecAuxiliar[1]) == "fast" {
				parametros[1] = "FAST"
				tipoOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "full" {
				parametros[1] = "FULL"
				tipoOpcional = true
			}
		} else if strings.ToLower(vecAuxiliar[0]) == "-add" {
			parametros[2] = vecAuxiliar[1]
			addOpcional = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-unit" {
			if strings.ToLower(vecAuxiliar[1]) == "b" {
				parametros[3] = "B"
				tipoOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "k" {
				parametros[3] = "K"
				tipoOpcional = true
			} else if strings.ToLower(vecAuxiliar[1]) == "m" {
				parametros[3] = "M"
				tipoOpcional = true
			}

		}
	}
	//[0]ID   [1]TIPO	[2]ADD	[3]UNIT
	if idObligatorio == true {
		if tipoOpcional == false {
			parametros[1] = "FULL"
		}
		if addOpcional == false {

		}
		if unitOpcional == false {
			parametros[3] = "K"
		}
		metodos.FormateLWH(parametros[0], parametros[1])
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "MKFS" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//LOGIN-----LOGIN-----FUNCIONES-----FUNCIONES-----LOGIN-----LOGIN-----FUNCIONES-----FUNCIONES-----LOGIN-----LOGIN-----FUNCIONES-----FUNCIONES
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionLOGIN =
func FuncionLOGIN(vector []string) {
	var parametros [3]string //[0]usr   [1]pwd	[2]id
	var vecAuxiliar []string = nil
	usrObligatorio := false
	pwdObligatorio := false
	idObligatorio := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-usr" {
			parametros[0] = vecAuxiliar[1]
			usrObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-pwd" {
			parametros[1] = vecAuxiliar[1]
			pwdObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[2] = vecAuxiliar[1]
			idObligatorio = true
		}
	}
	if usrObligatorio == true && pwdObligatorio == true && idObligatorio == true {
		metodos.Login(parametros[0], parametros[1], parametros[2])
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "LOGIN" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKGRP-----MKGRP-----FUNCIONES-----FUNCIONES-----MKGRP-----MKGRP-----FUNCIONES-----FUNCIONES-----MKGRP-----MKGRP-----FUNCIONES-----FUNCIONES
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionMKGRP =
func FuncionMKGRP(vector []string) {
	var parametros [2]string //[0]id   [1]name
	var vecAuxiliar []string = nil
	idObligatorio := false
	nameObligatorio := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-name" {
			parametros[1] = vecAuxiliar[1]
			nameObligatorio = true
		}
	}

	if idObligatorio == true && nameObligatorio == true {
		metodos.CrearGrupo(parametros[0], parametros[1])
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "MKGRP" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//RMGRP-----RMGRP-----FUNCIONES-----FUNCIONES-----RMGRP-----RMGRP-----FUNCIONES-----FUNCIONES-----RMGRP-----RMGRP-----FUNCIONES-----FUNCIONES
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionRMGRP =
func FuncionRMGRP(vector []string) {
	var parametros [2]string //[0]id   [1]name
	var vecAuxiliar []string = nil
	idObligatorio := false
	nameObligatorio := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-name" {
			parametros[1] = vecAuxiliar[1]
			nameObligatorio = true
		}
	}

	if idObligatorio == true && nameObligatorio == true {
		metodos.RemoverGrupo(parametros[0], parametros[1])
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "RMGRP" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKUSR-----MKUSR-----FUNCIONES-----FUNCIONES-----MKUSR-----MKUSR-----FUNCIONES-----FUNCIONES-----MKUSR-----MKUSR-----FUNCIONES-----FUNCIONES
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionMKUSR =
func FuncionMKUSR(vector []string) {
	var parametros [4]string //[0]id   [1]usr	[2]pwd	[3]grp
	var vecAuxiliar []string = nil
	idObligatorio := false
	usrObligatorio := false
	pwdObligatorio := false
	grpObligatorio := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-usr" {
			parametros[1] = vecAuxiliar[1]
			usrObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-pwd" {
			parametros[2] = vecAuxiliar[1]
			pwdObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-grp" {
			parametros[3] = vecAuxiliar[1]
			grpObligatorio = true
		}
	}

	if idObligatorio == true && usrObligatorio == true && pwdObligatorio == true && grpObligatorio == true {
		metodos.CrearUser(parametros[0], parametros[1], parametros[2], parametros[3])
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//RMUSR-----RMUSR-----FUNCIONES-----FUNCIONES-----RMUSR-----RMUSR-----FUNCIONES-----FUNCIONES-----RMUSR-----RMUSR-----FUNCIONES-----FUNCIONES
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionRMUSR =
func FuncionRMUSR(vector []string) {
	var parametros [2]string //[0]id   [1]usr
	var vecAuxiliar []string = nil
	idObligatorio := false
	usrObligatorio := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-usr" {
			parametros[1] = vecAuxiliar[1]
			usrObligatorio = true
		}
	}

	if idObligatorio == true && usrObligatorio == true {
		metodos.RemoverUser(parametros[0], parametros[1])
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "RMUSR" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKFILE-----MKFILE-----FUNCIONES-----FUNCIONES-----MKFILE-----MKFILE-----FUNCIONES-----FUNCIONES-----MKFILE-----MKFILE-----FUNCIONES-----FUN
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionMKFILE =
func FuncionMKFILE(vector []string) {
	var parametros [5]string //[0]id   [1]path	[2]-p	[3]size	[4]cont
	var vecAuxiliar []string = nil
	idObligatorio := false
	pathObligatorio := false
	pOpcional := false
	sizeOpcional := false
	contOpcional := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		//fmt.Println(vecAuxiliar)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-path" {
			parametros[1] = vecAuxiliar[1]
			pathObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-p" {
			pOpcional = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-size" {
			parametros[3] = vecAuxiliar[1]
			sizeOpcional = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-cont" {
			parametros[4] = vecAuxiliar[1]
			contOpcional = true
		}
	}

	if idObligatorio == true && pathObligatorio == true {
		if sizeOpcional == true {
			tamEnBytesArchivo, _ := strconv.ParseInt(parametros[3], 10, 64)
			if tamEnBytesArchivo < 0 {
				fmt.Println(red + "[ERROR]" + reset + "El parametro size tiene un valor negativo, por lo cual se procede a colocarse de tamaño 0 por defecto")
				parametros[3] = "0"
			}
		} else {
			parametros[3] = ""
		}
		if contOpcional == false {
			parametros[4] = ""
		}
		if pOpcional == true {
			metodos.CrearArchivo(parametros[0], parametros[1], true, parametros[3], parametros[4])
		} else {
			metodos.CrearArchivo(parametros[0], parametros[1], false, parametros[3], parametros[4])
		}
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "MKDIR" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//EDIT-----EDIT-----FUNCIONES-----FUNCIONES-----EDIT-----EDIT-----FUNCIONES-----FUNCIONES-----EDIT-----EDIT-----FUNCIONES-----FUNCIONES------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionEDIT =
func FuncionEDIT(vector []string) {
	var parametros [4]string //[0]id   [1]path	[2]size	[3]cont
	var vecAuxiliar []string = nil
	idObligatorio := false
	pathObligatorio := false
	sizeOpcional := false
	contOpcional := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		fmt.Println(vecAuxiliar)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-path" {
			parametros[1] = vecAuxiliar[1]
			pathObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-size" {
			parametros[2] = vecAuxiliar[1]
			sizeOpcional = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-cont" {
			parametros[3] = vecAuxiliar[1]
			contOpcional = true
		}
	}

	if idObligatorio == true && pathObligatorio == true {
		if sizeOpcional == false && contOpcional == false {
			fmt.Println(red + "[ERROR]" + reset + "Se debe ingresar al menos uno de estos parametros (size <-> cont)")
		} else {
			metodos.ModificarArchivo(parametros[0], parametros[1], parametros[2], parametros[3])
		}
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "EDIT" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//CAT-----CAT-----FUNCIONES-----FUNCIONES-----CAT-----CAT-----FUNCIONES-----FUNCIONES-----CAT-----CAT-----FUNCIONES-----FUNCIONES-----CAT----
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionCAT =
func FuncionCAT(vector []string) {
	var parametros [2]string //[0]PATH   [1]NAME
	var listadoArchivos []string
	var vecAuxiliar []string = nil
	idObligatorio := false
	fileObligatorio := false

	contador := 0
	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else {
			vecAuxiliar[0] = vecAuxiliar[0][:len(vecAuxiliar[0])-1]
			if strings.ToLower(vecAuxiliar[0]) == "-file" {
				parametros[1] = vecAuxiliar[1]
				listadoArchivos = append(listadoArchivos, parametros[1])
				fileObligatorio = true
				contador++
			}
		}
	}

	if idObligatorio == true && fileObligatorio == true {
		for i := 0; i < len(listadoArchivos); i++ {
			fmt.Println(magenta+"Contenido del archivo:"+cyan, listadoArchivos[i], reset)
			metodos.ImprimirContenidoArchivo(parametros[0], listadoArchivos[i])
		}
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "CAT" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKDIR-----MKDIR-----FUNCIONES-----FUNCIONES-----MKDIR-----MKDIR-----FUNCIONES-----FUNCIONES-----MKDIR-----MKDIR-----FUNCIONES-----FUNCIONES
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionMKDIR =
func FuncionMKDIR(vector []string) {
	var parametros [3]string //[0]id   [1]path	[2]-p
	var vecAuxiliar []string = nil
	idObligatorio := false
	pathObligatorio := false
	pOpcional := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		fmt.Println(vecAuxiliar)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-path" {
			parametros[1] = vecAuxiliar[1]
			pathObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-p" {
			pOpcional = true
		}
	}

	if idObligatorio == true && pathObligatorio == true {
		if pOpcional == true {
			metodos.CrearDirectorio(parametros[0], parametros[1], true)
		} else {
			metodos.CrearDirectorio(parametros[0], parametros[1], false)
		}
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "MKDIR" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//CHMOD-----CHMOD-----FUNCIONES-----FUNCIONES-----CHMOD-----CHMOD-----FUNCIONES-----FUNCIONES-----CHMOD-----CHMOD-----FUNCIONES-----FUNCIONES
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionCHMOD =
func FuncionCHMOD(vector []string) {
	var parametros [4]string //[0]id   [1]path	[2]ugo		[3]r
	var vecAuxiliar []string = nil
	idObligatorio := false
	pathObligatorio := false
	ugoObligatorio := false
	rOpcional := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		fmt.Println(vecAuxiliar)
		if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[0] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-path" {
			parametros[1] = vecAuxiliar[1]
			pathObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-ugo" {
			parametros[2] = vecAuxiliar[1]
			ugoObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-r" {
			rOpcional = true
		}
	}

	if idObligatorio == true && pathObligatorio == true && ugoObligatorio == true {
		if rOpcional == true {
			metodos.MetodoCHMOD(parametros[0], parametros[1], parametros[2], true)
		} else {
			metodos.MetodoCHMOD(parametros[0], parametros[1], parametros[2], false)
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//REP-----REP-----FUNCIONES-----FUNCIONES-----REP-----REP-----FUNCIONES-----FUNCIONES-----REP-----REP-----FUNCIONES-----FUNCIONES-----REP----
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionREP =
func FuncionREP(vector []string) {
	var parametros [4]string //[0]NOMBRE   [1]PATH	[2]ID	[3]RUTA
	var vecAuxiliar []string = nil
	nombreObligatorio := false
	pathObligatorio := false
	idObligatorio := false
	rutaOpcional := false

	for j := 1; j < len(vector); j++ {
		vecAuxiliar = strings.SplitN(vector[j], "->", -1)
		if strings.ToLower(vecAuxiliar[0]) == "-name" {
			if strings.ToLower(vecAuxiliar[1]) == "mbr" {
				parametros[0] = "MBR"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "disk" {
				parametros[0] = "DISK"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "sb" {
				parametros[0] = "SB"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "bm_arbdir" {
				parametros[0] = "BM_ARBDIR"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "bm_detdir" {
				parametros[0] = "BM_DETDIR"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "bm_inode" {
				parametros[0] = "BM_INODE"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "bm_block" {
				parametros[0] = "BM_BLOCK"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "bitacora" {
				parametros[0] = "BITACORA"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "directorio" {
				parametros[0] = "DIRECTORIO"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "tree_file" {
				parametros[0] = "TREE_FILE"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "tree_directorio" {
				parametros[0] = "TREE_DIRECTORIO"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "tree_complete" {
				parametros[0] = "TREE_COMPLETE"
				nombreObligatorio = true
			} else if strings.ToLower(vecAuxiliar[1]) == "ls" {
				parametros[0] = "LS"
				nombreObligatorio = true
			} else {
				fmt.Println(red + "[ERROR]" + reset + "El parametro ingresado no es una opcion valida")
			}
		} else if strings.ToLower(vecAuxiliar[0]) == "-path" {
			parametros[1] = vecAuxiliar[1]
			pathObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-id" {
			parametros[2] = vecAuxiliar[1]
			idObligatorio = true
		} else if strings.ToLower(vecAuxiliar[0]) == "-ruta" {
			parametros[3] = vecAuxiliar[1]
			rutaOpcional = true
		}

	}

	//[0]NOMBRE   [1]PATH	[2]ID	[3]RUTA
	if nombreObligatorio == true && pathObligatorio == true && idObligatorio == true {
		if parametros[0] == "MBR" {
			metodos.ReporteMBR(parametros[2], parametros[1])
		} else if parametros[0] == "DISK" {
			metodos.ReporteDISK(parametros[2], parametros[1])
		} else if parametros[0] == "SB" {
			metodos.ReporteSB(parametros[2], parametros[1])
		} else if parametros[0] == "BM_ARBDIR" {
			metodos.ReporteBitmap(parametros[2], parametros[1], 1)
		} else if parametros[0] == "BM_DETDIR" {
			metodos.ReporteBitmap(parametros[2], parametros[1], 2)
		} else if parametros[0] == "BM_INODE" {
			metodos.ReporteBitmap(parametros[2], parametros[1], 3)
		} else if parametros[0] == "BM_BLOCK" {
			metodos.ReporteBitmap(parametros[2], parametros[1], 4)
		} else if parametros[0] == "DIRECTORIO" {
			metodos.ReporteDirectorio(parametros[2], parametros[1])
		} else if parametros[0] == "TREE_FILE" {
			if rutaOpcional == true {
				metodos.ReporteTreeFile(parametros[2], parametros[1], parametros[3])
			} else {
				fmt.Println(red + "[ERROR]" + reset + "El parametro " + magenta + "ruta" + reset + " no fue ingresado para poder realizar el reporte Tree_File")
			}
		} else if parametros[0] == "TREE_COMPLETE" {
			metodos.ReporteTreeComplete(parametros[2], parametros[1])
		} else if parametros[0] == "TREE_DIRECTORIO" {
			metodos.ReporteTreeDirectorio(parametros[2], parametros[1])
		} else if parametros[0] == "LS" {
			//metodos.ReporteLS(parametros[2], parametros[1])
		}
	} else {
		fmt.Println(red + "[ERROR]" + reset + "Los parametros de " + magenta + "REP" + reset + " obligatorios no han sido completamente ingresados")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//PRUEBITA-----REP-----FUNCIONES-----FUNCIONES-----REP-----REP-----FUNCIONES-----FUNCIONES-----REP-----REP-----FUNCIONES-----FUNCIONES-----REP----
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionPRUEBITA =
func FuncionPRUEBITA() {

}
