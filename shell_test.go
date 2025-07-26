package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalizarEntrada(t *testing.T) {
	tests := []struct {
		linea           string
		comandoExp      string
		argsExp         []string
		segundoPlanoExp bool
	}{
		{"ls -l /tmp", "ls", []string{"-l", "/tmp"}, false},
		{"echo 'hola mundo'", "echo", []string{"'hola", "mundo'"}, false},
		{"sleep 5 &", "sleep", []string{"5"}, true},
		{"", "", nil, false},
		{"   ", "", nil, false},
	}

	for _, tt := range tests {
		comando, args, segundoPlano := AnalizarEntrada(tt.linea)

		if comando != tt.comandoExp {
			t.Errorf("Comando esperado: %q, obtenido: %q", tt.comandoExp, comando)
		}

		if !equal(args, tt.argsExp) {
			t.Errorf("Argumentos esperados: %v, obtenidos: %v", tt.argsExp, args)
		}

		if segundoPlano != tt.segundoPlanoExp {
			t.Errorf("Segundo plano esperado: %v, obtenido: %v", tt.segundoPlanoExp, segundoPlano)
		}
	}
}

func TestEjecutarCd(t *testing.T) {
	dirActual, _ := os.Getwd()
	defer os.Chdir(dirActual) // Restaurar el directorio original

	tempDir := t.TempDir()

	if err := ejecutarCd([]string{tempDir}); err != nil {
		t.Fatalf("Error al cambiar al directorio temporal: %v", err)
	}

	dirDespues, _ := os.Getwd()
	// Normalizar las rutas para resolver enlaces simbólicos
	tempDirEval, _ := filepath.EvalSymlinks(tempDir)
	dirDespuesEval, _ := filepath.EvalSymlinks(dirDespues)

	if dirDespuesEval != tempDirEval {
		t.Errorf("Directorio esperado: %q, obtenido: %q", tempDirEval, dirDespuesEval)
	}
}

// equal compara dos slices de strings
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
