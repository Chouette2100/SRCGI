package ShowroomCGIlib
import (
	"testing"
	"log"
)
func Test_GetWeightedCnt(t *testing.T) {
	pcl := make(Pclist, 0)
	pcl = append(pcl, P_c{50, 4, 0, 0})
	pcl = append(pcl, P_c{650, 3, 0, 0})
	pcl = append(pcl, P_c{1250, 1, 0, 0})
	pcl = append(pcl, P_c{1250, 2, 0, 0})
	pcl = append(pcl, P_c{1250, 3, 0, 0})
	pcl = append(pcl, P_c{1250, 4, 0, 0})
	pcl = append(pcl, P_c{1250, 5, 0, 0})
	pcl = append(pcl, P_c{1780, 1, 0, 0})
	pcl = append(pcl, P_c{1780, 2, 0, 0})
	pcl = append(pcl, P_c{1780, 3, 0, 0})
	pcl = append(pcl, P_c{1780, 4, 0, 0})
	pcl = append(pcl, P_c{1780, 5, 0, 0})
	nfr := 1
	GetWeightedCnt(pcl, nfr)
	//	t.Log(pcl)
	for p := range(pcl) {
		log.Println(pcl[p])
	}
}