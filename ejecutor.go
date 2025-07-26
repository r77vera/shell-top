package main

import (
	"fmt"
	"os"
	"os/exec"
)

// EjecutarComando se encarga de diferenciar entre comandos internos y externos y ejecutarlos.
func EjecutarComando(comando string, args []string, segundoPlano bool) error {
	switch comando {
	case "cd":
		return ejecutarCd(args)
	case "exit":
		return ejecutarExit()
	default:
		return ejecutarComandoExterno(comando, args, segundoPlano)
	}
}

// ejecutarCd cambia el directorio de trabajo actual.
func ejecutarCd(args []string) error {
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return os.Chdir(home)
	}
	return os.Chdir(args[0])
}

// ejecutarExit termina la shell.
func ejecutarExit() error {
	os.Exit(0)
	return nil
}

// ejecutarComandoExterno ejecuta comandos que no son internos de la shell.
func ejecutarComandoExterno(comando string, args []string, segundoPlano bool) error {
	cmd := exec.Command(comando, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if segundoPlano {
		err := cmd.Start()
		if err != nil {
			return err
		}
		fmt.Printf("[PID: %d] Proceso en segundo plano iniciado\n", cmd.Process.Pid)
		go func() {
			cmd.Wait()
		}()
		return nil
	}

	return cmd.Run()
}
