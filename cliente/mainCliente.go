package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Uso: cliente <IP> <Puerto> <PeriodoReporte>")
		return
	}
	IPserver := os.Args[1]
	puerto := os.Args[2]
	periodo := os.Args[3]

	fmt.Println("***********************************")
	fmt.Println("****** Cliente Windows ************")
	fmt.Println("***********************************")

	dirServerTCP, _ := net.ResolveTCPAddr("tcp4", IPserver+":"+puerto)
	socketCLiente, err := net.DialTCP("tcp", nil, dirServerTCP)
	if err != nil {
		fmt.Println("Error al conectarse con el server", err)
		return
	}
	defer socketCLiente.Close()
	fmt.Println("Conectado al servidor:", socketCLiente.RemoteAddr())

	teclado := bufio.NewReader(os.Stdin)
	bufferS := bufio.NewReader(socketCLiente)
	bufferC := bufio.NewWriter(socketCLiente)

	for {
		fmt.Print("usuario: ")
		user, _ := teclado.ReadString('\n')
		fmt.Print("Password:")
		pass, _ := teclado.ReadString('\n')
		bufferC.WriteString(strings.TrimSpace(user) + ":" + strings.TrimSpace(pass) + "\n")
		bufferC.Flush()

		resp, err := bufferS.ReadString('\n')
		if err != nil {
			fmt.Println("la conexion ha sido cerrada por el servidor, por muchos intentos fallidos")
			return
		}
		if strings.TrimSpace(resp) == "Recibido" {
			fmt.Println("Login correcto ;)")
			break
		}
		fmt.Println("Errorrr! Contraseña o usuer incorrecto:()")
	}

	bufferC.WriteString(periodo + "\n")
	bufferC.Flush()

	salir := make(chan bool)
	go envMensajeC(socketCLiente, teclado, salir)
	go recMensajeC(bufferS, salir)

	<-salir
	fmt.Println("Cliente cerrado")
}

func recMensajeC(bufferS *bufio.Reader, salir chan bool) {
	for {
		msjRec, err := bufferS.ReadString('\n')
		if err != nil {
			salir <- true
			return
		}
		msjRec = strings.TrimSpace(msjRec)
		if strings.HasPrefix(msjRec, "REP|") {
			fmt.Println("\n [REPORTE DEL SISTEMA]", strings.TrimPrefix(msjRec, "REP|"))
		} else if strings.HasPrefix(msjRec, "CMD|") {
			salida := strings.ReplaceAll(strings.TrimPrefix(msjRec, "CMD|"), " | ", "\n")
			fmt.Println("\n[Respuesta al comando]:\n" + salida)
		}
		fmt.Print("Digite el comando que desea enviar:")
	}
}

func envMensajeC(socketCLiente net.Conn, teclado *bufio.Reader, salir chan bool) {
	bufferC := bufio.NewWriter(socketCLiente)
	for {
		fmt.Print("Digite el comando a enviar: ")
		comando, _ := teclado.ReadString('\n')

		bufferC.WriteString(comando)
		bufferC.Flush()

		if strings.TrimSpace(comando) == "bye" {
			salir <- true
			return
		}
	}
}
