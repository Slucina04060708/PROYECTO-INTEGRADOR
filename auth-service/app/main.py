import json, os, re, secrets, time
import bcrypt
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI(title="Servicio de Autenticacion")
USUARIOS_DB = "base_usuarios.json"

tokens_activos = {}
sesion_por_usuario = {}
intentos_fallidos = {}

class RegistroReq(BaseModel):
    usuario: str
    password: str
    rol: str  # "operativo", "administrativo", "supervisor"

class LoginReq(BaseModel):
    usuario: str
    password: str

class TokenReq(BaseModel):
    token: str

def cargar_usuarios():
    if not os.path.exists(USUARIOS_DB):
        return {}
    with open(USUARIOS_DB, "r") as f:
        return json.load(f)

def guardar_usuarios(data):
    with open(USUARIOS_DB, "w") as f:
        json.dump(data, f, indent=4)

def validar_password_robusta(pwd: str, usuario: str):
    if len(pwd) < 10:
        return False, "La contraseña debe tener al menos 10 caracteres."
    if not re.search(r"[A-Z]", pwd) or not re.search(r"[a-z]", pwd):
        return False, "Debe incluir mayúsculas y minúsculas."
    if not re.search(r"\d", pwd) or not re.search(r"[^A-Za-z0-9]", pwd):
        return False, "Debe incluir al menos un número y un símbolo."
    if usuario.lower() in pwd.lower() or "123456" in pwd:
        return False, "La contraseña no puede incluir el nombre del usuario ni secuencias de valores predecibles."
    return True, "OK"

@app.post("/register")
def registrar(req: RegistroReq):
    usuarios = cargar_usuarios()
    if req.usuario in usuarios:
        raise HTTPException(status_code=400, detail="El usuario ya existe.")
    
    valido, msg = validar_password_robusta(req.password, req.usuario)
    if not valido:
        raise HTTPException(status_code=400, detail=f"[C-01] Política rechazada: {msg}")

    salt = bcrypt.gensalt(rounds=12)
    hashed = bcrypt.hashpw(req.password.encode("utf-8"), salt).decode("utf-8")

    usuarios[req.usuario] = {
        "usuario": req.usuario,
        "password_hash": hashed,
        "rol": req.rol
    }
    guardar_usuarios(usuarios)
    return {"mensaje": f"Usuario {req.usuario} creado con éxito", "rol": req.rol}

@app.post("/login")
def login(req: LoginReq):
    ahora = time.time()
    estado_bloqueo = intentos_fallidos.get(req.usuario, {"conteo": 0, "bloqueado_hasta": 0})

    if estado_bloqueo["bloqueado_hasta"] > ahora:
        restante = int(estado_bloqueo["bloqueado_hasta"] - ahora)
        raise HTTPException(status_code=403, detail=f"[C-03] Cuenta bloqueada temporalmente. Reintente en {restante}s.")

    usuarios = cargar_usuarios()
    user_data = usuarios.get(req.usuario)

    if not user_data or not bcrypt.checkpw(req.password.encode("utf-8"), user_data["password_hash"].encode("utf-8")):
        estado_bloqueo["conteo"] += 1
        if estado_bloqueo["conteo"] >= 3:
            estado_bloqueo["bloqueado_hasta"] = ahora + 900
            intentos_fallidos[req.usuario] = estado_bloqueo
            raise HTTPException(status_code=403, detail="[C-03] Límite de 3 intentos alcanzado. Cuenta bloqueada 15 min.")
        intentos_fallidos[req.usuario] = estado_bloqueo
        raise HTTPException(status_code=401, detail=f"Credenciales inválidas (Intento {estado_bloqueo['conteo']}/3).")

    intentos_fallidos[req.usuario] = {"conteo": 0, "bloqueado_hasta": 0}

    if req.usuario in sesion_por_usuario:
        token_viejo = sesion_por_usuario[req.usuario]
        tokens_activos.pop(token_viejo, None)

    nuevo_token = secrets.token_hex(24)
    tokens_activos[nuevo_token] = {
        "usuario": req.usuario,
        "rol": user_data["rol"],
        "expira": ahora + 3600
    }
    sesion_por_usuario[req.usuario] = nuevo_token

    return {"token": nuevo_token, "usuario": req.usuario, "rol": user_data["rol"]}

@app.post("/validar-token")
def validar_token(req: TokenReq):
    sesion = tokens_activos.get(req.token)
    if not sesion:
        raise HTTPException(status_code=401, detail="Token inválido o revocado por nueva sesión.")
    if sesion["expira"] < time.time():
        tokens_activos.pop(req.token, None)
        raise HTTPException(status_code=401, detail="Token expirado.")
    return {"valido": True, "usuario": sesion["usuario"], "rol": sesion["rol"]}

@app.get("/existe-usuario/{usuario}")
def existe_usuario(usuario: str):
    usuarios = cargar_usuarios()
    return {"existe": usuario in usuarios}