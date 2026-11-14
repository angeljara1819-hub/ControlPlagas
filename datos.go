package main

// datos.go carga datos de ejemplo al iniciar y contiene las funciones que
// coordinan varios paquetes: verifican en un paquete antes de llamar a otro,
// para que los paquetes no tengan que importarse entre sí.

import (
	"fmt"

	"controlplagas/clientes"
	"controlplagas/inspecciones"
	"controlplagas/productos"
	"controlplagas/servicios"
	"controlplagas/tecnicos"
)

func cargarDatosDeEjemplo() {
	c1 := registrarCliente("Panadería El Trigal", "1790012345001", "022456789", "Av. América N34-12", "comercial")
	c2 := registrarCliente("María Torres", "1712345678001", "0998765432", "Conocoto, calle B 15", "residencial")
	c3 := registrarCliente("Bodegas Andinas S.A.", "1791234567001", "022987654", "Parque Industrial Sur", "industrial")

	t1 := registrarTecnico("Luis Paredes", "1718293045", "rastreros")
	t2 := registrarTecnico("Andrea Cuenca", "1723456789", "roedores")

	registrarProducto("Insecticida piretroide", "REG-INS-0451", "litros", 12, 4)
	registrarProducto("Rodenticida en bloque", "REG-ROD-0213", "kg", 3.5, 5) // ya bajo el mínimo
	registrarProducto("Cebo en gel para cucarachas", "REG-INS-0388", "kg", 2, 0.5)

	s1 := programarServicio(c1, "rastreros", "2026-11-16", t1)
	s2 := programarServicio(c3, "roedores", "2026-11-16", t2)
	programarServicio(c2, "voladores", "2026-11-18", t1)
	programarServicio(99, "rastreros", "2026-11-18", t1) // cliente inexistente: se rechaza
	finalizarServicio(s1)
	finalizarServicio(s2)

	// Inspección conforme, con informe entregado dentro del plazo.
	i1 := programarInspeccion(c1, "2026-11-17")
	registrarResultado(i1, "", true)
	registrarEntregaInforme(i1, "2026-11-19")
	// Inspección no conforme, con el informe todavía pendiente.
	i2 := programarInspeccion(c3, "2026-11-17")
	registrarResultado(i2, "Rastros activos de roedores en bodega posterior", false)
	// Visita de seguimiento programada.
	programarInspeccion(c3, "2026-11-24")

	fmt.Println("Datos de ejemplo cargados.")
}

// ---------- Coordinación entre paquetes ----------

// programarServicio comprueba primero el cliente y el técnico en sus paquetes
// y solo después crea la orden en el paquete servicios. Devuelve 0 si falla.
func programarServicio(idCliente int, tipo, fecha string, idTecnico int) int {
	if !clientes.ExisteCliente(idCliente) {
		fmt.Printf("No se programó el servicio: el cliente #%d no existe\n", idCliente)
		return 0
	}
	if !tecnicos.EstaActivo(idTecnico) {
		fmt.Printf("No se programó el servicio: el técnico #%d no está activo\n", idTecnico)
		return 0
	}
	id, err := servicios.ProgramarServicio(idCliente, tipo, fecha)
	if err != nil {
		fmt.Println("No se programó el servicio:", err)
		return 0
	}
	if err := servicios.AsignarTecnico(id, idTecnico); err != nil {
		fmt.Println("No se asignó el técnico:", err)
	}
	return id
}

// programarInspeccion comprueba que el cliente exista y solo después agenda
// la visita en el paquete inspecciones. Devuelve 0 si falla.
func programarInspeccion(idCliente int, fecha string) int {
	if !clientes.ExisteCliente(idCliente) {
		fmt.Printf("No se programó la inspección: el cliente #%d no existe\n", idCliente)
		return 0
	}
	id, err := inspecciones.ProgramarInspeccion(idCliente, fecha)
	if err != nil {
		fmt.Println("No se programó la inspección:", err)
		return 0
	}
	return id
}

// finalizarServicio lleva el servicio de "programado" a "finalizado",
// pasando por "en curso" como exige el proceso.
func finalizarServicio(idServicio int) {
	for _, estado := range []string{servicios.EnCurso, servicios.Finalizado} {
		if err := servicios.CambiarEstadoServicio(idServicio, estado); err != nil {
			fmt.Println("No se actualizó el servicio:", err)
			return
		}
	}
}

// ---------- Registro de datos de ejemplo ----------

func registrarCliente(nombre, ruc, telefono, direccion, tipo string) int {
	id, err := clientes.RegistrarCliente(nombre, ruc, telefono, direccion, tipo)
	if err != nil {
		fmt.Println("No se registró el cliente:", err)
	}
	return id
}

func registrarTecnico(nombre, cedula, especialidad string) int {
	id, err := tecnicos.RegistrarTecnico(nombre, cedula, especialidad)
	if err != nil {
		fmt.Println("No se registró el técnico:", err)
	}
	return id
}

func registrarProducto(nombre, registro, unidad string, stock, minimo float64) {
	if _, err := productos.RegistrarProducto(nombre, registro, unidad, stock, minimo); err != nil {
		fmt.Println("No se registró el producto:", err)
	}
}

func registrarResultado(idInspeccion int, hallazgos string, conforme bool) {
	if err := inspecciones.RegistrarResultado(idInspeccion, hallazgos, conforme); err != nil {
		fmt.Println("No se registró el resultado:", err)
	}
}

func registrarEntregaInforme(idInspeccion int, fecha string) {
	if err := inspecciones.RegistrarEntregaInforme(idInspeccion, fecha); err != nil {
		fmt.Println("No se registró el informe:", err)
	}
}
