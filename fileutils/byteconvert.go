package fileutils

func ByteSizeConvert(fileBytes int64, units string) float64 {

	var oneMb float64 = 1024

	var bytes float64 = float64(fileBytes)

	var kilobytes float64 = bytes / oneMb

	var megabytes float64 = kilobytes / oneMb

	var gigabytes float64 = megabytes / oneMb

	var terabytes float64 = gigabytes / oneMb

	var petabytes float64 = terabytes / oneMb

	var exabytes float64 = petabytes / oneMb

	switch units {
	case "kb":
		return bytes / oneMb

	case "mb":
		return kilobytes / oneMb

	case "gb":
		return megabytes / oneMb

	case "tb":
		return gigabytes / oneMb

	case "pb":
		return terabytes / oneMb

	case "xb":
		return petabytes / oneMb

	case "zb":
		return exabytes / oneMb

	}

	return bytes
}
