// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package ets

import (
	"fmt"
	"github.com/yeka/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// InstallationFile is a file that contains zero or more project installations.
type InstallationFile struct {
	Path           string
	InstallationID string
}

// Decode the file in order to retrieve the project inside it.
func (i *InstallationFile) Decode() (p *Project, err error) {
	r, err := os.Open(i.Path)
	if err != nil {
		return
	}

	p, err = DecodeProject(r)
	r.Close()

	return
}

// ProjectFile is a file that contains project information.
type ProjectFile struct {
	Path string

	ProjectID         string
	InstallationFiles []InstallationFile
}

// Decode the file in order to retrieve the project info inside it.
func (pf *ProjectFile) Decode() (pi *ProjectInfo, err error) {
	r, err := os.Open(pf.Path)
	if err != nil {
		return
	}

	pi, err = DecodeProjectInfo(r)
	r.Close()

	return
}

var projectFileBaseRe = regexp.MustCompile("^(\\d).xml$")

func newProjectFile(archive *archive, file string) (projFile ProjectFile) {
	projectDir := path.Dir(file)

	projFile.Path = file
	projFile.ProjectID = projectDir

	// Search for the project installation file.
	for _, file := range archive.File {
		if path.Dir(file) != projectDir {
			continue
		}

		if matches := projectFileBaseRe.FindStringSubmatch(path.Base(file)); matches != nil {
			projFile.InstallationFiles = append(projFile.InstallationFiles, InstallationFile{
				Path:           file,
				InstallationID: matches[1],
			})
		}
	}

	return
}

var hardwareFileBaseRe = regexp.MustCompile("^(\\d).xml$")

func newHardwareFile(archive *archive, file string) (hwFile HardwareFile) {
	projectDir := path.Dir(file)

	hwFile.Path = file
	hwFile.ManufacturerID = ManufacturerID(filepath.Base(projectDir))

	return
}

// ManufacturerFile is a manufacturer file.
type ManufacturerFile struct {
	Path string

	ManufacturerID       ManufacturerID
	ApplicationProgramID ApplicationProgramID
}

// Decode the file in order to retrieve the manufacturer data inside it.
func (mf *ManufacturerFile) Decode() (md *ManufacturerData, err error) {
	r, err := os.Open(mf.Path)
	if err != nil {
		return
	}

	md, err = DecodeManufacturerData(r)
	r.Close()

	return
}

// HardwareFile is a hardware file.
type HardwareFile struct {
	Path           string
	ManufacturerID ManufacturerID
}

// Decode the file in order to retrieve the manufacturer data inside it.
func (hf *HardwareFile) Decode() (hd *HardwareData, err error) {
	r, err := os.Open(hf.Path)
	if err != nil {
		return
	}

	hd, err = DecodeHardwareData(r)
	r.Close()

	return
}

type archive struct {
	Dir  string
	File []string
}

// ExportArchive is a handle to an exported archive (.knxproj or .knxprod).
type ExportArchive struct {
	*archive

	ProjectFiles      []ProjectFile
	ManufacturerFiles []ManufacturerFile
	HardwareFiles     []HardwareFile
}

// OpenExportArchive opens the exported archive located at given path.
func OpenExportArchive(path, password string) (*ExportArchive, error) {
	tmpPath, err := ioutil.TempDir("", "ets")
	if err != nil {
		return nil, err
	}

	pwd := func(file string) string {
		return password
	}

	files, err := NestedUnzip(path, tmpPath, pwd)
	if err != nil {
		return nil, err
	}
	archive := &archive{tmpPath, files}

	ex := &ExportArchive{archive: archive}
	if err = ex.findFiles(); err != nil {
		return nil, err
	}

	return ex, nil
}

func (ex *ExportArchive) Delete() error {
	return os.RemoveAll(ex.Dir)
}

var (
	projectMetaFileRe  = regexp.MustCompile("(p|P)roject.xml$")
	projectZipFileRe   = regexp.MustCompile("(p|P)-([^.]+).zip$")
	manufacturerFileRe = regexp.MustCompile("(m|M)-([0-9a-zA-Z]+)/(m|M)-([^.]+).xml$")
	hardwareFileRe     = regexp.MustCompile("(m|M)-([0-9a-zA-Z]+)/(h|H)ardware.xml$")
	// TODO: Figure out if '/' is a universal path seperator in ZIP files.
)

func (ex *ExportArchive) findFiles() error {
	for _, file := range ex.File {
		if projectMetaFileRe.MatchString(file) {
			ex.ProjectFiles = append(ex.ProjectFiles, newProjectFile(ex.archive, file))
		} else if matches := manufacturerFileRe.FindStringSubmatch(file); matches != nil {
			fname := filepath.Base(file)
			fbase := strings.TrimSuffix(fname, filepath.Ext(file))
			ids := strings.Split(fbase, "_")
			if len(ids) != 2 {
				return fmt.Errorf("Invalid manufacturer file name %s", fname)
			}
			ex.ManufacturerFiles = append(ex.ManufacturerFiles, ManufacturerFile{
				Path:                 file,
				ManufacturerID:       ManufacturerID(ids[0]),
				ApplicationProgramID: ApplicationProgramID(ids[1]),
			})
		} else if hardwareFileRe.MatchString(file) {
			ex.HardwareFiles = append(ex.HardwareFiles, newHardwareFile(ex.archive, file))
		}
	}
	return nil
}

// Close the archive handle.
func (ex *ExportArchive) Close() error {
	// return ex.archive.Close()
	return nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func NestedUnzip(src, dest string, pwd func(string) string) ([]string, error) {
	files, err := Unzip(src, dest, pwd)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fext := filepath.Ext(file)
		if fext == ".zip" {
			fname := filepath.Base(file)
			name := strings.TrimSuffix(fname, fext)
			dest := filepath.Join(filepath.Dir(file), name)

			// Make Folder
			os.MkdirAll(dest, os.ModePerm)
			fs, err := NestedUnzip(file, dest, pwd)
			if err != nil {
				return nil, err
			}
			files = append(files, fs[:]...)
		}
	}

	return files, nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string, pwd func(string) string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(pwd(f.Name))
		}
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// only copy zip or xml files
		if !isZipOrXMLFile(fpath) {
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func isZipOrXMLFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".zip" || ext == ".xml"
}
