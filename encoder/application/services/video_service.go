package services

import (
	"context"
	"fmt"
	"enconder/application/repositories"
	"enconder/domain"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	if accountID == "" {
		return fmt.Errorf("R2_ACCOUNT_ID não configurado")
	}

	if accessKeyID == "" {
		return fmt.Errorf("R2_ACCESS_KEY_ID não configurado")
	}

	if secretAccessKey == "" {
		return fmt.Errorf("R2_SECRET_ACCESS_KEY não configurado")
	}

	// Endpoint S3 do Cloudflare R2.
	endpoint := fmt.Sprintf(
		"https://%s.r2.cloudflarestorage.com",
		accountID,
	)

	// Credenciais estáticas do R2.
	r2Credentials := credentials.NewStaticCredentialsProvider(
		accessKeyID,
		secretAccessKey,
		"",
	)

	// Configuração do AWS SDK para utilizar o R2.
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(r2Credentials),
	)

	if err != nil {
		return fmt.Errorf("erro ao carregar configuração do R2: %w", err)
	}

	// Cliente S3 apontando para o Cloudflare R2.
	client := s3.NewFromConfig(
		cfg,
		func(options *s3.Options) {
			options.BaseEndpoint = aws.String(endpoint)
		},
	)

	// Busca o vídeo no R2.
	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(v.Video.FilePath),
	})

	if err != nil {
		return fmt.Errorf(
			"erro ao baixar vídeo do R2: %w",
			err,
		)
	}

	defer result.Body.Close()

	// Diretório onde o encoder armazena os arquivos temporários.
	localStoragePath := os.Getenv("localStoragePath")

	if localStoragePath == "" {
		return fmt.Errorf("localStoragePath não configurado")
	}

	if err := os.MkdirAll(localStoragePath, os.ModePerm); err != nil {
		return fmt.Errorf(
			"erro ao criar diretório de armazenamento %s: %w",
			localStoragePath,
			err,
		)
	}


	localFile := filepath.Join(
		localStoragePath,
		v.Video.ID+".mp4",
	)

	// Cria o arquivo local.
	f, err := os.Create(localFile)

	if err != nil {
		return fmt.Errorf(
			"erro ao criar arquivo local %s: %w",
			localFile,
			err,
		)
	}

	defer f.Close()

	// Faz streaming do R2 diretamente para o disco.
	// Não carrega o vídeo inteiro na memória.
	_, err = io.Copy(f, result.Body)

	if err != nil {
		return fmt.Errorf(
			"erro ao salvar vídeo localmente: %w",
			err,
		)
	}

	log.Printf(
		"video %v has been stored",
		v.Video.ID,
	)

	return nil
}

func (v *VideoService) Fragment() error {
	localStoragePath := os.Getenv("localStoragePath")

	if localStoragePath == "" {
		return fmt.Errorf("localStoragePath não configurado")
	}

	videoDirectory := filepath.Join(
		localStoragePath,
		v.Video.ID,
	)

	// MkdirAll evita erro caso algum diretório pai ainda não exista.
	err := os.MkdirAll(videoDirectory, os.ModePerm)

	if err != nil {
		return fmt.Errorf(
			"erro ao criar diretório %s: %w",
			videoDirectory,
			err,
		)
	}

	source := filepath.Join(
		localStoragePath,
		v.Video.ID+".mp4",
	)

	target := filepath.Join(
		localStoragePath,
		v.Video.ID+".frag",
	)

	cmd := exec.Command(
		"mp4fragment",
		source,
		target,
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf(
			"erro ao executar mp4fragment: %w\noutput: %s",
			err,
			string(output),
		)
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/"+v.Video.ID+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("localStoragePath")+"/")
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")

	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(os.Getenv("localStoragePath")+"/"+v.Video.ID+".mp4")
	if err != nil {
		log.Println("error removing mp4 ", v.Video.ID, ".mp4")
		return  err
	}

	err = os.Remove(os.Getenv("localStoragePath")+"/"+v.Video.ID+".frag")
	if err != nil {
		log.Println("error removing frag ", v.Video.ID, ".frag")
		return  err
	}

	err = os.RemoveAll(os.Getenv("localStoragePath")+"/"+v.Video.ID)
	if err != nil {
		log.Println("error removing directory ", v.Video.ID, "/")
		return  err
	}

	log.Println("files have been removed: ", v.Video.ID)
	return  nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf(
			"======> Output: %s\n",
			string(out),
		)
	}
}

