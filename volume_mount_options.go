package volume_mount_options

import (
	"fmt"
	"strconv"
	"strings"
)

type MountOpts map[string]string

func NewMountOpts(userOpts map[string]interface{}, mask MountOptsMask) (MountOpts, error) {
	mountOpts := mask.Defaults
	errorList := []string{}

	for k, v := range userOpts {
		if inArray(mask.Ignored, k) {
			continue
		}

		if inArray(mask.Allowed, k) {
			uv := uniformKeyData(k, v)
			mountOpts[k] = uv
		} else if !mask.SloppyMount {
			errorList = append(errorList, k)
		}
	}

	if len(errorList) > 0 {
		return MountOpts{}, fmt.Errorf("Not allowed options: %s", strings.Join(errorList, ", "))
	}

	for _, k := range mask.Mandatory {
		if _, ok := userOpts[k]; !ok {
			errorList = append(errorList, k)
		}
	}

	if len(errorList) > 0 {
		return MountOpts{}, fmt.Errorf("Missing mandatory options: %s", strings.Join(errorList, ", "))
	}

	return mountOpts, nil
}

func inArray(list []string, key string) bool {
	for _, k := range list {
		if k == key {
			return true
		}
	}

	return false
}

func uniformKeyData(key string, data interface{}) string {
	switch key {
	case "auto-traverse-mounts":
		return uniformData(data, true)

	case "dircache":
		return uniformData(data, true)

	}

	return uniformData(data, false)
}

func uniformData(data interface{}, boolAsInt bool) string {

	switch data.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return fmt.Sprintf("%#v", data)

	case string:
		return data.(string)

	case bool:
		if boolAsInt {
			if data.(bool) {
				return "1"
			} else {
				return "0"
			}
		} else {
			return strconv.FormatBool(data.(bool))
		}
	}

	return ""
}