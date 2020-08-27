package metodos

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

//existeID = busca el id en la lista de particiones montadas para retonar el path de la particion y retorna un true o false
//segun sea el caso
func existeID(id string) (string, bool) {
	existe := false
	var path string
	if len(ContenedorMount) == 0 {
		fmt.Println(red + "[ERROR]" + reset + "No existen particiones montadas")
	} else {
		for i := 0; i < len(ContenedorMount); i++ {
			if ContenedorMount[i].ID == id {
				path = ContenedorMount[i].Path
				existe = true
				break
			}
		}
	}
	return path, existe
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

//SUPERBOOT es el struct que contiene el Super Bloque
type SUPERBOOT struct {
	SbNombreHd                      [16]byte
	SbArbolVirtualCount             int64
	SbDetalleDirectorioCount        int64
	SbInodosCount                   int64
	SbBloquesCount                  int64
	SbArbonVirtualFree              int64
	SbDetalleDirectorioFree         int64
	SbInodosFree                    int64
	SbBloquesFree                   int64
	SbDateCreacion                  [16]byte
	SbDateUltimoMontaje             [16]byte
	SbMontajeCount                  int64
	SbApBitmapArbolDirectorio       int64
	SbApArbolDirectorio             int64
	SbApBitmapDetalleDirectorio     int64
	SbApDetalleDirectorio           int64
	SbApBitmapTablaInodo            int64
	SbApTablaInodo                  int64
	SbApBitmapBloques               int64
	SbApBloques                     int64
	SbApLog                         int64
	SbSizeStructArbolDirectorio     int64
	SbSizeStructDetalleDirectorio   int64
	SbSizeStructInodo               int64
	SbSizeStructBloque              int64
	SbFirstFreeBitArbolDirectorio   int64
	SbFirstFreeBitDetalleDirectorio int64
	SbFirstFreeBitTablaInodo        int64
	SbFirstFreeBitBloques           int64
	SbMagicNum                      int64
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
	formatoTiempo := tiempo.Format("01-02-2006 15:04:05")
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
//FDISK-----FDISK-----METODOS-----METODOS-----FDISK-----FDISK-----METODOS-----METODOS-----FDISK-----FDISK-----METODOS-----METODOS-----FDISK--
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

		//Se obtiene el valor total de espacio libre en el disco
		if tamanoParticion > espacioLibreTotal {
			fmt.Println("T PARTICION:", tamanoParticion)
			fmt.Println("T ESPACIO LIBRE:", espacioLibreTotal)
			fmt.Println(red + "[ERROR]" + reset + "El tamaño de la particion supera el espacio disponible en el disco")
		} else {

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
	//fmt.Println(cyan+"tam 1: "+reset, int(binary.Size(ebr)))

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
//UNMOUNT-----UNMOUNT-----METODOS-----METODOS-----UNMOUNT-----UNMOUNT-----METODOS-----METODOS-----UNMOUNT-----UNMOUNT-----METODOS-----METODOS
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
//MKFS-----MKFS-----METODOS-----METODOS-----MKFS-----MKFS-----METODOS-----METODOS-----MKFS-----MKFS-----METODOS-----METODOS-----MKFS-----MKFS
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

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

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//REPORTE-----REPORTE-----REPORTE-----REPORTE-----REPORTE-----REPORTE-----REPORTE-----REPORTE-----REPORTE-----REPORTE-----REPORTE-----REPORTE
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//ReporteMBR = Metodo que genera el reporte mediante la funcion REP del MBR
func ReporteMBR(id, path string) {
	if len(ContenedorMount) == 0 {
		fmt.Println(red + "[ERROR]" + reset + "No existen particiones montadas")
	} else {
		ExisteID := false
		var cadenaReporteMBR string
		for i := 0; i < len(ContenedorMount); i++ {
			if ContenedorMount[i].ID == id {
				//******************************************************************
				//Se llena la cadena que sera ingresada en el .svg para uso de
				//Graphviz
				//******************************************************************
				mbr := leerMBR(ContenedorMount[i].Path)
				cadenaReporteMBR = "digraph MBR {\nnode [shape=plaintext]\nA [label=<\n<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\">\n"
				cadenaReporteMBR += "<TR>\n<TD BGCOLOR=\"#fa4251\" COLSPAN=\"2\">REPORTE MBR</TD>\n</TR>\n"
				cadenaReporteMBR += "<TR>\n<TD BGCOLOR=\"#ff6363\">NOMBRE</TD><TD BGCOLOR=\"#ff6363\">VALOR</TD>\n</TR>\n"
				cadenaReporteMBR += "<TR>\n<TD>MBR_Tamaño</TD><TD>" + strconv.FormatInt(mbr.MbrTamano, 10) + "</TD>\n</TR>\n"
				fechaConvertida := string(mbr.MbrFechaCreacion[:19])
				cadenaReporteMBR += "<TR>\n<TD>MBR_Fecha_Creacion</TD><TD>" + fechaConvertida + "</TD>\n</TR>\n"
				cadenaReporteMBR += "<TR>\n<TD>MBR_Disk_Signature</TD><TD>" + strconv.FormatInt(mbr.MbrDiskSignature, 10) + "</TD>\n</TR>\n"
				//******************************************************************
				//PARTICION 1
				//******************************************************************
				cadenaReporteMBR += "<TR>\n<TD BGCOLOR=\"#FF8585\" COLSPAN=\"2\">PARTICION 1</TD>\n</TR>\n"
				if mbr.MbrPartition1.PartStatus == 1 {
					cadenaReporteMBR += "<TR>\n<TD>Part_Status</TD><TD>1</TD>\n</TR>\n"
				} else {
					cadenaReporteMBR += "<TR>\n<TD>Part_Status</TD><TD>0</TD>\n</TR>\n"
					typeP1 := string(mbr.MbrPartition1.PartType)
					cadenaReporteMBR += "<TR>\n<TD>Part_Type</TD><TD>" + typeP1 + "</TD>\n</TR>\n"
					fitP1 := string(mbr.MbrPartition1.PartFit)
					cadenaReporteMBR += "<TR>\n<TD>Part_Fit</TD><TD>" + fitP1 + "</TD>\n</TR>\n"
					cadenaReporteMBR += "<TR>\n<TD>Part_Start</TD><TD>" + strconv.FormatInt(mbr.MbrPartition1.PartStart, 10) + "</TD>\n</TR>\n"
					cadenaReporteMBR += "<TR>\n<TD>Part_Size</TD><TD>" + strconv.FormatInt(mbr.MbrPartition1.PartSize, 10) + "</TD>\n</TR>\n"
					var nombreP1 string
					for i1, valor1 := range mbr.MbrPartition1.PartName {
						if mbr.MbrPartition1.PartName[i1] != 0 {
							nombreP1 += string(valor1)
						}
					}
					cadenaReporteMBR += "<TR>\n<TD>Part_Name</TD><TD>" + nombreP1 + "</TD>\n</TR>\n"
				}
				//******************************************************************
				//PARTICION 2
				//******************************************************************
				cadenaReporteMBR += "<TR>\n<TD BGCOLOR=\"#FF8585\" COLSPAN=\"2\">PARTICION 2</TD>\n</TR>\n"
				if mbr.MbrPartition2.PartStatus == 1 {
					cadenaReporteMBR += "<TR>\n<TD>Part_Status</TD><TD>1</TD>\n</TR>\n"
				} else {
					cadenaReporteMBR += "<TR>\n<TD>Part_Status</TD><TD>0</TD>\n</TR>\n"
					typeP2 := string(mbr.MbrPartition2.PartType)
					cadenaReporteMBR += "<TR>\n<TD>Part_Type</TD><TD>" + typeP2 + "</TD>\n</TR>\n"
					fitP2 := string(mbr.MbrPartition2.PartFit)
					cadenaReporteMBR += "<TR>\n<TD>Part_Fit</TD><TD>" + fitP2 + "</TD>\n</TR>\n"
					cadenaReporteMBR += "<TR>\n<TD>Part_Start</TD><TD>" + strconv.FormatInt(mbr.MbrPartition2.PartStart, 10) + "</TD>\n</TR>\n"
					cadenaReporteMBR += "<TR>\n<TD>Part_Size</TD><TD>" + strconv.FormatInt(mbr.MbrPartition2.PartSize, 10) + "</TD>\n</TR>\n"
					var nombreP2 string
					for i1, valor1 := range mbr.MbrPartition2.PartName {
						if mbr.MbrPartition2.PartName[i1] != 0 {
							nombreP2 += string(valor1)
						}
					}
					cadenaReporteMBR += "<TR>\n<TD>Part_Name</TD><TD>" + nombreP2 + "</TD>\n</TR>\n"
				}
				//******************************************************************
				//PARTICION 3
				//******************************************************************
				cadenaReporteMBR += "<TR>\n<TD BGCOLOR=\"#FF8585\" COLSPAN=\"2\">PARTICION 3</TD>\n</TR>\n"
				if mbr.MbrPartition3.PartStatus == 1 {
					cadenaReporteMBR += "<TR>\n<TD>Part_Status</TD><TD>1</TD>\n</TR>\n"
				} else {
					cadenaReporteMBR += "<TR>\n<TD>Part_Status</TD><TD>0</TD>\n</TR>\n"
					typeP3 := string(mbr.MbrPartition3.PartType)
					cadenaReporteMBR += "<TR>\n<TD>Part_Type</TD><TD>" + typeP3 + "</TD>\n</TR>\n"
					fitP3 := string(mbr.MbrPartition3.PartFit)
					cadenaReporteMBR += "<TR>\n<TD>Part_Fit</TD><TD>" + fitP3 + "</TD>\n</TR>\n"
					cadenaReporteMBR += "<TR>\n<TD>Part_Start</TD><TD>" + strconv.FormatInt(mbr.MbrPartition3.PartStart, 10) + "</TD>\n</TR>\n"
					cadenaReporteMBR += "<TR>\n<TD>Part_Size</TD><TD>" + strconv.FormatInt(mbr.MbrPartition3.PartSize, 10) + "</TD>\n</TR>\n"
					var nombreP3 string
					for i1, valor1 := range mbr.MbrPartition3.PartName {
						if mbr.MbrPartition3.PartName[i1] != 0 {
							nombreP3 += string(valor1)
						}
					}
					cadenaReporteMBR += "<TR>\n<TD>Part_Name</TD><TD>" + nombreP3 + "</TD>\n</TR>\n"
				}

				//******************************************************************
				//PARTICION 4
				//******************************************************************
				cadenaReporteMBR += "<TR>\n<TD BGCOLOR=\"#FF8585\" COLSPAN=\"2\">PARTICION 4</TD>\n</TR>\n"
				if mbr.MbrPartition4.PartStatus == 1 {
					cadenaReporteMBR += "<TR>\n<TD>Part_Status</TD><TD>1</TD>\n</TR>\n"
				} else {
					cadenaReporteMBR += "<TR>\n<TD>Part_Status</TD><TD>0</TD>\n</TR>\n"
					typeP4 := string(mbr.MbrPartition4.PartType)
					cadenaReporteMBR += "<TR>\n<TD>Part_Type</TD><TD>" + typeP4 + "</TD>\n</TR>\n"
					fitP4 := string(mbr.MbrPartition4.PartFit)
					cadenaReporteMBR += "<TR>\n<TD>Part_Fit</TD><TD>" + fitP4 + "</TD>\n</TR>\n"
					cadenaReporteMBR += "<TR>\n<TD>Part_Start</TD><TD>" + strconv.FormatInt(mbr.MbrPartition4.PartStart, 10) + "</TD>\n</TR>\n"
					cadenaReporteMBR += "<TR>\n<TD>Part_Size</TD><TD>" + strconv.FormatInt(mbr.MbrPartition4.PartSize, 10) + "</TD>\n</TR>\n"
					var nombreP4 string
					for i1, valor1 := range mbr.MbrPartition4.PartName {
						if mbr.MbrPartition4.PartName[i1] != 0 {
							nombreP4 += string(valor1)
						}
					}
					cadenaReporteMBR += "<TR>\n<TD>Part_Name</TD><TD>" + nombreP4 + "</TD>\n</TR>\n"
				}

				cadenaReporteMBR += "</TABLE>\n>];\n}"

				//******************************************************************
				//Se escribe la cadena en el archivo .svg que usara Graphviz
				//******************************************************************
				nombreGV, nombreExtension := crearArchivoParaReporte(path, cadenaReporteMBR)
				//******************************************************************
				//Aca se genera el la imagen, pdf segun sea ingresada
				//******************************************************************
				//cmd := exec.Command("dot", "-Tps", "/home/javier/Imágenes/graph1.gv", "-o", "/home/javier/Imágenes/gra.pdf")
				ruta, nombreArchivo := filepath.Split(path)
				nombreCompleto := ruta + nombreArchivo

				var cmd *exec.Cmd
				if nombreExtension == ".pdf" {
					cmd = exec.Command("dot", "-Tps", ruta+nombreGV, "-o", nombreCompleto)
				} else {
					cmd = exec.Command("dot", "-Tpng", ruta+nombreGV, "-o", nombreCompleto)
				}
				cmdOutput := &bytes.Buffer{}
				cmd.Stdout = cmdOutput
				err := cmd.Run()
				if err != nil {
					os.Stderr.WriteString(err.Error())
				}
				fmt.Print(string(cmdOutput.Bytes()))

				ExisteID = true
				break
			}
		}
		if ExisteID == false {
			fmt.Println(red + "[ERROR]" + reset + "El id " + cyan + id + reset + " no se encuentra montado")
		}
	}
}

//ReporteDISK = Metodo que genera el reporte mediante la funcion REP del MBR
func ReporteDISK(id, path string) {
	PathID, ExisteID := existeID(id)
	var cadenaReporteDISK string
	var cadenaReporteDISKAuxiliar string = ""

	if ExisteID == true {
		mbr := leerMBR(PathID)
		cadenaReporteDISK = "digraph {\n  A [\nshape=plaintext\nlabel=<\n"
		cadenaReporteDISK += "<table border='0' cellborder='1' BGCOLOR='#393939'  cellspacing='0'>\n"
		cadenaReporteDISK += "<tr>\n<td rowspan = '2'><font color='#00ad5f'>MBR</font></td>\n"

		if mbr.MbrPartition1.PartStatus == 1 {
			cadenaReporteDISK += "<td rowspan = '2'><font BGCOLOR = '#fa4251' color='#white'>Libre</font></td>\n"
		} else {
			if mbr.MbrPartition1.PartType == 'P' {
				cadenaReporteDISK += "<td rowspan = '2'><font color='#00ad5f'>Primaria</font></td>\n"
			} else {
				cadenaReporteDISK += "<td><font color='#00ad5f'>Extendida</font></td>\n"
				cadenaReporteDISKAuxiliar += "<table color='black' border='0' cellborder='1' cellpadding='10' cellspacing='0'>\n"
				cadenaReporteDISKAuxiliar += "<tr>\n"
				cadenaReporteDISKAuxiliar += "<td BGCOLOR = '#fa4251'><font color='white'>EBR</font></td>\n"
				cadenaReporteDISKAuxiliar += recorerParticionesLogicasReporte(PathID, mbr, 1)
				cadenaReporteDISKAuxiliar += "</tr>\n"
				cadenaReporteDISKAuxiliar += "</table>\n"
			}
		}
		if mbr.MbrPartition2.PartStatus == 1 {
			cadenaReporteDISK += "<td rowspan = '2'><font BGCOLOR = '#fa4251' color='#white'>Libre</font></td>\n"
		} else {
			if mbr.MbrPartition2.PartType == 'P' {
				cadenaReporteDISK += "<td rowspan = '2'><font color='#00ad5f'>Primaria</font></td>\n"
			} else {
				cadenaReporteDISK += "<td><font color='#00ad5f'>Extendida</font></td>\n"
				cadenaReporteDISKAuxiliar += "<table color='black' border='0' cellborder='1' cellpadding='10' cellspacing='0'>\n"
				cadenaReporteDISKAuxiliar += "<tr>\n"
				cadenaReporteDISKAuxiliar += "<td BGCOLOR = '#fa4251'><font color='white'>EBR</font></td>\n"
				cadenaReporteDISKAuxiliar += recorerParticionesLogicasReporte(PathID, mbr, 2)
				cadenaReporteDISKAuxiliar += "</tr>\n"
				cadenaReporteDISKAuxiliar += "</table>\n"
			}
		}
		if mbr.MbrPartition3.PartStatus == 1 {
			cadenaReporteDISK += "<td rowspan = '2'><font BGCOLOR = '#fa4251' color='white'>Libre</font></td>\n"
		} else {
			if mbr.MbrPartition3.PartType == 'P' {
				cadenaReporteDISK += "<td rowspan = '2'><font color='#00ad5f'>Primaria</font></td>\n"
			} else {
				cadenaReporteDISK += "<td><font color='#00ad5f'>Extendida</font></td>\n"
				cadenaReporteDISKAuxiliar += "<table color='black' border='0' cellborder='1' cellpadding='10' cellspacing='0'>\n"
				cadenaReporteDISKAuxiliar += "<tr>\n"
				cadenaReporteDISKAuxiliar += "<td BGCOLOR = '#fa4251'><font color='white'>EBR</font></td>\n"
				cadenaReporteDISKAuxiliar += recorerParticionesLogicasReporte(PathID, mbr, 3)
				cadenaReporteDISKAuxiliar += "</tr>\n"
				cadenaReporteDISKAuxiliar += "</table>\n"
			}
		}
		if mbr.MbrPartition4.PartStatus == 1 {
			cadenaReporteDISK += "<td rowspan = '2'><font BGCOLOR = '#fa4251' color='white'>Libre</font></td>\n"
		} else {
			if mbr.MbrPartition4.PartType == 'P' {
				cadenaReporteDISK += "<td rowspan = '2'><font color='#00ad5f'>Primaria</font></td>\n"
			} else {
				cadenaReporteDISK += "<td><font color='#00ad5f'>Extendida</font></td>\n"
				cadenaReporteDISKAuxiliar += "<table color='black' border='0' cellborder='1' cellpadding='10' cellspacing='0'>\n"
				cadenaReporteDISKAuxiliar += "<tr>\n"
				cadenaReporteDISKAuxiliar += "<td BGCOLOR = '#fa4251'><font color='white'>EBR</font></td>\n"
				cadenaReporteDISKAuxiliar += recorerParticionesLogicasReporte(PathID, mbr, 4)
				cadenaReporteDISKAuxiliar += "</tr>\n"
				cadenaReporteDISKAuxiliar += "</table>\n"
			}
		}
		cadenaReporteDISK += "</tr>\n"
		cadenaReporteDISK += "<tr>\n<td>\n"
		cadenaReporteDISK += cadenaReporteDISKAuxiliar
		cadenaReporteDISK += "</td>\n</tr>"
		cadenaReporteDISK += "</table>\n"
		cadenaReporteDISK += ">];\n}"

		//******************************************************************
		//Se escribe la cadena en el archivo .svg que usara Graphviz
		//******************************************************************
		nombreGV, nombreExtension := crearArchivoParaReporte(path, cadenaReporteDISK)
		//******************************************************************
		//Aca se genera el la imagen, pdf segun sea ingresada
		//******************************************************************
		//cmd := exec.Command("dot", "-Tps", "/home/javier/Imágenes/graph1.gv", "-o", "/home/javier/Imágenes/gra.pdf")
		ruta, nombreArchivo := filepath.Split(path)
		nombreCompleto := ruta + nombreArchivo

		var cmd *exec.Cmd
		if nombreExtension == ".pdf" {
			cmd = exec.Command("dot", "-Tps", ruta+nombreGV, "-o", nombreCompleto)
		} else {
			cmd = exec.Command("dot", "-Tpng", ruta+nombreGV, "-o", nombreCompleto)
		}
		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput
		err := cmd.Run()
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
		fmt.Print(string(cmdOutput.Bytes()))
	} else {

	}
}

func recorerParticionesLogicasReporte(path string, mbr MBR, queParticionExtendida int) string {
	var ebr EBR
	var cadena string = ""
	siExisteExtendida := false

	if queParticionExtendida == 1 {
		ebr = leerEBR(path, mbr.MbrPartition1.PartStart)
		siExisteExtendida = true
	} else if queParticionExtendida == 2 {
		ebr = leerEBR(path, mbr.MbrPartition2.PartStart)
		siExisteExtendida = true
	} else if queParticionExtendida == 3 {
		ebr = leerEBR(path, mbr.MbrPartition3.PartStart)
		siExisteExtendida = true
	} else if queParticionExtendida == 4 {
		ebr = leerEBR(path, mbr.MbrPartition4.PartStart)
		siExisteExtendida = true
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
				cadena += "<td BGCOLOR = '#393939'><font color='white'>Logica</font></td>\n"
				contadorLogica++
			} else {
				ebrAnterior = leerEBR(path, sumadorStart)
				sumadorStart += tamEBR + ebr.PartSize
				if ebrAnterior.PartNext == -1 {
					break
				} else {
					cadena += "<td BGCOLOR = '#fa4251'><font color='white'>EBR</font></td>\n"
					cadena += "<td BGCOLOR = '#393939'><font color='white'>Logica</font></td>\n"
					file.Seek(sumadorStart, 0)
				}
				contadorLogica++
			}
		}
	}
	return cadena
}

//crearArchivoParaReporte = Metodo que almacena la informacion en el archivo .dot o .svg (segun sea el caso)
//Primer parametro de retorno = retorna el nombre con la extension adecuada para generar el reporte en PDF o PNG, etc...
//Segundo parametro de retorno = regresa la extension, para saber que comando ejectura en el sistema
func crearArchivoParaReporte(path, cadena string) (string, string) {
	ruta, nombreArchivo := filepath.Split(path)
	//Ruta: /home/user/xxxx
	//Nombre Archivo: yyy.go
	var extension = filepath.Ext(nombreArchivo)                       //.go
	var nombre = nombreArchivo[0 : len(nombreArchivo)-len(extension)] //yyy
	//var nombreGV = nombre + ".gv"                                     //yyy.gv
	var nombreGV string
	if extension == ".pdf" {
		nombreGV = nombre + ".svg" //yyy.sgv
	} else {
		nombreGV = nombre + ".dot" //yyy.dot
	}

	if ExisteCarpeta(ruta) == false {
		CrearCarpeta(ruta)
	}
	if ExisteCarpeta(ruta) == true {
		//Se genera el archivo .gv (para uso de graphvz)
		file, err := os.Create(ruta + nombreGV)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		//Se escribe la informacion en el archivo
		err2 := ioutil.WriteFile(ruta+nombreGV, []byte(cadena), 0644)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	return nombreGV, extension
}
