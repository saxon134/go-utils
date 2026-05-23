package saOss

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

var _ SaOss = (*ossClient)(nil)

type fakeObjectProvider struct {
	putDestination     string
	putBody            string
	putOptions         uploadOptions
	putFileDestination string
	putFilePath        string
	putFileOptions     uploadOptions
	existingObjects    map[string]bool
	listObjects        []string
	deletedObjects     []string
	copiedObjects      [][2]string
	getBody            string
	stsErr             error
}

func (m *fakeObjectProvider) PutObject(destination string, reader io.Reader, options uploadOptions) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	m.putDestination = destination
	m.putBody = string(body)
	m.putOptions = options
	return nil
}

func (m *fakeObjectProvider) PutObjectFromFile(destination string, localPath string, options uploadOptions) error {
	m.putFileDestination = destination
	m.putFilePath = localPath
	m.putFileOptions = options
	return nil
}

func (m *fakeObjectProvider) IsObjectExist(key string) (bool, error) {
	return m.existingObjects[key], nil
}

func (m *fakeObjectProvider) DeleteObject(key string) error {
	m.deletedObjects = append(m.deletedObjects, key)
	return nil
}

func (m *fakeObjectProvider) ListObjects(prefix string) ([]string, error) {
	return m.listObjects, nil
}

func (m *fakeObjectProvider) CopyObject(src string, destination string) error {
	m.copiedObjects = append(m.copiedObjects, [2]string{src, destination})
	return nil
}

func (m *fakeObjectProvider) GetObject(key string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(m.getBody)), nil
}

func (m *fakeObjectProvider) StsToken(roleArn, roleSessionName string) (string, string, string, error) {
	if m.stsErr != nil {
		return "", "", "", m.stsErr
	}
	return "key-id", "key-secret", "token", nil
}

func TestProviderConstants(t *testing.T) {
	if ProviderAliyun != Provider("aliyun") {
		t.Fatalf("ProviderAliyun = %q", ProviderAliyun)
	}
	if ProviderTos != Provider("tos") {
		t.Fatalf("ProviderTos = %q", ProviderTos)
	}
}

func TestUploadDelegatesToProviderWithGeneratedPathAndOptions(t *testing.T) {
	provider := &fakeObjectProvider{}
	client := &ossClient{provider: provider}

	err := client.Upload("folder/", strings.NewReader("hello"), TageDeleteDay)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(provider.putDestination, "folder/") || provider.putDestination == "folder/" {
		t.Fatalf("destination was not generated under folder/: %q", provider.putDestination)
	}
	if provider.putBody != "hello" {
		t.Fatalf("body = %q", provider.putBody)
	}
	if !reflect.DeepEqual(provider.putOptions.Tags, TageDeleteDay) {
		t.Fatalf("tags = %#v", provider.putOptions.Tags)
	}
	if provider.putOptions.ObjectExpiresDays != 1 {
		t.Fatalf("ObjectExpiresDays = %d", provider.putOptions.ObjectExpiresDays)
	}
}

func TestUploadFromLocalFileDelegatesToProvider(t *testing.T) {
	provider := &fakeObjectProvider{}
	client := &ossClient{provider: provider}

	err := client.UploadFromLocalFile("avatar.png", "/tmp/avatar.png", TageDeleteWeek)
	if err != nil {
		t.Fatal(err)
	}

	if provider.putFileDestination != "avatar.png" {
		t.Fatalf("destination = %q", provider.putFileDestination)
	}
	if provider.putFilePath != "/tmp/avatar.png" {
		t.Fatalf("local path = %q", provider.putFilePath)
	}
	if provider.putFileOptions.ObjectExpiresDays != 7 {
		t.Fatalf("ObjectExpiresDays = %d", provider.putFileOptions.ObjectExpiresDays)
	}
}

func TestDeleteDeletesPrefixObjectsWhenPathIsDirectory(t *testing.T) {
	provider := &fakeObjectProvider{
		existingObjects: map[string]bool{},
		listObjects:     []string{"dir/a.txt", "dir/b.txt"},
	}
	client := &ossClient{provider: provider}

	err := client.Delete("dir/")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"dir/a.txt", "dir/b.txt"}
	if !reflect.DeepEqual(provider.deletedObjects, want) {
		t.Fatalf("deleted objects = %#v, want %#v", provider.deletedObjects, want)
	}
}

func TestCopyWithBucketCopiesObjectIntoDestinationFolder(t *testing.T) {
	provider := &fakeObjectProvider{
		existingObjects: map[string]bool{"source/file.txt": true},
	}
	client := &ossClient{provider: provider}

	err := client.CopyWithBucket("source/file.txt", "target/")
	if err != nil {
		t.Fatal(err)
	}

	want := [][2]string{{"source/file.txt", "target/file.txt"}}
	if !reflect.DeepEqual(provider.copiedObjects, want) {
		t.Fatalf("copied objects = %#v, want %#v", provider.copiedObjects, want)
	}
}

func TestGetTxtAndUploadTxtUseProvider(t *testing.T) {
	provider := &fakeObjectProvider{getBody: "stored text"}
	client := &ossClient{provider: provider}

	got, err := client.GetTxt("txt/path")
	if err != nil {
		t.Fatal(err)
	}
	if got != "stored text" {
		t.Fatalf("GetTxt = %q", got)
	}

	path, err := client.UploadTxt("txt/", "new text", TageDeleteMonth)
	if err != nil {
		t.Fatal(err)
	}
	if path != provider.putDestination {
		t.Fatalf("returned path = %q, provider path = %q", path, provider.putDestination)
	}
	if provider.putBody != "new text" {
		t.Fatalf("uploaded text = %q", provider.putBody)
	}
	if provider.putOptions.ObjectExpiresDays != 30 {
		t.Fatalf("ObjectExpiresDays = %d", provider.putOptions.ObjectExpiresDays)
	}
}

func TestSetAddAndDeleteUrlRoot(t *testing.T) {
	client := &ossClient{}
	client.SetUrlRoot("https://cdn.example.com/root")

	got := client.AddUrlRoot("/image.png")
	if got != "https://cdn.example.com/root/image.png" {
		t.Fatalf("AddUrlRoot = %q", got)
	}

	got = client.DeleteUrlRoot("https://cdn.example.com/root/image.png")
	if got != "image.png" {
		t.Fatalf("DeleteUrlRoot = %q", got)
	}
}

func TestStsTokenDelegatesToProvider(t *testing.T) {
	provider := &fakeObjectProvider{}
	client := &ossClient{provider: provider}

	keyID, keySecret, token, err := client.StsToken("role", "session")
	if err != nil {
		t.Fatal(err)
	}
	if keyID != "key-id" || keySecret != "key-secret" || token != "token" {
		t.Fatalf("unexpected token values: %q %q %q", keyID, keySecret, token)
	}

	provider.stsErr = errors.New("unsupported")
	_, _, _, err = client.StsToken("role", "session")
	if err == nil {
		t.Fatal("expected provider error")
	}
}

func TestParseUploadOptionsMergesTagsAndExtractsObjectExpiresDays(t *testing.T) {
	options := parseUploadOptions([]interface{}{
		TageKey{"delete-week": "7"},
		TageKey{"custom": "value"},
	})

	wantTags := TageKey{"delete-week": "7", "custom": "value"}
	if !reflect.DeepEqual(options.Tags, wantTags) {
		t.Fatalf("tags = %#v, want %#v", options.Tags, wantTags)
	}
	if options.ObjectExpiresDays != 7 {
		t.Fatalf("ObjectExpiresDays = %d", options.ObjectExpiresDays)
	}
}

func TestTosPutObjectBasicInputUsesObjectExpiresDays(t *testing.T) {
	input := tosPutObjectBasicInput("bucket-a", "path/file.txt", uploadOptions{ObjectExpiresDays: 7})

	if input.Bucket != "bucket-a" {
		t.Fatalf("Bucket = %q", input.Bucket)
	}
	if input.Key != "path/file.txt" {
		t.Fatalf("Key = %q", input.Key)
	}
	if input.ObjectExpires != 7 {
		t.Fatalf("ObjectExpires = %d", input.ObjectExpires)
	}
}
