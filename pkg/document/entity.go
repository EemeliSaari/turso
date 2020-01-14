package document

import "regexp"

func inspectEntity(ent string) bool {
	pattern := "&{1}#{0,1}[^;]*;{1}"

	match, err := regexp.MatchString(pattern, ent)
	if err != nil {
		panic(err)
	}
	return match
}
