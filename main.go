package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"icinga2.tf/lib/base"
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

type object struct {
	Name       string `json:"name"`
	Properties struct {
		Endpoints []string `json:"endpoints"`
		Parent    string   `json:"parent"`
	} `json:"properties"`
	Type string `json:"type"`
}

func main() {
	templateFiles := struct {
		root, branch, leaf *string
	}{
		flag.String("root", "", "FILE"),
		flag.String("branch", "", "FILE"),
		flag.String("leaf", "", "FILE"),
	}

	flag.Parse()

	var templates struct {
		root, branch, leaf *template.Template
	}

	for _, t := range []struct {
		n string
		f *string
		t **template.Template
	}{
		{"root", templateFiles.root, &templates.root},
		{"branch", templateFiles.branch, &templates.branch},
		{"leaf", templateFiles.leaf, &templates.leaf},
	} {
		if *t.f == "" {
			fmt.Fprintf(os.Stderr, "-%s missing", t.n)
			os.Exit(1)
		}

		txt, errRF := ioutil.ReadFile(*t.f)
		if errRF != nil {
			fmt.Fprintln(os.Stderr, errRF.Error())
			os.Exit(1)
		}

		var errPT error
		*t.t, errPT = template.New(t.n).Parse(string(txt))

		if errPT != nil {
			fmt.Fprintln(os.Stderr, errPT.Error())
			os.Exit(1)
		}
	}

	in := bufio.NewReader(os.Stdin)
	members := map[string][]string{}
	haveParents := map[string]struct{}{}
	haveChildren := map[string]struct{}{}

	for {
		ns, errRN := base.ReadNetStringFromStream(in, -1)
		if errRN != nil {
			if errRN == io.EOF {
				break
			}

			fmt.Fprintln(os.Stderr, errRN.Error())
			os.Exit(1)
		}

		var obj object
		if errUJ := json.Unmarshal(ns, &obj); errUJ != nil {
			fmt.Fprintln(os.Stderr, errUJ.Error())
			os.Exit(1)
		}

		if obj.Type == "Zone" {
			members[obj.Name] = obj.Properties.Endpoints

			if obj.Properties.Parent != "" {
				haveParents[obj.Name] = struct{}{}
				haveChildren[obj.Properties.Parent] = struct{}{}
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)

	for zone, endpoints := range members {
		var tmplt *template.Template
		if _, hasParents := haveParents[zone]; hasParents {
			if _, hasChildren := haveChildren[zone]; hasChildren {
				tmplt = templates.branch
			} else {
				tmplt = templates.leaf
			}
		} else {
			tmplt = templates.root
		}

		for _, endpoint := range endpoints {
			errET := tmplt.Execute(out, &struct{ Name, NameHex string }{endpoint, hex.EncodeToString([]byte(endpoint))})
			if errET != nil {
				fmt.Fprintln(os.Stderr, errET.Error())
				os.Exit(1)
			}
		}
	}

	out.Flush()
}
