package util

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFileCreate,
		ReadContext:   resourceFileRead,
		UpdateContext: resourceFileUpdate,
		DeleteContext: resourceFileDelete,

		Schema: map[string]*schema.Schema{
			"remote_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"file_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"local_image": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if err := createValidate(d); err != nil {
		return diag.FromErr(err)
	}

	fileName := d.Get("file_name").(string)
	imageAddress := d.Get("remote_address")

	err := DownloadFile(fileName, imageAddress.(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error downloading file: %s", err))
	}

	_ = d.Set("local_image", fileName)

	hash, err := hash_file_md5(fileName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting file hash: %s", err))
	}

	d.SetId(hash)

	return resourceFileRead(ctx, d, m)
}

func resourceFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	file := d.Get("local_image").(string)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		d.SetId("")
	}

	return nil
}

func resourceFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChange("remote_address") {
		if errDiag := resourceFileCreate(ctx, d, m); errDiag != nil {
			return errDiag
		}
	}

	return resourceFileRead(ctx, d, m)
}

func resourceFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	file := d.Get("local_image").(string)
	err := os.Remove(file)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting file %s :%s", file, err))
	}
	return nil
}

func hash_file_md5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

func DownloadFile(filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Println("")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Fprintf(getTerraformStdout(), "\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Fprintf(getTerraformStdout(), "\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func createValidate(d *schema.ResourceData) error {
	return nil
}
