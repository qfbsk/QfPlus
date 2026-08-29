package app

import "QfPlus/internal/storage"

func (a *App) writeJSONFile(path string, v interface{}) error {
	return storage.WriteJSONFile(path, v)
}
