package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type options struct {
	Dest     string `short:"d" long:"destination" description:"The path on disk to store the token. Use '-' for stdout." default:"/.vault-token" required:"true" env:"TOKEN_DEST_PATH"`
	Role     string `short:"r" long:"role" description:"The name of the Vault GCP role to use for authentication" required:"true" env:"VAULT_ROLE"`
	Path     string `short:"p" long:"path" description:"The name of the mount where the GCP auth method is enabled." default:"gcp" required:"true" env:"VAULT_GCP_MOUNT_PATH"`
	MetaAddr string `short:"m" long:"metadata-addr" description:"Hostname or IP of the GCP metadata API." default:"metadata.google.internal" required:"true" env:"METADATA_ADDR"`
}

func main() {
	opts := options{}

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(2)
	}

	jwt, err := readJwtToken(opts.MetaAddr, opts.Role)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(jwt) // TODO remove

	// Authenticate to vault using the jwt token
	token, err := authenticate(opts.Role, opts.Path, jwt)
	if err != nil {
		log.Fatal(err)
	}

	// Persist the vault token to disk or print to stdout if TOKEN_DEST_PATH == "-"
	if opts.Dest == "-" {
		fmt.Println(token)
		os.Exit(0)
	}
	if err := saveToken(token, opts.Dest); err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully stored vault token at %s", opts.Dest)

	os.Exit(0)
}

// readJwtToken fetchs the instance's default identity token from the google compute metadata API
// the aud value is set to 'vault/VAULT_ROLE' as required by vault
func readJwtToken(addr, role string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	audience := "vault/" + role
	url := fmt.Sprintf("http://%s/computeMetadata/v1/instance/service-accounts/default/identity?audience=%s&format=full", addr, audience)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create HTTP GET request")
	}
	req.Header.Add("Metadata-Flavor", "Google")
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to make HTTP request to metadata API")
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read jwt token")
	}

	return string(bytes.TrimSpace(data)), nil
}

// authenticate to a vault gcp auth backend with a gcp instance identity JWT
func authenticate(role, mountPath, jwt string) (string, error) {
	client, err := vault.NewClient(nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create Vault client")
	}
	path := "auth/" + mountPath + "/login"
	body := map[string]interface{}{
		"role": role,
		"jwt":  jwt,
	}
	resp, err := client.Logical().Write(path, body)
	if err != nil {
		return "", errors.Wrap(err, "failed to get successful response")
	}
	token, err := resp.TokenID()
	if err != nil {
		return "", errors.Wrap(err, "failed to get token")
	}

	return token, nil
}

func saveToken(token, dest string) error {
	if err := ioutil.WriteFile(dest, []byte(token), 0600); err != nil {
		return errors.Wrap(err, "failed to save token")
	}
	return nil
}
