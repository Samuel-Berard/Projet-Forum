package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// Extensions d'images autorisées à l'upload.
var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// UploadService gère l'enregistrement des fichiers dans le dossier uploads.
type UploadService struct {
	uploadDir string
}

func InitUploadService() *UploadService {
	uploadDir := "./uploads"
	// On s'assure que le dossier uploads existe au démarrage.
	os.MkdirAll(uploadDir, 0755)
	return &UploadService{uploadDir: uploadDir}
}

// SaveImage valide l'extension, enregistre le fichier sous un nom unique
// dans le dossier uploads, et retourne son chemin public (ex : /uploads/upload-123.png).
func (s *UploadService) SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExtensions[ext] {
		return "", fmt.Errorf("format non autorisé (jpg, jpeg, png, gif, webp uniquement)")
	}

	// Création d'un fichier au nom unique dans le dossier uploads.
	dest, err := os.CreateTemp(s.uploadDir, "upload-*"+ext)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du fichier : %s", err.Error())
	}
	defer dest.Close()

	// Copie du contenu reçu dans le fichier de destination.
	if _, err := io.Copy(dest, file); err != nil {
		return "", fmt.Errorf("erreur lors de l'enregistrement du fichier : %s", err.Error())
	}

	return "/uploads/" + filepath.Base(dest.Name()), nil
}
