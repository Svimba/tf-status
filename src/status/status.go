package status

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// TFStatus struct
type TFStatus struct {
	PodName   string     `json:"podName"`
	Groups    []*TFGroup `json:"groups"`
	PlainText string     `json:"-"`
}

// TFGroup stuct to handling TF service groups e.g. Control, Config, Analytics
type TFGroup struct {
	Name     string     `json:"name"`
	Services []*Service `json:"services"`
}

// Service struct
type Service struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (tfg *TFGroup) isDefined() bool {
	if len(tfg.Name) > 0 {
		return true
	}
	return false
}

// PrintGroups struct
func (tfs *TFStatus) PrintGroups() {
	for _, grp := range tfs.Groups {
		fmt.Printf(" --- %s ---\n", grp.Name)
		for _, svc := range grp.Services {
			fmt.Printf("\t %s : %s\n", svc.Name, svc.Status)
		}
	}
}

// GetContrailStatus func
func (tfs *TFStatus) GetContrailStatus() {
	cmd := exec.Command("python", "/root/contrail-status.py")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	tfs.PlainText = string(out)
}

// ParseToJSON func
func (tfs *TFStatus) ParseToJSON() {
	currentGroup := &TFGroup{}

	lines := strings.Split(tfs.PlainText, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "==") {
			group := getGroupName(line)
			if group != currentGroup.Name {
				if currentGroup.isDefined() {
					tfs.Groups = append(tfs.Groups, currentGroup)
					currentGroup = new(TFGroup)
				}
				currentGroup.Name = group
			}
			continue
		}
		if len(currentGroup.Name) == 0 {
			continue
		}
		srv := newService(line)
		if srv == nil {
			continue
		}
		currentGroup.Services = append(currentGroup.Services, srv)
	}
}

func getGroupName(group string) string {
	var re = regexp.MustCompile(`(^|[={2} \n])([ ={2}]|$)`)
	return re.ReplaceAllString(group, ``)
}

func newService(service string) *Service {
	s := strings.Split(service, ": ")

	if len(s) == 2 {
		re := regexp.MustCompile(`(^|\s{2,})($|\s)`)
		srv := re.ReplaceAllString(s[0], ``)
		stat := re.ReplaceAllString(s[1], ``)
		return &Service{Name: srv, Status: stat}
	}

	return nil
}
