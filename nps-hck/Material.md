---
centros -> [
    {
        "_id":"1",
        "nombre": "Bancoppel"
    },
    {
        "_id": "2",
        "nombre": "Mesa de ayuda"
    }
]

estatuses -> [
    {
        "_id" : "1",
        "nombre" : "Activo"
    },
    {
        "_id" : "2",
        "nombre" : "En atencion"
    },
    {
        "_id" : "3",
        "nombre" : "Atendido"
    },
    {
        "_id" : "4",
        "nombre" : "Cerrado"
    },
]

sentimientos -> [
    {
        "_id" : "1",
        "nombre": "Enojado"
    },
    {
        "_id" : "2",
        "nombre": "Feliz"
    },
    {
        "_id" : "3",
        "nombre": "Sarcastico"
    },
    {
        "_id" : "4",
        "nombre": "Neutral"
    },
    {
        "_id" : "5",
        "nombre": "Feliz"
    },
]

clasificacionesNPS -> [
    {
        "_id" : "1",
        "nombre" : "Detractores",
        "min": "1",
        "max": "6"
    },
    {
        "_id" : "2",
        "nombre" : "Pasivos",
        "min": "7",
        "max": "8"
    },
    {
        "_id" : "3",
        "nombre" : "Promotores",
        "min": "9",
        "max": "10"
    },
]


reseñas ->
[
    {
        "_id": "1",
        "centro": {
            "_id": "2",
            "nombre": "Mesa de ayuda"
        },
        "categoria": {
            "_id": "4",
            "nombre": "Atención y Soporte al Cliente",
        },
        "historia": "Llame a atencion a clientes en la app de bancoppel por problemas en mi cuenta y no me atendieron, estoy muy enojado",
        "problema_pricipal":  "Falta de atencion",
        "problema_secundario": "Problemas app bancoppel",
        "fecha_captura": "25/03/2026 7:30pm",
        "fecha_cierre": "25/03/2026 7:50pm",
        "nsp_score": "2",
        "estatus": "1",
        "nivel_1" : {
            "categoria" : {
                "_id": "1",
                "nombre": "Atención y Soporte al Cliente"
            },
            "sub_categoria" : {
                "_id": "2",
                "nombre": "Tiempos de espera"
            }
        },  
        "nivel_2": {
            "categoria": {
                "_id": "1",
                "nombre": "Experiencia Digital",
            },

            "sub_categoria" : {
                "_id": "2",
                "nombre": "APP"
            }
        },
        "sentimientoIA": {
            "id" : "1",
            "name" : "Enojado"
        },
        "problemaIA" : "Usuario tenia problemas en cuenta en aplicacion bancoppel",
        "npsIA" : "1",
        "razonIA" : "No atendieron al usuario",
        "clasificacionNPS": {
            "_id" : "1",
            "nombre" : "Detractores",
            "min": "1",
            "max": "6"
        }
        "iniciativa" : "",
        "iniciativaIA" : "Mejorar el area de atencion y soporte al cliente",
        "centroSoporte" : {
            "id" : "1",
            "nombre": "Bancoppel"
        },
        "respuesta" : "",
        "respuestaIA" : "Lamentamos lo sucedido trabajaremos para mejorar nuestra atencion y soporte" 
    }
]
---
baseConocimientoDB> db.feedback_nps_2026.find({}).limit(2)
[
  {
    _id: ObjectId('69c4b549f39e72ebe5ef28d8'),
    collection: 'feedback_nps_2026',
    centro: { _id: ObjectId('69c4b53af39e72ebe5ef28c7'), nombre: 'Bancoppel' },
    categoria: {
      _id: ObjectId('69c4b53af39e72ebe5ef28c0'),
      nombre: 'Experiencia Digital',
      subcategorias: [
        {
          nombre: 'APP',
          descripcion: 'Fallos, usabilidad o problemas dentro de la aplicación móvil'
        },
        {
          nombre: 'WEB',
          descripcion: 'Problemas al navegar en la página de escritorio'
        },
        {
          nombre: 'UI/UX',
          descripcion: 'Diseño confuso o botones que no se entienden'
        },
        {
          nombre: 'Performance',
          descripcion: 'Lentitud, tiempos de carga altos o congelamientos'
        },
        {
          nombre: 'Bugs',
          descripcion: 'Errores técnicos inesperados o funcionalidades rotas'
        }
      ]
    },
    historia: 'La aplicación de Bancoppel se cierra cada vez que intento hacer una transferencia. Es frustrante porque me urge hacer pagos y no puedo.',
    problema_pricipal: 'Cierre inesperado de la aplicación al transferir',
    problema_secundario: 'Imposibilidad de realizar pagos urgentes',
    fecha_captura: '2026-03-24 09:15:00',
    fecha_cierre: '',
    nsp_score: 2,
    estatus: { _id: ObjectId('69c4b53af39e72ebe5ef28cc'), nombre: 'Activo' },
    nivel_1: {
      categoria: {
        _id: ObjectId('69c4b53af39e72ebe5ef28c0'),
        nombre: 'Experiencia Digital',
        subcategorias: [
          {
            nombre: 'APP',
            descripcion: 'Fallos, usabilidad o problemas dentro de la aplicación móvil'
          },
          {
            nombre: 'WEB',
            descripcion: 'Problemas al navegar en la página de escritorio'
          },
          {
            nombre: 'UI/UX',
            descripcion: 'Diseño confuso o botones que no se entienden'
          },
          {
            nombre: 'Performance',
            descripcion: 'Lentitud, tiempos de carga altos o congelamientos'
          },
          {
            nombre: 'Bugs',
            descripcion: 'Errores técnicos inesperados o funcionalidades rotas'
          }
        ]
      },
      sub_categoria: { nombre: 'Bugs' }
    },
    nivel_2: {
      categoria: {
        _id: ObjectId('69c4b53af39e72ebe5ef28c0'),
        nombre: 'Experiencia Digital',
        subcategorias: [
          {
            nombre: 'APP',
            descripcion: 'Fallos, usabilidad o problemas dentro de la aplicación móvil'
          },
          {
            nombre: 'WEB',
            descripcion: 'Problemas al navegar en la página de escritorio'
          },
          {
            nombre: 'UI/UX',
            descripcion: 'Diseño confuso o botones que no se entienden'
          },
          {
            nombre: 'Performance',
            descripcion: 'Lentitud, tiempos de carga altos o congelamientos'
          },
          {
            nombre: 'Bugs',
            descripcion: 'Errores técnicos inesperados o funcionalidades rotas'
          }
        ]
      },
      sub_categoria: { nombre: 'Performance' }
    },
    sentimientoIA: { _id: ObjectId('69c4b53af39e72ebe5ef28d0'), nombre: 'Enojado' },
    problemaIA: 'Cierre inesperado de la aplicación al transferir',
    npsIA: '2',
    razonNpsIA: 'El usuario experimenta una falla crítica (bug) que impide la función principal de la app.',
    razonIA: 'Inestabilidad técnica en la aplicación móvil durante transacciones.',
    clasificacionNPS: {
      _id: ObjectId('69c4b53af39e72ebe5ef28d5'),
      nombre: 'Detractores',
      min: '1',
      max: '6'
    },
    iniciativa: '',
    iniciativaIA: 'Depuración de errores en el módulo de transferencias y optimización de estabilidad.',
    centroSoporte: { _id: ObjectId('69c4b53af39e72ebe5ef28c9'), nombre: 'Operaciones IT' },
    respuesta: '',
    respuestaIA: 'Lamentamos los inconvenientes con la app. Estamos trabajando en una actualización para corregir los cierres inesperados.'
  },
  {
    _id: ObjectId('69c4b549f39e72ebe5ef28d9'),
    collection: 'feedback_nps_2026',
    centro: { _id: ObjectId('69c4b53af39e72ebe5ef28c8'), nombre: 'Mesa de ayuda' },
    categoria: {
      _id: ObjectId('69c4b53af39e72ebe5ef28c3'),
      nombre: 'Atención y Soporte al Cliente',
      subcategorias: [
        {
          nombre: 'Interaccion',
          descripcion: 'Trato en general durante el contacto'
        },
        {
          nombre: 'Tiempos de espera',
          descripcion: 'Demora en contestar llamadas o chats'
        },
        {
          nombre: 'Efectividad',
          descripcion: 'Si el agente logró resolver o no el problema'
        },
        {
          nombre: 'Amabilidad',
          descripcion: 'Percepción de empatía y cortesía del personal'
        },
        {
          nombre: 'Claridad',
          descripcion: 'Información precisa, clara y sin contradicciones'
        },
        {
          nombre: 'Disponibilidad',
          descripcion: 'Facilidad para encontrar canales de ayuda'
        }
      ]
    },
    historia: 'El técnico que me atendió resolvió el problema de mi contraseña en menos de 5 minutos. Muy amable y directo al punto, excelente servicio.',
    problema_pricipal: 'Resolución rápida de acceso',
    problema_secundario: '',
    fecha_captura: '2026-03-25 11:30:00',
    fecha_cierre: '',
    nsp_score: 9,
    estatus: { _id: ObjectId('69c4b53af39e72ebe5ef28cc'), nombre: 'Activo' },
    nivel_1: {
      categoria: {
        _id: ObjectId('69c4b53af39e72ebe5ef28c3'),
        nombre: 'Atención y Soporte al Cliente',
        subcategorias: [
          {
            nombre: 'Interaccion',
            descripcion: 'Trato en general durante el contacto'
          },
          {
            nombre: 'Tiempos de espera',
            descripcion: 'Demora en contestar llamadas o chats'
          },
          {
            nombre: 'Efectividad',
            descripcion: 'Si el agente logró resolver o no el problema'
          },
          {
            nombre: 'Amabilidad',
            descripcion: 'Percepción de empatía y cortesía del personal'
          },
          {
            nombre: 'Claridad',
            descripcion: 'Información precisa, clara y sin contradicciones'
          },
          {
            nombre: 'Disponibilidad',
            descripcion: 'Facilidad para encontrar canales de ayuda'
          }
        ]
      },
      sub_categoria: { nombre: 'Efectividad' }
    },
    nivel_2: { categoria: { nombre: '' }, sub_categoria: { nombre: '' } },
    sentimientoIA: { _id: ObjectId('69c4b53af39e72ebe5ef28d1'), nombre: 'Feliz' },
    problemaIA: 'Resolución rápida de acceso',
    npsIA: '10',
    razonNpsIA: 'Atención eficiente, rápida y con trato amable.',
    razonIA: 'Excelente soporte técnico y resolución inmediata.',
    clasificacionNPS: {
      _id: ObjectId('69c4b53af39e72ebe5ef28d7'),
      nombre: 'Promotores',
      min: '9',
      max: '10'
    },
    iniciativa: '',
    iniciativaIA: 'Mantener los estándares de capacitación y tiempos de respuesta actuales.',
    centroSoporte: { nombre: '' },
    respuesta: '',
    respuestaIA: '¡Muchas gracias por tu comentario! Nos alegra saber que el equipo de soporte pudo ayudarte rápidamente.'
  }
]
---
si fue la entrada DATOS_PRUEBA = [
    {
        "fechaCaptura": "2026-03-24 09:15:00",
        "nps_score": 2,
        "comentario": "La aplicación de Bancoppel se cierra cada vez que intento hacer una transferencia. Es frustrante porque me urge hacer pagos y no puedo.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2026-03-25 11:30:00",
        "nps_score": 9,
        "comentario": "El técnico que me atendió resolvió el problema de mi contraseña en menos de 5 minutos. Muy amable y directo al punto, excelente servicio.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2026-03-25 14:20:00",
        "nps_score": 5,
        "comentario": "Llevamos toda la mañana sin internet en la sucursal. Los clientes se están yendo porque no podemos cobrar. Urge que arreglen los servidores.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2026-03-20 16:45:00",
        "nps_score": 1,
        "comentario": "Me marcan 10 veces al día de cobranza por un atraso de dos días. ¡Ya les dije que mañana pago! Es un acoso insoportable.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2026-03-21 10:00:00",
        "nps_score": 8,
        "comentario": "El curso de inducción estuvo muy bien explicado y el personal de RH fue muy claro con las prestaciones, aunque la plataforma de videos estaba un poco lenta.",
        "area": "Recursos humanos"
    },
    {
        "fechaCaptura": "2026-03-22 13:10:00",
        "nps_score": 4,
        "comentario": "Fui a solicitar un préstamo pero las tasas de interés que me mostraron en el contrato no coinciden con lo que vi en la publicidad de Facebook. Me sentí engañado.",
        "area": "Bancoppel"
    },
    {
        "fechaCaptura": "2026-03-23 15:50:00",
        "nps_score": 6,
        "comentario": "Pedí un monitor nuevo para mi estación de trabajo hace 3 semanas y sigo esperando. Nadie me da estatus de la entrega ni fecha estimada.",
        "area": "Mesa de ayuda"
    },
    {
        "fechaCaptura": "2026-03-24 08:05:00",
        "nps_score": 10,
        "comentario": "El nuevo sistema de facturación es una maravilla. Todo carga súper rápido y la interfaz es súper intuitiva, nos ahorró horas de trabajo.",
        "area": "Operaciones IT"
    },
    {
        "fechaCaptura": "2026-03-25 12:00:00",
        "nps_score": 3,
        "comentario": "Claro, qué bonito que te ofrezcan descuentos por pagar a tiempo, lástima que su sistema nunca registra mis abonos cuando los hago por la web.",
        "area": "Cobranza Digital"
    },
    {
        "fechaCaptura": "2026-03-25 17:30:00",
        "nps_score": 7,
        "comentario": "Tuve dudas sobre cómo registrar mis vacaciones en el portal. Me contestaron rápido, pero me mandaron un manual de 50 hojas en lugar de explicarme paso a paso.",
        "area": "Recursos humanos"
    }
]
---
baseConocimientoDB> show collections
categorias
centros
clasificacionesNPS
estatuses
feedback_nps_2026
sentimientos
baseConocimientoDB> db.categorias.find({})
[
  {
    _id: ObjectId('69c4d80b7b57be30197bb0da'),
    nombre: 'Experiencia Digital',
    subcategorias: [
      {
        nombre: 'APP',
        descripcion: 'Fallos, usabilidad o problemas dentro de la aplicación móvil'
      },
      {
        nombre: 'WEB',
        descripcion: 'Problemas al navegar en la página de escritorio'
      },
      {
        nombre: 'UI/UX',
        descripcion: 'Diseño confuso o botones que no se entienden'
      },
      {
        descripcion: 'Lentitud, tiempos de carga altos o congelamientos',
        nombre: 'Performance'
      },
      {
        nombre: 'Bugs',
        descripcion: 'Errores técnicos inesperados o funcionalidades rotas'
      }
    ]
  },
  {
    _id: ObjectId('69c4d80b7b57be30197bb0db'),
    nombre: 'Producto y Calidad',
    subcategorias: [
      {
        nombre: 'Producto',
        descripcion: 'Opiniones sobre el artículo físico en sí'
      },
      { nombre: 'Servicio', descripcion: 'Calidad del servicio prestado' },
      {
        nombre: 'Calidad',
        descripcion: 'Durabilidad, materiales o defectos de fábrica'
      },
      {
        nombre: 'Precios',
        descripcion: 'Relación calidad-precio o costo del artículo'
      },
      {
        nombre: 'Disponibilidad',
        descripcion: 'Falta de stock o productos agotados'
      },
      {
        nombre: 'Surtido',
        descripcion: 'Poca variedad de tallas, colores o modelos'
      }
    ]
  },
  {
    _id: ObjectId('69c4d80b7b57be30197bb0dc'),
    nombre: 'Proceso de Compra y Entrega',
    subcategorias: [
      {
        nombre: 'Entrega',
        descripcion: 'Tiempos de envío, demoras o entregas rápidas'
      },
      {
        nombre: 'Recepcion',
        descripcion: 'Estado del paquete al recibirlo (roto, dañado)'
      },
      {
        descripcion: 'Avisos o notificaciones de envío al cliente',
        nombre: 'Comunicacion'
      },
      {
        nombre: 'Rastreo',
        descripcion: 'Problemas con el tracking o guía de paquetería'
      },
      {
        nombre: 'Interaccion',
        descripcion: 'Trato o quejas específicas del repartidor'
      },
      {
        nombre: 'Costos',
        descripcion: 'Quejas sobre el costo del flete o envío'
      }
    ]
  },
  {
    _id: ObjectId('69c4d80b7b57be30197bb0dd'),
    nombre: 'Atención y Soporte al Cliente',
    subcategorias: [
      {
        nombre: 'Interaccion',
        descripcion: 'Trato en general durante el contacto'
      },
      {
        nombre: 'Tiempos de espera',
        descripcion: 'Demora en contestar llamadas o chats'
      },
      {
        nombre: 'Efectividad',
        descripcion: 'Si el agente logró resolver o no el problema'
      },
      {
        nombre: 'Amabilidad',
        descripcion: 'Percepción de empatía y cortesía del personal'
      },
      {
        nombre: 'Claridad',
        descripcion: 'Información precisa, clara y sin contradicciones'
      },
      {
        nombre: 'Disponibilidad',
        descripcion: 'Facilidad para encontrar canales de ayuda'
      }
    ]
  },
  {
    _id: ObjectId('69c4d80b7b57be30197bb0de'),
    nombre: 'Experiencia en Tienda Física',
    subcategorias: [
      {
        nombre: 'Visita',
        descripcion: 'Experiencia general al visitar una sucursal'
      },
      {
        nombre: 'Atencion',
        descripcion: 'Trato del personal de piso o vendedores'
      },
      {
        nombre: 'Tiempos de espera',
        descripcion: 'Filas largas en cajas o servicios'
      },
      {
        nombre: 'Orden / Limpieza',
        descripcion: 'Estado visual, limpieza de pasillos o baños'
      },
      {
        nombre: 'Claridad',
        descripcion: 'Señalización correcta y precios bien exhibidos'
      }
    ]
  },
  {
    _id: ObjectId('69c4d80b7b57be30197bb0df'),
    nombre: 'Servicios Financieros y Cobranza',
    subcategorias: [
      {
        nombre: 'Pagos',
        descripcion: 'Dificultad o facilidad para realizar pagos'
      },
      {
        nombre: 'Credito',
        descripcion: 'Condiciones y proceso para solicitar crédito'
      },
      {
        nombre: 'Cobranza',
        descripcion: 'Acoso, llamadas o visitas de gestores de cobranza'
      },
      {
        nombre: 'Claridad',
        descripcion: 'Dudas sobre estados de cuenta o tasas de interés'
      },
      {
        nombre: 'Abonos',
        descripcion: 'Problemas con el registro de abonos realizados'
      },
      {
        descripcion: 'Inconvenientes al solicitar préstamos personales',
        nombre: 'Prestamos'
      }
    ]
  },
  {
    _id: ObjectId('69c4d80b7b57be30197bb0e0'),
    nombre: 'Marketing y Comunicaciones',
    subcategorias: [
      { nombre: 'Cliente', descripcion: 'Percepción general de la marca' },
      {
        nombre: 'Promocion',
        descripcion: 'Promociones poco atractivas o mal aplicadas'
      },
      {
        nombre: 'Descuentos',
        descripcion: 'Quejas porque un descuento no pasó en caja/web'
      },
      {
        nombre: 'Publicidad',
        descripcion: 'Opiniones sobre anuncios en redes o TV'
      },
      {
        nombre: 'Engaños',
        descripcion: 'Publicidad engañosa o promesas falsas'
      },
      { descripcion: 'Términos y condiciones confusos', nombre: 'Claridad' },
      {
        nombre: 'Comunicaciones',
        descripcion: 'Exceso de correos (spam) o SMS'
      },
      {
        nombre: 'Politicas',
        descripcion: 'Problemas con políticas de devolución o cambios'
      }
    ]
  }
]
baseConocimientoDB> db.centros})
[
  { _id: ObjectId('69c4d80b7b57be30197bb0e1'), nombre: 'Bancoppel' },
  { _id: ObjectId('69c4d80b7b57be30197bb0e2'), nombre: 'Mesa de ayuda' },
  { _id: ObjectId('69c4d80b7b57be30197bb0e3'), nombre: 'Operaciones IT' },
  {
    _id: ObjectId('69c4d80b7b57be30197bb0e4'),
    nombre: 'Cobranza Digital'
  },
  {
    _id: ObjectId('69c4d80b7b57be30197bb0e5'),
    nombre: 'Recursos humanos'
  }
]
baseConocimientoDB> db.estatuses
[
  { _id: ObjectId('69c4d80b7b57be30197bb0e6'), nombre: 'Activo' },
  { _id: ObjectId('69c4d80b7b57be30197bb0e7'), nombre: 'En atencion' },
  { _id: ObjectId('69c4d80b7b57be30197bb0e8'), nombre: 'Atendido' },
  { _id: ObjectId('69c4d80b7b57be30197bb0e9'), nombre: 'Cerrado' }
]
baseConocimientoDB> db.sentimientos
[
  { _id: ObjectId('69c4d80b7b57be30197bb0ea'), nombre: 'Enojado' },
  { _id: ObjectId('69c4d80b7b57be30197bb0eb'), nombre: 'Feliz' },
  { _id: ObjectId('69c4d80b7b57be30197bb0ec'), nombre: 'Neutral' },
  { _id: ObjectId('69c4d80b7b57be30197bb0ed'), nombre: 'Triste' },
  { _id: ObjectId('69c4d80b7b57be30197bb0ee'), nombre: 'Sarcastico' }
]
baseConocimient
---
graph TD
    subgraph "Cliente (Usuario/Frontend)"
        A[Usuario/App]
    end

    subgraph "API Gateway (Go con Mux en localhost:8080)"
        B(Router: mux)
        B -- /nps/procesar-lote --> C{Handler: ProcessAndInsertBatch}
        B -- /nps/buscar --> D{Handler: SearchResenas}
        B -- /nps/stats --> E{Handler: GetNPSStats}
        B -- /nps/analisis/pareto --> F{Handler: GetParetoAnalysis}
        B -- /nps/recomendar --> G{Handler: RecommendStrategy}
        B -- /nps/prediccion-estrategica --> H{Handler: PrediccionEstrategica}
        B -- /catalogo/* --> I{Handlers CRUD}
        B -- /nps/resenas --> J{Handler: GetResenas}
    end

    subgraph "Lógica de IA (Servicio Externo)"
        K[LLM: Gemini Flash]
    end

    subgraph "Almacenamiento de Datos"
        L[(MongoDB)]
        M[(Qdrant: Vector DB)]
    end

    %% --- Flujos de Datos ---

    %% Flujo 1: Ingesta y Enriquecimiento (El más importante)
    A -- "1. POST con reseñas simples" --> C
    C -- "2. Llama a la IA para enriquecer datos" --> K
    K -- "3. Devuelve datos estructurados (sentimiento, categoría, etc.)" --> C
    C -- "4. Inserta reseña COMPLETA" --> L
    C -- "5. Genera vector y lo inserta con ID de Mongo" --> M

    %% Flujo 2: Búsqueda Semántica
    A -- "1. POST con texto de búsqueda" --> D
    D -- "2. Convierte texto a vector y busca en Qdrant" --> M
    M -- "3. Devuelve IDs de Mongo" --> D
    D -- "4. Usa IDs para obtener reseñas completas" --> L
    L -- "5. Devuelve documentos completos" --> D
    D -- "6. Responde al usuario" --> A

    %% Flujo 3: Análisis y Reportes
    A -- "GET/POST para análisis" --> E & F & H & J
    E & F & H & J -- "Consultas directas para métricas y filtros" --> L

    %% Flujo 4: IA para Estrategia
    G & H -- "Llama a la IA con evidencia de la BD" --> K
    K -- "Devuelve Análisis Predictivo / Recomendación" --> G & H
    G & H -- "Responde al usuario" --> A

    %% Flujo 5: CRUD de Catálogos
    A -- "GET/POST/PUT/DELETE" --> I
    I -- "Lee/Escribe en colecciones de catálogo" --> L


    %% Estilos
    classDef api fill:#e6f3ff,stroke:#007bff,stroke-width:2px
    classDef ia fill:#fff0e6,stroke:#ff8c00,stroke-width:2px
    classDef db fill:#e6ffed,stroke:#28a745,stroke-width:2px
    classDef client fill:#f0f0f0,stroke:#666,stroke-width:2px
    class B,C,D,E,F,G,H,I,J api
    class K ia
    class L,M db
    class A client

    ---
    mini descripcion del proyecto
    URL Base: http://localhost:8080

Sección 1: Endpoints de Ingesta y Gestión de Reseñas
Endpoints dedicados a la creación, búsqueda y gestión de las reseñas de clientes.

Procesar un Lote de Reseñas (Método Principal)
Endpoint: POST /nps/procesar-lote

Descripción: Es el método recomendado para añadir nuevas reseñas. Recibe un arreglo de objetos simples, los envía a un LLM para ser analizados y enriquecidos (determinando categoría, sentimiento, problema, etc.), y luego los inserta en la base de datos como documentos completos.

Cuerpo (Body) de la Petición: Un arreglo de objetos InputSimple.

Ejemplo curl:

bash
curl -X POST http://localhost:8080/nps/procesar-lote \
-H "Content-Type: application/json" \
-d '[
  {
    "fechaCaptura": "2026-03-27T10:00:00Z",
    "nps_score": 1,
    "comentario": "El servicio de entrega fue pésimo, mi paquete llegó roto.",
    "area": "Logística"
  },
  {
    "fechaCaptura": "2026-03-27T11:00:00Z",
    "nps_score": 10,
    "comentario": "La vendedora de la tienda fue increíblemente amable y eficiente.",
    "area": "Tienda Física"
  }
]'
Buscar Reseñas (Búsqueda Semántica)
Endpoint: POST /nps/buscar (y su alias POST /conocimiento/buscar)

Descripción: Realiza una búsqueda por similitud de texto utilizando la base de datos vectorial (Qdrant). Ideal para encontrar reseñas relacionadas con un tema específico, incluso si no usan las mismas palabras.

Cuerpo (Body) de la Petición: Un objeto SearchQuery.

Ejemplo curl:

bash
curl -X POST http://localhost:8080/nps/buscar \
-H "Content-Type: application/json" \
-d '{
  "collection": "feedback_nps_2026",
  "query": "el paquete no llegó a tiempo"
}'
Obtener Reseñas Recientes
Endpoint: GET /nps/resenas

Descripción: Devuelve una lista de las últimas 50 reseñas insertadas. Acepta parámetros en la URL para realizar filtros básicos sobre los campos de la reseña.

Parámetros de URL (Opcionales): collection, categoria, estatus, sentimientoIA, clasificacionNPS, npsIA.

Ejemplo (sin filtros):

bash
curl -X GET "http://localhost:8080/nps/resenas"
Ejemplo (con filtros): (Requiere la modificación del handler sugerida en nuestra conversación anterior)

bash
# Reemplaza los IDs con valores reales de tu base de datos
curl -X GET "http://localhost:8080/nps/resenas?categoria=ID_DE_CATEGORIA&npsIA=1"
Eliminar una Reseña
Endpoint: DELETE /nps/resenas/{id}

Descripción: Elimina permanentemente una reseña de MongoDB y Qdrant usando su _id de MongoDB.

Parámetros de URL: id (el ObjectID de la reseña).

Ejemplo curl:

bash
# Reemplaza 605c7d5a9b7e4a3b1c9d6b2f con un ID válido
curl -X DELETE "http://localhost:8080/nps/resenas/605c7d5a9b7e4a3b1c9d6b2f"
Insertar una Única Reseña (Manual)
Endpoint: POST /nps/insertar (y su alias POST /conocimiento/insertar)

Descripción: Permite insertar un único documento de reseña ya completo. Es útil para migraciones o inserciones manuales donde no se requiere el enriquecimiento de la IA.

Cuerpo (Body) de la Petición: Un objeto Resena completo.

Ejemplo curl:

bash
# Los valores de _id deben existir en tus colecciones de catálogo
curl -X POST http://localhost:8080/nps/insertar \
-H "Content-Type: application/json" \
-d '{
    "collection": "feedback_nps_2026",
    "centro": {"_id": "605c7b9a9b7e4a3b1c9d6b2c"},
    "historia": "Una reseña insertada manualmente.",
    "nsp_score": 8
}'
Sección 2: Endpoints de Análisis, Estrategia e IA
Endpoints que realizan cálculos y/o utilizan el LLM para generar insights a partir de los datos existentes.

Obtener Estadísticas de NPS (Global y por Área)
Endpoint: GET /nps/stats

Descripción: Calcula el NPS total actual, el porcentaje de promotores/detractores y desglosa estas mismas métricas para cada "área" o "centro" detectado en las reseñas.

Ejemplo curl:

bash
curl -X GET "http://localhost:8080/nps/stats"
Obtener Historial de NPS
Endpoint: GET /nps/stats/historico

Descripción: Devuelve los datos de NPS y porcentajes de promotores/detractores para cada uno de los últimos 6 meses, permitiendo visualizar la evolución.

Ejemplo curl:

bash
curl -X GET "http://localhost:8080/nps/stats/historico"
Realizar Análisis de Pareto de Detractores
Endpoint: GET /nps/analisis/pareto

Descripción: Analiza todas las reseñas de detractores y agrupa los problemas por categoría y subcategoría. Ordena los resultados para mostrar qué pocos problemas están causando la mayoría de las quejas, incluyendo las áreas más afectadas.

Ejemplo curl:

bash
curl -X GET "http://localhost:8080/nps/analisis/pareto"
Generar Recomendación Estratégica
Endpoint: POST /nps/recomendar

Descripción: La IA genera una hoja de ruta accionable (queja recurrente, impacto en negocio y plan de acción) basada en la evidencia real encontrada sobre un tema de focus.

Cuerpo (Body) de la Petición: Un objeto RecommendQuery.

Ejemplo curl:

bash
curl -X POST http://localhost:8080/nps/recomendar \
-H "Content-Type: application/json" \
-d '{
  "focus": "lentitud en la aplicación móvil"
}'
Generar Predicción Estratégica
Endpoint: POST /nps/prediccion-estrategica

Descripción: Analiza las reseñas más recientes para un centro y categoría específicos y utiliza la IA para generar un informe que predice la tendencia del NPS (alza/baja/estable), identifica el foco crítico y propone una estrategia proactiva.

Cuerpo (Body) de la Petición: Un objeto PrediccionQuery con los IDs de MongoDB del centro y la categoría.

Ejemplo curl:

bash
# Reemplaza los IDs con valores reales de tus catálogos
curl -X POST http://localhost:8080/nps/prediccion-estrategica \
-H "Content-Type: application/json" \
-d '{
  "centro_id": "605c7b9a9b7e4a3b1c9d6b2c",
  "categoria_id": "605c7ba99b7e4a3b1c9d6b2d"
}'
Sección 3: Endpoints de Administración y Catálogos
Endpoints para gestionar los datos maestros (catálogos) y las colecciones de la base de datos.

CRUD para Catálogos
Descripción: Conjunto de endpoints RESTful para gestionar las colecciones de catálogo: centros, estatuses, sentimientos, clasificacionesNPS, categorias.

Paths y Métodos:

GET /{catalogo}: Listar todos los elementos.

POST /{catalogo}: Crear un nuevo elemento.

GET /{catalogo}/{id}: Obtener un elemento por ID.

PUT /{catalogo}/{id}: Actualizar un elemento por ID.

DELETE /{catalogo}/{id}: Eliminar un elemento por ID.

Ejemplos curl:

bash
# Listar todos los centros
curl -X GET "http://localhost:8080/centros"

# Crear un nuevo sentimiento
curl -X POST "http://localhost:8080/sentimientos" -d '{"nombre": "Frustrado"}'

# Obtener una categoría por ID
curl -X GET "http://localhost:8080/categorias/605c7ba99b7e4a3b1c9d6b2d"

# Actualizar un estatus por ID
curl -X PUT "http://localhost:8080/estatuses/ID_DE_ESTATUS" -d '{"nombre": "Resuelto"}'

# Eliminar un centro por ID
curl -X DELETE "http://localhost:8080/centros/ID_DE_CENTRO"
Gestión de Colecciones Vectoriales (Qdrant)
Descripción: Endpoints para administrar las colecciones en la base de datos vectorial.

Paths y Métodos:

GET /colecciones: Listar todas las colecciones existentes.

POST /colecciones (o POST /registrar-coleccion): Crear una nueva colección.

DELETE /colecciones/{nombre}: Eliminar una colección por su nombre.

Ejemplos curl:

bash
# Listar colecciones
curl -X GET "http://localhost:8080/colecciones"

# Registrar una nueva colección
curl -X POST "http://localhost:8080/colecciones" -d '{"collection_name": "feedback_2025"}'

# Eliminar una colección
curl -X DELETE "http://localhost:8080/colecciones/feedback_a_eliminar"
Limpiar Base de Datos
Endpoint: DELETE /admin/limpiar

Descripción: ¡ACCIÓN DESTRUCTIVA! Borra completamente la base de datos baseConocimientoDB de MongoDB y todas las colecciones existentes en Qdrant. Diseñado solo para entornos de desarrollo.

Ejemplo curl:

bash
curl -X DELETE "http://localhost:8080/admin/limpiar"
---