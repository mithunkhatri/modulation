package radio

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func GetFavoritesPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config")
	}

	appDir := filepath.Join(configDir, "modulation")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		err = os.MkdirAll(appDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(appDir, "favorites.json"), nil
}

func SaveFavorites(favorites []Station) error {
	path, err := GetFavoritesPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(favorites, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadFavorites() ([]Station, error) {
	path, err := GetFavoritesPath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Station{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var favorites []Station
	err = json.Unmarshal(data, &favorites)
	if err != nil {
		return nil, err
	}
	return favorites, nil
}
