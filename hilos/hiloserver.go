package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
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

	salir := make(chan bool)
	go recMensajeS(&socketS, salir) //go vuelve la funcion un hilo
	go envMensajeS(&socketS)

	<-salir //tiene que haber un canal que lea cuando recMensajes vea que llga bye, cierre el proceso

	socketS.Close()
	fmt.Println("Servidor cerrado, Adiós")
}
func recMensajeS(socketS *net.Conn, salir chan bool) {
	bufferS := bufio.NewReader(*socketS) //conectamos el lector al socket

	for {
		msjRec, _ := bufferS.ReadString('\n')
		fmt.Println("Comando recibido: ", msjRec)
		if strings.TrimSpace(msjRec) == "bye" {
			salir <- true
			break
		}

	}

}

func envMensajeS(socketS *net.Conn) {
	num := 0

	for {
		num++
		time.Sleep(time.Second * 8)
		bufferS := bufio.NewWriter(*socketS)
		bufferS.WriteString("REPORTE(" + strconv.Itoa(num) + ")[DD:17%  -CPU:9%   -RAM:34%)]\n") //ENVIA la respueta al comando o incluso un reporte
		bufferS.Flush()
		fmt.Println("Reporte enviado C: !")
		fmt.Println("======================\n")
	}

}
