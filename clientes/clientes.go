// Package clientes gestiona los clientes de la empresa de control de plagas:
// residenciales, comerciales e industriales.
//
// En esta etapa (Aprendizaje Autónomo 1) los datos se guardan en maps con
// tipos básicos. En el Aprendizaje Autónomo 2 se reemplazarán por un struct
// Cliente con sus métodos.
package clientes

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrClienteNoEncontrado se devuelve cuando el ID no corresponde a ningún cliente.
var ErrClienteNoEncontrado = errors.New("cliente no encontrado")

// Tipos de cliente que acepta el sistema.
var tiposValidos = map[string]bool{
	"residencial": true,
	"comercial":   true,
	"industrial":  true,
}

// Almacenamiento en memoria: un map por dato, todos indexados por el ID.
var (
	nombres     = map[int]string{}
	rucs        = map[int]string{}
	telefonos   = map[int]string{}
	direcciones = map[int]string{}
	tipos       = map[int]string{}
	siguienteID = 1
)

// RegistrarCliente registra un cliente nuevo y devuelve el ID asignado.
// El RUC debe tener 13 dígitos y el tipo debe ser residencial, comercial o industrial.
func RegistrarCliente(nombre, ruc, telefono, direccion, tipo string) (int, error) {
	nombre = strings.TrimSpace(nombre)
	tipo = strings.ToLower(strings.TrimSpace(tipo))

	if nombre == "" {
		return 0, errors.New("el nombre del cliente es obligatorio")
	}
	if !esNumerico(ruc, 13) {
		return 0, errors.New("el RUC debe tener 13 dígitos")
	}
	if !tiposValidos[tipo] {
		return 0, fmt.Errorf("tipo de cliente %q no válido (use residencial, comercial o industrial)", tipo)
	}

	id := siguienteID
	siguienteID++
	nombres[id] = nombre
	rucs[id] = ruc
	telefonos[id] = telefono
	direcciones[id] = direccion
	tipos[id] = tipo
	return id, nil
}

// ConsultarCliente devuelve los datos del cliente en una línea de texto.
func ConsultarCliente(id int) (string, error) {
	if !ExisteCliente(id) {
		return "", ErrClienteNoEncontrado
	}
	return fmt.Sprintf("#%d %s | RUC %s | %s | tel. %s | %s",
		id, nombres[id], rucs[id], tipos[id], telefonos[id], direcciones[id]), nil
}

// ExisteCliente indica si hay un cliente registrado con ese ID.
// main la usa antes de programar un servicio.
func ExisteCliente(id int) bool {
	_, ok := nombres[id]
	return ok
}

// ActualizarTelefono cambia el teléfono de contacto del cliente.
func ActualizarTelefono(id int, telefono string) error {
	if !ExisteCliente(id) {
		return ErrClienteNoEncontrado
	}
	if strings.TrimSpace(telefono) == "" {
		return errors.New("el teléfono no puede quedar vacío")
	}
	telefonos[id] = telefono
	return nil
}

// EliminarCliente da de baja al cliente y borra sus datos.
func EliminarCliente(id int) error {
	if !ExisteCliente(id) {
		return ErrClienteNoEncontrado
	}
	delete(nombres, id)
	delete(rucs, id)
	delete(telefonos, id)
	delete(direcciones, id)
	delete(tipos, id)
	return nil
}

// ListarClientes devuelve los clientes del tipo indicado, ordenados por ID.
// Con tipo vacío devuelve todos.
func ListarClientes(tipo string) []string {
	tipo = strings.ToLower(strings.TrimSpace(tipo))
	var lista []string
	for _, id := range idsOrdenados() {
		if tipo == "" || tipos[id] == tipo {
			linea, _ := ConsultarCliente(id)
			lista = append(lista, linea)
		}
	}
	return lista
}

// idsOrdenados devuelve los IDs en orden, porque recorrer un map en Go
// no garantiza ningún orden.
func idsOrdenados() []int {
	ids := make([]int, 0, len(nombres))
	for id := range nombres {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	return ids
}

// esNumerico comprueba que el texto tenga exactamente n dígitos.
func esNumerico(texto string, n int) bool {
	if len(texto) != n {
		return false
	}
	for _, c := range texto {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
