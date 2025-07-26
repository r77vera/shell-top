package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
)

func main() {
	lector := bufio.NewReader(os.Stdin)

	for {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error al obtener el directorio actual:", err)
			wd = ""
		}

		currentUser, err := user.Current()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error al obtener el usuario actual:", err)
		}

		fmt.Printf("%s:%s goshell> ", currentUser.Username, wd)

		entrada, err := lector.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error al leer la entrada:", err)
			continue
		}

		comando, args, segundoPlano := AnalizarEntrada(entrada)

		if comando == "" {
			continue
		}

		if err := EjecutarComando(comando, args, segundoPlano); err != nil {
			fmt.Fprintln(os.Stderr, "Error al ejecutar el comando:", err)
		}
	}
}
