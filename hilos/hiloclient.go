package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("============================")
	fmt.Println("===CLIENTE OPER ============")
	fmt.Println("============================")

	IPServer := os.Args[1]
	puertoServer := os.Args[2] // Puerto recibido por parámetro (ej: 2306)
	// periodoReporte := os.Args[3] // Periodo en segundos (ej: 5)

	dirServerTCP, _ := net.ResolveTCPAddr("tcp4", IPServer+":"+puertoServer) //
	var socketCliente net.Conn
	socketCliente, _ = net.DialTCP("tcp", nil, dirServerTCP) //puntero de tipo tcp que seale a ese puerto
	fmt.Println("Conectado al server: ", socketCliente.RemoteAddr())

	bufferC := bufio.NewWriter(socketCliente)
	bufferC.WriteString("ls -l desde mariarami...\n")
	bufferC.Flush()
	fmt.Println("Comando enviado!")

	salir := make(chan bool) // canal para esperar sin quemar CPU

	go envMensajeC(&socketCliente, salir)
	go recMensajeC(&socketCliente)

	<-salir // se queda esperando aquí hasta que escribas "bye"
	socketCliente.Close()
}
func recMensajeC(socketCliente *net.Conn) {
	bufferS := bufio.NewReader(*socketCliente) // se crea UNA sola vez fuera del for
	for {
		msjRec, _ := bufferS.ReadString('\n') //lee hasta el salto de línea
		fmt.Println("\nReporte recibido: ", msjRec)
		fmt.Println("========================\n")
	}

}

func envMensajeC(socketCliente *net.Conn, salir chan bool) {
	lector := bufio.NewReader(os.Stdin) // se crea UNA sola vez fuera del for,del archivo shell.go
	for {
		fmt.Print("Digite el comando a Enviar: ")
		comando, _ := lector.ReadString('\n')
		//comando = strings.TrimRight(comando, "\n")

		bufferC := bufio.NewWriter(*socketCliente)
		bufferC.WriteString(comando)
		bufferC.Flush()
		fmt.Println("Comando enviado!")

		if strings.TrimSpace(comando) == "bye" {
			salir <- true // avisa al main para cerrar el socket y terminar
			break
		}
	}
}
