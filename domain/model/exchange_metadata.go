package model

type ExchangeMetadata struct {
	Id           uint64
	ExchangeName string
	Extra        map[string]string
}

type UpdateMetadataParam struct {
}

type UpdateMetadataData struct {
}

type QueryMetadataParam struct {
}

type QueryMetadataData struct {
}
