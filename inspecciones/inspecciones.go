// Package inspecciones gestiona las visitas de inspección a los clientes y
// los informes técnicos que se les entregan después de cada visita.
//
// Ciclo de una inspección: se programa con fecha de visita, se registra su
// resultado (conforme o no conforme) y se registra la entrega del informe
// técnico al cliente, que debe hacerse dentro del plazo establecido.
package inspecciones

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Estados de una inspección.
const (
	Programada = "programada"
	Realizada  = "realizada"
)

const formatoFecha = "2006-01-02"

// plazoInformeDias son los días que tiene la empresa para entregar el informe
// técnico al cliente, contados desde la fecha de la visita.
const plazoInformeDias = 5

// ErrInspeccionNoEncontrada se devuelve cuando el ID no corresponde a ninguna inspección.
var ErrInspeccionNoEncontrada = errors.New("inspección no encontrada")

// Almacenamiento en memoria, indexado por el ID de la inspección.
var (
	clienteDe      = map[int]int{}
	fechaVisitaDe  = map[int]string{}
	estadoDe       = map[int]string{}
	hallazgosDe    = map[int]string{}
	conformeDe     = map[int]bool{}
	fechaInformeDe = map[int]string{} // vacío = informe todavía no entregado
	siguienteID    = 1
)

// ProgramarInspeccion agenda una visita de inspección a un cliente en la fecha
// indicada (AAAA-MM-DD) y devuelve el ID de la inspección.
func ProgramarInspeccion(idCliente int, fechaVisita string) (int, error) {
	if idCliente <= 0 {
		return 0, errors.New("el ID del cliente no es válido")
	}
	if _, err := time.Parse(formatoFecha, fechaVisita); err != nil {
		return 0, errors.New("la fecha debe tener el formato AAAA-MM-DD")
	}

	id := siguienteID
	siguienteID++
	clienteDe[id] = idCliente
	fechaVisitaDe[id] = fechaVisita
	estadoDe[id] = Programada
	return id, nil
}

// RegistrarResultado registra lo encontrado en la visita y marca la inspección
// como realizada. Si no es conforme, los hallazgos son obligatorios.
func RegistrarResultado(idInspeccion int, hallazgos string, conforme bool) error {
	estado, ok := estadoDe[idInspeccion]
	if !ok {
		return ErrInspeccionNoEncontrada
	}
	if estado != Programada {
		return errors.New("la inspección ya tiene un resultado registrado")
	}
	hallazgos = strings.TrimSpace(hallazgos)
	if !conforme && hallazgos == "" {
		return errors.New("una inspección no conforme debe describir los hallazgos")
	}
	hallazgosDe[idInspeccion] = hallazgos
	conformeDe[idInspeccion] = conforme
	estadoDe[idInspeccion] = Realizada
	return nil
}

// RegistrarEntregaInforme registra la fecha (AAAA-MM-DD) en que se entregó al
// cliente el informe técnico de una inspección realizada.
// La fecha no puede ser anterior a la visita.
func RegistrarEntregaInforme(idInspeccion int, fechaEntrega string) error {
	estado, ok := estadoDe[idInspeccion]
	if !ok {
		return ErrInspeccionNoEncontrada
	}
	if estado != Realizada {
		return errors.New("primero registre el resultado de la inspección")
	}
	if fechaInformeDe[idInspeccion] != "" {
		return fmt.Errorf("el informe ya se entregó el %s", fechaInformeDe[idInspeccion])
	}
	entrega, err := time.Parse(formatoFecha, fechaEntrega)
	if err != nil {
		return errors.New("la fecha debe tener el formato AAAA-MM-DD")
	}
	visita, _ := time.Parse(formatoFecha, fechaVisitaDe[idInspeccion])
	if entrega.Before(visita) {
		return errors.New("el informe no puede entregarse antes de la visita")
	}
	fechaInformeDe[idInspeccion] = fechaEntrega
	return nil
}

// ConsultarInspeccion devuelve los datos de la inspección y el estado de su informe.
func ConsultarInspeccion(idInspeccion int) (string, error) {
	estado, ok := estadoDe[idInspeccion]
	if !ok {
		return "", ErrInspeccionNoEncontrada
	}
	linea := fmt.Sprintf("#%d visita %s | cliente #%d | %s",
		idInspeccion, fechaVisitaDe[idInspeccion], clienteDe[idInspeccion], estado)
	if estado == Programada {
		return linea, nil
	}

	if conformeDe[idInspeccion] {
		linea += " | conforme"
	} else {
		linea += " | NO conforme: " + hallazgosDe[idInspeccion]
	}
	if f := fechaInformeDe[idInspeccion]; f != "" {
		plazo := "a tiempo"
		if !informeATiempo(idInspeccion) {
			plazo = "fuera de plazo"
		}
		linea += fmt.Sprintf(" | informe entregado %s (%s)", f, plazo)
	} else {
		linea += " | informe PENDIENTE, límite " + fechaLimite(idInspeccion)
	}
	return linea, nil
}

// AgendaPorFecha devuelve las inspecciones programadas para una fecha (AAAA-MM-DD).
func AgendaPorFecha(fecha string) []string {
	return listar(func(id int) bool { return fechaVisitaDe[id] == fecha })
}

// ListarPorCliente devuelve todas las inspecciones de un cliente.
func ListarPorCliente(idCliente int) []string {
	return listar(func(id int) bool { return clienteDe[id] == idCliente })
}

// InformesPendientes devuelve las inspecciones realizadas cuyo informe técnico
// todavía no se ha entregado al cliente.
func InformesPendientes() []string {
	return listar(func(id int) bool { return estadoDe[id] == Realizada && fechaInformeDe[id] == "" })
}

// ListarNoConformes devuelve las inspecciones no conformes con sus hallazgos,
// para dar seguimiento a las acciones correctivas.
func ListarNoConformes() []string {
	return listar(func(id int) bool { return estadoDe[id] == Realizada && !conformeDe[id] })
}

// CumplimientoCalidad devuelve el porcentaje de inspecciones realizadas que
// resultaron conformes (0 a 100).
func CumplimientoCalidad() float64 {
	realizadas, conformes := 0, 0
	for id, estado := range estadoDe {
		if estado == Realizada {
			realizadas++
			if conformeDe[id] {
				conformes++
			}
		}
	}
	return porcentaje(conformes, realizadas)
}

// CumplimientoInformes devuelve el porcentaje de inspecciones realizadas cuyo
// informe técnico se entregó dentro del plazo (0 a 100). Con idCliente 0
// calcula el indicador de todos los clientes.
func CumplimientoInformes(idCliente int) float64 {
	realizadas, aTiempo := 0, 0
	for id, estado := range estadoDe {
		if estado != Realizada || (idCliente != 0 && clienteDe[id] != idCliente) {
			continue
		}
		realizadas++
		if fechaInformeDe[id] != "" && informeATiempo(id) {
			aTiempo++
		}
	}
	return porcentaje(aTiempo, realizadas)
}

// fechaLimite devuelve la fecha máxima para entregar el informe de una inspección.
func fechaLimite(id int) string {
	visita, _ := time.Parse(formatoFecha, fechaVisitaDe[id])
	return visita.AddDate(0, 0, plazoInformeDias).Format(formatoFecha)
}

// informeATiempo indica si el informe se entregó hasta la fecha límite.
// Las fechas AAAA-MM-DD se pueden comparar como texto.
func informeATiempo(id int) bool {
	return fechaInformeDe[id] <= fechaLimite(id)
}

// porcentaje calcula parte / total * 100; devuelve 0 si no hay datos.
func porcentaje(parte, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(parte) / float64(total) * 100
}

// listar devuelve, ordenadas por ID, las inspecciones que cumplen la condición.
func listar(incluir func(int) bool) []string {
	ids := make([]int, 0, len(estadoDe))
	for id := range estadoDe {
		if incluir(id) {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)

	lista := make([]string, 0, len(ids))
	for _, id := range ids {
		linea, _ := ConsultarInspeccion(id)
		lista = append(lista, linea)
	}
	return lista
}
