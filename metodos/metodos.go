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

	var tamanoEnBytes int = int(unsafe.Sizeof(ebr))

	file.Seek(start, 0)
	data := leerBytes(file, tamanoEnBytes)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &ebr)
	if err != nil {
		fmt.Println(err)
	}

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

//EliminarParticion =
func EliminarParticion() {
	fmt.Println("SOY ELIMINAR PARTICION")
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

	ebr.PartStatus = 1
	ebr.PartStart = start

	file.Seek(start, 0)

	//Se escribe el Struct EBR
	var binarioEBR bytes.Buffer
	binary.Write(&binarioEBR, binary.BigEndian, &ebr)
	escribirBytes(file, binarioEBR.Bytes())
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
			var nombreID string
			nombreID = "vd"

			var nmontaje NodoMontaje
			nmontaje.ID = nombreID
			nmontaje.Path = path
			nmontaje.NombreParticion = name
			nmontaje.Letra = ""
			nmontaje.Numero = 0
		}
	}
	//ContenedorMount = append(ContenedorMount, temporal)
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
		}
	}
}
