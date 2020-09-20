package fileoperation

import (
	"regexp"
	"strings"
)
//ReplaceFileName ...Eg: goldenfish.jpg --> 1.jpg( here 1 is the item_id in database)
func ReplaceFileName(actualFname,suggestedFname string)(string){
	re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
	match1 := re.FindStringSubmatch(actualFname)
	return strings.Replace(actualFname,match1[2],suggestedFname, 1)
}

