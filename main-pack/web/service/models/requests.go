package models

type PutNoteRequest struct {
	Credentials Credentials `json:"credentials"`
	Note        Note        `json:"note"`
}

type GetNoteRequest struct {
	Credentials Credentials `json:"credentials"`
	Name        string      `json:"name"`
}
