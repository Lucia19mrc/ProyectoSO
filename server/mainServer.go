package main

import (
	"bufio"
	"crypto/sha512"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("******************************")
	fmt.Println("****  SERVER DEBIAN       *****")
	fmt.Println("*******************************")
	registrarLog("=== INICIANDO SERVER DEBIAN ===")

	// 1. Leer configuración desde /etc/configProyecto.conf
	ipPermitida, puerto, maxintentos, usuariospermitidos, archivousuaris := leerConfig("/etc/configProyecto.conf")

	dirTCP, _ := net.ResolveTCPAddr("tcp4", ":"+puerto)
	oyente, err := net.ListenTCP("tcp", dirTCP)

	if err != nil {
		fmt.Println("Error al iniciar el server: ", err)
		registrarLog("Error al iniciar el server: " + err.Error())
		return
	}
	defer oyente.Close()

	// 2. Bucle infinito para que el servicio systemd siga escuchando después de cada cliente
	for {
		fmt.Println("Esperando conexiones en el puerto", puerto, "...")
		socketS, err := oyente.Accept()
		if err != nil {
			continue
		}

		ipCLient, _, _ := net.SplitHostPort(socketS.RemoteAddr().String())
		if ipCLient != ipPermitida {
			fmt.Println("Ip no auotrizada segun config:", ipCLient)
			registrarLog("IP no autorizada rechazada: " + ipCLient)
			socketS.Close()
			continue
		}
		fmt.Println("Cliente conectados desde Ip autorizada:", ipCLient)
		registrarLog("Cliente conectado desde IP autorizada: " + ipCLient)

		bufferS := bufio.NewReader(socketS)
		writerS := bufio.NewWriter(socketS)

		autenticacion := false
		for intento := 1; intento <= maxintentos; intento++ {
			linea, err := bufferS.ReadString('\n')
			if err != nil {
				break
			}
			datos := strings.Split(strings.TrimSpace(linea), ":")
			if len(datos) == 2 {
				if validarLogin(datos[0], datos[1], usuariospermitidos, archivousuaris) {
					writerS.WriteString("Recibido\n")
					writerS.Flush()
					fmt.Println("Usuaio autenticado:", datos[0])
					registrarLog("Usuario autenticado: " + datos[0])
					autenticacion = true
					break
				}
			}
			writerS.WriteString("Error\n")
			writerS.Flush()
			fmt.Println("Intento fallido :(", intento, "de", maxintentos)
			registrarLog(fmt.Sprintf("Intento fallido (%d/%d) desde %s", intento, maxintentos, ipCLient))
		}

		if !autenticacion {
			fmt.Println("3 intentos falllidos, se va a cerrar la conexión")
			registrarLog("Conexión cerrada por intentos fallidos: " + ipCLient)
			socketS.Close()
			continue
		}

		lineaPeriodo, _ := bufferS.ReadString('\n')
		periodo, err := strconv.Atoi(strings.TrimSpace(lineaPeriodo))
		if err != nil || periodo <= 0 {
			periodo = 5
		}
		fmt.Println("ENviando reportes cada ", periodo, "segundos")
		registrarLog(fmt.Sprintf("Monitoreo iniciado cada %d segundos", periodo))

		salir := make(chan bool)

		go recMensajeS(&socketS, bufferS, salir)
		go envReporteS(&socketS, periodo, salir)

		<-salir
		socketS.Close()
		fmt.Println("Sesión de cliente cerrada")
		registrarLog("Sesión de cliente cerrada")
	}
}

// 4. Función sencilla para guardar eventos en archivo de log
func registrarLog(mensaje string) {
	f, err := os.OpenFile("/var/log/proyectoOper/actividad.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	fecha := time.Now().Format("2006-01-02 15:04:05")
	f.WriteString("[" + fecha + "] " + mensaje + "\n")
}

func leerConfig(nombreArchivo string) (string, string, int, []string, string) {
	datos, err := os.ReadFile(nombreArchivo)
	if err != nil {
		fmt.Println("No se pudo leer", nombreArchivo)
		os.Exit(1)
	}
	var ip, puerto, archivousuaris string
	var maxintentos int
	var Usuaio []string
	lineas := strings.Split(string(datos), "\n")
	for _, linea := range lineas {
		partes := strings.Split(strings.TrimSpace(linea), "=")
		if len(partes) == 2 {
			switch partes[0] {
			case "ip_permitida":
				ip = partes[1]
			case "puerto":
				puerto = partes[1]
			case "max_intentos":
				maxintentos, _ = strconv.Atoi(partes[1])
			case "usuarios":
				Usuaio = strings.Split(partes[1], ",")
			case "archivo_usuarios":
				archivousuaris = partes[1]
			}
		}
	}
	return ip, puerto, maxintentos, Usuaio, archivousuaris
}

// 3. Validación SHA-512 con salt (formato en oper.users.db -> usuario:salt:hash)
func validarLogin(user string, pass string, permitidos []string, rutaDB string) bool {
	estaenConf := false
	for _, u := range permitidos {
		if strings.TrimSpace(u) == user {
			estaenConf = true
			break
		}
	}
	if !estaenConf {
		return false
	}

	contenido, err := os.ReadFile(rutaDB)
	if err != nil {
		fmt.Println("Error leyendo archivo de usuarios:", err)
		return false
	}
	lineas := strings.Split(string(contenido), "\n")
	for _, l := range lineas {
		campos := strings.Split(strings.TrimSpace(l), ":")
		if len(campos) == 3 && campos[0] == user {
			salt := campos[1]
			hashGuardado := campos[2]

			hash := sha512.Sum512([]byte(salt + pass))
			passHash := fmt.Sprintf("%x", hash)

			if passHash == hashGuardado {
				return true
			}
		}
	}
	return false
}

func recMensajeS(socketS *net.Conn, bufferS *bufio.Reader, salir chan bool) {
	writerS := bufio.NewWriter(*socketS)
	// Lista sencilla de comandos prohibidos
	prohibidos := []string{"rm", "rmdir", "shutdown", "reboot", "poweroff", "mkfs", "dd", "sudo", "su", "chmod", "chown"}

	for {
		msjRec, err := bufferS.ReadString('\n')
		if err != nil {
			salir <- true
			return
		}
		comando := strings.TrimSpace(msjRec)
		fmt.Println("Comando recibido:", comando)
		if comando == "bye" {
			registrarLog("Comando ejecutado: bye")
			salir <- true
			return
		}
		// Revisar si la primera palabra del comando está en la lista negra
		bloqueado := false
		palabras := strings.Fields(comando)
		if len(palabras) > 0 {
			for _, p := range prohibidos {
				if palabras[0] == p {
					bloqueado = true
					break
				}
			}
		}
		if bloqueado {
			fmt.Println("Comando bloqueado:", comando)
			registrarLog("COMANDO BLOQUEADO: " + comando)
			writerS.WriteString("CMD|Error: Comando no permitido por seguridad\n")
			writerS.Flush()
			continue
		}
		registrarLog("Comando ejecutado: " + comando)
		shell := exec.Command("bash", "-c", comando)
		salida, _ := shell.CombinedOutput()

		texto := strings.TrimSpace(string(salida))
		if texto == "" {
			texto = "(Comando ejecutado sin salida)"
		}
		texto = strings.ReplaceAll(texto, "\n", " | ")

		writerS.WriteString("CMD|" + texto + "\n")
		writerS.Flush()
	}
}
func envReporteS(socketS *net.Conn, periodo int, salir chan bool) {
	writerS := bufio.NewWriter(*socketS)
	num := 0
	for {
		select {
		case <-salir:
			return
		default:
			time.Sleep(time.Duration(periodo) * time.Second)
			num++

			//reporte del sistema
			cpu := ejecutarShell("top -bn1 | grep 'Cpu' | awk '{print 100 - $8\"%\"}'")
			ram := ejecutarShell("free -m | awk 'NR==2{printf \"%s/%sMB (%.1f%%)\", $3,$2,$3*100/$2}'")
			disco := ejecutarShell("df -h / | awk 'NR==2{printf \"%s/%s (%s)\", $3,$2,$5}'")
			procesos := ejecutarShell("ps -e --no-headers | wc -l")

			reporte := fmt.Sprintf("REP|(%d) CPU: %s | RAM: %s | DISCO: %s | PROCESOS: %s\n",
				num, cpu, ram, disco, procesos)

			_, err := writerS.WriteString(reporte)
			if err != nil {
				salir <- true
				return
			}
			writerS.Flush()
		}
	}
}

func ejecutarShell(cmd string) string {
	out, _ := exec.Command("bash", "-c", cmd).Output()
	res := strings.TrimSpace(string(out))
	if res == "" {
		return "0%"
	}
	return res
}
