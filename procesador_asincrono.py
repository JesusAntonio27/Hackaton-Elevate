import openai
import json
from tqdm.asyncio import tqdm as async_tqdm 
import time
import asyncio
import os
from dotenv import load_dotenv

# Cargar variables de entorno desde .env
load_dotenv()

BASE_URL = "https://api.genius.coppel.services/v1"
API_KEY = os.getenv("GENIUS_API_KEY")
MODELO = "gemini-3.1-pro-preview"
MAX_TAREAS_SIMULTANEAS = 50 

async_client = openai.AsyncOpenAI(
    base_url=BASE_URL,
    api_key=API_KEY,
    timeout=45.0, 
)

async def analizar_comentario_async(comentario: str, semaforo: asyncio.Semaphore) -> dict:
    """
    Función asíncrona que analiza un comentario. Usa un semáforo para limitar la concurrencia.
    """
    async with semaforo:
        if not isinstance(comentario, str) or not comentario.strip():
            return {
                "ia_area_tematica": "N/A", 
                "ia_dominio_general": "N/A", 
                "ia_problema_resumido": "Comentario vacío",
                "es_problematica": False
            }

        mensaje_sistema = """
Eres un asistente de IA experto en análisis de feedback para Coppel. Tu función es procesar un comentario y devolver un JSON con cuatro claves: `ia_area_tematica`, `ia_dominio_general`, `ia_problema_resumido` y `es_problematica`.

- `ia_area_tematica`: El sustantivo específico del comentario (Ej: "VPN", "Sistema de Nóminas").
- `ia_dominio_general`: Clasifica en: "Conectividad", "Rendimiento Sistemas", "Hardware", "Software", "Atención y Soporte", "Procesos", "Comunicación".
- `ia_problema_resumido`: Un resumen muy conciso (máx 10 palabras).
- `es_problematica`: Un booleano (true/false) que indica si el comentario expone una queja, falla, lentitud o algo que requiera atención. 
  IMPORTANTE: Etiqueta como true si el usuario menciona un problema, incluso si el comentario parece positivo en general o el score NPS es alto. 
  Ejemplo: "Todo bien solo que a veces si llegan a tardar en atender" -> es_problematica: true.

Tu respuesta DEBE ser únicamente el objeto JSON.
"""
        mensajes = [
            {"role": "system", "content": mensaje_sistema},
            {"role": "user", "content": comentario}
        ]
        
        try:
            response = await async_client.chat.completions.create(
                model=MODELO,
                messages=mensajes,
                temperature=0.1,
            )
            contenido = response.choices[0].message.content
            if not contenido: raise ValueError("Respuesta vacía")
            return json.loads(contenido)
        except Exception as e:
            return {
                "ia_area_tematica": "Error", 
                "ia_dominio_general": "Error", 
                "ia_problema_resumido": f"Falla en análisis: {type(e).__name__}",
                "es_problematica": False
            }


async def main():
    archivo_entrada = 'datos_nps.json'
    archivo_salida = 'datos_analizados.json'

    print(f"--- Iniciando análisis ASÍNCRONO con el modelo: {MODELO} ---")
    print(f"--- Concurrencia máxima: {MAX_TAREAS_SIMULTANEAS} tareas simultáneas ---")

    try:
        with open(archivo_entrada, 'r', encoding='utf-8') as f:
            datos_originales = json.load(f)
        print(f"Se cargaron {len(datos_originales)} registros de '{archivo_entrada}'.")
    except FileNotFoundError:
        print(f"Error: No se encontró el archivo '{archivo_entrada}'.")
        return

    semaforo = asyncio.Semaphore(MAX_TAREAS_SIMULTANEAS)
    
    tareas = []
    for registro in datos_originales:
        tarea = analizar_comentario_async(registro.get("comentario", ""), semaforo)
        tareas.append(tarea)

    print("Lanzando todas las tareas a la vez. El progreso se mostrará a medida que completen...")
    resultados_analisis = await async_tqdm.gather(*tareas, desc="Analizando concurrentemente")

    datos_enriquecidos = []
    for registro, analisis in zip(datos_originales, resultados_analisis):
        datos_enriquecidos.append({**registro, **analisis})

    #archivo final
    with open(archivo_salida, 'w', encoding='utf-8') as f:
        json.dump(datos_enriquecidos, f, ensure_ascii=False, indent=4)

    print(f"\n¡Proceso completado!")
    print(f"Se ha generado el archivo '{archivo_salida}' con los datos analizados.")


if __name__ == "__main__":
    start_time = time.time()
    asyncio.run(main())
    end_time = time.time()
    print(f"--- Tiempo total de ejecución: {end_time - start_time:.2f} segundos ---")
