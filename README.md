# PlagaControl

Sistema de gestión de servicios de control de plagas, escrito en Go.

| | |
|---|---|
| **Estudiante** | Angel Gabriel Jara Elizalde |
| **Asignatura** | Programación Orientada a Objetos (Go) · UIDE |
| **Facilitador** | Carlos Arévalo |
| **Etapa actual** | Aprendizaje Autónomo 1 · Planeación del software (noviembre 2026) |

## Objetivo

Organizar la operación de una empresa de control de plagas: registrar clientes y técnicos, programar y dar seguimiento a los servicios (control de rastreros, control de insectos voladores y control de roedores), controlar el inventario de productos químicos con alertas de desabastecimiento, y programar las inspecciones a cada cliente junto con la entrega de sus informes técnicos.

## Alcance

**Qué hace**

- Registra, consulta, modifica, elimina y lista clientes residenciales, comerciales e industriales.
- Registra técnicos y controla si están activos.
- Programa órdenes de servicio, asigna técnicos, controla sus estados y calcula su costo.
- Lleva el inventario de productos químicos: registro de productos con stock mínimo, ingresos (compras), consumos y alertas de próximo desabastecimiento.
- Programa inspecciones con fecha de visita a cada cliente, registra su resultado y la entrega del informe técnico, y calcula el cumplimiento de calidad y de entrega de informes.

**Qué no hace (fuera de alcance)**

- Facturación, cobros ni contabilidad.
- Rutas, mapas ni GPS.
- Aplicación móvil.
- Base de datos: los datos viven en memoria mientras el programa corre.
- Usuarios con contraseña ni roles.

## Módulos

| Módulo del sistema | Paquete Go | Responsabilidad |
|---|---|---|
| Clientes | `clientes` | Datos de los clientes y su tipo |
| Técnicos | `tecnicos` | Personal que ejecuta los servicios |
| Servicios | `servicios` | Órdenes de servicio, tipos, estados y costos |
| Productos | `productos` | Inventario de productos químicos y alertas de stock |
| Inspecciones | `inspecciones` | Visitas de inspección, informes técnicos e indicadores |

## Estructura del proyecto

```
controlplagas/
├── go.mod                         ← un solo módulo de Go
├── README.md
├── main.go                        ← package main: arranque y menú
├── menu.go                        ← package main: acciones de cada opción del menú
├── datos.go                       ← package main: datos de ejemplo y coordinación
├── clientes/clientes.go           ← package clientes
├── tecnicos/tecnicos.go           ← package tecnicos
├── servicios/servicios.go         ← package servicios
├── productos/productos.go         ← package productos
└── inspecciones/inspecciones.go   ← package inspecciones
```

Los cinco paquetes no se importan entre sí: `main` es el único que los importa y los coordina. Por ejemplo, antes de programar un servicio, `main` verifica con `clientes.ExisteCliente` y `tecnicos.EstaActivo`, y solo entonces llama a `servicios.ProgramarServicio`. Así no hay dependencias circulares, que Go no permite.

## Menú del programa

| Grupo | Opciones |
|---|---|
| Clientes y técnicos | 1 Registrar cliente · 2 Listar clientes · 3 Registrar técnico · 4 Listar técnicos activos |
| Servicios | 5 Programar servicio · 6 Agenda por fecha · 7 Calcular costo |
| Inventario | 8 Registrar producto · 9 Ingresar stock · 10 Registrar consumo · 11 Inventario y alertas |
| Inspecciones e informes | 12 Programar inspección · 13 Agenda por fecha · 14 Inspecciones de un cliente · 15 Registrar resultado · 16 Registrar entrega de informe · 17 Indicadores |

## Cómo ejecutarlo

Requiere Go 1.22 o superior.

```bash
go build ./...   # compila todos los paquetes
go run .         # ejecuta el menú con datos de ejemplo
```

## Paquetes utilizados

Solo la librería estándar de Go: `fmt`, `errors`, `strings`, `strconv`, `sort`, `time`, `bufio` y `os`. No se usan paquetes de terceros.

## Próximas etapas

- **Aprendizaje Autónomo 2:** structs (`Cliente`, `Tecnico`, `Producto`, `Inspeccion`…) con constructores, métodos y encapsulación, y una interfaz `Servicio` implementada por `ControlRastreros`, `ControlVoladores` y `ControlRoedores`.
- **Evaluación final:** el sistema como aplicación web con al menos 8 servicios JSON, concurrencia segura y pruebas.
