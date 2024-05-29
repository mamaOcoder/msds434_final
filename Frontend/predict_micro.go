package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"

	"google.golang.org/api/iterator"
)

// Define a struct to hold the prediction results
type PredictionResult struct {
	PredictedRecidivismWithin3years      string `bigquery:"predicted_RecidivismWithin3years"`
	PredictedRecidivismWithin3yearsProbs []struct {
		Label string `bigquery:"label"`
		Prob  string `bigquery:"prob"`
	} `bigquery:"predicted_RecidivismWithin3years_probs"`
	ID                                             int64
	Gender                                         string
	Race                                           string
	AgeAtRelease                                   string
	ResidencePUMA                                  string
	GangAffiliated                                 string
	SupervisionRiskScoreFirst                      int64
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
	PriorArrestEpisodesDVCharges                   bool
	PriorArrestEpisodesGunCharges                  bool
	PriorConvictionEpisodesFelony                  string
	PriorConvictionEpisodesMisd                    string
	PriorConvictionEpisodesViol                    bool
	PriorConvictionEpisodesProp                    string
	PriorConvictionEpisodesDrug                    string
	PriorConvictionEpisodesPPViolationCharges      bool
	PriorConvictionEpisodesDomesticViolenceCharges bool
	PriorConvictionEpisodesGunCharges              bool
	PriorRevocationsParole                         bool
	PriorRevocationsProbation                      bool
	ConditionMHSA                                  bool
	ConditionCogEd                                 bool
	ConditionOther                                 bool
	ViolationsElectronicMonitoring                 bool
	ViolationsInstruction                          bool
	ViolationsFailToReport                         bool
	ViolationsMoveWithoutPermission                bool
	DelinquencyReports                             string
	ProgramAttendances                             string
	ProgramUnexcusedAbsences                       string
	ResidenceChanges                               string
	AvgDaysPerDrugTest                             float64
	DrugTestsTHCPositive                           float64
	DrugTestsCocainePositive                       float64
	DrugTestsMethPositive                          float64
	DrugTestsOtherPositive                         float64
	PercentDaysEmployed                            float64
	JobsPerYear                                    float64
	EmploymentExempt                               bool
	RecidivismWithin3years                         bool
}

func predictQuery(recidId string) ([]PredictionResult, error) {
	projectID := "msds434-finalproj"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	qtxt := fmt.Sprintf(
		"SELECT * FROM ML.PREDICT(MODEL `recidivism.recid_xgb_model`, "+
			"(SELECT * FROM `recidivism.test_recid` "+
			"WHERE id = %s)) ", recidId)

	if recidId == "all" {
		qtxt = fmt.Sprintf(
			"SELECT * FROM ML.PREDICT(MODEL `recidivism.recid_xgb_model`, " +
				"(SELECT * FROM `recidivism.test_recid`))")
	}

	q := client.Query(qtxt)

	// Location must match that of the dataset(s) referenced in the query.
	q.Location = "US"
	//Run the query and print results when the query job is completed.
	job, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}

	var predictions []PredictionResult
	it, err := job.Read(ctx)
	for {
		// var row []bigquery.Value
		var row PredictionResult
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		log.Printf("Fetched row: %+v\n", row)
		predictions = append(predictions, row)
	}
	log.Printf("Total predictions fetched: %d\n", len(predictions))
	return predictions, nil
}
