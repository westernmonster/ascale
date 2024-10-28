package storage

import (
	"ascale/pkg/log"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gobuffalo/packr"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
)

func SignedGetURL(c context.Context, bucketName, objectKey string) (url string, err error) {
	box := packr.NewBox("./")
	jsonKey := box.String("./Done-eb54a73cbd2b.json")

	var conf *jwt.Config
	if conf, err = google.JWTConfigFromJSON([]byte(jsonKey)); err != nil {
		log.For(c).Errorf("invalid google account json. error(%+v)", err)
		return
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(7 * (24 * time.Hour)),
	}

	if url, err = storage.SignedURL(bucketName, objectKey, opts); err != nil {
		log.For(c).Errorf("Get gcloud storage signed url failed. error(%+v)", err)
		return
	}
	return
}

func SignedPutURL(c context.Context, bucketName, objectKey string) (url string, err error) {
	box := packr.NewBox("./")
	jsonKey := box.String("./Done-eb54a73cbd2b.json")

	var conf *jwt.Config
	if conf, err = google.JWTConfigFromJSON([]byte(jsonKey)); err != nil {
		log.For(c).Errorf("invalid google account json. error(%+v)", err)
		return
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "PUT",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	if url, err = storage.SignedURL(bucketName, objectKey, opts); err != nil {
		log.For(c).Errorf("Get gcloud storage signed url failed. error(%+v)", err)
		return
	}
	return
}

func DownloadFile(
	c context.Context,
	bucketName, objectKey string,
) (contentType, tmpFile string, err error) {
	box := packr.NewBox("./")
	jsonKey := box.String("./Done-eb54a73cbd2b.json")

	client, err := storage.NewClient(c, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil {
		log.For(c).Errorf("[storage] Get gcloud storage client failed. error(%+v)", err)
		return
	}
	defer client.Close()

	o := client.Bucket(bucketName).Object(objectKey)
	var attrs *storage.ObjectAttrs
	if attrs, err = o.Attrs(c); err != nil {
		log.For(c).Errorf("[storage] Get gcloud storage attrs. error(%+v)", err)
		return
	}

	contentType = attrs.ContentType
	// image/png
	// image/jpeg
	// application/pdf
	// image/heic

	var r *storage.Reader
	if r, err = client.Bucket(bucketName).Object(objectKey).NewReader(c); err != nil {
		log.For(c).Errorf("[storage] read from storage file. error(%+v)", err)
		return
	}
	defer r.Close()

	tmpFileName := filepath.Join(os.TempDir(), objectKey)

	f, err := os.Create(tmpFileName)
	if err != nil {
		log.For(c).Errorf("[storage] create temp file. error(%+v)", err)
		return
	}

	if _, err = io.Copy(f, r); err != nil {
		log.For(c).Errorf("[storage] io.Copy file. error(%+v)", err)
		return
	}

	if err = f.Close(); err != nil {
		return
	}

	tmpFile = tmpFileName

	return
}

func UploadFile(c context.Context, bucketName, objectKey, filePath string) (err error) {
	box := packr.NewBox("./")
	jsonKey := box.String("./Done-eb54a73cbd2b.json")

	client, err := storage.NewClient(c, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil {
		log.For(c).Errorf("Get gcloud storage client failed. error(%+v)", err)
		return
	}
	defer client.Close()

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	obj := client.Bucket(bucketName).Object(objectKey)
	w := obj.NewWriter(c)
	if _, err = io.Copy(w, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err = w.Close(); err != nil {
		log.For(c).Errorf("Close gcloud storage file failed. error(%+v)", err)
		return
	}
	return
}
