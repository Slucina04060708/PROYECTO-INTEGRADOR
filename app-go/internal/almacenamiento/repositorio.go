package almacenamiento

import (
	"errors"
	"sync"

	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/dominio"
)

type RepositorioCorreo interface {
	Guardar(c dominio.Correo) error
	BuscarPorID(id string) (dominio.Correo, error)
	ListarBandejaEntrada(usuario string) []dominio.Correo
	ListarBandejaSalida(usuario string) []dominio.Correo
	ActualizarEstadoLeido(id string) error
	EliminarLogico(id string) error
}

type RepositorioMemoria struct {
	mu      sync.RWMutex
	correos map[string]dominio.Correo
}

func NuevoRepositorioMemoria() *RepositorioMemoria {
	return &RepositorioMemoria{
		correos: make(map[string]dominio.Correo),
	}
}

func (r *RepositorioMemoria) Guardar(c dominio.Correo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.correos[c.ID] = c
	return nil
}

func (r *RepositorioMemoria) BuscarPorID(id string) (dominio.Correo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.correos[id]
	if !ok || c.Eliminado {
		return dominio.Correo{}, errors.New("mensaje no encontrado")
	}
	return c, nil
}

func (r *RepositorioMemoria) ListarBandejaEntrada(usuario string) []dominio.Correo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var lista []dominio.Correo
	for _, c := range r.correos {
		if c.Destinatario == usuario && !c.Eliminado {
			lista = append(lista, c)
		}
	}
	return lista
}

func (r *RepositorioMemoria) ListarBandejaSalida(usuario string) []dominio.Correo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var lista []dominio.Correo
	for _, c := range r.correos {
		if c.Remitente == usuario && !c.Eliminado {
			lista = append(lista, c)
		}
	}
	return lista
}

func (r *RepositorioMemoria) ActualizarEstadoLeido(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.correos[id]
	if !ok || c.Eliminado {
		return errors.New("mensaje no encontrado")
	}
	c.Leido = true
	r.correos[id] = c
	return nil
}

func (r *RepositorioMemoria) EliminarLogico(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.correos[id]
	if !ok || c.Eliminado {
		return errors.New("mensaje no encontrado")
	}
	c.Eliminado = true
	r.correos[id] = c
	return nil
}