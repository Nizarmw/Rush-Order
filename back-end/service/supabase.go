package service

import (
	"mime/multipart"
	"os"

	supabase_storage_uploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"github.com/joho/godotenv"
)

type SupabaseStorage struct {
	Client *supabase_storage_uploader.Client
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		panic("Error loading .env file")
	}
	if os.Getenv("SUPABASE_URL") == "" || os.Getenv("SUPABASE_TOKEN") == "" || os.Getenv("SUPABASE_BUCKET") == "" {
		panic("Supabase environment variables are not set")
	}
}

func NewSupabaseStorage() *SupabaseStorage {

	supClient := supabase_storage_uploader.New(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_TOKEN"),
		os.Getenv("SUPABASE_BUCKET"),
	)
	return &SupabaseStorage{
		Client: supClient,
	}
}

func (s *SupabaseStorage) Upload(file *multipart.FileHeader) (string, error) {
	link, err := s.Client.Upload(file)
	if err != nil {
		return link, err
	}
	return link, nil
}

func (s *SupabaseStorage) Delete(link string) error {
	err := s.Client.Delete(link)
	if err != nil {
		return err
	}
	return nil
}
