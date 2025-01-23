
# tmp

```
curl -X POST https://mycpe.cpe.fr/mobile/login \
-H "Accept-Encoding: gzip" \
-H "Accept-Language: en-US,en;q=0.5" \
-H "Connection: Keep-Alive" \
-H "Content-Type: application/json" \
-H "User-Agent: Dalvik/2.1.0 (Linux; U; Android 15; sdk_gphone64_x86_64 Build/AE3A.240806.005)" \
-H "Host: mycpe.cpe.fr" \
-d '{
  "login": "loan.maeght@cpe.fr",
  "password": "pass"
}'

```

reponse:
{"normal":"eyJhbGciOiJSUzI1NiJ9.eyJjcmVlIjoxNzM3NDQ2MjIzMTQ0LCJyb2xlcyI6WyJST0xFX2FwcHJlbmFudCJdLCJpZCI6MTgyOTkxODR9.Fo_QNOdAHSuxeM9S5g7JClamOkfT7w8CUxefpNZJk6hb7p5H3HFeo04GqZp-wZbghlC4ma9_MkOul-SS_Evb7k_0AigWGiA8DeRgykdL2Bxit71eDy3yEIpAFVLD0Fdlbj4D2jkSQnErpPx7UmKqRGwkOdGonPXMe7E_r6eyFh8CBgKsjHXPM79wCayExTeoqpxc-Q0DuCqWKr14G28pTU71silbZZNH9LQN8lA0YzgkHsIV55SV6PZn_-KQn5y3FIIwC5qB5TJsmgTmt78IHYDoHAAHYBpiyazUPHdqFHS955KfPG4cIPBz4AdIkbUtwuYMg9lcIZKiijR7nj6zFeZEXhJJpqrWUjL5XNSDWWXlOh7FKxelBVw6EfG_w5uylOqAjyjH3KZirPWoOe7fsGzZRkg-snk17AgMfTONV6u6VvhQlZ2KacYtqhHt3H_JSStvngFZQUAsqHMqEjZDh3dVykrgDVsskVGBR2DdeYUSAIJN9rHuloYWcF1lPLUVlyS9OMLOhZUcIiJD6CQUWA3liqjMnZ-_pX3dC2Xuw6IUZPSMa8FTrKthq0O8fU7Kg5kz_Y2o2FRo3Vjl6LBd8DDekf9M9q1TtFN9P2zzhj7XcZ2-e7tFYOCkDttME9LxiZ-9q81d5kQ0pdPrM90JlgTJ01fyMpIvoa-6sJk_0aU","comptage":"eyJhbGciOiJSUzI1NiJ9.eyJjcmVlIjoxNzM3NDQ2MjIzMjA3LCJpZCI6MTgyOTkxODQsInR5cGUiOiJDIn0.fF9AtkRwubFvx1fgAetRa3JEZGbcYjRdDkDv6fBNOyxdz3_Hiq-pI2-6IylC17RFElidVPJ_ocvSMZATcqsnOmh-asxVo9XSZecpCY4sqei7pWoR-R2uvKpFt5-9_48fM7_i1PUWd06hiIOYPEytj0-siVGZLLjvagPMsTgVWOuXTXu29Weaw-3AC9TIl36RLhGnxSHJP6WDowr_wQzdW9-w0xd4w1SnIGyKqJkxls6gdjFJMqlMaPR-podh287XHjvCVUzzRLsIeQxjap1epu-t9XwxxHYxaIClPuInFnDnTk14hhAPeHTSSfPwBkQ2NV8yfnglf8sQ0G4U8fsXg6URzLy8L1pesCj1jiObCVhrGcHQkXH_9GYCUj3BdgCdYPpWgQFxk-XKrjYKRRjrePetYlGLGtdh7taHFOimLyDqxd7HgWI3McqUgHH5isd0F9SoRpUM-V-SY4uJvu6uiAiQnGFuZzMJBo3SvIuoKCeQflxBQsQVGLYq4TPXGCnDNiu9B6kNDtKDcVxtt1ZQ8SPfOhi5R9_VyFJypMfTGebrUjxmbAJwUpiGogTZpM3ifN5OJZ484WMNrmj8WQy_2mhUV5SzlNHNWgOEpVjBbJ4izcbkcPTFWo6H6Z06RJ15X-XTXXNMfj8ezvSfphr4zc9FMpYpwE_liXsNz6b6rkE"}

```
curl -X GET "https://mycpe.cpe.fr/mobile/mon_planning?date_debut=2025-01-15&date_fin=2025-03-01" \
-H "Accept-Encoding: gzip" \
-H "Accept-Language: en-US,en;q=0.5" \
-H "Authorization: Bearer eyJhbGciOiJSUzI1NiJ9.eyJjcmVlIjoxNzM3NDgzMTkyNzM2LCJyb2xlcyI6WyJST0xFX2FwcHJlbmFudCJdLCJpZCI6MTgyOTkxODR9.FIe0OIio8kytov4dR3IBhTDGG_ZeIWlteGNj3ZUrbzLLZclvwJF1qpYwq02WKfTomeshwoj0BFQsjbov6841IshnY9AD-qURQkEFSHZQ6ZKZ-rqeMZBqSpyKdZNqBUg1fn9UvdA233VBqH1lAQI-IK7ity1qsv4sqnhf5ZiHjSkTQWqKyZ6fxJejVWpxVOzBEJ-Lrp_bTJAm_JDl3ZV_YKeFc-oaMq3urXk5qQaD9lblXgrXsXHWqJQk6Oi6HboJRX4elOArrBvGOmTy8RVuqGfznw2fi9SDjutFNvuLpTLbGp0hTirf28EzB3G8IQ0Cg0YhxQsGKELq04Rh8YhK_GZD6582mGltCjMfPS8p7pqU0ZTsb7H3VdUoCw1GrtXsqD2jgs0mVnhZy_qdM_VtxJ2UZ8Jz360sXVA1S78EsJOkfp37AVdDETxq-UA7LLCDIEvnUxTIaMo83-maiX1GyExC8C3p2nC6t7YUsGTdHThQuQCryL09JZ44m3gVIZL1zMThZ2-trhHV1Ua15Thw6SIc6_XfdMsnmgwZBsIs9pAd8KhhrhYyaW6RsalUipXUW9npkBmqBeVubh_t2emYeBnhkDUiNV2V5yYnoQYlPkOSK8Yv1RZYj6qQsDWB_lzrp_vMGwz82ufC6F22E-YllKS65OuoynpCCXCxRkzHMnM" \
-H "Connection: Keep-Alive" \
-H "Content-Type: application/json" \
-H "User-Agent: Dalvik/2.1.0 (Linux; U; Android 15; sdk_gphone64_x86_64 Build/AE3A.240806.005)" \
-H "Host: mycpe.cpe.fr"

```

reponse:
[{"est_intervention_planning_apprenant":true,"est_intervention_planning_intervenant":false,"id":null,"date_debut":"2025-02-15T08:00:00.000","date_fin":"2025-02-15T20:00:00.000","duree":"12:00","date_debut_multijours":null,"date_fin_multijours":null,"matiere":null,"type_activite":null,"validation_intervenant":null,"ressource":null,"statut_intervention":null,"intervenants":null,"is_break":false,"is_empty":true,"description":null,"favori":null,"est_derniere_intervention_planning_apprenant":true,"est_derniere_intervention_planning_intervenant":true,"est_derniere_intervention_planning_app_int":true}, 
 {"est_intervention_planning_apprenant":true,"est_intervention_planning_intervenant":false,"id":null,"date_debut":"2025-02-16T08:00:00.000","date_fin":"2025-02-16T20:00:00.000","duree":"12:00","date_debut_multijours":null,"date_fin_multijours":null,"matiere":null,"type_activite":null,"validation_intervenant":null,"ressource":null,"statut_intervention":null,"intervenants":null,"is_break":false,"is_empty":true,"description":null,"favori":null,"est_derniere_intervention_planning_apprenant":true,"est_derniere_intervention_planning_intervenant":true,"est_derniere_intervention_planning_app_int":true}, 
 {"est_intervention_planning_apprenant":true,"est_intervention_planning_intervenant":false,"id":19103545,"date_debut":"2025-02-17T08:00:00.000","date_fin":"2025-02-17T12:15:00.000","duree":"4:15","date_debut_multijours":null,"date_fin_multijours":null,"matiere":null,"type_activite":null,"validation_intervenant":null,"ressource":"","statut_intervention":"","intervenants":"COUY","is_break":false,"is_empty":false,"description":null,"favori":{"f1":19103545,"f2":" | I300","f3":"Architecture et Langages du Web","f4":"COUY","f5":"TD  "},"est_derniere_intervention_planning_apprenant":false,"est_derniere_intervention_planning_intervenant":false,"est_derniere_intervention_planning_app_int":false}, 


