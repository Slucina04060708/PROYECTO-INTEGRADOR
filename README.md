# PROYECTO-INTEGRADOR
## Integrantes del Equipo
* Camila Isabel Paucar Guaranda
* Santiago Alejandro Lucina Mena
* Pamela Stefania Toapanta Pila

## Descripción del Proyecto
Este proyecto consiste en el diseño e implementación de una aplicación cliente-servidor desarrollada modularmente en *Go* para la gestión de correos electrónicos internos entre dos sedes de una organización. El sistema opera sobre una infraestructura de red interconectada mediante *OSPF*, integrando un servicio de autenticación y validación de tokens de sesión mediante *FastAPI*.

## Funcionalidades Principales:
* **Arquitectura Cliente-Servidor:** Desacoplamiento total entre el cliente y el servidor mediante peticiones TCP/HTTP.
* **Gestión de Correos:** Redacción, envío, consulta de bandejas de entrada y salida, lectura con actualización de estado y eliminación de mensajes.
* **Seguridad y Autenticación:** Intercepción y validación de tokens de sesión emitidos por FastAPI en cada transacción.
* **Modularidad y POO:** Aplicación de principios de Programación Orientada a Objetos en *Go* como estructuras, encapsulamiento, composición y polimorfismo mediante interfaces para el repositorio de datos.
* **Concurrencia:** Atención concurrente y segura de múltiples clientes mediante el uso de *goroutines* y canales.
## Tablero de seguimiento de tareas
| Por Hacer (To Do) | En Proceso (In Progress) | Completado (Done) |
| :--- | :--- | :--- |
| - Diseño de casos de uso y flujos de mensajería (Camila Paucar)<br><br>- Definición del modelo de datos e integración con FastAPI (Camila Paucar)<br><br>- Implementación de la persistencia de datos en Go (Santiago Lucina)<br><br>- Pruebas de concurrencia mediante goroutines y canales (Santiago Lucina) | - Configuración del repositorio Git inicial con README, licencia y tablero (Pamela Toapanta)<br><br>- Despliegue de la infraestructura de red en dos sedes con OpenWrt, OSPF y VLANs (Pamela Toapanta)<br><br>- Modelado de POO en Go, structs y composición (Santiago Lucina) | - Selección del sistema de correos internos y análisis del problema (Camila Paucar)<br><br>- Redacción de los requisitos funcionales de software (Santiago Lucina)<br><br>- Elaboración del cronograma y diagrama de Gantt (Pamela Toapanta)<br><br>- Definición de actores, roles y consideraciones éticas (Equipo)<br><br>- Creación del diagrama de clases UML preliminar (Equipo) |

## Video Explicativo:
https://drive.google.com/file/d/1fxkeK9a3nD_wC4BHaLp9tCqehMOz-OfU/view?usp=sharing

