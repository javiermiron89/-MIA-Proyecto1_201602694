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
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			log.Fatal(err)
			return false
		}
		return true
	}
	return true
	/*
		os.MkdirAll(path, 0777)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err != nil {
				//panic(err)
				return false
			}
			return true
		}
		return true
	*/
}

//existeID = busca el id en la lista de particiones montadas para retonar el path y nombre de la particion y retorna un true o false
//segun sea el caso
func existeID(id string) (string, string, bool) {
	existe := false
	var path string
	var nombre string
	if len(ContenedorMount) == 0 {
		fmt.Println(red + "[ERROR]" + reset + "No existen particiones montadas")
	} else {
		for i := 0; i < len(ContenedorMount); i++ {
			if ContenedorMount[i].ID == id {
				path = ContenedorMount[i].Path
				nombre = ContenedorMount[i].NombreParticion
				existe = true
				break
			}
		}
	}
	return path, nombre, existe
}

//numeroDeEstructuras = funcion encargada de retornar el numero de estructuras a utilizar en la particion
func numeroDeEstructuras(tamanoDeParticion, tamanoDelSuperBloque, tamanoArbolVitual, tamanoDetalleDirectorio, tamanoInodo, tamanoBloque, bitacora int64) int64 {
	var numEstructuras float64 = 0
	var numerador float64 = (float64(tamanoDeParticion) - (2 * float64(tamanoDelSuperBloque)))
	var denominador float64 = (27 + float64(tamanoArbolVitual) + float64(tamanoDetalleDirectorio) + (5*float64(tamanoInodo) + (20 * float64(tamanoBloque)) + float64(bitacora)))
	numEstructuras = numerador / denominador
	fmt.Println(green+"Numero de Estructuras: "+reset, numEstructuras)
	return int64(numEstructuras)
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

func leerSB(path string, start int64) SUPERBOOT {
	var sb SUPERBOOT
	file, err := os.OpenFile(path, os.O_RDWR, 0755)

	defer file.Close()
	if err != nil {
		fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
	}

	var tamanoEnBytes int = int(binary.Size(sb))

	file.Seek(start, 0)
	data := leerBytes(file, tamanoEnBytes)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &sb)
	if err != nil {
		fmt.Println(err)
	}
	return sb
}

func leerAVD(file *os.File, start int64) ARBOLVIRTUALDIRECTORIO {
	var avdRecuperacion ARBOLVIRTUALDIRECTORIO
	var tamanoEnBytes int = int(binary.Size(avdRecuperacion))

	file.Seek(start, 0)
	data := leerBytes(file, tamanoEnBytes)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &avdRecuperacion)
	if err != nil {
		fmt.Println(err)
	}
	return avdRecuperacion
}

func leerDD(file *os.File, start int64) DETALLEDIRECTORIO {
	var ddRecuperacion DETALLEDIRECTORIO
	var tamanoEnBytes int = int(binary.Size(ddRecuperacion))

	file.Seek(start, 0)
	data := leerBytes(file, tamanoEnBytes)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &ddRecuperacion)
	if err != nil {
		fmt.Println(err)
	}
	return ddRecuperacion
}

func leerTABLAINODO(file *os.File, start int64) TABLAINODO {
	var inodoRecuperacion TABLAINODO
	var tamanoEnBytes int = int(binary.Size(inodoRecuperacion))

	file.Seek(start, 0)
	data := leerBytes(file, tamanoEnBytes)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &inodoRecuperacion)
	if err != nil {
		fmt.Println(err)
	}
	return inodoRecuperacion
}

func leerBLOQUEDATOS(file *os.File, start int64) BLOQUEDATOS {
	var bdRecuperacion BLOQUEDATOS
	var tamanoEnBytes int = int(binary.Size(bdRecuperacion))

	file.Seek(start, 0)
	data := leerBytes(file, tamanoEnBytes)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &bdRecuperacion)
	if err != nil {
		fmt.Println(err)
	}
	return bdRecuperacion
}

func leerBitmapResumen(path string, start int64, sb *SUPERBOOT) string {
	var cadena string
	fmt.Println(path)
	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer file.Close()
	if err != nil {
		fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
	}

	file.Seek(start, 0)
	fmt.Println("start: ", start)

	longitudBitmap := sb.SbArbolVirtualCount
	bitmapArbolVirtualDirectorio := make([]byte, longitudBitmap)

	data := leerBytes(file, int(sb.SbArbolVirtualCount))
	buffer := bytes.NewBuffer(data)
	err2 := binary.Read(buffer, binary.BigEndian, bitmapArbolVirtualDirectorio)
	if err2 != nil {
		fmt.Println(err2)
	}

	contador := 0
	cadena += "| "
	for i := 0; i < len(bitmapArbolVirtualDirectorio); i++ {
		//fmt.Println(strconv.Atoi(string(bitmapArbolVirtualDirectorio[i])))
		caracter, _ := strconv.Atoi(string(bitmapArbolVirtualDirectorio[i]))
		cadena += strconv.Itoa(caracter)
		if contador == 20 {
			cadena += "\n| "
			contador = 0
		} else {
			cadena += " | "
			contador++
		}
	}

	return cadena
}

//leerBitmap = Metodo que devuelve un arreglos de []byte con el contenido del bitmap especificado
func retornarBitmap(file *os.File, start int64, sb SUPERBOOT) []byte {
	file.Seek(start, 0)
	longitudBitmap := sb.SbArbolVirtualCount
	bitmapResultado := make([]byte, longitudBitmap)

	data := leerBytes(file, int(sb.SbArbolVirtualCount))
	buffer := bytes.NewBuffer(data)
	err2 := binary.Read(buffer, binary.BigEndian, bitmapResultado)
	if err2 != nil {
		fmt.Println(err2)
	}
	return bitmapResultado
}

func reescribirBitmap(file *os.File, start int64, arreglo []byte) {
	file.Seek(start, 0)
	var valorBinario bytes.Buffer
	binary.Write(&valorBinario, binary.BigEndian, &arreglo)
	escribirBytes(file, valorBinario.Bytes())
}

//SplitSubN = Metodo encargado de separa un string en subcadens
func SplitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}
	return subs
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
	SbArbolVirtualFree              int64
	SbDetalleDirectorioFree         int64
	SbInodosFree                    int64
	SbBloquesFree                   int64
	SbDateCreacion                  [19]byte
	SbDateUltimoMontaje             [19]byte
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

//ARBOLVIRTUALDIRECTORIO es el struct que contiene el arbol virtual de directorio
type ARBOLVIRTUALDIRECTORIO struct {
	AvdFechaCreacion            [19]byte
	AvdNombreDirectorio         [16]byte
	AvdApArraySubdirectorios    [6]int64
	AvdApDetalleDirectorio      int64
	AvdApArbolVirtualDirectorio int64
	AvdProper                   [10]byte
}

//DETALLEDIRECTORIO =
type DETALLEDIRECTORIO struct {
	DdArrayFiles          [5]SUBDETALLEDIRECTORIO
	DdApDetalleDirectorio int64
}

//SUBDETALLEDIRECTORIO =
type SUBDETALLEDIRECTORIO struct {
	DdFileNombre           [16]byte
	DdFileApInodo          int64
	DdFileDateCreacion     [19]byte
	DdFileDateModificacion [19]byte
}

//TABLAINODO =
type TABLAINODO struct {
	ICountInodo            int64
	ISizeArchivo           int64
	ICountBloquesAsignados int64
	IArrayBloques          [4]int64
	IApIndirecto           int64
	IIdProper              [10]byte
}

//BLOQUEDATOS =
type BLOQUEDATOS struct {
	DbDato [25]byte
}

//BITACORA =
type BITACORA struct {
	LogTipoOperacion string
	LogTipo          byte
	LogNombre        [16]byte
	LogContenido     byte
	LogFecha         [20]byte
}

//SESIONACTIVA =
type SESIONACTIVA struct {
	usuario string
}

//SesionActiva = Es un struct GLOBAL que almacena la sesion activa, para saber si en realidad
var SesionActiva SESIONACTIVA

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

	/*
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
	*/
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

//FormateLWH =
func FormateLWH(id, tipoFormateo string) {
	path, nameParticion, Existe := existeID(id)
	var sb SUPERBOOT
	var tamParticion int64 = 0
	var nombreAByte16 [16]byte
	var start int64

	if Existe == true {
		mbr := leerMBR(path)
		copy(nombreAByte16[:], nameParticion)
		if mbr.MbrPartition1.PartName == nombreAByte16 {
			tamParticion = mbr.MbrPartition1.PartSize
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition2.PartName == nombreAByte16 {
			tamParticion = mbr.MbrPartition2.PartSize
			start = mbr.MbrPartition2.PartStart
		} else if mbr.MbrPartition3.PartName == nombreAByte16 {
			tamParticion = mbr.MbrPartition3.PartSize
			start = mbr.MbrPartition3.PartStart
		} else if mbr.MbrPartition4.PartName == nombreAByte16 {
			tamParticion = mbr.MbrPartition4.PartSize
			start = mbr.MbrPartition4.PartStart
		}

		//********************************************************
		//Se abre el Archivo
		//********************************************************
		file, err := os.OpenFile(path, os.O_RDWR, 0755)
		defer file.Close()
		if err != nil {
			fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
		}

		//********************************************************
		//Se realiza el formateo si es de tipo FULL
		//********************************************************
		if tipoFormateo == "FULL" {
			//Se reescribe con 0's toda la particion
			file.Seek(start, 0)
			for i := 0; i < int(tamParticion); i++ {
				numero := byte(0)
				var valorBinario bytes.Buffer
				binary.Write(&valorBinario, binary.BigEndian, &numero)
				escribirBytes(file, valorBinario.Bytes())
			}
		}

		//********************************************************
		//Se obtienen los tamaños de los Struct
		//********************************************************
		fmt.Println(cyan+"Tamaño Particion: "+reset, tamParticion)
		tamSuperBoot := int64(binary.Size(SUPERBOOT{}))
		fmt.Println("Tamaño Super Boot: ", tamSuperBoot)
		tamArbolVirtualDirectorio := int64(unsafe.Sizeof(ARBOLVIRTUALDIRECTORIO{}))
		fmt.Println("Tamaño Arbol Virtual: ", tamArbolVirtualDirectorio)
		tamDetalleDirectorio := int64(binary.Size(DETALLEDIRECTORIO{}))
		fmt.Println("Tamaño Detalle Directorio: ", tamDetalleDirectorio)
		tamTablaInodo := int64(binary.Size(TABLAINODO{}))
		fmt.Println("Tamaño Tabla Inodo: ", tamTablaInodo)
		tamBloqueDatos := int64(binary.Size(BLOQUEDATOS{}))
		fmt.Println("Tamaño Bloque Datos: ", tamBloqueDatos)
		tamBitacora := int64(unsafe.Sizeof(BITACORA{}))
		fmt.Println("Tamaño Bitacora: ", tamBitacora)
		numEstructuras := numeroDeEstructuras(tamParticion, tamSuperBoot, tamArbolVirtualDirectorio, tamDetalleDirectorio, tamTablaInodo, tamBloqueDatos, tamBitacora)
		fmt.Println(red+"Numero de Estructuras: "+reset, numEstructuras)

		if numEstructuras == 0 {
			fmt.Println(red + "[ERROR]" + reset + "NO HAY ESPACIO SUFICIENTE PARA CREAR EL SISTEMA DE ARCHIVOS LWH")
			return
		}

		bitmapArbolVirtualDirectorio := make([]byte, numEstructuras)
		bitmapDetalleDirectorio := make([]byte, numEstructuras)
		bitmapTablaInodos := make([]byte, numEstructuras)
		bitmapBloqueDatos := make([]byte, numEstructuras)
		for i := 0; i < int(numEstructuras); i++ {
			bitmapArbolVirtualDirectorio[i] = '0'
			bitmapDetalleDirectorio[i] = '0'
			bitmapTablaInodos[i] = '0'
			bitmapBloqueDatos[i] = '0'
		}

		fmt.Println("Bitmap Arbol Virtual: ", binary.Size(bitmapArbolVirtualDirectorio))
		fmt.Println("Bitmap Detalle Directorio: ", binary.Size(bitmapDetalleDirectorio))
		fmt.Println("Bitmap Tabla Inodos: ", binary.Size(bitmapTablaInodos))
		fmt.Println("Bitmap Bloque Datos: ", binary.Size(bitmapBloqueDatos))

		//Se obtiene el nombre del disco
		_, nombreArchivo := filepath.Split(path)
		var nombreDiscoAByte16 [16]byte
		copy(nombreDiscoAByte16[:], nombreArchivo)
		//Se obtiene la fecha
		tiempo := time.Now()
		formatoTiempo := tiempo.Format("01-02-2006 15:04:05")
		//tiempoTemp := []byte(formatoTiempo)
		var conversorTiempo [19]byte
		copy(conversorTiempo[:], formatoTiempo)

		//********************************************************
		//Se llena el Super Boot
		//********************************************************
		sb.SbNombreHd = nombreDiscoAByte16
		sb.SbArbolVirtualCount = numEstructuras
		sb.SbDetalleDirectorioCount = numEstructuras
		sb.SbInodosCount = numEstructuras
		sb.SbBloquesCount = numEstructuras
		sb.SbArbolVirtualFree = numEstructuras
		sb.SbDetalleDirectorioFree = numEstructuras
		sb.SbInodosFree = numEstructuras
		sb.SbBloquesFree = numEstructuras
		sb.SbDateCreacion = conversorTiempo
		sb.SbDateUltimoMontaje = conversorTiempo
		sb.SbMontajeCount = 0
		sb.SbApBitmapArbolDirectorio = start + tamSuperBoot
		sb.SbApArbolDirectorio = sb.SbApBitmapArbolDirectorio + int64(binary.Size(bitmapArbolVirtualDirectorio))
		sb.SbApBitmapDetalleDirectorio = sb.SbApArbolDirectorio + (tamArbolVirtualDirectorio * numEstructuras)
		sb.SbApDetalleDirectorio = sb.SbApBitmapDetalleDirectorio + int64(binary.Size(bitmapDetalleDirectorio))
		sb.SbApBitmapTablaInodo = sb.SbApDetalleDirectorio + (tamDetalleDirectorio * numEstructuras)
		sb.SbApTablaInodo = sb.SbApBitmapTablaInodo + int64(binary.Size(bitmapTablaInodos))
		sb.SbApBitmapBloques = sb.SbApTablaInodo + (tamTablaInodo * numEstructuras)
		sb.SbApBloques = sb.SbApBitmapBloques + int64(binary.Size(bitmapBloqueDatos))
		sb.SbApLog = sb.SbApBloques + (tamBloqueDatos * numEstructuras)
		sb.SbSizeStructArbolDirectorio = tamArbolVirtualDirectorio
		sb.SbSizeStructDetalleDirectorio = tamDetalleDirectorio
		sb.SbSizeStructInodo = tamTablaInodo
		sb.SbSizeStructBloque = tamBloqueDatos
		sb.SbFirstFreeBitArbolDirectorio = 0
		sb.SbFirstFreeBitDetalleDirectorio = 0
		sb.SbFirstFreeBitTablaInodo = 0
		sb.SbFirstFreeBitBloques = 0
		sb.SbMagicNum = 201602694
		//********************************************************
		//Se escribe el Super Boot en el disco
		//********************************************************
		file.Seek(start, 0)
		var binarioSuperBoot bytes.Buffer
		binary.Write(&binarioSuperBoot, binary.BigEndian, &sb)
		escribirBytes(file, binarioSuperBoot.Bytes())
		//********************************************************
		//Se escribe el bitmap de Arbol de Directorio
		//********************************************************
		var numero0 byte
		numero0 = '0'

		file.Seek(sb.SbApBitmapArbolDirectorio, 0)
		for i := 0; i < int(numEstructuras); i++ {
			var valorBinario bytes.Buffer
			binary.Write(&valorBinario, binary.BigEndian, &numero0)
			escribirBytes(file, valorBinario.Bytes())
		}
		//********************************************************
		//Se escribe el bitmap de Detalle de Directorio
		//********************************************************
		file.Seek(sb.SbApBitmapDetalleDirectorio, 0)
		for i := 0; i < int(numEstructuras); i++ {
			var valorBinario bytes.Buffer
			binary.Write(&valorBinario, binary.BigEndian, &numero0)
			escribirBytes(file, valorBinario.Bytes())
		}
		//********************************************************
		//Se escribe el bitmap Tabla Inodo
		//********************************************************
		file.Seek(sb.SbApBitmapTablaInodo, 0)
		for i := 0; i < int(numEstructuras); i++ {
			var valorBinario bytes.Buffer
			binary.Write(&valorBinario, binary.BigEndian, &numero0)
			escribirBytes(file, valorBinario.Bytes())
		}
		//********************************************************
		//Se escribe el bitmap de Bloques
		//********************************************************
		file.Seek(sb.SbApBitmapBloques, 0)
		for i := 0; i < int(numEstructuras); i++ {
			var valorBinario bytes.Buffer
			binary.Write(&valorBinario, binary.BigEndian, &numero0)
			escribirBytes(file, valorBinario.Bytes())
		}

		//********************************************************
		//Se prepara el arbol con el directorio raiz '/'
		//Se escribe la raiz '/' y se actualiza el BITMAP
		//********************************************************
		var avd ARBOLVIRTUALDIRECTORIO
		avd.AvdFechaCreacion = conversorTiempo
		var nombreDirectorioRaiz [16]byte
		nombreDirectorioRaiz[0] = '/'
		avd.AvdNombreDirectorio = nombreDirectorioRaiz
		avd.AvdApArraySubdirectorios[0] = 0
		avd.AvdApArraySubdirectorios[1] = 0
		avd.AvdApArraySubdirectorios[2] = 0
		avd.AvdApArraySubdirectorios[3] = 0
		avd.AvdApArraySubdirectorios[4] = 0
		avd.AvdApArraySubdirectorios[5] = 0
		avd.AvdApDetalleDirectorio = 1
		avd.AvdApArbolVirtualDirectorio = 0
		avd.AvdProper[0] = 'r'
		avd.AvdProper[1] = 'o'
		avd.AvdProper[2] = 'o'
		avd.AvdProper[3] = 't'
		//Se posiciona en el inicio del Arbol Directorio, ya que como es la carpeta raiz, va en la primer posicion
		file.Seek(sb.SbApArbolDirectorio, 0)
		var valorBinarioArbolDirectorio bytes.Buffer
		binary.Write(&valorBinarioArbolDirectorio, binary.BigEndian, &avd)
		escribirBytes(file, valorBinarioArbolDirectorio.Bytes())
		//Se procese a actualizar el bitmap
		bitmapArbolVirtualDirectorio[0] = '1'
		reescribirBitmap(file, sb.SbApBitmapArbolDirectorio, bitmapArbolVirtualDirectorio)
		//********************************************************
		//Se prepara el DETALLE DIRECTORIO para user.txt
		//Se escribe el DD y se actualiza el BITMAP
		//********************************************************
		var dd DETALLEDIRECTORIO
		var nombreArchivoUserTXT [16]byte
		usertxt := "user.txt"
		var nombreUserTXTAByte16 [16]byte
		copy(nombreUserTXTAByte16[:], usertxt)
		nombreArchivoUserTXT = nombreUserTXTAByte16
		dd.DdArrayFiles[0].DdFileNombre = nombreArchivoUserTXT
		dd.DdArrayFiles[0].DdFileApInodo = 1
		dd.DdArrayFiles[0].DdFileDateCreacion = conversorTiempo
		dd.DdArrayFiles[0].DdFileDateModificacion = conversorTiempo
		dd.DdArrayFiles[1].DdFileApInodo = 0
		dd.DdArrayFiles[2].DdFileApInodo = 0
		dd.DdArrayFiles[3].DdFileApInodo = 0
		dd.DdArrayFiles[4].DdFileApInodo = 0
		dd.DdApDetalleDirectorio = 0
		//Se posiciona en el inicio del Detalle Directorio, ya que como es el archivo user.txt, va en la primer posicion
		file.Seek(sb.SbApDetalleDirectorio, 0)
		var valorBinarioDetalleDirectorio bytes.Buffer
		binary.Write(&valorBinarioDetalleDirectorio, binary.BigEndian, &dd)
		escribirBytes(file, valorBinarioDetalleDirectorio.Bytes())
		//Se procese a actualizar el bitmap
		bitmapDetalleDirectorio[0] = '1'
		reescribirBitmap(file, sb.SbApBitmapDetalleDirectorio, bitmapDetalleDirectorio)
		//********************************************************
		//Se prepara el INODO para user.txt
		//Se escribe el INODO y se actualiza el BITMAP
		//********************************************************
		var inodo TABLAINODO
		inodo.ICountInodo = 1
		inodo.ISizeArchivo = tamBloqueDatos * 2
		inodo.ICountBloquesAsignados = 2 //Se coloca 2 ya que la cadena inicial es mayor a 25 caracteres pero menor a 50
		inodo.IArrayBloques[0] = 1
		inodo.IArrayBloques[1] = 2
		inodo.IArrayBloques[2] = 0
		inodo.IArrayBloques[3] = 0
		inodo.IApIndirecto = 0
		inodo.IIdProper[0] = 'r'
		inodo.IIdProper[1] = 'o'
		inodo.IIdProper[2] = 'o'
		inodo.IIdProper[3] = 't'
		//Se posiciona en el inicio del Detalle Directorio, ya que como es el archivo user.txt, va en la primer posicion
		file.Seek(sb.SbApTablaInodo, 0)
		var valorBinarioTablaInodo bytes.Buffer
		binary.Write(&valorBinarioTablaInodo, binary.BigEndian, &inodo)
		escribirBytes(file, valorBinarioTablaInodo.Bytes())
		//Se procese a actualizar el bitmap
		bitmapTablaInodos[0] = '1'
		reescribirBitmap(file, sb.SbApBitmapTablaInodo, bitmapTablaInodos)
		//********************************************************
		//Se prepara el BLOQUE del archivo user.txt
		//Se escribe el user.txt y se actualiza el BITMAP
		//********************************************************
		var bloque BLOQUEDATOS
		contLongitud := 0
		var contCuantoBloquesVan int64 = 0
		cadenaUserTxt := "1,G,root\\n1,U,root,root,201602694\\n"
		for i := 0; i < len(cadenaUserTxt); i++ {
			bloque.DbDato[contLongitud] = cadenaUserTxt[i]
			if contLongitud == 24 {
				//Se escribe el bloque lleno
				fmt.Println(sb.SbApBloques + (tamBloqueDatos * contCuantoBloquesVan))
				file.Seek(sb.SbApBloques+(tamBloqueDatos*contCuantoBloquesVan), 0)
				var valorBinarioBloque bytes.Buffer
				binary.Write(&valorBinarioBloque, binary.BigEndian, &bloque)
				escribirBytes(file, valorBinarioBloque.Bytes())
				//Se reinicia el bloque ya que ya se agrego el completo
				var reinicio BLOQUEDATOS
				bloque = reinicio
				contLongitud = -1
				contCuantoBloquesVan++
			} else if i == len(cadenaUserTxt)-1 {
				//Se escribe el bloque lleno
				file.Seek(sb.SbApBloques+(tamBloqueDatos*contCuantoBloquesVan), 0)
				var valorBinarioBloque bytes.Buffer
				binary.Write(&valorBinarioBloque, binary.BigEndian, &bloque)
				escribirBytes(file, valorBinarioBloque.Bytes())
			}
			contLongitud++
		}
		bitmapBloqueDatos[0] = '1'
		bitmapBloqueDatos[1] = '1'
		reescribirBitmap(file, sb.SbApBitmapBloques, bitmapBloqueDatos)
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//LOGIN-----LOGIN-----METODOS-----METODOS-----LOGIN-----LOGIN-----METODOS-----METODOS-----LOGIN-----LOGIN-----METODOS-----METODOS-----LOGIN--
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//Login = metodo que verifica la existencia del usuario y contraseña en el archivo user.txt
func Login(usr, pwd, id string) {
	if SesionActiva.usuario != "" {
		fmt.Println(red + "[ERROR]" + reset + "Ya se encuentra una sesion activa!!")
	} else {
		pathParticion, nameParticion, Existe := existeID(id)
		var nombreAByte16 [16]byte
		var start int64

		if Existe == true {
			mbr := leerMBR(pathParticion)
			copy(nombreAByte16[:], nameParticion)
			if mbr.MbrPartition1.PartName == nombreAByte16 {
				start = mbr.MbrPartition1.PartStart
			} else if mbr.MbrPartition2.PartName == nombreAByte16 {
				start = mbr.MbrPartition2.PartStart
			} else if mbr.MbrPartition3.PartName == nombreAByte16 {
				start = mbr.MbrPartition3.PartStart
			} else if mbr.MbrPartition4.PartName == nombreAByte16 {
				start = mbr.MbrPartition4.PartStart
			}
			//********************************************************
			//Se abre el Archivo
			//********************************************************
			file, err := os.OpenFile(pathParticion, os.O_RDWR, 0755)
			defer file.Close()
			if err != nil {
				fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
			}
			//********************************************************
			//Se retorna el contenido del archivo users.txt
			//********************************************************
			sb := leerSB(pathParticion, start)
			numEstructuraTreeComplete = 0
			numEstructuraTreeCompleteMov = 0
			CadenaRetornoUserTXT = ""
			recorrerArbolRecursivoRetornarUsersTxt(file, sb, 3)
			fmt.Println(cyan + "Cadena USERS.TXT: " + reset + CadenaRetornoUserTXT)
			//********************************************************
			//Se separa el contenido (SPLIT)
			//********************************************************
			cadenaDivididaSlashN := strings.SplitN(CadenaRetornoUserTXT, "\\n", -1)
			tamCadena := len(cadenaDivididaSlashN)
			for i := 0; i < tamCadena; i++ {
				cadenaDivididaComas := strings.SplitN(cadenaDivididaSlashN[i], ",", -1)
				if cadenaDivididaComas[0] == "1" && cadenaDivididaComas[1] == "U" {
					if cadenaDivididaComas[3] == usr && cadenaDivididaComas[4] == pwd {
						SesionActiva.usuario = usr
						fmt.Println(green + "[EXITO]" + reset + "Sesion iniciada con exito!")
					} else if cadenaDivididaComas[3] == usr && cadenaDivididaComas[4] != pwd {
						fmt.Println(red + "[ERROR]" + reset + "Contraseña incorrecta")
					}
				}
			}
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//LOGOUT-----LOGOUT-----METODOS-----METODOS-----LOGOUT-----LOGOUT-----METODOS-----METODOS-----LOGOUT-----LOGOUT-----METODOS-----METODOS-----L
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//FuncionLOGOUT = funcion encargada de salir de la sesion
func FuncionLOGOUT() {
	if SesionActiva.usuario == "" {
		fmt.Println(red + "[ERROR]" + reset + "No se encuentra ninguna sesion activa")
	} else {
		SesionActiva.usuario = ""
		fmt.Println(green + "[EXITO]" + reset + "Sesion cerrada con exito!")
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKGRP-----MKGRP-----METODOS-----METODOS-----MKGRP-----MKGRP-----METODOS-----METODOS-----MKGRP-----MKGRP-----METODOS-----METODOS-----MKGRP--
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//CrearGrupo = metodo encargado de crear un nuevo grupo en el archivo user.txt
func CrearGrupo(id, name string) {
	if SesionActiva.usuario == "" {
		fmt.Println(red + "[ERROR]" + reset + "No se encuentra ninguna sesion activa")
	} else if SesionActiva.usuario != "root" {
		fmt.Println(red + "[ERROR]" + reset + "La funcion " + cyan + "MKGRP" + reset + " unicamente puede ser ejecutada por el usuario root")
	} else if SesionActiva.usuario == "root" {
		pathParticion, nameParticion, Existe := existeID(id)
		var nombreAByte16 [16]byte
		var start int64

		if Existe == true {
			mbr := leerMBR(pathParticion)
			copy(nombreAByte16[:], nameParticion)
			if mbr.MbrPartition1.PartName == nombreAByte16 {
				start = mbr.MbrPartition1.PartStart
			} else if mbr.MbrPartition2.PartName == nombreAByte16 {
				start = mbr.MbrPartition2.PartStart
			} else if mbr.MbrPartition3.PartName == nombreAByte16 {
				start = mbr.MbrPartition3.PartStart
			} else if mbr.MbrPartition4.PartName == nombreAByte16 {
				start = mbr.MbrPartition4.PartStart
			}
			//********************************************************
			//Se abre el Archivo
			//********************************************************
			file, err := os.OpenFile(pathParticion, os.O_RDWR, 0755)
			defer file.Close()
			if err != nil {
				fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
			}
			//********************************************************
			//Se retorna el contenido del archivo users.txt
			//********************************************************
			sb := leerSB(pathParticion, start)
			numEstructuraTreeComplete = 0
			CadenaRetornoUserTXT = ""
			recorrerArbolRecursivoRetornarUsersTxt(file, sb, 3)
			fmt.Println(cyan + "Cadena USERS.TXT: " + reset + CadenaRetornoUserTXT)
			//********************************************************
			//Se separa el contenido (SPLIT) y se verifica si el
			//usuario ya existe
			//********************************************************
			cadenaDivididaSlashN := strings.SplitN(CadenaRetornoUserTXT, "\\n", -1)
			tamCadena := len(cadenaDivididaSlashN)
			yaExisteGrupo := false
			for i := 0; i < tamCadena; i++ {
				cadenaDivididaComas := strings.SplitN(cadenaDivididaSlashN[i], ",", -1)
				if cadenaDivididaComas[0] == "1" && cadenaDivididaComas[1] == "G" {
					if cadenaDivididaComas[2] == name {
						fmt.Println(red + "[ERROR]" + reset + "El grupo ya existe")
						yaExisteGrupo = true
						break
					}
				}
			}
			//********************************************************
			//Se verifica primero si el nombre no supera los 10 byte
			//********************************************************
			if len(name) > 10 {
				fmt.Println(red + "[ERROR]" + reset + "El nombre del grupo supera los 10 caracteres")
				yaExisteGrupo = true
			}
			//********************************************************
			//Si no existe el usuario, se procede a crearlo
			//********************************************************
			if yaExisteGrupo == false {
				//	PARAMETRO 1 -> file: 		recibe el archivo
				//	PARAMETRO 2 -> sb: 			una estructura SUPERBOOT
				//	PARAMETRO 3 -> name:		recibe el nombre del GRUPO PARA USARLO EN GRUPO O USUARIO
				//	PARAMETRO 3 -> usr: 		recibe el nombre del USUARIO PARA USARLO EN USUARIO
				//	PARAMETRO 4 -> pwd: 		recibe el PASSWORD para el uso de la creacion del usuario
				//	PARAMETRO 5 -> tamCadena:	variable int que en princio lleva el len de el arreglo que separa por '\n'
				//  PARAMETRO 6 -> tipo:		1 indica que es grupo, 2 indica que es usuario, 3 y 4 indican que son remover
				modificarUSERTXT(file, sb, name, "", "", tamCadena, 1)
			}
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//RMGRP-----RMGRP-----METODOS-----METODOS-----RMGRP-----RMGRP-----METODOS-----METODOS-----RMGRP-----RMGRP-----METODOS-----METODOS-----RMGRP--
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//RemoverGrupo = metodo que se encarga de eliminar un grupo
func RemoverGrupo(id, name string) {
	if SesionActiva.usuario == "" {
		fmt.Println(red + "[ERROR]" + reset + "No se encuentra ninguna sesion activa")
	} else if SesionActiva.usuario != "root" {
		fmt.Println(red + "[ERROR]" + reset + "La funcion " + cyan + "RMGRP" + reset + " unicamente puede ser ejecutada por el usuario root")
	} else if SesionActiva.usuario == "root" {
		pathParticion, nameParticion, Existe := existeID(id)
		var nombreAByte16 [16]byte
		var start int64

		if Existe == true {
			mbr := leerMBR(pathParticion)
			copy(nombreAByte16[:], nameParticion)
			if mbr.MbrPartition1.PartName == nombreAByte16 {
				start = mbr.MbrPartition1.PartStart
			} else if mbr.MbrPartition2.PartName == nombreAByte16 {
				start = mbr.MbrPartition2.PartStart
			} else if mbr.MbrPartition3.PartName == nombreAByte16 {
				start = mbr.MbrPartition3.PartStart
			} else if mbr.MbrPartition4.PartName == nombreAByte16 {
				start = mbr.MbrPartition4.PartStart
			}
			//********************************************************
			//Se abre el Archivo
			//********************************************************
			file, err := os.OpenFile(pathParticion, os.O_RDWR, 0755)
			defer file.Close()
			if err != nil {
				fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
			}
			//********************************************************
			//Se retorna el contenido del archivo users.txt
			//********************************************************
			sb := leerSB(pathParticion, start)
			numEstructuraTreeComplete = 0
			CadenaRetornoUserTXT = ""
			recorrerArbolRecursivoRetornarUsersTxt(file, sb, 3)
			fmt.Println(cyan + "Cadena USERS.TXT: " + reset + CadenaRetornoUserTXT)
			//********************************************************
			//Se separa el contenido (SPLIT) y se verifica si el
			//grupo ya existe
			//********************************************************
			cadenaDivididaSlashN := strings.SplitN(CadenaRetornoUserTXT, "\\n", -1)
			cadenaReestructurada := ""
			tamCadena := len(cadenaDivididaSlashN)
			yaExisteGrupo := false
			for i := 0; i < tamCadena; i++ {
				cadenaDivididaComas := strings.SplitN(cadenaDivididaSlashN[i], ",", -1)
				if len(cadenaDivididaComas) > 1 {
					if cadenaDivididaComas[1] == "G" {
						if cadenaDivididaComas[0] == "1" && cadenaDivididaComas[2] == name {
							cadenaReestructurada += "0," + cadenaDivididaComas[1] + "," + cadenaDivididaComas[2] + "\\n"
							yaExisteGrupo = true
						} else {
							cadenaReestructurada += cadenaDivididaComas[0] + "," + cadenaDivididaComas[1] + "," + cadenaDivididaComas[2] + "\\n"
						}
					} else if cadenaDivididaComas[1] == "U" {
						cadenaReestructurada += cadenaDivididaComas[0] + "," + cadenaDivididaComas[1] + "," + cadenaDivididaComas[2] + "," + cadenaDivididaComas[3] + "," + cadenaDivididaComas[4] + "\\n"
					}
				}
			}
			fmt.Println(cyan + "Cadena REESTRUCTURADA: " + reset + cadenaReestructurada)
			fmt.Println(len(CadenaRetornoUserTXT))
			fmt.Println(len(cadenaReestructurada))
			CadenaRetornoUserTXT = cadenaReestructurada
			//********************************************************
			//Se verifica si el grupo no existe, para dar error
			//********************************************************
			if yaExisteGrupo == false {
				fmt.Println(red + "[ERROR]" + reset + "El grupo " + cyan + name + reset + " no se encuentra registrado")
			}
			//********************************************************
			//Si todo es correcto, se procede a eliminar el grupo
			//********************************************************
			if yaExisteGrupo == true {
				//	PARAMETRO 1 -> file: 		recibe el archivo
				//	PARAMETRO 2 -> sb: 			una estructura SUPERBOOT
				//	PARAMETRO 3 -> name:		recibe el nombre del GRUPO PARA USARLO EN GRUPO O USUARIO
				//	PARAMETRO 3 -> usr: 		recibe el nombre del USUARIO PARA USARLO EN USUARIO
				//	PARAMETRO 4 -> pwd: 		recibe el PASSWORD para el uso de la creacion del usuario
				//	PARAMETRO 5 -> tamCadena:	variable int que en princio lleva el len de el arreglo que separa por '\n'
				//  PARAMETRO 6 -> tipo:		1 indica que es grupo, 2 indica que es usuario, 3 y 4 indican que son remover
				modificarUSERTXT(file, sb, name, "", "", tamCadena, 3)
			}
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKUSR-----MKUSR-----METODOS-----METODOS-----MKUSR-----MKUSR-----METODOS-----METODOS-----MKUSR-----MKUSR-----METODOS-----METODOS-----MKUSR--
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//CrearUser = metodo encargado de crear un nuevo grupo en el archivo user.txt
func CrearUser(id, usr, pwd, grp string) {
	if SesionActiva.usuario == "" {
		fmt.Println(red + "[ERROR]" + reset + "No se encuentra ninguna sesion activa")
	} else if SesionActiva.usuario != "root" {
		fmt.Println(red + "[ERROR]" + reset + "La funcion " + cyan + "MKUSR" + reset + " unicamente puede ser ejecutada por el usuario root")
	} else if SesionActiva.usuario == "root" {
		pathParticion, nameParticion, Existe := existeID(id)
		var nombreAByte16 [16]byte
		var start int64

		if Existe == true {
			mbr := leerMBR(pathParticion)
			copy(nombreAByte16[:], nameParticion)
			if mbr.MbrPartition1.PartName == nombreAByte16 {
				start = mbr.MbrPartition1.PartStart
			} else if mbr.MbrPartition2.PartName == nombreAByte16 {
				start = mbr.MbrPartition2.PartStart
			} else if mbr.MbrPartition3.PartName == nombreAByte16 {
				start = mbr.MbrPartition3.PartStart
			} else if mbr.MbrPartition4.PartName == nombreAByte16 {
				start = mbr.MbrPartition4.PartStart
			}
			//********************************************************
			//Se abre el Archivo
			//********************************************************
			file, err := os.OpenFile(pathParticion, os.O_RDWR, 0755)
			defer file.Close()
			if err != nil {
				fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
			}
			//********************************************************
			//Se retorna el contenido del archivo users.txt
			//********************************************************
			sb := leerSB(pathParticion, start)
			numEstructuraTreeComplete = 0
			CadenaRetornoUserTXT = ""
			recorrerArbolRecursivoRetornarUsersTxt(file, sb, 3)
			fmt.Println(cyan + "Cadena USERS.TXT: " + reset + CadenaRetornoUserTXT)
			//********************************************************
			//Se separa el contenido (SPLIT) y se verifica si el
			//usuario ya existe
			//********************************************************
			cadenaDivididaSlashN := strings.SplitN(CadenaRetornoUserTXT, "\\n", -1)
			tamCadena := len(cadenaDivididaSlashN)
			yaExisteGrupo := false
			yaExisteUsuario := false
			for i := 0; i < tamCadena; i++ {
				cadenaDivididaComas := strings.SplitN(cadenaDivididaSlashN[i], ",", -1)
				if cadenaDivididaComas[0] == "1" && cadenaDivididaComas[1] == "G" {
					if cadenaDivididaComas[2] == grp {
						yaExisteGrupo = true
						break
					}
				} else if cadenaDivididaComas[0] == "1" && cadenaDivididaComas[1] == "U" {
					if cadenaDivididaComas[3] == usr {
						fmt.Println(red + "[ERROR]" + reset + "El usuario " + cyan + usr + reset + " ya existe")
						yaExisteUsuario = true
					}
				}
			}
			if yaExisteGrupo == false {
				fmt.Println(red + "[ERROR]" + reset + "El grupo " + cyan + grp + reset + " no existe")
			}
			//********************************************************
			//Se verifica primero si el usr y pwd no supera los 10 byte
			//********************************************************
			if len(usr) > 10 {
				fmt.Println(red + "[ERROR]" + reset + "El nombre del usuario supera los 10 caracteres")
				yaExisteGrupo = false
			}
			if len(pwd) > 10 {
				fmt.Println(red + "[ERROR]" + reset + "El password supera los 10 caracteres")
				yaExisteGrupo = false
			}
			//********************************************************
			//Si no existe el usuario, se procede a crearlo
			//********************************************************
			if yaExisteGrupo == true && yaExisteUsuario == false {
				//	PARAMETRO 1 -> file: 		recibe el archivo
				//	PARAMETRO 2 -> sb: 			una estructura SUPERBOOT
				//	PARAMETRO 3 -> name:		recibe el nombre del GRUPO PARA USARLO EN GRUPO O USUARIO
				//	PARAMETRO 3 -> usr: 		recibe el nombre del USUARIO PARA USARLO EN USUARIO
				//	PARAMETRO 4 -> pwd: 		recibe el PASSWORD para el uso de la creacion del usuario
				//	PARAMETRO 5 -> tamCadena:	variable int que en princio lleva el len de el arreglo que separa por '\n'
				//  PARAMETRO 6 -> tipo:		1 indica que es grupo, 2 indica que es usuario, 3 y 4 indican que son remover
				modificarUSERTXT(file, sb, grp, usr, pwd, tamCadena, 2)
			}
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//RMUSR-----RMUSR-----METODOS-----METODOS-----RMUSR-----RMUSR-----METODOS-----METODOS-----RMUSR-----RMUSR-----METODOS-----METODOS-----RMUSR--
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//RemoverUser = metodo encargado de crear un nuevo grupo en el archivo user.txt
func RemoverUser(id, usr string) {
	if SesionActiva.usuario == "" {
		fmt.Println(red + "[ERROR]" + reset + "No se encuentra ninguna sesion activa")
	} else if SesionActiva.usuario != "root" {
		fmt.Println(red + "[ERROR]" + reset + "La funcion " + cyan + "RMUSR" + reset + " unicamente puede ser ejecutada por el usuario root")
	} else if SesionActiva.usuario == "root" {
		pathParticion, nameParticion, Existe := existeID(id)
		var nombreAByte16 [16]byte
		var start int64

		if Existe == true {
			mbr := leerMBR(pathParticion)
			copy(nombreAByte16[:], nameParticion)
			if mbr.MbrPartition1.PartName == nombreAByte16 {
				start = mbr.MbrPartition1.PartStart
			} else if mbr.MbrPartition2.PartName == nombreAByte16 {
				start = mbr.MbrPartition2.PartStart
			} else if mbr.MbrPartition3.PartName == nombreAByte16 {
				start = mbr.MbrPartition3.PartStart
			} else if mbr.MbrPartition4.PartName == nombreAByte16 {
				start = mbr.MbrPartition4.PartStart
			}
			//********************************************************
			//Se abre el Archivo
			//********************************************************
			file, err := os.OpenFile(pathParticion, os.O_RDWR, 0755)
			defer file.Close()
			if err != nil {
				fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
			}
			//********************************************************
			//Se retorna el contenido del archivo users.txt
			//********************************************************
			sb := leerSB(pathParticion, start)
			numEstructuraTreeComplete = 0
			CadenaRetornoUserTXT = ""
			recorrerArbolRecursivoRetornarUsersTxt(file, sb, 3)
			fmt.Println(cyan + "Cadena USERS.TXT: " + reset + CadenaRetornoUserTXT)
			//********************************************************
			//Se separa el contenido (SPLIT) y se verifica si el
			//usuario ya existe
			//********************************************************
			cadenaDivididaSlashN := strings.SplitN(CadenaRetornoUserTXT, "\\n", -1)
			cadenaReestructurada := ""
			tamCadena := len(cadenaDivididaSlashN)
			yaExisteUsuario := false
			for i := 0; i < tamCadena; i++ {
				cadenaDivididaComas := strings.SplitN(cadenaDivididaSlashN[i], ",", -1)
				fmt.Println(cadenaDivididaComas)
				if len(cadenaDivididaComas) > 1 {
					if cadenaDivididaComas[1] == "G" {
						cadenaReestructurada += cadenaDivididaComas[0] + "," + cadenaDivididaComas[1] + "," + cadenaDivididaComas[2] + "\\n"
					} else if cadenaDivididaComas[1] == "U" {
						if cadenaDivididaComas[0] == "1" && cadenaDivididaComas[3] == usr {
							cadenaReestructurada += "0" + "," + cadenaDivididaComas[1] + "," + cadenaDivididaComas[2] + "," + cadenaDivididaComas[3] + "," + cadenaDivididaComas[4] + "\\n"
							yaExisteUsuario = true
						} else {
							cadenaReestructurada += cadenaDivididaComas[0] + "," + cadenaDivididaComas[1] + "," + cadenaDivididaComas[2] + "," + cadenaDivididaComas[3] + "," + cadenaDivididaComas[4] + "\\n"
						}
					}
				}
			}
			CadenaRetornoUserTXT = cadenaReestructurada
			//********************************************************
			//Se verifica si el usuario no existe, para dar error
			//********************************************************
			if yaExisteUsuario == false {
				fmt.Println(red + "[ERROR]" + reset + "El usuario " + cyan + usr + reset + " no se encuentra registrado")
			}
			//********************************************************
			//Si existe el usuario, se procede a removerlo
			//********************************************************
			if yaExisteUsuario == true {
				//	PARAMETRO 1 -> file: 		recibe el archivo
				//	PARAMETRO 2 -> sb: 			una estructura SUPERBOOT
				//	PARAMETRO 3 -> name:		recibe el nombre del GRUPO PARA USARLO EN GRUPO O USUARIO
				//	PARAMETRO 3 -> usr: 		recibe el nombre del USUARIO PARA USARLO EN USUARIO
				//	PARAMETRO 4 -> pwd: 		recibe el PASSWORD para el uso de la creacion del usuario
				//	PARAMETRO 5 -> tamCadena:	variable int que en princio lleva el len de el arreglo que separa por '\n'
				//  PARAMETRO 6 -> tipo:		1 indica que es grupo, 2 indica que es usuario, 3 y 4 indican que son remover
				modificarUSERTXT(file, sb, "", "", "", tamCadena, 4)
			}
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKDIR-----MKDIR-----METODOS-----METODOS-----MKDIR-----MKDIR-----METODOS-----METODOS-----MKDIR-----MKDIR-----METODOS-----METODOS-----MKDIR--
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//CrearDirectorio = metodo que se encarga de crear las carpetas dentro del archivo binario
func CrearDirectorio(id, path string, pActivo bool) {
	if SesionActiva.usuario == "" {
		fmt.Println(red + "[ERROR]" + reset + "No se encuentra ninguna sesion activa")
	} else {
		pathParticion, nameParticion, Existe := existeID(id)
		var nombreAByte16 [16]byte
		var start int64

		if Existe == true {
			mbr := leerMBR(pathParticion)
			copy(nombreAByte16[:], nameParticion)
			if mbr.MbrPartition1.PartName == nombreAByte16 {
				start = mbr.MbrPartition1.PartStart
			} else if mbr.MbrPartition2.PartName == nombreAByte16 {
				start = mbr.MbrPartition2.PartStart
			} else if mbr.MbrPartition3.PartName == nombreAByte16 {
				start = mbr.MbrPartition3.PartStart
			} else if mbr.MbrPartition4.PartName == nombreAByte16 {
				start = mbr.MbrPartition4.PartStart
			}
			//********************************************************
			//Se abre el Archivo
			//********************************************************
			file, err := os.OpenFile(pathParticion, os.O_RDWR, 0755)
			defer file.Close()
			if err != nil {
				fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
			}
			//********************************************************
			//Se retorna el contenido del SuperBoot
			//********************************************************
			sb := leerSB(pathParticion, start)
			//********************************************************
			//Se carga el bitmap del inodo y bloques
			//********************************************************
			bitmapArbolVirtualDirectorio := retornarBitmap(file, sb.SbApBitmapArbolDirectorio, sb)
			bitmapDetalleDirectorio := retornarBitmap(file, sb.SbApBitmapDetalleDirectorio, sb)
			fmt.Println(bitmapArbolVirtualDirectorio)
			fmt.Println(bitmapDetalleDirectorio)

			//********************************************************
			//Se reinician los arreglos de apuntadores, los cuales se
			//encargaran de almacenar los apuntadores que estan en uso
			//para poder reescribirlos luego
			//********************************************************
			ApuntadoresArbolVirtualCarpetasUso = nil
			ApuntadoresDetalleDirectorioCarpetasUso = nil
			ApuntadoresArbolVirtualCarpetasUso = append(ApuntadoresArbolVirtualCarpetasUso, 1)           //Se agrega el 1 ya que este siempre esta
			ApuntadoresDetalleDirectorioCarpetasUso = append(ApuntadoresDetalleDirectorioCarpetasUso, 1) //Se agrega el 1 ya que este siempre esta

			if pActivo == true {
				//********************************************************
				//Se separa el path para obtener las carpetas
				//********************************************************
				cadenaDivididaSlash := strings.SplitN(path, "/", -1)
				for i := range cadenaDivididaSlash {
					if cadenaDivididaSlash[i] == "" {
						cadenaDivididaSlash[i] = "/"
					}
				}

				RArbol = nil
				contadorRuta = 0

				var ra RARBOL
				ra.nivel = 0
				ra.nombre = "/"
				ra.puntero = 0
				RArbol = append(RArbol, ra)

				verificarNivelesRuta(file, sb, cadenaDivididaSlash, false, 0, 1)

				fmt.Println(red + "********************************************" + reset)
				for i := 0; i < len(RArbol); i++ {
					fmt.Println("Nivel: ", RArbol[i].nivel)
					fmt.Println("Nombre: ", RArbol[i].nombre, ", len: ", len(RArbol[i].nombre))
					fmt.Println("Puntero: ", RArbol[i].puntero)
					fmt.Println(cyan + "----------------" + reset)
				}
				fmt.Println(red + "********************************************" + reset)
				cuantoNivelesTieneLaRuta := len(cadenaDivididaSlash)
				cuantosNivelesCumple := 0

				//********************************************************
				//Se verifica cuandos padres existen y se valida para
				//escribir la cantidad necesaria
				//********************************************************
				var ArbolDeCumplimientos []RARBOL
				//fmt.Println(cuantoNivelesTieneLaRuta)
				for i := 0; i < cuantoNivelesTieneLaRuta; i++ {
					seCumplioCondicion := false
					//fmt.Println("valor: ", i, ", Dato: "+cadenaDivididaSlash[i]+", longitud nombre: ", len(cadenaDivididaSlash[i]))
					for j := 0; j < len(RArbol); j++ {
						if RArbol[j].nivel == int64(i) && RArbol[j].nombre == cadenaDivididaSlash[i] {
							//fmt.Println("CUMPLE ", i)
							ArbolDeCumplimientos = append(ArbolDeCumplimientos, RArbol[j])
							seCumplioCondicion = true
							cuantosNivelesCumple++
						}
					}
					if seCumplioCondicion == false {
						break
					}
				}
				fmt.Println(ArbolDeCumplimientos)
				fmt.Println("Niveles de la ruta: ", cuantoNivelesTieneLaRuta)
				fmt.Println("Niveles que se cumplen: ", cuantosNivelesCumple)
				var cuantosNivelesNuevos int = cuantoNivelesTieneLaRuta - cuantosNivelesCumple
				fmt.Println("Niveles nuevo a crear: ", cuantosNivelesNuevos)
				//********************************************************
				//Se verifica las posiciones disponibles en el bitmap
				//para la escritura
				//********************************************************
				posicionesLibresEnBitmapVirtualDirectorio := 0
				posicionesLibresEnBitmapDetalleDirectorio := 0
				for j := 0; j < len(bitmapArbolVirtualDirectorio); j++ {
					if bitmapArbolVirtualDirectorio[j] == '0' {
						posicionesLibresEnBitmapVirtualDirectorio++
					}
				}
				for j := 0; j < len(bitmapDetalleDirectorio); j++ {
					if bitmapDetalleDirectorio[j] == '0' {
						posicionesLibresEnBitmapDetalleDirectorio++
					}
				}
				posicionesDisponiblesArbol := posicionesLibresEnBitmapVirtualDirectorio - cuantosNivelesNuevos
				posicionesDisponiblesDetalle := posicionesLibresEnBitmapVirtualDirectorio - cuantosNivelesNuevos
				posicionAEscribirArbol := cuantosNivelesNuevos
				posicionAEscribirDirectorio := cuantosNivelesNuevos
				numeroRepitencias := cuantosNivelesNuevos
				fmt.Println("Posiciones a escribir: ", posicionAEscribirArbol)
				if posicionesDisponiblesArbol >= 0 && posicionesDisponiblesDetalle >= 0 {
					var todasLasPosicionesAEscribirArbol []int64
					var todasLasPosicionesAEscribirDirectorio []int64
					for i := 0; i < numeroRepitencias; i++ {
						//********************************************************
						//Se cargan todas las posiciones a la que se debe escribir
						//********************************************************
						for j := 0; j < len(bitmapArbolVirtualDirectorio); j++ {
							if posicionAEscribirArbol == 0 {
								break
							}
							if bitmapArbolVirtualDirectorio[j] == '0' {
								todasLasPosicionesAEscribirArbol = append(todasLasPosicionesAEscribirArbol[:], int64(j+1))
								posicionAEscribirArbol--
							}
						}
						for j := 0; j < len(bitmapDetalleDirectorio); j++ {
							if posicionAEscribirDirectorio == 0 {
								break
							}
							if bitmapDetalleDirectorio[j] == '0' {
								todasLasPosicionesAEscribirDirectorio = append(todasLasPosicionesAEscribirDirectorio[:], int64(j+1))
								posicionAEscribirDirectorio--
							}
						}

						fmt.Println(todasLasPosicionesAEscribirArbol)
						fmt.Println(todasLasPosicionesAEscribirDirectorio)

						//********************************************************
						//Se reescribe la informacion del padre, para agregar el
						//puntero del hijo a la posicion correcta
						//********************************************************
						romplerCicloInfinito := false
						posicionDeCumplimiento := ArbolDeCumplimientos[len(ArbolDeCumplimientos)-1].puntero
						for {
							var avdAnterior ARBOLVIRTUALDIRECTORIO
							//posicionUltimoDirectorio := sb.SbApArbolDirectorio + (sb.SbSizeStructArbolDirectorio * (ArbolDeCumplimientos[cuantosNivelesCumple-1].puntero))
							fmt.Println(cyan+"Puntero: ", ArbolDeCumplimientos[len(ArbolDeCumplimientos)-1].puntero, ", ArbolDeCumplimientos", ArbolDeCumplimientos[len(ArbolDeCumplimientos)-1], reset)
							posicionUltimoDirectorio := sb.SbApArbolDirectorio + (sb.SbSizeStructArbolDirectorio * (posicionDeCumplimiento))
							avdAnterior = leerAVD(file, posicionUltimoDirectorio)
							hayDisponibilidadDeEspacio := false
							for k := 0; k < 6; k++ {
								if avdAnterior.AvdApArraySubdirectorios[k] == 0 {
									avdAnterior.AvdApArraySubdirectorios[k] = todasLasPosicionesAEscribirArbol[i]
									hayDisponibilidadDeEspacio = true
									romplerCicloInfinito = true
									break
								}
							}
							apuntadorDeRecursividad := avdAnterior.AvdApArbolVirtualDirectorio
							if hayDisponibilidadDeEspacio == false {
								if apuntadorDeRecursividad == 0 {
									fmt.Println("Voy ENTRANDOOOOOOOOOOOOOOO")
									//********************************************************
									//Se buscan nuevas posiciones en el bitmap para escribir
									//sin interferir con las posiciones que ya estan
									//reservadas para escribir las carpetas
									//********************************************************
									var posicionesAEscribirArbolTemporal []int64
									//var posicionesBloqueadas []int64
									yaEncontreEspacio := false
									bitmapArbolTemporal := bitmapArbolVirtualDirectorio
									fmt.Println("BITMAP TEMPORAL INICIAL: ", bitmapArbolTemporal)
									for j := 0; j < len(bitmapArbolTemporal); j++ {
										if yaEncontreEspacio == true {
											break
										}
										if bitmapArbolTemporal[j] == '0' {
											for k := range todasLasPosicionesAEscribirArbol {
												if todasLasPosicionesAEscribirArbol[k] == int64(j+1) {
													fmt.Println("---------")
													fmt.Println("ESTA POSICION YA ESTA OCUPADA")
													fmt.Println("todasLasPosicionesAEscribirArbol[k]: ", todasLasPosicionesAEscribirArbol[k])
													fmt.Println("int64(j+1): ", int64(j+1))
													fmt.Println("---------")
													bitmapArbolTemporal[j] = '1'
													break
												}
											}
										}
									}
									for j := 0; j < len(bitmapArbolTemporal); j++ {
										if bitmapArbolTemporal[j] == '0' {
											posicionesAEscribirArbolTemporal = append(posicionesAEscribirArbolTemporal[:], int64(j+1))
											break
										}
									}
									fmt.Println("BITMAP TEMPORAL FINAL: ", bitmapArbolTemporal)
									fmt.Println("ARBOL TEMPORAL: ", posicionesAEscribirArbolTemporal)
									avdAnterior.AvdApArbolVirtualDirectorio = posicionesAEscribirArbolTemporal[0]
									file.Seek(posicionUltimoDirectorio, 0)
									var valorBinarioArbolVirtual bytes.Buffer
									binary.Write(&valorBinarioArbolVirtual, binary.BigEndian, &avdAnterior)
									escribirBytes(file, valorBinarioArbolVirtual.Bytes())

									var avdTemporal ARBOLVIRTUALDIRECTORIO
									posicionUltimoDirectorioTemporal := sb.SbApArbolDirectorio + (sb.SbSizeStructArbolDirectorio * (posicionesAEscribirArbolTemporal[i] - 1))
									avdTemporal = leerAVD(file, posicionUltimoDirectorioTemporal)
									avdTemporal.AvdFechaCreacion = avdAnterior.AvdFechaCreacion
									avdTemporal.AvdNombreDirectorio = avdAnterior.AvdNombreDirectorio
									for k := 0; k < 6; k++ {
										if avdTemporal.AvdApArraySubdirectorios[k] == 0 {
											avdTemporal.AvdApArraySubdirectorios[k] = 0
										}
									}
									avdTemporal.AvdApDetalleDirectorio = 0
									avdTemporal.AvdApArbolVirtualDirectorio = 0
									avdTemporal.AvdProper = avdAnterior.AvdProper

									file.Seek(posicionUltimoDirectorioTemporal, 0)
									var valorBinarioArbolVirtualTemporal bytes.Buffer
									binary.Write(&valorBinarioArbolVirtualTemporal, binary.BigEndian, &avdTemporal)
									escribirBytes(file, valorBinarioArbolVirtualTemporal.Bytes())

									bitmapArbolVirtualDirectorio[posicionesAEscribirArbolTemporal[i]-1] = '1'
									file.Seek(sb.SbApBitmapArbolDirectorio, 0)
									var valorBitmapVirtual bytes.Buffer
									binary.Write(&valorBitmapVirtual, binary.BigEndian, &bitmapArbolVirtualDirectorio)
									escribirBytes(file, valorBitmapVirtual.Bytes())
								} else {
									posicionDeCumplimiento = avdAnterior.AvdApArbolVirtualDirectorio - 1
								}
							} else {
								file.Seek(posicionUltimoDirectorio, 0)
								var valorBinarioArbolVirtual bytes.Buffer
								binary.Write(&valorBinarioArbolVirtual, binary.BigEndian, &avdAnterior)
								escribirBytes(file, valorBinarioArbolVirtual.Bytes())

								//********************************************************
								//Se escribe el nuevo AVD
								//********************************************************
								var avdNuevo ARBOLVIRTUALDIRECTORIO
								var nombreAByte16 [16]byte
								copy(nombreAByte16[:], cadenaDivididaSlash[len(cadenaDivididaSlash)-cuantosNivelesNuevos])
								fmt.Println("NOMBRE A ESCRIBIR: ", cadenaDivididaSlash[len(cadenaDivididaSlash)-cuantosNivelesNuevos])

								avdNuevo.AvdNombreDirectorio = nombreAByte16

								fmt.Println(red + "MAMARRE" + reset)

								tiempo := time.Now()
								formatoTiempo := tiempo.Format("01-02-2006 15:04:05")
								var conversorTiempo [19]byte
								copy(conversorTiempo[:], formatoTiempo)

								avdNuevo.AvdFechaCreacion = conversorTiempo
								for k := 0; k < 6; k++ {
									avdNuevo.AvdApArraySubdirectorios[k] = 0
								}
								avdNuevo.AvdApDetalleDirectorio = todasLasPosicionesAEscribirDirectorio[i]
								avdNuevo.AvdApArbolVirtualDirectorio = 0
								copy(avdNuevo.AvdProper[:], SesionActiva.usuario)

								fmt.Println(red, todasLasPosicionesAEscribirArbol[i]-1, reset)
								file.Seek(sb.SbApArbolDirectorio+(sb.SbSizeStructArbolDirectorio*(todasLasPosicionesAEscribirArbol[i]-1)), 0)
								var valorBinarioArbolVirtualNuevo bytes.Buffer
								binary.Write(&valorBinarioArbolVirtualNuevo, binary.BigEndian, &avdNuevo)
								escribirBytes(file, valorBinarioArbolVirtualNuevo.Bytes())

								//********************************************************
								//Se escribe el nuevo DD
								//********************************************************
								var ddNuevo DETALLEDIRECTORIO
								ddNuevo.DdArrayFiles[0].DdFileApInodo = 0
								ddNuevo.DdArrayFiles[1].DdFileApInodo = 0
								ddNuevo.DdArrayFiles[2].DdFileApInodo = 0
								ddNuevo.DdArrayFiles[3].DdFileApInodo = 0
								ddNuevo.DdArrayFiles[4].DdFileApInodo = 0
								ddNuevo.DdApDetalleDirectorio = 0

								file.Seek(sb.SbApDetalleDirectorio+(sb.SbSizeStructDetalleDirectorio*(todasLasPosicionesAEscribirDirectorio[i]-1)), 0)
								var valorBinarioDetalleDirectorio bytes.Buffer
								binary.Write(&valorBinarioDetalleDirectorio, binary.BigEndian, &ddNuevo)
								escribirBytes(file, valorBinarioDetalleDirectorio.Bytes())

								//********************************************************
								//Se reescriben los bitmaps
								//********************************************************
								bitmapArbolVirtualDirectorio[todasLasPosicionesAEscribirArbol[i]-1] = '1'
								bitmapDetalleDirectorio[todasLasPosicionesAEscribirDirectorio[i]-1] = '1'

								file.Seek(sb.SbApBitmapArbolDirectorio, 0)
								var valorBitmapArbol bytes.Buffer
								binary.Write(&valorBitmapArbol, binary.BigEndian, &bitmapArbolVirtualDirectorio)
								escribirBytes(file, valorBitmapArbol.Bytes())

								file.Seek(sb.SbApBitmapDetalleDirectorio, 0)
								var valorBitmapDirectorio bytes.Buffer
								binary.Write(&valorBitmapDirectorio, binary.BigEndian, &bitmapDetalleDirectorio)
								escribirBytes(file, valorBitmapDirectorio.Bytes())

								var rarbol2 RARBOL
								rarbol2.nivel = 50
								rarbol2.nombre = cadenaDivididaSlash[len(cadenaDivididaSlash)-cuantosNivelesNuevos]
								rarbol2.puntero = todasLasPosicionesAEscribirArbol[i] - 1
								ArbolDeCumplimientos = append(ArbolDeCumplimientos, rarbol2)

								cuantosNivelesNuevos--
							}

							if romplerCicloInfinito == true {
								break
							}
						}
					}
				}
			} else {
				//********************************************************
				//Se reinician los arreglos de apuntadores, los cuales se
				//encargaran de almacenar los apuntadores que estan en uso
				//para poder reescribirlos luego
				//********************************************************
				ApuntadoresArbolVirtualCarpetasUso = nil
				ApuntadoresDetalleDirectorioCarpetasUso = nil
				ApuntadoresArbolVirtualCarpetasUso = append(ApuntadoresArbolVirtualCarpetasUso, 1)           //Se agrega el 1 ya que este siempre esta
				ApuntadoresDetalleDirectorioCarpetasUso = append(ApuntadoresDetalleDirectorioCarpetasUso, 1) //Se agrega el 1 ya que este siempre esta

				//********************************************************
				//Se separa el path para obtener las carpetas
				//********************************************************
				cadenaDivididaSlash := strings.SplitN(path, "/", -1)
				for i := range cadenaDivididaSlash {
					if cadenaDivididaSlash[i] == "" {
						cadenaDivididaSlash[i] = "/"
					}
				}

				RArbol = nil
				contadorRuta = 0

				var ra RARBOL
				ra.nivel = 0
				ra.nombre = "/"
				ra.puntero = 0
				RArbol = append(RArbol, ra)

				verificarNivelesRuta(file, sb, cadenaDivididaSlash, false, 0, 1)

				fmt.Println(red + "********************************************" + reset)
				for i := 0; i < len(RArbol); i++ {
					fmt.Println("Nivel: ", RArbol[i].nivel)
					fmt.Println("Nombre: ", RArbol[i].nombre, ", len: ", len(RArbol[i].nombre))
					fmt.Println("Puntero: ", RArbol[i].puntero)
					fmt.Println(cyan + "----------------" + reset)
				}
				fmt.Println(red + "********************************************" + reset)
				cuantoNivelesTieneLaRuta := len(cadenaDivididaSlash)
				cuantosNivelesCumple := 0

				//********************************************************
				//Se verifica cuandos padres existen y se valida para
				//escribir la cantidad necesaria
				//********************************************************
				var ArbolDeCumplimientos []RARBOL
				//fmt.Println(cuantoNivelesTieneLaRuta)
				for i := 0; i < cuantoNivelesTieneLaRuta; i++ {
					seCumplioCondicion := false
					//fmt.Println("valor: ", i, ", Dato: "+cadenaDivididaSlash[i]+", longitud nombre: ", len(cadenaDivididaSlash[i]))
					for j := 0; j < len(RArbol); j++ {
						if RArbol[j].nivel == int64(i) && RArbol[j].nombre == cadenaDivididaSlash[i] {
							//fmt.Println("CUMPLE ", i)
							ArbolDeCumplimientos = append(ArbolDeCumplimientos, RArbol[j])
							seCumplioCondicion = true
							cuantosNivelesCumple++
						}
					}
					if seCumplioCondicion == false {
						break
					}
				}

				fmt.Println(ArbolDeCumplimientos)
				fmt.Println("Niveles de la ruta: ", cuantoNivelesTieneLaRuta)
				fmt.Println("Niveles que se cumplen: ", cuantosNivelesCumple)
				var cuantosNivelesNuevos int = cuantoNivelesTieneLaRuta - cuantosNivelesCumple
				fmt.Println("Niveles nuevo a crear: ", cuantosNivelesNuevos)
				existeElPadre := false
				if cuantosNivelesNuevos > 1 {
					fmt.Println(red + "[ERROR]" + reset + "El padre de la carpeta que desea crear, no existe")
				} else {
					existeElPadre = true
				}
				if existeElPadre == true {
					//********************************************************
					//Se verifica las posiciones disponibles en el bitmap
					//para la escritura
					//********************************************************
					posicionesLibresEnBitmapVirtualDirectorio := 0
					posicionesLibresEnBitmapDetalleDirectorio := 0
					for j := 0; j < len(bitmapArbolVirtualDirectorio); j++ {
						if bitmapArbolVirtualDirectorio[j] == '0' {
							posicionesLibresEnBitmapVirtualDirectorio++
						}
					}
					for j := 0; j < len(bitmapDetalleDirectorio); j++ {
						if bitmapDetalleDirectorio[j] == '0' {
							posicionesLibresEnBitmapDetalleDirectorio++
						}
					}

					posicionesDisponiblesArbol := posicionesLibresEnBitmapVirtualDirectorio - cuantosNivelesNuevos
					posicionesDisponiblesDetalle := posicionesLibresEnBitmapDetalleDirectorio - cuantosNivelesNuevos
					if posicionesDisponiblesArbol >= 0 && posicionesDisponiblesDetalle >= 0 {
						var todasLasPosicionesAEscribirArbol []int64
						var todasLasPosicionesAEscribirDirectorio []int64
						//********************************************************
						//Se cargan todas las posiciones a la que se debe escribir
						//********************************************************
						for j := 0; j < len(bitmapArbolVirtualDirectorio); j++ {
							if bitmapArbolVirtualDirectorio[j] == '0' {
								todasLasPosicionesAEscribirArbol = append(todasLasPosicionesAEscribirArbol[:], int64(j+1))
								break
							}
						}
						for j := 0; j < len(bitmapDetalleDirectorio); j++ {
							if bitmapDetalleDirectorio[j] == '0' {
								todasLasPosicionesAEscribirDirectorio = append(todasLasPosicionesAEscribirDirectorio[:], int64(j+1))
								break
							}
						}

						fmt.Println(todasLasPosicionesAEscribirArbol)
						fmt.Println(todasLasPosicionesAEscribirDirectorio)

						//********************************************************
						//Se reescribe la informacion del padre, para agregar el
						//puntero del hijo a la posicion correcta
						//********************************************************
						romplerCicloInfinito := false
						posicionDeCumplimiento := ArbolDeCumplimientos[len(ArbolDeCumplimientos)-1].puntero
						for {
							var avdAnterior ARBOLVIRTUALDIRECTORIO
							//posicionUltimoDirectorio := sb.SbApArbolDirectorio + (sb.SbSizeStructArbolDirectorio * (ArbolDeCumplimientos[cuantosNivelesCumple-1].puntero))
							posicionUltimoDirectorio := sb.SbApArbolDirectorio + (sb.SbSizeStructArbolDirectorio * (posicionDeCumplimiento))
							fmt.Println(cyan+"Puntero: ", ArbolDeCumplimientos[len(ArbolDeCumplimientos)-1].puntero, ", ArbolDeCumplimientos", ArbolDeCumplimientos[len(ArbolDeCumplimientos)-1], reset)
							avdAnterior = leerAVD(file, posicionUltimoDirectorio)
							hayDisponibilidadDeEspacio := false
							for k := 0; k < 6; k++ {
								if avdAnterior.AvdApArraySubdirectorios[k] == 0 {
									avdAnterior.AvdApArraySubdirectorios[k] = todasLasPosicionesAEscribirArbol[0]
									hayDisponibilidadDeEspacio = true
									romplerCicloInfinito = true
									break
								}
							}
							apuntadorDeRecursividad := avdAnterior.AvdApArbolVirtualDirectorio
							if hayDisponibilidadDeEspacio == false {
								if apuntadorDeRecursividad == 0 {
									//********************************************************
									//Se buscan nuevas posiciones en el bitmap para escribir
									//sin interferir con las posiciones que ya estan
									//reservadas para escribir las carpetas
									//********************************************************
									var posicionesAEscribirArbolTemporal []int64
									yaEncontreEspacio := false
									bitmapArbolTemporal := bitmapArbolVirtualDirectorio
									for j := 0; j < len(bitmapArbolTemporal); j++ {
										if yaEncontreEspacio == true {
											break
										}
										if bitmapArbolTemporal[j] == '0' {
											for k := range todasLasPosicionesAEscribirArbol {
												if todasLasPosicionesAEscribirArbol[k] == int64(j+1) {
													fmt.Println("ESTA POSICION YA ESTA OCUPADA")
												} else {
													posicionesAEscribirArbolTemporal = append(posicionesAEscribirArbolTemporal[:], int64(j+1))
													yaEncontreEspacio = true
													break
												}
											}
										}
									}
									avdAnterior.AvdApArbolVirtualDirectorio = posicionesAEscribirArbolTemporal[0]
									file.Seek(posicionUltimoDirectorio, 0)
									var valorBinarioArbolVirtual bytes.Buffer
									binary.Write(&valorBinarioArbolVirtual, binary.BigEndian, &avdAnterior)
									escribirBytes(file, valorBinarioArbolVirtual.Bytes())

									var avdTemporal ARBOLVIRTUALDIRECTORIO
									posicionUltimoDirectorioTemporal := sb.SbApArbolDirectorio + (sb.SbSizeStructArbolDirectorio * (posicionesAEscribirArbolTemporal[0] - 1))
									avdTemporal = leerAVD(file, posicionUltimoDirectorioTemporal)
									avdTemporal.AvdFechaCreacion = avdAnterior.AvdFechaCreacion
									avdTemporal.AvdNombreDirectorio = avdAnterior.AvdNombreDirectorio
									for k := 0; k < 6; k++ {
										if avdTemporal.AvdApArraySubdirectorios[k] == 0 {
											avdTemporal.AvdApArraySubdirectorios[k] = 0
										}
									}
									avdTemporal.AvdApDetalleDirectorio = 0
									avdTemporal.AvdApArbolVirtualDirectorio = 0
									avdTemporal.AvdProper = avdAnterior.AvdProper

									file.Seek(posicionUltimoDirectorioTemporal, 0)
									var valorBinarioArbolVirtualTemporal bytes.Buffer
									binary.Write(&valorBinarioArbolVirtualTemporal, binary.BigEndian, &avdTemporal)
									escribirBytes(file, valorBinarioArbolVirtualTemporal.Bytes())

									bitmapArbolVirtualDirectorio[posicionesAEscribirArbolTemporal[0]-1] = '1'
									file.Seek(sb.SbApBitmapArbolDirectorio, 0)
									var valorBitmapVirtual bytes.Buffer
									binary.Write(&valorBitmapVirtual, binary.BigEndian, &bitmapArbolVirtualDirectorio)
									escribirBytes(file, valorBitmapVirtual.Bytes())
									//fmt.Println(magenta, posicionesAEscribirArbolTemporal, reset)
								} else {
									posicionDeCumplimiento = avdAnterior.AvdApArbolVirtualDirectorio - 1
								}
							} else {
								file.Seek(posicionUltimoDirectorio, 0)
								var valorBinarioArbolVirtual bytes.Buffer
								binary.Write(&valorBinarioArbolVirtual, binary.BigEndian, &avdAnterior)
								escribirBytes(file, valorBinarioArbolVirtual.Bytes())

								//********************************************************
								//Se escribe el nuevo AVD
								//********************************************************
								var avdNuevo ARBOLVIRTUALDIRECTORIO
								var nombreAByte16 [16]byte
								copy(nombreAByte16[:], cadenaDivididaSlash[len(cadenaDivididaSlash)-cuantosNivelesNuevos])
								fmt.Println("NOMBRE A ESCRIBIR: ", cadenaDivididaSlash[len(cadenaDivididaSlash)-cuantosNivelesNuevos])

								avdNuevo.AvdNombreDirectorio = nombreAByte16

								tiempo := time.Now()
								formatoTiempo := tiempo.Format("01-02-2006 15:04:05")
								var conversorTiempo [19]byte
								copy(conversorTiempo[:], formatoTiempo)

								avdNuevo.AvdFechaCreacion = conversorTiempo
								for k := 0; k < 6; k++ {
									avdNuevo.AvdApArraySubdirectorios[k] = 0
								}
								avdNuevo.AvdApDetalleDirectorio = todasLasPosicionesAEscribirDirectorio[0]
								avdNuevo.AvdApArbolVirtualDirectorio = 0
								copy(avdNuevo.AvdProper[:], SesionActiva.usuario)

								fmt.Println(red, todasLasPosicionesAEscribirArbol[0]-1, reset)
								file.Seek(sb.SbApArbolDirectorio+(sb.SbSizeStructArbolDirectorio*(todasLasPosicionesAEscribirArbol[0]-1)), 0)
								var valorBinarioArbolVirtualNuevo bytes.Buffer
								binary.Write(&valorBinarioArbolVirtualNuevo, binary.BigEndian, &avdNuevo)
								escribirBytes(file, valorBinarioArbolVirtualNuevo.Bytes())

								//********************************************************
								//Se reescriben los bitmaps
								//********************************************************
								bitmapArbolVirtualDirectorio[todasLasPosicionesAEscribirArbol[0]-1] = '1'
								bitmapDetalleDirectorio[todasLasPosicionesAEscribirDirectorio[0]-1] = '1'

								file.Seek(sb.SbApBitmapArbolDirectorio, 0)
								var valorBitmapVirtual bytes.Buffer
								binary.Write(&valorBitmapVirtual, binary.BigEndian, &bitmapArbolVirtualDirectorio)
								escribirBytes(file, valorBitmapVirtual.Bytes())

								file.Seek(sb.SbApBitmapDetalleDirectorio, 0)
								var valorBitmapDirectorio bytes.Buffer
								binary.Write(&valorBitmapDirectorio, binary.BigEndian, &bitmapDetalleDirectorio)
								escribirBytes(file, valorBitmapDirectorio.Bytes())

								var rarbol2 RARBOL
								rarbol2.nivel = 5
								rarbol2.nombre = cadenaDivididaSlash[len(cadenaDivididaSlash)-cuantosNivelesNuevos]
								rarbol2.puntero = todasLasPosicionesAEscribirArbol[0] - 1
								ArbolDeCumplimientos = append(ArbolDeCumplimientos, rarbol2)
							}
							if romplerCicloInfinito == true {
								break
							}
						}
					}
				}
			}
		}
	}
}

//ApuntadoresArbolVirtualCarpetasUso = Arreglo que almacena los apuntadores de tipo Carpeta
var ApuntadoresArbolVirtualCarpetasUso []int64

//ApuntadoresDetalleDirectorioCarpetasUso = Arreglo que almacena los apuntadores de tipo Inodo que usa el archivo txt
var ApuntadoresDetalleDirectorioCarpetasUso []int64

//RESUMENARBOL =
type RESUMENARBOL struct {
	nivel   int64
	nombre  [16]byte
	puntero int64
	hijos   []*RESUMENARBOL
}

//RARBOL =
type RARBOL struct {
	nivel   int64
	nombre  string
	puntero int64
}

//ResumenArbol =
var ResumenArbol []RESUMENARBOL

//RArbol =
var RArbol []RARBOL
var contadorNormal int64
var contadorMov int64

func pruebaRecursividad(file *os.File, sb SUPERBOOT, posicion int64, tipoArchivo int) {
	if tipoArchivo == 1 { //SE UTILIZA PARA RECORRER AL PADRE
		var avd ARBOLVIRTUALDIRECTORIO
		avd = leerAVD(file, sb.SbApArbolDirectorio+(sb.SbSizeStructArbolDirectorio*posicion))
		var cont int = 0
		var posicionXD int64
		posicionXD = posicion
		for i := 0; i < 8; i++ {
			if i < 6 && avd.AvdApArraySubdirectorios[i] != 0 {
				fmt.Println("[TIPO1]pos ", i, ", grado: ", posicionXD, " {"+string(avd.AvdNombreDirectorio[:])+"} -> ", avd.AvdApArraySubdirectorios[i])
				var ra RARBOL
				ra.nivel = posicion
				var nombreP1 string
				for i1, valor1 := range avd.AvdNombreDirectorio {
					if avd.AvdNombreDirectorio[i1] != 0 {
						nombreP1 += string(valor1)
					}
				}
				ra.nombre = nombreP1
				ra.puntero = posicion
				yaExiste := false
				for j := 0; j < len(RArbol); j++ {
					if RArbol[j].nombre == ra.nombre {
						yaExiste = true
					}
				}
				if yaExiste == false {
					RArbol = append(RArbol, ra)
				}
				posicionNueva := avd.AvdApArraySubdirectorios[i]
				pruebaRecursividad(file, sb, posicionNueva-1, 1)
				cont++
			}
		}
		if cont == 0 {
			fmt.Println("[TIPO1]pos  X , grado: ", posicion, " {"+string(avd.AvdNombreDirectorio[:])+"} -> NADA")
			var ra RARBOL
			ra.nivel = posicion
			var nombreP1 string
			for i1, valor1 := range avd.AvdNombreDirectorio {
				if avd.AvdNombreDirectorio[i1] != 0 {
					nombreP1 += string(valor1)
				}
			}
			ra.nombre = nombreP1
			ra.puntero = posicion
			yaExiste := false
			for j := 0; j < len(RArbol); j++ {
				if RArbol[j].nombre == ra.nombre {
					yaExiste = true
				}
			}
			if yaExiste == false {
				RArbol = append(RArbol, ra)
			}
		}
	} else if tipoArchivo == 2 {
	}
}

var contadorRuta int64

func verificarNivelesRuta(file *os.File, sb SUPERBOOT, arregloRutas []string, soyYoMismo bool, posicion int64, tipoArchivo int) {
	if tipoArchivo == 1 { //SE UTILIZA PARA RECORRER AL PADRE
		var avd ARBOLVIRTUALDIRECTORIO
		avd = leerAVD(file, sb.SbApArbolDirectorio+(sb.SbSizeStructArbolDirectorio*posicion))
		if contadorRuta == 0 {
			/*
				var ra RARBOL
				ra.nivel = 0
				var nombreP1 string
				for i1, valor1 := range avd.AvdNombreDirectorio {
					if avd.AvdNombreDirectorio[i1] != 0 {
						nombreP1 += string(valor1)
					}
				}
				ra.nombre = nombreP1
				ra.puntero = posicion
				RArbol = append(RArbol, ra)
			*/
		}
		var cont int = 0
		var posicionXD int64
		posicionXD = posicion
		//fmt.Println("Contador ruta: ", contadorRuta)
		//fmt.Println("Nombre en contador:", arregloRutas[contadorRuta])
		if soyYoMismo == false {
			contadorRuta++
		}
		for i := 0; i < 8; i++ {
			if i < 6 && avd.AvdApArraySubdirectorios[i] != 0 {
				var apuntador int64
				apuntador = avd.AvdApArraySubdirectorios[i] - 1
				avdTemporal := leerAVD(file, sb.SbApArbolDirectorio+(sb.SbSizeStructArbolDirectorio*apuntador))
				var nombreDirectorio string
				for i1, valor1 := range avdTemporal.AvdNombreDirectorio {
					if avdTemporal.AvdNombreDirectorio[i1] != 0 {
						nombreDirectorio += string(valor1)
					}
				}
				if contadorRuta < int64(len(arregloRutas)) {
					fmt.Println(green + arregloRutas[contadorRuta] + reset)
					fmt.Println(yellow + nombreDirectorio + reset)
					if arregloRutas[contadorRuta] == nombreDirectorio {
						fmt.Println("[TIPO1]pos ", i, ", grado: ", posicionXD, " {"+string(avd.AvdNombreDirectorio[:])+"} -> ", avd.AvdApArraySubdirectorios[i])
						var ra RARBOL
						ra.nivel = contadorRuta
						ra.nombre = nombreDirectorio
						ra.puntero = apuntador
						verificarCopia := false
						for j := range RArbol {
							if RArbol[j].nivel == ra.nivel && RArbol[j].nombre == ra.nombre {
								fmt.Println("SOYYYYY COOOOOPIIIIIIAAAA")
								verificarCopia = true
							} else {

							}
						}
						if verificarCopia == false {
							RArbol = append(RArbol, ra)
						}
						//RArbol = append(RArbol, ra)
						verificarNivelesRuta(file, sb, arregloRutas, false, apuntador, 1)
					}
				}
				cont++
			} else if i == 7 && avd.AvdApArbolVirtualDirectorio != 0 {
				fmt.Println("[TIPO1]pos ind", i, ", grado: ", posicionXD, " {"+string(avd.AvdNombreDirectorio[:])+"} -> ", avd.AvdApArbolVirtualDirectorio)
				var apuntador int64
				apuntador = avd.AvdApArbolVirtualDirectorio
				verificarNivelesRuta(file, sb, arregloRutas, true, apuntador-1, 1)
			}
		}
		if cont == 0 && soyYoMismo == false {
			fmt.Println("[TIPO1]pos  X , grado: ", posicion, " {"+string(avd.AvdNombreDirectorio[:])+"} -> NADA")
			var ra RARBOL
			ra.nivel = contadorRuta
			var nombreDirectorio string
			for i1, valor1 := range avd.AvdNombreDirectorio {
				if avd.AvdNombreDirectorio[i1] != 0 {
					nombreDirectorio += string(valor1)
				}
			}
			fmt.Println(nombreDirectorio)
			ra.nombre = nombreDirectorio
			ra.puntero = posicion
			verificarCopia := false
			for i := range RArbol {
				if RArbol[i].nivel == ra.nivel && RArbol[i].nombre == ra.nombre {
					fmt.Println("SOYYYYY COOOOOPIIIIIIAAAA")
					verificarCopia = true
				}
			}
			if verificarCopia == false {
				RArbol = append(RArbol, ra)
			}
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//MKFILE-----MKFILE-----METODOS-----METODOS-----MKFILE-----MKFILE-----METODOS-----METODOS-----MKFILE-----MKFILE-----METODOS-----METODOS-----M
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//CrearArchivo = metodo que se encarga de crear los archivos dentro de las carpetas
func CrearArchivo(id, path string, pActivo bool, size, cont string) {
	if SesionActiva.usuario == "" {
		fmt.Println(red + "[ERROR]" + reset + "No se encuentra ninguna sesion activa")
	} else {
		pathParticion, nameParticion, Existe := existeID(id)
		var nombreAByte16 [16]byte
		var start int64

		if Existe == true {
			mbr := leerMBR(pathParticion)
			copy(nombreAByte16[:], nameParticion)
			if mbr.MbrPartition1.PartName == nombreAByte16 {
				start = mbr.MbrPartition1.PartStart
			} else if mbr.MbrPartition2.PartName == nombreAByte16 {
				start = mbr.MbrPartition2.PartStart
			} else if mbr.MbrPartition3.PartName == nombreAByte16 {
				start = mbr.MbrPartition3.PartStart
			} else if mbr.MbrPartition4.PartName == nombreAByte16 {
				start = mbr.MbrPartition4.PartStart
			}
			//********************************************************
			//Se abre el Archivo
			//********************************************************
			file, err := os.OpenFile(pathParticion, os.O_RDWR, 0755)
			defer file.Close()
			if err != nil {
				fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
			}
			//********************************************************
			//Se retorna el contenido del SuperBoot
			//********************************************************
			sb := leerSB(pathParticion, start)
			//********************************************************
			//Se carga el bitmap del inodo y bloques
			//********************************************************
			bitmapTablaInodo := retornarBitmap(file, sb.SbApBitmapTablaInodo, sb)
			fmt.Println(bitmapTablaInodo)

			//********************************************************
			//Se reinician los arreglos de apuntadores, los cuales se
			//encargaran de almacenar los apuntadores que estan en uso
			//para poder reescribirlos luego
			//********************************************************
			ApuntadoresArbolVirtualCarpetasUso = nil
			ApuntadoresDetalleDirectorioCarpetasUso = nil
			ApuntadoresArbolVirtualCarpetasUso = append(ApuntadoresArbolVirtualCarpetasUso, 1)           //Se agrega el 1 ya que este siempre esta
			ApuntadoresDetalleDirectorioCarpetasUso = append(ApuntadoresDetalleDirectorioCarpetasUso, 1) //Se agrega el 1 ya que este siempre esta

			//********************************************************
			//Se separa el path para obtener las carpetas
			//********************************************************
			cadenaDivididaSlash := strings.SplitN(path, "/", -1)
			for i := range cadenaDivididaSlash {
				if cadenaDivididaSlash[i] == "" {
					cadenaDivididaSlash[i] = "/"
				}
			}

			cadenaDivididaSinArchivo := cadenaDivididaSlash
			ultimoElementoEnCadenaDividida := len(cadenaDivididaSinArchivo) - 1
			copy(cadenaDivididaSinArchivo[ultimoElementoEnCadenaDividida:], cadenaDivididaSinArchivo[ultimoElementoEnCadenaDividida+1:])
			cadenaDivididaSinArchivo[len(cadenaDivididaSinArchivo)-1] = ""
			cadenaDivididaSinArchivo = cadenaDivididaSinArchivo[:len(cadenaDivididaSinArchivo)-1]
			fmt.Println("CADENA DIVIDIDA SIN ARCHIVOS: ", cadenaDivididaSinArchivo)

			//********************************************************
			//Se crea la carpeta de no existir
			//********************************************************
			rutaArchivo, NombreArchivo := filepath.Split(path)
			rutaArchivo = string([]rune(rutaArchivo)[:len(rutaArchivo)-1])
			fmt.Println("Ruta:", rutaArchivo)      // /ruta/hacia/algun/lado/
			fmt.Println("Archivo:", NombreArchivo) // nombre.tipo
			//SI LA CARPETA NO EXISTE Y EL COMANDO -P VIENE ACTIVO, SE PROCEDE A CREAR LA RUTA
			if pActivo == true {
				CrearDirectorio(id, rutaArchivo, true)
			}

			RArbol = nil
			contadorRuta = 0

			var ra RARBOL
			ra.nivel = 0
			ra.nombre = "/"
			ra.puntero = 0
			RArbol = append(RArbol, ra)

			verificarNivelesRuta(file, sb, cadenaDivididaSinArchivo, false, 0, 1)

			fmt.Println(red + "********************************************" + reset)
			for i := 0; i < len(RArbol); i++ {
				fmt.Println("Nivel: ", RArbol[i].nivel)
				fmt.Println("Nombre: ", RArbol[i].nombre, ", len: ", len(RArbol[i].nombre))
				fmt.Println("Puntero: ", RArbol[i].puntero)
				fmt.Println(cyan + "----------------" + reset)
			}
			fmt.Println(red + "********************************************" + reset)
			cuantoNivelesTieneLaRuta := len(cadenaDivididaSlash)
			cuantosNivelesCumple := 0

			//********************************************************
			//Se verifica cuandos padres existen y se valida para
			//escribir la cantidad necesaria
			//********************************************************
			var ArbolDeCumplimientos []RARBOL
			//fmt.Println(cuantoNivelesTieneLaRuta)
			for i := 0; i < cuantoNivelesTieneLaRuta; i++ {
				seCumplioCondicion := false
				for j := 0; j < len(RArbol); j++ {
					if RArbol[j].nivel == int64(i) && RArbol[j].nombre == cadenaDivididaSlash[i] {
						ArbolDeCumplimientos = append(ArbolDeCumplimientos, RArbol[j])
						seCumplioCondicion = true
						cuantosNivelesCumple++
					}
				}
				if seCumplioCondicion == false {
					break
				}
			}

			fmt.Println(ArbolDeCumplimientos)
			fmt.Println("Niveles de la ruta: ", cuantoNivelesTieneLaRuta)
			fmt.Println("Niveles que se cumplen: ", cuantosNivelesCumple)
			var cuantosNivelesNuevos int = cuantoNivelesTieneLaRuta - cuantosNivelesCumple
			fmt.Println("Niveles nuevo a crear: ", cuantosNivelesNuevos)
			existeElPadre := false
			if cuantosNivelesNuevos > 1 {
				fmt.Println(red + "[ERROR]" + reset + "El padre de la carpeta que desea crear, no existe")
			} else if cuantosNivelesNuevos == 1 {
				fmt.Println(green + "[EXITO]" + reset + "Si existe el padre")
				existeElPadre = true
			}
			if existeElPadre == true {
				fmt.Println(red + "********************************************" + reset)
				fmt.Println("Nombre Carpeta: ", RArbol[len(RArbol)-1].nombre)
				fmt.Println("Puntero a DD: ", RArbol[len(RArbol)-1].puntero)
				fmt.Println(red + "********************************************" + reset)
				//********************************************************
				//Se verifica la posicion disponible en el bitmap Inodos
				//para la escritura
				//********************************************************
				posicionLibreEnBitmapTablaInodo := 0
				for j := 0; j < len(bitmapTablaInodo); j++ {
					if bitmapTablaInodo[j] == '0' {
						posicionLibreEnBitmapTablaInodo = j
						bitmapTablaInodo[j] = '1'
						break
					}
				}
				//********************************************************
				//Se obtiene la posicion del DD
				//********************************************************
				romplerCicloInfinito := false
				avd := leerAVD(file, sb.SbApArbolDirectorio+(sb.SbSizeStructArbolDirectorio*(RArbol[len(RArbol)-1].puntero)))
				posicionDetalleDirectorio := avd.AvdApDetalleDirectorio - 1
				var queTablaInodo int64
				fmt.Println("DETALLE DIRECTORIO: ", avd.AvdApDetalleDirectorio)
				for {
					fmt.Println(green+"POSICION: ", posicionDetalleDirectorio, reset)
					//********************************************************
					//Se busca espacio disponible para escribir el nombre del
					//archivo
					//********************************************************
					ddAnterior := leerDD(file, sb.SbApDetalleDirectorio+(sb.SbSizeStructDetalleDirectorio*(posicionDetalleDirectorio)))
					hayDisponibilidadDeEspacio := false
					for k := 0; k < 5; k++ {
						if ddAnterior.DdArrayFiles[k].DdFileApInodo == 1 {

						} else if ddAnterior.DdArrayFiles[k].DdFileApInodo == 0 {
							fmt.Println("LIBRE: ", k)
							copy(ddAnterior.DdArrayFiles[k].DdFileNombre[:], NombreArchivo)
							ddAnterior.DdArrayFiles[k].DdFileApInodo = int64(posicionLibreEnBitmapTablaInodo + 1)
							queTablaInodo = int64(posicionLibreEnBitmapTablaInodo + 1)
							hayDisponibilidadDeEspacio = true
							romplerCicloInfinito = true
							break
						}
					}
					if hayDisponibilidadDeEspacio == false {
						if ddAnterior.DdApDetalleDirectorio == 0 {
							//********************************************************
							//Se buscan nuevas posiciones en el bitmap DD para
							//escribir sin interferir con las posiciones que ya estan
							//reservadas para escribir el nombre
							//********************************************************
							var posicionAEscribirNuevoDetalleDirectorio int64
							bitmapDetalleDirectorio := retornarBitmap(file, sb.SbApBitmapDetalleDirectorio, sb)
							//fmt.Println(red, bitmapDetalleDirectorio, reset)
							for j := 0; j < len(bitmapDetalleDirectorio); j++ {
								//fmt.Println(yellow+"ESTA ES LA POSICION EN USO: "+reset, j)
								if bitmapDetalleDirectorio[j] == '0' {
									fmt.Println(cyan+"ESTA ES LA POSICION DEL NUEVO DD: "+reset, j+1)
									posicionAEscribirNuevoDetalleDirectorio = int64(j + 1)
									bitmapDetalleDirectorio[j] = '1'
									break
								}
							}
							//fmt.Println(red, bitmapDetalleDirectorio, reset)
							file.Seek(sb.SbApBitmapDetalleDirectorio, 0)
							var valorBitmapDetalleDirectorio bytes.Buffer
							binary.Write(&valorBitmapDetalleDirectorio, binary.BigEndian, &bitmapDetalleDirectorio)
							escribirBytes(file, valorBitmapDetalleDirectorio.Bytes())
							//********************************************************
							//Se reapunta el detalle directorio al nuevo DD
							//********************************************************
							ddAnterior.DdApDetalleDirectorio = posicionAEscribirNuevoDetalleDirectorio

							file.Seek(sb.SbApDetalleDirectorio+(sb.SbSizeStructDetalleDirectorio*(posicionDetalleDirectorio)), 0)
							var valorBinarioAntiguoDetalleDirectorio bytes.Buffer
							binary.Write(&valorBinarioAntiguoDetalleDirectorio, binary.BigEndian, &ddAnterior)
							escribirBytes(file, valorBinarioAntiguoDetalleDirectorio.Bytes())
							//********************************************************
							//Se escribe el nuevo directorio
							//********************************************************
							var ddNuevo DETALLEDIRECTORIO
							ddNuevo.DdArrayFiles[0].DdFileApInodo = 0
							ddNuevo.DdArrayFiles[1].DdFileApInodo = 0
							ddNuevo.DdArrayFiles[2].DdFileApInodo = 0
							ddNuevo.DdArrayFiles[3].DdFileApInodo = 0
							ddNuevo.DdArrayFiles[4].DdFileApInodo = 0
							ddNuevo.DdApDetalleDirectorio = 0

							file.Seek(sb.SbApDetalleDirectorio+(sb.SbSizeStructDetalleDirectorio*(posicionAEscribirNuevoDetalleDirectorio-1)), 0)
							var valorBinarioNuevoDetalleDirectorio bytes.Buffer
							binary.Write(&valorBinarioNuevoDetalleDirectorio, binary.BigEndian, &ddNuevo)
							escribirBytes(file, valorBinarioNuevoDetalleDirectorio.Bytes())
						} else {
							posicionDetalleDirectorio = ddAnterior.DdApDetalleDirectorio - 1
						}
					} else {
						//********************************************************
						//Se reescribe el detalle directorio
						//********************************************************
						file.Seek(sb.SbApDetalleDirectorio+(sb.SbSizeStructDetalleDirectorio*(posicionDetalleDirectorio)), 0)
						var valorBinarioDetalleDirectorio bytes.Buffer
						binary.Write(&valorBinarioDetalleDirectorio, binary.BigEndian, &ddAnterior)
						escribirBytes(file, valorBinarioDetalleDirectorio.Bytes())
						//********************************************************
						//Se escribe el nuevo TABLA INODO
						//********************************************************
						var ti TABLAINODO
						ti.ICountInodo = int64(posicionLibreEnBitmapTablaInodo + 1)
						ti.ISizeArchivo = 0
						ti.ICountBloquesAsignados = 0
						ti.IArrayBloques[0] = 0
						ti.IArrayBloques[1] = 0
						ti.IArrayBloques[2] = 0
						ti.IArrayBloques[3] = 0
						ti.IApIndirecto = 0
						copy(ti.IIdProper[:], SesionActiva.usuario)
						//********************************************************
						//Se reescribe el bitmapTablaInodo
						//********************************************************
						file.Seek(sb.SbApBitmapTablaInodo, 0)
						var valorBitmapTablaInodo bytes.Buffer
						binary.Write(&valorBitmapTablaInodo, binary.BigEndian, &bitmapTablaInodo)
						escribirBytes(file, valorBitmapTablaInodo.Bytes())
					}
					if romplerCicloInfinito == true {
						break
					}
				}

				//********************************************************
				//Se procede a verificar si el size y cont estan activos,
				//de ser asi, se procede a ingresar el contenido del
				//archivo
				//********************************************************
				if size != "" {
					bitmapBloques := retornarBitmap(file, sb.SbApBitmapBloques, sb)
					//********************************************************
					//Se realiza el conteo de cuantos bloques se necesitan y
					//cuantos hay disponibles
					//********************************************************
					fmt.Println(queTablaInodo)
					tamEnBytesArchivo, _ := strconv.ParseInt(size, 10, 64)
					totalDeBloques := tamEnBytesArchivo / 25.0
					restoBloques := tamEnBytesArchivo % 25.0
					if restoBloques != 0 {
						totalDeBloques++
					}
					fmt.Println("Total de Bloques a usar: ", totalDeBloques)
					var contadorBloquesDisponibles int64 = 0
					for i := 0; i < len(bitmapBloques); i++ {
						if bitmapBloques[i] == '0' {
							contadorBloquesDisponibles++
						}
					}
					fmt.Println("BITMAP BLOQUES: ", bitmapBloques)
					fmt.Println("Bloques Disponibles: ", contadorBloquesDisponibles)
					if contadorBloquesDisponibles >= totalDeBloques {
						numEstructuraTreeComplete = 0 //Solo se reinicia
						numEstructuraTreeComplete = queTablaInodo - 1
						CadenaRetornoArchivo = ""
						if cont != "" {
							CadenaRetornoArchivo = cont
						} else {
							abecedario := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
							contadorLetra := 0
							for i := 0; i < int(tamEnBytesArchivo); i++ {
								if contadorLetra == 26 {
									contadorLetra = 0
								}
								CadenaRetornoArchivo += abecedario[contadorLetra]
								contadorLetra++
							}
						}
						ApuntadoresBloqueUsoArchivo = nil
						ApuntadoresInodosUsoArchivo = nil
						modificarArchivo(file, sb, queTablaInodo-1, 3)
					} else {
						fmt.Println(red + "[ERROR]" + reset + "No hay suficientes bloques disponibles para agregar los datos")
					}
				}
			}
		}
	}
}

//DIRECTORIOLIBRE =
type DIRECTORIOLIBRE struct {
	puntero int64
}

//DirectorioLibre =
var DirectorioLibre []DIRECTORIOLIBRE

func verificarDisponibilidadDirectorio(file *os.File, sb SUPERBOOT, posicion int64, tipoArchivo int) {
	if tipoArchivo == 2 { //SE UTILIZA PARA RECORRER AL PADRE
		var dd DETALLEDIRECTORIO
		dd = leerDD(file, sb.SbApDetalleDirectorio+(sb.SbSizeStructDetalleDirectorio*posicion))
		fmt.Println("ESTOY EN: ", posicion)
		var cont int = 0
		for i := 0; i < 6; i++ {
			if i < 5 && dd.DdArrayFiles[i].DdFileApInodo == 0 {
				var dl DIRECTORIOLIBRE
				dl.puntero = int64(i)
				DirectorioLibre = append(DirectorioLibre, dl)
				break
			} else if i < 5 && dd.DdArrayFiles[i].DdFileApInodo != 0 {
				fmt.Println("[TIPO2]pos ", i, " ->", dd.DdArrayFiles[i].DdFileApInodo)
				cont++
			} else if i == 5 && dd.DdApDetalleDirectorio != 0 && cont == 5 {
				fmt.Println("[TIPO2]pos ind", i)
				verificarDisponibilidadDirectorio(file, sb, dd.DdApDetalleDirectorio-1, 1)
			}
		}
	}
}

//CadenaRetornoArchivo =
var CadenaRetornoArchivo string

func retornarCadenaTextoArchivoArbolRecursivo(file *os.File, sb SUPERBOOT, tipoArchivo int) {
	//Se empieza a recorrer desde el inodo del archivo USER.TXT que siempre sera el primer INODO
	if tipoArchivo == 3 { //CUANDO ES TIPO TABLA INODO
		ti := leerTABLAINODO(file, sb.SbApTablaInodo+(sb.SbSizeStructInodo*numEstructuraTreeComplete))
		for i := 0; i < 5; i++ {
			if i < 4 && ti.IArrayBloques[i] != 0 {
				numEstructuraTreeCompleteMov = ti.ICountInodo - 1
				fmt.Println("[TIPO3]pos ", i, ",", numEstructuraTreeCompleteMov, ":", ti.IArrayBloques[i])
				numEstructuraTreeComplete = ti.IArrayBloques[i]
				numEstructuraTreeComplete--
				recorrerArbolRecursivoRetornarUsersTxt(file, sb, 4)
			} else if i == 4 && ti.IApIndirecto != 0 {
				fmt.Println("[TIPO3]pos ind", i, ":", ti.IApIndirecto)
				numEstructuraTreeComplete = ti.IApIndirecto
				numEstructuraTreeComplete--
				recorrerArbolRecursivoRetornarUsersTxt(file, sb, 3)
			}
		}
	} else if tipoArchivo == 4 { //CUANDO ES TIPO BLOQUE DE DATOS
		bd := leerBLOQUEDATOS(file, sb.SbApBloques+(sb.SbSizeStructBloque*numEstructuraTreeComplete))
		for i := 0; i < len(bd.DbDato); i++ {
			if bd.DbDato[i] != 0 {
				CadenaRetornoArchivo += string(bd.DbDato[i])
			}
		}
		fmt.Println("[TIPO4]pos: ", string(bd.DbDato[:]))
	}
}

func modificarArchivo(file *os.File, sb SUPERBOOT, inodoPrincipal int64, tipo int) {
	//********************************************************
	//Se carga el bitmap del inodo y bloques
	//********************************************************
	bitmapTablaInodos := retornarBitmap(file, sb.SbApBitmapTablaInodo, sb)
	bitmapBloques := retornarBitmap(file, sb.SbApBitmapBloques, sb)
	fmt.Println(cyan + "Cadena: " + reset + CadenaRetornoArchivo)
	tamCadena := len(CadenaRetornoArchivo)
	//********************************************************
	//Se reinician los arreglos de apuntadores, los cuales se
	//encargaran de almacenar los apuntadores que estan en uso
	//para poder reescribirlos luego
	//********************************************************
	ApuntadoresBloqueUsoArchivo = nil
	ApuntadoresInodosUsoArchivo = nil
	ApuntadoresInodosUsoArchivo = append(ApuntadoresInodosUsoArchivo, inodoPrincipal+1) //Se agrega el inodo que esta siempre en uso
	fmt.Println("Inodo Principal: ", inodoPrincipal)
	numEstructuraTreeComplete = inodoPrincipal
	recorrerArbolRecursivoRetornarApuntadoresArchivo(file, sb, 3)
	fmt.Println("Tamaño cadena: ", tamCadena)
	fmt.Println("Bitmap Inodos: ", bitmapTablaInodos)
	fmt.Println("Bitmap Bloques: ", bitmapBloques)
	//********************************************************
	//Se calculan cuantos bloques se usaran y se valida si hay
	//espacio suficiente en el bitmap
	//********************************************************
	totalDeBloques := tamCadena / 25.0
	restoBloques := tamCadena % 25.0
	if restoBloques != 0 {
		totalDeBloques++
	}
	totalDeBloquesNuevos := totalDeBloques - len(ApuntadoresBloqueUsoArchivo)
	fmt.Println("Total Bloques: ", totalDeBloques)
	fmt.Println("Bloques nuevos: ", totalDeBloquesNuevos)
	fmt.Println("Tamano Bitmap Bloques: ", len(bitmapBloques))
	fmt.Println("Apuntadores bloques en uso: ", ApuntadoresBloqueUsoArchivo[:])
	fmt.Println("Apuntadores indodos en uso: ", ApuntadoresInodosUsoArchivo[:])

	var todasLasPosicionesAEscribirBloques []int64
	if ApuntadoresBloqueUsoArchivo != nil {
		todasLasPosicionesAEscribirBloques = append(todasLasPosicionesAEscribirBloques[:], ApuntadoresBloqueUsoArchivo[:]...)
	}
	contadorBloquesLibres := totalDeBloquesNuevos

	existenBloquesSuficientes := false
	if len(bitmapBloques) < totalDeBloquesNuevos {
		fmt.Println(red + "[ERROR]" + reset + "No hay suficientes BLOQUES para insertar el archivo")
	} else {
		for i := 0; i < len(bitmapBloques); i++ {
			if contadorBloquesLibres == 0 {
				break
			}
			if bitmapBloques[i] == '0' {
				todasLasPosicionesAEscribirBloques = append(todasLasPosicionesAEscribirBloques[:], int64(i+1))
				contadorBloquesLibres--
			}
		}
		fmt.Println(magenta + "TODAS LAS POSICIONES DE BLOQUES: " + reset)
		fmt.Println(todasLasPosicionesAEscribirBloques)
		existenBloquesSuficientes = true
	}

	//********************************************************
	//Se calculan cuantos inodos se usaran y se valida si hay
	//espacio suficiente en el bitmap
	//********************************************************

	totalDeInodos := totalDeBloques / 4
	restoInodos := totalDeBloques % 4
	if restoInodos != 0 {
		totalDeInodos++
	}
	totalDeInodosNuevos := totalDeInodos - len(ApuntadoresInodosUsoArchivo)
	fmt.Println("Total de inodos nuevos: ", totalDeInodosNuevos)

	var todasLasPosicionesAEscribirInodos []int64
	if ApuntadoresInodosUsoArchivo != nil {
		todasLasPosicionesAEscribirInodos = append(todasLasPosicionesAEscribirInodos[:], ApuntadoresInodosUsoArchivo[:]...)
	}
	contadorInodosLibres := totalDeInodosNuevos

	existenInodosSuficientes := false
	if len(bitmapTablaInodos) < totalDeInodosNuevos {
		fmt.Println(red + "[ERROR]" + reset + "No hay suficientes INODOS para insertar el archivo")
	} else {
		for i := 0; i < len(bitmapTablaInodos); i++ {
			if contadorInodosLibres == 0 {
				break
			}
			if bitmapTablaInodos[i] == '0' {
				todasLasPosicionesAEscribirInodos = append(todasLasPosicionesAEscribirInodos[:], int64(i+1))
				contadorInodosLibres--
			}
		}
		fmt.Println(magenta + "TODAS LAS POSICIONES DE INODOS: " + reset)
		fmt.Println(todasLasPosicionesAEscribirInodos)
		existenInodosSuficientes = true
	}

	if existenBloquesSuficientes == true && existenInodosSuficientes == true {
		//********************************************************
		//Se escriben las nuevas estructuras
		//********************************************************
		fmt.Println("T: ", totalDeBloques)
		contPosBloque := totalDeBloques
		numPosBloque := 0
		for i := 0; i < len(todasLasPosicionesAEscribirInodos); i++ {
			var nuevoTI TABLAINODO
			nuevoTI.ICountInodo = todasLasPosicionesAEscribirInodos[i]
			nuevoTI.ISizeArchivo = int64(totalDeBloques) * sb.SbSizeStructBloque
			nuevoTI.ICountBloquesAsignados = int64(totalDeBloques)
			for j := 0; j < 4; j++ {
				if contPosBloque != 0 {
					nuevoTI.IArrayBloques[j] = todasLasPosicionesAEscribirBloques[numPosBloque]
					numPosBloque++
					contPosBloque--
				}
			}
			if contPosBloque != 0 {
				nuevoTI.IApIndirecto = todasLasPosicionesAEscribirInodos[i+1]
			}
			copy(nuevoTI.IIdProper[:], SesionActiva.usuario)
			file.Seek(sb.SbApTablaInodo+(sb.SbSizeStructInodo*(todasLasPosicionesAEscribirInodos[i]-1)), 0)
			var valorBinarioTablaInodo bytes.Buffer
			binary.Write(&valorBinarioTablaInodo, binary.BigEndian, &nuevoTI)
			escribirBytes(file, valorBinarioTablaInodo.Bytes())
		}

		cadenaUserTxtSeparada := SplitSubN(CadenaRetornoArchivo, 25)

		for i := 0; i < len(todasLasPosicionesAEscribirBloques); i++ {
			var nuevoBD BLOQUEDATOS
			var contenidoAByte25 [25]byte
			copy(contenidoAByte25[:], cadenaUserTxtSeparada[i])
			nuevoBD.DbDato = contenidoAByte25
			//CadenaRetornoArchivo
			//fmt.Println(string(contenidoAByte25[:]))
			//fmt.Println(todasLasPosicionesAEscribirBloques[i])
			file.Seek(sb.SbApBloques+(sb.SbSizeStructBloque*(todasLasPosicionesAEscribirBloques[i]-1)), 0)
			var valorBinarioBloque bytes.Buffer
			binary.Write(&valorBinarioBloque, binary.BigEndian, &nuevoBD)
			escribirBytes(file, valorBinarioBloque.Bytes())
		}

		//********************************************************
		//Se escriben las nuevas estructuras
		//********************************************************
		nuevoBitmapTablaInodos := bitmapTablaInodos
		for i := 0; i < len(todasLasPosicionesAEscribirInodos); i++ {
			nuevoBitmapTablaInodos[todasLasPosicionesAEscribirInodos[i]-1] = '1'
		}
		fmt.Println(red, nuevoBitmapTablaInodos, reset)
		reescribirBitmap(file, sb.SbApBitmapTablaInodo, nuevoBitmapTablaInodos)
		nuevoBitmapBloques := bitmapBloques
		for i := 0; i < len(todasLasPosicionesAEscribirBloques); i++ {
			nuevoBitmapBloques[todasLasPosicionesAEscribirBloques[i]-1] = '1'
		}
		reescribirBitmap(file, sb.SbApBitmapBloques, nuevoBitmapBloques)
	}
}

//ApuntadoresBloqueUsoArchivo = Arreglo que almacena los apuntadores de tipo bloque que usa el archivo
var ApuntadoresBloqueUsoArchivo []int64

//ApuntadoresInodosUsoArchivo = Arreglo que almacena los apuntadores de tipo Inodo que usa el archivo
var ApuntadoresInodosUsoArchivo []int64

func recorrerArbolRecursivoRetornarApuntadoresArchivo(file *os.File, sb SUPERBOOT, tipoArchivo int) {
	//Se empieza a recorrer desde el inodo del archivo USER.TXT que siempre sera el primer INODO
	if tipoArchivo == 3 { //CUANDO ES TIPO TABLA INODO
		ti := leerTABLAINODO(file, sb.SbApTablaInodo+(sb.SbSizeStructInodo*numEstructuraTreeComplete))
		for i := 0; i < 5; i++ {
			if i < 4 && ti.IArrayBloques[i] != 0 {
				numEstructuraTreeComplete = ti.IArrayBloques[i]
				ApuntadoresBloqueUsoArchivo = append(ApuntadoresBloqueUsoUSERTXT, ti.IArrayBloques[i])
				numEstructuraTreeComplete--
			} else if i == 4 && ti.IApIndirecto != 0 {
				numEstructuraTreeComplete = ti.IApIndirecto
				ApuntadoresInodosUsoArchivo = append(ApuntadoresInodosUsoUSERTXT, ti.IApIndirecto)
				numEstructuraTreeComplete--
				recorrerArbolRecursivoRetornarApuntadoresArchivo(file, sb, 3)
			}
		}
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
	PathID, _, ExisteID := existeID(id)
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

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//SUPERBOOT-----SUPERBOOT-----REPORTE-----REPORTE-----SUPERBOOT-----SUPERBOOT-----REPORTE-----REPORTE-----SUPERBOOT-----SUPERBOOT-----REPORTE
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//ReporteSB = Genera el reporte del Super Boot
func ReporteSB(id, path string) {
	PathID, NombreParticion, ExisteID := existeID(id) //Path, nombre particion, bool si existe
	var cadenaReporteSB string
	var start int64
	//var cadenaReporteDISKAuxiliar string = ""

	if ExisteID == true {
		mbr := leerMBR(PathID)
		var nombreAByte16 [16]byte
		copy(nombreAByte16[:], NombreParticion)
		if mbr.MbrPartition1.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition2.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition3.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition4.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		}
		//******************************************************
		//Se inicia el el llenado de los datos de Graphviz
		//******************************************************

		sb := leerSB(PathID, start)
		cadenaReporteSB = "digraph MBR {\nnode [shape=plaintext]\nA [label=<\n<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#5A69D6' COLSPAN='2'><font color='white'>REPORTE SB</font></TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#6C7AE0'><font color='white'>NOMBRE</font></TD><TD BGCOLOR='#6C7AE0'><font color='white'>VALOR</font></TD>\n</TR>\n"
		var nombreDiscoSin0s string
		for i1, valor1 := range sb.SbNombreHd {
			if sb.SbNombreHd[i1] != 0 {
				nombreDiscoSin0s += string(valor1)
			}
		}
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_nombre_hd</TD><TD>" + nombreDiscoSin0s + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_arbol_virtual_count</TD><TD>" + strconv.FormatInt(sb.SbArbolVirtualCount, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_detalle_directorio_count</TD><TD>" + strconv.FormatInt(sb.SbDetalleDirectorioCount, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_inodos_count</TD><TD>" + strconv.FormatInt(sb.SbInodosCount, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_bloques_count</TD><TD>" + strconv.FormatInt(sb.SbBloquesCount, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_arbol_virtual_free</TD><TD>" + strconv.FormatInt(sb.SbArbolVirtualFree, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_detalle_directorio_free</TD><TD>" + strconv.FormatInt(sb.SbDetalleDirectorioFree, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_inodos_free</TD><TD>" + strconv.FormatInt(sb.SbInodosFree, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_bloques_free</TD><TD>" + strconv.FormatInt(sb.SbBloquesFree, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_date_creacion</TD><TD>" + string(sb.SbDateCreacion[:19]) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_date_ultimo_montaje</TD><TD>" + string(sb.SbDateUltimoMontaje[:19]) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_montaje_count</TD><TD>" + strconv.FormatInt(sb.SbMontajeCount, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_ap_bitmap_arbol_directorio</TD><TD>" + strconv.FormatInt(sb.SbApBitmapArbolDirectorio, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_ap_arbol_directorio</TD><TD>" + strconv.FormatInt(sb.SbApArbolDirectorio, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_ap_bitmap_detalle_directorio</TD><TD>" + strconv.FormatInt(sb.SbApBitmapDetalleDirectorio, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_ap_detalle_directorio</TD><TD>" + strconv.FormatInt(sb.SbApDetalleDirectorio, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_ap_bitmap_tabla_inodo</TD><TD>" + strconv.FormatInt(sb.SbApBitmapTablaInodo, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_ap_tabla_inodo</TD><TD>" + strconv.FormatInt(sb.SbApTablaInodo, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_ap_bitmap_bloques</TD><TD>" + strconv.FormatInt(sb.SbApBitmapBloques, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_ap_bloques</TD><TD>" + strconv.FormatInt(sb.SbApBloques, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_ap_log</TD><TD>" + strconv.FormatInt(sb.SbApLog, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_size_struct_arbol_directorio</TD><TD>" + strconv.FormatInt(sb.SbSizeStructArbolDirectorio, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_size_struct_detalle_directorio</TD><TD>" + strconv.FormatInt(sb.SbSizeStructDetalleDirectorio, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_size_struct_inodo</TD><TD>" + strconv.FormatInt(sb.SbSizeStructInodo, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_size_struct_bloque</TD><TD>" + strconv.FormatInt(sb.SbSizeStructBloque, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_first_free_bit_arbol_directorio</TD><TD>" + strconv.FormatInt(sb.SbFirstFreeBitArbolDirectorio, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_first_free_bit_detalle_directorio</TD><TD>" + strconv.FormatInt(sb.SbFirstFreeBitDetalleDirectorio, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_first_free_bit_tabla_inodo</TD><TD>" + strconv.FormatInt(sb.SbFirstFreeBitTablaInodo, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_first_free_bit_bloques</TD><TD>" + strconv.FormatInt(sb.SbFirstFreeBitBloques, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "<TR>\n<TD BGCOLOR='#A1ACFC'>sb_magic_num</TD><TD>" + strconv.FormatInt(sb.SbMagicNum, 10) + "</TD>\n</TR>\n"
		cadenaReporteSB += "</TABLE>\n>];\n}"

		//******************************************************************
		//Se escribe la cadena en el archivo .svg que usara Graphviz
		//******************************************************************
		nombreGV, nombreExtension := crearArchivoParaReporte(path, cadenaReporteSB)
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
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//BITMAP-----BITMAP-----REPORTE-----REPORTE-----BITMAP-----BITMAP-----REPORTE-----REPORTE-----BITMAP-----BITMAP-----REPORTE-----REPORTE------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//ReporteBitmap = Genera el reporte de bitmap de Arbol de Directorios
func ReporteBitmap(id, path string, tipoReporte int) {
	PathID, NombreParticion, ExisteID := existeID(id) //Path, nombre particion, bool si existe
	var cadenaReporteBitmap string = ""
	var start int64

	if ExisteID == true {
		mbr := leerMBR(PathID)
		var nombreAByte16 [16]byte
		copy(nombreAByte16[:], NombreParticion)
		if mbr.MbrPartition1.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition2.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition3.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition4.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		}
		//**********************************************************************
		//Se inicia la declaracion de variables para llenar el arreglo de bitmap
		//segun sea el tipo de reporte especificado
		//**********************************************************************
		//TIPO 1: es el reporte del ARBOL VIRTUAL DE DIRECTORIO
		//TIPO 2: es el reporte del DETALLE DE DIRECTORIO
		//TIPO 3: es el reporte de la TABLA DE INODOS
		//TIPO 4: es el reporte de los BLOQUES
		sb := leerSB(PathID, start)
		if tipoReporte == 1 {
			cadenaReporteBitmap = leerBitmapResumen(PathID, sb.SbApBitmapArbolDirectorio, &sb)
		} else if tipoReporte == 2 {
			cadenaReporteBitmap = leerBitmapResumen(PathID, sb.SbApBitmapDetalleDirectorio, &sb)
		} else if tipoReporte == 3 {
			cadenaReporteBitmap = leerBitmapResumen(PathID, sb.SbApBitmapTablaInodo, &sb)
		} else if tipoReporte == 4 {
			cadenaReporteBitmap = leerBitmapResumen(PathID, sb.SbApBitmapBloques, &sb)
		}
		crearArchivoParaReporteBitmap(path, cadenaReporteBitmap)
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//ARBOL-----ARBOL-----REPORTE-----REPORTE-----ARBOL-----ARBOL-----REPORTE-----REPORTE-----ARBOL-----ARBOL-----REPORTE-----REPORTE-----ARBOL--
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//ReporteTreeComplete = Metodo que crea el reporte del arbol completo
func ReporteTreeComplete(id, path string) {
	PathID, NombreParticion, ExisteID := existeID(id) //Path, nombre particion, bool si existe
	//var cadenaReporteTreeComplete string = ""
	var start int64

	if ExisteID == true {
		mbr := leerMBR(PathID)
		var nombreAByte16 [16]byte
		copy(nombreAByte16[:], NombreParticion)
		if mbr.MbrPartition1.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition2.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition3.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		} else if mbr.MbrPartition4.PartName == nombreAByte16 {
			start = mbr.MbrPartition1.PartStart
		}
		sb := leerSB(PathID, start)

		//Se abre el archivo para su uso
		file, err := os.OpenFile(PathID, os.O_RDWR, 0755)
		defer file.Close()
		if err != nil {
			fmt.Println(red + "[ERROR]" + reset + "No se ha podido abrir el archivo")
		}

		//**********************************************************************
		//Se preparan los bitmap para saber como moverse en el archivo
		//**********************************************************************
		bitmapArbolVirtualDirectorio := retornarBitmap(file, sb.SbApBitmapArbolDirectorio, sb)
		//fmt.Println(bitmapArbolVirtualDirectorio)
		bitmapDetalleDirectorio := retornarBitmap(file, sb.SbApBitmapDetalleDirectorio, sb)
		//fmt.Println(bitmapDetalleDirectorio)
		bitmapTablaInodos := retornarBitmap(file, sb.SbApBitmapTablaInodo, sb)
		//fmt.Println(bitmapTablaInodos)
		bitmapBloques := retornarBitmap(file, sb.SbApBitmapBloques, sb)
		//fmt.Println(bitmapBloques)

		//**********************************************************************
		//Se preparan los tamaños de los archivos
		//**********************************************************************
		//tamSuperBoot := int64(binary.Size(SUPERBOOT{}))
		//tamArbolVirtualDirectorio := sb.SbSizeStructArbolDirectorio
		//tamDetalleDirectorio := sb.SbSizeStructDetalleDirectorio
		//tamTablaInodo := sb.SbSizeStructInodo
		//tamBloqueDatos := sb.SbSizeStructBloque
		var cadenaReporteTreeComplete string
		if bitmapArbolVirtualDirectorio[0] == '1' {
			cadenaReporteTreeComplete = "digraph MBR {\nnode [shape=plaintext]\nrankdir=LR;\n"
			//aDondeVoy := sb.SbApArbolDirectorio
			subCadenaReporteTreeComplete = ""
			numEstructuraTreeComplete = 0
			numEstructuraTreeCompleteMov = 0
			numEstructuraTreeCompleteAux = 0
			//recorrerArbolRecursivoReporte(file, sb, 1)
			for i := 0; i < len(bitmapArbolVirtualDirectorio); i++ {
				if bitmapArbolVirtualDirectorio[i] == '1' {
					recorrerArbolRecursivoReportePorBitmap(file, sb, int64(i), 1)
				}
			}
			for i := 0; i < len(bitmapDetalleDirectorio); i++ {
				if bitmapDetalleDirectorio[i] == '1' {
					recorrerArbolRecursivoReportePorBitmap(file, sb, int64(i), 2)
				}
			}
			for i := 0; i < len(bitmapTablaInodos); i++ {
				if bitmapTablaInodos[i] == '1' {
					recorrerArbolRecursivoReportePorBitmap(file, sb, int64(i), 3)
				}
			}
			for i := 0; i < len(bitmapBloques); i++ {
				if bitmapBloques[i] == '1' {
					recorrerArbolRecursivoReportePorBitmap(file, sb, int64(i), 4)
				}
			}
			//recorrerArbolRecursivoReporte(file, sb, 0, 1)
			cadenaReporteTreeComplete += subCadenaReporteTreeComplete
			cadenaReporteTreeComplete += "}"
			//fmt.Println(red + "------------" + reset)
			//numEstructuraTreeComplete = 0
			//numEstructuraTreeCompleteMov = 0
			//recorrerArbolRecursivo(file, sb, 1)
			//******************************************************************
			//Se escribe la cadena en el archivo .svg que usara Graphviz
			//******************************************************************
			nombreGV, nombreExtension := crearArchivoParaReporte(path, cadenaReporteTreeComplete)
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
		}
	}
	//leerBitmapResumen(path, )
}

var subCadenaReporteTreeComplete string
var numEstructuraTreeComplete int64 = 0
var numEstructuraTreeCompleteMov int64 = 0
var numEstructuraTreeCompleteAux int64 = 0
var static int64 = 0
var contadorReporteTree int

func recorrerArbolRecursivoReporte(file *os.File, sb SUPERBOOT, movimiento int64, tipoArchivo int) {
	//Se empieza a recorrer desde root '/'
	if tipoArchivo == 1 { //CUANDO ES TIPO ARBOL DE DIRECTORIO
		avd := leerAVD(file, sb.SbApArbolDirectorio+(sb.SbSizeStructArbolDirectorio*numEstructuraTreeComplete))
		subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + "[label=<\n<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"
		var nombreDirectorio string
		for i1, valor1 := range avd.AvdNombreDirectorio {
			if avd.AvdNombreDirectorio[i1] != 0 {
				nombreDirectorio += string(valor1)
			}
		}
		subCadenaReporteTreeComplete += "<TR port='0'>\n<TD BGCOLOR='#99ccff' COLSPAN='2'><font color='black'>" + nombreDirectorio + "</font></TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>Fecha Creacion</TD><TD>" + string(avd.AvdFechaCreacion[:]) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 1</TD><TD port='1'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[0], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 2</TD><TD port='2'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[1], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 3</TD><TD port='3'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[2], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 4</TD><TD port='4'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[3], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 5</TD><TD port='5'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[4], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 6</TD><TD port='6'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[5], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>Detalle Directorio</TD><TD port='7'>" + strconv.FormatInt(avd.AvdApDetalleDirectorio, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APDI 1</TD><TD port='8'>" + strconv.FormatInt(avd.AvdApArbolVirtualDirectorio, 10) + "</TD>\n</TR>\n"
		var nombreProper string
		for i1, valor1 := range avd.AvdProper {
			if avd.AvdProper[i1] != 0 {
				nombreProper += string(valor1)
			}
		}
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>PROPER</TD><TD>" + nombreProper + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "</TABLE>\n>];\n\n"
		for i := 0; i < 8; i++ {
			numEstructuraTreeCompleteAux = numEstructuraTreeComplete
			if i < 6 && avd.AvdApArraySubdirectorios[i] != 0 {
				numEstructuraTreeCompleteMov = numEstructuraTreeCompleteAux
				//numEstructuraTreeCompleteMov = numEstructuraTreeComplete
				fmt.Println("[TIPO1]pos ", i, ":", avd.AvdNombreDirectorio)
				subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(numEstructuraTreeCompleteMov+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				numEstructuraTreeComplete = avd.AvdApArraySubdirectorios[i]
				movimiento = avd.AvdApArraySubdirectorios[i]
				subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + "\n"
				//numEstructuraTreeComplete--
				//recorrerArbolRecursivoReporte(file, sb, 1)
				recorrerArbolRecursivoReporte(file, sb, movimiento, 1)
			} else if i == 6 && avd.AvdApDetalleDirectorio != 0 {
				fmt.Println("[TIPO1]pos det", i, ":", avd.AvdApDetalleDirectorio)
				subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(numEstructuraTreeCompleteMov+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				numEstructuraTreeComplete = avd.AvdApDetalleDirectorio
				subCadenaReporteTreeComplete += "DD" + strconv.FormatInt(numEstructuraTreeComplete, 10) + "\n"
				numEstructuraTreeComplete--
				//recorrerArbolRecursivoReporte(file, sb, 2)
				recorrerArbolRecursivoReporte(file, sb, movimiento, 2)
			} else if i == 7 && avd.AvdApArbolVirtualDirectorio != 0 {
				fmt.Println("[TIPO1]pos ind", i, ":", avd.AvdApArbolVirtualDirectorio)
				subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				numEstructuraTreeComplete = avd.AvdApArbolVirtualDirectorio
				subCadenaReporteTreeComplete += "AVD:" + strconv.FormatInt(numEstructuraTreeComplete, 10) + "\n"
				numEstructuraTreeComplete--
				//recorrerArbolRecursivoReporte(file, sb, 1)
				recorrerArbolRecursivoReporte(file, sb, movimiento, 1)
			}
		}
	} else if tipoArchivo == 2 { //CUANDO ES TIPO DETALLE DIRECTORIO
		dd := leerDD(file, sb.SbApDetalleDirectorio+(sb.SbSizeStructDetalleDirectorio*numEstructuraTreeComplete))
		subCadenaReporteTreeComplete += "\nDD" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + "[label=<\n<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#a3d977'>DETALLE<br/>DIRECTORIO</TD><TD BGCOLOR='#a3d977'>" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + "</TD>\n</TR>"
		for i := 0; i < 5; i++ {
			if dd.DdArrayFiles[i].DdFileApInodo != 0 {
				var nombreDirectorio string
				for i1, valor1 := range dd.DdArrayFiles[i].DdFileNombre {
					if dd.DdArrayFiles[i].DdFileNombre[i1] != 0 {
						nombreDirectorio += string(valor1)
					}
				}
				subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#a3d977'>" + nombreDirectorio + "</TD><TD port='1'>" + strconv.FormatInt(dd.DdArrayFiles[i].DdFileApInodo, 10) + "</TD>\n</TR>\n"
			} else {
				subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#a3d977'>APD" + strconv.Itoa(i+1) + "</TD><TD port='" + strconv.Itoa(i+1) + "'>-1</TD>\n</TR>\n"
			}
		}
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#a3d977'>API1</TD><TD port='6'>" + strconv.FormatInt(dd.DdApDetalleDirectorio, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "</TABLE>\n>];\n\n"
		for i := 0; i < 6; i++ {
			if i < 5 && dd.DdArrayFiles[i].DdFileApInodo != 0 {
				fmt.Println("[TIPO2]pos ", i, ":", dd.DdArrayFiles[i].DdFileApInodo)
				subCadenaReporteTreeComplete += "DD" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				numEstructuraTreeComplete = dd.DdArrayFiles[i].DdFileApInodo
				subCadenaReporteTreeComplete += "TI" + strconv.FormatInt(numEstructuraTreeComplete, 10) + "\n"
				//aDondeVoy = sb.SbApDetalleDirectorio + (sb.SbSizeStructDetalleDirectorio * numEstructura)
				numEstructuraTreeComplete--
				//recorrerArbolRecursivoReporte(file, sb, 3)
				recorrerArbolRecursivoReporte(file, sb, 0, 3)
			} else if i == 5 && dd.DdApDetalleDirectorio != 0 {
				fmt.Println("[TIPO2]pos ", i, ":", dd.DdApDetalleDirectorio)
				subCadenaReporteTreeComplete += "DD" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				numEstructuraTreeComplete = dd.DdApDetalleDirectorio
				subCadenaReporteTreeComplete += "DD:" + strconv.FormatInt(numEstructuraTreeComplete, 10) + "\n"
				numEstructuraTreeComplete--
				//recorrerArbolRecursivoReporte(file, sb, 2)
				recorrerArbolRecursivoReporte(file, sb, movimiento, 2)
			}
		}
	} else if tipoArchivo == 3 { //CUANDO ES TIPO TABLA INODO
		ti := leerTABLAINODO(file, sb.SbApTablaInodo+(sb.SbSizeStructInodo*numEstructuraTreeComplete))
		subCadenaReporteTreeComplete += "\nTI" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + " [label=<\n<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>TABLA INODO</TD><TD BGCOLOR='#ffc374'>" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>Tamaño</TD><TD>" + strconv.FormatInt(ti.ISizeArchivo, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>Bloques</TD><TD>" + strconv.FormatInt(ti.ICountBloquesAsignados, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>APD 1</TD><TD port='1'>" + strconv.FormatInt(ti.IArrayBloques[0], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>APD 2</TD><TD port='2'>" + strconv.FormatInt(ti.IArrayBloques[1], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>APD 3</TD><TD port='3'>" + strconv.FormatInt(ti.IArrayBloques[2], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>APD 4</TD><TD port='4'>" + strconv.FormatInt(ti.IArrayBloques[3], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>API</TD><TD port='5'>" + strconv.FormatInt(ti.IApIndirecto, 10) + "</TD>\n</TR>\n"
		var nombreProper string
		for i1, valor1 := range ti.IIdProper {
			if ti.IIdProper[i1] != 0 {
				nombreProper += string(valor1)
			}
		}
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>Proper</TD><TD>" + nombreProper + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "</TABLE>\n>];\n\n"
		for i := 0; i < 5; i++ {
			if i < 4 && ti.IArrayBloques[i] != 0 {
				numEstructuraTreeCompleteMov = ti.ICountInodo - 1
				fmt.Println("[TIPO3]pos ", i, ",", numEstructuraTreeCompleteMov, ":", ti.IArrayBloques[i])
				subCadenaReporteTreeComplete += "TI" + strconv.FormatInt(numEstructuraTreeCompleteMov+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				numEstructuraTreeComplete = ti.IArrayBloques[i]
				subCadenaReporteTreeComplete += "B" + strconv.FormatInt(numEstructuraTreeComplete, 10) + "\n"
				numEstructuraTreeComplete--
				//recorrerArbolRecursivoReporte(file, sb, 4)
				recorrerArbolRecursivoReporte(file, sb, movimiento, 4)
			} else if i == 4 && ti.IApIndirecto != 0 {
				fmt.Println("[TIPO3]pos ind", i, ":", ti.IApIndirecto)
				numEstructuraTreeCompleteMov = ti.ICountInodo
				subCadenaReporteTreeComplete += "TI" + strconv.FormatInt(numEstructuraTreeCompleteMov, 10) + ":" + strconv.Itoa(i+1) + " -> "
				numEstructuraTreeComplete = ti.IApIndirecto
				tiTemp := leerTABLAINODO(file, sb.SbApTablaInodo+(sb.SbSizeStructInodo*(numEstructuraTreeComplete-1)))
				subCadenaReporteTreeComplete += "TI" + strconv.FormatInt(tiTemp.ICountInodo, 10) + "\n"
				numEstructuraTreeComplete--
				//recorrerArbolRecursivoReporte(file, sb, 3)
				recorrerArbolRecursivoReporte(file, sb, movimiento, 3)
			}
		}
	} else if tipoArchivo == 4 { //CUANDO ES TIPO BLOQUE DE DATOS
		bd := leerBLOQUEDATOS(file, sb.SbApBloques+(sb.SbSizeStructBloque*numEstructuraTreeComplete))
		subCadenaReporteTreeComplete += "\nB" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + "[label=<\n\n<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ff8f80'>BLOQUE</TD><TD BGCOLOR='#ff8f80'>" + strconv.FormatInt(numEstructuraTreeComplete+1, 10) + "</TD>\n</TR>\n"
		var contenidoBloque string
		for i1, valor1 := range bd.DbDato {
			if bd.DbDato[i1] != 0 {
				contenidoBloque += string(valor1)
			}
		}
		subCadenaReporteTreeComplete += "<TR>\n<TD COLSPAN='2'>" + contenidoBloque + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "</TABLE>\n>];\n\n"
		fmt.Println("[TIPO4]pos: ", string(bd.DbDato[:]))
	}
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

//crearArchivoParaReporteBitmap = Metodo que se encarga de escribir un .TXT con el contenido de los bitmap
func crearArchivoParaReporteBitmap(path, cadena string) {
	ruta, nombreArchivo := filepath.Split(path)
	//Ruta: /home/user/xxxx
	//Nombre Archivo: yyy.go
	var extension = filepath.Ext(nombreArchivo)                       //.go
	var nombre = nombreArchivo[0 : len(nombreArchivo)-len(extension)] //yyy
	nombre += ".txt"

	if ExisteCarpeta(ruta) == false {
		CrearCarpeta(ruta)
	}
	if ExisteCarpeta(ruta) == true {
		//Se genera el archivo .gv (para uso de graphvz)
		file, err := os.Create(ruta + nombre)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		//Se escribe la informacion en el archivo
		err2 := ioutil.WriteFile(ruta+nombre, []byte(cadena), 0644)
		if err2 != nil {
			log.Fatal(err2)
		}
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL-----ARBOL----
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

func recorrerArbolRecursivo(file *os.File, sb SUPERBOOT, tipoArchivo int) {
	//Se empieza a recorrer desde root '/'
	if tipoArchivo == 1 { //CUANDO ES TIPO ARBOL DE DIRECTORIO
		avd := leerAVD(file, sb.SbApArbolDirectorio+(sb.SbSizeStructArbolDirectorio*numEstructuraTreeComplete))

		for i := 0; i < 8; i++ {
			if i < 6 && avd.AvdApArraySubdirectorios[i] != 0 {
				fmt.Println("[TIPO1]pos ", i, ":", avd.AvdNombreDirectorio)
				numEstructuraTreeComplete = avd.AvdApArraySubdirectorios[i]
				numEstructuraTreeComplete--
				recorrerArbolRecursivo(file, sb, 1)
			} else if i == 6 && avd.AvdApDetalleDirectorio != 0 {
				fmt.Println("[TIPO1]pos ", i, ":", avd.AvdApDetalleDirectorio)
				numEstructuraTreeComplete = avd.AvdApDetalleDirectorio
				numEstructuraTreeComplete--
				recorrerArbolRecursivo(file, sb, 2)
			} else if i == 7 && avd.AvdApArbolVirtualDirectorio != 0 {
				fmt.Println("[TIPO1]pos ", i, ":", avd.AvdApArbolVirtualDirectorio)
				numEstructuraTreeComplete = avd.AvdApArbolVirtualDirectorio
				numEstructuraTreeComplete--
				recorrerArbolRecursivo(file, sb, 1)
			}
		}
	} else if tipoArchivo == 2 { //CUANDO ES TIPO DETALLE DIRECTORIO
		dd := leerDD(file, sb.SbApDetalleDirectorio+(sb.SbSizeStructDetalleDirectorio*numEstructuraTreeComplete))
		for i := 0; i < 6; i++ {
			if i < 5 && dd.DdArrayFiles[i].DdFileApInodo != 0 {
				fmt.Println("[TIPO2]pos ", i, ":", dd.DdArrayFiles[i].DdFileApInodo)
				numEstructuraTreeComplete = dd.DdArrayFiles[i].DdFileApInodo
				numEstructuraTreeComplete--
				recorrerArbolRecursivo(file, sb, 3)
			} else if i == 5 && dd.DdApDetalleDirectorio != 0 {
				fmt.Println("[TIPO2]pos ", i, ":", dd.DdApDetalleDirectorio)
				numEstructuraTreeComplete = dd.DdApDetalleDirectorio
				numEstructuraTreeComplete--
				recorrerArbolRecursivo(file, sb, 2)
			}
		}
	} else if tipoArchivo == 3 { //CUANDO ES TIPO TABLA INODO
		ti := leerTABLAINODO(file, sb.SbApTablaInodo+(sb.SbSizeStructInodo*numEstructuraTreeComplete))
		for i := 0; i < 5; i++ {
			if i < 4 && ti.IArrayBloques[i] != 0 {
				numEstructuraTreeCompleteMov = ti.ICountInodo - 1
				fmt.Println("[TIPO3]pos ", i, ",", numEstructuraTreeCompleteMov, ":", ti.IArrayBloques[i])
				//fmt.Println("[TIPO3]pos ", i, ",", numEstructuraTreeComplete, ":", ti.IArrayBloques[i])
				numEstructuraTreeComplete = ti.IArrayBloques[i]
				numEstructuraTreeComplete--
				recorrerArbolRecursivo(file, sb, 4)
			} else if i == 4 && ti.IApIndirecto != 0 {
				fmt.Println("[TIPO3]pos ind ", i, ":", ti.IApIndirecto)
				numEstructuraTreeComplete = ti.IApIndirecto
				numEstructuraTreeComplete--
				recorrerArbolRecursivo(file, sb, 3)
			}
		}
	} else if tipoArchivo == 4 { //CUANDO ES TIPO BLOQUE DE DATOS
		bd := leerBLOQUEDATOS(file, sb.SbApBloques+(sb.SbSizeStructBloque*numEstructuraTreeComplete))
		fmt.Println("[TIPO4]pos: ", string(bd.DbDato[:]))
	}
}

//CadenaRetornoUserTXT =
var CadenaRetornoUserTXT string

func recorrerArbolRecursivoRetornarUsersTxt(file *os.File, sb SUPERBOOT, tipoArchivo int) {
	//Se empieza a recorrer desde el inodo del archivo USER.TXT que siempre sera el primer INODO
	if tipoArchivo == 3 { //CUANDO ES TIPO TABLA INODO
		ti := leerTABLAINODO(file, sb.SbApTablaInodo+(sb.SbSizeStructInodo*numEstructuraTreeComplete))
		for i := 0; i < 5; i++ {
			if i < 4 && ti.IArrayBloques[i] != 0 {
				numEstructuraTreeCompleteMov = ti.ICountInodo - 1
				fmt.Println("[TIPO3]pos ", i, ",", numEstructuraTreeCompleteMov, ":", ti.IArrayBloques[i])
				numEstructuraTreeComplete = ti.IArrayBloques[i]
				numEstructuraTreeComplete--
				recorrerArbolRecursivoRetornarUsersTxt(file, sb, 4)
			} else if i == 4 && ti.IApIndirecto != 0 {
				fmt.Println("[TIPO3]pos ind", i, ":", ti.IApIndirecto)
				numEstructuraTreeComplete = ti.IApIndirecto
				numEstructuraTreeComplete--
				recorrerArbolRecursivoRetornarUsersTxt(file, sb, 3)
			}
		}
	} else if tipoArchivo == 4 { //CUANDO ES TIPO BLOQUE DE DATOS
		bd := leerBLOQUEDATOS(file, sb.SbApBloques+(sb.SbSizeStructBloque*numEstructuraTreeComplete))
		for i := 0; i < len(bd.DbDato); i++ {
			if bd.DbDato[i] != 0 {
				CadenaRetornoUserTXT += string(bd.DbDato[i])
			}
		}
		//CadenaRetornoUserTXT += string(bd.DbDato[:])
		fmt.Println("[TIPO4]pos: ", string(bd.DbDato[:]))
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//USER.TXT-----USER.TXT-----USER.TXT-----USER.TXT-----USER.TXT-----USER.TXT-----USER.TXT-----USER.TXT-----USER.TXT-----USER.TXT-----USER.TXT-
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------

//modificarUSERTXT =
//	PARAMETRO 1 -> file: 		recibe el archivo
//	PARAMETRO 2 -> sb: 			una estructura SUPERBOOT
//	PARAMETRO 3 -> name:		recibe el nombre del GRUPO PARA USARLO EN GRUPO O USUARIO
//	PARAMETRO 3 -> usr: 		recibe el nombre del USUARIO PARA USARLO EN USUARIO
//	PARAMETRO 4 -> pwd: 		recibe el PASSWORD para el uso de la creacion del usuario
//	PARAMETRO 5 -> tamCadena:	variable int que en princio lleva el len de el arreglo que separa por '\n'
//  PARAMETRO 6 -> tipo:		1 indica que es grupo, 2 indica que es usuario, 3 y 4 indican que son remover
func modificarUSERTXT(file *os.File, sb SUPERBOOT, name, usr, pwd string, tamCadena, tipo int) {
	//********************************************************
	//Se carga el bitmap del inodo y bloques
	//********************************************************
	bitmapTablaInodos := retornarBitmap(file, sb.SbApBitmapTablaInodo, sb)
	bitmapBloques := retornarBitmap(file, sb.SbApBitmapBloques, sb)
	if tipo == 1 {
		CadenaRetornoUserTXT += "1,G," + name + "\\n"
	} else if tipo == 2 {
		CadenaRetornoUserTXT += "1,U," + name + "," + usr + "," + pwd + "\\n"
	}
	fmt.Println(cyan + "Cadena USERS.TXT Nueva: " + reset + CadenaRetornoUserTXT)
	tamCadena = len(CadenaRetornoUserTXT)
	//********************************************************
	//Se reinician los arreglos de apuntadores, los cuales se
	//encargaran de almacenar los apuntadores que estan en uso
	//para poder reescribirlos luego
	//********************************************************
	ApuntadoresBloqueUsoUSERTXT = nil
	ApuntadoresInodosUsoUSERTXT = nil
	ApuntadoresInodosUsoUSERTXT = append(ApuntadoresInodosUsoUSERTXT, 1) //Se agrega el 1 ya que este siempre esta
	numEstructuraTreeComplete = 0
	recorrerArbolRecursivoRetornarApuntadoresUSERTXT(file, sb, 3)
	fmt.Println(tamCadena)
	fmt.Println(bitmapTablaInodos)
	fmt.Println(bitmapBloques)
	//********************************************************
	//Se calculan cuantos bloques se usaran y se valida si hay
	//espacio suficiente en el bitmap
	//********************************************************
	totalDeBloques := tamCadena / 25.0
	restoBloques := tamCadena % 25.0
	if restoBloques != 0 {
		totalDeBloques++
	}
	totalDeBloquesNuevos := totalDeBloques - len(ApuntadoresBloqueUsoUSERTXT)
	fmt.Println("Bloques nuevos: ", totalDeBloquesNuevos)
	//var cuantosBloquesLibres int64
	//var posicionBloquesLibres []int64
	fmt.Println("Tamano Bitmap Bloques: ", len(bitmapBloques))
	fmt.Println("Arreglo posiciones Uso: ", ApuntadoresBloqueUsoUSERTXT)
	fmt.Println("Apuntadores bloques en uso: ", ApuntadoresBloqueUsoUSERTXT[:])
	fmt.Println("Apuntadores indodos en uso: ", ApuntadoresInodosUsoUSERTXT[:])

	var todasLasPosicionesAEscribirBloques []int64
	todasLasPosicionesAEscribirBloques = append(todasLasPosicionesAEscribirBloques[:], ApuntadoresBloqueUsoUSERTXT[:]...)
	contadorBloquesLibres := totalDeBloquesNuevos

	if len(bitmapBloques) < totalDeBloquesNuevos {
		fmt.Println(red + "[ERROR]" + reset + "No hay suficientes BLOQUES para insertar el grupo")
	} else {
		for i := 0; i < len(bitmapBloques); i++ {
			if contadorBloquesLibres == 0 {
				break
			}
			if bitmapBloques[i] == '0' {
				todasLasPosicionesAEscribirBloques = append(todasLasPosicionesAEscribirBloques[:], int64(i+1))
				contadorBloquesLibres--
			}
		}
		fmt.Println(magenta + "TODAS LAS POSICIONES DE BLOQUES: " + reset)
		fmt.Println(todasLasPosicionesAEscribirBloques)
	}

	//********************************************************
	//Se calculan cuantos inodos se usaran y se valida si hay
	//espacio suficiente en el bitmap
	//********************************************************

	totalDeInodos := totalDeBloques / 4
	restoInodos := totalDeBloques % 4
	if restoInodos != 0 {
		totalDeInodos++
	}
	totalDeInodosNuevos := totalDeInodos - len(ApuntadoresInodosUsoUSERTXT)
	fmt.Println("Total de inodos nuevos: ", totalDeInodosNuevos)

	var todasLasPosicionesAEscribirInodos []int64
	todasLasPosicionesAEscribirInodos = append(todasLasPosicionesAEscribirInodos[:], ApuntadoresInodosUsoUSERTXT[:]...)
	contadorInodosLibres := totalDeInodosNuevos

	if len(bitmapTablaInodos) < totalDeInodosNuevos {
		fmt.Println(red + "[ERROR]" + reset + "No hay suficientes INODOS para insertar el grupo")
	} else {
		for i := 0; i < len(bitmapTablaInodos); i++ {
			if contadorInodosLibres == 0 {
				break
			}
			if bitmapTablaInodos[i] == '0' {
				todasLasPosicionesAEscribirInodos = append(todasLasPosicionesAEscribirInodos[:], int64(i+1))
				contadorInodosLibres--
			}
		}
		fmt.Println(magenta + "TODAS LAS POSICIONES DE INODOS: " + reset)
		fmt.Println(todasLasPosicionesAEscribirInodos)
	}

	//ti := leerTABLAINODO(file, sb.SbApTablaInodo+(sb.SbSizeStructInodo*numEstructuraTreeComplete))
	//********************************************************
	//Se escriben las nuevas estructuras
	//********************************************************
	fmt.Println("T: ", totalDeBloques)
	contPosBloque := totalDeBloques
	numPosBloque := 0
	for i := 0; i < len(todasLasPosicionesAEscribirInodos); i++ {
		var nuevoTI TABLAINODO
		nuevoTI.ICountInodo = todasLasPosicionesAEscribirInodos[i]
		nuevoTI.ISizeArchivo = int64(totalDeBloques) * sb.SbSizeStructBloque
		nuevoTI.ICountBloquesAsignados = int64(totalDeBloques)
		for j := 0; j < 4; j++ {
			if contPosBloque != 0 {
				nuevoTI.IArrayBloques[j] = todasLasPosicionesAEscribirBloques[numPosBloque]
				numPosBloque++
				contPosBloque--
			}
		}
		if contPosBloque != 0 {
			nuevoTI.IApIndirecto = todasLasPosicionesAEscribirInodos[i+1]
		}
		nuevoTI.IIdProper[0] = 'r'
		nuevoTI.IIdProper[1] = 'o'
		nuevoTI.IIdProper[2] = 'o'
		nuevoTI.IIdProper[3] = 't'
		file.Seek(sb.SbApTablaInodo+(sb.SbSizeStructInodo*(todasLasPosicionesAEscribirInodos[i]-1)), 0)
		var valorBinarioTablaInodo bytes.Buffer
		binary.Write(&valorBinarioTablaInodo, binary.BigEndian, &nuevoTI)
		escribirBytes(file, valorBinarioTablaInodo.Bytes())
	}

	cadenaUserTxtSeparada := SplitSubN(CadenaRetornoUserTXT, 25)
	//fmt.Println("mamarre VAMOS A BLOQUES")
	//fmt.Println("Longitud: ", len(todasLasPosicionesAEscribirBloques))

	for i := 0; i < len(todasLasPosicionesAEscribirBloques); i++ {
		var nuevoBD BLOQUEDATOS
		var contenidoAByte25 [25]byte
		copy(contenidoAByte25[:], cadenaUserTxtSeparada[i])
		nuevoBD.DbDato = contenidoAByte25
		//CadenaRetornoUserTXT
		fmt.Println(string(contenidoAByte25[:]))
		fmt.Println(todasLasPosicionesAEscribirBloques[i])
		file.Seek(sb.SbApBloques+(sb.SbSizeStructBloque*(todasLasPosicionesAEscribirBloques[i]-1)), 0)
		var valorBinarioBloque bytes.Buffer
		binary.Write(&valorBinarioBloque, binary.BigEndian, &nuevoBD)
		escribirBytes(file, valorBinarioBloque.Bytes())
	}

	//********************************************************
	//Se escriben las nuevas estructuras
	//********************************************************
	nuevoBitmapTablaInodos := bitmapTablaInodos
	for i := 0; i < len(todasLasPosicionesAEscribirInodos); i++ {
		nuevoBitmapTablaInodos[todasLasPosicionesAEscribirInodos[i]-1] = '1'
	}
	fmt.Println(red, nuevoBitmapTablaInodos, reset)
	reescribirBitmap(file, sb.SbApBitmapTablaInodo, nuevoBitmapTablaInodos)
	nuevoBitmapBloques := bitmapBloques
	for i := 0; i < len(todasLasPosicionesAEscribirBloques); i++ {
		nuevoBitmapBloques[todasLasPosicionesAEscribirBloques[i]-1] = '1'
	}
	reescribirBitmap(file, sb.SbApBitmapBloques, nuevoBitmapBloques)
}

//ApuntadoresBloqueUsoUSERTXT = Arreglo que almacena los apuntadores de tipo bloque que usa el archivo txt
var ApuntadoresBloqueUsoUSERTXT []int64

//ApuntadoresInodosUsoUSERTXT = Arreglo que almacena los apuntadores de tipo Inodo que usa el archivo txt
var ApuntadoresInodosUsoUSERTXT []int64

func recorrerArbolRecursivoRetornarApuntadoresUSERTXT(file *os.File, sb SUPERBOOT, tipoArchivo int) {
	//Se empieza a recorrer desde el inodo del archivo USER.TXT que siempre sera el primer INODO
	if tipoArchivo == 3 { //CUANDO ES TIPO TABLA INODO
		ti := leerTABLAINODO(file, sb.SbApTablaInodo+(sb.SbSizeStructInodo*numEstructuraTreeComplete))
		for i := 0; i < 5; i++ {
			if i < 4 && ti.IArrayBloques[i] != 0 {
				//fmt.Println("[TIPO3]pos ", i, ",", numEstructuraTreeComplete, ":", ti.IArrayBloques[i])
				numEstructuraTreeComplete = ti.IArrayBloques[i]
				ApuntadoresBloqueUsoUSERTXT = append(ApuntadoresBloqueUsoUSERTXT, ti.IArrayBloques[i])
				numEstructuraTreeComplete--
				//recorrerArbolRecursivoRetornarApuntadoresUSERTXT(file, sb, 4)
			} else if i == 4 && ti.IApIndirecto != 0 {
				//fmt.Println("[TIPO3]pos ", i, ":", ti.IArrayBloques[i])
				numEstructuraTreeComplete = ti.IApIndirecto
				ApuntadoresInodosUsoUSERTXT = append(ApuntadoresInodosUsoUSERTXT, ti.IApIndirecto)
				numEstructuraTreeComplete--
				recorrerArbolRecursivoRetornarApuntadoresUSERTXT(file, sb, 3)
			}
		}
	} else if tipoArchivo == 4 { //CUANDO ES TIPO BLOQUE DE DATOS
		//bd := leerBLOQUEDATOS(file, sb.SbApBloques+(sb.SbSizeStructBloque*numEstructuraTreeComplete))
		//fmt.Println("[TIPO4]pos: ", string(bd.DbDato[:]))
	}
}

func recorrerArbolRecursivoReportePorBitmap(file *os.File, sb SUPERBOOT, posicion int64, tipoArchivo int) {
	if tipoArchivo == 1 {
		avd := leerAVD(file, sb.SbApArbolDirectorio+(sb.SbSizeStructArbolDirectorio*posicion))
		subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(posicion+1, 10) + "[label=<\n<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"
		var nombreDirectorio string
		for i1, valor1 := range avd.AvdNombreDirectorio {
			if avd.AvdNombreDirectorio[i1] != 0 {
				nombreDirectorio += string(valor1)
			}
		}
		subCadenaReporteTreeComplete += "<TR port='0'>\n<TD BGCOLOR='#99ccff' COLSPAN='2'><font color='black'>" + nombreDirectorio + "</font></TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>Fecha Creacion</TD><TD>" + string(avd.AvdFechaCreacion[:]) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 1</TD><TD port='1'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[0], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 2</TD><TD port='2'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[1], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 3</TD><TD port='3'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[2], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 4</TD><TD port='4'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[3], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 5</TD><TD port='5'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[4], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APD 6</TD><TD port='6'>" + strconv.FormatInt(avd.AvdApArraySubdirectorios[5], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>Detalle Directorio</TD><TD port='7'>" + strconv.FormatInt(avd.AvdApDetalleDirectorio, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>APDI 1</TD><TD port='8'>" + strconv.FormatInt(avd.AvdApArbolVirtualDirectorio, 10) + "</TD>\n</TR>\n"
		var nombreProper string
		for i1, valor1 := range avd.AvdProper {
			if avd.AvdProper[i1] != 0 {
				nombreProper += string(valor1)
			}
		}
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#99ccff'>PROPER</TD><TD>" + nombreProper + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "</TABLE>\n>];\n\n"
		for i := 0; i < 8; i++ {
			if i < 6 && avd.AvdApArraySubdirectorios[i] != 0 {
				subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(posicion+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(avd.AvdApArraySubdirectorios[i], 10) + "\n"
			} else if i == 6 && avd.AvdApDetalleDirectorio != 0 {
				fmt.Println("[TIPO1]pos det", i, ":", avd.AvdApDetalleDirectorio)
				subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(posicion+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				subCadenaReporteTreeComplete += "DD" + strconv.FormatInt(avd.AvdApDetalleDirectorio, 10) + "\n"
				numEstructuraTreeComplete = avd.AvdApDetalleDirectorio
			} else if i == 7 && avd.AvdApArbolVirtualDirectorio != 0 {
				fmt.Println("[TIPO1]pos ind", i, ":", avd.AvdApArbolVirtualDirectorio)
				subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(posicion+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				numEstructuraTreeComplete = avd.AvdApArbolVirtualDirectorio
				subCadenaReporteTreeComplete += "AVD" + strconv.FormatInt(avd.AvdApArbolVirtualDirectorio, 10) + "\n"
			}
		}
	} else if tipoArchivo == 2 {
		//fmt.Println(posicion)
		dd := leerDD(file, sb.SbApDetalleDirectorio+(sb.SbSizeStructDetalleDirectorio*posicion))
		subCadenaReporteTreeComplete += "\nDD" + strconv.FormatInt(posicion+1, 10) + "[label=<\n<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#a3d977'>DETALLE<br/>DIRECTORIO</TD><TD BGCOLOR='#a3d977'>" + strconv.FormatInt(posicion+1, 10) + "</TD>\n</TR>"
		for i := 0; i < 5; i++ {
			if dd.DdArrayFiles[i].DdFileApInodo != 0 {
				//fmt.Println("mamarre")
				var nombreDirectorio string
				for i1, valor1 := range dd.DdArrayFiles[i].DdFileNombre {
					if dd.DdArrayFiles[i].DdFileNombre[i1] != 0 {
						nombreDirectorio += string(valor1)
					}
				}
				subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#a3d977'>" + nombreDirectorio + "</TD><TD port='" + strconv.Itoa(i+1) + "'>" + strconv.FormatInt(dd.DdArrayFiles[i].DdFileApInodo, 10) + "</TD>\n</TR>\n"
			} else {
				subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#a3d977'>APD" + strconv.Itoa(i+1) + "</TD><TD port='" + strconv.Itoa(i+1) + "'>0</TD>\n</TR>\n"
			}
		}
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#a3d977'>API1</TD><TD port='6'>" + strconv.FormatInt(dd.DdApDetalleDirectorio, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "</TABLE>\n>];\n\n"
		for i := 0; i < 6; i++ {
			if i < 5 && dd.DdArrayFiles[i].DdFileApInodo != 0 {
				subCadenaReporteTreeComplete += "DD" + strconv.FormatInt(posicion+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				subCadenaReporteTreeComplete += "TI" + strconv.FormatInt(dd.DdArrayFiles[i].DdFileApInodo, 10) + "\n"
			} else if i == 5 && dd.DdApDetalleDirectorio != 0 {
				fmt.Println("[TIPO2]pos ", i, ":", dd.DdApDetalleDirectorio)
				subCadenaReporteTreeComplete += "DD" + strconv.FormatInt(posicion+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				subCadenaReporteTreeComplete += "DD" + strconv.FormatInt(dd.DdApDetalleDirectorio, 10) + "\n"
			}
		}
	} else if tipoArchivo == 3 {
		ti := leerTABLAINODO(file, sb.SbApTablaInodo+(sb.SbSizeStructInodo*posicion))
		subCadenaReporteTreeComplete += "\nTI" + strconv.FormatInt(posicion+1, 10) + " [label=<\n<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>TABLA INODO</TD><TD BGCOLOR='#ffc374'>" + strconv.FormatInt(posicion+1, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>Tamaño</TD><TD>" + strconv.FormatInt(ti.ISizeArchivo, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>Bloques</TD><TD>" + strconv.FormatInt(ti.ICountBloquesAsignados, 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>APD 1</TD><TD port='1'>" + strconv.FormatInt(ti.IArrayBloques[0], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>APD 2</TD><TD port='2'>" + strconv.FormatInt(ti.IArrayBloques[1], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>APD 3</TD><TD port='3'>" + strconv.FormatInt(ti.IArrayBloques[2], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>APD 4</TD><TD port='4'>" + strconv.FormatInt(ti.IArrayBloques[3], 10) + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>API</TD><TD port='5'>" + strconv.FormatInt(ti.IApIndirecto, 10) + "</TD>\n</TR>\n"
		var nombreProper string
		for i1, valor1 := range ti.IIdProper {
			if ti.IIdProper[i1] != 0 {
				nombreProper += string(valor1)
			}
		}
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ffc374'>Proper</TD><TD>" + nombreProper + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "</TABLE>\n>];\n\n"
		for i := 0; i < 5; i++ {
			if i < 4 && ti.IArrayBloques[i] != 0 {
				subCadenaReporteTreeComplete += "TI" + strconv.FormatInt(posicion+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				subCadenaReporteTreeComplete += "B" + strconv.FormatInt(ti.IArrayBloques[i], 10) + "\n"
			} else if i == 4 && ti.IApIndirecto != 0 {
				subCadenaReporteTreeComplete += "TI" + strconv.FormatInt(posicion+1, 10) + ":" + strconv.Itoa(i+1) + " -> "
				subCadenaReporteTreeComplete += "TI" + strconv.FormatInt(ti.IApIndirecto, 10) + "\n"
			}
		}
	} else if tipoArchivo == 4 {
		bd := leerBLOQUEDATOS(file, sb.SbApBloques+(sb.SbSizeStructBloque*posicion))
		subCadenaReporteTreeComplete += "\nB" + strconv.FormatInt(posicion+1, 10) + "[label=<\n\n<TABLE BORDER='0' CELLBORDER='1' CELLSPACING='0'>\n"
		subCadenaReporteTreeComplete += "<TR>\n<TD BGCOLOR='#ff8f80'>BLOQUE</TD><TD BGCOLOR='#ff8f80'>" + strconv.FormatInt(posicion+1, 10) + "</TD>\n</TR>\n"
		var contenidoBloque string
		for i1, valor1 := range bd.DbDato {
			if bd.DbDato[i1] != 0 {
				contenidoBloque += string(valor1)
			}
		}
		subCadenaReporteTreeComplete += "<TR>\n<TD COLSPAN='2'>" + contenidoBloque + "</TD>\n</TR>\n"
		subCadenaReporteTreeComplete += "</TABLE>\n>];\n\n"
	}
}
