package utility

import (
	"fmt"
	"math"
	"runtime"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
)

// CheckOS is a function to check the OS of the user
func CheckOS() string {
	os := runtime.GOOS
	var returnValue string
	switch os {
	case "windows":
		returnValue = "windows"
	case "darwin":
		returnValue = "darwin"
	case "linux":
		returnValue = "linux"
	default:
		fmt.Printf("%s.\n", os)
	}

	return returnValue
}

// CheckQuotaPercent function to check the percent of the quota
func CheckQuotaPercent(limit int, usage int) string {

	var returnText string

	calculation := float64(usage) / float64(limit) * 100
	percent := math.Round(calculation)

	switch {
	case percent >= 80 && percent < 100:
		returnText = Orange(fmt.Sprintf("%d/%d", usage, limit))
	case percent == 100:
		returnText = Red(fmt.Sprintf("%d/%d", usage, limit))
	default:
		returnText = Green(fmt.Sprintf("%d/%d", usage, limit))
	}

	return returnText
}

// CheckAvailability is a function to check if the user can
// create Iaas and k8s cluster base on the result of region
func CheckAvailability(resource string, regionSet string) (bool, string, error) {
	var defaultRegion *civogo.Region
	client, err := config.CivoAPIClient()
	if err != nil {
		return false, "", err
	}

	if regionSet != "" {
		client.Region = regionSet
	}

	switch {
	case config.Current.Meta.DefaultRegion == "" && regionSet != "" || config.Current.Meta.DefaultRegion != "" && regionSet != "":
		defaultRegion, err = client.FindRegion(regionSet)
		if err != nil {
			return false, "", err
		}
	case config.Current.Meta.DefaultRegion != "" && regionSet == "":
		defaultRegion, err = client.FindRegion(config.Current.Meta.DefaultRegion)
		if err != nil {
			return false, "", err
		}
	default:
		defaultRegion, err = client.GetDefaultRegion()
		if err != nil {
			return false, "", err
		}
	}

	if resource == "kubernetes" {
		if defaultRegion.Features.Kubernetes && defaultRegion.OutOfCapacity == false {
			return true, "", nil
		}
	}
	if resource == "instance" {
		if defaultRegion.Features.Iaas && defaultRegion.OutOfCapacity == false {
			return true, "", nil
		}
	}

	return false, defaultRegion.Code, nil
}
