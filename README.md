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

## Video Explicativo:
