package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"time"

	"golang.org/x/term"
)

func main() {
	//https://emn178.github.io/online-tools/sha256.html

	intentos := 1

	for {

		const hspasswd string = "caf90169eefa5f807d577486b9f795ab86ae2983c5c20806cff959117e90af18" //password del sistema
		//pssword :linux

		fmt.Printf("[intento %d] Digite la password : ", intentos)
		//var upasswd string
		//fmt.Scanf("%s", &upasswd) //lee como string lo ingresado y lo define cmo la variable upasswd
		upasswd, _ := term.ReadPassword(int(os.Stdin.Fd()))
		// la _ variable anonima para evitar errores que salien antes, que era que las funciones que usabamos reciben dos datos pero el segundo no lo necesitamos y no lo ha entobnces con eso no se toma ecuenta el segundp dato

		hupasswd := sha256.Sum256([]byte(upasswd)) //el string lo comvierta a unos bytes

		//fmt.Printf("hupasswd= %s \n", hupasswd) //imrime el hupasswd como string
		shupasswd := fmt.Sprintf("%x", hupasswd)

		if shupasswd == hspasswd {
			fmt.Println("\nLogin correctooo ;D")
			break
		}
		fmt.Println("\nERROR !!Login innnncorrectooo :c")

		if intentos%3 == 0 {
			time.Sleep(1 * time.Second)
			cmd := exec.Command("clear") //lo que hac es prepara el clear
			cmd.Stdout = os.Stdout       // conecta este programa con la termina
			cmd.Run()                    //corre clear en la terminal
		}

		intentos++
		time.Sleep(time.Duration(intentos) * 500 * time.Millisecond)
	}
}
