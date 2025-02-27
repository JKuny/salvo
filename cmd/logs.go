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
			fmt.Printf("Using namespace \"%s\"\n", namespace)
			fmt.Printf("Writing to directory \"%s\"\n", directory)
		}

		// Start grabbing Kubernetes information
		getK8sInfo()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

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

// handleError
func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// getK8sInfo Assembles the needed parts to get pod logs.
func getK8sInfo() {
	if verbose {
		fmt.Printf("Getting Kubernetes pods for namespace %s\n", namespace)
	}
	userHomeDir, err := os.UserHomeDir()
	handleError(err)

	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	if verbose {
		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	handleError(err)

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	handleError(err)

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{})
	handleError(err)

	// Go through and write all the logs for the pods found
	for _, pod := range pods.Items {
		if verbose {
			fmt.Printf("Pod name: %s\n", pod.Name)
		}
		processPodLogs(clientset, pod)
	}
}

// writeLogs write `.log` files to the directory targeted
func writeLogs(content string, pod corev1.Pod) {
	if verbose {
		fmt.Printf("Writing files to directory %s\n", directory)
	}

	// Check if the directory exists before writing to it, created it if it doesn't
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, 0755)
		handleError(err)
	}

	file, err := os.Create(filepath.Join(directory, pod.Name+".log"))
	handleError(err)
	defer func(file *os.File) {
		err := file.Close()
		handleError(err)
	}(file)

	_, err = file.WriteString(content)
	handleError(err)

	if verbose {
		fmt.Printf("Created file %s\n", file.Name())
	}
}

// processPodLogs handles streaming, reading, and saving logs for a single pod
func processPodLogs(clientset *kubernetes.Clientset, pod corev1.Pod) {
	podLogOptions := corev1.PodLogOptions{}
	req := clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &podLogOptions)

	// Stream logs
	logStream, err := req.Stream(context.TODO())
	handleError(err)
	defer func(logStream io.ReadCloser) {
		err := logStream.Close()
		handleError(err)
	}(logStream)

	// Process logs
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logStream)
	handleError(err)

	writeLogs(buf.String(), pod)
}
