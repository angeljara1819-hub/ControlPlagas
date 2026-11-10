// Package inspecciones gestiona el control de calidad posterior a cada
// servicio: verificación en sitio, hallazgos, conformidad y seguimiento.
package inspecciones

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

const formatoFecha = "2006-01-02"

// Almacenamiento en memoria, indexado por el ID de la inspección.
var (
	servicioDe  = map[int]int{}
	hallazgosDe = map[int]string{}
	conformeDe  = map[int]bool{}
	seguimiento = map[int]string{} // ID del servicio → fecha de la próxima visita
	siguienteID = 1
)

// RegistrarInspeccion registra la verificación de calidad de un servicio y devuelve su ID.
// Si el servicio no es conforme, los hallazgos son obligatorios.
func RegistrarInspeccion(idServicio int, hallazgos string, conforme bool) (int, error) {
	if idServicio <= 0 {
		return 0, errors.New("el ID del servicio no es válido")
	}
	hallazgos = strings.TrimSpace(hallazgos)
	if !conforme && hallazgos == "" {
		return 0, errors.New("una inspección no conforme debe describir los hallazgos")
	}

	id := siguienteID
	siguienteID++
	servicioDe[id] = idServicio
	hallazgosDe[id] = hallazgos
	conformeDe[id] = conforme
	return id, nil
}

// ProgramarSeguimiento calcula la fecha de la próxima visita sumando días a la
// fecha del servicio (AAAA-MM-DD), la guarda y la devuelve.
func ProgramarSeguimiento(idServicio int, fechaBase string, dias int) (string, error) {
	if dias <= 0 {
		return "", errors.New("los días de seguimiento deben ser mayores que cero")
	}
	base, err := time.Parse(formatoFecha, fechaBase)
	if err != nil {
		return "", errors.New("la fecha debe tener el formato AAAA-MM-DD")
	}
	proxima := base.AddDate(0, 0, dias).Format(formatoFecha)
	seguimiento[idServicio] = proxima
	return proxima, nil
}

// CalcularCumplimiento devuelve el porcentaje de inspecciones conformes (0 a 100).
// Es el indicador principal de calidad del servicio.
func CalcularCumplimiento() float64 {
	if len(conformeDe) == 0 {
		return 0
	}
	conformes := 0
	for _, ok := range conformeDe {
		if ok {
			conformes++
		}
	}
	return float64(conformes) / float64(len(conformeDe)) * 100
}

// ListarNoConformes devuelve las inspecciones no conformes con sus hallazgos,
// para dar seguimiento a las acciones correctivas.
func ListarNoConformes() []string {
	ids := make([]int, 0)
	for id, ok := range conformeDe {
		if !ok {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)

	lista := make([]string, 0, len(ids))
	for _, id := range ids {
		lista = append(lista, fmt.Sprintf("inspección #%d | servicio #%d | %s",
			id, servicioDe[id], hallazgosDe[id]))
	}
	return lista
}
