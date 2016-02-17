package commands_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf-experimental/mkman/commands"
)

var _ = Describe("CreateManifestsCommand", func() {
	var (
		args                []string
		cmd                 commands.CreateManifestsCommand
		configPathContents  string
		configPath          string
		tempDirPath         string
		outputManifest      *bytes.Buffer
		fixturesDir         string
		exampleManifestPath string
		stemcellPath        string
		stubPath            string
	)

	BeforeEach(func() {
		By("Locating fixtures dir")
		testDir := getDirOfCurrentFile()
		fixturesDir = filepath.Join(testDir, "..", "fixtures")

		By("Ensuring $CF_RELEASE_DIR is set")
		cfReleasePath := os.Getenv("CF_RELEASE_DIR")
		Expect(cfReleasePath).NotTo(BeEmpty(), "$CF_RELEASE_DIR must be provided")

		var err error
		tempDirPath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		By("Creating manifest template")
		manifestTemplatePath := filepath.Join(fixturesDir, "manifest.yml.template")
		templateContents, err := ioutil.ReadFile(manifestTemplatePath)
		Expect(err).NotTo(HaveOccurred())
		templateContents2 := strings.Replace(string(templateContents), "$CF_RELEASE_DIR", cfReleasePath, -1)

		stemcellPath = filepath.Join(fixturesDir, "no-image-stemcell.tgz")
		templateContents3 := strings.Replace(string(templateContents2), "$STEMCELL_PATH", stemcellPath, -1)

		exampleManifestPath = filepath.Join(tempDirPath, "manifest.yml")
		err = ioutil.WriteFile(exampleManifestPath, []byte(templateContents3), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		stubPath = filepath.Join(fixturesDir, "stub.yml")

		configPathContents = fmt.Sprintf(`
cf: %s
stemcell: %s
stubs:
- %s
`,
			cfReleasePath,
			stemcellPath,
			stubPath,
		)
		configPath = filepath.Join(tempDirPath, "config.yml")
		args = []string{configPath}

		outputManifest = &bytes.Buffer{}

		cmd = commands.CreateManifestsCommand{
			OutputWriter: outputManifest,
		}
	})

	AfterEach(func() {
		err := os.RemoveAll(tempDirPath)
		Expect(err).ShouldNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		err := ioutil.WriteFile(configPath, []byte(configPathContents), os.ModePerm)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("creates manifest without error", func() {
		err := cmd.Execute(args)
		Expect(err).NotTo(HaveOccurred())

		manifestPath := filepath.Join(tempDirPath, "output_manifest.yml")
		err = ioutil.WriteFile(manifestPath, outputManifest.Bytes(), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		diffCommand := exec.Command("diff", "-C3", manifestPath, exampleManifestPath)
		diffSession, err := gexec.Start(diffCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(diffSession).Should(gexec.Exit())
		Expect(diffSession.Out.Contents()).To(BeEmpty())
		Expect(diffSession.Err.Contents()).To(BeEmpty())
	})

	Context("when path is not provided", func() {
		BeforeEach(func() {
			args = []string{}
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("create-manifests requires PATH_TO_CONFIG"))
		})
	})

	Context("when path is not valid", func() {
		BeforeEach(func() {
			args = []string{"/bad/path"}
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("open /bad/path: no such file or directory"))
		})
	})

	Context("when the path points to an invalid config", func() {
		BeforeEach(func() {
			configPathContents = "{{"
		})

		It("returns an error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("yaml: line 1: did not find expected node content"))
		})
	})

	Context("when the manifest generator returns an error", func() {
		BeforeEach(func() {
			// force an error by giving a bad cfReleasePath
			configPathContents = fmt.Sprintf(`
cf: /not/a/valid/path
stemcell: %s
stubs:
- %s
`,
				stemcellPath,
				stubPath,
			)
		})

		It("forwards the error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when writing the output fails", func() {
		BeforeEach(func() {
			cmd = commands.CreateManifestsCommand{
				OutputWriter: &alwaysErrorWriter{},
			}
		})

		It("forwards the error", func() {
			err := cmd.Execute(args)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("writer error"))
		})
	})
})

func getDirOfCurrentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

type alwaysErrorWriter struct{}

func (w *alwaysErrorWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("writer error")
}
