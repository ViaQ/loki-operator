package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ViaQ/logerr/log"
	"github.com/ViaQ/loki-operator/internal/manifests"
	"sigs.k8s.io/yaml"
)

// Define the manifest options here as structured objects
type Config struct {
	Name         string
	Namespace    string
	Replicas     int `yaml:"replicas"`
	printVersion bool
	verifyConfig bool
	configFile   string
	writeToDir   string
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&c.Name, "name", "", "the name of the stack")
	f.StringVar(&c.Namespace, "namespace", "", "Namespace to deploy to")
	f.BoolVar(&c.printVersion, "version", false, "Print this builds version information")
	f.BoolVar(&c.verifyConfig, "verify-config", false, "Verify config file and exits")
	f.StringVar(&c.configFile, "config.file", "", "yaml file to load")
	f.StringVar(&c.writeToDir, "output.write-dir", "", "write each file to the specified directory")
}

var cfg *Config

func init() {
	log.Init("loki-broker")
	cfg = &Config{}
}

func main() {
	f := flag.NewFlagSet("", flag.ExitOnError)
	cfg.RegisterFlags(f)
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Error(err, "failed to parse flags")
	}

	if cfg.Name == "" {
		log.Info("-name flag is required")
		os.Exit(1)
	}

	// Convert Config to manifest.Options
	opts := manifests.Options{
		Name:      cfg.Name,
		Namespace: cfg.Namespace,
	}

	objects, err := manifests.BuildAll(opts)
	if err != nil {
		log.Error(err, "failed to build manifests")
		os.Exit(1)
	}

	for _, o := range objects {
		b, err := yaml.Marshal(o)
		if err != nil {
			log.Error(err, "failed to marshal manifest", "name", o.GetName(), "kind", o.GetObjectKind())
			continue
		}

		if cfg.writeToDir != "" {
			basename := fmt.Sprintf("%s-%s.yaml", o.GetObjectKind().GroupVersionKind().Kind, o.GetName())
			fname := strings.ToLower(path.Join(cfg.writeToDir, basename))
			if err := ioutil.WriteFile(fname, b, 0644); err != nil {
				log.Error(err, "failed to write file to directory", "path", fname)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stdout, "---\n%s", b)
		}
	}
}
