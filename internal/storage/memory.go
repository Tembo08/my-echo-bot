package storage

// MemoryStorage - временное хранилище в памяти
type MemoryStorage struct {
	userCities map[int64]string // chatID -> city
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		userCities: make(map[int64]string),
	}
}

func (m *MemoryStorage) SaveUserCity(chatID int64, city string) {
	m.userCities[chatID] = city
}

func (m *MemoryStorage) GetUserCity(chatID int64) (string, bool) {
	city, exists := m.userCities[chatID]
	return city, exists
}
