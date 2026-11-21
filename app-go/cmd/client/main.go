package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/dominio"
	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/red/cliente"
)

func leerLinea(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func registrarEnFastAPI(usuario, password, rol string) {
	body, _ := json.Marshal(map[string]string{
		"usuario":  usuario,
		"password": password,
		"rol":      rol,
	})
	resp, err := http.Post("http://127.0.0.1:8000/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("[ERROR 503] No se pudo conectar con el servicio de autenticación.")
		return
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if resp.StatusCode == 200 {
		fmt.Printf("[OK 200] %v (Rol: %v)\n", res["mensaje"], res["rol"])
	} else {
		fmt.Printf("[ERROR %d] %v\n", resp.StatusCode, res["detail"])
	}
}

func main() {
	cli := cliente.NuevoClienteCorreo("http://127.0.0.1:8080")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if cli.TokenSesion == "" {
			fmt.Println("\n==========================================")
			fmt.Println(" SISTEMA DE CORREOS INTERNOS | ACCESO")
			fmt.Println("==========================================")
			fmt.Println("[1] Iniciar sesión")
			fmt.Println("[2] Registrar nuevo usuario")
			fmt.Println("[0] Salir")
			opcionInicio := leerLinea(scanner, "> Seleccione opción: ")

			if opcionInicio == "0" {
				fmt.Println("Saliendo del sistema...")
				return
			}

			if opcionInicio == "2" {
				fmt.Println("\n--- REGISTRO DE USUARIO (Política C-01) ---")
				usr := leerLinea(scanner, "Nuevo usuario: ")
				pwd := leerLinea(scanner, "Contraseña: ")
				rol := leerLinea(scanner, "Rol (operativo / administrativo / supervisor): ")
				registrarEnFastAPI(usr, pwd, rol)
				continue
			}

			if opcionInicio == "1" {
				usr := leerLinea(scanner, "Usuario: ")
				pwd := leerLinea(scanner, "Contraseña: ")
				if usr == "" || pwd == "" {
					fmt.Println("[ERROR] Complete usuario y contraseña.")
					continue
				}
				res := cli.EnviarSolicitud(dominio.SolicitudCorreo{Accion: "login", Usuario: usr, Password: pwd})
				if res.Exito {
					fmt.Printf("[OK] %s | Rol: %s | Token asignado en memoria.\n", res.Mensaje, res.Rol)
				} else {
					fmt.Printf("[ERROR %d] %s\n", res.Codigo, res.Mensaje)
				}
			}
			continue
		}

		fmt.Printf("\n--- USUARIO: %s (%s) | SESIÓN: [ACTIVA] ---\n", cli.UsuarioActivo, cli.RolActivo)
		fmt.Println("[1] Redactar y enviar correo")
		fmt.Println("[2] Consultar bandeja de entrada")
		fmt.Println("[3] Consultar bandeja de salida")
		fmt.Println("[4] Leer correo por ID (ej. MSG-001)")
		fmt.Println("[5] Eliminar correo por ID")
		fmt.Println("[6] Panel de Auditoría (Solo Supervisor)")
		fmt.Println("[0] Cerrar sesión")
		op := leerLinea(scanner, "> Seleccione opción: ")

		switch op {
		case "1":
			para := leerLinea(scanner, "Para: ")
			asunto := leerLinea(scanner, "Asunto: ")
			cuerpo := leerLinea(scanner, "Cuerpo: ")
			res := cli.EnviarSolicitud(dominio.SolicitudCorreo{
				Accion: "enviar",
				Correo: dominio.Correo{Destinatario: para, Asunto: asunto, Cuerpo: cuerpo},
			})
			fmt.Printf("[%d] %s\n", res.Codigo, res.Mensaje)
		case "2":
			res := cli.EnviarSolicitud(dominio.SolicitudCorreo{Accion: "entrada"})
			if !res.Exito {
				fmt.Printf("[ERROR %d] %s\n", res.Codigo, res.Mensaje)
				continue
			}
			for _, m := range res.Correos {
				est := "NO LEÍDO"
				if m.Leido {
					est = "LEÍDO"
				}
				fmt.Printf("- %s | [%s] | De: %s | Asunto: %s\n", m.ID, est, m.Remitente, m.Asunto)
			}
		case "3":
			res := cli.EnviarSolicitud(dominio.SolicitudCorreo{Accion: "salida"})
			if !res.Exito {
				fmt.Printf("[ERROR %d] %s\n", res.Codigo, res.Mensaje)
				continue
			}
			for _, m := range res.Correos {
				fmt.Printf("- %s | Para: %s | Asunto: %s\n", m.ID, m.Destinatario, m.Asunto)
			}
		case "4":
			id := leerLinea(scanner, "ID del correo: ")
			res := cli.EnviarSolicitud(dominio.SolicitudCorreo{Accion: "leer", IDCorreo: id})
			if res.Exito && res.Correo != nil {
				fmt.Printf("\nDe: %s | Para: %s\nAsunto: %s\nCuerpo: %s\n", res.Correo.Remitente, res.Correo.Destinatario, res.Correo.Asunto, res.Correo.Cuerpo)
			} else {
				fmt.Printf("[ERROR %d] %s\n", res.Codigo, res.Mensaje)
			}
		case "5":
			id := leerLinea(scanner, "ID del correo a eliminar: ")
			res := cli.EnviarSolicitud(dominio.SolicitudCorreo{Accion: "eliminar", IDCorreo: id})
			fmt.Printf("[%d] %s\n", res.Codigo, res.Mensaje)
		case "6":
			res := cli.EnviarSolicitud(dominio.SolicitudCorreo{Accion: "ver_auditoria"})
			fmt.Printf("[%d] %s\n", res.Codigo, res.Mensaje)
		case "0":
			cli.TokenSesion = ""
			fmt.Println("[OK] Sesión cerrada.")
		}
	}
}