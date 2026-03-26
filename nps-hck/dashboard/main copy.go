// main.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
// MODELOS DE DATOS (CATÁLOGOS Y RESEÑAS)
// =============================================================================

// ItemCatalogo sirve para Centros, Estatus, Sentimientos, Categorías y Subcategorías
type ItemCatalogo struct {
	ID     string `bson:"_id,omitempty" json:"_id,omitempty"`
	Nombre string `bson:"nombre" json:"nombre"`
}

type ClasificacionNPS struct {
	ID     string `bson:"_id,omitempty" json:"_id,omitempty"`
	Nombre string `bson:"nombre" json:"nombre"`
	Min    string `bson:"min" json:"min"`
	Max    string `bson:"max" json:"max"`
}

type Nivel struct {
	Categoria    ItemCatalogo `bson:"categoria" json:"categoria"`
	SubCategoria ItemCatalogo `bson:"sub_categoria" json:"sub_categoria"`
}

type Resena struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Collection         string             `bson:"collection" json:"collection"` // Para Qdrant
	Centro             ItemCatalogo       `bson:"centro" json:"centro"`
	Categoria          ItemCatalogo       `bson:"categoria" json:"categoria"`
	Historia           string             `bson:"historia" json:"historia"`
	ProblemaPrincipal  string             `bson:"problema_pricipal" json:"problema_pricipal"`
	ProblemaSecundario string             `bson:"problema_secundario" json:"problema_secundario"`
	FechaCaptura       string             `bson:"fecha_captura" json:"fecha_captura"`
	FechaCierre        string             `bson:"fecha_cierre" json:"fecha_cierre"`
	NpsScore           string             `bson:"nsp_score" json:"nsp_score"`
	Estatus            ItemCatalogo       `bson:"estatus" json:"estatus"`

	// Clasificaciones Profundas
	Nivel1 Nivel `bson:"nivel_1" json:"nivel_1"`
	Nivel2 Nivel `bson:"nivel_2" json:"nivel_2"`

	// Enriquecimiento IA
	SentimientoIA    ItemCatalogo     `bson:"sentimientoIA" json:"sentimientoIA"`
	ProblemaIA       string           `bson:"problemaIA" json:"problemaIA"`
	NpsIA            string           `bson:"npsIA" json:"npsIA"`
	RazonIA          string           `bson:"razonIA" json:"razonIA"`
	ClasificacionNPS ClasificacionNPS `bson:"clasificacionNPS" json:"clasificacionNPS"`
	Iniciativa       string           `bson:"iniciativa" json:"iniciativa"`
	IniciativaIA     string           `bson:"iniciativaIA" json:"iniciativaIA"`
	CentroSoporte    ItemCatalogo     `bson:"centroSoporte" json:"centroSoporte"`
	Respuesta        string           `bson:"respuesta" json:"respuesta"`
	RespuestaIA      string           `bson:"respuestaIA" json:"respuestaIA"`
}

// Payload esperado por la IA
type IAEnrichment struct {
	IDRespuesta      string `json:"id_respuesta"`
	ProblemaIA       string `json:"problemaIA"`
	RazonIA          string `json:"razonIA"`
	IniciativaIA     string `json:"iniciativaIA"`
	RespuestaIA      string `json:"respuestaIA"`
	Nivel1Categoria  string `json:"nivel1_categoria"`
	Nivel1SubCat     string `json:"nivel1_subcat"`
	Nivel2Categoria  string `json:"nivel2_categoria"`
	Nivel2SubCat     string `json:"nivel2_subcat"`
	SentimientoIANom string `json:"sentimiento_nombre"`
	ClasificacionNom string `json:"clasificacion_nombre"`
}

type RecommendQuery struct {
	Collection string `json:"collection"`
	Focus      string `json:"focus"`
}

// Estructuras LLM
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
// LÓGICA DE ALMACENAMIENTO
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

func getVectorForText(text string, vectorSize uint64) []float32 {
	vec := make([]float32, vectorSize)
	for i := range vec {
		vec[i] = rand.Float32()
	}
	return vec
}

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

func (s *Storage) InsertResena(ctx context.Context, resena *Resena) (string, error) {
	if resena.Collection == "" {
		resena.Collection = "resenas_globales"
	}
	if err := s.CreateCollection(ctx, resena.Collection); err != nil {
		return "", err
	}

	res, err := s.mongoDB.Collection(resena.Collection).InsertOne(ctx, resena)
	if err != nil {
		return "", err
	}
	objectID := res.InsertedID.(primitive.ObjectID)
	resena.ID = objectID

	qdrantID := uuid.New().String()
	combinedText := fmt.Sprintf("Historia: %s. Problema: %s. Cat: %s. Sentimiento: %s. Iniciativa IA: %s",
		resena.Historia, resena.ProblemaIA, resena.Nivel1.Categoria.Nombre, resena.SentimientoIA.Nombre, resena.IniciativaIA)

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

func (s *Storage) SearchResena(ctx context.Context, collectionName, query string) ([]*Resena, error) {
	queryVector := getVectorForText(query, s.vectorSize)
	searchResult, err := s.qdrantPoints.Search(ctx, &pb.SearchPoints{
		CollectionName: collectionName, Vector: queryVector, Limit: 10,
		WithPayload: &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return nil, err
	}

	var mongoIDs []primitive.ObjectID
	for _, point := range searchResult.GetResult() {
		mongoIDHex := point.GetPayload()["mongo_id"].GetStringValue()
		objectID, _ := primitive.ObjectIDFromHex(mongoIDHex)
		mongoIDs = append(mongoIDs, objectID)
	}

	if len(mongoIDs) == 0 {
		return []*Resena{}, nil
	}

	cursor, err := s.mongoDB.Collection(collectionName).Find(ctx, bson.M{"_id": bson.M{"$in": mongoIDs}})
	if err != nil {
		return nil, err
	}
	var results []*Resena
	err = cursor.All(ctx, &results)
	return results, err
}

// =============================================================================
// LÓGICA DE IA Y ESTRATEGIA
// =============================================================================

func llamarIAEnriquecerResena(lote []Resena) (map[string]IAEnrichment, error) {
	var inputIA []map[string]interface{}
	for _, r := range lote {
		inputIA = append(inputIA, map[string]interface{}{"id": r.ID.Hex(), "historia": r.Historia})
	}
	datosJSON, _ := json.Marshal(inputIA)

	prompt := fmt.Sprintf(`Actúa como analista de CX. Clasifica estas historias basándote ESTRICTAMENTE en esta taxonomía:
1. Experiencia Digital (APP, WEB, UI/UX, Performance, Bugs)
2. Producto y Calidad (Producto, Servicio, Calidad, Precios, Disponibilidad, Surtido)
3. Proceso de Compra y Entrega (Entrega, Recepcion, Comunicacion, Rastreo, Interaccion, Costos)
4. Atención y Soporte al Cliente (Interaccion, Tiempos de espera, Efectividad, Amabilidad, Claridad, Disponibilidad)
5. Experiencia en Tienda Física (Visita, Atencion, Tiempos de espera, Orden / Limpieza, Claridad)
6. Servicios Financieros y Cobranza (Pagos, Credito, Cobranza, Claridad, Abonos, Prestamos)
7. Marketing y Comunicaciones (Cliente, Promocion, Descuentos, Publicidad, Engaños, Claridad, Comunicaciones, Politicas)

Para cada historia, extrae la información y devuelve UN ARREGLO JSON válido con esta estructura exacta:
[{"id_respuesta": "string_id", "problemaIA": "resumen 1 linea", "razonIA": "por qué falló", "iniciativaIA": "qué hacer para arreglarlo de raiz", "respuestaIA": "borrador de respuesta empatica al cliente", "nivel1_categoria": "Categoría Principal", "nivel1_subcat": "Subcategoría", "nivel2_categoria": "Categoría Secundaria (si aplica)", "nivel2_subcat": "Subcat Sec.", "sentimiento_nombre": "Enojado|Feliz|Neutral|Triste", "clasificacion_nombre": "Detractores|Pasivos|Promotores"}]
Datos: %s`, string(datosJSON))

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
	for i, fb := range evidencia {
		evidenciaTexto.WriteString(fmt.Sprintf("%d. Historia: '%s' (Sentimiento: %s, Problema: %s)\n", i+1, fb.Historia, fb.SentimientoIA.Nombre, fb.ProblemaIA))
	}

	prompt := fmt.Sprintf(`Eres un Director de Estrategia. Propon una hoja de ruta accionable para mejorar el NPS basándote ÚNICAMENTE en esta evidencia real.
Tema: "%s"
Evidencia:
%s
Responde en Markdown:
1. **Queja Recurrente:**
2. **Impacto en Negocio:**
3. **Plan de Acción (3 iniciativas):**`, focus, evidenciaTexto.String())

	reqBody := LLMRequest{
		Model:       LLM_MODEL,
		Messages:    []LLMMessage{{Role: "user", Content: prompt}},
		Temperature: 0.2,
	}
	jb, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", LLM_API_URL, bytes.NewBuffer(jb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-litellm-api-key", LLM_API_KEY)

	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var llmResp LLMResponse
	json.NewDecoder(resp.Body).Decode(&llmResp)
	return llmResp.Choices[0].Message.Content, nil
}

// =============================================================================
// CONTROLADORES (HANDLERS)
// =============================================================================

type Handler struct{ storage *Storage }

// --- CRUD GENÉRICO PARA CATÁLOGOS ---
func (h *Handler) GetCatalog(collectionName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cursor, err := h.storage.mongoDB.Collection(collectionName).Find(r.Context(), bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var results []bson.M
		cursor.All(r.Context(), &results)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func (h *Handler) PostCatalog(collectionName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload interface{}
		json.NewDecoder(r.Body).Decode(&payload)
		res, err := h.storage.mongoDB.Collection(collectionName).InsertOne(r.Context(), payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"inserted_id": res.InsertedID})
	}
}

// --- HANDLERS PARA RESEÑAS ---
func (h *Handler) SearchResenas(w http.ResponseWriter, r *http.Request) {
	var q struct {
		Collection string `json:"collection"`
		Query      string `json:"query"`
	}
	json.NewDecoder(r.Body).Decode(&q)
	if q.Collection == "" {
		q.Collection = "resenas_globales"
	}
	res, _ := h.storage.SearchResena(r.Context(), q.Collection, q.Query)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) RecommendStrategy(w http.ResponseWriter, r *http.Request) {
	var q RecommendQuery
	json.NewDecoder(r.Body).Decode(&q)
	if q.Collection == "" {
		q.Collection = "resenas_globales"
	}

	evidencia, err := h.storage.SearchResena(r.Context(), q.Collection, q.Focus)
	if err != nil || len(evidencia) == 0 {
		http.Error(w, "No hay evidencia suficiente", http.StatusNotFound)
		return
	}

	recomendacion, err := llamarIARecomendar(q.Focus, evidencia)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"recomendacion": recomendacion})
}

func (h *Handler) ProcessAndInsertBatchResenas(w http.ResponseWriter, r *http.Request) {
	var resenas []Resena
	json.NewDecoder(r.Body).Decode(&resenas)

	for i := 0; i < len(resenas); i += 10 {
		fin := i + 10
		if fin > len(resenas) {
			fin = len(resenas)
		}
		lote := resenas[i:fin]

		// Pre-generamos IDs para poder mapear el resultado de la IA
		for j := range lote {
			if lote[j].ID.IsZero() {
				lote[j].ID = primitive.NewObjectID()
			}
		}

		datosIA, _ := llamarIAEnriquecerResena(lote)

		for j := range lote {
			registro := &lote[j]
			if enriquecido, ok := datosIA[registro.ID.Hex()]; ok {
				registro.ProblemaIA = enriquecido.ProblemaIA
				registro.RazonIA = enriquecido.RazonIA
				registro.IniciativaIA = enriquecido.IniciativaIA
				registro.RespuestaIA = enriquecido.RespuestaIA

				registro.Nivel1.Categoria.Nombre = enriquecido.Nivel1Categoria
				registro.Nivel1.SubCategoria.Nombre = enriquecido.Nivel1SubCat
				registro.Nivel2.Categoria.Nombre = enriquecido.Nivel2Categoria
				registro.Nivel2.SubCategoria.Nombre = enriquecido.Nivel2SubCat

				registro.SentimientoIA.Nombre = enriquecido.SentimientoIANom
				registro.ClasificacionNPS.Nombre = enriquecido.ClasificacionNom
			}
			h.storage.InsertResena(r.Context(), registro)
		}
		time.Sleep(500 * time.Millisecond)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "lote procesado e insertado"})
}

// =============================================================================
// FUNCIÓN PRINCIPAL (MAIN)
// =============================================================================

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatalf("Error Mongo: %v", err)
	}
	defer mongoClient.Disconnect(ctx)
	log.Println("Conectado a MongoDB.")
	db := mongoClient.Database("baseConocimientoDB")

	conn, err := grpc.Dial("localhost:6334", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error Qdrant: %v", err)
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

	// --- 1. RUTAS CRUD PARA CATÁLOGOS ---
	catalogos := []string{"centros", "estatuses", "sentimientos", "clasificacionesNPS", "categorias"}
	for _, cat := range catalogos {
		r.HandleFunc(fmt.Sprintf("/%s", cat), h.GetCatalog(cat)).Methods("GET", "OPTIONS")
		r.HandleFunc(fmt.Sprintf("/%s", cat), h.PostCatalog(cat)).Methods("POST", "OPTIONS")
	}

	// --- 2. RUTAS PARA RESEÑAS E IA ---
	r.HandleFunc("/resenas/procesar-lote", h.ProcessAndInsertBatchResenas).Methods("POST", "OPTIONS")
	r.HandleFunc("/resenas/buscar", h.SearchResenas).Methods("POST", "OPTIONS")
	r.HandleFunc("/resenas/recomendar", h.RecommendStrategy).Methods("POST", "OPTIONS")

	log.Println("🚀 Servidor API Enterprise iniciado en http://localhost:8080 (CORS Habilitado)")
	if err := http.ListenAndServe(":8080", corsMiddleware(r)); err != nil {
		log.Fatal(err)
	}
}