package negocio

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/almacenamiento"
	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/dominio"
	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/seguridad"
)

type ServicioCorreo struct {
	repo     almacenamiento.RepositorioCorreo
	auth     seguridad.ValidadorAuth
	contador uint64
	logs     []string
}

func NuevoServicioCorreo(repo almacenamiento.RepositorioCorreo, auth seguridad.ValidadorAuth) *ServicioCorreo {
	return &ServicioCorreo{repo: repo, auth: auth}
}

func (s *ServicioCorreo) registrarAuditoria(evento string) {
	linea := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), evento)
	s.logs = append(s.logs, linea)
	log.Println(linea)
}

func (s *ServicioCorreo) Procesar(req dominio.SolicitudCorreo) dominio.RespuestaServidor {
	if req.Accion == "login" {
		token, rol, code, msg := s.auth.Autenticar(req.Usuario, req.Password)
		s.registrarAuditoria(fmt.Sprintf("LOGIN usuario=%s codigo=%d msg=%s", req.Usuario, code, msg))
		return dominio.RespuestaServidor{Exito: code == 200, Codigo: code, Mensaje: msg, Token: token, Rol: rol}
	}

	usr, err := s.auth.ValidarToken(req.Token)
	if err != nil {
		s.registrarAuditoria(fmt.Sprintf("ACCESO_DENEGADO accion=%s motivo=%s", req.Accion, err.Error()))
		return dominio.RespuestaServidor{Exito: false, Codigo: 401, Mensaje: err.Error()}
	}

	switch req.Accion {
	case "enviar":
		if usr.Rol == "supervisor" {
			s.registrarAuditoria(fmt.Sprintf("403_FORBIDDEN usuario=%s rol=%s intento=enviar", usr.ID, usr.Rol))
			return dominio.RespuestaServidor{Exito: false, Codigo: 403, Mensaje: "[C-04] Acceso denegado: el rol supervisor no puede enviar correos operativos."}
		}
		if req.Correo.Destinatario == "" || req.Correo.Asunto == "" {
			return dominio.RespuestaServidor{Exito: false, Codigo: 400, Mensaje: "Destinatario y asunto son obligatorios."}
		}
		if !s.auth.ExisteUsuario(req.Correo.Destinatario) {
			return dominio.RespuestaServidor{Exito: false, Codigo: 404, Mensaje: "Destinatario no existe."}
		}
		id := fmt.Sprintf("MSG-%03d", atomic.AddUint64(&s.contador, 1))
		req.Correo.ID = id
		req.Correo.Remitente = usr.ID
		req.Correo.FechaEnvio = time.Now()
		s.repo.Guardar(req.Correo)
		s.registrarAuditoria(fmt.Sprintf("CORREO_ENVIADO id=%s de=%s para=%s", id, usr.ID, req.Correo.Destinatario))
		return dominio.RespuestaServidor{Exito: true, Codigo: 200, Mensaje: "Correo enviado con ID " + id}

	case "entrada":
		if usr.Rol == "supervisor" {
			return dominio.RespuestaServidor{Exito: false, Codigo: 403, Mensaje: "[C-04] Acceso denegado: supervisor no accede a buzones personales."}
		}
		lista := s.repo.ListarBandejaEntrada(usr.ID)
		return dominio.RespuestaServidor{Exito: true, Codigo: 200, Mensaje: "Bandeja de entrada", Correos: lista}

	case "salida":
		if usr.Rol == "supervisor" {
			return dominio.RespuestaServidor{Exito: false, Codigo: 403, Mensaje: "[C-04] Acceso denegado: supervisor no accede a buzones personales."}
		}
		lista := s.repo.ListarBandejaSalida(usr.ID)
		return dominio.RespuestaServidor{Exito: true, Codigo: 200, Mensaje: "Bandeja de salida", Correos: lista}

	case "leer":
		correo, err := s.repo.BuscarPorID(req.IDCorreo)
		if err != nil {
			return dominio.RespuestaServidor{Exito: false, Codigo: 404, Mensaje: "Mensaje no encontrado."}
		}
		if correo.Remitente != usr.ID && correo.Destinatario != usr.ID {
			s.registrarAuditoria(fmt.Sprintf("403_FORBIDDEN usuario=%s intento leer correo ajeno %s", usr.ID, req.IDCorreo))
			return dominio.RespuestaServidor{Exito: false, Codigo: 403, Mensaje: "[C-04] Acceso denegado (403 Forbidden): este correo pertenece a otro colaborador."}
		}
		if correo.Destinatario == usr.ID {
			s.repo.ActualizarEstadoLeido(req.IDCorreo)
			correo.Leido = true
		}
		return dominio.RespuestaServidor{Exito: true, Codigo: 200, Mensaje: "Mensaje recuperado", Correo: &correo}

	case "eliminar":
		correo, err := s.repo.BuscarPorID(req.IDCorreo)
		if err != nil {
			return dominio.RespuestaServidor{Exito: false, Codigo: 404, Mensaje: "Mensaje no encontrado."}
		}
		if correo.Remitente != usr.ID && correo.Destinatario != usr.ID {
			s.registrarAuditoria(fmt.Sprintf("403_FORBIDDEN usuario=%s intento eliminar correo ajeno %s", usr.ID, req.IDCorreo))
			return dominio.RespuestaServidor{Exito: false, Codigo: 403, Mensaje: "[C-04] Acceso denegado (403 Forbidden): no puede eliminar correos ajenos."}
		}
		s.repo.EliminarLogico(req.IDCorreo)
		return dominio.RespuestaServidor{Exito: true, Codigo: 200, Mensaje: "Mensaje eliminado lógicamente"}

	case "ver_auditoria":
		if usr.Rol != "supervisor" {
			s.registrarAuditoria(fmt.Sprintf("403_FORBIDDEN usuario=%s rol=%s intento=ver_auditoria", usr.ID, usr.Rol))
			return dominio.RespuestaServidor{Exito: false, Codigo: 403, Mensaje: "[C-04] Acceso denegado (403 Forbidden): módulo exclusivo para Supervisor Técnico."}
		}
		return dominio.RespuestaServidor{Exito: true, Codigo: 200, Mensaje: fmt.Sprintf("Registros de auditoría (%d eventos)", len(s.logs))}
	}

	return dominio.RespuestaServidor{Exito: false, Codigo: 400, Mensaje: "Acción no reconocida"}
}