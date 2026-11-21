package dominio

import "time"

type Usuario struct {
	ID  string `json:"id"`
	Rol string `json:"rol"`
}

type Correo struct {
	ID           string    `json:"id"`
	Remitente    string    `json:"remitente"`
	Destinatario string    `json:"destinatario"`
	Asunto       string    `json:"asunto"`
	Cuerpo       string    `json:"cuerpo"`
	FechaEnvio   time.Time `json:"fecha_envio"`
	Leido        bool      `json:"leido"`
	Eliminado    bool      `json:"eliminado"`
}

type SolicitudCorreo struct {
	Accion   string `json:"accion"` // "login", "enviar", "entrada", "salida", "leer", "eliminar", "panel_auditoria"
	Token    string `json:"token"`
	Usuario  string `json:"usuario,omitempty"`
	Password string `json:"password,omitempty"`
	Correo   Correo `json:"correo,omitempty"`
	IDCorreo string `json:"id_correo,omitempty"`
}

type RespuestaServidor struct {
	Exito   bool     `json:"exito"`
	Codigo  int      `json:"codigo"` // 200, 401, 403, 404, 503
	Mensaje string   `json:"mensaje"`
	Token   string   `json:"token,omitempty"`
	Rol     string   `json:"rol,omitempty"`
	Correos []Correo `json:"correos,omitempty"`
	Correo  *Correo  `json:"correo,omitempty"`
}