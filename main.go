package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
	"icomm/kafkaintegration/models"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	db := initDb()
	esClient := initESClient()
	mqChan := initRabbitMQ()
	consumerConfig := &kafka.ConfigMap{
		"bootstrap.servers":  os.Getenv("BOOTSTRAP_SERVERS"),
		"security.protocol":  os.Getenv("SECURITY_PROTOCOL"),
		"sasl.mechanism":     os.Getenv("SASL_MECHANISM"),
		"group.id":           os.Getenv("GROUP_ID"),
		"sasl.username":      os.Getenv("SASL_USERNAME"),
		"sasl.password":      os.Getenv("SASL_PASSWORD"),
		"client.id":          os.Getenv("CLIENT_ID"),
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	}
	consumer, err := kafka.NewConsumer(consumerConfig)

	if err != nil {
		log.Fatal("Failed to create consumer: ", err)
	}

	defer consumer.Close()

	err = consumer.Subscribe(os.Getenv("TOPIC"), nil)
	if err != nil {
		log.Fatal("Failed to subscribe to topic: ", err)
	}

	log.Println("Waiting for messages...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigchan:
			log.Printf("Received signal: %s, exiting...", sig)
			return
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Printf("Received message: %s", string(e.Value))
				var receivedMessage models.ReceivedMessage
				// Unmarshal the message
				err := json.Unmarshal(e.Value, &receivedMessage)
				if err != nil {
					log.Fatalf("Failed to unmarshal message: %v", err)
				}
				//Process data
				processData(db, esClient, mqChan, &receivedMessage)
			case kafka.Error:
				log.Printf("Error: %v", e)
				if e.IsFatal() {
					log.Fatalf("Fatal error: %v", e)
				}
			}
		}
	}
}

func processData(db *sql.DB, esClient *elasticsearch.TypedClient, mqChan *amqp091.Channel, data *models.ReceivedMessage) {
	doc, isExisted := saveDoc(db, esClient, data)
	if isExisted {
		return
	}

	var detailContent []models.DetailContent = []models.DetailContent{}

	var content = parseContent(data.Content)

	if content != "" {
		detailContent = append(detailContent, models.DetailContent{
			Id:          uuid.NewString(),
			CreatedTime: time.Now().Format(time.RFC3339),
			Index:       0,
			Content:     content,
		})
	}

	//push to rabbitmq
	req := models.ProcessOcrRequest{
		DocumentId:          doc.ID,
		Priority:            doc.Priority,
		DocumentCreatedTime: doc.CreatedTime.Format(time.RFC3339),
		IsDetectFace:        doc.IsDetectFace,
		FileType:            doc.FileType,
		OriginalLangCode:    doc.OriginalLangCode,
		TranslateLangCode:   doc.TranslateLangCode,
		Title:               doc.Title,
		Subject:             doc.Subject,
		DetailContent:       detailContent,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Error marshalling request: %v", err)
	}

	err = mqChan.PublishWithContext(context.Background(), "", "process-ocr-requests-priority", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        reqBytes,
	})

	if err != nil {
		log.Fatalf("Error publishing message: %v", err)
	}
}

// Insert into postgres, return true if data already exists, else false
func saveDoc(db *sql.DB, esClient *elasticsearch.TypedClient, data *models.ReceivedMessage) (*models.Document, bool) {
	createdTime := time.Now()
	systemKeyId := os.Getenv("SYSTEM_KEY_ID")

	//Get title based on type

	var fileType models.FileTypes

	switch data.Type {
	case "DOC":
		fileType = models.FileTypeDoc
	case "PIC":
		fileType = models.FileTypeImage
	case "MEDIA":
		fileType = models.FileTypeVideo
	case "FILE":
		fileType = models.FileTypeDoc
	}

	issuedTime := parseDateStringToTime(data.Metadata.IssuedDate)
	inputSourceType := "tich_hop_gd_1"
	language := ParseLangCode(data.Metadata.Language)
	privacy := ParsePrivacy(data.Metadata.Mode)
	physicalState := ParsePhysicalState(data.Metadata.Format)
	reliabilityLevel := ParseReliability(data.Metadata.ConfidenceLevel)
	keywords := strings.Split(data.Metadata.Keyword, ",")
	metadata := data
	metadataBytes, err := json.Marshal(metadata)
	metadataStr := string(metadataBytes)

	query := `
    INSERT INTO documents (
    id,
    title,
    subject,
    description,
    file_type,
    created_time,
    inserted_time,
    issued_time,
    document_code,
    creator_id,
    creator_name,
    metadata,
    input_source_type,
    original_lang_code,
    translate_lang_code,
    autograph,
    privacy,
    keywords,
    physical_state,
    has_attachment,
    reliability_level,
    integration_id,
    is_detect_face,
    priority,
    input_file_urls,
    configs,
    can_find_document_by_image,
    status,
    approve_status,
    ocr_process_status,
    face_detect_process_status,
    extract_pure_info_process_status,
    extract_content_process_status,
    legal_document_process_status
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34)
    ON CONFLICT (integration_id) DO NOTHING
    RETURNING id;
    `

	if err != nil {
		log.Fatalf("Error marshalling metadata: %v", err)
	}
	creatorName := "system"

	agrs := []any{
		uuid.NewString(),
		"",
		data.Metadata.Subject,
		data.Metadata.Description,
		fileType,
		createdTime,
		createdTime,
		issuedTime,
		data.Metadata.ArcDocCode,
		systemKeyId,
		creatorName,
		metadataBytes,
		inputSourceType,
		language,
		"org",
		data.Metadata.Autograph,
		privacy,
		pq.Array(keywords),
		physicalState,
		nil,
		reliabilityLevel,
		data.ID,
		true,
		8,
		pq.Array([]string{}),
		[]byte(`[]`),
		false,
		1,
		0,
		0,
		0,
		0,
		0,
		0,
	}

	var id string
	err = db.QueryRow(query, agrs...).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows || err.Error() == "sql: no rows in result set" {
			log.Printf("Document with Integration ID %s already exists in the database", data.ID)
			return nil, true
		} else {
			log.Fatalf("Error inserting document: %v", err)
		}
	}

	//Save to elastic search
	document := models.Document{
		ID:                           id,
		Status:                       models.DocStatusNotStart,
		OcrProcessStatus:             models.Pending,
		FaceDetectProcessStatus:      models.Pending,
		ExtractPureInfoProcessStatus: models.Pending,
		ExtractContentProcessStatus:  models.Pending,
		LegalDocumentProcessStatus:   models.Pending,
		DetailsUpdated:               false,
		BackupStatus:                 models.NotBackedUp,
		HasBackup:                    false,
		ApproveStatus:                models.ApproveStatusDraft,
		Title:                        "",
		Subject:                      &data.Metadata.Subject,
		Description:                  &data.Metadata.Description,
		FileType:                     fileType,
		CreatedTime:                  createdTime,
		InsertedTime:                 createdTime,
		IssuedTime:                   issuedTime,
		DocumentCode:                 &data.Metadata.ArcDocCode,
		CreatorID:                    systemKeyId,
		CreatorName:                  &creatorName,
		Metadata:                     &metadataStr,
		InputSourceType:              &inputSourceType,
		OriginalLangCode:             language,
		TranslateLangCode:            "org",
		Autograph:                    &data.Metadata.Autograph,
		Privacy:                      privacy,
		Keywords:                     keywords,
		PhysicalState:                physicalState,
		HasAttachment:                nil,
		ReliabilityLevel:             reliabilityLevel,
		IntegrationID:                &data.ID,
		IsDetectFace:                 true,
		Priority:                     1,
		InputFileURLs:                []string{},
	}

	_, err = esClient.Index("icocr.staging.document").Document(document).Id(document.ID).Do(context.TODO())
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}

	return &document, false
}

func initDb() *sql.DB {
	psqlInfo := os.Getenv("DATABASE_URL")
	if psqlInfo == "" {
		panic("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")
	return db
}

func initESClient() *elasticsearch.TypedClient {
	addresses := strings.Split(os.Getenv("ES_ADDRESSES"), ",")
	// Initialize your Elasticsearch client here
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: addresses,
		Username:  os.Getenv("ES_USERNAME"),
		Password:  os.Getenv("ES_PASSWORD"),
	})

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return client
}

func initRabbitMQ() *amqp091.Channel {
	conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	log.Println("Successfully connected to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	_, err = ch.QueueDeclare("process-ocr-requests", true, false, false, false, amqp091.Table{
		"x-queue-type": "classic",
	})
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	return ch
}

func ParseLangCode(languages []string) string {
	if len(languages) == 0 || languages[0] == "" {
		return ""
	}

	baseLanguage := languages[0]

	switch baseLanguage {
	case "01":
		return "vi"
	case "02":
		return "en"
	case "03":
		return "fr"
	case "04":
		return "ru"
	case "05":
		return "zh"
	case "06":
		return "vi-en"
	case "07":
		return "vi-ru"
	case "08":
		return "vi-fr"
	case "09":
		return "sino_vn"
	case "10":
		return "vi-zh"
	default:
		return ""
	}
}

func ParsePrivacy(mode string) models.Privacy {
	switch mode {
	case "01":
		return models.Public
	case "02":
		return models.Conditional
	case "03":
		return models.Private
	default:
		return models.Private
	}
}

func ParsePhysicalState(state string) *models.PhysicalState {
	switch state {
	case "01":
		s := models.Good
		return &s
	case "02":
		s := models.Normal
		return &s
	case "03":
		s := models.Damaged
		return &s
	default:
		return nil
	}
}

func ParseReliability(level string) *models.ReliabilityLevel {
	switch level {
	case "01":
		r := models.ElectronicOriginal
		return &r
	case "02":
		r := models.Digitalization
		return &r
	case "03":
		r := models.Mixed
		return &r
	default:
		return nil
	}
}

func parseContent(unknownTypeContent any) string {
	switch val := unknownTypeContent.(type) {
	case string:
		return val
	case []string:
		return strings.Join(val, "\n")
	default:
		return ""
	}
}

func parseDateStringToTime(raw string) *time.Time {
	t, err := time.Parse("26/01/2006", raw)
	if err != nil {
		return nil
	}
	year := t.Year()
	if year >= 0 && year <= 9999 {
		return &t
	}

	return nil
}
