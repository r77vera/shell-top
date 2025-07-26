package main

import (
	"strings"
)

// AnalizarEntrada procesa la línea de entrada del usuario para separar el comando,
// los argumentos y detectar si debe ejecutarse en segundo plano.
func AnalizarEntrada(entrada string) (string, []string, bool) {
	entrada = strings.TrimSpace(entrada)
	if entrada == "" {
		return "", nil, false
	}

	segundoPlano := false
	if strings.HasSuffix(entrada, "&") {
		segundoPlano = true
		entrada = strings.TrimSpace(strings.TrimSuffix(entrada, "&"))
	}

	partes := strings.Fields(entrada)
	comando := partes[0]
	args := partes[1:]

	return comando, args, segundoPlano
}
