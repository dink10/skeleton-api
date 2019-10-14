package models

type ContextKey string

func (c ContextKey) String() string {
	return string(c) + "ContextKey"
}

type StorageKey string

func (c StorageKey) String() string {
	return string(c) + "StorageKey"
}
