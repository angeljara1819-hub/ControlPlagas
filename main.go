// PlagaControl: sistema de gestión de servicios de control de plagas.
//
// main.go solo arranca el programa y muestra el menú. Las acciones de cada
// opción están en menu.go y la lógica del negocio vive en los paquetes
// clientes, tecnicos, servicios, productos e inspecciones.
package main

import "fmt"

func main() {
	fmt.Println("=== PlagaControl · Gestión de servicios de control de plagas ===")
	cargarDatosDeEjemplo()
	mostrarAlertasDeStock()

	for {
		mostrarMenu()
		switch leer("Opción: ") {
		case "1":
			opcionRegistrarCliente()
		case "2":
			opcionListarClientes()
		case "3":
			opcionRegistrarTecnico()
		case "4":
			opcionListarTecnicos()
		case "5":
			opcionProgramarServicio()
		case "6":
			opcionAgendaServicios()
		case "7":
			opcionCalcularCosto()
		case "8":
			opcionRegistrarProducto()
		case "9":
			opcionIngresarStock()
		case "10":
			opcionRegistrarConsumo()
		case "11":
			opcionVerInventario()
		case "12":
			opcionProgramarInspeccion()
		case "13":
			opcionAgendaInspecciones()
		case "14":
			opcionInspeccionesDeCliente()
		case "15":
			opcionRegistrarResultado()
		case "16":
			opcionRegistrarEntregaInforme()
		case "17":
			opcionIndicadores()
		case "0":
			fmt.Println("Hasta pronto.")
			return
		default:
			fmt.Println("Opción no válida.")
		}
	}
}

func mostrarMenu() {
	fmt.Println("\n--- CLIENTES Y TÉCNICOS ---")
	fmt.Println(" 1. Registrar cliente           2. Listar clientes")
	fmt.Println(" 3. Registrar técnico           4. Listar técnicos activos")
	fmt.Println("--- SERVICIOS ---")
	fmt.Println(" 5. Programar servicio          6. Agenda de servicios por fecha")
	fmt.Println(" 7. Calcular costo de un servicio")
	fmt.Println("--- INVENTARIO DE PRODUCTOS ---")
	fmt.Println(" 8. Registrar producto nuevo    9. Ingresar stock (compra)")
	fmt.Println("10. Registrar consumo          11. Ver inventario y alertas")
	fmt.Println("--- INSPECCIONES E INFORMES TÉCNICOS ---")
	fmt.Println("12. Programar inspección       13. Agenda de inspecciones por fecha")
	fmt.Println("14. Inspecciones de un cliente 15. Registrar resultado de inspección")
	fmt.Println("16. Registrar entrega de informe técnico")
	fmt.Println("17. Indicadores de calidad e informes")
	fmt.Println(" 0. Salir")
}
