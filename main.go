package main

import (
	"bufio"  // Para leer línea por línea desde la entrada estándar
	"fmt"    // Para formatear y mostrar salida
	"os"     // Para interactuar con el sistema operativo
	"os/user" // Para obtener información del usuario actual
)

// main es la función principal que implementa el bucle REPL (Bucle de lectura-evaluación-impresión)
// de la shell. Ejecuta indefinidamente hasta que el usuario ejecute el comando "exit"
func main() {
	// Mostrar mensaje de bienvenida al iniciar la shell
	mostrarBienvenida()
	
	// Crear un lector para capturar la entrada del usuario desde stdin
	// bufio.NewReader es más eficiente que fmt.Scan para leer líneas completas
	lector := bufio.NewReader(os.Stdin)

	// Bucle infinito que implementa el REPL de la shell
	for {
		// PASO 1: Obtener información para mostrar en el prompt
		
		// Obtener el directorio de trabajo actual para mostrarlo en el prompt
		wd, err := os.Getwd()
		if err != nil {
			// Si hay error obteniendo el directorio, mostrar el error y usar cadena vacía
			fmt.Fprintln(os.Stderr, "Error al obtener el directorio actual:", err)
			wd = ""
		}

		// Obtener información del usuario actual para personalizar el prompt
		currentUser, err := user.Current()
		if err != nil {
			// Si hay error obteniendo el usuario, mostrar el error pero continuar
			fmt.Fprintln(os.Stderr, "Error al obtener el usuario actual:", err)
		}

		// PASO 2: Mostrar el prompt y leer la entrada del usuario
		
		// Mostrar prompt colorizado en formato "usuario:directorio goshell> "
		// Usando códigos ANSI para colores
		mostrarPrompt(currentUser.Username, wd)

		// Leer una línea completa de entrada hasta encontrar '\n' (Enter)
		// ReadString incluye el carácter delimitador en el resultado
		entrada, err := lector.ReadString('\n')
		if err != nil {
			// Si hay error leyendo (ej: EOF), mostrar error y continuar el bucle
			fmt.Fprintln(os.Stderr, "Error al leer la entrada:", err)
			continue
		}

		// PASO 3: Analizar y procesar la entrada del usuario
		
		// Parsear la línea de entrada para extraer:
		// - comando: el programa o comando a ejecutar
		// - args: lista de argumentos para el comando
		// - segundoPlano: boolean que indica si debe ejecutarse en background (&)
		comando, args, segundoPlano := AnalizarEntrada(entrada)

		// Si no hay comando (línea vacía o solo espacios), continuar al siguiente ciclo
		// Esto evita errores al intentar ejecutar comandos vacíos
		// También se ignoran líneas que solo contienen espacios o tabulaciones
		if comando == "" {
			continue
		}

		// PASO 4: Ejecutar el comando
		
		// Intentar ejecutar el comando (interno o externo) con sus argumentos
		// Si hay error en la ejecución, mostrarlo pero no terminar la shell
		if err := EjecutarComando(comando, args, segundoPlano); err != nil {
			fmt.Fprintln(os.Stderr, "Error al ejecutar el comando:", err)
		}
		
		// El bucle continúa para procesar el siguiente comando
		// Solo se rompe cuando se ejecuta el comando interno "exit"
	}
}

// mostrarBienvenida muestra un mensaje de bienvenida colorizado al iniciar la shell
// Incluye información sobre los comandos disponibles y ejemplos de uso
func mostrarBienvenida() {
	// Definir códigos de color ANSI
	const (
		ColorReset  = "\033[0m"
		ColorRojo   = "\033[31m"
		ColorVerde  = "\033[32m"
		ColorAmarillo = "\033[33m"
		ColorAzul   = "\033[34m"
		ColorMagenta = "\033[35m"
		ColorCian   = "\033[36m"
		ColorBlanco = "\033[37m"
		ColorNegrita = "\033[1m"
	)

	fmt.Printf("%s%s", ColorCian, ColorNegrita)
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                          🐚 GoShell v1.0                         ║")
	fmt.Println("║                    Shell Básica implementada en Go               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Printf("%s", ColorReset)

	fmt.Printf("\n%s%s¡Bienvenido a GoShell!%s\n", ColorVerde, ColorNegrita, ColorReset)
	fmt.Printf("Una shell simple pero poderosa construida en Go.\n\n")

	// Mostrar información de autores con enlaces de GitHub
	fmt.Printf("%s%s👥 DESARROLLADORES:%s\n", ColorMagenta, ColorNegrita, ColorReset)
	fmt.Printf("\n")
	
	// ASCII Art para r77vera
	fmt.Printf("%s%s", ColorRojo, ColorNegrita)
	fmt.Println("  ██████╗ ███████╗███████╗██╗   ██╗███████╗██████╗  █████╗ ")
	fmt.Println("  ██╔══██╗╚════██║╚════██║██║   ██║██╔════╝██╔══██╗██╔══██╗")
	fmt.Println("  ██████╔╝    ██╔╝    ██╔╝██║   ██║█████╗  ██████╔╝███████║")
	fmt.Println("  ██╔══██╗   ██╔╝    ██╔╝ ╚██╗ ██╔╝██╔══╝  ██╔══██╗██╔══██║")
	fmt.Println("  ██║  ██║   ██║     ██║   ╚████╔╝ ███████╗██║  ██║██║  ██║")
	fmt.Println("  ╚═╝  ╚═╝   ╚═╝     ╚═╝    ╚═══╝  ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝")
	fmt.Printf("%s", ColorReset)
	fmt.Printf("                    %s🔗 https://github.com/r77vera%s\n\n", ColorCian, ColorReset)

	// ASCII Art para Bjohan23
	fmt.Printf("%s%s", ColorAzul, ColorNegrita)
	fmt.Println("  ██████╗      ██╗ ██████╗ ██╗  ██╗ █████╗ ███╗   ██╗██████╗ ██████╗ ")
	fmt.Println("  ██╔══██╗     ██║██╔═══██╗██║  ██║██╔══██╗████╗  ██║╚════██╗╚════██╗")
	fmt.Println("  ██████╔╝     ██║██║   ██║███████║███████║██╔██╗ ██║ █████╔╝ █████╔╝")
	fmt.Println("  ██╔══██╗██   ██║██║   ██║██╔══██║██╔══██║██║╚██╗██║██╔═══╝  ╚═══██╗")
	fmt.Println("  ██████╔╝╚█████╔╝╚██████╔╝██║  ██║██║  ██║██║ ╚████║███████╗██████╔╝")
	fmt.Println("  ╚═════╝  ╚════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚══════╝╚═════╝ ")
	fmt.Printf("%s", ColorReset)
	fmt.Printf("                    %s🔗 https://github.com/Bjohan23%s\n\n", ColorCian, ColorReset)
	
	fmt.Printf("%s%s🎓 Proyecto académico - Taller de Lenguajes de Programación%s\n", ColorAmarillo, ColorNegrita, ColorReset)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	fmt.Printf("%s%sComandos Internos Disponibles:%s\n", ColorAmarillo, ColorNegrita, ColorReset)
	fmt.Printf("  %s• cd [directorio]%s  - Cambiar directorio (sin args = ir a home)\n", ColorCian, ColorReset)
	fmt.Printf("  %s• exit%s             - Salir de la shell\n", ColorCian, ColorReset)

	fmt.Printf("\n%s%sComandos Externos:%s\n", ColorAmarillo, ColorNegrita, ColorReset)
	fmt.Printf("  %s• ls, cat, echo, grep, etc.%s - Cualquier programa en tu PATH\n", ColorCian, ColorReset)

	fmt.Printf("\n%s%sEjecución en Segundo Plano:%s\n", ColorAmarillo, ColorNegrita, ColorReset)
	fmt.Printf("  %s• comando &%s         - Ejecutar comando en background\n", ColorCian, ColorReset)

	fmt.Printf("\n%s%sEjemplos:%s\n", ColorMagenta, ColorNegrita, ColorReset)
	fmt.Printf("  %sgoshell>%s ls -la\n", ColorVerde, ColorReset)
	fmt.Printf("  %sgoshell>%s cd /tmp\n", ColorVerde, ColorReset)
	fmt.Printf("  %sgoshell>%s sleep 5 &\n", ColorVerde, ColorReset)
	fmt.Printf("  %sgoshell>%s exit\n", ColorVerde, ColorReset)

	fmt.Printf("\n%s%s💡 Tip:%s Usa 'exit' para salir de la shell\n", ColorAmarillo, ColorNegrita, ColorReset)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// mostrarPrompt muestra el prompt colorizado de la shell
// Formato: usuario:directorio goshell>
func mostrarPrompt(usuario, directorio string) {
	// Definir códigos de color ANSI
	const (
		ColorReset    = "\033[0m"
		ColorVerde    = "\033[32m"
		ColorAzul     = "\033[34m"
		ColorMagenta  = "\033[35m"
		ColorCian     = "\033[36m"
		ColorNegrita  = "\033[1m"
	)

	// Acortar el directorio si es muy largo (mostrar solo los últimos 40 caracteres)
	dirMostrar := directorio
	if len(directorio) > 40 {
		dirMostrar = "..." + directorio[len(directorio)-37:]
	}

	// Mostrar prompt colorizado: usuario en verde, directorio en azul, "goshell>" en magenta
	fmt.Printf("%s%s%s%s:%s%s%s%s %s%sgoshell>%s ",
		ColorVerde, ColorNegrita, usuario, ColorReset,
		ColorAzul, ColorNegrita, dirMostrar, ColorReset,
		ColorMagenta, ColorNegrita, ColorReset)
}
