// Package cmd
/*
Copyright © 2025 James Kuny <james.kuny@yahoo.com>
*/
package cmd

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"

    v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
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
Or just write to the current directory.
	salvo logs
`,
    Run: func(cmd *cobra.Command, args []string) {
        // Gather flags and display them
        namespace := cmd.Flag("namespace").Value.String()
        directory := cmd.Flag("directory").Value.String()
        fmt.Printf("Using namespace \"%s\"\n", namespace)
        fmt.Printf("Writing to directory \"%s\"\n", directory)

        // Start grabbing Kubernetes information
        getK8sInfo(namespace)
    },
}

func init() {
    rootCmd.AddCommand(logsCmd)

    // Setup any persistent flags
    logsCmd.PersistentFlags().StringP("namespace", "n", "default", "The namespace to get logs from")

    // Setup any local flags
    // Get the current directory running in as where to place the log files
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
    logsCmd.Flags().StringP("directory", "d", exPath, "The file path to write the logs to")
}

// getK8sInfo Assembles the needed parts to get pod logs.
func getK8sInfo(optionalArgs ...string) {
    // Set to default namespace if somehow blank
    namespace := "default"
    if len(optionalArgs) > 0 {
        namespace = optionalArgs[0]
    }

    fmt.Printf("Getting Kubernetes pods for namespace %s\n", namespace)

    userHomeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("error getting user home dir: %v\n", err)
        os.Exit(1)
    }
    kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
    fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

    kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
    if err != nil {
        fmt.Printf("error getting Kubernetes config: %v\n", err)
        os.Exit(1)
    }

    clientset, err := kubernetes.NewForConfig(kubeConfig)
    if err != nil {
        fmt.Printf("error getting Kubernetes clientset: %v\n", err)
        os.Exit(1)
    }

    pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{})
    if err != nil {
        fmt.Printf("error getting pods: %v\n", err)
        os.Exit(1)
    }
    for _, pod := range pods.Items {
        fmt.Printf("Pod name: %s\n", pod.Name)
    }
}

// writeLogs outputs the
func writeLogs(directory string) {
    fmt.Printf("Writing files to directory %s\n", directory)
}
