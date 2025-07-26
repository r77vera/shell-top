# Shell Básica en Go

Este proyecto es una implementación de una shell de comandos simple utilizando Go, como parte de un reto de programación.

## Enfoque del Análisis de la Línea de Comandos

La función `AnalizarEntrada` en `analizador.go` se encarga de procesar la entrada del usuario. Primero, limpia los espacios en blanco y verifica si la entrada está vacía. Luego, detecta si el comando debe ejecutarse en segundo plano buscando el sufijo `&`. Finalmente, divide la entrada en el comando principal y sus argumentos.

## Ejecución de Comandos Externos y Redirección de E/S

Para los comandos externos, se utiliza el paquete `os/exec`. La función `ejecutarComandoExterno` en `ejecutor.go` crea un nuevo `exec.Cmd` con el comando y sus argumentos. La entrada estándar (`stdin`), salida estándar (`stdout`) y error estándar (`stderr`) del proceso hijo se redirigen a los de la shell principal, permitiendo una interacción directa.

## Implementación de Comandos Internos

Se implementaron dos comandos internos:

- `cd <directorio>`: Cambia el directorio de trabajo actual de la shell. Utiliza la función `os.Chdir` para esta operación. Si no se proporciona un directorio, cambia al directorio `home` del usuario.
- `exit`: Termina la ejecución de la shell con `os.Exit(0)`.

## Estrategia para la Ejecución en Segundo Plano y Concurrencia

Si un comando termina con `&`, se ejecuta en segundo plano. En `ejecutarComandoExterno`, en lugar de `cmd.Run()`, se utiliza `cmd.Start()` para iniciar el comando sin bloquear. El PID del proceso hijo se muestra al usuario. Se lanza una goroutine separada que llama a `cmd.Wait()`, permitiendo que la shell principal continúe aceptando comandos mientras el proceso en segundo plano finaliza de forma independiente.
