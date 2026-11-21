package cliente

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/dominio"
)

type ClienteCorreo struct {
	URLServidor     string
	TokenSesion     string
	UsuarioActivo   string
	RolActivo       string
	UltimaActividad time.Time
}

func NuevoClienteCorreo(url string) *ClienteCorreo {
	return &ClienteCorreo{URLServidor: url}
}

func (c *ClienteCorreo) EnviarSolicitud(req dominio.SolicitudCorreo) dominio.RespuestaServidor {
	if c.TokenSesion != "" && time.Since(c.UltimaActividad) > 10*time.Minute {
		c.TokenSesion = ""
		c.UsuarioActivo = ""
		return dominio.RespuestaServidor{Exito: false, Codigo: 401, Mensaje: "[C-02] Sesión cerrada por 10 minutos de inactividad."}
	}
	c.UltimaActividad = time.Now()
	req.Token = c.TokenSesion

	body, _ := json.Marshal(req)
	resp, err := http.Post(c.URLServidor+"/api", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return dominio.RespuestaServidor{Exito: false, Codigo: 503, Mensaje: "No se pudo conectar al servidor Go"}
	}
	defer resp.Body.Close()

	var res dominio.RespuestaServidor
	json.NewDecoder(resp.Body).Decode(&res)
	if req.Accion == "login" && res.Exito {
		c.TokenSesion = res.Token
		c.UsuarioActivo = req.Usuario
		c.RolActivo = res.Rol
	}
	return res
}