package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository *repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (service *VideoService) Download(bucketName string) error {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(service.Video.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("localStoragePath") + "/" + service.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("o video %v foi salvo", service.Video.ID)

	return nil
}

func (service *VideoService) Fragment() error {

	if err := os.MkdirAll(os.Getenv("localStorageFragPath")+"/"+service.Video.ID, os.ModePerm); err != nil {
		log.Fatalln(err.Error())
		return err
	}

	source := service.Video.FilePath
	target := os.Getenv("localStorageFragPath") + "/" + service.Video.ID + "/" + "video-danca.frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	println(output)
	return nil
}

func (service *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, os.Getenv("localStorageFragPath")+"/"+service.Video.ID+"/"+"video-danca.frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, os.Getenv("localStorageFragPath")+"/"+service.Video.ID)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin")

	cmd := exec.Command("mp4dash", cmdArgs...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}

	return nil
}

func (service *VideoService) Finish() error {
	// if err := os.Remove(os.Getenv("localStoragePath") + "/" + service.Video.ID + "/video-danca.mp4"); err != nil {
	// 	log.Fatalf("erro ao excluir o MP4: %v", err.Error())
	// 	return err
	// }

	if err := os.Remove(os.Getenv("localStorageFragPath") + "/" + service.Video.ID + "/video-danca.frag"); err != nil {
		log.Fatalf("erro ao excluir o Frag: %v", err.Error())
		return err
	}

	return nil
}
