package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type recidData struct {
	ID                                             string  `json:"id"`
	Gender                                         string  `json:"gender"`
	Race                                           string  `json:"race"`
	AgeAtRelease                                   string  `json:"age_at_release"`
	ResidencePUMA                                  string  `json:"residence_puma"`
	GangAffiliated                                 bool    `json:"gang_affiliated"`
	SupervisionRiskScoreFirst                      string  `json:"supervision_risk_score_first"`
	SupervisionLevelFirst                          string  `json:"supervision_level_first"`
	EducationLevel                                 string  `json:"education_level"`
	Dependents                                     string  `json:"dependents"`
	PrisonOffense                                  string  `json:"prison_offense"`
	PrisonYears                                    string  `json:"prison_years"`
	PriorArrestEpisodesFelony                      string  `json:"prior_arrest_episodes_felony"`
	PriorArrestEpisodesMisd                        string  `json:"prior_arrest_episodes_misd"`
	PriorArrestEpisodesViolent                     string  `json:"prior_arrest_episodes_violent"`
	PriorArrestEpisodesProperty                    string  `json:"prior_arrest_episodes_property"`
	PriorArrestEpisodesDrug                        string  `json:"prior_arrest_episodes_drug"`
	PriorArrestEpisodesPPViolationCharges          string  `json:"prior_arrest_episodes"`
	PriorArrestEpisodesDVCharges                   bool    `json:"prior_arrest_episodes_1"`
	PriorArrestEpisodesGunCharges                  bool    `json:"prior_arrest_episodes_2"`
	PriorConvictionEpisodesFelony                  string  `json:"prior_conviction_episodes"`
	PriorConvictionEpisodesMisd                    string  `json:"prior_conviction_episodes_1"`
	PriorConvictionEpisodesViol                    bool    `json:"prior_conviction_episodes_2"`
	PriorConvictionEpisodesProp                    string  `json:"prior_conviction_episodes_3"`
	PriorConvictionEpisodesDrug                    string  `json:"prior_conviction_episodes_4"`
	PriorConvictionEpisodesPPViolationCharges      bool    `json:"prior_conviction_episodes_5"`
	PriorConvictionEpisodesDomesticViolenceCharges bool    `json:"prior_conviction_episodes_6"`
	PriorConvictionEpisodesGunCharges              bool    `json:"prior_conviction_episodes_7"`
	PriorRevocationsParole                         bool    `json:"prior_revocations_parole"`
	PriorRevocationsProbation                      bool    `json:"prior_revocations_probation"`
	ConditionMHSA                                  bool    `json:"condition_mh_sa"`
	ConditionCogEd                                 bool    `json:"condition_cog_ed"`
	ConditionOther                                 bool    `json:"condition_other"`
	ViolationsElectronicMonitoring                 bool    `json:"violations"`
	ViolationsInstruction                          bool    `json:"violations_instruction"`
	ViolationsFailToReport                         bool    `json:"violations_failtoreport"`
	ViolationsMoveWithoutPermission                bool    `json:"violations_1"`
	DelinquencyReports                             string  `json:"delinquency_reports"`
	ProgramAttendances                             string  `json:"program_attendances"`
	ProgramUnexcusedAbsences                       string  `json:"program_unexcusedabsences"`
	ResidenceChanges                               string  `json:"residence_changes"`
	AvgDaysperDrugTest                             float64 `json:"avg_days_per_drugtest"`
	DrugTestsTHCPositive                           float64 `json:"drugtests_thc_positive"`
	DrugTestsCocainePositive                       float64 `json:"drugtests_cocaine_positive"`
	DrugTestsMethPositive                          float64 `json:"drugtests_meth_positive"`
	DrugTestsOtherPositive                         float64 `json:"drugtests_other_positive"`
	PercentDaysEmployed                            float64 `json:"percent_days_employed"`
	JobsPerYear                                    float64 `json:"jobs_per_year"`
	EmploymentExempt                               bool    `json:"employment_exempt"`
	RecidivismWithin3years                         bool    `json:"recidivism_within_3years"`
	TrainingSample                                 string  `json:"training_sample"`
	// RecidivismArrestYear1                          bool   `json:"recidivism_arrest_year1"`
	// RecidivismArrestYear2                          bool   `json:"recidivism_arrest_year2"`
	// RecidivismArrestYear3                          bool   `json:"recidivism_arrest_year3"`
}

func buildUrls(base_url string) []string {
	// API has a limit of 1,000 rows, so need to create list of url in 1000 increment chunks
	var query_urls []string
	limit := 1000
	totaldata := 10000 //25835

	for page, completed := 0, 0; completed < totaldata; page++ {
		url := fmt.Sprintf("%s?$limit=%d&$offset=%d", base_url, limit, page*limit)
		query_urls = append(query_urls, url)
		completed += limit
	}

	return query_urls
}

func getData(url string, wg *sync.WaitGroup, ch chan<- []recidData) {
	defer wg.Done()

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	defer response.Body.Close()

	var recid []recidData
	b, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	if err := json.Unmarshal(b, &recid); err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return
	}

	ch <- recid

}

func getFullDataset() ([]recidData, error) {
	base_url := "https://data.ojp.usdoj.gov/resource/ynf5-u8nk.json"

	fmt.Println("Starting Data Pull")

	query_urls := buildUrls(base_url)

	recidChan := make(chan []recidData, 26)

	var wg sync.WaitGroup
	for _, url := range query_urls {
		wg.Add(1)
		go getData(url, &wg, recidChan)
	}

	// Close channel after all goroutines have finished
	go func() {
		wg.Wait()
		close(recidChan)
	}()

	var allRecid []recidData
	for records := range recidChan {
		allRecid = append(allRecid, records...)
	}

	fmt.Printf("Total records fetched: %d\n", len(allRecid))

	return allRecid, nil

}

func splitTrainTest(allRecid []recidData) ([]recidData, []recidData, error) {

	var trainSet, testSet []recidData
	for _, rec := range allRecid {
		if rec.TrainingSample == "1" {
			trainSet = append(trainSet, rec)
		} else {
			testSet = append(testSet, rec)
		}
	}
	fmt.Printf("Training samples: %d\n", len(trainSet))
	fmt.Printf("Non-training samples: %d\n", len(testSet))

	return trainSet, testSet, nil

}
