package seguridad

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/dominio"
)

type ValidadorAuth interface {
	Autenticar(usuario, clave string) (string, string, int, string)
	ValidarToken(token string) (dominio.Usuario, error)
	ExisteUsuario(usuario string) bool
}

type ClienteFastAPI struct {
	urlBase string
	cliente *http.Client
}

func NuevoClienteFastAPI(urlBase string) *ClienteFastAPI {
	return &ClienteFastAPI{
		urlBase: urlBase,
		cliente: &http.Client{Timeout: 3 * time.Second},
	}
}

func (c *ClienteFastAPI) Autenticar(usuario, clave string) (string, string, int, string) {
	body, _ := json.Marshal(map[string]string{"usuario": usuario, "password": clave})
	resp, err := c.cliente.Post(c.urlBase+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", "", 503, "Servicio FastAPI no disponible"
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if resp.StatusCode != 200 {
		detalle, _ := res["detail"].(string)
		return "", "", resp.StatusCode, detalle
	}
	return res["token"].(string), res["rol"].(string), 200, "Autenticación exitosa"
}

func (c *ClienteFastAPI) ValidarToken(token string) (dominio.Usuario, error) {
	if token == "" {
		return dominio.Usuario{}, errors.New("token vacío")
	}
	body, _ := json.Marshal(map[string]string{"token": token})
	resp, err := c.cliente.Post(c.urlBase+"/validar-token", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return dominio.Usuario{}, errors.New("503: servicio de autenticación no disponible")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return dominio.Usuario{}, errors.New("401: token inválido o expirado")
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	return dominio.Usuario{
		ID:  res["usuario"].(string),
		Rol: res["rol"].(string),
	}, nil
}

func (c *ClienteFastAPI) ExisteUsuario(usuario string) bool {
	resp, err := c.cliente.Get(fmt.Sprintf("%s/existe-usuario/%s", c.urlBase, usuario))
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	defer resp.Body.Close()
	var res map[string]bool
	json.NewDecoder(resp.Body).Decode(&res)
	return res["existe"]
}