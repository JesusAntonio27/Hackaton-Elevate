// main.go

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	pb "github.com/qdrant/go-client/qdrant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// =============================================================================
// CONFIGURACIÓN DE IA
// =============================================================================

const LLM_API_URL = "https://api.genius.coppel.services/openai/deployments/gemini%2Fgemini-3-flash-preview/chat/completions"
const LLM_API_KEY = "sk-vOavmAHlaCWHfMVWLxfJig"
const LLM_MODEL = "gemini/gemini-3-flash-preview"

// =============================================================================
// MODELOS DE DATOS
// =============================================================================

type ResenaFilters struct {
	Categoria        string
	Estatus          string
	SentimientoIA    string
	ClasificacionNPS string
	NpsIA            int
}

type ParetoResult struct {
	Categoria           string   `json:"categoria"`
	Subcategoria        string   `json:"subcategoria"`
	Frecuencia          int      `json:"frecuencia"`
	Porcentaje          float64  `json:"porcentaje"`
	PorcentajeAcumulado float64  `json:"porcentaje_acumulado"`
	AreasAfectadas      []string `json:"areas_afectadas"` // <-- NUEVO CAMPO
}

type NPSMetrics struct {
	NPS             float64 `json:"nps"`
	PromotoresPct   float64 `json:"promotores_pct"`
	DetractoresPct  float64 `json:"detractores_pct"`
	ResenasTotales  int     `json:"resenas_totales"`
	Detalles        map[string]int `json:"detalles"`
}

type GlobalNPSMetrics struct {
	NPSMetrics
	NPSDelta           float64 `json:"nps_delta"`
	PromotoresPctDelta float64 `json:"promotores_pct_delta"`
	DetractoresPctDelta float64 `json:"detractores_pct_delta"`
}

type MonthlyStat struct {
	Month          string  `json:"month"`
	Year           int     `json:"year"`
	NPS            float64 `json:"nps"`
	PromotoresPct  float64 `json:"promotores_pct"`
	DetractoresPct float64 `json:"detractores_pct"`
}

type PrediccionQuery struct {
	Collection string `json:"collection"`
	CentroID   string `json:"centro_id"`
	CategoriaID string `json:"categoria_id"`
}

type InputSimple struct {
	FechaCaptura string `json:"fechaCaptura"`
	NpsScore     int    `json:"nps_score"`
	Comentario   string `json:"comentario"`
	Area         string `json:"area"`
}

type Subcategoria struct {
	Nombre      string `bson:"nombre" json:"nombre"`
	Descripcion string `bson:"descripcion" json:"descripcion"`
}

type ItemCatalogo struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nombre        string             `bson:"nombre" json:"nombre"`
	Subcategorias []Subcategoria     `bson:"subcategorias,omitempty" json:"subcategorias,omitempty"`
}

type ClasificacionNPS struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nombre string             `bson:"nombre" json:"nombre"`
	Min    string             `bson:"min" json:"min"`
	Max    string             `bson:"max" json:"max"`
}

type Nivel struct {
	Categoria    ItemCatalogo `bson:"categoria" json:"categoria"`
	SubCategoria ItemCatalogo `bson:"sub_categoria" json:"sub_categoria"`
}

type Resena struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Collection         string             `bson:"collection" json:"collection"`
	Centro             ItemCatalogo       `bson:"centro" json:"centro"`
	Categoria          ItemCatalogo       `bson:"categoria" json:"categoria"` // 🟢 CORREGIDO
	Historia           string             `bson:"historia" json:"historia"`
	ProblemaPrincipal  string             `bson:"problema_pricipal" json:"problema_pricipal"`
	ProblemaSecundario string             `bson:"problema_secundario" json:"problema_secundario"`
	FechaCaptura       string             `bson:"fecha_captura" json:"fecha_captura"`
	FechaCierre        string             `bson:"fecha_cierre" json:"fecha_cierre"`
	NpsScore           int                `bson:"nsp_score" json:"nsp_score"`
	Estatus            ItemCatalogo       `bson:"estatus" json:"estatus"`
	Nivel1             Nivel              `bson:"nivel_1" json:"nivel_1"`
	Nivel2             Nivel              `bson:"nivel_2" json:"nivel_2"`
	SentimientoIA      ItemCatalogo       `bson:"sentimientoIA" json:"sentimientoIA"`
	ProblemaIA         string             `bson:"problemaIA" json:"problemaIA"`
	NpsIA              string             `bson:"npsIA" json:"npsIA"`
	RazonNpsIA         string             `bson:"razonNpsIA" json:"razonNpsIA"` // 🟢 NUEVO CAMPO
	RazonIA            string             `bson:"razonIA" json:"razonIA"`
	ClasificacionNPS   ClasificacionNPS   `bson:"clasificacionNPS" json:"clasificacionNPS"`
	Iniciativa         string             `bson:"iniciativa" json:"iniciativa"`
	IniciativaIA       string             `bson:"iniciativaIA" json:"iniciativaIA"`
	CentroSoporte      ItemCatalogo       `bson:"centroSoporte" json:"centroSoporte"`
	Respuesta          string             `bson:"respuesta" json:"respuesta"`
	RespuestaIA        string             `bson:"respuestaIA" json:"respuestaIA"`
}

type IAEnrichment struct {
	IDRespuesta        string `json:"id_respuesta"`
	ProblemaPrincipal  string `json:"problema_principal"`
	ProblemaSecundario string `json:"problema_secundario"`
	NpsIA              string `json:"npsIA"`
	RazonNPSIA         string `json:"razonNpsIA"` // 🟢 NUEVO CAMPO
	RazonIA            string `json:"razonIA"`
	IniciativaIA       string `json:"iniciativaIA"`
	RespuestaIA        string `json:"respuestaIA"`
	Nivel1Categoria    string `json:"nivel1_categoria"`
	Nivel1SubCat       string `json:"nivel1_subcat"`
	Nivel2Categoria    string `json:"nivel2_categoria"`
	Nivel2SubCat       string `json:"nivel2_subcat"`
	SentimientoIANom   string `json:"sentimiento_nombre"`
	ClasificacionNom   string `json:"clasificacion_nombre"`
	CentroInferido     string `json:"centro_inferido"`
	CentroSoporteDep   string `json:"centro_soporte_dependencia"`
}

type RecommendQuery struct {
	Collection string `json:"collection"`
	Focus      string `json:"focus"`
}

type SearchQuery struct {
	Collection string `json:"collection"`
	Query      string `json:"query"`
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMRequest struct {
	Model       string       `json:"model"`
	Messages    []LLMMessage `json:"messages"`
	Temperature float32      `json:"temperature"`
}

type LLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// =============================================================================
// LÓGICA DE ALMACENAMIENTO Y CATÁLOGOS
// =============================================================================

type Storage struct {
	mongoDB           *mongo.Database
	qdrantCollections pb.CollectionsClient
	qdrantPoints      pb.PointsClient
	vectorSize        uint64
	distance          pb.Distance
	knownCollections  sync.Map
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Mapas para Join en Memoria
func (s *Storage) getCatalogoMap(ctx context.Context, collectionName string) map[string]ItemCatalogo {
	cursor, err := s.mongoDB.Collection(collectionName).Find(ctx, bson.M{})
	m := make(map[string]ItemCatalogo)
	if err != nil {
		return m
	}
	var items []ItemCatalogo
	cursor.All(ctx, &items)
	for _, item := range items {
		m[item.Nombre] = item
	}
	return m
}

func (s *Storage) getClasificacionNPSMap(ctx context.Context) map[string]ClasificacionNPS {
	cursor, err := s.mongoDB.Collection("clasificacionesNPS").Find(ctx, bson.M{})
	m := make(map[string]ClasificacionNPS)
	if err != nil {
		return m
	}
	var items []ClasificacionNPS
	cursor.All(ctx, &items)
	for _, item := range items {
		m[item.Nombre] = item
	}
	return m
}

func getVectorForText(text string, vectorSize uint64) []float32 {
	vec := make([]float32, vectorSize)
	for i := range vec {
		vec[i] = rand.Float32()
	}
	return vec
}

// --- GESTIÓN QDRANT ---
func (s *Storage) CreateCollection(ctx context.Context, collectionName string) error {
	if _, ok := s.knownCollections.Load(collectionName); ok {
		return nil
	}
	_, err := s.qdrantCollections.Create(ctx, &pb.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig:  &pb.VectorsConfig{Config: &pb.VectorsConfig_Params{Params: &pb.VectorParams{Size: s.vectorSize, Distance: s.distance}}},
	})
	if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
		s.knownCollections.Store(collectionName, true)
		return nil
	}
	if err == nil {
		s.knownCollections.Store(collectionName, true)
	}
	return err
}

func (s *Storage) ListCollections(ctx context.Context) ([]string, error) {
	res, err := s.qdrantCollections.List(ctx, &pb.ListCollectionsRequest{})
	if err != nil {
		return nil, err
	}
	var names []string
	for _, c := range res.GetCollections() {
		names = append(names, c.GetName())
	}
	return names, nil
}

func (s *Storage) DeleteCollection(ctx context.Context, collectionName string) error {
	_, err := s.qdrantCollections.Delete(ctx, &pb.DeleteCollection{
		CollectionName: collectionName,
	})
	if err == nil {
		s.knownCollections.Delete(collectionName)
	}
	return err
}

func llamarIAPrediccion(centro string, categoria string, evidencia []*Resena) (string, error) {
	var evidenciaTexto strings.Builder
	for _, fb := range evidencia {
		evidenciaTexto.WriteString(fmt.Sprintf("- [NPS: %d] Problema: '%s' (Sentimiento: %s, Razón: %s)\n",
			fb.NpsScore, fb.ProblemaPrincipal, fb.SentimientoIA.Nombre, fb.RazonIA))
	}
	prompt := fmt.Sprintf(`Eres un Analista Predictivo de Customer Success. Analiza el comportamiento reciente del Centro "%s" en la Categoría "%s".
EVIDENCIA RECIENTE:
%s
Basado en esta evidencia, redacta un informe ejecutivo en formato Markdown estructurado exactamente así:
1. **Predicción de NPS a Corto Plazo:** (Indica claramente si la tendencia es a la BAJA, ALZA o ESTABLE, y justifica por qué basándote en la gravedad y recurrencia de los problemas o aciertos recientes).
2. **Foco Crítico:** (Cuál es el problema/acierto que más peso tiene en esta predicción).
3. **Estrategia Proactiva (3 Pasos):** (Qué acciones inmediatas y tangibles debe ejecutar el equipo de "%s" para revertir la tendencia negativa o capitalizar la positiva).`, centro, categoria, evidenciaTexto.String(), centro)

	reqBody := LLMRequest{
		Model:       LLM_MODEL,
		Messages:    []LLMMessage{{Role: "user", Content: prompt}},
		Temperature: 0.3,
	}

	jb, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", LLM_API_URL, bytes.NewBuffer(jb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-litellm-api-key", LLM_API_KEY)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	var llmResp LLMResponse
	json.NewDecoder(resp.Body).Decode(&llmResp)

	if len(llmResp.Choices) > 0 {
		return llmResp.Choices[0].Message.Content, nil
	}
	return "No se pudo generar la predicción", nil
}

// --- GESTIÓN RESEÑAS DUAL (MONGO + QDRANT) ---
func (s *Storage) InsertResena(ctx context.Context, resena *Resena) (string, error) {
	if err := s.CreateCollection(ctx, resena.Collection); err != nil {
		return "", err
	}
	if resena.ID.IsZero() {
		resena.ID = primitive.NewObjectID()
	}

	res, err := s.mongoDB.Collection(resena.Collection).InsertOne(ctx, resena)
	if err != nil {
		return "", err
	}

	objectID := res.InsertedID.(primitive.ObjectID)
	qdrantID := uuid.New().String()
	combinedText := fmt.Sprintf("Historia: %s. Problema: %s. Cat: %s. Sentimiento: %s.",
		resena.Historia, resena.ProblemaPrincipal, resena.Nivel1.Categoria.Nombre, resena.SentimientoIA.Nombre)
	vector := getVectorForText(combinedText, s.vectorSize)

	wait := true
	upsertReq := &pb.UpsertPoints{
		CollectionName: resena.Collection, Wait: &wait, Points: []*pb.PointStruct{
			{
				Id:      &pb.PointId{PointIdOptions: &pb.PointId_Uuid{Uuid: qdrantID}},
				Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: vector}}},
				Payload: map[string]*pb.Value{"mongo_id": {Kind: &pb.Value_StringValue{StringValue: objectID.Hex()}}},
			},
		},
	}

	_, err = s.qdrantPoints.Upsert(ctx, upsertReq)
	if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
		s.knownCollections.Delete(resena.Collection)
		s.CreateCollection(ctx, resena.Collection)
		_, err = s.qdrantPoints.Upsert(ctx, upsertReq)
	}
	return objectID.Hex(), err
}

func (s *Storage) SearchResenas(ctx context.Context, collectionName, query string, filters ResenaFilters) ([]*Resena, error) {
	// 1. Mantienes tu búsqueda vectorial igual
	queryVector := getVectorForText(query, s.vectorSize)
	searchResult, err := s.qdrantPoints.Search(ctx, &pb.SearchPoints{
		CollectionName: collectionName,
		Vector:         queryVector,
		Limit:          50, // Tip: Aumenta el límite si esperas filtrar mucho en Mongo
		WithPayload:    &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
	})

	if err != nil {
		return nil, err
	}

	var mongoIDs []primitive.ObjectID
	for _, point := range searchResult.GetResult() {
		if val, ok := point.GetPayload()["mongo_id"]; ok {
			objectID, _ := primitive.ObjectIDFromHex(val.GetStringValue())
			mongoIDs = append(mongoIDs, objectID)
		}
	}

	if len(mongoIDs) == 0 {
		return []*Resena{}, nil
	}

	// 2. Construcción Dinámica del Filtro de MongoDB
	// Empezamos con el filtro de IDs que ya tenías
	mongoQuery := bson.M{
		"_id": bson.M{"$in": mongoIDs},
	}

	// Agregamos filtros solo si vienen con valor
	if filters.Categoria != "" {
		id, _ := primitive.ObjectIDFromHex(filters.Categoria)
		mongoQuery["categoria"] = id
	}
	if filters.Estatus != "" {
		id, _ := primitive.ObjectIDFromHex(filters.Estatus)
		mongoQuery["estatus"] = id
	}
	if filters.SentimientoIA != "" {
		id, _ := primitive.ObjectIDFromHex(filters.SentimientoIA)
		mongoQuery["sentimientoIA"] = id
	}
	if filters.ClasificacionNPS != "" {
		id, _ := primitive.ObjectIDFromHex(filters.ClasificacionNPS)
		mongoQuery["clasificacionNPS"] = id
	}
	// Para números, si el default es 0 y 0 es un valor válido,
	// podrías necesitar un puntero *int, pero aquí asumo > 0
	if filters.NpsIA > 0 {
		mongoQuery["npsIA"] = filters.NpsIA
	}

	// 3. Ejecución de la consulta filtrada
	cursor, err := s.mongoDB.Collection(collectionName).Find(ctx, mongoQuery)
	if err != nil {
		return nil, err
	}

	var results []*Resena
	err = cursor.All(ctx, &results)
	return results, err
}

func (s *Storage) DeleteResena(ctx context.Context, collectionName string, id primitive.ObjectID) error {
	_, err := s.mongoDB.Collection(collectionName).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	wait := true
	_, err = s.qdrantPoints.Delete(ctx, &pb.DeletePoints{
		CollectionName: collectionName,
		Wait:           &wait,
		Points: &pb.PointsSelector{
			PointsSelectorOneOf: &pb.PointsSelector_Filter{
				Filter: &pb.Filter{
					Must: []*pb.Condition{
						{
							ConditionOneOf: &pb.Condition_Field{
								Field: &pb.FieldCondition{
									Key: "mongo_id",
									Match: &pb.Match{
										MatchValue: &pb.Match_Text{Text: id.Hex()},
									},
								},
							},
						},
					},
				},
			},
		},
	})
	return err
}

// =============================================================================
// LÓGICA DE IA (RAG Y CLASIFICACIÓN DINÁMICA CON DESCRIPCIONES)
// =============================================================================

func (s *Storage) llamarIAEnriquecerResena(ctx context.Context, lote []InputSimple) (map[string]IAEnrichment, error) {
	catMap := s.getCatalogoMap(ctx, "categorias")
	centrosMap := s.getCatalogoMap(ctx, "centros")

	var catNames []string
	for _, c := range catMap {
		if len(c.Subcategorias) > 0 {
			var subcats []string
			for _, sub := range c.Subcategorias {
				subcats = append(subcats, fmt.Sprintf("%s (%s)", sub.Nombre, sub.Descripcion))
			}
			catNames = append(catNames, fmt.Sprintf("%s [Subcategorías: %s]", c.Nombre, strings.Join(subcats, ", ")))
		} else {
			catNames = append(catNames, c.Nombre)
		}
	}
	categoriasStr := strings.Join(catNames, "\n")

	var centrosNames []string
	for _, c := range centrosMap {
		centrosNames = append(centrosNames, c.Nombre)
	}
	centrosStr := strings.Join(centrosNames, ", ")

	var inputIA []map[string]interface{}
	for i, r := range lote {
		inputIA = append(inputIA, map[string]interface{}{
			"id_respuesta": fmt.Sprintf("req_%d", i),
			"historia":     r.Comentario,
			"score":        r.NpsScore,
			"area_origen":  r.Area,
		})
	}
	datosJSON, _ := json.Marshal(inputIA)

	// 🟢 SUPER PROMPT ACTUALIZADO CON RAZON NPS IA
	prompt := fmt.Sprintf(`Actúa como un experto Analista de Datos de Customer Experience. Clasifica las historias usando ESTRICTAMENTE estas opciones.
OPCIONES PERMITIDAS:
- Categorías y su contexto: 
%s
- Sentimientos: Enojado, Feliz, Neutral, Triste, Sarcastico
- Clasificación NPS: Detractores, Pasivos, Promotores
- Centros válidos: %s
REGLAS:
1. Si 'area_origen' viene vacía o no tiene sentido, infiere el Centro correcto y ponlo en 'centro_inferido'. Si ya es correcto, repítelo.
2. Si el problema de un centro depende de que OTRO centro lo arregle, pon el culpable en 'centro_soporte_dependencia'. Si no hay dependencia, déjalo vacío "".
3. Identifica el 'problema_principal' y si aplica, un 'problema_secundario'.
4. Si hay problema secundario, llena 'nivel2_categoria' y 'nivel2_subcat'. Si no, déjalos vacíos "".
5. En 'nivel1_categoria' escribe SOLO el nombre exacto de la categoría (ej. "Experiencia Digital"). No incluyas descripciones.
Devuelve ÚNICAMENTE un arreglo JSON válido con esta estructura:
[{"id_respuesta": "string", "problema_principal": "string", "problema_secundario": "string", "npsIA": "1 a 10", "razonNpsIA": "por qué asignas este NPS", "razonIA": "razón principal de la queja", "iniciativaIA": "solucion sugerida", "respuestaIA": "borrador para el usuario", "nivel1_categoria": "cat permitida", "nivel1_subcat": "subcat", "nivel2_categoria": "cat permitida o vacio", "nivel2_subcat": "subcat o vacio", "sentimiento_nombre": "sentimiento", "clasificacion_nombre": "clasificacion", "centro_inferido": "centro valido", "centro_soporte_dependencia": "centro culpable o vacio"}]
Datos: %s`, categoriasStr, centrosStr, string(datosJSON))

	reqBody := LLMRequest{
		Model:       LLM_MODEL,
		Messages:    []LLMMessage{{Role: "user", Content: prompt}},
		Temperature: 0.1,
	}

	jb, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", LLM_API_URL, bytes.NewBuffer(jb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-litellm-api-key", LLM_API_KEY)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var llmResp LLMResponse
	json.NewDecoder(resp.Body).Decode(&llmResp)

	txt := strings.TrimPrefix(strings.TrimSpace(llmResp.Choices[0].Message.Content), "```json")
	txt = strings.TrimSuffix(txt, "```")

	var enriquecimientos []IAEnrichment
	json.Unmarshal([]byte(txt), &enriquecimientos)

	mapa := make(map[string]IAEnrichment)
	for _, e := range enriquecimientos {
		mapa[e.IDRespuesta] = e
	}
	return mapa, nil
}

func llamarIARecomendar(focus string, evidencia []*Resena) (string, error) {
	var evidenciaTexto strings.Builder
	for _, fb := range evidencia {
		evidenciaTexto.WriteString(fmt.Sprintf("- Historia: '%s' (Cat: %s, Sentimiento: %s)\n", fb.Historia, fb.Nivel1.Categoria.Nombre, fb.SentimientoIA.Nombre))
	}

	systemPrompt := `Eres un agente de Generative UI para un dashboard NPS. Tu trabajo es analizar evidencia de reseñas y devolver EXCLUSIVAMENTE un objeto JSON válido (sin bloques markdown, sin backticks, sin texto antes ni después del JSON).

El JSON DEBE tener esta estructura exacta:
{
  "ui_type": "pareto_insight",
  "mensaje_principal": "Texto en markdown con tu análisis narrativo completo. Usa **negritas**, listas con - y ## subtítulos.",
  "ui_data": {
    "titulo_tarjeta": "Título corto para la tarjeta visual",
    "series": [
      {"etiqueta": "Causa 1", "porcentaje": 45},
      {"etiqueta": "Causa 2", "porcentaje": 30},
      {"etiqueta": "Causa 3", "porcentaje": 15}
    ],
    "filas": [
      {"etiqueta": "Queja recurrente", "valor": "Descripción breve"},
      {"etiqueta": "Impacto en negocio", "valor": "Descripción breve"},
      {"etiqueta": "Reseñas analizadas", "valor": "N"}
    ],
    "analisis": "Párrafo breve de conclusión del análisis visual.",
    "etiquetas_tags": [
      {"color": "red", "texto": "Crítico"},
      {"color": "blue", "texto": "Plan Estratégico"}
    ]
  }
}

REGLAS:
- "series" debe tener entre 3 y 5 elementos, ordenados de mayor a menor porcentaje. Los porcentajes deben sumar ~100.
- "mensaje_principal" es tu análisis narrativo completo en markdown: incluye las quejas recurrentes, el impacto en el negocio, y un plan de acción con 3 iniciativas concretas.
- "filas" resume los datos clave en pares etiqueta/valor.
- "etiquetas_tags" usa colores: "red" para crítico, "orange" para moderado, "blue" para informativo, "green" para positivo.
- NO incluyas texto fuera del JSON. NO uses bloques de código. Solo el JSON puro.`

	userPrompt := fmt.Sprintf("Analiza este foco rojo y genera los artefactos visuales.\nTema: \"%s\"\nEvidencia real:\n%s", focus, evidenciaTexto.String())

	reqBody := LLMRequest{
		Model: LLM_MODEL,
		Messages: []LLMMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
	}

	jb, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", LLM_API_URL, bytes.NewBuffer(jb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-litellm-api-key", LLM_API_KEY)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	var llmResp LLMResponse
	json.NewDecoder(resp.Body).Decode(&llmResp)
	return llmResp.Choices[0].Message.Content, nil
}

// =============================================================================
// CONTROLADORES (HANDLERS)
// =============================================================================

type Handler struct{ storage *Storage }

func (h *Handler) PrediccionEstrategica(w http.ResponseWriter, r *http.Request) {
	var q PrediccionQuery
	json.NewDecoder(r.Body).Decode(&q)
	if q.Collection == "" {
		q.Collection = "feedback_nps_2026"
	}

	// Filtramos en Mongo por el Centro y Categoría específicos
	centroObjID, _ := primitive.ObjectIDFromHex(q.CentroID)
	catObjID, _ := primitive.ObjectIDFromHex(q.CategoriaID)
	filtro := bson.M{
		"centro._id":    centroObjID,
		"categoria._id": catObjID,
	}

	// Buscamos las últimas 20 reseñas de ese centro/categoría para ver la tendencia reciente
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(20)
	cursor, err := h.storage.mongoDB.Collection(q.Collection).Find(r.Context(), filtro, opts)
	if err != nil {
		http.Error(w, "Error consultando base de datos", http.StatusInternalServerError)
		return
	}

	var evidencia []*Resena
	cursor.All(r.Context(), &evidencia)

	if len(evidencia) == 0 {
		http.Error(w, "No hay evidencia suficiente para este Centro y Categoría", http.StatusNotFound)
		return
	}

	// Usamos el nombre del centro y categoría del primer resultado para el prompt
	nombreCentro := evidencia[0].Centro.Nombre
	nombreCat := evidencia[0].Categoria.Nombre

	prediccion, _ := llamarIAPrediccion(nombreCentro, nombreCat, evidencia)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"centro":              nombreCentro,
		"categoria":           nombreCat,
		"analisis_predictivo": prediccion,
	})
}

func (h *Handler) LimpiarBase(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.mongoDB.Drop(r.Context()); err != nil {
		http.Error(w, "Error al borrar Mongo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	res, _ := h.storage.qdrantCollections.List(r.Context(), &pb.ListCollectionsRequest{})
	for _, c := range res.GetCollections() {
		h.storage.qdrantCollections.Delete(r.Context(), &pb.DeleteCollection{CollectionName: c.GetName()})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Base de datos limpia"})
}

func (h *Handler) GetCatalog(collectionName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cursor, _ := h.storage.mongoDB.Collection(collectionName).Find(r.Context(), bson.M{})
		var results []bson.M
		cursor.All(r.Context(), &results)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func (h *Handler) GetCatalogByID(collectionName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idHex := mux.Vars(r)["id"]
		objID, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		var result bson.M
		err = h.storage.mongoDB.Collection(collectionName).FindOne(r.Context(), bson.M{"_id": objID}).Decode(&result)
		if err != nil {
			http.Error(w, "No encontrado", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func (h *Handler) PostCatalog(collectionName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload interface{}
		json.NewDecoder(r.Body).Decode(&payload)
		res, _ := h.storage.mongoDB.Collection(collectionName).InsertOne(r.Context(), payload)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"inserted_id": res.InsertedID})
	}
}

func (h *Handler) PutCatalog(collectionName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idHex := mux.Vars(r)["id"]
		objID, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		var payload bson.M
		json.NewDecoder(r.Body).Decode(&payload)
		delete(payload, "_id")
		_, err = h.storage.mongoDB.Collection(collectionName).UpdateOne(r.Context(), bson.M{"_id": objID}, bson.M{"$set": payload})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"mensaje": "Actualizado"})
	}
}

func (h *Handler) DeleteCatalog(collectionName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idHex := mux.Vars(r)["id"]
		objID, err := primitive.ObjectIDFromHex(idHex)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		_, err = h.storage.mongoDB.Collection(collectionName).DeleteOne(r.Context(), bson.M{"_id": objID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"mensaje": "Eliminado"})
	}
}

func (h *Handler) RegisterCollection(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CollectionName string `json:"collection_name"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	h.storage.CreateCollection(r.Context(), body.CollectionName)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Colección registrada."})
}

func (h *Handler) ListCollections(w http.ResponseWriter, r *http.Request) {
	names, _ := h.storage.ListCollections(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"colecciones": names})
}

func (h *Handler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	h.storage.DeleteCollection(r.Context(), vars["nombre"])
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Colección eliminada."})
}

// En tu archivo main.go, reemplaza la función GetResenas con esta:
func (h *Handler) GetResenas(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    collection := r.URL.Query().Get("collection")
    if collection == "" {
        collection = "feedback_nps_2026"
    }

    // 1. Construir el filtro dinámicamente desde los query params de la URL
    filtro := bson.M{}

    if categoriaID := r.URL.Query().Get("categoria"); categoriaID != "" {
        objID, err := primitive.ObjectIDFromHex(categoriaID)
        if err == nil {
            // Nota: El campo en MongoDB es "categoria._id" si es un subdocumento,
            // o "categoria" si solo guardas el ID. Asumiré que es "categoria._id".
            // Ajusta esto según tu struct Resena.
            filtro["categoria._id"] = objID
        }
    }

    if estatusID := r.URL.Query().Get("estatus"); estatusID != "" {
        objID, err := primitive.ObjectIDFromHex(estatusID)
        if err == nil {
            filtro["estatus._id"] = objID
        }
    }

    if sentimientoID := r.URL.Query().Get("sentimientoIA"); sentimientoID != "" {
        objID, err := primitive.ObjectIDFromHex(sentimientoID)
        if err == nil {
            filtro["sentimientoIA._id"] = objID
        }
    }
    
    if clasificacionID := r.URL.Query().Get("clasificacionNPS"); clasificacionID != "" {
        objID, err := primitive.ObjectIDFromHex(clasificacionID)
        if err == nil {
            filtro["clasificacionNPS._id"] = objID
        }
    }

    if npsIA := r.URL.Query().Get("npsIA"); npsIA != "" {
        nps, err := strconv.Atoi(npsIA)
        if err == nil {
            // En MongoDB, el campo se llama "npsIA" y parece ser un string,
            // pero en tu filtro lo quieres como número. Asegúrate de que el tipo coincida.
            // Si en la base de datos es string, el filtro sería filtro["npsIA"] = npsIA
            filtro["npsIA"] = fmt.Sprintf("%d", nps) // Asumiendo que se guarda como string
        }
    }

    // 2. Ejecutar la consulta con el filtro construido
    opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(50)
    cursor, err := h.storage.mongoDB.Collection(collection).Find(ctx, filtro, opts)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var results []Resena
    cursor.All(ctx, &results)
    if results == nil {
        results = []Resena{}
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}


func (h *Handler) DeleteResena(w http.ResponseWriter, r *http.Request) {
	idHex := mux.Vars(r)["id"]
	objID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	collection := r.URL.Query().Get("collection")
	if collection == "" {
		collection = "feedback_nps_2026"
	}
	err = h.storage.DeleteResena(r.Context(), collection, objID)
	if err != nil {
		http.Error(w, "Error al borrar: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Reseña eliminada"})
}

func (h *Handler) InsertSingleResena(w http.ResponseWriter, r *http.Request) {
	var resena Resena
	if err := json.NewDecoder(r.Body).Decode(&resena); err != nil {
		http.Error(w, "Cuerpo inválido", http.StatusBadRequest)
		return
	}
	if resena.Collection == "" {
		resena.Collection = "feedback_nps_2026"
	}
	id, err := h.storage.InsertResena(r.Context(), &resena)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *Handler) SearchResenas(w http.ResponseWriter, r *http.Request) {
	var q SearchQuery
	json.NewDecoder(r.Body).Decode(&q)
	if q.Collection == "" {
		q.Collection = "feedback_nps_2026"
	}
	res, _ := h.storage.SearchResenas(r.Context(), q.Collection, q.Query, ResenaFilters{})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) RecommendStrategy(w http.ResponseWriter, r *http.Request) {
	var q RecommendQuery
	json.NewDecoder(r.Body).Decode(&q)
	if q.Collection == "" {
		q.Collection = "feedback_nps_2026"
	}
	evidencia, _ := h.storage.SearchResenas(r.Context(), q.Collection, q.Focus, ResenaFilters{})
	if len(evidencia) == 0 {
		http.Error(w, "No hay evidencia suficiente", http.StatusNotFound)
		return
	}
	recomendacion, _ := llamarIARecomendar(q.Focus, evidencia)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"recomendacion": recomendacion})
}

// GetNPSStats calcula el NPS global y por área.
func (h *Handler) GetNPSStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. OBTENER RANGOS DINÁMICOS DESDE clasificacionesNPS
	classCol := h.storage.mongoDB.Collection("clasificacionesNPS")
	cursorClas, err := classCol.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, `{"error": "Error al obtener clasificaciones NPS"}`, http.StatusInternalServerError)
		return
	}
	defer cursorClas.Close(ctx)

	type Rango struct {
		Nombre string `bson:"nombre"`
		Min    string `bson:"min"`
		Max    string `bson:"max"`
	}
	var rangos []Rango

	// CORRECCIÓN: Usar un bucle para iterar el cursor, es la práctica moderna.
	for cursorClas.Next(ctx) {
		var r Rango
		if err := cursorClas.Decode(&r); err != nil {
			log.Printf("Error al decodificar rango NPS: %v", err)
			continue
		}
		rangos = append(rangos, r)
	}

	if err := cursorClas.Err(); err != nil {
		http.Error(w, `{"error": "Error en el cursor de clasificaciones"}`, http.StatusInternalServerError)
		return
	}

	getLimits := func(nombre string) (int, int) {
		for _, r := range rangos {
			if r.Nombre == nombre {
				min, _ := strconv.Atoi(r.Min)
				max, _ := strconv.Atoi(r.Max)
				return min, max
			}
		}
		return 0, 0
	}

	minD, maxD := getLimits("Detractores")
	minP, maxP := getLimits("Pasivos")
	minPr, maxPr := getLimits("Promotores")

	// 2. CONSULTAR RESEÑAS
	collection := h.storage.mongoDB.Collection("feedback_nps_2026")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Error al consultar feedback: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var arrPromotores, arrDetractores, arrNeutrales []interface{}
	type AreaGroup struct {
		P, D, N []interface{}
	}
	areas := make(map[string]*AreaGroup)

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("Error al decodificar feedback: %v", err)
			continue
		}

		// CORRECCIÓN 1: Usar "nsp_score" en lugar de "calificacion".
		var score int
		if scoreVal, ok := doc["nsp_score"]; ok {
			switch v := scoreVal.(type) {
			case int32:
				score = int(v)
			case float64:
				score = int(v)
			case int:
				score = v
			// Añadir soporte para string si es necesario
			case string:
				score, _ = strconv.Atoi(v)
			}
		}

		// CORRECCIÓN 2: Extraer el "nombre" del objeto "centro".
		var centroNombre string
		if centroObj, ok := doc["centro"].(primitive.M); ok {
			if nombre, ok := centroObj["nombre"].(string); ok {
				centroNombre = nombre
			}
		}
		if centroNombre == "" {
			centroNombre = "sin_asignar"
		}

		if _, ok := areas[centroNombre]; !ok {
			areas[centroNombre] = &AreaGroup{P: []interface{}{}, D: []interface{}{}, N: []interface{}{}}
		}

		// 3. CLASIFICACIÓN
		if score >= minPr && score <= maxPr {
			arrPromotores = append(arrPromotores, doc["_id"])
			areas[centroNombre].P = append(areas[centroNombre].P, doc["_id"])
		} else if score >= minD && score <= maxD {
			arrDetractores = append(arrDetractores, doc["_id"])
			areas[centroNombre].D = append(areas[centroNombre].D, doc["_id"])
		} else if score >= minP && score <= maxP {
			arrNeutrales = append(arrNeutrales, doc["_id"])
			areas[centroNombre].N = append(areas[centroNombre].N, doc["_id"])
		}
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, `{"error": "Error en el cursor de feedback"}`, http.StatusInternalServerError)
		return
	}

	// 4. FUNCIÓN DE CÁLCULO FINAL (Sin cambios, la lógica es correcta)
	procesarMetricas := func(p, d, n []interface{}) map[string]interface{} {
		total := len(p) + len(d) + len(n)
		if total == 0 {
			return map[string]interface{}{"nps": 0, "promotores_pct": 0, "detractores_pct": 0, "total": 0}
		}
		pctP := (float64(len(p)) / float64(total)) * 100
		pctD := (float64(len(d)) / float64(total)) * 100
		nps := pctP - pctD
		return map[string]interface{}{
			"nps":             math.Round(nps),
			"promotores_pct":  math.Round(pctP),
			"detractores_pct": math.Round(pctD),
			"resenas_totales": total,
			"detalles":        map[string]int{"p": len(p), "d": len(d), "n": len(n)},
		}
	}

	// Construcción de respuesta
	respAreas := []map[string]interface{}{}
	for id, g := range areas {
		m := procesarMetricas(g.P, g.D, g.N)
		// CORRECCIÓN 3: El identificador ahora es el nombre del centro.
		m["centro_nombre"] = id
		respAreas = append(respAreas, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"global": procesarMetricas(arrPromotores, arrDetractores, arrNeutrales),
		"areas":  respAreas,
	})
}

func (h *Handler) GetNPSHistory(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	// 1. OBTENER RANGOS NPS (Lógica necesaria para la clasificación)
	classCol := h.storage.mongoDB.Collection("clasificacionesNPS")
	cursorClas, _ := classCol.Find(ctx, bson.M{})
	type Rango struct {
		Nombre string `bson:"nombre"`
		Min    string `bson:"min"`
		Max    string `bson:"max"`
	}
	var rangos []Rango
	cursorClas.All(ctx, &rangos)
	defer cursorClas.Close(ctx)

	getLimits := func(nombre string) (int, int) {
		for _, r := range rangos {
			if r.Nombre == nombre {
				min, _ := strconv.Atoi(r.Min)
				max, _ := strconv.Atoi(r.Max)
				return min, max
			}
		}
		return 0, 0
	}
	minD, maxD := getLimits("Detractores")
	minP, maxP := getLimits("Pasivos")
	minPr, maxPr := getLimits("Promotores")

	// 2. PREPARAR ESTRUCTURAS PARA AGRUPAR EN MEMORIA
	type MonthlyDataGroup struct{ P, D, N []interface{} }
	monthlyData := make(map[string]*MonthlyDataGroup)

	// 3. CONSULTAR TODAS LAS RESEÑAS Y AGRUPAR POR MES EN MEMORIA
	collection := h.storage.mongoDB.Collection("feedback_nps_2026")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error al consultar feedback", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		var score int
		if scoreVal, ok := doc["nsp_score"]; ok {
			switch v := scoreVal.(type) {
			case int32:
				score = int(v)
			case float64:
				score = int(v)
			case int:
				score = v
			}
		}
		fechaStr, _ := doc["fecha_captura"].(string)
		if len(fechaStr) < 7 {
			continue
		}
		monthKey := fechaStr[:7] // "YYYY-MM"

		if _, ok := monthlyData[monthKey]; !ok {
			monthlyData[monthKey] = &MonthlyDataGroup{}
		}
		if score >= minPr && score <= maxPr {
			monthlyData[monthKey].P = append(monthlyData[monthKey].P, doc["_id"])
		}
		if score >= minD && score <= maxD {
			monthlyData[monthKey].D = append(monthlyData[monthKey].D, doc["_id"])
		}
		if score >= minP && score <= maxP {
			monthlyData[monthKey].N = append(monthlyData[monthKey].N, doc["_id"])
		}
	}

	// 4. FUNCIÓN AUXILIAR DE CÁLCULO
	procesarMetricas := func(p, d, n []interface{}) NPSMetrics {
		total := len(p) + len(d) + len(n)
		if total == 0 {
			return NPSMetrics{NPS: 0, PromotoresPct: 0, DetractoresPct: 0}
		}
		pctP := (float64(len(p)) / float64(total)) * 100
		pctD := (float64(len(d)) / float64(total)) * 100
		nps := pctP - pctD
		return NPSMetrics{
			NPS: math.Round(nps), PromotoresPct: math.Round(pctP), DetractoresPct: math.Round(pctD),
		}
	}

	// 5. CONSTRUIR RESPUESTA CON LOS DATOS HISTÓRICOS
	var historico []MonthlyStat
	for i := 0; i < 6; i++ {
		now := time.Now()
		targetDate := now.AddDate(0, -i, 0)
		monthKey := targetDate.Format("2006-01")
		var monthlyMetrics NPSMetrics
		if data, ok := monthlyData[monthKey]; ok {
			monthlyMetrics = procesarMetricas(data.P, data.D, data.N)
		} else {
			monthlyMetrics = procesarMetricas(nil, nil, nil) // Mes sin datos
		}
		historico = append(historico, MonthlyStat{
			Month: targetDate.Month().String()[:3], Year: targetDate.Year(),
			NPS: monthlyMetrics.NPS, PromotoresPct: monthlyMetrics.PromotoresPct, DetractoresPct: monthlyMetrics.DetractoresPct,
		})
	}
	// Invertir el resultado para que quede en orden cronológico
	for i, j := 0, len(historico)-1; i < j; i, j = i+1, j-1 {
		historico[i], historico[j] = historico[j], historico[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(historico) // Devolvemos directamente el arreglo del histórico
}

func (h *Handler) GetParetoAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := h.storage.mongoDB.Collection("feedback_nps_2026")
	filter := bson.M{"clasificacionNPS.nombre": "Detractores"}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error consultando detractores", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// Estructura para agrupar datos, incluyendo las áreas
	type problemaKey struct{ Categoria, Subcategoria string }
	type paretoData struct {
		Frecuencia int
		Areas      map[string]bool // Usamos un mapa como un 'set' para evitar áreas duplicadas
	}
	frecuencias := make(map[problemaKey]*paretoData)
	totalDetractores := 0

	for cursor.Next(ctx) {
		var resena Resena
		if err := cursor.Decode(&resena); err != nil {
			continue
		}
		cat := resena.Nivel1.Categoria.Nombre
		subcat := resena.Nivel1.SubCategoria.Nombre
		area := resena.Centro.Nombre

		if cat != "" && subcat != "" {
			key := problemaKey{Categoria: cat, Subcategoria: subcat}
			// Si es la primera vez que vemos este problema, inicializamos la estructura
			if _, ok := frecuencias[key]; !ok {
				frecuencias[key] = &paretoData{Areas: make(map[string]bool)}
			}
			// Incrementamos la frecuencia y añadimos el área al 'set'
			frecuencias[key].Frecuencia++
			if area != "" {
				frecuencias[key].Areas[area] = true
			}
			totalDetractores++
		}
	}

	if totalDetractores == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ParetoResult{})
		return
	}

	// Convertimos el mapa a una lista para poder ordenarlo
	type problemaFrecuencia struct {
		Key   problemaKey
		Datos *paretoData
	}
	listaProblemas := make([]problemaFrecuencia, 0, len(frecuencias))
	for k, v := range frecuencias {
		listaProblemas = append(listaProblemas, problemaFrecuencia{Key: k, Datos: v})
	}

	// Ordenar de mayor a menor frecuencia
	sort.Slice(listaProblemas, func(i, j int) bool {
		return listaProblemas[i].Datos.Frecuencia > listaProblemas[j].Datos.Frecuencia
	})

	// Construir la respuesta final
	resultadosPareto := make([]ParetoResult, len(listaProblemas))
	var acumulado float64 = 0.0
	for i, problema := range listaProblemas {
		porcentaje := (float64(problema.Datos.Frecuencia) / float64(totalDetractores)) * 100
		acumulado += porcentaje
		// Convertir el mapa de áreas a una lista de strings
		areasList := make([]string, 0, len(problema.Datos.Areas))
		for area := range problema.Datos.Areas {
			areasList = append(areasList, area)
		}

		resultadosPareto[i] = ParetoResult{
			Categoria:           problema.Key.Categoria,
			Subcategoria:        problema.Key.Subcategoria,
			Frecuencia:          problema.Datos.Frecuencia,
			Porcentaje:          math.Round(porcentaje*100) / 100,
			PorcentajeAcumulado: math.Round(acumulado*100) / 100,
			AreasAfectadas:      areasList, // <-- AÑADIMOS LA LISTA DE ÁREAS
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultadosPareto)
}

func (h *Handler) ProcessAndInsertBatch(w http.ResponseWriter, r *http.Request) {
	var inputList []InputSimple
	if err := json.NewDecoder(r.Body).Decode(&inputList); err != nil {
		http.Error(w, "Cuerpo inválido", http.StatusBadRequest)
		return
	}

	catMap := h.storage.getCatalogoMap(r.Context(), "categorias")
	sentMap := h.storage.getCatalogoMap(r.Context(), "sentimientos")
	centroMap := h.storage.getCatalogoMap(r.Context(), "centros")
	estatusMap := h.storage.getCatalogoMap(r.Context(), "estatuses")
	clasMap := h.storage.getClasificacionNPSMap(r.Context())

	var exitosos, fallidos int
	const tamanoLote = 10
	for i := 0; i < len(inputList); i += tamanoLote {
		fin := i + tamanoLote
		if fin > len(inputList) {
			fin = len(inputList)
		}
		loteInput := inputList[i:fin]
		log.Printf("Procesando lote %d al %d...", i+1, fin)

		datosIA, err := h.storage.llamarIAEnriquecerResena(r.Context(), loteInput)
		if err != nil {
			log.Printf("Error IA: %v", err)
		}

		for j, inData := range loteInput {
			resena := Resena{
				Collection:   "feedback_nps_2026",
				Historia:     inData.Comentario,
				FechaCaptura: inData.FechaCaptura,
				NpsScore:     inData.NpsScore,
				Estatus:      estatusMap["Activo"],
			}

			idReq := fmt.Sprintf("req_%d", j)
			if e, ok := datosIA[idReq]; ok {
				resena.ProblemaPrincipal = e.ProblemaPrincipal
				resena.ProblemaSecundario = e.ProblemaSecundario
				resena.ProblemaIA = e.ProblemaPrincipal
				resena.NpsIA = e.NpsIA
				resena.RazonNpsIA = e.RazonNPSIA // 🟢 SE GUARDA LA NUEVA EXPLICACION DEL NPS
				resena.RazonIA = e.RazonIA
				resena.Iniciativa = ""
				resena.IniciativaIA = e.IniciativaIA
				resena.RespuestaIA = e.RespuestaIA
				resena.Nivel1.Categoria = catMap[e.Nivel1Categoria]
				resena.Nivel1.SubCategoria = ItemCatalogo{Nombre: e.Nivel1SubCat}
				// 🟢 SE ASEGURA QUE LA CATEGORIA GENERAL (RAÍZ) TAMBIÉN SE LLENE
				resena.Categoria = resena.Nivel1.Categoria

				if e.Nivel2Categoria != "" {
					resena.Nivel2.Categoria = catMap[e.Nivel2Categoria]
					resena.Nivel2.SubCategoria = ItemCatalogo{Nombre: e.Nivel2SubCat}
				}
				resena.SentimientoIA = sentMap[e.SentimientoIANom]
				resena.ClasificacionNPS = clasMap[e.ClasificacionNom]

				if inData.Area != "" {
					resena.Centro = centroMap[inData.Area]
				} else if e.CentroInferido != "" {
					resena.Centro = centroMap[e.CentroInferido]
				}

				if e.CentroSoporteDep != "" {
					resena.CentroSoporte = centroMap[e.CentroSoporteDep]
				}
			} else {
				resena.Centro = centroMap[inData.Area]
			}

			_, err := h.storage.InsertResena(r.Context(), &resena)
			if err != nil {
				fallidos++
			} else {
				exitosos++
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje":  "Lote procesado",
		"exitosos": exitosos,
		"fallidos": fallidos,
	})
}

// =============================================================================
// FUNCIÓN PRINCIPAL (MAIN)
// =============================================================================

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://host.docker.internal:27017"))
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)
	log.Println("Conectado a MongoDB.")

	db := mongoClient.Database("baseConocimientoDB")

	conn, err := grpc.Dial("host.docker.internal:6334", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error conectando a Qdrant: %v", err)
	}
	defer conn.Close()
	log.Println("Conectado a Qdrant.")

	storage := &Storage{
		mongoDB:           db,
		qdrantCollections: pb.NewCollectionsClient(conn),
		qdrantPoints:      pb.NewPointsClient(conn),
		vectorSize:        384,
		distance:          pb.Distance_Cosine,
	}

	h := &Handler{storage: storage}
	r := mux.NewRouter()

	r.HandleFunc("/nps/prediccion-estrategica", h.PrediccionEstrategica).Methods("POST", "OPTIONS")
	r.HandleFunc("/nps/stats", h.GetNPSStats).Methods("GET", "OPTIONS")
	r.HandleFunc("/nps/stats/historico", h.GetNPSHistory).Methods("GET", "OPTIONS")
	r.HandleFunc("/nps/analisis/pareto", h.GetParetoAnalysis).Methods("GET", "OPTIONS")
	r.HandleFunc("/admin/limpiar", h.LimpiarBase).Methods("DELETE", "OPTIONS")

	catalogos := []string{"centros", "estatuses", "sentimientos", "clasificacionesNPS", "categorias"}
	for _, cat := range catalogos {
		r.HandleFunc(fmt.Sprintf("/%s", cat), h.GetCatalog(cat)).Methods("GET", "OPTIONS")
		r.HandleFunc(fmt.Sprintf("/%s/{id}", cat), h.GetCatalogByID(cat)).Methods("GET", "OPTIONS")
		r.HandleFunc(fmt.Sprintf("/%s", cat), h.PostCatalog(cat)).Methods("POST", "OPTIONS")
		r.HandleFunc(fmt.Sprintf("/%s/{id}", cat), h.PutCatalog(cat)).Methods("PUT", "OPTIONS")
		r.HandleFunc(fmt.Sprintf("/%s/{id}", cat), h.DeleteCatalog(cat)).Methods("DELETE", "OPTIONS")
	}

	r.HandleFunc("/colecciones", h.ListCollections).Methods("GET", "OPTIONS")
	r.HandleFunc("/colecciones", h.RegisterCollection).Methods("POST", "OPTIONS")
	r.HandleFunc("/colecciones/{nombre}", h.DeleteCollection).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/registrar-coleccion", h.RegisterCollection).Methods("POST", "OPTIONS")

	r.HandleFunc("/nps/resenas", h.GetResenas).Methods("GET", "OPTIONS")
	r.HandleFunc("/nps/resenas/{id}", h.DeleteResena).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/nps/insertar", h.InsertSingleResena).Methods("POST", "OPTIONS")
	r.HandleFunc("/conocimiento/insertar", h.InsertSingleResena).Methods("POST", "OPTIONS")
	r.HandleFunc("/nps/buscar", h.SearchResenas).Methods("POST", "OPTIONS")
	r.HandleFunc("/nps/recomendar", h.RecommendStrategy).Methods("POST", "OPTIONS")
	r.HandleFunc("/nps/procesar-lote", h.ProcessAndInsertBatch).Methods("POST", "OPTIONS")
	r.HandleFunc("/conocimiento/buscar", h.SearchResenas).Methods("POST", "OPTIONS")

	log.Println("Servidor API Enterprise iniciado en http://localhost:8080 (CORS Habilitado)")
	if err := http.ListenAndServe(":8080", corsMiddleware(r)); err != nil {
		log.Fatal(err)
	}
}

