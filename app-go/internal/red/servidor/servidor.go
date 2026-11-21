package servidor

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/dominio"
	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/negocio"
)

type ServidorCorreo struct {
	direccion string
	servicio  *negocio.ServicioCorreo
}

func NuevoServidorCorreo(direccion string, servicio *negocio.ServicioCorreo) *ServidorCorreo {
	return &ServidorCorreo{direccion: direccion, servicio: servicio}
}

func (s *ServidorCorreo) Iniciar() error {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		var req dominio.SolicitudCorreo
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(dominio.RespuestaServidor{Exito: false, Codigo: 400, Mensaje: "JSON inválido"})
			return
		}
		resp := s.servicio.Procesar(req)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Codigo)
		json.NewEncoder(w).Encode(resp)
	})
	log.Println("Servidor Go listo y escuchando peticiones en", s.direccion)
	return http.ListenAndServe(s.direccion, nil)
}