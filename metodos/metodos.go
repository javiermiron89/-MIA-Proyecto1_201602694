package metodos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var red = "\033[1;31m"
var green = "\033[1;32m"
var yellow = "\033[1;33m"
var blue = "\033[1;34m"
var magenta = "\033[1;35m"
var cyan = "\033[1;36m"
var white = "\033[1;37m"
var reset = "\u001B[0m"

//ContenedorMount es la lista de todos los parametros mount
var ContenedorMount []NodoMontaje

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//RAPIDOS-----RAPIDOS-----RAPIDOS-----RAPIDOS-----RAPIDOS-----RAPIDOS-----RAPIDOS-----RAPIDOS-----RAPIDOS-----RAPIDOS-----RAPIDOS-----STRUCTS
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//InitApp =
func InitApp() {

}

//ExisteCarpeta =
func ExisteCarpeta(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//CrearCarpeta =
func CrearCarpeta(path string) bool {
	os.Mkdir(path, 0755)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err != nil {
			//panic(err)
			return false
		}
		return true
	}
	return true
}

func leerMBR(path string) MBR {
	var mbr MBR
	file, err := os.OpenFile(path, os.O_RDWR, 0755)

	defer file.Close()
	if err != nil {
		fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
	}

	var tamanoEnBytes int = int(unsafe.Sizeof(mbr))

	file.Seek(0, 0)
	data := leerBytes(file, tamanoEnBytes)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println(err)
	}

	/*
		fmt.Println(magenta + "*************" + reset)
		fmt.Println(mbr.MbrTamano)
		//fmt.Println(mbr.MbrFechaCreacion)
		mamarre := string(mbr.MbrFechaCreacion[:])
		fmt.Println(mamarre)
		fmt.Println(mbr.MbrDiskSignature)
		fmt.Println(mbr.MbrPartition1.PartStatus)
		fmt.Println(mbr.MbrPartition2.PartStatus)
		fmt.Println(mbr.MbrPartition3.PartStatus)
		fmt.Println(mbr.MbrPartition4.PartStatus)
		fmt.Println(magenta + "*************" + reset)
	*/

	return mbr
}

func leerEBR(path string, start int64) EBR {
	var ebr EBR
	file, err := os.OpenFile(path, os.O_RDWR, 0755)

	defer file.Close()
	if err != nil {
		fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
	}

	var tamanoEnBytes int = int(binary.Size(ebr))

	file.Seek(start, 0)
	data := leerBytes(file, tamanoEnBytes)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &ebr)
	if err != nil {
		fmt.Println(err)
	}
	/*
		fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
		fmt.Println(cyan+"Status: "+reset, ebr.PartStatus)
		fitP1 := string(ebr.PartFit)
		fmt.Println(cyan+"Fit: "+reset, fitP1)
		fmt.Println(cyan+"Start: "+reset, ebr.PartStart)
		fmt.Println(cyan+"Size: "+reset, ebr.PartSize)
		fmt.Println(cyan+"Next: "+reset, ebr.PartNext)
		mamarreP1 := string(ebr.PartName[:])
		fmt.Println(cyan + "Name: " + reset + mamarreP1)
		fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	*/
	return ebr
}

func escribirBytes(file *os.File, bytes []byte) {
	_, error := file.Write(bytes)
	if error != nil {
		fmt.Println(red + " [ERROR DE ESCRITURA]" + reset)
		log.Fatal(error)
	}
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		fmt.Println(red + " [ERROR DE LECTURA]" + reset)
		log.Fatal(err)
	}

	return bytes
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//STRUCTS-----STRUCTS-----STRUCTS-----STRUCTS-----STRUCTS-----STRUCTS-----STRUCTS-----STRUCTS-----STRUCTS-----STRUCTS-----STRUCTS-----STRUCTS
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//NodoMontaje es el strcut que almacena los discos montados en memoria
type NodoMontaje struct {
	ID              string
	Path            string
	NombreParticion string
	Letra           string
	Numero          int
}

//EBR es el struct de la particion logica
type EBR struct {
	PartStatus byte
	PartFit    byte
	PartStart  int64
	PartSize   int64
	PartNext   int64
	PartName   [16]byte
}

//Partition es el struct que contiene los datos de las particiones
type Partition struct {
	PartStatus byte
	PartType   byte
	PartFit    byte
	PartStart  int64
	PartSize   int64
	PartName   [16]byte
}

//MBR es el struct del master boot recorder
type MBR struct {
	MbrTamano        int64
	MbrFechaCreacion [20]byte
	MbrDiskSignature int64
	MbrPartition1    Partition
	MbrPartition2    Partition
	MbrPartition3    Partition
	MbrPartition4    Partition
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

//CrearDisco =
func CrearDisco(size, path, name, unit string) {
	var sizeFinal int64 = 0

	fmt.Println(cyan + "size:	" + reset + size)
	fmt.Println(cyan + "path:	" + reset + path)
	fmt.Println(cyan + "name:	" + reset + name)
	fmt.Println(cyan + "unit:	" + reset + unit)

	if path[len(path)-1:] != "/" {
		path = path + "/"
	}

	if unit == "k" {
		sizeFinal, _ = strconv.ParseInt(size, 10, 64)
		//fmt.Println(sizeFinal)
		sizeFinal = sizeFinal * 1024
		//fmt.Println(sizeFinal)
	} else if unit == "m" {
		sizeFinal, _ = strconv.ParseInt(size, 10, 64)
		//fmt.Println(sizeFinal)
		sizeFinal = sizeFinal * 1024 * 1024
		//fmt.Println(sizeFinal)
	}

	//Aca se verifica la existencia de la carpeta
	if ExisteCarpeta(path) == true {
		if name[len(name)-4:] == ".dsk" {
			CrearArchivoBinario(path+name, sizeFinal)
			InsertarMBR(sizeFinal, path+name)
			//leerMBR(path + name)
			fmt.Println(green + "[EXITO]" + reset + "El archivo " + magenta + name + reset + " se ha creado con exito")
		} else {
			fmt.Println(red + "[ERROR]" + reset + "La extension en el parametro main no coincide con " + cyan + "[.dsk]" + reset)
		}
	} else {
		fmt.Println(yellow + "[AVISO]" + reset + "La ruta especificada no existe, se procedera a crearla")
		if CrearCarpeta(path) == true {
			fmt.Println(green + "[EXITO]" + reset + "La carpeta se ha creado en la ruta especificada")
			CrearArchivoBinario(path+name, sizeFinal)
			InsertarMBR(sizeFinal, path+name)
			fmt.Println(green + "[EXITO]" + reset + "El archivo " + magenta + name + reset + " se ha creado con exito")
		} else {
			fmt.Println(red + "[ERROR]" + reset + "La carpeta no se ha podido crear en la ruta especificada")
		}
	}

}

//CrearArchivoBinario =
func CrearArchivoBinario(nombre string, size int64) {
	file, error := os.Create(nombre)
	defer file.Close() //defer es la encargada de asegurar que una funcion es llamada (obliga a funcionar)
	if error != nil {
		fmt.Println(red + "[ERROR]" + reset + " NO SE HA PODIDO GENERAR EL DISCO: " + yellow + strings.ToUpper(nombre) + reset)
		//log.Fatal(error)
	}

	var cero int64 = 0
	//fmt.Println(unsafe.Sizeof(&cero))

	var valorBinario bytes.Buffer
	binary.Write(&valorBinario, binary.BigEndian, &cero)
	escribirBytes(file, valorBinario.Bytes())

	file.Seek(size-8, 0)

	escribirBytes(file, valorBinario.Bytes())

}

//InsertarMBR =
func InsertarMBR(size int64, path string) {
	var mbr MBR
	file, err := os.OpenFile(path, os.O_RDWR, 0755)

	defer file.Close()
	if err != nil {
		fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
	}

	tiempo := time.Now()
	formatoTiempo := tiempo.Format("01-02-2006")
	//tiempoTemp := []byte(formatoTiempo)
	var conversorTiempo [20]byte
	copy(conversorTiempo[:], formatoTiempo)
	/*
		for i := 0; i < len(tiempoTemp); i++ {
			if i < 20 {
				conversorTiempo[i] = tiempoTemp[i]
			}
		}
	*/
	mbr.MbrTamano = size
	mbr.MbrFechaCreacion = conversorTiempo
	mbr.MbrDiskSignature = rand.Int63()
	mbr.MbrPartition1.PartSize = 0
	mbr.MbrPartition1.PartStatus = 1
	mbr.MbrPartition2.PartSize = 0
	mbr.MbrPartition2.PartStatus = 1
	mbr.MbrPartition3.PartSize = 0
	mbr.MbrPartition3.PartStatus = 1
	mbr.MbrPartition4.PartSize = 0
	mbr.MbrPartition4.PartStatus = 1

	file.Seek(0, 0)
	s1 := &mbr

	//Se escribe el Struct MBR
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	escribirBytes(file, binario3.Bytes())
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//RMDISK-----RMDISK-----METODOS-----METODOS-----RMDISK-----RMDISK-----METODOS-----METODOS-----RMDISK-----RMDISK-----METODOS-----METODOS------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//EliminarDisco =
func EliminarDisco(path string) {
	fmt.Println(cyan + "path:	" + reset + path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(red + "[ERROR]" + reset + "El archivo en la ruta especificada no existe")
	} else {
		for {
			fmt.Println(yellow + "[ADVERTENCIA]" + reset + "REALMENTE DESEA ELIMINAR EL ARCHIVO?[y/n] ")
			var tecla string
			fmt.Scanln(&tecla)
			if tecla == "y" {
				err := os.Remove(path)
				if err != nil {
					//println(err)
					fmt.Println(red + "[ERROR]" + reset + "El archivo no se ha podido eliminar")
					break
				} else {
					fmt.Println(green + "[EXITO]" + reset + "El archivo ha sido eliminado con exito")
					break
				}
			} else if tecla == "n" {
				break
			} else {
				fmt.Println(red + "-" + tecla + "-" + " no es una opcion valida, ingresa alguna de las opciones indicadas" + reset)
			}
		}
	}

}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//RMDISK-----RMDISK-----METODOS-----METODOS-----RMDISK-----RMDISK-----METODOS-----METODOS-----RMDISK-----RMDISK-----METODOS-----METODOS------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//CrearParticion =
func CrearParticion(size, unit, path, type2, fit, name string) {
	var mbr MBR = leerMBR(path)
	var espacioLibreTotal int64 = 0
	var existeExtendida bool = false
	var existeNombre bool = false
	espacioLibreTotal = mbr.MbrTamano - int64(unsafe.Sizeof(mbr))
	//Se recupera el espacio total del disco
	tamanoParticion, _ := strconv.ParseInt(size, 10, 64)

	//Se verifica si no se ha llegado al limite de 4 particiones
	if mbr.MbrPartition1.PartStatus == 0 && mbr.MbrPartition2.PartStatus == 0 && mbr.MbrPartition3.PartStatus == 0 && mbr.MbrPartition4.PartStatus == 0 {
		fmt.Println(red + "[ERROR]" + reset + "Se ha llegado al limite de particiones")
	} else {
		//Si es una extendida, se verifica si ya existe alguna extendida
		if type2 == "E" {
			if mbr.MbrPartition1.PartType == 'E' || mbr.MbrPartition2.PartType == 'E' || mbr.MbrPartition3.PartType == 'E' || mbr.MbrPartition4.PartType == 'E' {
				existeExtendida = true
				fmt.Println(red + "[ERROR]" + reset + "Ya existe una particion extendida")
			}
		}

		var nombreAByte16 [16]byte
		copy(nombreAByte16[:], name)

		if nombreAByte16 == mbr.MbrPartition1.PartName || nombreAByte16 == mbr.MbrPartition2.PartName || nombreAByte16 == mbr.MbrPartition3.PartName || nombreAByte16 == mbr.MbrPartition4.PartName {
			existeNombre = true
			fmt.Println(red + "[ERROR]" + reset + "Ya existe una particion con este nombre")
		}

		//Si no existe extendida puede pasar, independientemente de si es Logica, Primario o Extendida
		if existeExtendida == false && existeNombre == false {
			//Verifica el espacio disponible total
			if mbr.MbrPartition1.PartStatus == 0 {
				espacioLibreTotal = espacioLibreTotal - mbr.MbrPartition1.PartSize
			}
			if mbr.MbrPartition2.PartStatus == 0 {
				espacioLibreTotal = espacioLibreTotal - mbr.MbrPartition2.PartSize
			}
			if mbr.MbrPartition3.PartStatus == 0 {
				espacioLibreTotal = espacioLibreTotal - mbr.MbrPartition3.PartSize
			}
			if mbr.MbrPartition4.PartStatus == 0 {
				espacioLibreTotal = espacioLibreTotal - mbr.MbrPartition4.PartSize
			}

			//Se hace el calculo del tamaño de la particion a crear
			if unit == "K" {
				tamanoParticion = tamanoParticion * 1024
			} else if unit == "M" {
				tamanoParticion = tamanoParticion * 1024 * 1024
			}

			if mbr.MbrPartition1.PartStatus == 1 {
				//INGRESA SEGUN SEA SU TIPO DE PARTICION
				esExtendida := false
				mbr.MbrPartition1.PartStatus = 0
				if type2 == "P" {
					mbr.MbrPartition1.PartType = 'P'
				} else if type2 == "E" {
					esExtendida = true
					mbr.MbrPartition1.PartType = 'E'
				}
				var conversorFit [1]byte
				copy(conversorFit[:], fit)
				mbr.MbrPartition1.PartFit = conversorFit[0]
				mbr.MbrPartition1.PartStart = int64(unsafe.Sizeof(mbr) + 1)
				mbr.MbrPartition1.PartSize = tamanoParticion
				var conversorName [16]byte
				copy(conversorName[:], name)
				mbr.MbrPartition1.PartName = conversorName

				/*
					fmt.Println(cyan + "----------" + reset)
					fmt.Println(mbr.MbrTamano)
					fmt.Println(mbr.MbrFechaCreacion)
					fmt.Println(mbr.MbrDiskSignature)
					fmt.Println(magenta + "----------" + reset)
					fmt.Println(mbr.MbrPartition1.PartStatus)
					fmt.Println(mbr.MbrPartition1.PartType)
					fmt.Println(mbr.MbrPartition1.PartFit)
					fmt.Println(mbr.MbrPartition1.PartStart)
					fmt.Println(mbr.MbrPartition1.PartSize)
					fmt.Println(mbr.MbrPartition1.PartName)
					fmt.Println(magenta + "----------" + reset)
				*/

				ActualizarMBREInsertarParticion(path, &mbr, '1')
				fmt.Println(green + "[EXITO]" + reset + "La particion ha sido creada con exito")

				if esExtendida == true {
					CrearPrimerEBR(path, mbr.MbrPartition1.PartStart)
					var ebr2 EBR = leerEBR(path, mbr.MbrPartition1.PartStart)

					if ebr2.PartStatus == 1 {
						fmt.Println("MAMARRRRRRREEEEEEEEEEEEEEEEE")
					}
				}
			} else if mbr.MbrPartition2.PartStatus == 1 {
				//INGRESA SEGUN SEA SU TIPO DE PARTICION
				esExtendida := false
				mbr.MbrPartition2.PartStatus = 0
				if type2 == "P" {
					mbr.MbrPartition2.PartType = 'P'
				} else if type2 == "E" {
					esExtendida = true
					mbr.MbrPartition2.PartType = 'E'
				}
				var conversorFit [1]byte
				copy(conversorFit[:], fit)
				mbr.MbrPartition2.PartFit = conversorFit[0]
				mbr.MbrPartition2.PartStart = mbr.MbrPartition1.PartStart + mbr.MbrPartition1.PartSize + 1
				mbr.MbrPartition2.PartSize = tamanoParticion
				var conversorName [16]byte
				copy(conversorName[:], name)
				mbr.MbrPartition2.PartName = conversorName

				ActualizarMBREInsertarParticion(path, &mbr, '2')
				fmt.Println(green + "[EXITO]" + reset + "La particion ha sido creada con exito")

				if esExtendida == true {
					CrearPrimerEBR(path, mbr.MbrPartition2.PartStart)
				}
			} else if mbr.MbrPartition3.PartStatus == 1 {
				//INGRESA SEGUN SEA SU TIPO DE PARTICION]
				esExtendida := false
				mbr.MbrPartition3.PartStatus = 0
				if type2 == "P" {
					mbr.MbrPartition3.PartType = 'P'
				} else if type2 == "E" {
					esExtendida = true
					mbr.MbrPartition3.PartType = 'E'
				}
				var conversorFit [1]byte
				copy(conversorFit[:], fit)
				mbr.MbrPartition3.PartFit = conversorFit[0]
				mbr.MbrPartition3.PartStart = mbr.MbrPartition2.PartStart + mbr.MbrPartition2.PartSize + 1
				mbr.MbrPartition3.PartSize = tamanoParticion
				var conversorName [16]byte
				copy(conversorName[:], name)
				mbr.MbrPartition3.PartName = conversorName

				ActualizarMBREInsertarParticion(path, &mbr, '3')
				fmt.Println(green + "[EXITO]" + reset + "La particion ha sido creada con exito")

				if esExtendida == true {
					CrearPrimerEBR(path, mbr.MbrPartition3.PartStart)
				}
			} else if mbr.MbrPartition4.PartStatus == 1 {
				//INGRESA SEGUN SEA SU TIPO DE PARTICION
				esExtendida := false
				mbr.MbrPartition4.PartStatus = 0
				if type2 == "P" {
					mbr.MbrPartition4.PartType = 'P'
				} else if type2 == "E" {
					esExtendida = true
					mbr.MbrPartition4.PartType = 'E'
				}
				var conversorFit [1]byte
				copy(conversorFit[:], fit)
				mbr.MbrPartition4.PartFit = conversorFit[0]
				mbr.MbrPartition4.PartStart = mbr.MbrPartition3.PartStart + mbr.MbrPartition3.PartSize + 1
				mbr.MbrPartition4.PartSize = tamanoParticion
				var conversorName [16]byte
				copy(conversorName[:], name)
				mbr.MbrPartition4.PartName = conversorName

				ActualizarMBREInsertarParticion(path, &mbr, '4')
				fmt.Println(green + "[EXITO]" + reset + "La particion ha sido creada con exito")

				if esExtendida == true {
					CrearPrimerEBR(path, mbr.MbrPartition4.PartStart)
				}
			}
		}

		//Se obtiene el valor total de espacio libre en el disco

		if tamanoParticion > espacioLibreTotal {
			fmt.Println(red + "[ERROR]" + reset + "El tamaño de la particion supera el espacio disponible en el disco")
		} else {

		}
	}
}

//EliminarParticion = Metodo encargado de eliminar las particiones segun sean requeridas
func EliminarParticion(path, name, fit string) {
	var mbr MBR = leerMBR(path)
	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer file.Close()
	if err != nil {
		fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
	}

	var nombreAByte16 [16]byte
	copy(nombreAByte16[:], name)
	if nombreAByte16 == mbr.MbrPartition1.PartName {
		if fit == "FAST" {
			//Se reescribe el Struct MBR
			mbr.MbrPartition1.PartStatus = 1
			mbr.MbrPartition1.PartType = 0
			mbr.MbrPartition1.PartFit = 0
			mbr.MbrPartition1.PartStart = 0
			mbr.MbrPartition1.PartSize = 0
			var reinicio [16]byte
			mbr.MbrPartition1.PartName = reinicio

			file.Seek(0, 0)
			//Se reescribe el Struct MBR
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, mbr)
			escribirBytes(file, binario3.Bytes())
		} else if fit == "FULL" {
			//Se reescribe con 0's toda la particion
			file.Seek(mbr.MbrPartition1.PartStart, 0)
			for i := 0; i < int(mbr.MbrPartition1.PartSize); i++ {
				numero := byte(0)
				var valorBinario bytes.Buffer
				binary.Write(&valorBinario, binary.BigEndian, &numero)
				escribirBytes(file, valorBinario.Bytes())
			}

			//Se reescribe el Struct MBR
			mbr.MbrPartition1.PartStatus = 1
			mbr.MbrPartition1.PartType = 0
			mbr.MbrPartition1.PartFit = 0
			mbr.MbrPartition1.PartStart = 0
			mbr.MbrPartition1.PartSize = 0
			var reinicio [16]byte
			mbr.MbrPartition1.PartName = reinicio

			file.Seek(0, 0)
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, mbr)
			escribirBytes(file, binario3.Bytes())
		} else {
			fmt.Println(red + "[ERROR]" + reset + "El tipo de formateo no esta r")
		}
	} else if nombreAByte16 == mbr.MbrPartition2.PartName {
		if fit == "FAST" {
			//Se reescribe el Struct MBR
			mbr.MbrPartition2.PartStatus = 1
			mbr.MbrPartition2.PartType = 0
			mbr.MbrPartition2.PartFit = 0
			mbr.MbrPartition2.PartStart = 0
			mbr.MbrPartition2.PartSize = 0
			var reinicio [16]byte
			mbr.MbrPartition2.PartName = reinicio

			file.Seek(0, 0)
			//Se reescribe el Struct MBR
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, mbr)
			escribirBytes(file, binario3.Bytes())
		} else if fit == "FULL" {
			//Se reescribe con 0's toda la particion
			file.Seek(mbr.MbrPartition2.PartStart, 0)
			for i := 0; i < int(mbr.MbrPartition2.PartSize); i++ {
				numero := byte(0)
				var valorBinario bytes.Buffer
				binary.Write(&valorBinario, binary.BigEndian, &numero)
				escribirBytes(file, valorBinario.Bytes())
			}

			//Se reescribe el Struct MBR
			mbr.MbrPartition2.PartStatus = 1
			mbr.MbrPartition2.PartType = 0
			mbr.MbrPartition2.PartFit = 0
			mbr.MbrPartition2.PartStart = 0
			mbr.MbrPartition2.PartSize = 0
			var reinicio [16]byte
			mbr.MbrPartition2.PartName = reinicio

			file.Seek(0, 0)
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, mbr)
			escribirBytes(file, binario3.Bytes())
		} else {
			fmt.Println(red + "[ERROR]" + reset + "El tipo de formateo no esta r")
		}
	} else if nombreAByte16 == mbr.MbrPartition3.PartName {
		if fit == "FAST" {
			//Se reescribe el Struct MBR
			mbr.MbrPartition3.PartStatus = 1
			mbr.MbrPartition3.PartType = 0
			mbr.MbrPartition3.PartFit = 0
			mbr.MbrPartition3.PartStart = 0
			mbr.MbrPartition3.PartSize = 0
			var reinicio [16]byte
			mbr.MbrPartition3.PartName = reinicio

			file.Seek(0, 0)
			//Se reescribe el Struct MBR
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, mbr)
			escribirBytes(file, binario3.Bytes())
		} else if fit == "FULL" {
			//Se reescribe con 0's toda la particion
			file.Seek(mbr.MbrPartition3.PartStart, 0)
			for i := 0; i < int(mbr.MbrPartition3.PartSize); i++ {
				numero := byte(0)
				var valorBinario bytes.Buffer
				binary.Write(&valorBinario, binary.BigEndian, &numero)
				escribirBytes(file, valorBinario.Bytes())
			}

			//Se reescribe el Struct MBR
			mbr.MbrPartition3.PartStatus = 1
			mbr.MbrPartition3.PartType = 0
			mbr.MbrPartition3.PartFit = 0
			mbr.MbrPartition3.PartStart = 0
			mbr.MbrPartition3.PartSize = 0
			var reinicio [16]byte
			mbr.MbrPartition3.PartName = reinicio

			file.Seek(0, 0)
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, mbr)
			escribirBytes(file, binario3.Bytes())
		} else {
			fmt.Println(red + "[ERROR]" + reset + "El tipo de formateo no esta r")
		}
	} else if nombreAByte16 == mbr.MbrPartition4.PartName {
		if fit == "FAST" {
			//Se reescribe el Struct MBR
			mbr.MbrPartition4.PartStatus = 1
			mbr.MbrPartition4.PartType = 0
			mbr.MbrPartition4.PartFit = 0
			mbr.MbrPartition4.PartStart = 0
			mbr.MbrPartition4.PartSize = 0
			var reinicio [16]byte
			mbr.MbrPartition4.PartName = reinicio

			file.Seek(0, 0)
			//Se reescribe el Struct MBR
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, mbr)
			escribirBytes(file, binario3.Bytes())
		} else if fit == "FULL" {
			//Se reescribe con 0's toda la particion
			file.Seek(mbr.MbrPartition4.PartStart, 0)
			for i := 0; i < int(mbr.MbrPartition4.PartSize); i++ {
				numero := byte(0)
				var valorBinario bytes.Buffer
				binary.Write(&valorBinario, binary.BigEndian, &numero)
				escribirBytes(file, valorBinario.Bytes())
			}

			//Se reescribe el Struct MBR
			mbr.MbrPartition4.PartStatus = 1
			mbr.MbrPartition4.PartType = 0
			mbr.MbrPartition4.PartFit = 0
			mbr.MbrPartition4.PartStart = 0
			mbr.MbrPartition4.PartSize = 0
			var reinicio [16]byte
			mbr.MbrPartition4.PartName = reinicio

			file.Seek(0, 0)
			var binario3 bytes.Buffer
			binary.Write(&binario3, binary.BigEndian, mbr)
			escribirBytes(file, binario3.Bytes())
		} else {
			fmt.Println(red + "[ERROR]" + reset + "El tipo de formateo no esta r")
		}
	} else {
		fmt.Println(red + "[ERROR]" + reset + "No existe la particion a eliminar")
	}
}

//ModificarParticion =
func ModificarParticion() {
	fmt.Println("SOY MODIFICAR PARTICION")
}

//ActualizarMBREInsertarParticion =
func ActualizarMBREInsertarParticion(path string, mbr *MBR, numero byte) {
	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer file.Close()
	if err != nil {
		fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
	}

	file.Seek(0, 0)

	//Se reescribe el Struct MBR
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, mbr)
	escribirBytes(file, binario3.Bytes())

	if numero == '1' {
		file.Seek(mbr.MbrPartition1.PartStart, 0)
		for i := 0; i < int(mbr.MbrPartition1.PartSize); i++ {
			var valorBinario bytes.Buffer
			binary.Write(&valorBinario, binary.BigEndian, &numero)
			escribirBytes(file, valorBinario.Bytes())
		}
	} else if numero == '2' {
		file.Seek(mbr.MbrPartition2.PartStart, 0)
		for i := 0; i < int(mbr.MbrPartition2.PartSize); i++ {
			var valorBinario bytes.Buffer
			binary.Write(&valorBinario, binary.BigEndian, &numero)
			escribirBytes(file, valorBinario.Bytes())
		}
	} else if numero == '3' {
		file.Seek(mbr.MbrPartition3.PartStart, 0)
		for i := 0; i < int(mbr.MbrPartition3.PartSize); i++ {
			var valorBinario bytes.Buffer
			binary.Write(&valorBinario, binary.BigEndian, &numero)
			escribirBytes(file, valorBinario.Bytes())
		}
	} else if numero == '4' {
		file.Seek(mbr.MbrPartition4.PartStart, 0)
		for i := 0; i < int(mbr.MbrPartition4.PartSize); i++ {
			var valorBinario bytes.Buffer
			binary.Write(&valorBinario, binary.BigEndian, &numero)
			escribirBytes(file, valorBinario.Bytes())
		}
	}
}

//CrearPrimerEBR = Funcion encargada de crear el 1er EBR
func CrearPrimerEBR(path string, start int64) {
	var ebr EBR

	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer file.Close()
	if err != nil {
		fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
	}

	//fmt.Println(cyan+"tam 1: "+reset, int(unsafe.Sizeof(ebr)))
	fmt.Println(cyan+"tam 1: "+reset, int(binary.Size(ebr)))

	ebr.PartStatus = 1
	ebr.PartStart = start
	ebr.PartNext = -1

	file.Seek(start, 0)

	//Se escribe el Struct EBR
	var binarioEBR bytes.Buffer
	binary.Write(&binarioEBR, binary.BigEndian, &ebr)
	escribirBytes(file, binarioEBR.Bytes())
}

//InsertarParticionLogica = Metodo que se encarga de insertar las particiones logicas
func InsertarParticionLogica(path, name, size, fit string) {
	mbr := leerMBR(path)
	var ebr EBR
	siExisteExtendida := false
	var ajuste byte

	if fit == "B" {
		ajuste = 'B'
	} else if fit == "F" {
		ajuste = 'F'
	} else if fit == "W" {
		ajuste = 'W'
	}

	var nombreAByte16 [16]byte
	copy(nombreAByte16[:], name)
	if mbr.MbrPartition1.PartType == 'E' {
		ebr = leerEBR(path, mbr.MbrPartition1.PartStart)
		siExisteExtendida = true
	} else if mbr.MbrPartition2.PartType == 'E' {
		ebr = leerEBR(path, mbr.MbrPartition2.PartStart)
		siExisteExtendida = true
	} else if mbr.MbrPartition3.PartType == 'E' {
		ebr = leerEBR(path, mbr.MbrPartition3.PartStart)
		siExisteExtendida = true
	} else if mbr.MbrPartition4.PartType == 'E' {
		ebr = leerEBR(path, mbr.MbrPartition4.PartStart)
		siExisteExtendida = true
	} else {
		fmt.Println(red + "[ERROR]" + reset + "No existe una particion extendida para insertar los datos")
	}

	var tamEBR int64 = int64(binary.Size(ebr))
	var contadorLogica int64 = 0
	var sumadorStart int64 = ebr.PartStart
	var ebrAnterior EBR
	siguienteEBRVacio := false
	if siExisteExtendida == true {
		file, err := os.OpenFile(path, os.O_RDWR, 0755)
		defer file.Close()
		if err != nil {
			fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
		}
		for {
			if ebr.PartStatus == 1 { // Solo sirve para llenar el primer EBR
				fmt.Println("soy logica libre")

				ebr.PartStatus = 0
				ebr.PartFit = ajuste
				ebr.PartStart = sumadorStart
				sizeConvertido, _ := strconv.ParseInt(size, 10, 64)
				ebr.PartSize = sizeConvertido
				ebr.PartNext = -1
				ebr.PartName = nombreAByte16

				file.Seek(sumadorStart, 0)
				fmt.Println(cyan+"Sumador: "+reset, sumadorStart)
				//Se escribe el Struct EBR
				var binarioEBR bytes.Buffer
				binary.Write(&binarioEBR, binary.BigEndian, &ebr)
				escribirBytes(file, binarioEBR.Bytes())
				break
			} else if siguienteEBRVacio == true { //
				fmt.Println("Soy logica nueva")

				file.Seek(sumadorStart, 0)
				fmt.Println(cyan+"Sumador: "+reset, sumadorStart)

				var ebrTemporal EBR
				ebrTemporal.PartStatus = 0
				fmt.Println(cyan+"Ajuste: "+reset, ajuste)
				ebrTemporal.PartFit = ajuste
				ebrTemporal.PartStart = sumadorStart
				sizeConvertido, _ := strconv.ParseInt(size, 10, 64)
				ebrTemporal.PartSize = sizeConvertido
				ebrTemporal.PartNext = -1
				ebrTemporal.PartName = nombreAByte16

				var binarioEBR2 bytes.Buffer
				binary.Write(&binarioEBR2, binary.BigEndian, &ebrTemporal)
				escribirBytes(file, binarioEBR2.Bytes())

				file.Seek(ebrAnterior.PartStart, 0)
				fmt.Println(cyan+"posicion anterior particion logica: "+reset, ebrAnterior.PartStart)
				ebrAnterior.PartNext = sumadorStart
				var binarioEBRAnterior bytes.Buffer
				binary.Write(&binarioEBRAnterior, binary.BigEndian, &ebrAnterior)
				escribirBytes(file, binarioEBRAnterior.Bytes())
				break
			} else { //Recorre los ebr para encontra uno nuevo
				ebrAnterior = leerEBR(path, sumadorStart)

				sumadorStart += tamEBR + ebr.PartSize
				contadorLogica++

				if ebrAnterior.PartNext == -1 {
					siguienteEBRVacio = true
				} else {
					//numPartNext := ebrAnterior.PartNext
					//ebrTransitorio = leerEBR(path, numPartNext)
					fmt.Println(cyan+"Sumador: "+reset, sumadorStart)
					file.Seek(sumadorStart, 0)
				}
				fmt.Println("soy logica ocupada")
			}
		}
	}
}

func cuantoEspacioDisponibleHayParaLogica(path, name string, size int64) int64 {
	var resultado int64

	return resultado
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MOUNT-----MOUNT-----METODOS-----METODOS-----MOUNT-----MOUNT-----METODOS-----METODOS-----MOUNT-----MOUNT-----METODOS-----METODOS------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//MontarParticion = metodo encargada de montar las particiones mediante el comando MOUNT
func MontarParticion(path, name string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(red + "[ERROR]" + reset + "El disco a montar en la ruta especificada no existe")
	} else {
		if existeNombreParticion(path, name) == false {
			fmt.Println(red + "[ERROR]" + reset + "El nombre de la particion a montar no existe")
		} else {
			//************************************************************
			//Si el arreglo esta vacio agrega el primer valor a esta mamada
			//************************************************************
			if len(ContenedorMount) == 0 {
				nombreID := "vd"
				nombreID += "a1"
				var nmontaje NodoMontaje
				nmontaje.ID = nombreID
				nmontaje.Path = path
				nmontaje.NombreParticion = name
				nmontaje.Letra = "a"
				nmontaje.Numero = 1

				ContenedorMount = append(ContenedorMount, nmontaje)
			} else {
				nombreID := "vd"
				//************************************************************
				//Verificar si ya existe la letra o cuantas van para una nueva
				//************************************************************
				abecedario := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
				contadorLetrasDiferentes := 0
				QueLetraAsigno := ""
				yaEncontre := false
				yaEstaLaParticionMontada := false
				for i := 0; i < len(ContenedorMount); i++ {
					if ContenedorMount[i].Path == path {
						//************************************************************
						//Si encuentra que ya existe la particion montada
						//retorna un mensaje de error
						//************************************************************
						var retornarID string
						for j := 0; j < len(ContenedorMount); j++ {
							if ContenedorMount[j].Path == path && ContenedorMount[j].NombreParticion == name {
								yaEstaLaParticionMontada = true
								retornarID = ContenedorMount[j].ID
							}
						}
						if yaEstaLaParticionMontada == true {
							yaEncontre = true
							fmt.Println(red + "[ERROR]" + reset + "La particion indicada ya se encuentra montada con el ID: " + cyan + retornarID + reset)
						} else {
							//************************************************************
							//Entra aca cuando ya existe el path montado, verifica el No.
							//de particion y lo agrega.
							//************************************************************
							nombreID += ContenedorMount[i].Letra //vd'letra'
							QueLetraAsigno = ContenedorMount[i].Letra
							var vectorLetras []string
							for j := 0; j < len(ContenedorMount); j++ {
								vectorLetras = append(vectorLetras, ContenedorMount[j].Path)
							}
							num := contarNumeroMismaParticion(ContenedorMount, path) + 1
							nombreID += strconv.Itoa(num) //vdx'numero'
							contadorLetrasDiferentes = num
							yaEncontre = true
							break
						}
					}
				}
				if yaEncontre == false {
					//************************************************************
					//Entra aca cuando no existe el path y genera una nueva letra
					//para estas particiones
					//************************************************************
					var vectorLetras []string
					for i := 0; i < len(ContenedorMount); i++ {
						vectorLetras = append(vectorLetras, ContenedorMount[i].Path)
					}
					contadorLetrasDiferentes = contarCuantasLetras(vectorLetras) //Retorna el numero para el arreglo de letras en uso
					QueLetraAsigno = abecedario[contadorLetrasDiferentes]
					nombreID += QueLetraAsigno + "1"
				}

				if yaEstaLaParticionMontada == false {
					var nmontaje NodoMontaje
					nmontaje.ID = nombreID
					nmontaje.Path = path
					nmontaje.NombreParticion = name
					nmontaje.Letra = QueLetraAsigno
					nmontaje.Numero = contadorLetrasDiferentes

					ContenedorMount = append(ContenedorMount, nmontaje)
				}
			}
		}
	}
}

//contarCuantasLetras = Recorre el vector de path que recive para verificar cuantas veces se repite la ruta
//y devuelve un int con la cantidad de veces encontrado
func contarCuantasLetras(vector []string) int {
	Total := 0
	var yaEncontrados []string
	for i := 0; i < len(vector); i++ {
		//fmt.Println("ruta: " + vector[i])
		ya := false
		for j := 0; j < len(yaEncontrados); j++ {
			if yaEncontrados[j] == vector[i] {
				ya = true
			}
		}
		if ya == false {
			cuantos := 0
			for j := 0; j < len(vector); j++ {
				if vector[i] == vector[j] {
					cuantos++
					yaEncontrados = append(yaEncontrados, vector[i])
				}
			}
			if cuantos != 0 {
				Total++
				//fmt.Println("Total: ", Total)
			}
		}
	}
	return Total
}

func contarNumeroMismaParticion(vector []NodoMontaje, path string) int {
	Total := 0
	for i := 0; i < len(vector); i++ {
		if path == vector[i].Path {
			Total++
		}
	}
	return Total
}

func existeNombreParticion(path, name string) bool {
	var mbr MBR = leerMBR(path)
	var nombreAByte16 [16]byte
	copy(nombreAByte16[:], name)

	if nombreAByte16 == mbr.MbrPartition1.PartName || nombreAByte16 == mbr.MbrPartition2.PartName || nombreAByte16 == mbr.MbrPartition3.PartName || nombreAByte16 == mbr.MbrPartition4.PartName {
		return true
	}

	return false
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MOUNT-----MOUNT-----METODOS-----METODOS-----MOUNT-----MOUNT-----METODOS-----METODOS-----MOUNT-----MOUNT-----METODOS-----METODOS------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//DesmontarParticion = funcion encargada de desmontar la particion
func DesmontarParticion(idParticion string) {
	fueEncontrado := false
	for i := 0; i < len(ContenedorMount); i++ {
		if ContenedorMount[i].ID == idParticion {
			var nmontaje NodoMontaje
			//************************************************************
			//Sustituye el indice deseado para mover el resto de elementos
			//un espacio hacia la izquierda
			//************************************************************
			copy(ContenedorMount[i:], ContenedorMount[i+1:])
			//Se elimina el ultimo elemento con un NodoMontaje vacio
			ContenedorMount[len(ContenedorMount)-1] = nmontaje
			//Trunca el ultimo valor
			ContenedorMount = ContenedorMount[:len(ContenedorMount)-1]
			fueEncontrado = true
			fmt.Println(green + "[EXITO]" + reset + "El id " + cyan + idParticion + reset + " fue desmontado con exito")
		}
	}
	if fueEncontrado == false {
		fmt.Println(red + "[ERROR]" + reset + "El id " + cyan + idParticion + reset + " no se encuentra montado")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//RESUMEN-----RESUMEN-----RESUMEN-----RESUMEN-----RESUMEN-----RESUMEN-----RESUMEN-----RESUMEN-----RESUMEN-----RESUMEN-----RESUMEN-----RESUMEN
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//ResumenMBR = Genera un reporte en consola de todas las propiedades del mbr
func ResumenMBR(path string) {
	var mbr MBR
	file, err := os.OpenFile(path, os.O_RDWR, 0755)

	defer file.Close()
	if err != nil {
		fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
	}

	var tamanoEnBytes int = int(unsafe.Sizeof(mbr))

	file.Seek(0, 0)
	data := leerBytes(file, tamanoEnBytes)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(magenta + "|                           RESUMEN MBR                             |" + reset)
	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(cyan+"Tamaño Disco: "+reset, mbr.MbrTamano)
	//fmt.Println(mbr.MbrFechaCreacion)
	mamarre := string(mbr.MbrFechaCreacion[:])
	fmt.Println(cyan + "Fecha Creacion: " + reset + mamarre)
	fmt.Println(cyan+"Disk Signature: "+reset, mbr.MbrDiskSignature)
	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(magenta + "|                        RESUMEN PARTICION 1                        |" + reset)
	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(cyan+"Status: "+reset, mbr.MbrPartition1.PartStatus)
	typeP1 := string(mbr.MbrPartition1.PartType)
	fmt.Println(cyan+"Type: "+reset, typeP1)
	fitP1 := string(mbr.MbrPartition1.PartFit)
	fmt.Println(cyan+"Fit: "+reset, fitP1)
	fmt.Println(cyan+"Start: "+reset, mbr.MbrPartition1.PartStart)
	fmt.Println(cyan+"Size: "+reset, mbr.MbrPartition1.PartSize)
	mamarreP1 := string(mbr.MbrPartition1.PartName[:])
	fmt.Println(cyan + "Name: " + reset + mamarreP1)
	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(magenta + "|                        RESUMEN PARTICION 2                        |" + reset)
	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(cyan+"Status: "+reset, mbr.MbrPartition2.PartStatus)
	typeP2 := string(mbr.MbrPartition2.PartType)
	fmt.Println(cyan+"Type: "+reset, typeP2)
	fitP2 := string(mbr.MbrPartition2.PartFit)
	fmt.Println(cyan+"Fit: "+reset, fitP2)
	fmt.Println(cyan+"Start: "+reset, mbr.MbrPartition2.PartStart)
	fmt.Println(cyan+"Size: "+reset, mbr.MbrPartition2.PartSize)
	mamarreP2 := string(mbr.MbrPartition2.PartName[:])
	fmt.Println(cyan + "Name: " + reset + mamarreP2)
	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(magenta + "|                        RESUMEN PARTICION 3                        |" + reset)
	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(cyan+"Status: "+reset, mbr.MbrPartition3.PartStatus)
	typeP3 := string(mbr.MbrPartition3.PartType)
	fmt.Println(cyan+"Type: "+reset, typeP3)
	fitP3 := string(mbr.MbrPartition3.PartFit)
	fmt.Println(cyan+"Fit: "+reset, fitP3)
	fmt.Println(cyan+"Start: "+reset, mbr.MbrPartition3.PartStart)
	fmt.Println(cyan+"Size: "+reset, mbr.MbrPartition3.PartSize)
	mamarreP3 := string(mbr.MbrPartition3.PartName[:])
	fmt.Println(cyan + "Name: " + reset + mamarreP3)
	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(magenta + "|                        RESUMEN PARTICION 4                        |" + reset)
	fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
	fmt.Println(cyan+"Status: "+reset, mbr.MbrPartition4.PartStatus)
	typeP4 := string(mbr.MbrPartition4.PartType)
	fmt.Println(cyan+"Type: "+reset, typeP4)
	fitP4 := string(mbr.MbrPartition4.PartFit)
	fmt.Println(cyan+"Fit: "+reset, fitP4)
	fmt.Println(cyan+"Start: "+reset, mbr.MbrPartition4.PartStart)
	fmt.Println(cyan+"Size: "+reset, mbr.MbrPartition4.PartSize)
	mamarreP4 := string(mbr.MbrPartition4.PartName[:])
	fmt.Println(cyan + "Name: " + reset + mamarreP4)
}

//ResumenParticionesMontadas =
func ResumenParticionesMontadas() {
	if len(ContenedorMount) == 0 {
		fmt.Println(red + "[ERROR]" + reset + "No existen particiones montadas")
	} else {
		for i := 0; i < len(ContenedorMount); i++ {
			fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
			fmt.Print(magenta + "|                RESUMEN PARTICION MONTADA ")
			fmt.Print(i)
			fmt.Println("                        |" + reset)
			fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
			fmt.Println(cyan + "ID: 	" + reset + ContenedorMount[i].ID)
			fmt.Println(cyan + "PATH: 	" + reset + ContenedorMount[i].Path)
			fmt.Println(cyan + "NAME: 	" + reset + ContenedorMount[i].NombreParticion)
			fmt.Println(cyan + "LETRA: 	" + reset + ContenedorMount[i].Letra)
			fmt.Println(cyan+"NUMERO:"+reset, ContenedorMount[i].Numero)
		}
	}
}

//ResumenEBR =
func ResumenEBR(path, name string) {
	mbr := leerMBR(path)
	var ebr EBR
	siExisteExtendida := false

	var nombreAByte16 [16]byte
	copy(nombreAByte16[:], name)
	if mbr.MbrPartition1.PartType == 'E' {
		ebr = leerEBR(path, mbr.MbrPartition1.PartStart)
		siExisteExtendida = true
	} else if mbr.MbrPartition2.PartType == 'E' {
		ebr = leerEBR(path, mbr.MbrPartition2.PartStart)
		siExisteExtendida = true
	} else if mbr.MbrPartition3.PartType == 'E' {
		ebr = leerEBR(path, mbr.MbrPartition3.PartStart)
		siExisteExtendida = true
	} else if mbr.MbrPartition4.PartType == 'E' {
		ebr = leerEBR(path, mbr.MbrPartition4.PartStart)
		siExisteExtendida = true
	} else {
		fmt.Println(red + "[ERROR]" + reset + "No existe una particion extendida para generar el resumen de particiones logicas")
	}

	var tamEBR int64 = int64(binary.Size(ebr))
	var contadorLogica int64 = 1
	var sumadorStart int64 = ebr.PartStart
	var ebrAnterior EBR
	if siExisteExtendida == true {
		file, err := os.OpenFile(path, os.O_RDWR, 0755)
		defer file.Close()
		if err != nil {
			fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
		}
		for {
			if contadorLogica == 1 { // Solo sirve para llenar el primer EBR
				fmt.Println("IF 1")
				fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
				fmt.Println(magenta + "|                          RESUMEN LOGICA " + strconv.FormatInt(contadorLogica, 10) + "                         |" + reset)
				fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
				fmt.Println(cyan+"Status: "+reset, ebr.PartStatus)
				FitPL := string(ebr.PartFit)
				fmt.Println(cyan+"Fit: "+reset, FitPL)
				fmt.Println(cyan+"Start: "+reset, ebr.PartStart)
				fmt.Println(cyan+"Size: "+reset, ebr.PartSize)
				fmt.Println(cyan+"Next: "+reset, ebr.PartNext)
				mamarreP1 := string(ebr.PartName[:])
				fmt.Println(cyan + "Name: " + reset + mamarreP1)
				contadorLogica++
			} else {
				ebrAnterior = leerEBR(path, sumadorStart)

				sumadorStart += tamEBR + ebr.PartSize

				if ebrAnterior.PartNext == -1 {
					fmt.Println("IF 2 EXIT")
					break
				} else {
					fmt.Println("IF 2")
					ebrAnterior = leerEBR(path, sumadorStart)
					fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
					fmt.Println(magenta + "|                          RESUMEN LOGICA " + strconv.FormatInt(contadorLogica, 10) + "                         |" + reset)
					fmt.Println(magenta + "---------------------------------------------------------------------" + reset)
					fmt.Println(cyan+"Status: "+reset, ebrAnterior.PartStatus)
					FitPL := string(ebrAnterior.PartFit)
					fmt.Println(cyan+"Fit: "+reset, FitPL)
					fmt.Println(cyan+"Start: "+reset, ebrAnterior.PartStart)
					fmt.Println(cyan+"Size: "+reset, ebrAnterior.PartSize)
					fmt.Println(cyan+"Next: "+reset, ebrAnterior.PartNext)
					mamarreP1 := string(ebrAnterior.PartName[:])
					fmt.Println(cyan + "Name: " + reset + mamarreP1)
					file.Seek(sumadorStart, 0)
				}
				contadorLogica++
			}
		}
	}
}
