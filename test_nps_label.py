import json
import asyncio
from procesador_asincrono import analizar_comentario_async

async def test_label():
    # Mocking the client or using the real one if environment permits
    # This is a manual test script
    semaforo = asyncio.Semaphore(1)
    
    # Test case 1: High score but has a problem
    comentario = "Todo bien solo que a veces si llegan a tardar en atender"
    print(f"Testing comment: {comentario}")
    resultado = await analizar_comentario_async(comentario, semaforo)
    print(f"Result: {json.dumps(resultado, indent=2)}")
    
    # Test case 2: Negative comment
    comentario = "El sistema falló y no pude completar mi trámite."
    print(f"\nTesting comment: {comentario}")
    resultado = await analizar_comentario_async(comentario, semaforo)
    print(f"Result: {json.dumps(resultado, indent=2)}")

    # Test case 3: Pure positive
    comentario = "Excelente servicio, el agente fue muy amable."
    print(f"\nTesting comment: {comentario}")
    resultado = await analizar_comentario_async(comentario, semaforo)
    print(f"Result: {json.dumps(resultado, indent=2)}")

if __name__ == "__main__":
    asyncio.run(test_label())
