package fileoperation

import (
	"regexp"
	"strings"
)
//ReplaceFileName ...Eg: goldenfish.jpg --> 8.jpg( here 8 is the item_id in database)
func ReplaceFileName(actualFname,suggestedFname string)(finalName string){
	re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
	match1 := re.FindStringSubmatch(actualFname)
	finalName = strings.Replace(actualFname,match1[2],suggestedFname, 1)
	return 
}

