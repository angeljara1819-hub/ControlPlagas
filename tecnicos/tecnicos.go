// Package tecnicos gestiona al personal técnico que ejecuta los servicios.
package tecnicos

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrTecnicoNoEncontrado se devuelve cuando el ID no corresponde a ningún técnico.
var ErrTecnicoNoEncontrado = errors.New("técnico no encontrado")

// Almacenamiento en memoria, indexado por el ID del técnico.
var (
	nombres        = map[int]string{}
	cedulas        = map[int]string{}
	especialidades = map[int]string{}
	activos        = map[int]bool{}
	siguienteID    = 1
)

// RegistrarTecnico registra un técnico nuevo (activo) y devuelve su ID.
// La cédula debe tener 10 dígitos.
func RegistrarTecnico(nombre, cedula, especialidad string) (int, error) {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return 0, errors.New("el nombre del técnico es obligatorio")
	}
	if len(cedula) != 10 || strings.Trim(cedula, "0123456789") != "" {
		return 0, errors.New("la cédula debe tener 10 dígitos")
	}

	id := siguienteID
	siguienteID++
	nombres[id] = nombre
	cedulas[id] = cedula
	especialidades[id] = strings.ToLower(strings.TrimSpace(especialidad))
	activos[id] = true
	return id, nil
}

// ConsultarTecnico devuelve los datos del técnico en una línea de texto.
func ConsultarTecnico(id int) (string, error) {
	if _, ok := nombres[id]; !ok {
		return "", ErrTecnicoNoEncontrado
	}
	estado := "activo"
	if !activos[id] {
		estado = "inactivo"
	}
	return fmt.Sprintf("#%d %s | CI %s | %s | %s",
		id, nombres[id], cedulas[id], especialidades[id], estado), nil
}

// EstaActivo indica si el técnico existe y puede recibir servicios.
// main la usa antes de asignar un técnico a un servicio.
func EstaActivo(id int) bool {
	return activos[id]
}

// CambiarEstado activa o desactiva a un técnico (por ejemplo, por vacaciones).
func CambiarEstado(id int, activo bool) error {
	if _, ok := nombres[id]; !ok {
		return ErrTecnicoNoEncontrado
	}
	activos[id] = activo
	return nil
}

// ListarActivos devuelve los técnicos activos, ordenados por ID.
func ListarActivos() []string {
	ids := make([]int, 0, len(nombres))
	for id := range nombres {
		if activos[id] {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)

	lista := make([]string, 0, len(ids))
	for _, id := range ids {
		linea, _ := ConsultarTecnico(id)
		lista = append(lista, linea)
	}
	return lista
}
