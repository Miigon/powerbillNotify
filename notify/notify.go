package notify

import (
	"fmt"
	"strings"

	"github.com/miigon/powerbillNotify/bill"
	"github.com/miigon/powerbillNotify/conf"
)

func TryNotify(bills []bill.PowerBill) {
	// assuming the bills are chronologically continuous
	var pastUsages []float64
	for i, v := range bills {
		if i == 0 {
			continue
		}
		pastUsages = append(pastUsages, v.TotalUsage-bills[i-1].TotalUsage)
	}
	fmt.Printf("pastUsages: %v\n", pastUsages)

	var nDayAverage, nMinusOneDayAverage float64
	var sum float64 = 0.0
	for i, v := range pastUsages {
		sum += v
		if i == len(pastUsages)-2 {
			nMinusOneDayAverage = sum / float64(len(pastUsages)-1)
		}
	}
	nDayAverage = sum / float64(len(pastUsages))

	fmt.Printf("nDayAverage: %v\n", nDayAverage)
	fmt.Printf("nMinusOneDayAverage: %v\n", nMinusOneDayAverage)

	var templateValues = make(map[string]string)
	templateValues["BUILDING"] = conf.Config.Building
	templateValues["ROOM"] = conf.Config.RoomName
	templateValues["NDAYSAVERAGE"] = fmt.Sprintf("%.2f", nDayAverage)
	templateValues["N_MINUS_1_DAYSAVERAGE"] = fmt.Sprintf("%.2f", nMinusOneDayAverage)
	templateValues["PASTNDAYS"] = fmt.Sprintf("%d", conf.Config.Conditions.NDaysAveragePastUsage)
	templateValues["PASTNDAYS_MINUS_1"] = fmt.Sprintf("%d", conf.Config.Conditions.NDaysAveragePastUsage-1)
	powerYesterday := pastUsages[len(pastUsages)-1]
	powerLeft := bills[len(bills)-1].PowerLeft
	daysLeft := int(powerLeft / nDayAverage)
	templateValues["POWERYESTERDAY"] = fmt.Sprintf("%.2f", powerYesterday)
	templateValues["POWERLEFT"] = fmt.Sprintf("%.2f", powerLeft)
	templateValues["DAYSLEFT"] = fmt.Sprintf("%d", daysLeft)

	fmt.Printf("powerLeft: %.2f\n", powerLeft)
	fmt.Printf("daysLeft: %v\n", daysLeft)
	if conf.Config.Conditions.NotifyNearRunningOut {
		fmt.Printf("conf.Config.Conditions.DaysLeftThreshold: %v\n", conf.Config.Conditions.DaysLeftThreshold)
		if daysLeft <= conf.Config.Conditions.DaysLeftThreshold {
			SendNotification(
				fillTemplate(conf.Config.Templates.NearRunningOut.Title, templateValues),
				fillTemplate(conf.Config.Templates.NearRunningOut.Content, templateValues))
		}
	}

	fmt.Printf("powerYesterday: %v\n", powerYesterday)
	if conf.Config.Conditions.NotifyAbnormallyHigh {
		if powerYesterday >= nMinusOneDayAverage*conf.Config.Conditions.HighUsageThreshold {
			SendNotification(
				fillTemplate(conf.Config.Templates.AbnormallyHigh.Title, templateValues),
				fillTemplate(conf.Config.Templates.AbnormallyHigh.Content, templateValues))
		}
	}
}

func fillTemplate(template string, values map[string]string) string {
	result := template
	for k, v := range values {
		result = strings.ReplaceAll(result, "${"+k+"}", v)
	}
	return result
}
