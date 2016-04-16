package pods

import (
    "os"
    "bufio"
    "path"
    "strings"
    "io/ioutil"
    "errors"
)

var (
    ErrNoCgroupSubsystem = errors.New("rkt-inspect: Subsystem cgroup not found")
)


func GetPodPpid(root string, uuid string) (string, error) {
    ppidPath := path.Join(getPodDir(root, uuid, ""), "ppid")
    ppidRaw,err := ioutil.ReadFile(ppidPath)
    if err != nil {
        return "", err
    }

    return strings.TrimSpace(string(ppidRaw)), nil
}

func GetPodCgroup(root string, uuid string, subsystem string) (string, error) {
    ppid,err := GetPodPpid(root, uuid)
    if err != nil {
        return "", err
    }

    fd,err2 := os.Open("/proc/"+ppid+"/cgroup")
    if err2 != nil {
        return "", err2
    }
    defer fd.Close()

    cgroups := make(map[string]string)

    scanner := bufio.NewScanner(fd)
    for scanner.Scan() {
        parts := strings.Split(scanner.Text(), ":")
        subsystems := strings.Split(parts[1], ",")
        for i := range subsystems {
            cgroups[subsystems[i]] = parts[2]
        }
    }

    if cgroup,found := cgroups[subsystem]; found {
        return cgroup,nil
    } else {
        return "",ErrNoCgroupSubsystem
    }
}
