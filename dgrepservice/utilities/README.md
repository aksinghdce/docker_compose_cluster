# Documentation for utilities package

```
PS C:\go_work\src\docker_compose_cluster> docker exec dockercomposecluster_grepservice1_1 go doc ./utilities
package utilities // import "app/utilities"

func Decorate(f http.HandlerFunc) http.HandlerFunc
func LocalGrep(arguments []string) string
func Log(ctx context.Context, msg ...string)
func ReadConfig(path string) *list.List
func RemoteGrep(machine string, cmd url.Values) <-chan string
type Cluster struct{ ... }
```
