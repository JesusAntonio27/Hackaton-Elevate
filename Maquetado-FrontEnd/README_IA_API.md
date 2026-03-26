# Integracion IA (OpenAI Compatible)

Base URL:

- `https://api.genius.coppel.services`

Protocolo:

- OpenAI-compatible (LiteLLM API)

## Endpoints clave para el frontend

- `POST /v1/chat/completions` (chat principal)
- `POST /v1/embeddings` (si se requiere RAG del lado servidor)
- `GET /v1/models` (listar modelos disponibles)
- `POST /v1/responses` (alternativa moderna a chat/completions)

## Modelos recomendados (Google Gemini)

Para tu caso de uso, los IDs que si estan disponibles en el endpoint son:

- `gemini/gemini-3-flash-preview` (rapido y mas economico para chat operativo)
- `gemini-3.1-pro-preview` (mejor calidad de razonamiento para analisis complejos)

Sugerencia practica:

- usar `gemini/gemini-3-flash-preview` por defecto;
- escalar a `gemini-3.1-pro-preview` para consultas estrategicas o reportes ejecutivos.

## Autenticacion

Header requerido:

- `Authorization: Bearer <GENIUS_API_KEY>`
- `Content-Type: application/json`

## Ejemplo request (chat completions)

```json
{
  "model": "gemini/gemini-3-flash-preview",
  "messages": [
    { "role": "system", "content": "Eres un asistente de NPS para soporte TI." },
    { "role": "user", "content": "Resume las causas principales de detraccion." }
  ],
  "temperature": 0.3
}
```

## Recomendacion de arquitectura (segura)

No exponer `GENIUS_API_KEY` en `chatbot.html` ni en JavaScript del navegador.

1. Frontend (`chatbot.html`) llama a un endpoint interno:
   - `POST /api/chat`
2. Un backend/proxy privado agrega el header `Authorization` con la key del `.env`.
3. El backend reenvia a:
   - `https://api.genius.coppel.services/v1/chat/completions`

## Ejemplo rapido de fetch (frontend -> backend propio)

```js
const response = await fetch("/api/chat", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    message: "Cuales son las principales causas de detraccion?"
  })
});

const data = await response.json();
```

## Nota TLS

Durante la validacion desde terminal, el endpoint presento problema de validacion de certificado SSL en algunas herramientas cliente.
Si ocurre en runtime, revisar cadena de certificados del servidor o usar cliente con trust store actualizado.
