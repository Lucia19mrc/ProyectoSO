package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Uso: proyOperGUI <IP> <Puerto> <PeriodoReporte>")
		return
	}
	IPserver := os.Args[1]
	puerto := os.Args[2]
	periodo := os.Args[3]

	// 1. Conectar al servidor igual que en mainCliente.go
	dirServerTCP, _ := net.ResolveTCPAddr("tcp4", IPserver+":"+puerto)
	socketCliente, err := net.DialTCP("tcp", nil, dirServerTCP)
	if err != nil {
		fmt.Println("Error al conectarse con el server:", err)
		return
	}
	defer socketCliente.Close()

	bufferS := bufio.NewReader(socketCliente)
	bufferC := bufio.NewWriter(socketCliente)

	// 2. Crear ventana Fyne
	prog := app.New()
	ventPpal := prog.NewWindow("Proyecto Sistemas Operativos")

	// --- ELEMENTOS DE LOGIN ---
	entryUser := widget.NewEntry()
	entryUser.SetPlaceHolder("Usuario")

	entryPass := widget.NewPasswordEntry()
	entryPass.SetPlaceHolder("Password")

	labelEstado := widget.NewLabel("Ingrese usuario y contraseña")

	// --- ELEMENTOS DEL PANEL PRINCIPAL ---
	textComando := widget.NewEntry()
	textComando.SetPlaceHolder("Digite el comando a enviar (o bye)...")

	labelRes := widget.NewLabel("Esperando comando...")
	labelRes.Wrapping = fyne.TextWrapWord

	labelReporte := widget.NewLabel("Esperando reporte del servidor...")
	labelReporte.Wrapping = fyne.TextWrapWord

	barraDisco := widget.NewProgressBar()
	barraDisco.Max = 100

	btnEjecutar := widget.NewButton("Ejecutar", func() {
		comando := strings.TrimSpace(textComando.Text)
		if comando == "" {
			return
		}
		bufferC.WriteString(comando + "\n")
		bufferC.Flush()
		textComando.SetText("")

		if comando == "bye" {
			socketCliente.Close()
			ventPpal.Close()
		}
	})

	// --- BOTÓN DE LOGIN ---
	btnLogin := widget.NewButton("Iniciar Sesión", func() {
		user := strings.TrimSpace(entryUser.Text)
		pass := strings.TrimSpace(entryPass.Text)

		bufferC.WriteString(user + ":" + pass + "\n")
		bufferC.Flush()

		resp, err := bufferS.ReadString('\n')
		if err != nil {
			labelEstado.SetText("Conexión cerrada por el servidor (intentos agotados o IP no permitida).")
			return
		}

		if strings.TrimSpace(resp) == "Recibido" {
			// Enviar periodo de reporte
			bufferC.WriteString(periodo + "\n")
			bufferC.Flush()

			// Cambiar contenido de la ventana al panel de comandos y reportes
			ventPpal.SetContent(container.NewVBox(
				labelReporte,
				barraDisco,
				textComando,
				btnEjecutar,
				labelRes,
			))
			ventPpal.Resize(fyne.NewSize(500, 320))

			// Iniciar hilo receptor
			go recMensajeGUI(bufferS, labelReporte, labelRes, barraDisco)
		} else {
			labelEstado.SetText("Error! Usuario o contraseña incorrectos :(")
			entryPass.SetText("")
		}
	})

	// Mostrar primero pantalla de login
	ventPpal.SetContent(container.NewVBox(
		widget.NewLabel("Login Servidor Debian"),
		entryUser,
		entryPass,
		btnLogin,
		labelEstado,
	))
	ventPpal.Resize(fyne.NewSize(380, 240))
	ventPpal.ShowAndRun()
}

func recMensajeGUI(bufferS *bufio.Reader, labelReporte *widget.Label, labelRes *widget.Label, barraDisco *widget.ProgressBar) {
	for {
		msjRec, err := bufferS.ReadString('\n')
		if err != nil {
			fyne.Do(func() {
				labelReporte.SetText("Desconectado del servidor.")
			})
			return
		}

		msjRec = strings.TrimSpace(msjRec)

		if strings.HasPrefix(msjRec, "REP|") {
			textoRep := strings.TrimPrefix(msjRec, "REP|")
			porcentaje := sacarPorcentajeDisco(textoRep)

			fyne.Do(func() {
				labelReporte.SetText(textoRep)
				if porcentaje >= 0 {
					barraDisco.SetValue(porcentaje)
				}
			})
		} else if strings.HasPrefix(msjRec, "CMD|") {
			salida := strings.ReplaceAll(strings.TrimPrefix(msjRec, "CMD|"), " | ", "\n")
			fyne.Do(func() {
				labelRes.SetText("Resultado:\n" + salida)
			})
		}
	}
}

func sacarPorcentajeDisco(texto string) float64 {
	// Busca el número entre "(" y "%" dentro de DISCO: 4.2G/20G (22%)
	idx := strings.Index(texto, "DISCO:")
	if idx == -1 {
		return -1
	}
	sub := texto[idx:]
	ini := strings.Index(sub, "(")
	fin := strings.Index(sub, "%")
	if ini == -1 || fin == -1 || fin <= ini {
		return -1
	}
	valor, err := strconv.ParseFloat(sub[ini+1:fin], 64)
	if err != nil {
		return -1
	}
	return valor
}
