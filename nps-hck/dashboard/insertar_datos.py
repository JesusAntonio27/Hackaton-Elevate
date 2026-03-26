import json
import requests
import time

# --- CONFIGURACIÓN BASE ---
API_BASE_URL = "http://localhost:8080"

# --- TUS DATOS DE PRUEBA (10 casos variados con "area") ---
DATOS_PRUEBA = [
    # --- Octubre 2025 (10 Registros) ---
    {
        "fechaCaptura": "2025-10-05 09:15:00",
        "nps_score": 9,
        "comentario": "¡El mejor servicio que he recibido! Resolvieron todo en menos de 5 minutos.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2025-10-08 11:30:00",
        "nps_score": 2,
        "comentario": "La aplicación se bloquea constantemente, es inútil.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2025-10-12 14:00:00",
        "nps_score": 7,
        "comentario": "Funciona, pero la interfaz es muy confusa.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2025-10-15 16:45:00",
        "nps_score": 1,
        "comentario": "Es un acoso constante de llamadas de cobranza, incluso cuando ya pagué.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2025-10-18 10:00:00",
        "nps_score": 8,
        "comentario": "Me resolvieron la duda, pero tuve que explicarla tres veces.",
        "area": "Recursos humanos"
    },
    {
        "fechaCaptura": "2025-10-21 13:10:00",
        "nps_score": 10,
        "comentario": "Superó mis expectativas, el producto es de excelente calidad.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2025-10-23 15:50:00",
        "nps_score": 5,
        "comentario": "Llevo semanas esperando una solución y nadie me responde. Pésimo servicio.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2025-10-25 08:05:00",
        "nps_score": 10,
        "comentario": "La nueva actualización de la app es una maravilla, súper intuitiva y rápida.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2025-10-28 12:00:00",
        "nps_score": 3,
        "comentario": "Me hicieron un cargo doble y ahora tengo que perder mi tiempo para que lo corrijan.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2025-10-30 17:30:00",
        "nps_score": 7,
        "comentario": "El producto está bien, sin más. No me sorprende.",
        "area": "Recursos humanos"
    },
    # --- Noviembre 2025 (10 Registros) ---
    {
        "fechaCaptura": "2025-11-04 09:15:00",
        "nps_score": 4,
        "comentario": "La publicidad es engañosa, el producto final no tiene nada que ver.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2025-11-07 11:30:00",
        "nps_score": 9,
        "comentario": "El personal fue extremadamente amable y me ayudó en todo momento.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2025-11-11 14:00:00",
        "nps_score": 8,
        "comentario": "Es difícil encontrar la opción que necesito en la página web.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2025-11-14 16:45:00",
        "nps_score": 2,
        "comentario": "La aplicación se bloquea constantemente, es inútil.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2025-11-17 10:00:00",
        "nps_score": 7,
        "comentario": "El tiempo de espera fue aceptable, ni rápido ni lento.",
        "area": "Recursos humanos"
    },
    {
        "fechaCaptura": "2025-11-20 13:10:00",
        "nps_score": 6,
        "comentario": "Llevo semanas esperando una solución y nadie me responde. Pésimo servicio.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2025-11-22 15:50:00",
        "nps_score": 10,
        "comentario": "Recomendaría esta empresa a todos mis amigos y familiares sin dudarlo.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2025-11-24 08:05:00",
        "nps_score": 9,
        "comentario": "La nueva actualización de la app es una maravilla, súper intuitiva y rápida.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2025-11-27 12:00:00",
        "nps_score": 3,
        "comentario": "Es un acoso constante de llamadas de cobranza, incluso cuando ya pagué.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2025-11-29 17:30:00",
        "nps_score": 8,
        "comentario": "Me resolvieron la duda, pero tuve que explicarla tres veces.",
        "area": "Recursos humanos"
    },
    # --- Diciembre 2025 (10 Registros) ---
    {
        "fechaCaptura": "2025-12-03 09:15:00",
        "nps_score": 10,
        "comentario": "¡El mejor servicio que he recibido! Resolvieron todo en menos de 5 minutos.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2025-12-06 11:30:00",
        "nps_score": 1,
        "comentario": "Me hicieron un cargo doble y ahora tengo que perder mi tiempo para que lo corrijan.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2025-12-10 14:00:00",
        "nps_score": 7,
        "comentario": "Funciona, pero la interfaz es muy confusa.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2025-12-13 16:45:00",
        "nps_score": 8,
        "comentario": "El producto está bien, sin más. No me sorprende.",
        "area": "Recursos humanos"
    },
    {
        "fechaCaptura": "2025-12-16 10:00:00",
        "nps_score": 5,
        "comentario": "Llevo semanas esperando una solución y nadie me responde. Pésimo servicio.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2025-12-19 13:10:00",
        "nps_score": 9,
        "comentario": "Superó mis expectativas, el producto es de excelente calidad.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2025-12-21 15:50:00",
        "nps_score": 2,
        "comentario": "La aplicación se bloquea constantemente, es inútil.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2025-12-23 08:05:00",
        "nps_score": 10,
        "comentario": "Recomendaría esta empresa a todos mis amigos y familiares sin dudarlo.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2025-12-28 12:00:00",
        "nps_score": 4,
        "comentario": "La publicidad es engañosa, el producto final no tiene nada que ver.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2025-12-30 17:30:00",
        "nps_score": 7,
        "comentario": "El tiempo de espera fue aceptable, ni rápido ni lento.",
        "area": "Recursos humanos"
    },
    # --- Enero 2026 (10 Registros) ---
    {
        "fechaCaptura": "2026-01-04 09:15:00",
        "nps_score": 6,
        "comentario": "Llevo semanas esperando una solución y nadie me responde. Pésimo servicio.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2026-01-07 11:30:00",
        "nps_score": 8,
        "comentario": "Es difícil encontrar la opción que necesito en la página web.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2026-01-11 14:00:00",
        "nps_score": 9,
        "comentario": "¡El mejor servicio que he recibido! Resolvieron todo en menos de 5 minutos.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2026-01-14 16:45:00",
        "nps_score": 1,
        "comentario": "Es un acoso constante de llamadas de cobranza, incluso cuando ya pagué.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2026-01-17 10:00:00",
        "nps_score": 10,
        "comentario": "Superó mis expectativas, el producto es de excelente calidad.",
        "area": "Recursos humanos"
    },
    {
        "fechaCaptura": "2026-01-20 13:10:00",
        "nps_score": 3,
        "comentario": "Me hicieron un cargo doble y ahora tengo que perder mi tiempo para que lo corrijan.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2026-01-22 15:50:00",
        "nps_score": 7,
        "comentario": "El producto está bien, sin más. No me sorprende.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2026-01-24 08:05:00",
        "nps_score": 9,
        "comentario": "El personal fue extremadamente amable y me ayudó en todo momento.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2026-01-27 12:00:00",
        "nps_score": 2,
        "comentario": "La aplicación se bloquea constantemente, es inútil.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2026-01-30 17:30:00",
        "nps_score": 8,
        "comentario": "Me resolvieron la duda, pero tuve que explicarla tres veces.",
        "area": "Recursos humanos"
    },
    # --- Febrero 2026 (10 Registros) ---
    {
        "fechaCaptura": "2026-02-02 09:15:00",
        "nps_score": 10,
        "comentario": "La nueva actualización de la app es una maravilla, súper intuitiva y rápida.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2026-02-05 11:30:00",
        "nps_score": 7,
        "comentario": "Funciona, pero la interfaz es muy confusa.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2026-02-09 14:00:00",
        "nps_score": 4,
        "comentario": "La publicidad es engañosa, el producto final no tiene nada que ver.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2026-02-12 16:45:00",
        "nps_score": 9,
        "comentario": "¡El mejor servicio que he recibido! Resolvieron todo en menos de 5 minutos.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2026-02-15 10:00:00",
        "nps_score": 8,
        "comentario": "El tiempo de espera fue aceptable, ni rápido ni lento.",
        "area": "Recursos humanos"
    },
    {
        "fechaCaptura": "2026-02-18 13:10:00",
        "nps_score": 1,
        "comentario": "Me hicieron un cargo doble y ahora tengo que perder mi tiempo para que lo corrijan.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2026-02-20 15:50:00",
        "nps_score": 10,
        "comentario": "Recomendaría esta empresa a todos mis amigos y familiares sin dudarlo.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2026-02-23 08:05:00",
        "nps_score": 5,
        "comentario": "Llevo semanas esperando una solución y nadie me responde. Pésimo servicio.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2026-02-25 12:00:00",
        "nps_score": 8,
        "comentario": "El producto está bien, sin más. No me sorprende.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2026-02-27 17:30:00",
        "nps_score": 9,
        "comentario": "Superó mis expectativas, el producto es de excelente calidad.",
        "area": "Recursos humanos"
    },
    # --- Marzo 2026 (10 Registros) ---
    {
        "fechaCaptura": "2026-03-02 09:15:00",
        "nps_score": 2,
        "comentario": "La aplicación de Bancoppel se cierra cada vez que intento hacer una transferencia. Es frustrante porque me urge hacer pagos y no puedo.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2026-03-05 11:30:00",
        "nps_score": 9,
        "comentario": "El técnico que me atendió resolvió el problema de mi contraseña en menos de 5 minutos. Muy amable y directo al punto, excelente servicio.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2026-03-09 14:20:00",
        "nps_score": 5,
        "comentario": "Llevamos toda la mañana sin internet en la sucursal. Los clientes se están yendo porque no podemos cobrar. Urge que arreglen los servidores.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2026-03-12 16:45:00",
        "nps_score": 1,
        "comentario": "Me marcan 10 veces al día de cobranza por un atraso de dos días. ¡Ya les dije que mañana pago! Es un acoso insoportable.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2026-03-15 10:00:00",
        "nps_score": 8,
        "comentario": "El curso de inducción estuvo muy bien explicado y el personal de RH fue muy claro con las prestaciones, aunque la plataforma de videos estaba un poco lenta.",
        "area": "Recursos humanos"
    },
    {
        "fechaCaptura": "2026-03-18 13:10:00",
        "nps_score": 4,
        "comentario": "Fui a solicitar un préstamo pero las tasas de interés que me mostraron no coinciden con lo que vi en la publicidad. Me sentí engañado.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2026-03-20 15:50:00",
        "nps_score": 6,
        "comentario": "Pedí un monitor nuevo para mi estación de trabajo hace 3 semanas y sigo esperando. Nadie me da estatus.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2026-03-23 08:05:00",
        "nps_score": 10,
        "comentario": "El nuevo sistema de facturación es una maravilla. Todo carga súper rápido y la interfaz es súper intuitiva.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2026-03-25 12:00:00",
        "nps_score": 3,
        "comentario": "Claro, qué bonito que te ofrezcan descuentos, lástima que su sistema nunca registra mis abonos por la web.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2026-03-28 17:30:00",
        "nps_score": 7,
        "comentario": "Tuve dudas sobre mis vacaciones. Me contestaron rápido, pero me mandaron un manual de 50 hojas en lugar de explicarme.",
        "area": "Recursos humanos"
    }
]


# --- DEFINICIÓN DE CATÁLOGOS (SEEDING) CON DESCRIPCIONES ---
CATALOGOS_A_REGISTRAR = {
    "categorias": [
        {
            "nombre": "Experiencia Digital",
            "subcategorias": [
                {"nombre": "APP", "descripcion": "Fallos, usabilidad o problemas dentro de la aplicación móvil"},
                {"nombre": "WEB", "descripcion": "Problemas al navegar en la página de escritorio"},
                {"nombre": "UI/UX", "descripcion": "Diseño confuso o botones que no se entienden"},
                {"nombre": "Performance", "descripcion": "Lentitud, tiempos de carga altos o congelamientos"},
                {"nombre": "Bugs", "descripcion": "Errores técnicos inesperados o funcionalidades rotas"}
            ]
        },
        {
            "nombre": "Producto y Calidad",
            "subcategorias": [
                {"nombre": "Producto", "descripcion": "Opiniones sobre el artículo físico en sí"},
                {"nombre": "Servicio", "descripcion": "Calidad del servicio prestado"},
                {"nombre": "Calidad", "descripcion": "Durabilidad, materiales o defectos de fábrica"},
                {"nombre": "Precios", "descripcion": "Relación calidad-precio o costo del artículo"},
                {"nombre": "Disponibilidad", "descripcion": "Falta de stock o productos agotados"},
                {"nombre": "Surtido", "descripcion": "Poca variedad de tallas, colores o modelos"}
            ]
        },
        {
            "nombre": "Proceso de Compra y Entrega",
            "subcategorias": [
                {"nombre": "Entrega", "descripcion": "Tiempos de envío, demoras o entregas rápidas"},
                {"nombre": "Recepcion", "descripcion": "Estado del paquete al recibirlo (roto, dañado)"},
                {"nombre": "Comunicacion", "descripcion": "Avisos o notificaciones de envío al cliente"},
                {"nombre": "Rastreo", "descripcion": "Problemas con el tracking o guía de paquetería"},
                {"nombre": "Interaccion", "descripcion": "Trato o quejas específicas del repartidor"},
                {"nombre": "Costos", "descripcion": "Quejas sobre el costo del flete o envío"}
            ]
        },
        {
            "nombre": "Atención y Soporte al Cliente",
            "subcategorias": [
                {"nombre": "Interaccion", "descripcion": "Trato en general durante el contacto"},
                {"nombre": "Tiempos de espera", "descripcion": "Demora en contestar llamadas o chats"},
                {"nombre": "Efectividad", "descripcion": "Si el agente logró resolver o no el problema"},
                {"nombre": "Amabilidad", "descripcion": "Percepción de empatía y cortesía del personal"},
                {"nombre": "Claridad", "descripcion": "Información precisa, clara y sin contradicciones"},
                {"nombre": "Disponibilidad", "descripcion": "Facilidad para encontrar canales de ayuda"}
            ]
        },
        {
            "nombre": "Experiencia en Tienda Física",
            "subcategorias": [
                {"nombre": "Visita", "descripcion": "Experiencia general al visitar una sucursal"},
                {"nombre": "Atencion", "descripcion": "Trato del personal de piso o vendedores"},
                {"nombre": "Tiempos de espera", "descripcion": "Filas largas en cajas o servicios"},
                {"nombre": "Orden / Limpieza", "descripcion": "Estado visual, limpieza de pasillos o baños"},
                {"nombre": "Claridad", "descripcion": "Señalización correcta y precios bien exhibidos"}
            ]
        },
        {
            "nombre": "Servicios Financieros y Cobranza",
            "subcategorias": [
                {"nombre": "Pagos", "descripcion": "Dificultad o facilidad para realizar pagos"},
                {"nombre": "Credito", "descripcion": "Condiciones y proceso para solicitar crédito"},
                {"nombre": "Cobranza", "descripcion": "Acoso, llamadas o visitas de gestores de cobranza"},
                {"nombre": "Claridad", "descripcion": "Dudas sobre estados de cuenta o tasas de interés"},
                {"nombre": "Abonos", "descripcion": "Problemas con el registro de abonos realizados"},
                {"nombre": "Prestamos", "descripcion": "Inconvenientes al solicitar préstamos personales"}
            ]
        },
        {
            "nombre": "Marketing y Comunicaciones",
            "subcategorias": [
                {"nombre": "Cliente", "descripcion": "Percepción general de la marca"},
                {"nombre": "Promocion", "descripcion": "Promociones poco atractivas o mal aplicadas"},
                {"nombre": "Descuentos", "descripcion": "Quejas porque un descuento no pasó en caja/web"},
                {"nombre": "Publicidad", "descripcion": "Opiniones sobre anuncios en redes o TV"},
                {"nombre": "Engaños", "descripcion": "Publicidad engañosa o promesas falsas"},
                {"nombre": "Claridad", "descripcion": "Términos y condiciones confusos"},
                {"nombre": "Comunicaciones", "descripcion": "Exceso de correos (spam) o SMS"},
                {"nombre": "Politicas", "descripcion": "Problemas con políticas de devolución o cambios"}
            ]
        }
    ],
    "centros": [
        {"nombre": "Bancoppel"}, {"nombre": "Mesa de ayuda"}, {"nombre": "Operaciones IT"},
        {"nombre": "Cobranza Digital"}, {"nombre": "Recursos humanos"}
    ],
    "estatuses": [
        {"nombre": "Activo"}, {"nombre": "En atencion"}, {"nombre": "Atendido"}, {"nombre": "Cerrado"}
    ],
    "sentimientos": [
        {"nombre": "Enojado"}, {"nombre": "Feliz"}, {"nombre": "Neutral"}, {"nombre": "Triste"}, {"nombre": "Sarcastico"}
    ],
    "clasificacionesNPS": [
        {"nombre": "Detractores", "min": "1", "max": "6"},
        {"nombre": "Pasivos", "min": "7", "max": "8"},
        {"nombre": "Promotores", "min": "9", "max": "10"}
    ]
}

def limpiar_base():
    print("==================================================")
    print("🧹 FASE 0: Limpiando base de datos Mongo y Qdrant...")
    print("==================================================")
    try:
        response = requests.delete(f"{API_BASE_URL}/admin/limpiar")
        if response.status_code == 200:
            print("✅ ¡Base de datos reseteada a estado de fábrica!")
        else:
            print(f"⚠️ Error al limpiar: {response.text}")
    except Exception as e:
        print(f"❌ Error de conexión: {e}")
        exit()
    time.sleep(1)

def registrar_catalogos():
    print("\n==================================================")
    print("🌱 FASE 1: Sembrando Catálogos de Negocio...")
    print("==================================================")
    for endpoint, items in CATALOGOS_A_REGISTRAR.items():
        url = f"{API_BASE_URL}/{endpoint}"
        exitos = 0
        for item in items:
            try:
                if requests.post(url, json=item).status_code == 201: exitos += 1
            except Exception: pass
        print(f"✅ Registrados {exitos}/{len(items)} en '{endpoint}'.")
    time.sleep(1)

def procesar_nps():
    print("\n==================================================")
    print("🚀 FASE 2: Procesando datos con Gemini IA...")
    print("==================================================")
    url_procesar = f"{API_BASE_URL}/nps/procesar-lote"
    print(f"Enviando {len(DATOS_PRUEBA)} registros a Go...")

    try:
        response = requests.post(url_procesar, json=DATOS_PRUEBA)
        if response.status_code == 200:
            res_json = response.json()
            print(f"✅ ¡Éxito! Insertados: {res_json.get('exitosos', 0)} | Fallidos: {res_json.get('fallidos', 0)}")
        else:
            print(f"❌ Error del servidor: {response.text}")
    except Exception as e:
        print(f"❌ Error de red: {e}")

if __name__ == "__main__":
    limpiar_base()
    registrar_catalogos()
    procesar_nps()