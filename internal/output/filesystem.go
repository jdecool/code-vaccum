package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/jdecool/github-vacuum/internal/provider"
	log "github.com/sirupsen/logrus"
	gossh "golang.org/x/crypto/ssh"
)

type filesystemOutputFormatter struct {
	opts FilesystemOptions
}

type FilesystemOptions struct {
	Folder     string
	SSHKeyPath string
}

func newFilesystemOutput(opts FilesystemOptions) (*filesystemOutputFormatter, error) {
	return &filesystemOutputFormatter{opts}, nil
}

func (o filesystemOutputFormatter) Handle(r provider.Repository) {
	path := o.opts.Folder
	if strings.TrimSpace(path) != "" {
		path += "/"
	}
	path += r.Owner + "/" + r.Name

	if err := o.tryClone(r, path, r.SSHUrl, "SSH"); err != nil {
		if err := o.tryClone(r, path, r.CloneURL, "HTTPS"); err != nil {
			log.Errorf("Failed to clone repository %s with both SSH and HTTPS: %v", r.Fullname(), err)
			return
		}
	}

	log.Debugf("Successfully cloned repository %s to %s", r.Fullname(), path)
}

func (o filesystemOutputFormatter) tryClone(r provider.Repository, path, url, method string) error {
	if strings.TrimSpace(url) == "" {
		return fmt.Errorf("%s URL not available", method)
	}

	log.Debugf("Attempting to clone %s using %s: %s", r.Fullname(), method, url)

	cloneOptions := &git.CloneOptions{
		URL: url,
	}

	if method == "SSH" && strings.TrimSpace(o.opts.SSHKeyPath) != "" {
		auth, err := o.createSSHAuth()
		if err != nil {
			log.Debugf("Failed to create SSH auth for %s: %v", r.Fullname(), err)
			return err
		}
		cloneOptions.Auth = auth
		log.Debugf("Using SSH key from %s for %s", o.opts.SSHKeyPath, r.Fullname())
	}

	_, err := git.PlainClone(path, false, cloneOptions)

	if err != nil {
		if method == "SSH" && isSSHAuthError(err) {
			log.Debugf("SSH authentication failed for %s, will try HTTPS fallback", r.Fullname())
			return err
		}

		log.Debugf("Clone failed with %s for %s: %v", method, r.Fullname(), err)
		return err
	}

	log.Debugf("Successfully cloned %s using %s", r.Fullname(), method)
	return nil
}

func (o filesystemOutputFormatter) createSSHAuth() (transport.AuthMethod, error) {
	privateKey, err := os.ReadFile(o.opts.SSHKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key file %s: %w", o.opts.SSHKeyPath, err)
	}

	auth, err := ssh.NewPublicKeys("git", privateKey, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH auth from key %s: %w", o.opts.SSHKeyPath, err)
	}

	// Disable host key verification for simplicity (similar to ssh -o StrictHostKeyChecking=no)
	auth.HostKeyCallback = gossh.InsecureIgnoreHostKey()

	return auth, nil
}

func isSSHAuthError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "ssh: handshake failed") ||
		strings.Contains(errStr, "ssh: unable to authenticate") ||
		strings.Contains(errStr, "no supported methods remain") ||
		strings.Contains(errStr, "permission denied (publickey)")
}

func (o filesystemOutputFormatter) Flush() error {
	return nil
}
