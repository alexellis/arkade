// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package config

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func GetUserDir() string {
	home := os.Getenv("HOME")
	root := fmt.Sprintf("%s/.arkade/", home)
	return root
}

func InitUserDir() (string, error) {
	home := os.Getenv("HOME")
	root := fmt.Sprintf("%s/.arkade/", home)

	if len(home) == 0 {
		return home, fmt.Errorf("env-var HOME, not set")
	}

	binPath := path.Join(root, "/bin/")
	err := os.MkdirAll(binPath, 0700)
	if err != nil {
		return binPath, err
	}

	helmPath := path.Join(root, "/.helm/")
	helmErr := os.MkdirAll(helmPath, 0700)
	if helmErr != nil {
		return helmPath, helmErr
	}

	return root, nil
}

func GetDefaultKubeconfig() string {
	kubeConfigPath := path.Join(os.Getenv("HOME"), ".kube/config")

	if val, ok := os.LookupEnv("KUBECONFIG"); ok {
		kubeConfigPath = val
	}

	return kubeConfigPath
}

func MergeFlags(existingMap map[string]string, setOverrides []string) error {
	for _, setOverride := range setOverrides {
		// Limit the number of parts to 2 to keep `=` characters in the value.
		flag := strings.SplitN(setOverride, "=", 2)
		if len(flag) != 2 {
			return fmt.Errorf("incorrect format for custom flag `%s`", setOverride)
		}

		if strings.HasPrefix(flag[1], "'") && strings.HasSuffix(flag[1], "'") {
			flag[1] = flag[1][1 : len(flag[1])-1]
		}

		existingMap[flag[0]] = flag[1]
	}
	return nil
}

func SetKubeconfig(kubeconfigPath string) error {
	// Favour explicitly set kubeconfig
	if len(kubeconfigPath) > 0 {
		err := os.Setenv("KUBECONFIG", kubeconfigPath)
		if err != nil {
			return err
		}
	}

	kubeconfig := GetDefaultKubeconfig()

	fmt.Printf("Using Kubeconfig: %s\n", kubeconfig)
	return nil
}
