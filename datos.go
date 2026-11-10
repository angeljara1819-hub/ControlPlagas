package main

// datos.go carga datos de ejemplo al iniciar y muestra cómo main coordina los
// paquetes: verifica en un paquete antes de llamar a otro, para que los
// paquetes no tengan que importarse entre sí.

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
	registrarCliente("Bodegas Andinas S.A.", "1791234567001", "022987654", "Parque Industrial Sur", "industrial")

	t1 := registrarTecnico("Luis Paredes", "1718293045", "fumigacion")
	t2 := registrarTecnico("Andrea Cuenca", "1723456789", "desratizacion")

	registrarProducto("Insecticida piretroide", "REG-INS-0451", 12)
	registrarProducto("Rodenticida en bloque", "REG-ROD-0213", 3.5)
	registrarProducto("Desinfectante amonio cuaternario", "REG-DES-0877", 20)

	s1 := programarServicio(c1, "fumigacion", "2026-11-16", t1)
	s2 := programarServicio(c2, "desratizacion", "2026-11-16", t2)
	programarServicio(c1, "desinfeccion", "2026-11-18", t1)
	programarServicio(99, "fumigacion", "2026-11-18", t1) // cliente inexistente: se rechaza

	finalizarEInspeccionar(s1, "", true, "2026-11-16", 30)
	finalizarEInspeccionar(s2, "Rastros activos en bodega posterior", false, "2026-11-16", 15)

	fmt.Println("Datos de ejemplo cargados.")
}

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

func registrarProducto(nombre, registro string, stock float64) {
	if _, err := productos.RegistrarProducto(nombre, registro, stock); err != nil {
		fmt.Println("No se registró el producto:", err)
	}
}

// programarServicio comprueba primero el cliente y el técnico en sus paquetes
// y solo después crea la orden en el paquete servicios.
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

// finalizarEInspeccionar lleva el servicio hasta "finalizado" y registra su
// inspección de calidad con la fecha de seguimiento.
func finalizarEInspeccionar(idServicio int, hallazgos string, conforme bool, fecha string, diasSeguimiento int) {
	for _, estado := range []string{servicios.EnCurso, servicios.Finalizado} {
		if err := servicios.CambiarEstadoServicio(idServicio, estado); err != nil {
			fmt.Println("No se actualizó el servicio:", err)
			return
		}
	}
	if _, err := inspecciones.RegistrarInspeccion(idServicio, hallazgos, conforme); err != nil {
		fmt.Println("No se registró la inspección:", err)
		return
	}
	if _, err := inspecciones.ProgramarSeguimiento(idServicio, fecha, diasSeguimiento); err != nil {
		fmt.Println("No se programó el seguimiento:", err)
	}
}
