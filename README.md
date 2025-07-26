# GoShell - Shell Básica en Go

Una implementación simple pero funcional de una shell de comandos construida en Go. Este proyecto demuestra los conceptos fundamentales de sistemas operativos como la ejecución de procesos, redirección de E/S y manejo de concurrencia.

![Go Version](https://img.shields.io/badge/Go-1.16+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

## 🚀 Características

- **Bucle REPL interactivo** con prompt personalizado
- **Comandos externos** - ejecuta cualquier programa disponible en el PATH
- **Comandos internos** - `cd` y `exit` implementados nativamente
- **Ejecución en segundo plano** - soporte para comandos con `&`
- **Redirección completa de E/S** - stdin, stdout y stderr
- **Manejo robusto de errores** y validación de entrada
- **Concurrencia segura** usando goroutines para procesos en background

## 📋 Requisitos

- **Go 1.16 o superior**
- Sistema operativo: Linux, macOS, o Windows
- Terminal compatible con ANSI (para el prompt colorizado)

## 🛠️ Instalación

### Opción 1: Clonar el repositorio

```bash
# Clonar el repositorio
git clone https://github.com/r77vera/shell-top.git

# Navegar al directorio del proyecto
cd shell-top

# Construir el ejecutable
go build -o goshell

# Ejecutar la shell
./goshell
```

### Opción 2: Instalación directa con Go

```bash
# Instalar directamente desde GitHub
go install github.com/r77vera/shell-top@latest

# Ejecutar (asegurate de que $GOPATH/bin esté en tu PATH)
shell-top
```

### Opción 3: Ejecución sin construcción

```bash
# Clonar el repositorio
git clone https://github.com/r77vera/shell-top.git
cd shell-top

# Ejecutar directamente
go run *.go
```

## 🎯 Uso

### Iniciando la Shell

```bash
$ ./goshell
usuario:/ruta/actual goshell> 
```

### Comandos Soportados

#### Comandos Externos
Cualquier programa disponible en tu PATH:

```bash
goshell> ls -la
goshell> cat archivo.txt
goshell> grep "patron" archivo.txt
goshell> python script.py
```

#### Comandos Internos

**Cambio de directorio:**
```bash
goshell> cd /ruta/destino     # Cambiar a directorio específico
goshell> cd                   # Cambiar al directorio home
```

**Salir de la shell:**
```bash
goshell> exit
```

#### Ejecución en Segundo Plano

Agrega `&` al final del comando para ejecutarlo en background:

```bash
goshell> sleep 10 &
[PID: 12345] Proceso en segundo plano iniciado
goshell> # La shell continúa disponible inmediatamente
```

### Ejemplos de Uso

```bash
# Navegación básica
goshell> pwd
/home/usuario
goshell> ls
archivo1.txt  archivo2.txt  directorio/
goshell> cd directorio
goshell> pwd
/home/usuario/directorio

# Comandos con argumentos múltiples
goshell> ls -la /tmp
goshell> find . -name "*.go"

# Procesos en segundo plano
goshell> ping google.com &
[PID: 12346] Proceso en segundo plano iniciado
goshell> ps aux | grep ping  # La shell sigue funcionando
```

## 🧪 Ejecutar Tests

```bash
# Ejecutar todas las pruebas
go test

# Ejecutar con información detallada
go test -v

# Ejecutar con cobertura
go test -cover
```

### Pruebas Incluidas

- **TestAnalizarEntrada**: Valida el parsing de comandos y argumentos
- **TestEjecutarCd**: Verifica la funcionalidad del comando `cd`

## 🏗️ Arquitectura del Proyecto

```
shell-top/
├── main.go          # Bucle REPL principal
├── analizador.go    # Parsing de la entrada del usuario
├── ejecutor.go      # Ejecución de comandos internos y externos
├── shell_test.go    # Pruebas unitarias
├── README.md        # Este archivo
└── go.mod          # Dependencias del módulo Go
```

### Flujo de Ejecución

1. **Lectura** - El prompt solicita entrada del usuario
2. **Análisis** - Se parsea la línea para extraer comando, argumentos y flags
3. **Ejecución** - Se determina si es comando interno o externo
4. **Salida** - Se muestra el resultado y se vuelve al paso 1

## 🔧 Detalles Técnicos

### Enfoque del Análisis de la Línea de Comandos

La función `AnalizarEntrada` en `analizador.go` procesa la entrada del usuario:
1. Limpia espacios en blanco al inicio y final
2. Detecta ejecución en segundo plano (sufijo `&`)
3. Separa comando principal de sus argumentos usando `strings.Fields`

### Ejecución de Comandos Externos y Redirección de E/S

Para comandos externos se utiliza `os/exec`:
- `exec.Command` crea la estructura del proceso
- Se redirige stdin, stdout y stderr al proceso padre
- `cmd.Run()` para ejecución síncrona, `cmd.Start()` para asíncrona

### Implementación de Comandos Internos

**Comandos internos implementados:**
- `cd <directorio>`: Usa `os.Chdir` para cambiar directorio
- `exit`: Termina la shell con `os.Exit(0)`

### Estrategia para Ejecución en Segundo Plano

Los comandos con `&` se ejecutan asincrónicamente:
- `cmd.Start()` inicia el proceso sin bloquear
- Se muestra el PID al usuario
- Una goroutine separada ejecuta `cmd.Wait()` para limpiar recursos

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver `LICENSE` para más detalles.

## 👥 Autores

- **r77vera** - *Desarrollo inicial* - [r77vera](https://github.com/r77vera)

## 🙏 Agradecimientos

- Inspirado en shells clásicas como Bash y Zsh
- Proyecto educativo para entender sistemas operativos
- Comunidad de Go por excelente documentación
