package main

import "testing"

func TestRunEvaluationLogic(t *testing.T) {
	app := &App{}

	tests := []struct {
		name   string
		info   *CombinedFlagInfo
		userID string
		want   bool
	}{
		{
			name: "flag desligada",
			info: &CombinedFlagInfo{
				Flag: &Flag{Name: "checkout", IsEnabled: false},
			},
			want: false,
		},
		{
			name: "flag ligada sem regra",
			info: &CombinedFlagInfo{
				Flag: &Flag{Name: "checkout", IsEnabled: true},
			},
			want: true,
		},
		{
			name: "regra desligada",
			info: &CombinedFlagInfo{
				Flag: &Flag{Name: "checkout", IsEnabled: true},
				Rule: &TargetingRule{IsEnabled: false},
			},
			want: true,
		},
		{
			name: "percentagem 100 ativa",
			info: &CombinedFlagInfo{
				Flag: &Flag{Name: "checkout", IsEnabled: true},
				Rule: &TargetingRule{
					IsEnabled: true,
					Rules:     Rule{Type: "PERCENTAGE", Value: float64(100)},
				},
			},
			want: true,
		},
		{
			name: "percentagem zero desativa",
			info: &CombinedFlagInfo{
				Flag: &Flag{Name: "checkout", IsEnabled: true},
				Rule: &TargetingRule{
					IsEnabled: true,
					Rules:     Rule{Type: "PERCENTAGE", Value: float64(0)},
				},
			},
			want: false,
		},
		{
			name:   "50 por cento usuario fora",
			userID: "user-1",
			info: &CombinedFlagInfo{
				Flag: &Flag{Name: "checkout", IsEnabled: true},
				Rule: &TargetingRule{
					IsEnabled: true,
					Rules:     Rule{Type: "PERCENTAGE", Value: float64(50)},
				},
			},
			want: false,
		},
		{
			name:   "50 por cento usuario dentro",
			userID: "user-2",
			info: &CombinedFlagInfo{
				Flag: &Flag{Name: "checkout", IsEnabled: true},
				Rule: &TargetingRule{
					IsEnabled: true,
					Rules:     Rule{Type: "PERCENTAGE", Value: float64(50)},
				},
			},
			want: true,
		},
		{
			name: "valor de percentagem invalido",
			info: &CombinedFlagInfo{
				Flag: &Flag{Name: "checkout", IsEnabled: true},
				Rule: &TargetingRule{
					IsEnabled: true,
					Rules:     Rule{Type: "PERCENTAGE", Value: "50"},
				},
			},
			want: false,
		},
		{
			name: "flag inexistente",
			info: &CombinedFlagInfo{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.runEvaluationLogic(tt.info, tt.userID)
			if got != tt.want {
				t.Errorf("resultado = %v; esperado = %v", got, tt.want)
			}
		})
	}
}

func TestGetDeterministicBucket(t *testing.T) {
	input := "user-1checkout"

	first := getDeterministicBucket(input)
	second := getDeterministicBucket(input)

	if first != second {
		t.Error("o mesmo usuario recebeu resultados diferentes")
	}

	if first < 0 || first >= 100 {
		t.Errorf("bucket fora do intervalo 0-99: %d", first)
	}
}
