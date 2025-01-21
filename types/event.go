package types

// Favori struct to hold the favori data
type Favori struct {
	F1 int    `json:"f1"`
	F2 string `json:"f2"`
	F3 string `json:"f3"`
	F4 string `json:"f4"`
	F5 string `json:"f5"`
}

// Event struct to hold individual event data
type Event struct {
	ID                                         *int64  `json:"id"`
	DateDebut                                  string  `json:"date_debut"`
	DateFin                                    string  `json:"date_fin"`
	Duree                                      string  `json:"duree"`
	DateDebutMultijours                        *string `json:"date_debut_multijours"`
	DateFinMultijours                          *string `json:"date_fin_multijours"`
	Matiere                                    *string `json:"matiere"`
	TypeActivite                               *string `json:"type_activite"`
	ValidationIntervenant                      *string `json:"validation_intervenant"`
	Ressource                                  *string `json:"ressource"`
	StatutIntervention                         *string `json:"statut_intervention"`
	Intervenants                               string  `json:"intervenants"`
	IsBreak                                    bool    `json:"is_break"`
	IsEmpty                                    bool    `json:"is_empty"`
	Description                                *string `json:"description"`
	Favori                                     *Favori `json:"favori"` // added the nested Favori struct
	EstInterventionPlanningApprenant           bool    `json:"est_intervention_planning_apprenant"`
	EstInterventionPlanningIntervenant         bool    `json:"est_intervention_planning_intervenant"`
	EstDerniereInterventionPlanningApprenant   bool    `json:"est_derniere_intervention_planning_apprenant"`
	EstDerniereInterventionPlanningIntervenant bool    `json:"est_derniere_intervention_planning_intervenant"`
	EstDerniereInterventionPlanningAppInt      bool    `json:"est_derniere_intervention_planning_app_int"`
}

// {
//     "est_intervention_planning_apprenant": true,
//     "est_intervention_planning_intervenant": false,
//     "id": 19156166,
//     "date_debut": "2025-02-28T13:30:00.000",
//     "date_fin": "2025-02-28T17:45:00.000",
//     "duree": "4:15",
//     "date_debut_multijours": null,
//     "date_fin_multijours": null,
//     "matiere": null,
//     "type_activite": null,
//     "validation_intervenant": null,
//     "ressource": null,
//     "statut_intervention": "",
//     "intervenants": "LANNEL",
//     "is_break": false,
//     "is_empty": false,
//     "description": null,
//     "favori": {
//         "f1": 19156166,
//         "f2": " | ",
//         "f3": "Droit ",
//         "f4": "LANNEL",
//         "f5": "Cours FHES  "
//     },
//     "est_derniere_intervention_planning_apprenant": false,
//     "est_derniere_intervention_planning_intervenant": false,
//     "est_derniere_intervention_planning_app_int": true
// }
