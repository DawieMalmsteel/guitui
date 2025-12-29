package lesson

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadLessons đọc file JSON và trả về danh sách bài học
func LoadLessons(path string) ([]Lesson, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("đéo đọc được file %s: %w", path, err)
	}

	var lessons []Lesson
	if err := json.Unmarshal(data, &lessons); err != nil {
		return nil, fmt.Errorf("json format ngu học: %w", err)
	}

	// Xử lý dữ liệu sau khi load
	for i := range lessons {
		l := &lessons[i]
		// Parse Key từ string ("A", "C#") sang theory.Note (int)
		// Hàm parseNote này nằm trong file generator.go (vì cùng package lesson nên gọi được nhau)
		l.ActualKey = parseNote(l.KeyStr)
	}

	return lessons, nil
}
