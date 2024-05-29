package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"

	"google.golang.org/api/iterator"
)

// Define a struct to hold the prediction results
type PredictionResult struct {
	predicted_RecidivismWithin3years       string
	predicted_RecidivismWithin3years_probs []struct {
		label string
		prob  string
	}
	ID                                             string
	Gender                                         string
	Race                                           string
	AgeAtRelease                                   string
	ResidencePUMA                                  string
	GangAffiliated                                 string
	SupervisionRiskScoreFirst                      string
	SupervisionLevelFirst                          string
	EducationLevel                                 string
	Dependents                                     string
	PrisonOffense                                  string
	PrisonYears                                    string
	PriorArrestEpisodesFelony                      string
	PriorArrestEpisodesMisd                        string
	PriorArrestEpisodesViolent                     string
	PriorArrestEpisodesProperty                    string
	PriorArrestEpisodesDrug                        string
	PriorArrestEpisodesPPViolationCharges          string
	PriorArrestEpisodesDVCharges                   string
	PriorArrestEpisodesGunCharges                  string
	PriorConvictionEpisodesFelony                  string
	PriorConvictionEpisodesMisd                    string
	PriorConvictionEpisodesViol                    string
	PriorConvictionEpisodesProp                    string
	PriorConvictionEpisodesDrug                    string
	PriorConvictionEpisodesPPViolationCharges      string
	PriorConvictionEpisodesDomesticViolenceCharges string
	PriorConvictionEpisodesGunCharges              string
	PriorRevocationsParole                         string
	PriorRevocationsProbation                      string
	ConditionMHSA                                  string
	ConditionCogEd                                 string
	ConditionOther                                 string
	ViolationsElectronicMonitoring                 string
	ViolationsInstruction                          string
	ViolationsFailToReport                         string
	ViolationsMoveWithoutPermission                string
	DelinquencyReports                             string
	ProgramAttendances                             string
	ProgramUnexcusedAbsences                       string
	ResidenceChanges                               string
	AvgDaysPerDrugTest                             string
	DrugTestsTHCPositive                           string
	DrugTestsCocainePositive                       string
	DrugTestsMethPositive                          string
	DrugTestsOtherPositive                         string
	PercentDaysEmployed                            string
	JobsPerYear                                    string
	EmploymentExempt                               string
	RecidivismWithin3years                         string
}

func predictQuery(recidId string) ([]PredictionResult, error) {
	projectID := "msds434-mod7"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	qtxt := fmt.Sprintf(
		"SELECT * FROM ML.PREDICT(MODEL `recidivism.recid_xgb_model`, "+
			"(SELECT * FROM `recidivism.test` "+
			"WHERE id = %s)) ", recidId)

	if recidId == "all" {
		qtxt = fmt.Sprintf(
			"SELECT * FROM ML.PREDICT(MODEL `recidivism.recid_xgb_model`, " +
				"(SELECT * FROM `recidivism.test`))")
	}

	q := client.Query(qtxt)

	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"
	// Run the query and print results when the query job is completed.
	// job, err := q.Run(ctx)
	// if err != nil {
	// 	return err
	// }
	// status, err := job.Wait(ctx)
	// if err != nil {
	// 	return err
	// }
	// if err := status.Err(); err != nil {
	// 	return err
	// }

	var predictions []PredictionResult
	it, err := q.Read(ctx)
	for {
		//var row []bigquery.Value
		var row PredictionResult
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, row)
	}
	return predictions, nil
}
