// Package servicios gestiona las órdenes de servicio: control de rastreros,
// control de insectos voladores y control de roedores. Es el núcleo del sistema.
//
// Este paquete no importa clientes ni tecnicos: main verifica que el
// cliente y el técnico existan antes de llamar a estas funciones. Así se
// evitan dependencias circulares entre paquetes.
package servicios

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Estados posibles de una orden de servicio.
const (
	Programado = "programado"
	EnCurso    = "en curso"
	Finalizado = "finalizado"
	Cancelado  = "cancelado"
)

// formatoFecha es el formato de fecha que usa todo el sistema (AAAA-MM-DD).
const formatoFecha = "2006-01-02"

// costoMinimo es lo mínimo que se cobra por una visita, sin importar el área (USD).
const costoMinimo = 25.0

// ErrServicioNoEncontrado se devuelve cuando el ID no corresponde a ningún servicio.
var ErrServicioNoEncontrado = errors.New("servicio no encontrado")

// tarifas guarda el precio por metro cuadrado de cada tipo de servicio (USD).
// En el Aprendizaje Autónomo 2 cada tipo será un struct que implemente
// la interfaz Servicio con su propio cálculo de costo.
var tarifas = map[string]float64{
	"rastreros": 0.80,
	"voladores": 0.70,
	"roedores":  0.50,
}

// nombresTipo guarda el nombre completo de cada tipo, para mostrarlo al usuario.
var nombresTipo = map[string]string{
	"rastreros": "control de rastreros",
	"voladores": "control de insectos voladores",
	"roedores":  "control de roedores",
}

// ordenTipos fija el orden en que se muestran los tipos de servicio.
var ordenTipos = []string{"rastreros", "voladores", "roedores"}

// transiciones define a qué estados se puede pasar desde cada estado.
var transiciones = map[string][]string{
	Programado: {EnCurso, Cancelado},
	EnCurso:    {Finalizado},
}

// Almacenamiento en memoria, indexado por el ID del servicio.
var (
	clienteDe   = map[int]int{}
	tipoDe      = map[int]string{}
	fechaDe     = map[int]string{}
	estadoDe    = map[int]string{}
	tecnicoDe   = map[int]int{}
	siguienteID = 1
)

// ProgramarServicio crea una orden de servicio en estado "programado" y devuelve su ID.
// La fecha va en formato AAAA-MM-DD.
func ProgramarServicio(idCliente int, tipo, fecha string) (int, error) {
	tipo = normalizarTipo(tipo)
	if _, ok := tarifas[tipo]; !ok {
		return 0, fmt.Errorf("tipo de servicio %q no válido", tipo)
	}
	if _, err := time.Parse(formatoFecha, fecha); err != nil {
		return 0, errors.New("la fecha debe tener el formato AAAA-MM-DD")
	}

	id := siguienteID
	siguienteID++
	clienteDe[id] = idCliente
	tipoDe[id] = tipo
	fechaDe[id] = fecha
	estadoDe[id] = Programado
	return id, nil
}

// AsignarTecnico asigna el técnico que ejecutará el servicio.
// Solo se puede asignar mientras el servicio está programado.
func AsignarTecnico(idServicio, idTecnico int) error {
	estado, ok := estadoDe[idServicio]
	if !ok {
		return ErrServicioNoEncontrado
	}
	if estado != Programado {
		return fmt.Errorf("no se puede asignar técnico a un servicio %s", estado)
	}
	tecnicoDe[idServicio] = idTecnico
	return nil
}

// CambiarEstadoServicio mueve el servicio al estado indicado si la transición es válida:
// programado → en curso o cancelado; en curso → finalizado.
// Para pasar a "en curso" el servicio debe tener un técnico asignado.
func CambiarEstadoServicio(idServicio int, estado string) error {
	actual, ok := estadoDe[idServicio]
	if !ok {
		return ErrServicioNoEncontrado
	}
	estado = strings.ToLower(strings.TrimSpace(estado))
	if !puedePasar(actual, estado) {
		return fmt.Errorf("no se puede pasar de %q a %q", actual, estado)
	}
	if estado == EnCurso && tecnicoDe[idServicio] == 0 {
		return errors.New("asigne un técnico antes de iniciar el servicio")
	}
	estadoDe[idServicio] = estado
	return nil
}

// EstadoServicio devuelve el estado actual del servicio.
// main la usa antes de registrar una inspección.
func EstadoServicio(idServicio int) (string, error) {
	estado, ok := estadoDe[idServicio]
	if !ok {
		return "", ErrServicioNoEncontrado
	}
	return estado, nil
}

// CalcularCosto calcula el costo de un servicio según su tipo y el área en m².
// Nunca cobra menos del costo mínimo por visita.
func CalcularCosto(tipo string, metrosCuadrados float64) (float64, error) {
	tarifa, ok := tarifas[normalizarTipo(tipo)]
	if !ok {
		return 0, fmt.Errorf("tipo de servicio %q no válido", tipo)
	}
	if metrosCuadrados <= 0 {
		return 0, errors.New("el área debe ser mayor que cero")
	}
	costo := tarifa * metrosCuadrados
	if costo < costoMinimo {
		costo = costoMinimo
	}
	return costo, nil
}

// ConsultarServicio devuelve los datos del servicio en una línea de texto.
func ConsultarServicio(idServicio int) (string, error) {
	if _, ok := estadoDe[idServicio]; !ok {
		return "", ErrServicioNoEncontrado
	}
	tecnico := "sin asignar"
	if t := tecnicoDe[idServicio]; t != 0 {
		tecnico = fmt.Sprintf("técnico #%d", t)
	}
	return fmt.Sprintf("#%d %s | %s | cliente #%d | %s | %s",
		idServicio, fechaDe[idServicio], nombresTipo[tipoDe[idServicio]],
		clienteDe[idServicio], tecnico, estadoDe[idServicio]), nil
}

// TiposDisponibles devuelve los tipos de servicio con su tarifa, para mostrarlos en el menú.
func TiposDisponibles() []string {
	lista := make([]string, 0, len(ordenTipos))
	for _, t := range ordenTipos {
		lista = append(lista, fmt.Sprintf("%s: %s (USD %.2f por m²)", t, nombresTipo[t], tarifas[t]))
	}
	return lista
}

// ListarPorFecha devuelve la agenda de servicios de un día (AAAA-MM-DD).
func ListarPorFecha(fecha string) []string {
	ids := make([]int, 0)
	for id, f := range fechaDe {
		if f == fecha {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)

	agenda := make([]string, 0, len(ids))
	for _, id := range ids {
		linea, _ := ConsultarServicio(id)
		agenda = append(agenda, linea)
	}
	return agenda
}

// puedePasar indica si existe la transición de un estado a otro.
func puedePasar(desde, hacia string) bool {
	for _, permitido := range transiciones[desde] {
		if permitido == hacia {
			return true
		}
	}
	return false
}

// normalizarTipo convierte lo que escribe el usuario en la clave del tipo de
// servicio. Acepta formas como "rastreros", "Control de rastreros",
// "insectos voladores" o "ROEDORES".
func normalizarTipo(tipo string) string {
	reemplazos := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u")
	t := reemplazos.Replace(strings.ToLower(strings.TrimSpace(tipo)))
	switch {
	case strings.Contains(t, "rastrero"):
		return "rastreros"
	case strings.Contains(t, "volador"):
		return "voladores"
	case strings.Contains(t, "roedor"):
		return "roedores"
	}
	return t
}
