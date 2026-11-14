package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("****************************")
	fmt.Println("***SERVER OPER**************")
	fmt.Println("*************************")

	dirTCP, _ := net.ResolveTCPAddr("tcp4", ":2710") //así atiende  una solo direccion Ip,
	// si lo dejo vavio el seguno espacio va a todas las ips
	oyente, _ := net.ListenTCP("tcp", dirTCP) //puntero de tipo tcp que seale a ese puerto
	fmt.Println("Esperando Conexiones !!")

	socketS, _ := oyente.Accept()
	fmt.Println("Cliente COnectado", socketS.RemoteAddr())

	bufferS := bufio.NewReader(socketS)
	msjRec, _ := bufferS.ReadString('\n')
	fmt.Println("Comando recibido: ", msjRec)
	for {
		recMensajeS(&socketS)
		envMensajeS(&socketS)

	}

	socketS.Close()
}
func recMensajeS(socketS *net.Conn) {
	bufferS := bufio.NewReader(*socketS) //conectamos el lector al socket
	msjRec, _ := bufferS.ReadString('\n')
	fmt.Println("Comando recibido: ", msjRec)
}

func envMensajeS(socketS *net.Conn) {
	time.Sleep(time.Second * 5)
	bufferS := bufio.NewWriter(*socketS)
	bufferS.WriteString("REPORTE [DD:17%  -CPU:9%   -RAM:34%)]\n") //ENVIA la respueta al comando o incluso un reporte
	bufferS.Flush()
	fmt.Println("Reporte enviado C: !")
	fmt.Println("======================\n")
}
