package models

import (
	"time"
)

type Document struct {
	ID                           string            `json:"id"` // UUID
	CreatedTime                  time.Time         `json:"created_time"`
	DetailsLastUpdatedTime       *time.Time        `json:"details_last_updated_time,omitempty"`
	DetailsLastUpdatedBy         *string           `json:"details_last_updated_by,omitempty"`
	DetailsLastUpdatedByName     *string           `json:"details_last_updated_by_name,omitempty"`
	DetailsUpdated               bool              `json:"details_updated"`
	Title                        string            `json:"title"`
	Subject                      *string           `json:"subject,omitempty"`
	TitleTranslate               *string           `json:"title_translate,omitempty"`
	SubjectTranslate             *string           `json:"subject_translate,omitempty"`
	DocumentCode                 *string           `json:"document_code,omitempty"`
	Description                  *string           `json:"description,omitempty"`
	IssuingAuthority             *string           `json:"issuing_authority,omitempty"`
	DocumentTemplateID           *string           `json:"document_template_id,omitempty"` // UUID
	Priority                     int               `json:"priority"`
	DocumentTemplateName         *string           `json:"document_template_name,omitempty"`
	InputFileURLs                []string          `json:"input_file_urls"`
	Snippet                      *string           `json:"snippet,omitempty"`
	OriginalLangCode             string            `json:"original_lang_code"`
	TranslateLangCode            string            `json:"translate_lang_code"`
	Metadata                     *string           `json:"metadata,omitempty"`
	EffectiveStartTime           *time.Time        `json:"effective_start_time,omitempty"`
	EffectiveEndTime             *time.Time        `json:"effective_end_time,omitempty"`
	Keywords                     []string          `json:"keywords,omitempty"`
	TopicIDs                     []string          `json:"topic_ids,omitempty"`           // UUIDs
	TopicRecommendIDs            []string          `json:"topic_recommend_ids,omitempty"` // UUIDs
	Signer                       *string           `json:"signer,omitempty"`
	IsDetectFace                 bool              `json:"is_detect_face"`
	CanFindDocumentByImage       bool              `json:"can_find_document_by_image"`
	UploadConfigID               *string           `json:"upload_config_id,omitempty"` // UUID
	Publisher                    *string           `json:"publisher,omitempty"`
	FileType                     FileTypes         `json:"file_type"`
	Status                       DocumentStatus    `json:"status"`
	OcrProcessStatus             ProcessStatuses   `json:"ocr_process_status"`
	FaceDetectProcessStatus      ProcessStatuses   `json:"face_detect_process_status"`
	ExtractPureInfoProcessStatus ProcessStatuses   `json:"extract_pure_info_process_status"`
	ExtractContentProcessStatus  ProcessStatuses   `json:"extract_content_process_status"`
	LegalDocumentProcessStatus   ProcessStatuses   `json:"legal_document_process_status"`
	ApproveStatus                ApproveStatus     `json:"approve_status"`
	CreatorID                    string            `json:"creator_id"`
	CreatorName                  *string           `json:"creator_name,omitempty"`
	InsertedTime                 time.Time         `json:"inserted_time"`
	IssuedTime                   *time.Time        `json:"issued_time,omitempty"`
	ExtractionTypeKey            *string           `json:"extraction_type_key,omitempty"`
	SummaryAll                   *string           `json:"summary_all,omitempty"`
	Author                       *string           `json:"author,omitempty"`
	Autograph                    *string           `json:"autograph,omitempty"`
	ReliabilityLevel             *ReliabilityLevel `json:"reliability_level,omitempty"`
	PhysicalState                *PhysicalState    `json:"physical_state,omitempty"`
	Classification               *Classification   `json:"classification,omitempty"`
	ColorType                    *ColorType        `json:"color_type,omitempty"`
	HasAttachment                *HasAttachment    `json:"has_attachment,omitempty"`
	FontID                       *string           `json:"font_id,omitempty"` // UUID
	FontCode                     *string           `json:"font_code,omitempty"`
	ResumeID                     *string           `json:"resume_id,omitempty"` // UUID
	ResumeCode                   *string           `json:"resume_code,omitempty"`
	DocumentTypeCodes            []string          `json:"document_type_codes,omitempty"`
	Privacy                      Privacy           `json:"privacy"`
	HasBackup                    bool              `json:"has_backup"`
	BackupStatus                 BackupStatus      `json:"backup_status"`
	Genre                        *DocumentGenre    `json:"genre,omitempty"`
	Form                         *DocumentForm     `json:"form,omitempty"`
	Note                         *string           `json:"note,omitempty"`
	InputSourceType              *string           `json:"input_source_type,omitempty"`
	IntegrationID                *string           `json:"integration_id,omitempty"`
}

type DocumentConfig struct {
	Name         string    `json:"name"`
	Key          string    `json:"key"`
	DataType     string    `json:"data_type"`
	Regex        string    `json:"regex"`
	BBox         []float32 `json:"bbox"`
	IsRequired   bool      `json:"is_required"`
	StartContent string    `json:"start_content"`
	EndContent   string    `json:"end_content"`
}

// Enums
type FileTypes int

const (
	FileTypeImage FileTypes = 1
	FileTypeVideo FileTypes = 2
	FileTypePdf   FileTypes = 3
	FileTypeDoc   FileTypes = 4
)

type DocumentStatus int

const (
	DocStatusNotStart DocumentStatus = 1
	DocStatusPending  DocumentStatus = 2
	DocStatusDone     DocumentStatus = 3
	DocStatusFail     DocumentStatus = 4

	// Add more statuses as needed
)

type ProcessStatuses int

const (
	Error      ProcessStatuses = -1
	Pending    ProcessStatuses = 0
	Processing ProcessStatuses = 1
	Done       ProcessStatuses = 2
)

type ApproveStatus int

const (
	ApproveStatusDraft ApproveStatus = iota
	ApproveStatusPending
	ApproveStatusApproved
	ApproveStatusRejected
)

type ReliabilityLevel int

const (
	ElectronicOriginal ReliabilityLevel = iota + 1
	Digitalization
	Mixed
)

type PhysicalState int

const (
	Good PhysicalState = iota + 1
	Normal
	Damaged
)

type Classification int

const (
	Movie Classification = iota + 1
	Image
)

type ColorType int

const (
	Color ColorType = iota + 1
	BlackWhite
)

type HasAttachment int

const (
	Yes HasAttachment = iota + 1
	No
)

type Privacy int

const (
	Public Privacy = iota
	Conditional
	Private
)

type BackupStatus int

const (
	NotBackedUp BackupStatus = iota
	BackedUp
)

type DocumentGenre int

const (
	Resolution DocumentGenre = iota
	Decree
	Directive
	Regulation
	Rule
	Announcement
	Notice
	Instruction
)

type DocumentForm int

const (
	Manuscript DocumentForm = iota
	TechnicalDocument
	AudioVisualDocument
	AudioRecordingDocument
)
