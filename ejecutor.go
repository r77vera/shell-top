// Módulo ejecutor: Maneja la ejecución de comandos tanto internos como externos
// Se encarga de la lógica de ejecución, redirección de E/S y manejo de procesos
package main

import (
	"fmt"     // Para formatear salida y mostrar mensajes
	"os"      // Para operaciones del sistema operativo
	"os/exec" // Para ejecutar programas externos
)

// EjecutarComando es la función principal que actúa como dispatcher de comandos.
// Determina si un comando es interno (built-in) o externo y delega su ejecución.
//
// Comandos internos implementados:
//   - cd: cambio de directorio
//   - exit: salir de la shell
// 
// Todos los demás comandos se consideran externos y se buscan en el PATH del sistema.
//
// Parámetros:
//   - comando: string con el nombre del comando a ejecutar
//   - args: slice de strings con los argumentos del comando
//   - segundoPlano: boolean que indica si debe ejecutarse en background
//
// Retorna:
//   - error: nil si la ejecución fue exitosa, error específico en caso contrario
func EjecutarComando(comando string, args []string, segundoPlano bool) error {
	// Usar switch para determinar el tipo de comando y delegarlo
	switch comando {
	case "cd":
		// Comando interno: cambio de directorio
		return ejecutarCd(args)
	case "exit":
		// Comando interno: salir de la shell
		return ejecutarExit()
	default:
		// Comando externo: delegar a ejecutarComandoExterno
		return ejecutarComandoExterno(comando, args, segundoPlano)
	}
}

// ejecutarCd implementa el comando interno 'cd' para cambiar el directorio de trabajo.
// 
// Comportamiento:
//   - Sin argumentos: cambia al directorio home del usuario
//   - Con argumento: cambia al directorio especificado
//
// Parámetros:
//   - args: slice de argumentos del comando cd
//
// Retorna:
//   - error: nil si el cambio fue exitoso, error si el directorio no existe o no es accesible
func ejecutarCd(args []string) error {
	// Si no se proporcionan argumentos, ir al directorio home
	if len(args) == 0 {
		// Obtener el directorio home del usuario actual
		home, err := os.UserHomeDir()
		if err != nil {
			// Error obteniendo el directorio home
			return err
		}
		// Cambiar al directorio home
		return os.Chdir(home)
	}
	
	// Si hay argumentos, usar el primer argumento como destino
	// os.Chdir cambia el directorio de trabajo del proceso actual
	return os.Chdir(args[0])
}

// ejecutarExit implementa el comando interno 'exit' para terminar la shell.
// Termina inmediatamente el programa con código de salida 0 (éxito).
//
// Retorna:
//   - error: en teoría nunca retorna porque os.Exit termina el programa,
//            pero se mantiene la signatura para consistencia
func ejecutarExit() error {
	// os.Exit(0) termina inmediatamente el programa con código de salida 0
	// No ejecuta defer statements ni finalizers
	os.Exit(0)
	return nil // Esta línea nunca se ejecuta
}

// ejecutarComandoExterno maneja la ejecución de programas externos del sistema.
// Utiliza os/exec para crear y ejecutar procesos hijo.
//
// Funcionalidad:
//   - Busca el programa en el PATH del sistema
//   - Configura redirección de stdin, stdout y stderr
//   - Maneja ejecución síncrona (foreground) y asíncrona (background)
//   - Para procesos en background, usa goroutines para no bloquear la shell
//
// Parámetros:
//   - comando: nombre del programa a ejecutar (ej: "ls", "cat", "grep")
//   - args: argumentos para pasar al programa
//   - segundoPlano: true para ejecución asíncrona, false para síncrona
//
// Retorna:
//   - error: nil si la ejecución inició correctamente, error en caso contrario
func ejecutarComandoExterno(comando string, args []string, segundoPlano bool) error {
	// PASO 1: Crear el comando usando exec.Command
	// exec.Command busca el programa en el PATH del sistema
	// y prepara la estructura para la ejecución
	cmd := exec.Command(comando, args...)
	
	// PASO 2: Configurar redirección de E/S
	// Conectar los streams del proceso hijo con los de la shell padre
	// Esto permite que la salida del comando aparezca en la terminal
	cmd.Stdin = os.Stdin   // Entrada estándar: teclado → proceso hijo
	cmd.Stdout = os.Stdout // Salida estándar: proceso hijo → pantalla
	cmd.Stderr = os.Stderr // Error estándar: proceso hijo → pantalla

	// PASO 3: Determinar modo de ejecución (foreground vs background)
	if segundoPlano {
		// EJECUCIÓN EN SEGUNDO PLANO (ASÍNCRONA)
		
		// cmd.Start() inicia el proceso pero NO espera a que termine
		// Retorna inmediatamente, permitiendo que la shell continúe
		err := cmd.Start()
		if err != nil {
			// Error iniciando el proceso (ej: comando no encontrado)
			return err
		}
		
		// Mostrar información del proceso en background al usuario
		// cmd.Process.Pid contiene el Process ID del proceso hijo
		fmt.Printf("[PID: %d] Proceso en segundo plano iniciado\n", cmd.Process.Pid)
		
		// Lanzar una goroutine para esperar la terminación del proceso
		// Esto evita procesos zombie y libera recursos cuando el proceso termina
		go func() {
			// cmd.Wait() espera a que el proceso termine y libera recursos
			// Se ejecuta en una goroutine separada para no bloquear la shell principal
			cmd.Wait()
			// Nota: En una implementación más avanzada, aquí podríamos notificar
			// al usuario cuando el proceso en background termine
		}()
		
		// Retornar inmediatamente para que la shell pueda procesar más comandos
		return nil
	}

	// EJECUCIÓN EN PRIMER PLANO (SÍNCRONA)
	// cmd.Run() es equivalente a cmd.Start() seguido de cmd.Wait()
	// Inicia el proceso y espera hasta que termine antes de retornar
	// La shell se bloquea hasta que el comando complete su ejecución
	return cmd.Run()
}
