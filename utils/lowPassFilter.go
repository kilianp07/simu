package utils

// Structure pour stocker l'état du filtre passe-bas
type LowPassFilter struct {
	alpha float64 // Facteur d'atténuation (0 < alpha < 1)
	yPrev float64 // Dernière valeur de sortie
}

// Fonction pour créer un nouveau filtre passe-bas avec un facteur d'atténuation donné
func NewLowPassFilter(alpha float64) *LowPassFilter {
	return &LowPassFilter{alpha, 0.0}
}

// Fonction pour mettre à jour le filtre avec une nouvelle valeur cible p_target_w
func (lpf *LowPassFilter) Update(p_target_w float64) float64 {
	lpf.yPrev = p_target_w*lpf.alpha + lpf.yPrev*(1.0-lpf.alpha)
	return lpf.yPrev
}
