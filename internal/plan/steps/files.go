package steps

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kwizyHQ/irex/internal/plan"
)

type CopyFilesStep struct {
	// Source filesystem (embed.FS, os.DirFS, etc.)
	FS fs.FS

	// map[sourceDir]destinationDir (paths inside FS → relative to ctx.TargetDir)
	DirectoryCopy map[string]string

	// map[sourceFile]destinationFile (paths inside FS → relative to ctx.TargetDir)
	FilesCopy map[string]string

	// Optional ReadOnly flag (if true, copied files will have read-only permissions)
	ReadOnly bool
}

func (c *CopyFilesStep) ID() string {
	return "copy:files"
}

func (c *CopyFilesStep) Name() string {
	return "Copy files and directories"
}

func (c *CopyFilesStep) Description() string {
	return "Copies files and directories into the target project"
}

func (c *CopyFilesStep) Run(ctx *plan.PlanContext) error {
	if c.FS == nil {
		return fmt.Errorf("source filesystem is nil")
	}

	// copy directories
	for srcDir, destDir := range c.DirectoryCopy {
		if err := c.copyDir(
			srcDir,
			filepath.Join(ctx.TargetDir, destDir),
			c.ReadOnly,
		); err != nil {
			return err
		}
	}

	// copy individual files
	for srcFile, destFile := range c.FilesCopy {
		if err := copyFileFromFS(
			c.FS,
			srcFile,
			filepath.Join(ctx.TargetDir, destFile),
			c.ReadOnly,
		); err != nil {
			return err
		}
	}

	return nil
}

//
// ===== Internal helpers =====
//

func (c *CopyFilesStep) copyDir(srcDir, destRoot string, readOnly bool) error {
	return fs.WalkDir(c.FS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destRoot, rel)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		return copyFileFromFS(c.FS, path, destPath, readOnly)
	})
}

func copyFileFromFS(srcFS fs.FS, srcPath, destPath string, readOnly bool) error {
	in, err := srcFS.Open(srcPath)
	if err != nil {
		return err
	}
	// fmt.Printf("srcPath: %s, destPath: %s\n", srcPath, destPath)
	defer in.Close()

	_, err = fs.Stat(srcFS, srcPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	var permission os.FileMode = 0644
	if readOnly {
		permission = 0444
	}

	out, err := os.OpenFile(
		destPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		permission,
	)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
