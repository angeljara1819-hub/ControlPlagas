package main

// menu.go contiene lo que hace cada opción del menú: pide los datos al
// usuario, llama a las funciones de los paquetes y muestra el resultado o el
// error. No contiene reglas del negocio: esas viven en cada paquete.

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

// ---------- Clientes y técnicos ----------

// opcionRegistrarCliente pide los datos de un cliente nuevo y lo registra.
// Si algún dato no es válido, muestra el error y no guarda nada.
func opcionRegistrarCliente() {
	nombre := leer("Nombre o razón social: ")
	ruc := leer("RUC (13 dígitos): ")
	telefono := leer("Teléfono: ")
	direccion := leer("Dirección: ")
	tipo := leer("Tipo (residencial, comercial, industrial): ")

	id, err := clientes.RegistrarCliente(nombre, ruc, telefono, direccion, tipo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	datos, _ := clientes.ConsultarCliente(id)
	fmt.Println("Cliente registrado:", datos)
}

func opcionListarClientes() {
	tipo := leer("Tipo (residencial, comercial, industrial o Enter = todos): ")
	imprimir("Clientes", clientes.ListarClientes(tipo))
}

// opcionRegistrarTecnico pide los datos de un técnico nuevo y lo registra como activo.
func opcionRegistrarTecnico() {
	nombre := leer("Nombre completo: ")
	cedula := leer("Cédula (10 dígitos): ")
	especialidad := leer("Especialidad (rastreros, voladores, roedores): ")

	id, err := tecnicos.RegistrarTecnico(nombre, cedula, especialidad)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	datos, _ := tecnicos.ConsultarTecnico(id)
	fmt.Println("Técnico registrado:", datos)
}

func opcionListarTecnicos() {
	imprimir("Técnicos activos", tecnicos.ListarActivos())
}

// ---------- Servicios ----------

// opcionProgramarServicio muestra clientes, tipos y técnicos, pide los datos
// de la orden y la programa con programarServicio (datos.go), que verifica
// cliente y técnico antes de crearla.
func opcionProgramarServicio() {
	imprimir("Clientes", clientes.ListarClientes(""))
	idCliente := leerEntero("ID del cliente: ")
	imprimir("Tipos de servicio", servicios.TiposDisponibles())
	tipo := leer("Tipo de servicio: ")
	fecha := leer("Fecha del servicio (AAAA-MM-DD): ")
	imprimir("Técnicos activos", tecnicos.ListarActivos())
	idTecnico := leerEntero("ID del técnico: ")

	if id := programarServicio(idCliente, tipo, fecha, idTecnico); id != 0 {
		datos, _ := servicios.ConsultarServicio(id)
		fmt.Println("Servicio programado:", datos)
	}
}

func opcionAgendaServicios() {
	imprimir("Agenda de servicios", servicios.ListarPorFecha(leer("Fecha (AAAA-MM-DD): ")))
}

func opcionCalcularCosto() {
	imprimir("Tipos de servicio", servicios.TiposDisponibles())
	tipo := leer("Tipo de servicio: ")
	costo, err := servicios.CalcularCosto(tipo, leerNumero("Área en m²: "))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Costo estimado: USD %.2f\n", costo)
}

// ---------- Inventario de productos ----------

// opcionRegistrarProducto agrega un producto nuevo al inventario con su stock
// inicial y su stock mínimo.
func opcionRegistrarProducto() {
	nombre := leer("Nombre del producto: ")
	registro := leer("Registro sanitario: ")
	unidad := leer("Unidad de medida (litros, kg, unidades): ")
	inicial := leerNumero("Stock inicial: ")
	minimo := leerNumero("Stock mínimo (nivel de alerta): ")

	id, err := productos.RegistrarProducto(nombre, registro, unidad, inicial, minimo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	datos, _ := productos.ConsultarProducto(id)
	fmt.Println("Producto registrado:", datos)
	alertarSiHayQueReponer(id)
}

// opcionIngresarStock suma una compra o reposición al stock de un producto.
func opcionIngresarStock() {
	imprimir("Inventario", productos.ListarInventario())
	id := leerEntero("ID del producto: ")
	stock, err := productos.IngresarStock(id, leerNumero("Cantidad que ingresa: "))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Ingreso registrado. Stock actual: %.2f\n", stock)
	alertarSiHayQueReponer(id)
}

// opcionRegistrarConsumo descuenta lo usado en un servicio y alerta si el
// producto quedó en su stock mínimo o por debajo.
func opcionRegistrarConsumo() {
	imprimir("Inventario", productos.ListarInventario())
	id := leerEntero("ID del producto: ")
	stock, err := productos.RegistrarConsumo(id, leerNumero("Cantidad usada: "))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Consumo registrado. Stock actual: %.2f\n", stock)
	alertarSiHayQueReponer(id)
}

func opcionVerInventario() {
	imprimir("Inventario", productos.ListarInventario())
	imprimir("Productos por reponer", productos.ProductosBajoMinimo())
}

// alertarSiHayQueReponer avisa un próximo desabastecimiento del producto.
func alertarSiHayQueReponer(idProducto int) {
	if reponer, _ := productos.NecesitaReposicion(idProducto); reponer {
		datos, _ := productos.ConsultarProducto(idProducto)
		fmt.Println("ALERTA de desabastecimiento:", datos)
	}
}

// mostrarAlertasDeStock lista al iniciar los productos que hay que reponer.
func mostrarAlertasDeStock() {
	if bajos := productos.ProductosBajoMinimo(); len(bajos) > 0 {
		imprimir("ALERTA: productos por reponer", bajos)
	}
}

// ---------- Inspecciones e informes técnicos ----------

// opcionProgramarInspeccion agenda una visita de inspección a un cliente.
func opcionProgramarInspeccion() {
	imprimir("Clientes", clientes.ListarClientes(""))
	idCliente := leerEntero("ID del cliente: ")
	fecha := leer("Fecha de la visita (AAAA-MM-DD): ")
	if id := programarInspeccion(idCliente, fecha); id != 0 {
		datos, _ := inspecciones.ConsultarInspeccion(id)
		fmt.Println("Inspección programada:", datos)
	}
}

func opcionAgendaInspecciones() {
	imprimir("Agenda de inspecciones", inspecciones.AgendaPorFecha(leer("Fecha (AAAA-MM-DD): ")))
}

// opcionInspeccionesDeCliente muestra las visitas de un cliente y su
// cumplimiento en la entrega de informes técnicos.
func opcionInspeccionesDeCliente() {
	idCliente := leerEntero("ID del cliente: ")
	datos, err := clientes.ConsultarCliente(idCliente)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Cliente:", datos)
	imprimir("Inspecciones", inspecciones.ListarPorCliente(idCliente))
	fmt.Printf("Informes entregados a tiempo: %.1f %%\n", inspecciones.CumplimientoInformes(idCliente))
}

// opcionRegistrarResultado registra si la visita fue conforme y sus hallazgos.
func opcionRegistrarResultado() {
	id := leerEntero("ID de la inspección: ")
	conforme := strings.ToLower(leer("¿Conforme? (s/n): ")) == "s"
	hallazgos := leer("Hallazgos (obligatorio si no es conforme): ")
	if err := inspecciones.RegistrarResultado(id, hallazgos, conforme); err != nil {
		fmt.Println("Error:", err)
		return
	}
	datos, _ := inspecciones.ConsultarInspeccion(id)
	fmt.Println("Resultado registrado:", datos)
}

// opcionRegistrarEntregaInforme registra la fecha de entrega del informe técnico.
func opcionRegistrarEntregaInforme() {
	imprimir("Informes pendientes", inspecciones.InformesPendientes())
	id := leerEntero("ID de la inspección: ")
	fecha := leer("Fecha de entrega del informe (AAAA-MM-DD): ")
	if err := inspecciones.RegistrarEntregaInforme(id, fecha); err != nil {
		fmt.Println("Error:", err)
		return
	}
	datos, _ := inspecciones.ConsultarInspeccion(id)
	fmt.Println("Entrega registrada:", datos)
}

// opcionIndicadores muestra los indicadores de calidad y de informes.
func opcionIndicadores() {
	fmt.Printf("Inspecciones conformes: %.1f %%\n", inspecciones.CumplimientoCalidad())
	fmt.Printf("Informes técnicos entregados a tiempo: %.1f %%\n", inspecciones.CumplimientoInformes(0))
	imprimir("No conformidades", inspecciones.ListarNoConformes())
	imprimir("Informes pendientes de entrega", inspecciones.InformesPendientes())
}

// ---------- Lectura e impresión ----------

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
	n, err := strconv.ParseFloat(strings.ReplaceAll(leer(mensaje), ",", "."), 64)
	if err != nil {
		return 0
	}
	return n
}

// leerEntero pide un número entero, como un ID.
func leerEntero(mensaje string) int {
	return int(leerNumero(mensaje))
}

func imprimir(titulo string, lineas []string) {
	fmt.Printf("--- %s (%d) ---\n", titulo, len(lineas))
	if len(lineas) == 0 {
		fmt.Println("(sin resultados)")
	}
	for _, l := range lineas {
		fmt.Println(l)
	}
}
