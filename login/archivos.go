package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	nomArchivo := "oper.users.db"
	fmt.Print("ARCHIVOS OPER\n1.crear\n2.leer\n3.modificar\n4.salir\nDigite Opcion: ")
	var Opcion int
	fmt.Scanf("%d", &Opcion) //lee las opciones por teclado
	switch Opcion {
	case 1:
		crearArchivo(nomArchivo)
	case 2:
		leerArchivo(nomArchivo)
	case 3:
		modiArchivo(nomArchivo)
	case 4:
		fmt.Println("Saliendo...Adióss")
	default:
		fmt.Println("Opcion invalida  :(")
	}
	//crearArchivo(nomArchivo)
	//leerArchivo(nomArchivo)
	//modiArchivo(nomArchivo)

}
func crearArchivo(nomArchivo string) {
	credUser := "\nmramirez:caf90169eefa5f807d577486b9f795ab86ae2983c5c20806cff959117e90af18"
	os.WriteFile(nomArchivo, []byte(credUser), 0777)
}
func leerArchivo(nomArchivo string) {
	archivo, _ := os.ReadFile(nomArchivo)
	strArchivo := string(archivo)
	//fmt.Println(strArchivo)
	usuarios := strings.Split(strArchivo, "\n") //split devuelve iun vestoy tipo byte sin tamaño predefinido
	//usuarios es un vector que lee los usuarios de acuerdo a sus saltos d línea

	for i := 0; i < len(usuarios); i++ {
		fmt.Println("Usuario[", i, "] = ", usuarios[i])
	}

}
func modiArchivo(nomArchivo string) {
	nuevoUsuario := "luis :340d600392818df2413382dc7d8325c360d83ea49a262d31760348484bbc10b5" //windows
	//nuevoUsuario := "\nmaria :94aec9fbed989ece189a7e172c9cf41669050495152bc4c1dbf2a38d7fd85627" //maria
	//nuevoUsuario := "\njose:1ec4ed037766aa181d8840ad04b9fc6e195fd37dedc04c98a5767a67d3758ece" //jose
	archivo, _ := os.OpenFile(nomArchivo, os.O_APPEND|os.O_WRONLY, 0777)
	archivo.WriteString(nuevoUsuario)
	archivo.Close()
}

//hacer un cicli anidado para ñeer los nombres y las credenciales
//diseñar la estructura de archivos, crear carpeta en ec, llamada ope que a a llevar la informacion de archivos
//ejercicios complementario, poner un swicth case para elegir que
//para el proyecto se le debe dar seguridad (KMODO DE ACCESO a esos archivos)
