package storage

import (
	"context"

	"github.com/debondarenkq/itmo-is-medsoft-7sem-y27/services/clinical/internal/models"
)

func (repository *Repository) ListPatients(ctx context.Context, page models.Page) ([]models.Patient, int, error) {
	filter := " FROM patients p WHERE ($1='' OR (p.first_name || ' ' || p.last_name || ' ' || COALESCE(p.middle_name,'')) ILIKE '%' || $1 || '%' OR EXISTS(SELECT 1 FROM records r WHERE r.patient_id=p.id AND (r.id::text ILIKE '%' || $1 || '%' OR ('ЭМК-'||lpad(r.number::text,GREATEST(6,length(r.number::text)),'0')) ILIKE '%'||$1||'%'))) AND ($2::boolean IS NULL OR EXISTS(SELECT 1 FROM records r WHERE r.patient_id=p.id)=$2)"
	var total int
	if err := repository.pool.QueryRow(ctx, "SELECT count(*)"+filter, page.Search, page.HasRecord).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := repository.pool.Query(ctx, "SELECT "+patientColumns+filter+" ORDER BY p.last_name,p.first_name,p.id LIMIT $3 OFFSET $4", page.Search, page.HasRecord, page.Limit, page.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []models.Patient{}
	for rows.Next() {
		patient, err := scanPatient(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, patient)
	}
	return items, total, rows.Err()
}
