package config

// SSOConfig определяет контракт для конфигурационного объекта
type SSOConfig interface {
	Get() Config
}
