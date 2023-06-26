package main

import (
	"fmt"
	"os"
	"path/filepath"
        "unicode/utf8"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 根据当前用户的 Kubernetes 配置文件创建一个客户端配置
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	// 创建一个 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 生成逻辑树状结构
	root := &Node{
		Name:      "/",
		Namespace: "",
		Children:  []*Node{},
	}

	if err := generateLogicalTree(clientset, root); err != nil {
		panic(err)
	}

	// 打印逻辑树状结构
	printLogicalTree(root, 0, "")
}

// 生成逻辑树状结构
func generateLogicalTree(clientset *kubernetes.Clientset, node *Node) error {
	// 获取 Kubernetes API 的资源组列表
	groupList, err := clientset.Discovery().ServerGroups()
	if err != nil {
		return err
	}

	for _, group := range groupList.Groups {
		// 创建顶层节点
		groupNode := &Node{
			Name:      "/" + group.Name,
			Namespace: "",
			Parent:    node,
			Children:  []*Node{},
		}

		node.Children = append(node.Children, groupNode)

		// 获取组内的 API 资源列表
		apiResourceList, err := clientset.Discovery().ServerResourcesForGroupVersion(group.PreferredVersion.GroupVersion)
		if err != nil {
			return err
		}

		for _, apiResource := range apiResourceList.APIResources {
			// 创建子节点
			resourceNode := &Node{
				Name:      "/" + group.Name + "/" + apiResource.Name,
				Namespace: "",
				Parent:    groupNode,
				Children:  []*Node{},
			}

			groupNode.Children = append(groupNode.Children, resourceNode)
		}
	}

	return nil
}

// Node 表示一个节点
type Node struct {
	Name      string    // 节点名称
	Namespace string    // 节点所属的命名空间
	Parent    *Node     // 父节点
	Children  []*Node   // 子节点列表
}


// 打印逻辑树状结构
func printLogicalTree(node *Node, level int, prefix string) {
	if level > 0 {
		for i := 0; i < level-1; i++ {
			prefix += "│   "
		}
		prefix += "├── "
	}

	fmt.Printf("%s%s\n", prefix, node.Name)

	prefixLast := prefix
	if len(node.Children) > 0 {
		_, width := utf8.DecodeLastRuneInString(prefixLast)
		prefixLast = prefixLast[:len(prefixLast)-width] + "└── "
	}

	for i, child := range node.Children {
		if i == len(node.Children)-1 {
			printLogicalTree(child, level+1, prefixLast)
		} else {
			printLogicalTree(child, level+1, prefix)
		}
	}
}
