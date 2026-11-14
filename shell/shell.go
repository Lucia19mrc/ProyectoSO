package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// VER PQ EL CODIGO NO ME ESTA MOSTRANDO LOS COMANDO O BNO PQ NO LOS EJECUT EN LA TERMINAL
func main() {
	for {
		//var comando string
		fmt.Print("Oper# : ")
		//fmt.Scanf("%s", &comando)

		lector := bufio.NewReader(os.Stdin)   //tomo la entrada estndar sel didtema y genero un buffer(menmoria temporal), re tiene lo que llega y leo esa memoria
		comando, _ := lector.ReadString('\n') //el \n es el enter que dice cuando termino de ingresar un valor
		comando = strings.TrimRight(comando, "\n")
		//trimright borra caracteres a la derechha al final

		if comando == "bye" {
			fmt.Println("Se ha terminado el programa ;)")
			break
		}
		fmt.Println("Ejecutando", comando, "...")
		//------codigo para ejecutar comandos
		var shell *exec.Cmd
		shell = exec.Command("bash", "-c", comando)
		res, _ := shell.Output()
		//------comando ejecutado ---------
		fmt.Println("Oper#", string(res))

	}
}
