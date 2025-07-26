// Módulo analizador: Se encarga del parsing y análisis de la entrada del usuario
// Implementa la lógica para convertir una línea de texto en componentes ejecutables
package main

import (
	"strings" // Para manipulación de cadenas de texto
)

// AnalizarEntrada es la función principal de parsing que procesa la línea de entrada del usuario.
// 
// Funcionalidad:
// 1. Limpia espacios en blanco al inicio y final
// 2. Detecta si el comando debe ejecutarse en segundo plano (sufijo &)
// 3. Separa el comando principal de sus argumentos
// 
// Parámetros:
//   - entrada: string que contiene la línea completa ingresada por el usuario
// 
// Retorna:
//   - string: el comando principal a ejecutar (ej: "ls", "cat", "cd")
//   - []string: slice con los argumentos del comando (ej: ["-l", "/tmp"])
//   - bool: true si debe ejecutarse en segundo plano (encontró &), false en caso contrario
//
// Ejemplos:
//   - "ls -l /tmp" → ("ls", ["-l", "/tmp"], false)
//   - "sleep 5 &" → ("sleep", ["5"], true)
//   - "   " → ("", nil, false)
func AnalizarEntrada(entrada string) (string, []string, bool) {
	// PASO 1: Limpiar la entrada removiendo espacios al inicio y final
	// strings.TrimSpace elimina espacios, tabs, saltos de línea, etc.
	entrada = strings.TrimSpace(entrada)
	
	// Si la entrada está vacía después de limpiar, no hay nada que procesar
	if entrada == "" {
		return "", nil, false
	}

	// PASO 2: Detectar ejecución en segundo plano
	// Inicializar la bandera de segundo plano como false
	segundoPlano := false
	
	// Verificar si la línea termina con el carácter "&"
	// strings.HasSuffix verifica si la cadena termina con el sufijo dado
	if strings.HasSuffix(entrada, "&") {
		// Marcar que debe ejecutarse en segundo plano
		segundoPlano = true
		
		// Remover el "&" de la entrada y limpiar espacios nuevamente
		// strings.TrimSuffix remueve el sufijo especificado del final
		entrada = strings.TrimSpace(strings.TrimSuffix(entrada, "&"))
	}

	// PASO 3: Separar comando y argumentos
	// strings.Fields divide la cadena en palabras separadas por espacios
	// Automáticamente maneja múltiples espacios consecutivos
	partes := strings.Fields(entrada)
	
	// El primer elemento es siempre el comando principal
	comando := partes[0]
	
	// Los elementos restantes son los argumentos (puede ser slice vacío)
	// partes[1:] crea un slice con todos los elementos desde el índice 1 en adelante
	args := partes[1:]

	// Retornar los tres componentes analizados
	return comando, args, segundoPlano
}
