package models

type ProcessOcrRequest struct {
	DocumentId          string          `json:"document_id"`
	Priority            int             `json:"priority"`
	DocumentCreatedTime string          `json:"document_created_time"`
	IsDetectFace        bool            `json:"is_detect_face"`
	FileType            FileTypes       `json:"file_type"`
	OriginalLangCode    string          `json:"original_lang_code"`
	TranslateLangCode   string          `json:"translate_lang_code"`
	Title               string          `json:"title"`
	Subject             *string         `json:"subject"`
	DetailContent       []DetailContent `json:"detail_content"`
}

type DetailContent struct {
	Id          string `json:"id"`
	CreatedTime string `json:"created_time"`
	Index       int    `json:"index"`
	Content     string `json:"content"`
}
