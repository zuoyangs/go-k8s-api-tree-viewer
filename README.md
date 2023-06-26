# go-k8s-api-tree-viewer

这个程序通过使用 Kubernetes 客户端库来获取 Kubernetes API 的资源组和资源信息，并使用递归方法构建和打印逻辑树状结构。

首先，程序获取了当前用户的 Kubernetes 配置文件路径，并使用该路径创建了一个 Kubernetes 客户端配置。接着，使用客户端配置创建了一个 Kubernetes 客户端。然后，定义了一个名为 Node 的结构体，表示逻辑树状结构的节点。

generateLogicalTree 函数用于生成逻辑树状结构。它首先获取 Kubernetes API 的资源组列表，然后遍历每个资源组，创建一个顶层节点。接着，获取资源组内的 API 资源列表，为每个 API 资源创建一个子节点，并将子节点添加到顶层节点的子节点列表中。

printLogicalTree 函数用于打印逻辑树状结构，它采用递归的方式遍历节点树。在打印每个节点时，使用了一个前缀字符串来表示树的层次结构。最后，程序在主函数中调用了 generateLogicalTree 函数生成逻辑树状结构，并调用 printLogicalTree 函数打印树形结构。
