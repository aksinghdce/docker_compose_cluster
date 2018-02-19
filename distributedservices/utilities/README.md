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

## Need to change the functions to accept interfaces only
```
func LocalGrep(cmd GrepCommandInterface) GrepResultInterface
```

Similarly for other functions


# Membership service functionality
![Design of a membership service node](![Design of Membership service](https://github.com/aksinghdce/docker_compose_cluster/blob/assignment2/doc/images/Overall%20design%20of%20membership%20service.png)
)