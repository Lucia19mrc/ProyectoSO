package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("============================")
	fmt.Println("===CLIENTE OPER ============")
	fmt.Println("============================")

	IPServer := os.Args[1]
	dirServerTCP, _ := net.ResolveTCPAddr("tcp4", IPServer+":2710") //
	var socketCliente net.Conn
	socketCliente, _ = net.DialTCP("tcp", nil, dirServerTCP) //puntero de tipo tcp que seale a ese puerto
	fmt.Println("Conectado al server: ", socketCliente.RemoteAddr())

	bufferC := bufio.NewWriter(socketCliente)
	bufferC.WriteString("ls -l desde mariarami...\n")
	bufferC.Flush()
	fmt.Println("Comando enviado!")
	for {
		envMensajeC(&socketCliente)
		recMensajeC(&socketCliente)
	}

	socketCliente.Close()
}
func recMensajeC(socketCliente *net.Conn) {
	bufferS := bufio.NewReader(*socketCliente)
	msjRec, _ := bufferS.ReadString('\n') //lee hasta el salto de línea
	fmt.Println("Reporte recibido: ", msjRec)
	fmt.Println("========================\n")
}

func envMensajeC(socketCliente *net.Conn) {
	//del archivo shell.go
	fmt.Print("Digite el comando a Enviar: ")
	lector := bufio.NewReader(os.Stdin)
	comando, _ := lector.ReadString('\n')
	//comando = strings.TrimRight(comando, "\n")

	bufferC := bufio.NewWriter(*socketCliente)
	bufferC.WriteString(comando)
	bufferC.Flush()
	fmt.Println("Comando enviado!")
}
