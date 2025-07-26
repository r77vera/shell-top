// Módulo de testing: Contiene pruebas unitarias para validar el funcionamiento
// correcto de los componentes principales de la shell
package main

import (
	"os"           // Para operaciones del sistema operativo en tests
	"path/filepath" // Para manipulación de rutas de archivos
	"testing"      // Framework de testing estándar de Go
)

// TestAnalizarEntrada prueba la función de parsing de la entrada del usuario.
// Verifica que la función AnalizarEntrada procese correctamente diferentes tipos
// de comandos y detecte apropiadamente la ejecución en segundo plano.
//
// Casos de prueba incluidos:
//   - Comando con argumentos múltiples
//   - Comando con argumentos que contienen espacios/comillas
//   - Comando con sufijo & (segundo plano)
//   - Entrada vacía
//   - Entrada con solo espacios
func TestAnalizarEntrada(t *testing.T) {
	// Definir tabla de casos de prueba
	// Cada caso contiene: entrada, comando esperado, argumentos esperados, flag segundo plano esperado
	tests := []struct {
		linea           string   // Línea de entrada a procesar
		comandoExp      string   // Comando esperado después del parsing
		argsExp         []string // Argumentos esperados después del parsing
		segundoPlanoExp bool     // Flag de segundo plano esperado
	}{
		// Caso 1: Comando típico con argumentos múltiples
		{"ls -l /tmp", "ls", []string{"-l", "/tmp"}, false},
		
		// Caso 2: Comando con argumentos que incluyen comillas
		// Nota: strings.Fields no maneja comillas como shell real,
		// divide por espacios independientemente de las comillas
		{"echo 'hola mundo'", "echo", []string{"'hola", "mundo'"}, false},
		
		// Caso 3: Comando con sufijo & para ejecución en segundo plano
		{"sleep 5 &", "sleep", []string{"5"}, true},
		
		// Caso 4: Entrada completamente vacía
		{"", "", nil, false},
		
		// Caso 5: Entrada con solo espacios en blanco
		{"   ", "", nil, false},
	}

	// Iterar sobre cada caso de prueba
	for _, tt := range tests {
		// Ejecutar la función bajo prueba
		comando, args, segundoPlano := AnalizarEntrada(tt.linea)

		// VERIFICACIÓN 1: Comando obtenido vs esperado
		if comando != tt.comandoExp {
			t.Errorf("Comando esperado: %q, obtenido: %q", tt.comandoExp, comando)
		}

		// VERIFICACIÓN 2: Argumentos obtenidos vs esperados
		// Usar función auxiliar equal para comparar slices
		if !equal(args, tt.argsExp) {
			t.Errorf("Argumentos esperados: %v, obtenidos: %v", tt.argsExp, args)
		}

		// VERIFICACIÓN 3: Flag de segundo plano obtenido vs esperado
		if segundoPlano != tt.segundoPlanoExp {
			t.Errorf("Segundo plano esperado: %v, obtenido: %v", tt.segundoPlanoExp, segundoPlano)
		}
	}
}

// TestEjecutarCd prueba la funcionalidad del comando interno 'cd'.
// Verifica que el comando cd cambie efectivamente el directorio de trabajo
// y que el directorio actual de la shell se actualice correctamente.
//
// El test:
//   1. Guarda el directorio actual para restaurarlo después
//   2. Crea un directorio temporal para las pruebas
//   3. Ejecuta el comando cd hacia el directorio temporal
//   4. Verifica que el directorio actual haya cambiado correctamente
//   5. Restaura el directorio original usando defer
func TestEjecutarCd(t *testing.T) {
	// PASO 1: Guardar el directorio actual para restaurarlo al final
	dirActual, _ := os.Getwd()
	// defer asegura que se ejecute al final de la función, sin importar cómo termine
	defer os.Chdir(dirActual) // Restaurar el directorio original

	// PASO 2: Crear un directorio temporal para la prueba
	// t.TempDir() crea un directorio temporal que se limpia automáticamente
	// al final del test
	tempDir := t.TempDir()

	// PASO 3: Ejecutar el comando cd con el directorio temporal
	if err := ejecutarCd([]string{tempDir}); err != nil {
		// Si hay error, fallar inmediatamente el test
		t.Fatalf("Error al cambiar al directorio temporal: %v", err)
	}

	// PASO 4: Verificar que el directorio efectivamente cambió
	dirDespues, _ := os.Getwd()
	
	// Normalizar las rutas para resolver enlaces simbólicos
	// Esto es necesario porque algunos sistemas operativos usan enlaces simbólicos
	// en rutas temporales (ej: /tmp -> /private/tmp en macOS)
	tempDirEval, _ := filepath.EvalSymlinks(tempDir)
	dirDespuesEval, _ := filepath.EvalSymlinks(dirDespues)

	// Comparar las rutas normalizadas
	if dirDespuesEval != tempDirEval {
		t.Errorf("Directorio esperado: %q, obtenido: %q", tempDirEval, dirDespuesEval)
	}
}

// equal es una función auxiliar que compara dos slices de strings para igualdad.
// Go no tiene una función built-in para comparar slices, por lo que implementamos
// nuestra propia función de comparación.
//
// Parámetros:
//   - a, b: slices de strings a comparar
//
// Retorna:
//   - bool: true si los slices son idénticos (mismo tamaño y mismos elementos
//           en el mismo orden), false en caso contrario
//
// Algoritmo:
//   1. Verificar que tengan la misma longitud
//   2. Comparar elemento por elemento en el mismo índice
func equal(a, b []string) bool {
	// Si tienen longitudes diferentes, no pueden ser iguales
	if len(a) != len(b) {
		return false
	}
	
	// Comparar cada elemento en la misma posición
	for i, v := range a {
		if v != b[i] {
			return false // Encontramos una diferencia
		}
	}
	
	// Si llegamos aquí, todos los elementos son iguales
	return true
}
