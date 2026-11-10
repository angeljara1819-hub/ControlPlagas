// PlagaControl: sistema de gestión de servicios de control de plagas.
//
// main.go solo arranca el programa y muestra el menú; la lógica vive en los
// paquetes clientes, tecnicos, servicios, productos e inspecciones.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"controlplagas/clientes"
	"controlplagas/inspecciones"
	"controlplagas/productos"
	"controlplagas/servicios"
	"controlplagas/tecnicos"
)

var entrada = bufio.NewScanner(os.Stdin)

func main() {
	fmt.Println("=== PlagaControl · Gestión de servicios de control de plagas ===")
	cargarDatosDeEjemplo()

	for {
		mostrarMenu()
		switch leer("Opción: ") {
		case "1":
			imprimir("Clientes", clientes.ListarClientes(leer("Tipo (vacío = todos): ")))
		case "2":
			imprimir("Técnicos activos", tecnicos.ListarActivos())
		case "3":
			imprimir("Agenda", servicios.ListarPorFecha(leer("Fecha (AAAA-MM-DD): ")))
		case "4":
			opcionCalcularCosto()
		case "5":
			opcionRegistrarConsumo()
		case "6":
			imprimir("Productos bajo el mínimo", productos.ProductosBajoMinimo(leerNumero("Stock mínimo: ")))
		case "7":
			imprimir("Inspecciones no conformes", inspecciones.ListarNoConformes())
			fmt.Printf("Cumplimiento de calidad: %.1f %%\n", inspecciones.CalcularCumplimiento())
		case "0":
			fmt.Println("Hasta pronto.")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

func mostrarMenu() {
	fmt.Println("\n1. Listar clientes")
	fmt.Println("2. Listar técnicos activos")
	fmt.Println("3. Agenda de servicios por fecha")
	fmt.Println("4. Calcular costo de un servicio")
	fmt.Println("5. Registrar consumo de producto")
	fmt.Println("6. Productos bajo el stock mínimo")
	fmt.Println("7. Calidad: no conformidades y cumplimiento")
	fmt.Println("0. Salir")
}

func opcionCalcularCosto() {
	tipo := leer("Tipo (fumigacion, desratizacion, desinfeccion): ")
	costo, err := servicios.CalcularCosto(tipo, leerNumero("Área en m²: "))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Costo estimado: USD %.2f\n", costo)
}

func opcionRegistrarConsumo() {
	id := int(leerNumero("ID del producto: "))
	if err := productos.RegistrarConsumo(id, leerNumero("Cantidad usada: ")); err != nil {
		fmt.Println("Error:", err)
		return
	}
	stock, _ := productos.ConsultarStock(id)
	fmt.Printf("Consumo registrado. Stock actual: %.2f\n", stock)
}

// leer muestra un mensaje y devuelve la línea que escribe el usuario.
func leer(mensaje string) string {
	fmt.Print(mensaje)
	if !entrada.Scan() {
		os.Exit(0) // fin de la entrada
	}
	return strings.TrimSpace(entrada.Text())
}

// leerNumero pide un número; si no es válido devuelve 0 y la función
// que lo recibe se encarga de rechazarlo con su propio error.
func leerNumero(mensaje string) float64 {
	n, err := strconv.ParseFloat(leer(mensaje), 64)
	if err != nil {
		return 0
	}
	return n
}

func imprimir(titulo string, lineas []string) {
	fmt.Printf("--- %s (%d) ---\n", titulo, len(lineas))
	for _, l := range lineas {
		fmt.Println(l)
	}
}
