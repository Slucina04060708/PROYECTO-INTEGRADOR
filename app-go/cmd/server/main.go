package main

import (
	"log"

	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/almacenamiento"
	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/negocio"
	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/red/servidor"
	"github.com/Slucina04060708/PROYECTO-INTEGRADOR/app-go/internal/seguridad"
)

func main() {
	repo := almacenamiento.NuevoRepositorioMemoria()
	auth := seguridad.NuevoClienteFastAPI("http://127.0.0.1:8000")
	servicio := negocio.NuevoServicioCorreo(repo, auth)
	srv := servidor.NuevoServidorCorreo(":8080", servicio)
	log.Fatal(srv.Iniciar())
}