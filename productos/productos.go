// Package productos gestiona el inventario de productos químicos
// (insecticidas, rodenticidas, cebos) que se usan en los servicios.
//
// Cada producto tiene un stock mínimo: cuando el stock llega a ese nivel o
// baja de él, el sistema alerta un próximo desabastecimiento.
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
	unidades    = map[int]string{}
	stocks      = map[int]float64{}
	minimos     = map[int]float64{}
	siguienteID = 1
)

// RegistrarProducto agrega un producto nuevo al inventario y devuelve su ID.
// La unidad es la forma de medirlo (litros, kg, unidades) y el stock mínimo es
// el nivel desde el cual el sistema alerta que hay que reponerlo.
// No se aceptan productos sin registro sanitario ni registros repetidos.
func RegistrarProducto(nombre, registroSanitario, unidad string, stockInicial, stockMinimo float64) (int, error) {
	nombre = strings.TrimSpace(nombre)
	registroSanitario = strings.ToUpper(strings.TrimSpace(registroSanitario))
	unidad = strings.ToLower(strings.TrimSpace(unidad))

	if nombre == "" {
		return 0, errors.New("el nombre del producto es obligatorio")
	}
	if registroSanitario == "" {
		return 0, errors.New("el producto debe tener registro sanitario")
	}
	if id := buscarPorRegistro(registroSanitario); id != 0 {
		return 0, fmt.Errorf("el registro %s ya pertenece al producto #%d: use la opción de ingresar stock", registroSanitario, id)
	}
	if unidad == "" {
		return 0, errors.New("indique la unidad de medida (litros, kg, unidades)")
	}
	if stockInicial < 0 || stockMinimo < 0 {
		return 0, errors.New("el stock inicial y el stock mínimo no pueden ser negativos")
	}

	id := siguienteID
	siguienteID++
	nombres[id] = nombre
	registros[id] = registroSanitario
	unidades[id] = unidad
	stocks[id] = stockInicial
	minimos[id] = stockMinimo
	return id, nil
}

// IngresarStock suma al inventario una compra o reposición del producto
// y devuelve el stock resultante.
func IngresarStock(idProducto int, cantidad float64) (float64, error) {
	if _, ok := stocks[idProducto]; !ok {
		return 0, ErrProductoNoEncontrado
	}
	if cantidad <= 0 {
		return 0, errors.New("la cantidad a ingresar debe ser mayor que cero")
	}
	stocks[idProducto] += cantidad
	return stocks[idProducto], nil
}

// RegistrarConsumo descuenta del inventario lo que se usó en un servicio
// y devuelve el stock que queda. Da error si no hay stock suficiente.
func RegistrarConsumo(idProducto int, cantidad float64) (float64, error) {
	stock, ok := stocks[idProducto]
	if !ok {
		return 0, ErrProductoNoEncontrado
	}
	if cantidad <= 0 {
		return 0, errors.New("la cantidad consumida debe ser mayor que cero")
	}
	if cantidad > stock {
		return stock, fmt.Errorf("stock insuficiente de %s: hay %.2f %s y se piden %.2f",
			nombres[idProducto], stock, unidades[idProducto], cantidad)
	}
	stocks[idProducto] = stock - cantidad
	return stocks[idProducto], nil
}

// NecesitaReposicion indica si el stock del producto está en su mínimo o por debajo.
func NecesitaReposicion(idProducto int) (bool, error) {
	stock, ok := stocks[idProducto]
	if !ok {
		return false, ErrProductoNoEncontrado
	}
	return stock <= minimos[idProducto], nil
}

// ConsultarStock devuelve la cantidad disponible del producto.
func ConsultarStock(idProducto int) (float64, error) {
	stock, ok := stocks[idProducto]
	if !ok {
		return 0, ErrProductoNoEncontrado
	}
	return stock, nil
}

// ConsultarProducto devuelve los datos del producto y el estado de su stock.
func ConsultarProducto(idProducto int) (string, error) {
	if _, ok := stocks[idProducto]; !ok {
		return "", ErrProductoNoEncontrado
	}
	u := unidades[idProducto]
	return fmt.Sprintf("#%d %s | reg. %s | stock %.2f %s | mínimo %.2f %s | %s",
		idProducto, nombres[idProducto], registros[idProducto],
		stocks[idProducto], u, minimos[idProducto], u, nivelStock(idProducto)), nil
}

// ListarInventario devuelve todos los productos con su stock, ordenados por ID.
func ListarInventario() []string {
	return listar(func(int) bool { return true })
}

// ProductosBajoMinimo devuelve los productos en su stock mínimo o por debajo,
// es decir, los que hay que reponer para evitar un desabastecimiento.
func ProductosBajoMinimo() []string {
	return listar(func(id int) bool { return stocks[id] <= minimos[id] })
}

// nivelStock clasifica el stock del producto: AGOTADO, REPONER (en el mínimo
// o por debajo) o normal.
func nivelStock(id int) string {
	switch {
	case stocks[id] == 0:
		return "AGOTADO"
	case stocks[id] <= minimos[id]:
		return "REPONER"
	default:
		return "normal"
	}
}

// listar devuelve, ordenados por ID, los productos que cumplen la condición.
func listar(incluir func(int) bool) []string {
	ids := make([]int, 0, len(stocks))
	for id := range stocks {
		if incluir(id) {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)

	lista := make([]string, 0, len(ids))
	for _, id := range ids {
		linea, _ := ConsultarProducto(id)
		lista = append(lista, linea)
	}
	return lista
}

// buscarPorRegistro devuelve el ID del producto con ese registro sanitario, o 0 si no existe.
func buscarPorRegistro(registro string) int {
	for id, r := range registros {
		if r == registro {
			return id
		}
	}
	return 0
}
