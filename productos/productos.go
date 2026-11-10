// Package productos gestiona el inventario de productos químicos
// (insecticidas, rodenticidas, desinfectantes) que se usan en los servicios.
//
// Las cantidades se manejan en la unidad de cada producto (litros o kilogramos).
package productos

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrProductoNoEncontrado se devuelve cuando el ID no corresponde a ningún producto.
var ErrProductoNoEncontrado = errors.New("producto no encontrado")

// Almacenamiento en memoria, indexado por el ID del producto.
var (
	nombres     = map[int]string{}
	registros   = map[int]string{}
	stocks      = map[int]float64{}
	siguienteID = 1
)

// RegistrarProducto agrega un producto al inventario y devuelve su ID.
// El registro sanitario es obligatorio: no se usan productos sin registro.
func RegistrarProducto(nombre, registroSanitario string, stock float64) (int, error) {
	nombre = strings.TrimSpace(nombre)
	if nombre == "" {
		return 0, errors.New("el nombre del producto es obligatorio")
	}
	if strings.TrimSpace(registroSanitario) == "" {
		return 0, errors.New("el producto debe tener registro sanitario")
	}
	if stock < 0 {
		return 0, errors.New("el stock inicial no puede ser negativo")
	}

	id := siguienteID
	siguienteID++
	nombres[id] = nombre
	registros[id] = registroSanitario
	stocks[id] = stock
	return id, nil
}

// RegistrarConsumo descuenta del inventario lo que se usó en un servicio.
// Devuelve error si no hay stock suficiente.
func RegistrarConsumo(idProducto int, cantidad float64) error {
	stock, ok := stocks[idProducto]
	if !ok {
		return ErrProductoNoEncontrado
	}
	if cantidad <= 0 {
		return errors.New("la cantidad consumida debe ser mayor que cero")
	}
	if cantidad > stock {
		return fmt.Errorf("stock insuficiente de %s: hay %.2f y se piden %.2f",
			nombres[idProducto], stock, cantidad)
	}
	stocks[idProducto] = stock - cantidad
	return nil
}

// ReponerStock suma al inventario una compra o reposición del producto.
func ReponerStock(idProducto int, cantidad float64) error {
	if _, ok := stocks[idProducto]; !ok {
		return ErrProductoNoEncontrado
	}
	if cantidad <= 0 {
		return errors.New("la cantidad a reponer debe ser mayor que cero")
	}
	stocks[idProducto] += cantidad
	return nil
}

// ConsultarStock devuelve la cantidad disponible del producto.
func ConsultarStock(idProducto int) (float64, error) {
	stock, ok := stocks[idProducto]
	if !ok {
		return 0, ErrProductoNoEncontrado
	}
	return stock, nil
}

// ProductosBajoMinimo devuelve los productos cuyo stock es menor que el mínimo,
// para saber qué hay que reponer.
func ProductosBajoMinimo(minimo float64) []string {
	ids := make([]int, 0)
	for id, stock := range stocks {
		if stock < minimo {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)

	alertas := make([]string, 0, len(ids))
	for _, id := range ids {
		alertas = append(alertas, fmt.Sprintf("#%d %s | reg. %s | stock %.2f",
			id, nombres[id], registros[id], stocks[id]))
	}
	return alertas
}
