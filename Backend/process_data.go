package main

import (
	"fmt"
	"strconv"
)

type processedRecidData struct {
	ID                                             string `json:"id"`
	Gender                                         string `json:"gender"`
	Race                                           string `json:"race"`
	AgeAtRelease                                   string `json:"age_at_release"`
	ResidencePUMA                                  string `json:"residence_puma"`
	GangAffiliated                                 string `json:"gang_affiliated"`
	SupervisionRiskScoreFirst                      string `json:"supervision_risk_score_first"`
	SupervisionLevelFirst                          string `json:"supervision_level_first"`
	EducationLevel                                 string `json:"education_level"`
	Dependents                                     string `json:"dependents"`
	PrisonOffense                                  string `json:"prison_offense"`
	PrisonYears                                    string `json:"prison_years"`
	PriorArrestEpisodesFelony                      string `json:"prior_arrest_episodes_felony"`
	PriorArrestEpisodesMisd                        string `json:"prior_arrest_episodes_misd"`
	PriorArrestEpisodesViolent                     string `json:"prior_arrest_episodes_violent"`
	PriorArrestEpisodesProperty                    string `json:"prior_arrest_episodes_property"`
	PriorArrestEpisodesDrug                        string `json:"prior_arrest_episodes_drug"`
	PriorArrestEpisodesPPViolationCharges          string `json:"prior_arrest_episodes"`
	PriorArrestEpisodesDVCharges                   bool   `json:"prior_arrest_episodes_1"`
	PriorArrestEpisodesGunCharges                  bool   `json:"prior_arrest_episodes_2"`
	PriorConvictionEpisodesFelony                  string `json:"prior_conviction_episodes"`
	PriorConvictionEpisodesMisd                    string `json:"prior_conviction_episodes_1"`
	PriorConvictionEpisodesViol                    bool   `json:"prior_conviction_episodes_2"`
	PriorConvictionEpisodesProp                    string `json:"prior_conviction_episodes_3"`
	PriorConvictionEpisodesDrug                    string `json:"prior_conviction_episodes_4"`
	PriorConvictionEpisodesPPViolationCharges      bool   `json:"prior_conviction_episodes_5"`
	PriorConvictionEpisodesDomesticViolenceCharges bool   `json:"prior_conviction_episodes_6"`
	PriorConvictionEpisodesGunCharges              bool   `json:"prior_conviction_episodes_7"`
	PriorRevocationsParole                         bool   `json:"prior_revocations_parole"`
	PriorRevocationsProbation                      bool   `json:"prior_revocations_probation"`
	ConditionMHSA                                  bool   `json:"condition_mh_sa"`
	ConditionCogEd                                 bool   `json:"condition_cog_ed"`
	ConditionOther                                 bool   `json:"condition_other"`
	ViolationsElectronicMonitoring                 bool   `json:"violations"`
	ViolationsInstruction                          bool   `json:"violations_instruction"`
	ViolationsFailToReport                         bool   `json:"violations_failtoreport"`
	ViolationsMoveWithoutPermission                bool   `json:"violations_1"`
	DelinquencyReports                             string `json:"delinquency_reports"`
	ProgramAttendances                             string `json:"program_attendances"`
	ProgramUnexcusedAbsences                       string `json:"program_unexcusedabsences"`
	ResidenceChanges                               string `json:"residence_changes"`
	AvgDaysPerDrugTest                             string `json:"avg_days_per_drugtest"`
	DrugTestsTHCPositive                           string `json:"drugtests_thc_positive"`
	DrugTestsCocainePositive                       string `json:"drugtests_cocaine_positive"`
	DrugTestsMethPositive                          string `json:"drugtests_meth_positive"`
	DrugTestsOtherPositive                         string `json:"drugtests_other_positive"`
	PercentDaysEmployed                            string `json:"percent_days_employed"`
	JobsPerYear                                    string `json:"jobs_per_year"`
	EmploymentExempt                               bool   `json:"employment_exempt"`
	RecidivismWithin3years                         bool   `json:"recidivism_within_3years"`
	TrainingSample                                 string `json:"training_sample"`
}

// impute function to handle missing values
func impute(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func cleanRecid(data []recidData) ([]processedRecidData, error) {
	var processedData []processedRecidData
	// Impute missing values
	for _, entry := range data {
		processed := processedRecidData{
			ID:                                        entry.ID,
			Gender:                                    entry.Gender,
			Race:                                      entry.Race,
			AgeAtRelease:                              entry.AgeAtRelease,
			ResidencePUMA:                             entry.ResidencePUMA,
			GangAffiliated:                            impute(strconv.FormatBool(entry.GangAffiliated), "Unknown"),
			SupervisionRiskScoreFirst:                 impute(entry.SupervisionRiskScoreFirst, "-1"),
			SupervisionLevelFirst:                     entry.SupervisionLevelFirst,
			EducationLevel:                            entry.EducationLevel,
			Dependents:                                entry.Dependents,
			PrisonOffense:                             entry.PrisonOffense,
			PrisonYears:                               entry.PrisonYears,
			PriorArrestEpisodesFelony:                 entry.PriorArrestEpisodesFelony,
			PriorArrestEpisodesMisd:                   entry.PriorArrestEpisodesMisd,
			PriorArrestEpisodesViolent:                entry.PriorArrestEpisodesViolent,
			PriorArrestEpisodesProperty:               entry.PriorArrestEpisodesProperty,
			PriorArrestEpisodesDrug:                   entry.PriorArrestEpisodesDrug,
			PriorArrestEpisodesPPViolationCharges:     entry.PriorArrestEpisodesPPViolationCharges,
			PriorArrestEpisodesDVCharges:              entry.PriorArrestEpisodesDVCharges,
			PriorArrestEpisodesGunCharges:             entry.PriorArrestEpisodesGunCharges,
			PriorConvictionEpisodesFelony:             entry.PriorConvictionEpisodesFelony,
			PriorConvictionEpisodesMisd:               entry.PriorConvictionEpisodesMisd,
			PriorConvictionEpisodesViol:               entry.PriorConvictionEpisodesViol,
			PriorConvictionEpisodesProp:               entry.PriorConvictionEpisodesProp,
			PriorConvictionEpisodesDrug:               entry.PriorConvictionEpisodesDrug,
			PriorConvictionEpisodesPPViolationCharges: entry.PriorConvictionEpisodesPPViolationCharges,
			PriorConvictionEpisodesDomesticViolenceCharges: entry.PriorConvictionEpisodesDomesticViolenceCharges,
			PriorConvictionEpisodesGunCharges:              entry.PriorConvictionEpisodesGunCharges,
			PriorRevocationsParole:                         entry.PriorRevocationsParole,
			PriorRevocationsProbation:                      entry.PriorRevocationsProbation,
			ConditionMHSA:                                  entry.ConditionMHSA,
			ConditionCogEd:                                 entry.ConditionCogEd,
			ConditionOther:                                 entry.ConditionOther,
			ViolationsElectronicMonitoring:                 entry.ViolationsElectronicMonitoring,
			ViolationsInstruction:                          entry.ViolationsInstruction,
			ViolationsFailToReport:                         entry.ViolationsFailToReport,
			ViolationsMoveWithoutPermission:                entry.ViolationsMoveWithoutPermission,
			DelinquencyReports:                             entry.DelinquencyReports,
			ProgramAttendances:                             entry.ProgramAttendances,
			ProgramUnexcusedAbsences:                       entry.ProgramUnexcusedAbsences,
			ResidenceChanges:                               entry.ResidenceChanges,
			AvgDaysPerDrugTest:                             impute(entry.AvgDaysPerDrugTest, "-1"),
			DrugTestsTHCPositive:                           impute(entry.DrugTestsTHCPositive, "-1"),
			DrugTestsCocainePositive:                       impute(entry.DrugTestsCocainePositive, "-1"),
			DrugTestsMethPositive:                          impute(entry.DrugTestsMethPositive, "-1"),
			DrugTestsOtherPositive:                         impute(entry.DrugTestsOtherPositive, "-1"),
			PercentDaysEmployed:                            impute(entry.PercentDaysEmployed, "-1"),
			JobsPerYear:                                    impute(entry.JobsPerYear, "-1"),
			EmploymentExempt:                               entry.EmploymentExempt,
			RecidivismWithin3years:                         entry.RecidivismWithin3years,
			TrainingSample:                                 entry.TrainingSample,
		}
		processedData = append(processedData, processed)
	}

	return processedData, nil
}

func splitTrainTest(allRecid []processedRecidData) ([]processedRecidData, []processedRecidData, error) {

	var trainSet, testSet []processedRecidData
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
