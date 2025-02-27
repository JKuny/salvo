// Package cmd
/*
Copyright © 2025 James Kuny <james.kuny@yahoo.com>
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespace string
	directory string
	verbose   bool
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Retrieve logs for all pods in a namespace",
	Long: `Retrieve logs for all pods in a namespace.
By default, the command will look at the "default" namespace. 
   salvo logs
You can specify a different namespace with the -n/--namespace flag.
   salvo logs -n my-namespace
You can also specify a target directory to write the logs to with the -d/--target-directory flag.
	salvo logs -n my-namespace -d /tmp/logs
Or just write to the current directory with them writing to ./logs/<namespace>/.
	salvo logs
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Gather flags and display them
		namespace, _ = cmd.Flags().GetString("namespace")
		directory, _ = cmd.Flags().GetString("directory")
		if directory == "" {
			directory = "./logs/" + namespace + "/"
		}
		verbose, _ = cmd.Flags().GetBool("verbose")

		if verbose {
			cmd.Printf("Using namespace \"%s\"\n", namespace)
			cmd.Printf("Writing to directory \"%s\"\n", directory)
		}

		// Pass a logger function instead of directly referencing logsCmd
		logger := func(format string, a ...interface{}) {
			if verbose {
				cmd.Printf(format, a...)
			}
		}

		// Start grabbing Kubernetes information
		err := getK8sInfo(namespace, directory, logger)
		if err != nil {
			cmd.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)

	// Setup any local flags
	logsCmd.Flags().StringP(
		"namespace",
		"n",
		"default",
		"The namespace to get logs from")
	logsCmd.Flags().StringP("directory",
		"d",
		"",
		"The file path to write the logs to")
}

// getK8sInfo Assembles the needed parts to get pod logs.
func getK8sInfo(namespace, directory string, logger func(string, ...interface{})) error {
	logger("Getting Kubernetes pods for namespace %s\n", namespace)

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	logger("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes clientset: %v", err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods in namespace %s: %v", namespace, err)
	}

	// Go through and write all the logs for the pods found
	for _, pod := range pods.Items {
		logger("Pod name: %s\n", pod.Name)
		if err := processPodLogs(clientset, pod, namespace, directory, logger); err != nil {
			return fmt.Errorf("failed to process logs for pod %s: %v", pod.Name, err)
		}
	}
	return nil
}

// processPodLogs handles streaming, reading, and saving logs for a single pod
func processPodLogs(clientset *kubernetes.Clientset, pod corev1.Pod, namespace, directory string, logger func(string, ...interface{})) error {
	podLogOptions := corev1.PodLogOptions{}
	req := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &podLogOptions)

	// Stream logs
	logStream, err := req.Stream(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to stream logs for pod %s: %v", pod.Name, err)
	}
	defer logStream.Close()

	// Process logs
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logStream)
	if err != nil {
		return fmt.Errorf("failed to copy log stream for pod %s: %v", pod.Name, err)
	}

	return writeLogs(buf.String(), pod, directory, logger)
}

// writeLogs writes `.log` files to the directory targeted
func writeLogs(content string, pod corev1.Pod, directory string, logger func(string, ...interface{})) error {
	logger("Writing files to directory %s\n", directory)

	// Check if the directory exists before writing to it, created it if it doesn't
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", directory, err)
		}
	}

	file, err := os.Create(filepath.Join(directory, pod.Name+".log"))
	if err != nil {
		return fmt.Errorf("failed to create log file for pod %s: %v", pod.Name, err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to log file for pod %s: %v", pod.Name, err)
	}

	logger("Created file %s\n", file.Name())
	return nil
}
